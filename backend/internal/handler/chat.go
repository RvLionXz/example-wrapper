package handler

import (
	"encoding/json"
	"fmt"
	"goclientside/backend/internal/model"
	"goclientside/backend/internal/service"
	"io"
	"net/http"
	"os"
)

// ChatCompletions adalah handler untuk endpoint /v1/chat/completions.
func ChatCompletions(w http.ResponseWriter, r *http.Request) {
	geminiAPIKey := os.Getenv("GEMINI_API_KEY")
	if geminiAPIKey == "" {
		http.Error(w, "GEMINI_API_KEY tidak di-set", http.StatusInternalServerError)
		return
	}

	// Decode request body dari client
	var apiReq model.APIRequest
	if err := json.NewDecoder(r.Body).Decode(&apiReq); err != nil {
		http.Error(w, "Bad Request: Gagal decode JSON", http.StatusBadRequest)
		return
	}

	// Panggil service untuk berinteraksi dengan Gemini API
	resp, err := service.CallGeminiAPI(geminiAPIKey, apiReq)
	if err != nil {
		http.Error(w, fmt.Sprintf("Gagal saat memanggil Gemini API: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Salin header dari respons Gemini ke respons kita
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)

	// Salin body (stream atau non-stream) dari respons Gemini ke respons kita
	io.Copy(w, resp.Body)
}
