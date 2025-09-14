package model

// APIRequest adalah struktur untuk request yang masuk dari client.
type APIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Stream      bool      `json:"stream"`
	Temperature *float64  `json:"temperature,omitempty"`
}

// Message adalah struktur untuk pesan di dalam APIRequest.
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// EmbeddingAPIRequest adalah struktur untuk request embedding yang masuk dari client.
type EmbeddingAPIRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

// EmbeddingAPIResponse adalah struktur untuk response embedding yang dikirim ke client.
type EmbeddingAPIResponse struct {
	Object string      `json:"object"`
	Data   []Embedding `json:"data"`
	Model  string      `json:"model"`
}

// Embedding adalah struktur data untuk satu item embedding.
type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int       `json:"index"`
}
