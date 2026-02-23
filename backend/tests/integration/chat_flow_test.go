package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"fnchatbot/internal/api"
	"fnchatbot/internal/api/ws"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"github.com/gorilla/websocket"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

var (
	testBaseURL string
	testAPIKey  string
	testModelID string
)

func init() {
	// Load environment variables for testing
	testBaseURL = os.Getenv("TEST_BASE_URL")
	if testBaseURL == "" {
		testBaseURL = "https://api.openai.com/v1"
	}
	testAPIKey = os.Getenv("TEST_API_KEY")
	testModelID = os.Getenv("TEST_MODEL_ID")
	if testModelID == "" {
		testModelID = "gpt-4o-mini"
	}
}

func setupTestDB() {
	var err error
	// Use in-memory SQLite for testing
	db.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		// Fallback to on-disk temp file if memory fails (windows issues sometimes)
		db.DB, err = gorm.Open(sqlite.Open("test_integration.db"), &gorm.Config{})
		if err != nil {
			panic(fmt.Sprintf("Failed to connect to test database: %v", err))
		}
	}

	// Migrate schema
	err = db.DB.AutoMigrate(
		&models.ModelConfig{},
		&models.Conversation{},
		&models.Message{},
		&models.Skill{},
		&models.MCPConfig{},
		&models.AgentTask{},
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to migrate test database: %v", err))
	}

	// Clean up data
	db.DB.Exec("DELETE FROM model_configs")
	db.DB.Exec("DELETE FROM conversations")
	db.DB.Exec("DELETE FROM messages")
	db.DB.Exec("DELETE FROM skills")
	db.DB.Exec("DELETE FROM mcp_configs")
}

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.Default()

	apiGroup := r.Group("/api")
	api.RegisterRoutes(apiGroup)

	r.GET("/ws/chat/:id", ws.HandleWebSocket)
	return r
}

func TestChatFlow_Basic(t *testing.T) {
	if testAPIKey == "" {
		t.Skip("TEST_API_KEY not set, skipping integration test")
	}

	setupTestDB()
	r := setupRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	// 1. Add Model
	modelConfig := models.ModelConfig{
		Name:      "Test Model",
		Provider:  "openai",
		BaseURL:   testBaseURL,
		ApiKey:    testAPIKey,
		Model:     testModelID,
		IsDefault: true,
	}

	modelID := createModel(t, srv.URL, modelConfig)

	// 2. Create Conversation
	convID := createConversation(t, srv.URL, "Test Chat", modelID)

	// 3. WebSocket Chat
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat/" + fmt.Sprintf("%d", convID)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer wsConn.Close()

	// 4. Send Message
	userMsg := ws.WSMessage{
		Type:    ws.TypeUserMessage,
		Content: "Hello, this is a test message. Please reply with 'Received'.",
	}
	if err := wsConn.WriteJSON(userMsg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// 5. Receive Stream
	receivedContent := ""
	for {
		var msg ws.WSMessage
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			// Expected close is ok
			break
		}

		if msg.Type == ws.TypeMessage {
			receivedContent += msg.Delta
		} else if msg.Type == ws.TypeMessageEnd {
			break
		}
	}

	if receivedContent == "" {
		t.Error("Received empty response content")
	}
	t.Logf("Full response: %s", receivedContent)
}

func TestChatFlow_WithSkill(t *testing.T) {
	if testAPIKey == "" {
		t.Skip("TEST_API_KEY not set, skipping integration test")
	}

	setupTestDB()
	r := setupRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	// 1. Add Model
	modelConfig := models.ModelConfig{
		Name:      "Test Model",
		Provider:  "openai",
		BaseURL:   testBaseURL,
		ApiKey:    testAPIKey,
		Model:     testModelID,
		IsDefault: true,
	}
	modelID := createModel(t, srv.URL, modelConfig)

	// 2. Add Skill
	skillConfig := map[string]interface{}{
		"parameters": map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"format": map[string]interface{}{
					"type":        "string",
					"description": "The format of the time (e.g. 'iso')",
				},
			},
		},
	}
	skillConfigBytes, _ := json.Marshal(skillConfig)

	skill := models.Skill{
		Name:        "get_current_time",
		Description: "Get the current time.",
		Enabled:     true,
		Config:      datatypes.JSON(skillConfigBytes),
	}
	if err := db.DB.Create(&skill).Error; err != nil {
		t.Fatalf("Failed to create skill: %v", err)
	}

	// 3. Create Conversation
	convID := createConversation(t, srv.URL, "Skill Chat", modelID)

	// 4. Connect WS
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat/" + fmt.Sprintf("%d", convID)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer wsConn.Close()

	// 5. Send Message that triggers skill
	userMsg := ws.WSMessage{
		Type:    ws.TypeUserMessage,
		Content: "What is the current time?",
	}
	if err := wsConn.WriteJSON(userMsg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// 6. Verify Tool Call and Response
	receivedContent := ""

	for {
		var msg ws.WSMessage
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			break
		}

		if msg.Type == ws.TypeMessage {
			receivedContent += msg.Delta
		} else if msg.Type == ws.TypeMessageEnd {
			break
		}
	}

	if !strings.Contains(receivedContent, "2023-10-27") && !strings.Contains(receivedContent, "10:00:00") {
		t.Fatalf("Response content: %s, expected time", receivedContent)
	} else {
		t.Log("Skill execution verified via response content")
	}
}

func TestChatFlow_WithMCP(t *testing.T) {
	if testAPIKey == "" {
		t.Skip("TEST_API_KEY not set, skipping integration test")
	}

	setupTestDB()
	r := setupRouter()
	srv := httptest.NewServer(r)
	defer srv.Close()

	// 1. Mock MCP Server
	mcpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/tools" {
			// Return tools
			tools := []services.Tool{
				{
					Type: "function",
					Function: services.ToolFunction{
						Name:        "mcp_echo",
						Description: "Echo back the input",
						Parameters: json.RawMessage(`{
							"type": "object",
							"properties": {
								"text": {"type": "string"}
							}
						}`),
					},
				},
			}
			json.NewEncoder(w).Encode(tools)
			return
		}
		if req.URL.Path == "/tools/execute" {
			// Execute tool
			var input struct {
				Name string `json:"name"`
				Args string `json:"args"`
			}
			json.NewDecoder(req.Body).Decode(&input)

			if input.Name == "mcp_echo" {
				w.Write([]byte("Echo: " + input.Args))
			} else {
				w.WriteHeader(http.StatusNotFound)
			}
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mcpServer.Close()

	// 2. Add Model
	modelConfig := models.ModelConfig{
		Name:      "Test Model",
		Provider:  "openai",
		BaseURL:   testBaseURL,
		ApiKey:    testAPIKey,
		Model:     testModelID,
		IsDefault: true,
	}
	modelID := createModel(t, srv.URL, modelConfig)

	// 3. Add MCP Config
	mcpConfig := models.MCPConfig{
		Name:    "Test MCP",
		BaseURL: mcpServer.URL,
		Enabled: true,
	}
	if err := db.DB.Create(&mcpConfig).Error; err != nil {
		t.Fatalf("Failed to create MCP config: %v", err)
	}

	// 4. Create Conversation
	convID := createConversation(t, srv.URL, "MCP Chat", modelID)

	// 5. Connect WS
	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/ws/chat/" + fmt.Sprintf("%d", convID)
	wsConn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket: %v", err)
	}
	defer wsConn.Close()

	// 6. Send Message that triggers MCP
	userMsg := ws.WSMessage{
		Type:    ws.TypeUserMessage,
		Content: "Please use mcp_echo to echo 'HelloMCP'.",
	}
	if err := wsConn.WriteJSON(userMsg); err != nil {
		t.Fatalf("Failed to send message: %v", err)
	}

	// 7. Verify
	receivedContent := ""
	for {
		var msg ws.WSMessage
		err := wsConn.ReadJSON(&msg)
		if err != nil {
			break
		}

		if msg.Type == ws.TypeMessage {
			receivedContent += msg.Delta
		} else if msg.Type == ws.TypeMessageEnd {
			break
		}
	}

	if !strings.Contains(receivedContent, "Echo") && !strings.Contains(receivedContent, "HelloMCP") {
		t.Fatalf("Response content: %s, expected Echo", receivedContent)
	} else {
		t.Log("MCP execution verified via response content")
	}
}

// Helpers (Same as before)
func createModel(t *testing.T, baseURL string, config models.ModelConfig) uint {
	body, _ := json.Marshal(config)
	resp, err := http.Post(baseURL+"/api/models", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create model: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create model returned status: %d", resp.StatusCode)
	}

	var model models.ModelConfig
	json.NewDecoder(resp.Body).Decode(&model)
	return model.ID
}

func createConversation(t *testing.T, baseURL string, title string, modelID uint) uint {
	input := map[string]interface{}{
		"title":    title,
		"model_id": modelID,
	}
	body, _ := json.Marshal(input)
	resp, err := http.Post(baseURL+"/api/conversations", "application/json", bytes.NewBuffer(body))
	if err != nil {
		t.Fatalf("Failed to create conversation: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Create conversation returned status: %d", resp.StatusCode)
	}

	var conv models.Conversation
	json.NewDecoder(resp.Body).Decode(&conv)
	return conv.ID
}
