package api

import (
	"net/http"

	"fnchatbot/internal/auth"
	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

// GetMCPs returns all MCP configurations
func GetMCPs(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var mcps []models.MCPConfig
	if err := db.DB.Where("user_id = ?", user.ID).Find(&mcps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcps)
}

// GetMCP returns a single MCP configuration
func GetMCP(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	var mcp models.MCPConfig
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&mcp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP config not found"})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

// CreateMCP creates a new MCP configuration
func CreateMCP(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var mcp models.MCPConfig
	if err := c.ShouldBindJSON(&mcp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	mcp.UserID = user.ID

	if err := db.DB.Create(&mcp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

// UpdateMCP updates an existing MCP configuration
func UpdateMCP(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	var mcp models.MCPConfig
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).First(&mcp).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP config not found"})
		return
	}

	var input models.MCPConfig
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update fields
	mcp.Name = input.Name
	mcp.BaseURL = input.BaseURL
	mcp.ApiKey = input.ApiKey
	mcp.Enabled = input.Enabled

	if err := db.DB.Save(&mcp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

// DeleteMCP deletes an MCP configuration
func DeleteMCP(c *gin.Context) {
	user, ok := auth.CurrentUser(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	id := c.Param("id")
	if err := db.DB.Where("id = ? AND user_id = ?", id, user.ID).Delete(&models.MCPConfig{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MCP deleted"})
}
