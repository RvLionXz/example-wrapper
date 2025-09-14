package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"goclientside/backend/internal/core/model"
	"net/http"
)

// CallGeminiAPI membuat dan mengirim request ke Google Gemini API.
func CallGeminiAPI(apiKey string, apiReq model.APIRequest) (*http.Response, error) {
	var prompt string
	for i := len(apiReq.Messages) - 1; i >= 0; i-- {
		if apiReq.Messages[i].Role == "user" {
			prompt = apiReq.Messages[i].Content
			break
		}
	}

	geminiReq := model.GeminiRequest{
		Contents: []model.Content{{Parts: []model.Part{{Text: prompt}}}},
	}

	if apiReq.Temperature != nil {
		geminiReq.GenerationConfig = &model.GenerationConfig{Temperature: apiReq.Temperature}
	} else {
		defaultTemp := 0.7
		geminiReq.GenerationConfig = &model.GenerationConfig{Temperature: &defaultTemp}
	}

	geminiReqBody, err := json.Marshal(geminiReq)
	if err != nil {
		return nil, fmt.Errorf("gagal marshal request gemini: %w", err)
	}

	var endpoint string
	if apiReq.Stream {
		endpoint = "streamGenerateContent"
	} else {
		endpoint = "generateContent"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:%s?key=%s", apiReq.Model, endpoint, apiKey)
	if apiReq.Stream {
		url += "&alt=sse"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(geminiReqBody))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	return http.DefaultClient.Do(req)
}
