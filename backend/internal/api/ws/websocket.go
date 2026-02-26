package ws

import (
	"context"
	"encoding/json"
	"fmt"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"
	"fnchatbot/internal/services/llm"
	"fnchatbot/internal/services/memory"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tmc/langchaingo/llms"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const (
	TypeUserMessage        = "user_message"
	TypeTaskUpdate         = "task_update"
	TypeMessage            = "message"
	TypeMessageEnd         = "message_end"
	TypePermissionResponse = "permission_response"
	TypeImage              = "image"
)

type WSMessage struct {
	Type          string         `json:"type"`
	Content       string         `json:"content,omitempty"`
	Images        []ImagePayload `json:"images,omitempty"`
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

type ImagePayload struct {
	Data string `json:"data"`
	Type string `json:"type"`
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

func HandleWebSocket(c *gin.Context) {
	sessionID := c.Param("id")
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
			handleUserMessage(conn, sessionID, msg)
		case TypeImage:
			handleUserMessage(conn, sessionID, msg)
		case TypePermissionResponse:
			HandlePermissionResponse(msg)
		}
	}
}

func handleUserMessage(conn *websocket.Conn, sessionIDStr string, msg WSMessage) {
	sessionID, err := strconv.ParseUint(sessionIDStr, 10, 32)
	if err != nil {
		log.Printf("Invalid session ID: %v", err)
		return
	}

	var session models.Session
	if err := db.DB.Preload("Model").Preload("Model.ProviderRef").First(&session, uint(sessionID)).Error; err != nil {
		log.Printf("Session not found: %v", err)
		if err := sendJSON(conn, WSMessage{Type: TypeMessage, Content: "Error: Session not found."}); err != nil {
			log.Printf("Failed to send error message: %v", err)
		}
		if err := sendJSON(conn, WSMessage{Type: TypeMessageEnd}); err != nil {
			log.Printf("Failed to send message end: %v", err)
		}
		return
	}

	if session.Model.ID == 0 {
		log.Printf("Model configuration not found for session %d", sessionID)
		if err := sendJSON(conn, WSMessage{Type: TypeMessage, Content: "Error: Model configuration not found."}); err != nil {
			log.Printf("Failed to send error message: %v", err)
		}
		if err := sendJSON(conn, WSMessage{Type: TypeMessageEnd}); err != nil {
			log.Printf("Failed to send message end: %v", err)
		}
		return
	}

	// Determine provider
	var provider models.Provider
	if session.Model.ProviderRef != nil {
		provider = *session.Model.ProviderRef
	} else if session.Model.ProviderID != 0 {
		if err := db.DB.First(&provider, session.Model.ProviderID).Error; err != nil {
			log.Printf("Provider not found: %v", err)
			if err := sendJSON(conn, WSMessage{Type: TypeMessage, Content: "Error: Provider not found."}); err != nil {
				log.Printf("Failed to send error message: %v", err)
			}
			if err := sendJSON(conn, WSMessage{Type: TypeMessageEnd}); err != nil {
				log.Printf("Failed to send message end: %v", err)
			}
			return
		}
	} else {
		if err := sendJSON(conn, WSMessage{Type: TypeMessage, Content: "Error: Provider not configured."}); err != nil {
			log.Printf("Failed to send error message: %v", err)
		}
		if err := sendJSON(conn, WSMessage{Type: TypeMessageEnd}); err != nil {
			log.Printf("Failed to send message end: %v", err)
		}
		return
	}

	ctx := context.Background()
	llmService := llm.NewService(db.DB)

	// Save User Message
	if err := llmService.SaveUserMessage(ctx, uint(sessionID), msg.Content); err != nil {
		log.Printf("Failed to save user message: %v", err)
	}

	// Prepare Tools
	toolService := services.NewToolService()
	svcTools, _ := toolService.GetAvailableTools()
	lcTools := convertToLangChainTools(svcTools)

	// Loop for Multi-turn (Tool Execution)
	maxTurns := 5
	currentTurn := 0

	for currentTurn < maxTurns {
		currentTurn++

		// Load History (including just saved user message or previous tool outputs)
		history, err := llmService.GetHistory(ctx, uint(sessionID))
		if err != nil {
			log.Printf("Failed to load history: %v", err)
			break
		}

		// Convert History to []llms.MessageContent
		var contentMessages []llms.MessageContent

		// Find the index of the last human message to attach images to
		lastHumanIdx := -1
		for j := len(history) - 1; j >= 0; j-- {
			if history[j].GetType() == llms.ChatMessageTypeHuman {
				lastHumanIdx = j
				break
			}
		}

		for i, m := range history {
			parts := []llms.ContentPart{}

			// Handle Tool Calls
			if aiMsg, ok := m.(llms.AIChatMessage); ok && len(aiMsg.ToolCalls) > 0 {
				// For AI message with tool calls, we usually don't have text content in LangChain's view if it's purely a tool call,
				// but sometimes it has both.
				if aiMsg.Content != "" {
					parts = append(parts, llms.TextPart(aiMsg.Content))
				}
				for _, tc := range aiMsg.ToolCalls {
					parts = append(parts, llms.ToolCall{
						ID:           tc.ID,
						Type:         tc.Type,
						FunctionCall: tc.FunctionCall,
					})
				}
			} else if toolMsg, ok := m.(llms.ToolChatMessage); ok {
				// For Tool output
				parts = append(parts, llms.ToolCallResponse{
					ToolCallID: toolMsg.ID,
					Content:    toolMsg.Content,
					Name:       "", // Name is not stored in ToolChatMessage
				})
			} else {
				// Normal text message or Human message with images
				if i == lastHumanIdx && len(msg.Images) > 0 && m.GetType() == llms.ChatMessageTypeHuman {
					content := m.GetContent()
					if content != "" {
						parts = append(parts, llms.TextPart(content))
					}
					for _, img := range msg.Images {
						mimeType := img.Type
						if mimeType == "" {
							mimeType = "image/png"
						}
						url := fmt.Sprintf("data:%s;base64,%s", mimeType, img.Data)
						parts = append(parts, llms.ImageURLPart(url))
					}
				} else {
					parts = append(parts, llms.TextPart(m.GetContent()))
				}
			}

			contentMessages = append(contentMessages, llms.MessageContent{
				Role:  m.GetType(),
				Parts: parts,
			})
		}

		// Stream Chat
		resp, err := llmService.StreamChat(ctx, provider, session.Model.Model, contentMessages, lcTools, func(ctx context.Context, chunk []byte) error {
			if err := sendJSON(conn, WSMessage{
				Type:  TypeMessage,
				Delta: string(chunk),
			}); err != nil {
				log.Printf("Failed to send chunk: %v", err)
			}
			return nil
		})

		if err != nil {
			log.Printf("StreamChat error: %v", err)
			if err := sendJSON(conn, WSMessage{Type: TypeMessage, Content: fmt.Sprintf("\nError: %v", err)}); err != nil {
				log.Printf("Failed to send stream error: %v", err)
			}
			break
		}

		if len(resp.Choices) == 0 {
			break
		}

		choice := resp.Choices[0]

		// If there are tool calls
		if len(choice.ToolCalls) > 0 {
			// Save AI Message with Tool Calls
			// We need to manually construct AIChatMessage with ToolCalls and save it
			// Note: We might need to handle partial text content too if any
			aiMsg := llms.AIChatMessage{
				Content:   choice.Content,
				ToolCalls: choice.ToolCalls,
			}
			// Use internal memory helper to save complex message (since SaveAIMessage only takes string)
			// We need to expose a SaveMessage method or use memory package directly.
			// Since llmService exposes GetHistory which returns memory.SQLiteHistory, let's use memory package directly here or add helper.
			// Ideally llmService should have `SaveMessage`.
			// Let's assume we can access `memory` package since we imported it (but we didn't import it in this file yet, let's check imports).
			// We imported `fnchatbot/internal/services/memory`? No, we need to.
			// Or we can add `SaveMessage` to `llm.Service`.
			// I'll add `SaveMessage` to `llm.Service` via a separate edit or just rely on `SaveAIMessage` for text and fail for tools?
			// No, tools are critical.

			// Let's hack it: `llmService.SaveAIMessage` only saves text.
			// I need to use `memory.NewSQLiteHistory(db.DB, ...).AddMessage(...)`.
			hist := memory.NewSQLiteHistory(db.DB, uint(sessionID))
			if err := hist.AddMessage(ctx, aiMsg); err != nil {
				log.Printf("Failed to save tool-call message: %v", err)
			}

			// Execute Tools
			for _, tc := range choice.ToolCalls {
				if err := sendJSON(conn, WSMessage{
					Type:    TypeMessage,
					Content: fmt.Sprintf("\n\n> Calling tool: %s...\n", tc.FunctionCall.Name),
				}); err != nil {
					log.Printf("Failed to send tool notice: %v", err)
				}

				// Execute
				result, err := toolService.ExecuteSkill(tc.FunctionCall.Name, tc.FunctionCall.Arguments)
				if err != nil {
					result = fmt.Sprintf("Error: %v", err)
				}

				// Handle specific tool UI updates (TodoWrite, etc) - Copied from old code
				handleToolUIUpdates(conn, tc.FunctionCall.Name, tc.FunctionCall.Arguments)

				// Save Tool Output
				toolMsg := llms.ToolChatMessage{
					ID:      tc.ID,
					Content: result,
				}
				if err := hist.AddMessage(ctx, toolMsg); err != nil {
					log.Printf("Failed to save tool result: %v", err)
				}
			}

			// Continue loop to generate next response
		} else {
			// No tool calls, just text response
			// Save it
			if err := llmService.SaveAIMessage(ctx, uint(sessionID), choice.Content); err != nil {
				log.Printf("Failed to save AI message: %v", err)
			}
			break
		}
	}

	if err := sendJSON(conn, WSMessage{Type: TypeMessageEnd}); err != nil {
		log.Printf("Failed to send message end: %v", err)
	}
}

func convertToLangChainTools(tools []services.Tool) []llms.Tool {
	var lcTools []llms.Tool
	for _, t := range tools {
		lcTools = append(lcTools, llms.Tool{
			Type: t.Type,
			Function: &llms.FunctionDefinition{
				Name:        t.Function.Name,
				Description: t.Function.Description,
				Parameters:  t.Function.Parameters,
			},
		})
	}
	return lcTools
}

func handleToolUIUpdates(conn *websocket.Conn, name, args string) {
	if name == "TodoWrite" {
		var todoArgs struct {
			Items []TaskDTO `json:"items"`
		}
		if err := json.Unmarshal([]byte(args), &todoArgs); err == nil {
			if err := sendJSON(conn, WSMessage{
				Type:  TypeTaskUpdate,
				Tasks: todoArgs.Items,
			}); err != nil {
				log.Printf("Failed to send task update: %v", err)
			}
		}
	}
	if name == "Task" {
		var taskArgs struct {
			Subagent string `json:"subagent_type"`
			Desc     string `json:"description"`
		}
		if err := json.Unmarshal([]byte(args), &taskArgs); err == nil {
			if err := sendJSON(conn, WSMessage{
				Type:    TypeMessage,
				Content: fmt.Sprintf("\n*Subagent [%s] started: %s*\n", taskArgs.Subagent, taskArgs.Desc),
			}); err != nil {
				log.Printf("Failed to send task start: %v", err)
			}
		}
	}
	if name == "Skill" {
		var skillArgs struct {
			Name string `json:"name"`
		}
		if err := json.Unmarshal([]byte(args), &skillArgs); err == nil {
			if err := sendJSON(conn, WSMessage{
				Type:    TypeMessage,
				Content: fmt.Sprintf("\n*Loaded Skill: %s*\n", skillArgs.Name),
			}); err != nil {
				log.Printf("Failed to send skill message: %v", err)
			}
		}
	}
}

func HandlePermissionResponse(msg WSMessage) {
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

func sendJSON(conn *websocket.Conn, v interface{}) error {
	return conn.WriteJSON(v)
}
