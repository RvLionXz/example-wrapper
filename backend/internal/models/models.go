package models

// ChatRequest represents the request structure for chat completions
type ChatRequest struct {
	Model       string    `json:"model" binding:"required"`
	Messages    []Message `json:"messages" binding:"required"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// Message represents a single message in a chat
type Message struct {
	Role    string `json:"role" binding:"required"`
	Content string `json:"content" binding:"required"`
}

// EmbeddingRequest represents the request structure for embeddings
type EmbeddingRequest struct {
	Model   string  `json:"model" binding:"required"`
	Content Content `json:"content" binding:"required"`
}

// Content represents the content structure for embeddings
type Content struct {
	Parts []Part `json:"parts" binding:"required"`
}

// Part represents a part of the content
type Part struct {
	Text string `json:"text" binding:"required"`
}

// Gemini API request/response models
type GeminiChatRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

type GenerationConfig struct {
	Temperature *float64 `json:"temperature,omitempty"`
}

type GeminiEmbeddingRequest struct {
	Content Content `json:"content"`
}

type GeminiEmbeddingResponse struct {
	Embedding Embedding `json:"embedding"`
}

type Embedding struct {
	Values []float64 `json:"values"`
}