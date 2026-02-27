package memory

import (
	"context"
	"encoding/json"
	"testing"

	"fnchatbot/internal/models"

	"github.com/glebarez/sqlite"
	"github.com/stretchr/testify/assert"
	"github.com/tmc/langchaingo/llms"
	"gorm.io/gorm"
)

func TestImageStorage(t *testing.T) {
	// Setup in-memory SQLite
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Migrate schema
	err = db.AutoMigrate(&models.Message{}, &models.Part{})
	assert.NoError(t, err)

	ctx := context.Background()
	sessionID := uint(1)
	history := NewSQLiteHistory(db, sessionID)

	// Create a multimodal message with text and image
	imageContent := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg=="
	imageMime := "image/png"
	imageURL := "data:" + imageMime + ";base64," + imageContent

	msg := MultiModalMessage{
		Type:    llms.ChatMessageTypeHuman,
		Content: "Describe this image",
		Parts: []llms.ContentPart{
			llms.TextPart("Describe this image"),
			llms.ImageURLPart(imageURL),
		},
	}

	// Test 1: AddMessage
	err = history.AddMessage(ctx, msg)
	assert.NoError(t, err)

	// Test 2: Verify stored Part records in DB
	var storedMsg models.Message
	err = db.Preload("Parts").First(&storedMsg, "session_id = ?", sessionID).Error
	assert.NoError(t, err)
	assert.Equal(t, models.RoleUser, storedMsg.Role)
	assert.Len(t, storedMsg.Parts, 2)

	var textPart, filePart models.Part
	for _, p := range storedMsg.Parts {
		if p.Type == models.PartTypeText {
			textPart = p
		} else if p.Type == models.PartTypeFile {
			filePart = p
		}
	}

	// Verify text part
	assert.Equal(t, models.PartTypeText, textPart.Type)
	assert.Equal(t, "Describe this image", textPart.Content)

	// Verify file part
	assert.Equal(t, models.PartTypeFile, filePart.Type)
	assert.Equal(t, imageContent, filePart.Content) // Content should be raw base64

	var meta models.FilePartMeta
	err = json.Unmarshal(filePart.Meta, &meta)
	assert.NoError(t, err)
	assert.Equal(t, imageMime, meta.Mime)
	assert.Equal(t, "image.jpg", meta.Filename) // Default filename in AddMessage

	// Test 3: Retrieve messages
	messages, err := history.Messages(ctx)
	assert.NoError(t, err)
	assert.Len(t, messages, 1)

	retrievedMsg, ok := messages[0].(MultiModalMessage)
	assert.True(t, ok)
	assert.Equal(t, llms.ChatMessageTypeHuman, retrievedMsg.Type)
	assert.Len(t, retrievedMsg.Parts, 2)

	var retrievedTextPart llms.TextContent
	var retrievedImagePart llms.ImageURLContent

	for _, p := range retrievedMsg.Parts {
		switch part := p.(type) {
		case llms.TextContent:
			retrievedTextPart = part
		case llms.ImageURLContent:
			retrievedImagePart = part
		}
	}

	assert.Equal(t, "Describe this image", retrievedTextPart.Text)
	// The retrieved URL should be reconstructed as data URL
	assert.Equal(t, imageURL, retrievedImagePart.URL)
}
