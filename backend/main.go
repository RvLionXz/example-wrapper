package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Structs untuk menerima request dari client (format OpenAI)
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
	Stream   bool            `json:"stream"` // Menandakan apakah client meminta streaming
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Structs internal untuk mengirim request ke Google Gemini
type GeminiRequest struct {
	Contents []Content `json:"contents"`
}

type Content struct {
	Parts []Part `json:"parts"`
}

type Part struct {
	Text string `json:"text"`
}

// handleChatRequest menangani semua logika, termasuk streaming.
func handleChatRequest(w http.ResponseWriter, r *http.Request) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")

	var openAIReq OpenAIRequest
	if err := json.NewDecoder(r.Body).Decode(&openAIReq); err != nil {
		http.Error(w, "Bad Request: Gagal membaca JSON.", http.StatusBadRequest)
		return
	}

	var prompt string
	if len(openAIReq.Messages) > 0 {
		prompt = openAIReq.Messages[len(openAIReq.Messages)-1].Content
	}
	if prompt == "" {
		http.Error(w, "Bad Request: Prompt tidak boleh kosong.", http.StatusBadRequest)
		return
	}

	// Siapkan request untuk Gemini
	geminiReqBody, err := json.Marshal(GeminiRequest{
		Contents: []Content{{Parts: []Part{{Text: prompt}}}},
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Gunakan endpoint :streamGenerateContent untuk streaming
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:streamGenerateContent?key=" + geminiAPIKey
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(geminiReqBody))
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Kirim request ke Gemini
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set header untuk client agar tahu ini adalah stream
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	// Siapkan flusher untuk mendorong data ke client secara real-time
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming tidak didukung!", http.StatusInternalServerError)
		return
	}

	// Baca response stream dari Gemini baris per baris
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			// Kirim data mentah (chunk) ke client
			fmt.Fprintf(w, "%s\n\n", line)
			flusher.Flush() // Dorong data ke client saat itu juga
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error membaca stream dari Gemini: %v", err)
	}
}

func main() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Error: Environment variable GEMINI_API_KEY tidak di-set.")
	}

	http.HandleFunc("/v1/chat/completions", handleChatRequest)

	port := "8080"
	log.Printf("Server Streaming (v5) berjalan di port %s...", port)
	log.Println("Endpoint: POST http://localhost:8080/v1/chat/completions")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}