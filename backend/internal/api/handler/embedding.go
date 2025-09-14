package handler

import (
	"goclientside/backend/internal/core/model"
	"goclientside/backend/internal/core/service"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// CreateEmbedding adalah handler untuk endpoint /v1/embeddings.
func CreateEmbedding(c *gin.Context) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GEMINI_API_KEY tidak di-set"})
		return
	}

	var apiReq model.EmbeddingAPIRequest
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid: " + err.Error()})
		return
	}

	// Panggil service untuk mendapatkan embedding
	embeddingResp, err := service.CallGeminiEmbeddingAPI(geminiAPIKey, apiReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal saat memanggil Gemini Embedding API: " + err.Error()})
		return
	}

	// Kirim response
	c.JSON(http.StatusOK, embeddingResp)
}
