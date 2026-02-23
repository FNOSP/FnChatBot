package main

import (
	"log"
	"net/http"
	"os"

	"fnchatbot/internal/api"
	"fnchatbot/internal/api/ws"
	"fnchatbot/internal/db"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize Database
	db.InitDB("fnchatbot.db")

	// Initialize Gin
	r := gin.Default()

	// CORS Configuration
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	r.Use(cors.New(config))

	// Routes
	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"message": "pong",
			})
		})

		// Register REST API handlers
		api.RegisterRoutes(apiGroup)
	}

	// WebSocket Route
	r.GET("/ws/chat/:id", ws.HandleWebSocket)

	// Serve Frontend Static Files (for production)
	// r.StaticFS("/", http.Dir("./dist"))

	// Start Server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
