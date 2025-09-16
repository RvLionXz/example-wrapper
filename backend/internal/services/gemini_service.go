package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"goclientside/backend-refactored/internal/models"
)

type GeminiService struct {
	apiKey string
}

func NewGeminiService(apiKey string) *GeminiService {
	return &GeminiService{
		apiKey: apiKey,
	}
}

// CallGeminiChatAPI makes a request to the Google Gemini Chat API
func (s *GeminiService) CallGeminiChatAPI(request models.ChatRequest) (*http.Response, error) {
	var prompt string
	for i := len(request.Messages) - 1; i >= 0; i-- {
		if request.Messages[i].Role == "user" {
			prompt = request.Messages[i].Content
			break
		}
	}

	geminiReq := models.GeminiChatRequest{
		Contents: []models.Content{{Parts: []models.Part{{Text: prompt}}}},
	}

	if request.Temperature != nil {
		geminiReq.GenerationConfig = &models.GenerationConfig{Temperature: request.Temperature}
	} else {
		defaultTemp := 0.7
		geminiReq.GenerationConfig = &models.GenerationConfig{Temperature: &defaultTemp}
	}

	geminiReqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal Gemini request: %w", err)
	}

	var endpoint string
	if request.Stream {
		endpoint = "streamGenerateContent"
	} else {
		endpoint = "generateContent"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:%s?key=%s", request.Model, endpoint, s.apiKey)
	if request.Stream {
		url += "&alt=sse"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(geminiReqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}

// CallGeminiEmbeddingAPI makes a request to the Google Gemini Embedding API
func (s *GeminiService) CallGeminiEmbeddingAPI(request models.EmbeddingRequest) (*models.GeminiEmbeddingResponse, error) {
	geminiReq := models.GeminiEmbeddingRequest{
		Content: request.Content,
	}

	geminiReqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal embedding request: %w", err)
	}

	// Use the correct model name format for the URL
	modelName := strings.TrimPrefix(request.Model, "models/")
	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:embedContent", modelName)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(geminiReqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-goog-api-key", s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request to Gemini: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gemini API returned error: %s", string(body))
	}

	var geminiResp models.GeminiEmbeddingResponse
	if err := json.Unmarshal(body, &geminiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &geminiResp, nil
}