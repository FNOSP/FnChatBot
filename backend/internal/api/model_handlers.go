package api

import (
	"net/http"
	"strconv"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

func GetProviderModels(c *gin.Context) {
	providerID := c.Param("id")

	var provider models.Provider
	if err := db.DB.Preload("Models").First(&provider, providerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	c.JSON(http.StatusOK, provider.Models)
}

type AddModelRequest struct {
	ModelID                string                   `json:"model_id" binding:"required"`
	Name                   string                   `json:"name" binding:"required"`
	Group                  string                   `json:"group"`
	Description            string                   `json:"description"`
	OwnedBy                string                   `json:"owned_by"`
	Capabilities           []models.ModelCapability `json:"capabilities"`
	SupportedEndpointTypes []models.EndpointType    `json:"supported_endpoint_types"`
	EndpointType           models.EndpointType      `json:"endpoint_type"`
	MaxTokens              int                      `json:"max_tokens"`
	InputPrice             float64                  `json:"input_price"`
	OutputPrice            float64                  `json:"output_price"`
	SupportedTextDelta     bool                     `json:"supported_text_delta"`
	Enabled                bool                     `json:"enabled"`
}

func AddModelToProvider(c *gin.Context) {
	providerID := c.Param("id")

	var provider models.Provider
	if err := db.DB.First(&provider, providerID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Provider not found"})
		return
	}

	var req AddModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	model := models.Model{
		ProviderID:             provider.ID,
		ModelID:                req.ModelID,
		Name:                   req.Name,
		Group:                  req.Group,
		Description:            req.Description,
		OwnedBy:                req.OwnedBy,
		Capabilities:           req.Capabilities,
		SupportedEndpointTypes: req.SupportedEndpointTypes,
		EndpointType:           req.EndpointType,
		MaxTokens:              req.MaxTokens,
		InputPrice:             req.InputPrice,
		OutputPrice:            req.OutputPrice,
		SupportedTextDelta:     req.SupportedTextDelta,
		Enabled:                req.Enabled,
		IsDefault:              false,
	}

	if model.MaxTokens == 0 {
		model.MaxTokens = 4096
	}

	if err := db.DB.Create(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, model)
}

type UpdateModelRequest struct {
	ModelID                string                   `json:"model_id"`
	Name                   string                   `json:"name"`
	Group                  string                   `json:"group"`
	Description            string                   `json:"description"`
	OwnedBy                string                   `json:"owned_by"`
	Capabilities           []models.ModelCapability `json:"capabilities"`
	SupportedEndpointTypes []models.EndpointType    `json:"supported_endpoint_types"`
	EndpointType           models.EndpointType      `json:"endpoint_type"`
	MaxTokens              int                      `json:"max_tokens"`
	InputPrice             float64                  `json:"input_price"`
	OutputPrice            float64                  `json:"output_price"`
	SupportedTextDelta     *bool                    `json:"supported_text_delta"`
	Enabled                *bool                    `json:"enabled"`
}

func UpdateModelHandler(c *gin.Context) {
	modelID := c.Param("id")

	var model models.Model
	if err := db.DB.First(&model, modelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	var req UpdateModelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ModelID != "" {
		model.ModelID = req.ModelID
	}
	if req.Name != "" {
		model.Name = req.Name
	}
	if req.Group != "" {
		model.Group = req.Group
	}
	if req.Description != "" {
		model.Description = req.Description
	}
	if req.OwnedBy != "" {
		model.OwnedBy = req.OwnedBy
	}
	if req.Capabilities != nil {
		model.Capabilities = req.Capabilities
	}
	if req.SupportedEndpointTypes != nil {
		model.SupportedEndpointTypes = req.SupportedEndpointTypes
	}
	if req.EndpointType != "" {
		model.EndpointType = req.EndpointType
	}
	if req.MaxTokens > 0 {
		model.MaxTokens = req.MaxTokens
	}
	if req.InputPrice > 0 {
		model.InputPrice = req.InputPrice
	}
	if req.OutputPrice > 0 {
		model.OutputPrice = req.OutputPrice
	}
	if req.SupportedTextDelta != nil {
		model.SupportedTextDelta = *req.SupportedTextDelta
	}
	if req.Enabled != nil {
		model.Enabled = *req.Enabled
	}

	if err := db.DB.Save(&model).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}

func DeleteModelHandler(c *gin.Context) {
	modelID := c.Param("id")

	id, err := strconv.ParseUint(modelID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid model ID"})
		return
	}

	result := db.DB.Delete(&models.Model{}, uint(id))
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Model deleted successfully"})
}

func SetDefaultModel(c *gin.Context) {
	modelID := c.Param("id")

	var model models.Model
	if err := db.DB.First(&model, modelID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Model not found"})
		return
	}

	tx := db.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Model(&models.Model{}).Where("provider_id = ?", model.ProviderID).Update("is_default", false).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	model.IsDefault = true
	if err := tx.Save(&model).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Commit().Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, model)
}
