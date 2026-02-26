package memory

import (
	"context"
	"encoding/json"

	"fnchatbot/internal/models"

	"github.com/tmc/langchaingo/llms"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// SQLiteHistory implements schema.ChatMessageHistory using GORM and SQLite
type SQLiteHistory struct {
	DB        *gorm.DB
	SessionID uint
}

// NewSQLiteHistory creates a new SQLiteHistory
func NewSQLiteHistory(db *gorm.DB, sessionID uint) *SQLiteHistory {
	return &SQLiteHistory{
		DB:        db,
		SessionID: sessionID,
	}
}

// AddMessage adds a message to the history
func (h *SQLiteHistory) AddMessage(ctx context.Context, message llms.ChatMessage) error {
	var role models.MessageRole
	switch message.GetType() {
	case llms.ChatMessageTypeAI:
		role = models.RoleAssistant
	case llms.ChatMessageTypeSystem:
		role = models.RoleSystem
	case llms.ChatMessageTypeHuman:
		role = models.RoleUser
	case llms.ChatMessageTypeGeneric:
		role = models.RoleUser
	case llms.ChatMessageTypeTool:
		role = "tool"
	}

	msg := models.Message{
		SessionID: h.SessionID,
		Role:      role,
	}

	if err := h.DB.Create(&msg).Error; err != nil {
		return err
	}

	var parts []models.Part

	// Add content part
	if content := message.GetContent(); content != "" {
		partType := "text"
		var meta datatypes.JSON

		// Handle ToolCallID for Tool messages
		if toolMsg, ok := message.(llms.ToolChatMessage); ok && toolMsg.ID != "" {
			partType = "tool_result"
			metaJSON, _ := json.Marshal(map[string]interface{}{
				"tool_call_id": toolMsg.ID,
			})
			meta = datatypes.JSON(metaJSON)
		}

		parts = append(parts, models.Part{
			MessageID: msg.ID,
			Type:      partType,
			Content:   content,
			Meta:      meta,
		})
	}

	// Handle ToolCalls for AI messages
	if aiMsg, ok := message.(llms.AIChatMessage); ok && len(aiMsg.ToolCalls) > 0 {
		toolCallsJSON, err := json.Marshal(aiMsg.ToolCalls)
		if err == nil {
			parts = append(parts, models.Part{
				MessageID: msg.ID,
				Type:      "tool_calls",
				Content:   string(toolCallsJSON),
			})
		}
	}

	if len(parts) > 0 {
		return h.DB.Create(&parts).Error
	}
	return nil
}

// AddUserMessage adds a user message to the history
func (h *SQLiteHistory) AddUserMessage(ctx context.Context, message string) error {
	return h.AddMessage(ctx, llms.HumanChatMessage{Content: message})
}

// AddAIMessage adds an AI message to the history
func (h *SQLiteHistory) AddAIMessage(ctx context.Context, message string) error {
	return h.AddMessage(ctx, llms.AIChatMessage{Content: message})
}

// Clear clears the history
func (h *SQLiteHistory) Clear(ctx context.Context) error {
	return h.DB.Where("session_id = ?", h.SessionID).Delete(&models.Message{}).Error
}

// Messages retrieves all messages from the history
func (h *SQLiteHistory) Messages(ctx context.Context) ([]llms.ChatMessage, error) {
	var dbMessages []models.Message
	if err := h.DB.Where("session_id = ?", h.SessionID).Preload("Parts").Order("created_at asc").Find(&dbMessages).Error; err != nil {
		return nil, err
	}

	var chatMessages []llms.ChatMessage
	for _, msg := range dbMessages {
		var content string
		var toolCalls []llms.ToolCall
		var toolCallID string

		for _, part := range msg.Parts {
			switch part.Type {
			case "text":
				content += part.Content
			case "tool_calls":
				_ = json.Unmarshal([]byte(part.Content), &toolCalls)
			case "tool_result":
				content += part.Content
				if len(part.Meta) > 0 {
					var meta map[string]interface{}
					_ = json.Unmarshal(part.Meta, &meta)
					if id, ok := meta["tool_call_id"].(string); ok {
						toolCallID = id
					}
				}
			}
		}

		var chatMsg llms.ChatMessage
		switch msg.Role {
		case models.RoleUser:
			chatMsg = llms.HumanChatMessage{Content: content}
		case models.RoleAssistant:
			aiMsg := llms.AIChatMessage{Content: content}
			if len(toolCalls) > 0 {
				aiMsg.ToolCalls = toolCalls
			}
			chatMsg = aiMsg
		case models.RoleSystem:
			chatMsg = llms.SystemChatMessage{Content: content}
		case "tool":
			toolMsg := llms.ToolChatMessage{Content: content}
			if toolCallID != "" {
				toolMsg.ID = toolCallID
			}
			chatMsg = toolMsg
		default:
			chatMsg = llms.HumanChatMessage{Content: content}
		}
		chatMessages = append(chatMessages, chatMsg)
	}

	return chatMessages, nil
}

// SetMessages sets the messages in the history
func (h *SQLiteHistory) SetMessages(ctx context.Context, messages []llms.ChatMessage) error {
	if err := h.Clear(ctx); err != nil {
		return err
	}
	for _, msg := range messages {
		if err := h.AddMessage(ctx, msg); err != nil {
			return err
		}
	}
	return nil
}
