package omnic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apiKey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: &http.Client{},
	}
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type Choice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

func (c *Client) GenerateContent(request OpenAIRequest) (*OpenAIResponse, error) {

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("gagal mengubah request ke json: %w", err)
	}

	URL := c.baseURL + "/v1/chat/completions"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("request gagal: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server memberikan error dengan status code: %d", resp.StatusCode)
	}

	var openAIResp OpenAIResponse

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("gagal decode respon json: %w", err)
	}

	return &openAIResp, nil

}
