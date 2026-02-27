package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"fnchatbot/internal/auth"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

// --- Model Handlers ---

func GetModels(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var config []models.ModelConfig
	if err := db.DB.Preload("ProviderRef").Where("user_id = ?", user.ID).Find(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

type CreateModelConfigRequest struct {
	Name        string  `json:"name" binding:"required"`
	Provider    string  `json:"provider"`
	ProviderID  uint    `json:"provider_id"`
	Model       string  `json:"model" binding:"required"`
	Temperature float32 `json:"temperature"`
	MaxTokens   int     `json:"max_tokens"`
	IsDefault   bool    `json:"is_default"`
}

func CreateModel(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var req CreateModelConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	providerID, err := resolveProviderID(req.ProviderID, req.Provider)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config := models.ModelConfig{
		Name:        req.Name,
		Provider:    req.Provider,
		ProviderID:  providerID,
		Model:       req.Model,
		Temperature: req.Temperature,
		MaxTokens:   req.MaxTokens,
		IsDefault:   req.IsDefault,
		UserID:      user.ID,
	}

	if err := db.DB.Create(&config).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, config)
}

func resolveProviderID(providerID uint, providerKey string) (uint, error) {
	if providerID > 0 {
		return providerID, nil
	}
	if providerKey == "" {
		return 0, fmt.Errorf("provider_id or provider is required")
	}
	var provider models.Provider
	if err := db.DB.Where("provider_id = ?", providerKey).First(&provider).Error; err != nil {
		return 0, fmt.Errorf("provider not found")
	}
	return provider.ID, nil
}

// --- Session Handlers ---

func GetSessions(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var sessions []models.Session
	if err := db.DB.Where("user_id = ?", user.ID).Order("created_at desc").Find(&sessions).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sessions)
}

func CreateSession(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	var input struct {
		Title   string `json:"title"`
		ModelID uint   `json:"model_id"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session := models.Session{
		Title:   input.Title,
		ModelID: input.ModelID,
		UserID:  user.ID,
	}
	if err := db.DB.Create(&session).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, session)
}

func GetSessionMessages(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	// Ensure session belongs to current user.
	var session models.Session
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&session).Error; err != nil {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	var messages []models.Message
	if err := db.DB.Where("session_id = ?", id).Preload("Parts").Order("created_at asc").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func DeleteSession(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.Session{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Session deleted"})
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
	defer func() {
		_ = resp.Body.Close()
	}()

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

	modelList := make([]ModelInfo, 0, len(openAIResp.Data))
	for _, m := range openAIResp.Data {
		modelList = append(modelList, ModelInfo{
			ID:   m.ID,
			Name: m.ID,
		})
	}

	c.JSON(http.StatusOK, AvailableModelsResponse{Models: modelList})
}
