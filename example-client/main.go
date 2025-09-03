package main

import (
	"fmt"
	"log"

	"goclientside/omnic"
)

func main() {
	backendURL := "http://localhost:8080"
	clientAPIKey := "supersecret-client-key-123"

	client := omnic.NewClient(backendURL, clientAPIKey)

	prompt := "Jelaskan apa itu goroutine dalam satu kalimat"

	fmt.Printf("-> Mengirim prompt: \"%s\"...\n\n", prompt)

	generatedText, err := client.GenerateContent(prompt)

	if err != nil {
		log.Fatalf("ERROR: Gagal mendapatkan jawaban dari backend: %v", err)
	}

	fmt.Println("---\n JAWABAN ---")
	fmt.Println(generatedText)
}
