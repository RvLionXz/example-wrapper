package main

import (
	"fmt"
	"goclientside/omnic"
	"log"
)

func main() {
	baseURL := "http://localhost:8080"

	client := omnic.NewClient(baseURL)

	request := omnic.OpenAIRequest{
		Model: "gemini-1.5-flash-latest",
		Messages: []omnic.OpenAIMessage{
			{
				Role:    "user",
				Content: "Apa itu bahasa pemrograman Go? Jelaskan dalam satu paragraf singkat.",
			},
		},
	}

	fmt.Println("-> Mengirim Request...")

	response, err := client.GenerateContent(request)

	if err != nil {
		log.Fatalf("ERROR: Gagal mendapatkan jawaban dari backend: %v", err)
	}

	if len(response.Choices) > 0 {
		fmt.Println("--- Jawaban dari Server ---")
		fmt.Println(response.Choices[0].Message.Content)
	} else {
		fmt.Println("Tidak ada jawaban yang diterima dari server.")
	}

}