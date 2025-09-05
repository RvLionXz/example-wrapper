package omnic

import (
	"bytes"
	"encoding/json"
	"fmt"
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

// Struc untuk request
type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

// struc untuk respon
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

// method untuk request dan response (generate konten)
func (c *Client) GenerateContent(request OpenAIRequest) (*OpenAIResponse, error) {

	// encode struc -> json
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("Gagal mengubah request ke json: %w", err)
	}

	// objek post request
	URL := c.baseURL + "/v1/chat/completions"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("request gagal", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// mengirim request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Gagal melakukan request", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Server memberikan error: ", resp.StatusCode)
	}

	// decode body json -> struc
	var openAIResp OpenAIResponse

	if err := json.NewDecoder(resp.Body).Decode(&openAIResp); err != nil {
		return nil, fmt.Errorf("Gagal decode respon json: ", err)
	}

	return &openAIResp, nil

}
