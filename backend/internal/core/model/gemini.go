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
