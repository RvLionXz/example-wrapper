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

// Structs untuk request dari client (format OpenAI)
type OpenAIRequest struct {
	Model    string          `json:"model"`
	Messages []OpenAIMessage `json:"messages"`
	Stream   bool            `json:"stream"`
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

// Struct untuk parsing stream dari Gemini
type GeminiStreamResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}
	}
}

// Struct untuk mengirim stream ke client (format OpenAI)
type OpenAIStreamChoice struct {
	Delta struct {
		Content string `json:"content"`
	}
}
type OpenAIStreamResponse struct {
	Choices []OpenAIStreamChoice `json:"choices"`
}

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

	geminiReqBody, err := json.Marshal(GeminiRequest{
		Contents: []Content{{Parts: []Part{{Text: prompt}}}},
	})
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// PERBAIKAN: Menambahkan &alt=sse untuk mendapatkan format stream yang benar
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:streamGenerateContent?key=" + geminiAPIKey + "&alt=sse"
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

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming tidak didukung!", http.StatusInternalServerError)
		return
	}

	// Logika scanner ini sekarang akan bekerja karena format stream dari Gemini sudah benar
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			jsonData := strings.TrimPrefix(line, "data: ")

			var geminiChunk GeminiStreamResponse
			if err := json.Unmarshal([]byte(jsonData), &geminiChunk); err != nil {
				continue
			}

			var textChunk string
			if len(geminiChunk.Candidates) > 0 && len(geminiChunk.Candidates[0].Content.Parts) > 0 {
				textChunk = geminiChunk.Candidates[0].Content.Parts[0].Text
			}

			openAIChunk := OpenAIStreamResponse{
				Choices: []OpenAIStreamChoice{
					{Delta: struct {
						Content string `json:"content"`
					}{Content: textChunk}},
				},
			}

			jsonResponse, err := json.Marshal(openAIChunk)
			if err != nil {
				continue
			}

			fmt.Fprintf(w, "data: %s\n\n", jsonResponse)
			flusher.Flush()
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error membaca stream dari Gemini: %v", err)
	}

	fmt.Fprintf(w, "data: [DONE]\n\n")
	flusher.Flush()
}

func main() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Error: Environment variable GEMINI_API_KEY tidak di-set.")
	}

	http.HandleFunc("/v1/chat/completions", handleChatRequest)

	port := "8080"
	log.Printf("Server Streaming Final (v9) berjalan di port %s...", port)
	log.Println("Endpoint: POST http://localhost:8080/v1/chat/completions")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
