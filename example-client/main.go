package main

import (
	"fmt"
	"goclientside/omnic"
	"log"
)

func main() {
	baseURl := "http://localhost:8080"
	apiKey := "kunci-rahasia-client-A-123"

	client := omnic.NewClient(baseURl, apiKey)

	request := omnic.OpenAiRequest{
		Model: "gemini-1.5-flash-latest",
		Messages: []omnic.OpenAiMessage{
			omnic.OpenAiMessage{
				Role:    "user",
				Content: "Halo saya baru belajar bahasa golangn",
			},
		},
		Stream: true,
	}

	response, err := client.GenerateContent(request)

	if err != nil {
		log.Fatal("ERROR: gagal mendapatkan jawaban dari server: %w", err)
	}

	for textChunk := range response {
		fmt.Print(textChunk)
		fmt.Println()
	}
}
