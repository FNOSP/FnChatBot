package api

import (
	"net/http"

	"fnchatbot/internal/db"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"

	"github.com/gin-gonic/gin"
)

var sandboxService *services.SandboxService

func InitSandboxService() {
	sandboxService = services.NewSandboxService(db.DB)
}

func GetSandboxConfig(c *gin.Context) {
	enabled := sandboxService.IsEnabled()
	paths := sandboxService.GetAllPaths()

	c.JSON(http.StatusOK, gin.H{
		"enabled": enabled,
		"paths":   paths,
	})
}

func UpdateSandboxConfig(c *gin.Context) {
	var input struct {
		Enabled bool `json:"enabled"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := sandboxService.SetEnabled(input.Enabled); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"enabled": input.Enabled,
	})
}

func AddSandboxPath(c *gin.Context) {
	var input models.SandboxPathInfo
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	if err := sandboxService.AddPath(input.Path, input.Description); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	paths := sandboxService.GetAllPaths()
	c.JSON(http.StatusOK, gin.H{
		"paths": paths,
	})
}

func RemoveSandboxPath(c *gin.Context) {
	path := c.Param("path")
	if path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "path is required"})
		return
	}

	if err := sandboxService.RemovePath(path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Path removed"})
}
