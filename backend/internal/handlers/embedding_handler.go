package handlers

import (
	"net/http"

	"goclientside/backend-refactored/internal/models"
	"goclientside/backend-refactored/internal/services"

	"github.com/gin-gonic/gin"
)

type EmbeddingHandler struct {
	geminiService *services.GeminiService
}

func NewEmbeddingHandler(geminiService *services.GeminiService) *EmbeddingHandler {
	return &EmbeddingHandler{
		geminiService: geminiService,
	}
}

// CreateEmbedding handles the /v1/embeddings endpoint
func (h *EmbeddingHandler) CreateEmbedding(c *gin.Context) {
	var request models.EmbeddingRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	response, err := h.geminiService.CallGeminiEmbeddingAPI(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Gemini API: " + err.Error()})
		return
	}

	// Return the raw response from Gemini API
	c.JSON(http.StatusOK, response)
}