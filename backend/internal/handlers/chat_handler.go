package handlers

import (
	"io"
	"net/http"

	"goclientside/backend-refactored/internal/models"
	"goclientside/backend-refactored/internal/services"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	geminiService *services.GeminiService
}

func NewChatHandler(geminiService *services.GeminiService) *ChatHandler {
	return &ChatHandler{
		geminiService: geminiService,
	}
}

// ChatCompletions handles the /v1/chat/completions endpoint
func (h *ChatHandler) ChatCompletions(c *gin.Context) {
	var request models.ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body: " + err.Error()})
		return
	}

	response, err := h.geminiService.CallGeminiChatAPI(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to call Gemini API: " + err.Error()})
		return
	}

	// Stream the response if requested
	if request.Stream {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		c.Header("Transfer-Encoding", "chunked")
	}

	// Copy the response from Gemini API to the client
	c.Header("Content-Type", response.Header.Get("Content-Type"))
	c.Status(response.StatusCode)
	
	// Stream or copy the response body
	_, err = io.Copy(c.Writer, response.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to stream response: " + err.Error()})
		return
	}
}