// File: plugins/ai/cerebras/cerebras_test.go
package cerebras

import (
	"testing"
)

// Test the client initialization
func TestNewClient_EmbeddedClientNotNil(t *testing.T) {
	client := NewClient()
	if client.Client == nil {
		t.Fatalf("Expected embedded openai.Client to be non-nil, got nil")
	}
}

// Test the client name and URL configuration
func TestNewClient_ConfiguredCorrectly(t *testing.T) {
	client := NewClient()
	if client.GetName() != "Cerebras" {
		t.Errorf("Expected client name to be 'Cerebras', got '%s'", client.GetName())
	}

	// Check if the ApiBaseURL is set correctly
	if client.ApiBaseURL.Value != "https://api.cerebras.ai/v1" {
		t.Errorf("Expected base URL to be 'https://api.cerebras.ai/v1', got '%s'", client.ApiBaseURL.Value)
	}
}
