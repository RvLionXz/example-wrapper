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
