package handler

import (
	"goclientside/backend/internal/core/model"
	"goclientside/backend/internal/core/service"
	"io"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

// ChatCompletions adalah handler untuk endpoint /v1/chat/completions.
func ChatCompletions(c *gin.Context) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "GEMINI_API_KEY tidak di-set"})
		return
	}

	var apiReq model.APIRequest
	if err := c.ShouldBindJSON(&apiReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request body tidak valid: " + err.Error()})
		return
	}

	resp, err := service.CallGeminiAPI(geminiAPIKey, apiReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal saat memanggil Gemini API: " + err.Error()})
		return
	}
	defer resp.Body.Close()

	// Set header dari respons Gemini
	for key, values := range resp.Header {
		for _, value := range values {
			c.Header(key, value)
		}
	}

	// Stream response body ke client
	c.Stream(func(w io.Writer) bool {
		_, err := io.Copy(w, resp.Body)
		// Jika error atau selesai, hentikan streaming
		return err != nil
	})
}
