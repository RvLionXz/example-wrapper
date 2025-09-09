package main

import (
	"fmt"
	"goclientside/omnic"
	"log"
)

func main() {

	client := omnic.NewClient("http://localhost:8080")

	apiRequest := omnic.APIRequest{
		Model: "gemini-1.5-flash-latest",
		Messages: []omnic.Message{
			{
				Role:    "user",
				Content: "Apa itu goroutine di Go? Jelaskan seolah saya anak 5 tahun.",
			},
		},
		Stream: false,
	}

	streamChan, err := client.ChatCompletionCreate(apiRequest)
	if err != nil {
		log.Fatalf("ERROR: Gagal memulai koneksi stream: %v", err)
	}

	for textChunk := range streamChan {
		fmt.Print(textChunk)
	}

}
