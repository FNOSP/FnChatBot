package ws

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"gorm.io/datatypes"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	TypeUserMessage        = "user_message"
	TypeThinking           = "thinking"
	TypeTaskUpdate         = "task_update"
	TypeMessage            = "message"
	TypeMessageEnd         = "message_end"
	TypePermissionRequest  = "permission_request"
	TypePermissionResponse = "permission_response"
	TypeCommandBlocked     = "command_blocked"
)

type WSMessage struct {
	Type          string         `json:"type"`
	Content       string         `json:"content,omitempty"`
	ModelID       uint           `json:"model_id,omitempty"`
	Options       map[string]any `json:"options,omitempty"`
	Delta         string         `json:"delta,omitempty"`
	Tasks         []TaskDTO      `json:"tasks,omitempty"`
	RequestID     string         `json:"request_id,omitempty"`
	RequestedPath string         `json:"requested_path,omitempty"`
	Command       string         `json:"command,omitempty"`
	BlockedPaths  []string       `json:"blocked_paths,omitempty"`
	Approved      bool           `json:"approved,omitempty"`
	Remember      bool           `json:"remember,omitempty"`
}

type ChatMessage struct {
	Role       string     `json:"role"`
	Content    string     `json:"content"`
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`
	ToolCallID string     `json:"tool_call_id,omitempty"`
	Name       string     `json:"name,omitempty"`
}

type ToolCall struct {
	Index    int      `json:"index"`
	ID       string   `json:"id,omitempty"`
	Type     string   `json:"type"`
	Function ToolFunc `json:"function"`
}

type ToolFunc struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionRequest struct {
	Model       string          `json:"model"`
	Messages    []ChatMessage   `json:"messages"`
	Stream      bool            `json:"stream"`
	Temperature float32         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Tools       []services.Tool `json:"tools,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int         `json:"index"`
		Delta        StreamDelta `json:"delta"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
}

type StreamDelta struct {
	Content   string     `json:"content"`
	ToolCalls []ToolCall `json:"tool_calls,omitempty"`
}

type TaskDTO struct {
	Name        string `json:"name"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

type PermissionResponse struct {
	Approved bool
	Remember bool
}

type PermissionManager struct {
	mu        sync.RWMutex
	pending   map[string]chan PermissionResponse
	responses map[string]PermissionResponse
}

var permissionManager = &PermissionManager{
	pending:   make(map[string]chan PermissionResponse),
	responses: make(map[string]PermissionResponse),
}

func (pm *PermissionManager) CreateRequest(requestID string) chan PermissionResponse {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	ch := make(chan PermissionResponse, 1)
	pm.pending[requestID] = ch
	return ch
}

func (pm *PermissionManager) SetResponse(requestID string, response PermissionResponse) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	if ch, exists := pm.pending[requestID]; exists {
		ch <- response
		delete(pm.pending, requestID)
	}
	pm.responses[requestID] = response
}

func (pm *PermissionManager) GetResponse(requestID string) (PermissionResponse, bool) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	resp, exists := pm.responses[requestID]
	return resp, exists
}

func (pm *PermissionManager) RemoveRequest(requestID string) {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	delete(pm.pending, requestID)
	delete(pm.responses, requestID)
}

type SandboxService struct {
	blockedPaths map[string]bool
}

func NewSandboxService() *SandboxService {
	return &SandboxService{
		blockedPaths: map[string]bool{
			"/etc/passwd":       true,
			"/etc/shadow":       true,
			"/root":             true,
			"/var/log":          true,
			"C:\\Windows":       true,
			"C:\\Program Files": true,
		},
	}
}

func (s *SandboxService) CheckCommandPermission(command string, paths []string) (allowed bool, blockedPaths []string) {
	for _, path := range paths {
		if s.blockedPaths[path] {
			blockedPaths = append(blockedPaths, path)
		}
	}
	return len(blockedPaths) == 0, blockedPaths
}

func HandleWebSocket(c *gin.Context) {
	conversationID := c.Param("id")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade websocket: %v", err)
		return
	}
	defer conn.Close()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}

		var msg WSMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Printf("Invalid message format: %v", err)
			continue
		}

		switch msg.Type {
		case TypeUserMessage:
			handleUserMessage(conn, conversationID, msg)
		case TypePermissionResponse:
			HandlePermissionResponse(conn, msg)
		}
	}
}

func handleUserMessage(conn *websocket.Conn, conversationIDStr string, msg WSMessage) {
	// Parse conversation ID
	conversationID, err := strconv.ParseUint(conversationIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid conversation ID: %v", err)
		return
	}

	// Fetch conversation and model config
	var conversation models.Conversation
	if err := db.DB.Preload("Model").Preload("Model.ProviderRef").First(&conversation, uint(conversationID)).Error; err != nil {
		log.Printf("Conversation not found: %v", err)
		sendJSON(conn, WSMessage{
			Type:    TypeMessage,
			Content: "Error: Conversation not found.",
		})
		sendJSON(conn, WSMessage{Type: TypeMessageEnd})
		return
	}

	// Check if model config exists
	if conversation.Model.ID == 0 {
		log.Printf("Model configuration not found for conversation %d", conversationID)
		sendJSON(conn, WSMessage{
			Type:    TypeMessage,
			Content: "Error: Model configuration not found. Please select a model for this conversation.",
		})
		sendJSON(conn, WSMessage{Type: TypeMessageEnd})
		return
	}

	// Save user message
	saveMessage(uint(conversationID), "user", msg.Content)

	// Fetch conversation history
	var historyMessages []models.Message
	if err := db.DB.Where("conversation_id = ?", conversationID).Order("created_at asc").Find(&historyMessages).Error; err != nil {
		log.Printf("Failed to fetch history: %v", err)
	}

	// Convert history to ChatMessage
	var messages []ChatMessage
	for _, m := range historyMessages {
		// Note: Detailed tool calls reconstruction from DB is complex if not stored properly.
		// For simplicity, we assume simple text messages here, or valid JSON in Meta if we implemented that.
		// Current DB schema stores Content as string.
		messages = append(messages, ChatMessage{
			Role:    string(m.Role),
			Content: m.Content,
		})
	}

	// Stream chat completion
	responseContent, err := streamChatCompletion(conversation.Model, messages, conn)
	if err != nil {
		log.Printf("Error streaming chat completion: %v", err)
		sendJSON(conn, WSMessage{
			Type:    TypeMessage,
			Content: fmt.Sprintf("\n\nError: %v", err),
		})
	}

	// Save assistant message
	if responseContent != "" {
		saveMessage(uint(conversationID), "assistant", responseContent)
	}

	sendJSON(conn, WSMessage{
		Type: TypeMessageEnd,
	})
}

func streamChatCompletion(config models.ModelConfig, messages []ChatMessage, conn *websocket.Conn) (string, error) {
	var provider models.Provider
	if config.ProviderRef != nil {
		provider = *config.ProviderRef
	} else if config.ProviderID != 0 {
		if err := db.DB.First(&provider, config.ProviderID).Error; err != nil {
			return "", fmt.Errorf("provider not found: %v", err)
		}
	} else {
		return "", fmt.Errorf("provider not configured for model")
	}

	if provider.BaseURL == "" {
		return "", fmt.Errorf("provider base URL is empty")
	}
	if provider.APIKey == "" {
		return "", fmt.Errorf("provider API key is empty")
	}

	baseURL := strings.TrimSuffix(provider.BaseURL, "/")
	if !strings.HasSuffix(baseURL, "/v1") {
		baseURL += "/v1"
	}
	url := fmt.Sprintf("%s/chat/completions", baseURL)

	// Get Tools
	toolService := services.NewToolService()
	tools, err := toolService.GetAvailableTools()
	if err != nil {
		log.Printf("Failed to get tools: %v", err)
		// Proceed without tools if failed
	}

	client := &http.Client{}
	var fullResponse strings.Builder

	// Maximum turns to prevent infinite loops
	maxTurns := 5
	currentTurn := 0

	for currentTurn < maxTurns {
		currentTurn++

		reqBody := ChatCompletionRequest{
			Model:       config.Model,
			Messages:    messages,
			Stream:      true,
			Temperature: config.Temperature,
			MaxTokens:   config.MaxTokens,
		}

		if len(tools) > 0 {
			reqBody.Tools = tools
		}

		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return "", fmt.Errorf("failed to marshal request: %v", err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return "", fmt.Errorf("failed to create request: %v", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.APIKey))

		resp, err := client.Do(req)
		if err != nil {
			return "", fmt.Errorf("failed to send request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(bodyBytes))
		}

		reader := bufio.NewReader(resp.Body)

		var currentToolCalls map[int]*ToolCall = make(map[int]*ToolCall)

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					break
				}
				return fullResponse.String(), fmt.Errorf("error reading stream: %v", err)
			}

			line = strings.TrimSpace(line)
			if !strings.HasPrefix(line, "data: ") {
				continue
			}

			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				break
			}

			var streamResp ChatCompletionStreamResponse
			if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
				log.Printf("Error unmarshaling stream response: %v", err)
				continue
			}

			if len(streamResp.Choices) > 0 {
				choice := streamResp.Choices[0]

				// Handle Content
				delta := choice.Delta.Content
				if delta != "" {
					fullResponse.WriteString(delta)
					sendJSON(conn, WSMessage{
						Type:  TypeMessage,
						Delta: delta,
					})
				}

				// Handle Tool Calls
				for _, tc := range choice.Delta.ToolCalls {
					if _, exists := currentToolCalls[tc.Index]; !exists {
						currentToolCalls[tc.Index] = &ToolCall{
							Index: tc.Index,
							ID:    tc.ID,
							Type:  tc.Type,
							Function: ToolFunc{
								Name:      tc.Function.Name,
								Arguments: "",
							},
						}
					}

					// Append arguments
					currentToolCalls[tc.Index].Function.Arguments += tc.Function.Arguments
					// Update name if present (usually only in first chunk)
					if tc.Function.Name != "" {
						currentToolCalls[tc.Index].Function.Name = tc.Function.Name
					}
				}
			}
		}

		// If no tool calls, we are done
		if len(currentToolCalls) == 0 {
			break
		}

		// If we have tool calls, we need to execute them and continue the conversation
		log.Printf("Tool calls detected: %d", len(currentToolCalls))

		// 1. Add assistant message with tool calls to history
		var toolCalls []ToolCall
		for i := 0; i < len(currentToolCalls); i++ {
			if tc, ok := currentToolCalls[i]; ok {
				toolCalls = append(toolCalls, *tc)
			}
		}

		messages = append(messages, ChatMessage{
			Role:      "assistant",
			Content:   "", // Tool calls usually have empty content
			ToolCalls: toolCalls,
		})

		// 2. Execute tools and add tool result messages
		for _, tc := range toolCalls {
			log.Printf("Executing tool: %s args: %s", tc.Function.Name, tc.Function.Arguments)

			// Notify frontend about tool execution
			sendJSON(conn, WSMessage{
				Type:    TypeMessage,
				Content: fmt.Sprintf("\n\n> Calling tool: %s...\n", tc.Function.Name),
			})

			// Special handling for UI updates based on tool type
			if tc.Function.Name == "TodoWrite" {
				var todoArgs struct {
					Items []TaskDTO `json:"items"`
				}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &todoArgs); err == nil {
					sendJSON(conn, WSMessage{
						Type:  TypeTaskUpdate,
						Tasks: todoArgs.Items,
					})
				}
			}

			if tc.Function.Name == "Task" {
				var taskArgs struct {
					Subagent string `json:"subagent_type"`
					Desc     string `json:"description"`
				}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &taskArgs); err == nil {
					sendJSON(conn, WSMessage{
						Type:    TypeMessage,
						Content: fmt.Sprintf("\n*Subagent [%s] started: %s*\n", taskArgs.Subagent, taskArgs.Desc),
					})
				}
			}

			if tc.Function.Name == "Skill" {
				var skillArgs struct {
					Name string `json:"name"`
				}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &skillArgs); err == nil {
					sendJSON(conn, WSMessage{
						Type:    TypeMessage,
						Content: fmt.Sprintf("\n*Loaded Skill: %s*\n", skillArgs.Name),
					})
				}
			}

			result, err := toolService.ExecuteSkill(tc.Function.Name, tc.Function.Arguments)
			if err != nil {
				result = fmt.Sprintf("Error: %v", err)
			}

			messages = append(messages, ChatMessage{
				Role:       "tool",
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
				Content:    result,
			})
		}

		// Continue loop to get next response from model
	}

	return fullResponse.String(), nil
}

func SendPermissionRequest(conn *websocket.Conn, requestID, command string, blockedPaths []string) error {
	msg := WSMessage{
		Type:         TypePermissionRequest,
		RequestID:    requestID,
		Command:      command,
		BlockedPaths: blockedPaths,
		Content:      "The command you are trying to execute requires access to restricted paths. Do you want to allow this operation?",
	}
	return conn.WriteJSON(msg)
}

func HandlePermissionResponse(conn *websocket.Conn, msg WSMessage) {
	if msg.RequestID == "" {
		log.Printf("Permission response missing request_id")
		return
	}

	response := PermissionResponse{
		Approved: msg.Approved,
		Remember: msg.Remember,
	}

	permissionManager.SetResponse(msg.RequestID, response)

	log.Printf("Received permission response for request %s: approved=%v, remember=%v",
		msg.RequestID, msg.Approved, msg.Remember)
}

func generateRequestID() string {
	return "req_" + time.Now().Format("20060102150405") + "_" + randomString(8)
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
		time.Sleep(1 * time.Nanosecond)
	}
	return string(b)
}

func sendJSON(conn *websocket.Conn, v interface{}) {
	conn.WriteJSON(v)
}

func saveMessage(conversationID uint, role string, content string) {
	msg := models.Message{
		ConversationID: conversationID,
		Role:           models.MessageRole(role),
		Content:        content,
		Meta:           datatypes.JSON{},
	}
	db.DB.Create(&msg)
}
