package main

import (
	"log"
	"os"

	"goclientside/backend-refactored/internal/handlers"
	"goclientside/backend-refactored/internal/middleware"
	"goclientside/backend-refactored/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title Gemini API Wrapper
// @version 1.0
// @description A simple wrapper for Google Gemini API
// @host localhost:8080
// @BasePath /v1
func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get API key from environment
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	// Initialize services
	geminiService := services.NewGeminiService(apiKey)

	// Initialize handlers
	chatHandler := handlers.NewChatHandler(geminiService)
	embeddingHandler := handlers.NewEmbeddingHandler(geminiService)

	// Set up Gin router
	router := gin.New()
	
	// Add middlewares
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.CORSMiddleware())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"message": "Gemini API Wrapper is running",
		})
	})

	// API v1 routes
	v1 := router.Group("/v1")
	{
		v1.POST("/chat/completions", chatHandler.ChatCompletions)
		v1.POST("/embeddings", embeddingHandler.CreateEmbedding)
	}

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}