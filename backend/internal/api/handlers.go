package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

// --- Model Handlers ---

func GetModels(c *gin.Context) {
	var config []models.ModelConfig
	if err := db.DB.Find(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func CreateModel(c *gin.Context) {
	var config models.ModelConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func UpdateModel(c *gin.Context) {
	id := c.Param("id")
	var config models.ModelConfig
	if err := db.DB.First(&config, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}
	if err := c.ShouldBindJSON(&config); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Save(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func DeleteModel(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.ModelConfig{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Model deleted"})
}

// --- Conversation Handlers ---

func GetConversations(c *gin.Context) {
	var conversations []models.Conversation
	if err := db.DB.Order("created_at desc").Find(&conversations).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conversations)
}

func CreateConversation(c *gin.Context) {
	var input struct {
		Title   string `json:"title"`
		ModelID uint   `json:"model_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv := models.Conversation{
		Title:   input.Title,
		ModelID: input.ModelID,
	}
	if err := db.DB.Create(&conv).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, conv)
}

func GetConversationMessages(c *gin.Context) {
	id := c.Param("id")
	var messages []models.Message
	if err := db.DB.Where("conversation_id = ?", id).Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func DeleteConversation(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Conversation{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Conversation deleted"})
}

type AvailableModelsRequest struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type ModelInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type AvailableModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

type OpenAIModelsResponse struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int64  `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func GetAvailableModels(c *gin.Context) {
	var req AvailableModelsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.BaseURL == "" || req.APIKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "base_url and api_key are required"})
		return
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	url := fmt.Sprintf("%s/v1/models", req.BaseURL)
	httpReq, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to create request: %v", err)})
		return
	}

	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", req.APIKey))
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": fmt.Sprintf("failed to connect to API: %v", err)})
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to read response: %v", err)})
		return
	}

	if resp.StatusCode == http.StatusUnauthorized {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid API key"})
		return
	}

	if resp.StatusCode != http.StatusOK {
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("API returned status %d: %s", resp.StatusCode, string(body))})
		return
	}

	var openAIResp OpenAIModelsResponse
	if err := json.Unmarshal(body, &openAIResp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("failed to parse response: %v", err)})
		return
	}

	models := make([]ModelInfo, 0, len(openAIResp.Data))
	for _, m := range openAIResp.Data {
		models = append(models, ModelInfo{
			ID:   m.ID,
			Name: m.ID,
		})
	}

	c.JSON(http.StatusOK, AvailableModelsResponse{Models: models})
}
