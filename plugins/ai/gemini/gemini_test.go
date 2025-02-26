package gemini

import (
	"testing"

	"github.com/google/generative-ai-go/genai"
)

// Test generated using Keploy
func TestBuildModelNameSimple(t *testing.T) {
	client := &Client{}
	fullModelName := "models/chat-bison-001"
	expected := "chat-bison-001"

	result := client.buildModelNameSimple(fullModelName)

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

// Test generated using Keploy
func TestExtractText(t *testing.T) {
	client := &Client{}
	response := &genai.GenerateContentResponse{
		Candidates: []*genai.Candidate{
			{
				Content: &genai.Content{
					Parts: []genai.Part{
						genai.Text("Hello, "),
						genai.Text("world!"),
					},
				},
			},
		},
	}
	expected := "Hello, world!"

	result := client.extractText(response)

	if result != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}
