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
				Content: "Apa itu bahasa golang?",
			},
		},
	}

	fmt.Println("Mengirim Request....")

	response, err := client.GenerateContent(request)

	if err != nil {
		log.Fatalf("ERROR: Gagal mendapatkan jawaban dari backend: ", err)
	}

	if len(response.Choices) > 0 {
		fmt.Println(response.Choices[0].Message.Content)
	} else {
		fmt.Println("Tidak ada jawaban yg diterima oleh server")
	}

}
