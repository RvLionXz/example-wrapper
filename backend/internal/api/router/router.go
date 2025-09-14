package router

import (
	"goclientside/backend/internal/api/handler"

	"github.com/gin-gonic/gin"
)

// SetupRouter mengkonfigurasi dan mengembalikan instance dari Gin engine.
func SetupRouter() *gin.Engine {
	// Membuat router dengan middleware default (logger, recovery)
	r := gin.Default()

	// Grup rute untuk v1
	v1 := r.Group("/v1")
	{
		// Endpoint untuk chat completions
		v1.POST("/chat/completions", handler.ChatCompletions)
		// Endpoint untuk embeddings
		v1.POST("/embeddings", handler.CreateEmbedding)
	}

	return r
}
