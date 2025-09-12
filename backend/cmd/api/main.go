package main

import (
	"goclientside/backend/internal/handler"
	"log"
	"net/http"
	"os"
)

func main() {
	// Pastikan API key sudah di-set sebelum memulai server
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Environment variable GEMINI_API_KEY tidak di-set.")
	}

	// Daftarkan handler untuk endpoint kita
	http.HandleFunc("/v1/chat/completions", handler.ChatCompletions)

	port := "8080"
	log.Printf("Server berjalan di port %s...", port)
	log.Println("Endpoint: POST http://localhost:8080/v1/chat/completions")

	// Jalankan server
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
