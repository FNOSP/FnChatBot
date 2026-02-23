package api

import (
	"net/http"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"

	"github.com/gin-gonic/gin"
)

// GetMCPs returns all MCP configurations
func GetMCPs(c *gin.Context) {
	var mcps []models.MCPConfig
	if err := db.DB.Find(&mcps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcps)
}

// GetMCP returns a single MCP configuration
func GetMCP(c *gin.Context) {
	id := c.Param("id")
	var mcp models.MCPConfig
	if err := db.DB.First(&mcp, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP config not found"})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

// CreateMCP creates a new MCP configuration
func CreateMCP(c *gin.Context) {
	var mcp models.MCPConfig
	if err := c.ShouldBindJSON(&mcp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := db.DB.Create(&mcp).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, mcp)
}

// UpdateMCP updates an existing MCP configuration
func UpdateMCP(c *gin.Context) {
	id := c.Param("id")
	var mcp models.MCPConfig
	if err := db.DB.First(&mcp, id).Error; err != nil {
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
	id := c.Param("id")
	if err := db.DB.Delete(&models.MCPConfig{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MCP deleted"})
}
