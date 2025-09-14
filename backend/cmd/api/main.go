package main

import (
	"goclientside/backend/internal/api/router"
	"log"
	"os"
)

func main() {
	if os.Getenv("GEMINI_API_KEY") == "" {
		log.Fatal("Environment variable GEMINI_API_KEY tidak di-set.")
	}

	// Setup router dari package router
	r := router.SetupRouter()

	port := "8080"
	log.Printf("Server Gin berjalan di port %s...", port)

	// Jalankan server
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server Gin: %v", err)
	}
}
