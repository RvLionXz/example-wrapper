package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GeminiRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GenerationConfig struct {
	Temperature *float64 `json:"temperature,omitempty"`
}

func handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	var apiReq APIRequest
	if err := json.NewDecoder(r.Body).Decode(&apiReq); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	var prompt string
	for i := len(apiReq.Messages) - 1; i >= 0; i-- {
		if apiReq.Messages[i].Role == "user" {
			prompt = apiReq.Messages[i].Content
			break
		}
	}

	geminiReq := GeminiRequest{
		Contents: []Content{{Parts: []Part{{Text: prompt}}}},
	}
	if apiReq.Temperature != nil {
		geminiReq.GenerationConfig = &GenerationConfig{Temperature: apiReq.Temperature}
	} else {
		defaultTemp := 0.7
		geminiReq.GenerationConfig = &GenerationConfig{Temperature: &defaultTemp}
	}

	geminiReqBody, err := json.Marshal(geminiReq)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	var endpoint string
	if apiReq.Stream {
		endpoint = "streamGenerateContent"
	} else {
		endpoint = "generateContent"
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:%s?key=%s", apiReq.Model, endpoint, geminiAPIKey)
	if apiReq.Stream {
		url += "&alt=sse"
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(geminiReqBody))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func main() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("GEMINI_API_KEY tidak di-set.")
	}

	http.HandleFunc("/v1/chat/completions", handleChatCompletions)

	port := "8080"
	log.Printf("Server Hybrid Sederhana (v14) berjalan di port %s...", port)
	log.Println("Endpoint: POST http://localhost:8080/v1/chat/completions")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
