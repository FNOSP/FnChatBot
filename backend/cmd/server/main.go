package main

import (
	"log"
	"net/http"
	"os"

	"fnchatbot/internal/api"
	"fnchatbot/internal/api/ws"
	"fnchatbot/internal/auth"
	"fnchatbot/internal/config"
	"fnchatbot/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load application configuration.
	appCfg := config.LoadConfig()

	// Initialize Database
	db.InitDB("fnchatbot.db")

	// Initialize authentication and authorization services.
	if err := auth.Init(db.DB, appCfg.Auth); err != nil {
		log.Fatalf("Failed to initialize auth: %v", err)
	}

	// Initialize Gin
	r := gin.Default()

	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Public API routes (no auth required).
	publicAPI := r.Group("/api")
	{
		publicAPI.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Login is public; me/reset-password need JWT but not Casbin.
		api.RegisterPublicAuthRoutes(publicAPI)
	}

	// Auth-required but Casbin-exempt routes (for password reset flow).
	authOnlyAPI := publicAPI.Group("")
	authOnlyAPI.Use(auth.JWTMiddleware())
	api.RegisterAuthProtectedRoutes(authOnlyAPI)

	// Protected API routes (require auth & authorization).
	protectedAPI := publicAPI.Group("")
	protectedAPI.Use(auth.JWTMiddleware(), auth.CasbinMiddleware())
	api.RegisterRoutes(protectedAPI)

	// WebSocket Route
	r.GET("/ws/chat/:id", auth.WebSocketAuthMiddleware(), ws.HandleWebSocket)

	// Serve Frontend Static Files (for production)
	// r.StaticFS("/", http.Dir("./dist"))

	// Start Server
	port := appCfg.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
