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

var clientKeys = map[string]bool{
	"supersecret-client-key-123": true,
	"another-client-key-456":     true,
}

var geminiAPIKey string

const geminiAPIURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash-latest:generateContent?key="

type ClientRequest struct {
	Prompt string `json:"prompt"`
}

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

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientKey := r.Header.Get("X-Client-Api-Key")
		if _, valid := clientKeys[clientKey]; !valid {
			http.Error(w, "Unauthorized: API Key tidak valid atau tidak ada.", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func handleGeminiRequest(w http.ResponseWriter, r *http.Request) {
	var clientReq ClientRequest
	if err := json.NewDecoder(r.Body).Decode(&clientReq); err != nil {
		http.Error(w, "Bad Request: Gagal membaca JSON.", http.StatusBadRequest)
		return
	}
	if clientReq.Prompt == "" {
		http.Error(w, "Bad Request: Field 'prompt' tidak boleh kosong.", http.StatusBadRequest)
		return
	}

	geminiReq := GeminiRequest{
		Contents: []Content{
			{
				Parts: []Part{
					{Text: clientReq.Prompt},
				},
			},
		},
	}

	reqBody, err := json.Marshal(geminiReq)
	if err != nil {
		http.Error(w, "Internal Server Error: Gagal membuat request body.", http.StatusInternalServerError)
		return
	}

	fullURL := geminiAPIURL + geminiAPIKey
	resp, err := http.Post(fullURL, "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Internal Server Error: Gagal menghubungi API Gemini.", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "Status: OK. Gemini Wrapper Backend is running.")
}

func main() {
	geminiAPIKey = os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		log.Fatal("Error: Environment variable GEMINI_API_KEY tidak di-set.")
	}

	http.HandleFunc("/", handleHealthCheck)
	http.HandleFunc("/api/generate", authMiddleware(handleGeminiRequest))

	port := "8080"
	log.Printf("Server berjalan di port %s...", port)
	log.Println("Endpoint GET: http://localhost:8080/")
	log.Println("Endpoint POST: http://localhost:8080/api/generate")
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}