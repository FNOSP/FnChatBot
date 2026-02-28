package api

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有 API 路由
func RegisterRoutes(r *gin.RouterGroup) {
	// 供应商管理
	r.GET("/providers", GetProviders)
	r.GET("/providers/:id", GetProvider)
	r.POST("/providers", CreateProvider)
	r.PUT("/providers/:id", UpdateProvider)
	r.DELETE("/providers/:id", DeleteProvider)
	r.PATCH("/providers/:id/toggle", ToggleProvider)
	r.POST("/providers/:id/fetch-models", FetchModels)

	// 模型管理
	r.GET("/providers/:id/models", GetProviderModels)
	r.POST("/providers/:id/models", AddModelToProvider)
	r.PUT("/models/:id", UpdateModelHandler)
	r.DELETE("/models/:id", DeleteModelHandler)
	r.PATCH("/models/:id/default", SetDefaultModel)

	// 保留旧的 API 以兼容（标记为废弃）
	r.GET("/models", GetModels)                     // 废弃：使用 GET /providers
	r.POST("/models", CreateModel)                  // 废弃：使用 POST /providers/:id/models
	r.POST("/models/available", GetAvailableModels) // 废弃：使用 POST /providers/:id/fetch-models

	// 对话
	r.GET("/conversations", GetSessions)
	r.POST("/conversations", CreateSession)
	r.GET("/conversations/:id/messages", GetSessionMessages)
	r.DELETE("/conversations/:id", DeleteSession)

	// Skills
	r.GET("/skills", GetSkills)
	r.POST("/skills/upload", UploadSkill)
	r.PATCH("/skills/:id", ToggleSkill)
	r.DELETE("/skills/:id", DeleteSkill)

	// MCP (name-based; config in mcp.json)
	r.GET("/mcp", GetMCPs)
	r.POST("/mcp/check", CheckAllMCPs)
	r.GET("/mcp/:name", GetMCP)
	r.POST("/mcp", CreateMCP)
	r.PUT("/mcp/:name", UpdateMCP)
	r.DELETE("/mcp/:name", DeleteMCP)
	r.POST("/mcp/:name/check", CheckMCP)

	// Sandbox
	r.GET("/sandbox", GetSandboxConfig)
	r.PUT("/sandbox", UpdateSandboxConfig)
	r.POST("/sandbox/paths", AddSandboxPath)
	r.DELETE("/sandbox/paths/:path", RemoveSandboxPath)

	// User management
	r.GET("/users", GetUsers)
	r.GET("/users/me", GetCurrentUser)
	r.POST("/users", CreateUser)
	r.PUT("/users/:id", UpdateUser)
	r.PATCH("/users/:id", UpdateUser)
	r.PATCH("/users/:id/password", ChangeUserPassword)
}
