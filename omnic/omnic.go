package omnic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type APIRequest struct {
	Messages    []Message `json:"messages"`
	Model       string    `json:"model"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// Methood utama
func (c *Client) ChatCompletionCreate(request APIRequest) (<-chan string, error) {

	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("gagal membuat json request: %w", err)
	}

	URL := c.baseURL + "/v1/chat/completions"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal melakukan request: %w", err)
	}

	// Buat channel yang akan selalu dikembalikan
	resultChan := make(chan string)

	go func() {
		defer resp.Body.Close()
		defer close(resultChan)

		// Cek apakah client meminta streaming
		if request.Stream {
			scanner := bufio.NewScanner(resp.Body)
			for scanner.Scan() {
				line := scanner.Text()
				resultChan <- line
			}
			if err := scanner.Err(); err != nil {
				log.Printf("error membaca stream: %v", err)
			}
		} else {
			geminiResp, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Printf("gagal membaca Response dari gemini %v", err)
				return
			}
			resultChan <- string(geminiResp)
		}
	}()

	return resultChan, nil
}
