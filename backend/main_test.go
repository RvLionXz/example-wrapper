package main

import (
	"testing"

	"goclientside/backend-refactored/internal/services"
)

func TestGeminiServiceCreation(t *testing.T) {
	// Test that we can create a GeminiService
	service := services.NewGeminiService("test-key")
	if service == nil {
		t.Error("Expected GeminiService to be created, got nil")
	}
}