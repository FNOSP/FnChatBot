package memory

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

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

	// Helper to add part
	addPart := func(pType models.PartType, content string, meta interface{}) {
		var metaJSON datatypes.JSON
		if meta != nil {
			b, _ := json.Marshal(meta)
			metaJSON = datatypes.JSON(b)
		}
		parts = append(parts, models.Part{
			MessageID: msg.ID,
			Type:      pType,
			Content:   content,
			Meta:      metaJSON,
		})
	}

	// Iterate message parts if available
	var messageParts []llms.ContentPart
	switch m := message.(type) {
	case MultiModalMessage:
		messageParts = m.Parts
	}

	hasContent := false

	if len(messageParts) > 0 {
		for _, p := range messageParts {
			switch part := p.(type) {
			case llms.TextContent:
				addPart(models.PartTypeText, part.Text, nil)
				hasContent = true
			case llms.ImageURLContent:
				content := part.URL
				mime := "image/jpeg"
				filename := "image.jpg"

				// Check for data URL
				if strings.HasPrefix(content, "data:") {
					parts := strings.Split(content, ",")
					if len(parts) == 2 {
						metaPart := parts[0]
						dataPart := parts[1]
						if strings.Contains(metaPart, ";") {
							mime = strings.TrimPrefix(strings.Split(metaPart, ";")[0], "data:")
						}
						content = dataPart
					}
				}

				addPart(models.PartTypeFile, content, models.FilePartMeta{
					Mime:     mime,
					Filename: filename,
				})
				hasContent = true
			case llms.BinaryContent:
				encoded := base64.StdEncoding.EncodeToString(part.Data)
				addPart(models.PartTypeFile, encoded, models.FilePartMeta{
					Mime:     part.MIMEType,
					Filename: "file.bin",
				})
				hasContent = true
			}
		}
	}

	// Fallback to simple content if no parts processed
	if !hasContent {
		if content := message.GetContent(); content != "" {
			partType := models.PartTypeText
			var meta interface{}

			// Handle ToolCallID for Tool messages
			if toolMsg, ok := message.(llms.ToolChatMessage); ok && toolMsg.ID != "" {
				partType = models.PartTypeToolResult
				meta = map[string]interface{}{
					"tool_call_id": toolMsg.ID,
				}
			}

			addPart(partType, content, meta)
		}
	}

	// Handle ToolCalls for AI messages
	if aiMsg, ok := message.(llms.AIChatMessage); ok && len(aiMsg.ToolCalls) > 0 {
		toolCallsJSON, err := json.Marshal(aiMsg.ToolCalls)
		if err == nil {
			addPart(models.PartTypeToolCall, string(toolCallsJSON), nil)
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
		var parts []llms.ContentPart
		var toolCalls []llms.ToolCall
		var toolCallID string
		var contentStr string

		for _, part := range msg.Parts {
			switch part.Type {
			case models.PartTypeText:
				parts = append(parts, llms.TextPart(part.Content))
				contentStr += part.Content
			case models.PartTypeFile:
				var meta models.FilePartMeta
				_ = json.Unmarshal(part.Meta, &meta)

				url := part.Content
				if meta.Mime != "" {
					url = fmt.Sprintf("data:%s;base64,%s", meta.Mime, part.Content)
				}
				parts = append(parts, llms.ImageURLPart(url))
			case models.PartTypeToolCall:
				_ = json.Unmarshal([]byte(part.Content), &toolCalls)
			case models.PartTypeToolResult:
				parts = append(parts, llms.TextPart(part.Content))
				contentStr += part.Content
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
			if len(parts) > 0 {
				chatMsg = MultiModalMessage{
					Type:    llms.ChatMessageTypeHuman,
					Content: contentStr,
					Parts:   parts,
				}
			} else {
				chatMsg = llms.HumanChatMessage{
					Content: contentStr,
				}
			}
		case models.RoleAssistant:
			aiMsg := llms.AIChatMessage{
				Content: contentStr,
			}
			if len(toolCalls) > 0 {
				aiMsg.ToolCalls = toolCalls
			}
			chatMsg = aiMsg
		case models.RoleSystem:
			if len(parts) > 0 {
				chatMsg = MultiModalMessage{
					Type:    llms.ChatMessageTypeSystem,
					Content: contentStr,
					Parts:   parts,
				}
			} else {
				chatMsg = llms.SystemChatMessage{
					Content: contentStr,
				}
			}
		case "tool":
			toolMsg := llms.ToolChatMessage{
				Content: contentStr,
			}
			if toolCallID != "" {
				toolMsg.ID = toolCallID
			}
			chatMsg = toolMsg
		default:
			chatMsg = llms.HumanChatMessage{
				Content: contentStr,
			}
		}
		chatMessages = append(chatMessages, chatMsg)
	}

	return chatMessages, nil
}

// MultiModalMessage represents a message with multiple parts (text, image, etc.)
type MultiModalMessage struct {
	Type    llms.ChatMessageType
	Content string
	Parts   []llms.ContentPart
}

func (m MultiModalMessage) GetType() llms.ChatMessageType {
	return m.Type
}

func (m MultiModalMessage) GetContent() string {
	return m.Content
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
