package api

import (
	"net/http"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

type CreateProviderRequest struct {
	ProviderID string                    `json:"provider_id" binding:"required"`
	Name       string                    `json:"name" binding:"required"`
	Type       models.ProviderType       `json:"type" binding:"required"`
	BaseURL    string                    `json:"base_url" binding:"required"`
	APIKey     string                    `json:"api_key"`
	Enabled    bool                      `json:"enabled"`
	ApiOptions models.ProviderApiOptions `json:"api_options"`
}

type UpdateProviderRequest struct {
	ProviderID string                    `json:"provider_id"`
	Name       string                    `json:"name"`
	Type       models.ProviderType       `json:"type"`
	BaseURL    string                    `json:"base_url"`
	APIKey     string                    `json:"api_key"`
	Enabled    bool                      `json:"enabled"`
	ApiOptions models.ProviderApiOptions `json:"api_options"`
}

type FetchModelsResponse struct {
	Models []ModelDetail `json:"models"`
}

type FetchModelsRequest struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

type ModelDetail struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	OwnedBy string `json:"owned_by"`
}

func GetProviders(c *gin.Context) {
	var providers []models.Provider
	if err := db.DB.Preload("Models").Find(&providers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, providers)
}

func GetProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider
	if err := db.DB.Preload("Models").First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}
	c.JSON(http.StatusOK, provider)
}

func CreateProvider(c *gin.Context) {
	var req CreateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider := models.Provider{
		ProviderID: req.ProviderID,
		Name:       req.Name,
		Type:       req.Type,
		BaseURL:    req.BaseURL,
		APIKey:     req.APIKey,
		Enabled:    req.Enabled,
		IsSystem:   false,
		ApiOptions: req.ApiOptions,
	}

	if err := db.DB.Create(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func UpdateProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider
	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	var req UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ProviderID != "" {
		provider.ProviderID = req.ProviderID
	}
	if req.Name != "" {
		provider.Name = req.Name
	}
	if req.Type != "" {
		provider.Type = req.Type
	}
	if req.BaseURL != "" {
		provider.BaseURL = req.BaseURL
	}
	if req.APIKey != "" {
		provider.APIKey = req.APIKey
	}
	provider.Enabled = req.Enabled
	provider.ApiOptions = req.ApiOptions

	if err := db.DB.Save(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func DeleteProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider
	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	if provider.IsSystem {
		c.JSON(http.StatusForbidden, gin.H{"error": "Cannot delete system provider"})
		return
	}

	if err := db.DB.Delete(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Provider deleted"})
}

func ToggleProvider(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider
	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	provider.Enabled = !provider.Enabled
	if err := db.DB.Save(&provider).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func FetchModels(c *gin.Context) {
	id := c.Param("id")
	var provider models.Provider
	if err := db.DB.First(&provider, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	var req FetchModelsRequest
	_ = c.ShouldBindJSON(&req)

	// baseURL := provider.BaseURL
	apiKey := provider.APIKey
	if req.BaseURL != "" {
		// baseURL = req.BaseURL
	}
	if req.APIKey != "" {
		apiKey = req.APIKey
	}

	// Prefer request credentials so model fetching uses current input values
	if apiKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Provider has no API key configured"})
		return
	}

	// adapter := adapters.GetAdapter(provider.Type)
	// if adapter == nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{"error": "No adapter available for provider type: " + string(provider.Type)})
	// 	return
	// }

	// ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	// defer cancel()

	// modelsInfo, err := adapter.FetchModels(ctx, baseURL, apiKey)
	// if err != nil {
	// 	c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	// 	return
	// }

	// result := make([]ModelDetail, 0, len(modelsInfo))
	// for _, m := range modelsInfo {
	// 	result = append(result, ModelDetail{
	// 		ID:      m.ID,
	// 		Name:    m.Name,
	// 		OwnedBy: m.OwnedBy,
	// 	})
	// }

	c.JSON(http.StatusOK, FetchModelsResponse{Models: []ModelDetail{}}) // result})
}
