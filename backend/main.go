package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// --- Structs untuk Format OpenAI (Request dari Client) ---

type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// --- Structs untuk Format OpenAI (Response ke Client) ---

type OpenAIResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

type Choice struct {
	Index        int           `json:"index"`
	Message      OpenAIMessage `json:"message"`
	FinishReason string        `json:"finish_reason"`
}

// --- Structs Internal untuk Komunikasi dengan Google Gemini ---

type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

// --- FUNGSI WRAPPER UTAMA (PENERJEMAH) ---

func chatCompletion(request OpenAIRequest, apiKey string) (*OpenAIResponse, error) {
	// 1. Ekstrak prompt dari format OpenAI.
	// Untuk simulasi ini, kita ambil konten dari pesan terakhir.
	var prompt string
	if len(request.Messages) > 0 {
		prompt = request.Messages[len(request.Messages)-1].Content
	}
	if prompt == "" {
		return nil, fmt.Errorf("prompt tidak boleh kosong")
	}

	// 2. Buat dan kirim request ke API Gemini (logika ini tetap sama).
	geminiReqBody, err := json.Marshal(GeminiRequest{
		Contents: []Content{{Parts: []Part{{Text: prompt}}}},
	})
	if err != nil {
		return nil, fmt.Errorf("gagal marshal request ke gemini: %w", err)
	}

	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key=" + apiKey
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(geminiReqBody))
	if err != nil {
		return nil, fmt.Errorf("gagal mengirim request ke gemini: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("gemini merespons dengan error %d: %s", resp.StatusCode, string(body))
	}

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, fmt.Errorf("gagal decode respons dari gemini: %w", err)
	}

	// 3. Terjemahkan respons dari format Gemini ke format OpenAI.
	var responseText string
	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		responseText = geminiResp.Candidates[0].Content.Parts[0].Text
	}

	openAIResp := OpenAIResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()), // Buat ID dummy
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   request.Model, // Kembalikan nama model yang diminta client
		Choices: []Choice{
			{
				Index: 0,
				Message: OpenAIMessage{
					Role:    "assistant",
					Content: responseText,
				},
				FinishReason: "stop",
			},
		},
	}

	return &openAIResp, nil
}

// --- HTTP Handlers ---

func handleChatRequest(w http.ResponseWriter, r *http.Request) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Bad Request: Gagal membaca JSON.", http.StatusBadRequest)
		return
	}

	finalResp, err := chatCompletion(openAIReq, geminiAPIKey)
	if err != nil {
		log.Printf("Error dari fungsi chatCompletion: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(finalResp)
}

// --- MAIN FUNCTION ---

func main() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Error: Environment variable GEMINI_API_KEY tidak di-set.")
	}

	http.HandleFunc("/v1/chat/completions", handleChatRequest)

	port := "8080"
	log.Printf("Server Simulasi OpenAI (v4) berjalan di port %s...", port)
	log.Println("Endpoint: POST http://localhost:8080/v1/chat/completions")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}