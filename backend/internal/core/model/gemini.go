package model

// GeminiRequest adalah struktur request yang dikirim ke Google Gemini API.
type GeminiRequest struct {
	Contents         []Content         `json:"contents"`
	GenerationConfig *GenerationConfig `json:"generationConfig,omitempty"`
}

// Content adalah bagian dari GeminiRequest.
type Content struct {
	Parts []Part `json:"parts"`
}

// Part adalah bagian dari Content.
type Part struct {
	Text string `json:"text"`
}

// GenerationConfig mengatur parameter seperti temperature.
type GenerationConfig struct {
	Temperature *float64 `json:"temperature,omitempty"`
}

// GeminiEmbeddingRequest adalah struktur request untuk embedding ke Gemini.
type GeminiEmbeddingRequest struct {
	Content Content `json:"content"`
}

// GeminiEmbeddingResponse adalah struktur response embedding dari Gemini.
type GeminiEmbeddingResponse struct {
	Embedding GeminiEmbedding `json:"embedding"`
}

// GeminiEmbedding berisi nilai-nilai vektor embedding.
type GeminiEmbedding struct {
	Values []float64 `json:"values"`
}
