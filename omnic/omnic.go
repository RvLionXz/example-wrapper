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

type backendRequest struct {
	Prompt string `json:"prompt"`
}

type Part struct {
	Text string `json:"text"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Candidate struct {
	Content Content `json:"content"`
}

type backendResponse struct {
	Candidates []Candidate `json:"candidates"`
}

// NAMA FUNGSI DIPERBAIKI: GenerateContent
func (c *Client) GenerateContent(prompt string) (string, error) {
	reqData := backendRequest{Prompt: prompt}
	jsonBody, err := json.Marshal(reqData)

	if err != nil {
		return "", fmt.Errorf("gagal membuat json request: %w", err)
	}

	fullURL := c.baseURL + "/api/generate"
	req, err := http.NewRequest("POST", fullURL, bytes.NewBuffer(jsonBody))

	if err != nil {
		return "", fmt.Errorf("gagal membuat http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Api-Key", c.apiKey)

	resp, err := c.httpClient.Do(req)

	if err != nil {
		return "", fmt.Errorf("gagal mengirim request server: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("error dari server, status code: %d", resp.StatusCode)
	}

	var backendResp backendResponse

	if err := json.NewDecoder(resp.Body).Decode(&backendResp); err != nil {
		return "", fmt.Errorf("gagal decode json response: %w", err)
	}

	if len(backendResp.Candidates) > 0 && len(backendResp.Candidates[0].Content.Parts) > 0 {
		return backendResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", fmt.Errorf("tidak ada content di dalam response")
}