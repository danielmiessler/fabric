package anthropic

import (
	"testing"
)

// Test generated using Keploy
func TestNewClient_DefaultInitialization(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("Expected client to be initialized, got nil")
	}

	if client.ApiBaseURL.Value != defaultBaseUrl {
		t.Errorf("Expected default API Base URL to be %s, got %s", defaultBaseUrl, client.ApiBaseURL.Value)
	}

	if client.maxTokens != 4096 {
		t.Errorf("Expected default maxTokens to be 4096, got %d", client.maxTokens)
	}

	if len(client.models) == 0 {
		t.Error("Expected models to be initialized with default values, got empty list")
	}
}

// Test generated using Keploy
func TestClientListModels(t *testing.T) {
	client := NewClient()

	models, err := client.ListModels()
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(models) != len(client.models) {
		t.Errorf("Expected %d models, got %d", len(client.models), len(models))
	}

	for i, model := range models {
		if model != client.models[i] {
			t.Errorf("Expected model at index %d to be %s, got %s", i, client.models[i], model)
		}
	}
}

func TestClient_ListModels_ReturnsCorrectModels(t *testing.T) {
	client := NewClient()
	models, err := client.ListModels()

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(models) != len(client.models) {
		t.Errorf("Expected %d models, got %d", len(client.models), len(models))
	}

	for i, model := range models {
		if model != client.models[i] {
			t.Errorf("Expected model %s at index %d, got %s", client.models[i], i, model)
		}
	}
}
