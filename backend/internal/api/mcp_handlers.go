package api

import (
	"context"
	"net/http"
	"time"

	"fnchatbot/internal/auth"
	"fnchatbot/internal/models"
	"fnchatbot/internal/services"

	"github.com/gin-gonic/gin"
)

// GetMCPs returns all MCP configs with runtime status (system-level, no user filter).
func GetMCPs(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusOK, []models.MCPServerInfo{})
		return
	}
	f, err := services.DefaultMCPService.LoadFile()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	status := services.DefaultMCPService.GetStatus()
	list := make([]models.MCPServerInfo, 0, len(f.Servers))
	for name, cfg := range f.Servers {
		st := status[name]
		if st.Status == "" {
			st = models.MCPStatus{Status: models.MCPStatusUnknown}
		}
		list = append(list, models.MCPServerInfo{
			Name:            name,
			MCPServerConfig: cfg,
			MCPStatus:       st,
		})
	}
	c.JSON(http.StatusOK, list)
}

// GetMCP returns a single MCP config and status by name.
func GetMCP(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "MCP server not found"})
		return
	}
	cfg, err := services.DefaultMCPService.GetServerConfig(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	status := services.DefaultMCPService.GetStatus()
	st := status[name]
	if st.Status == "" {
		st = models.MCPStatus{Status: models.MCPStatusUnknown}
	}
	c.JSON(http.StatusOK, models.MCPServerInfo{
		Name:            name,
		MCPServerConfig: *cfg,
		MCPStatus:       st,
	})
}

// CreateMCP creates or overwrites an MCP server in mcp.json. Body: { "name": "...", ... MCPServerConfig } or name in body.
func CreateMCP(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MCP service not initialized"})
		return
	}
	var body struct {
		Name string `json:"name"`
		models.MCPServerConfig
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	name := body.Name
	if name == "" {
		name = c.Param("name")
	}
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if err := services.DefaultMCPService.SetServer(name, body.MCPServerConfig); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	status := services.DefaultMCPService.GetStatus()
	st := status[name]
	if st.Status == "" {
		st = models.MCPStatus{Status: models.MCPStatusUnknown}
	}
	c.JSON(http.StatusOK, models.MCPServerInfo{
		Name:            name,
		MCPServerConfig: body.MCPServerConfig,
		MCPStatus:       st,
	})
}

// UpdateMCP updates an MCP server by name. Path param "name" identifies the server.
func UpdateMCP(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "MCP service not initialized"})
		return
	}
	var cfg models.MCPServerConfig
	if err := c.ShouldBindJSON(&cfg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := services.DefaultMCPService.SetServer(name, cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	status := services.DefaultMCPService.GetStatus()
	st := status[name]
	if st.Status == "" {
		st = models.MCPStatus{Status: models.MCPStatusUnknown}
	}
	c.JSON(http.StatusOK, models.MCPServerInfo{
		Name:            name,
		MCPServerConfig: cfg,
		MCPStatus:       st,
	})
}

// DeleteMCP removes an MCP server by name.
func DeleteMCP(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusOK, gin.H{"message": "MCP deleted"})
		return
	}
	if err := services.DefaultMCPService.DeleteServer(name); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "MCP deleted"})
}

// CheckAllMCPs runs status check for all enabled MCP servers and returns status map.
func CheckAllMCPs(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusOK, gin.H{})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
	defer cancel()
	status := services.DefaultMCPService.CheckAllEnabled(ctx)
	c.JSON(http.StatusOK, status)
}

// CheckMCP runs status check for a single MCP server by name.
func CheckMCP(c *gin.Context) {
	if _, ok := auth.CurrentUser(c); !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name required"})
		return
	}
	if services.DefaultMCPService == nil {
		c.JSON(http.StatusOK, models.MCPStatus{Status: models.MCPStatusUnknown})
		return
	}
	ctx, cancel := context.WithTimeout(c.Request.Context(), 30*time.Second)
	defer cancel()
	st := services.DefaultMCPService.CheckServer(ctx, name)
	c.JSON(http.StatusOK, st)
}
