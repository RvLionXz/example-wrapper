package omnic

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

func NewClient(baseURL, apikey string) *Client {
	return &Client{
		baseURL:    baseURL,
		apiKey:     apikey,
		httpClient: &http.Client{},
	}
}

type OpenAiMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAiRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAiMessage `json:"messages"`
	Stream   bool            `json:"stream"`
}

type Choices struct {
	Index        int             `json:"index"`
	Messages     []OpenAiMessage `json:"messages"`
	FinishReason string          `json:"finish_reason"`
}

type OpenAiResponse struct {
	ID      string    `json:"id"`
	Object  string    `json:"object"`
	Created int       `json:"created"`
	Model   string    `json:"model"`
	Choices []Choices `json:"choices"`
}

func (client *Client) GenerateContent(request OpenAiRequest) (<-chan string, error) {

	request.Stream = true

	// mengubah struc -> json
	jsonBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("gagal mengubah struc ke json: %w", err)
	}

	// membuat http request
	URL := client.baseURL + "/v1/chat/completions"
	req, err := http.NewRequest("POST", URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("gagal membuat http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Client-Api-Key", client.apiKey)

	// mengirim request ke server
	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request ke server: %w", err)
	}

	// Mengirim data response.body dengan stream (go routine)
	streamChan := make(chan string) // pipa untuk mengirim data string
	go func() {

		defer close(streamChan) // memastikan pipa pengirim ditutup ketika response selesai
		defer resp.Body.Close() // memastikan koneksi ke server ditutup ketika response selesai

		reader := bufio.NewScanner(resp.Body)

		for reader.Scan() {
			line := reader.Text() // baca baris satu per satu

			// buat kondisi untuk mengecek baris yg kosong
			if line == "" {
				continue
			}

			jsonData := strings.TrimPrefix(line, "data: ") // kita buang prefix "data :"
			// buat kondisi berhenti ketika server selesai mengirim sinyal
			if jsonData == "[DONE]" {
				break
			}

			// membuat struct sementara untuk kita parsing sebagai potongan potongan json
			type Delta struct {
				Content string `json:"content"`
			}

			type StreamChoice struct {
				Delta Delta `json:"delta"`
			}

			type StreamResponse struct {
				Choices []StreamChoice `json:"choices"`
			}

			// membuat variable streamResp dengan tipe StreamResponse
			var streamResp StreamResponse
			// parsing json ke struct kemudian konversi ke byte
			err := json.Unmarshal([]byte(jsonData), &streamResp)
			if err != nil {
				continue
			}

			// cek jika ada teks di potongan data streamResp maka masukan teks tersebut kedalam streamChan
			if len(streamResp.Choices) > 0 {
				streamChan <- streamResp.Choices[0].Delta.Content
			}

		}
	}()

	return streamChan, nil
}
