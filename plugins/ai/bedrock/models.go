package bedrock

import "strings"

// ModelProvider represents the provider of a Bedrock model
type ModelProvider string

const (
	ProviderAnthropic ModelProvider = "anthropic"
	ProviderMeta      ModelProvider = "meta"
	ProviderMistral   ModelProvider = "mistral"
	ProviderCohere    ModelProvider = "cohere"
	ProviderAI21      ModelProvider = "ai21"
)

// ModelInfo contains information about a Bedrock model
type ModelInfo struct {
	ID                      string
	Provider                ModelProvider
	DisplayName             string
	MaxTokens               int
	DefaultMaxTokens        int
	SupportsStreaming       bool
	SupportsSystemMessage   bool
	InputPricePerMillion    float64
	OutputPricePerMillion   float64
}

// ModelRegistry contains information about all supported Bedrock models
var ModelRegistry = map[string]ModelInfo{
	// Anthropic Claude models
	"anthropic.claude-3-opus-20240229-v1:0": {
		ID:                    "anthropic.claude-3-opus-20240229-v1:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3 Opus",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  15.0,
		OutputPricePerMillion: 75.0,
	},
	"anthropic.claude-3-sonnet-20240229-v1:0": {
		ID:                    "anthropic.claude-3-sonnet-20240229-v1:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3 Sonnet",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  3.0,
		OutputPricePerMillion: 15.0,
	},
	"us.anthropic.claude-3-5-sonnet-20240620-v1:0": {
		ID:                    "us.anthropic.claude-3-5-sonnet-20240620-v1:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3.5 Sonnet v1 (US)",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  3.0,
		OutputPricePerMillion: 15.0,
	},
	"us.anthropic.claude-3-5-sonnet-20241022-v2:0": {
		ID:                    "us.anthropic.claude-3-5-sonnet-20241022-v2:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3.5 Sonnet v2 (US)",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  3.0,
		OutputPricePerMillion: 15.0,
	},
	"eu.anthropic.claude-3-5-sonnet-20240620-v1:0": {
		ID:                    "eu.anthropic.claude-3-5-sonnet-20240620-v1:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3.5 Sonnet v1 (EU)",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  3.0,
		OutputPricePerMillion: 15.0,
	},
	"anthropic.claude-3-haiku-20240307-v1:0": {
		ID:                    "anthropic.claude-3-haiku-20240307-v1:0",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 3 Haiku",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  0.25,
		OutputPricePerMillion: 1.25,
	},
	"anthropic.claude-v2:1": {
		ID:                    "anthropic.claude-v2:1",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 2.1",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  8.0,
		OutputPricePerMillion: 24.0,
	},
	"anthropic.claude-v2": {
		ID:                    "anthropic.claude-v2",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude 2",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  8.0,
		OutputPricePerMillion: 24.0,
	},
	"anthropic.claude-instant-v1": {
		ID:                    "anthropic.claude-instant-v1",
		Provider:              ProviderAnthropic,
		DisplayName:           "Claude Instant",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  0.8,
		OutputPricePerMillion: 2.4,
	},

	// Meta Llama models
	"meta.llama3-70b-instruct-v1:0": {
		ID:                    "meta.llama3-70b-instruct-v1:0",
		Provider:              ProviderMeta,
		DisplayName:           "Llama 3 70B Instruct",
		MaxTokens:             2048,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  2.65,
		OutputPricePerMillion: 3.5,
	},
	"meta.llama3-8b-instruct-v1:0": {
		ID:                    "meta.llama3-8b-instruct-v1:0",
		Provider:              ProviderMeta,
		DisplayName:           "Llama 3 8B Instruct",
		MaxTokens:             2048,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: true,
		InputPricePerMillion:  0.3,
		OutputPricePerMillion: 0.6,
	},
	"meta.llama2-70b-chat-v1": {
		ID:                    "meta.llama2-70b-chat-v1",
		Provider:              ProviderMeta,
		DisplayName:           "Llama 2 70B Chat",
		MaxTokens:             2048,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  1.95,
		OutputPricePerMillion: 2.56,
	},
	"meta.llama2-13b-chat-v1": {
		ID:                    "meta.llama2-13b-chat-v1",
		Provider:              ProviderMeta,
		DisplayName:           "Llama 2 13B Chat",
		MaxTokens:             2048,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  0.75,
		OutputPricePerMillion: 1.0,
	},

	// Mistral models
	"mistral.mistral-7b-instruct-v0:2": {
		ID:                    "mistral.mistral-7b-instruct-v0:2",
		Provider:              ProviderMistral,
		DisplayName:           "Mistral 7B Instruct",
		MaxTokens:             8192,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  0.15,
		OutputPricePerMillion: 0.2,
	},
	"mistral.mixtral-8x7b-instruct-v0:1": {
		ID:                    "mistral.mixtral-8x7b-instruct-v0:1",
		Provider:              ProviderMistral,
		DisplayName:           "Mixtral 8x7B Instruct",
		MaxTokens:             4096,
		DefaultMaxTokens:      4096,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  0.45,
		OutputPricePerMillion: 0.7,
	},

	// Cohere models
	"cohere.command-text-v14": {
		ID:                    "cohere.command-text-v14",
		Provider:              ProviderCohere,
		DisplayName:           "Command",
		MaxTokens:             4096,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  1.5,
		OutputPricePerMillion: 2.0,
	},
	"cohere.command-light-text-v14": {
		ID:                    "cohere.command-light-text-v14",
		Provider:              ProviderCohere,
		DisplayName:           "Command Light",
		MaxTokens:             4096,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     true,
		SupportsSystemMessage: false,
		InputPricePerMillion:  0.3,
		OutputPricePerMillion: 0.6,
	},

	// AI21 models
	"ai21.j2-ultra-v1": {
		ID:                    "ai21.j2-ultra-v1",
		Provider:              ProviderAI21,
		DisplayName:           "Jurassic-2 Ultra",
		MaxTokens:             8192,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     false,
		SupportsSystemMessage: false,
		InputPricePerMillion:  12.5,
		OutputPricePerMillion: 12.5,
	},
	"ai21.j2-mid-v1": {
		ID:                    "ai21.j2-mid-v1",
		Provider:              ProviderAI21,
		DisplayName:           "Jurassic-2 Mid",
		MaxTokens:             8192,
		DefaultMaxTokens:      2048,
		SupportsStreaming:     false,
		SupportsSystemMessage: false,
		InputPricePerMillion:  1.25,
		OutputPricePerMillion: 1.25,
	},
}

// GetModelInfo returns information about a specific model
func GetModelInfo(modelID string) (ModelInfo, bool) {
	info, exists := ModelRegistry[modelID]
	return info, exists
}

// GetModelProvider returns the provider for a given model ID
func GetModelProvider(modelID string) ModelProvider {
	info, exists := ModelRegistry[modelID]
	if !exists {
		// Try to extract provider from model ID
		// Handle both direct model IDs (e.g., "anthropic.claude-v2") 
		// and inference profile IDs (e.g., "us.anthropic.claude-3-5-sonnet")
		parts := strings.Split(modelID, ".")
		if len(parts) >= 2 {
			// For inference profiles like "us.anthropic.claude-3-5-sonnet"
			if parts[0] == "us" || parts[0] == "eu" {
				return ModelProvider(parts[1])
			}
			// For direct model IDs like "anthropic.claude-v2"
			return ModelProvider(parts[0])
		}
		return ""
	}
	return info.Provider
}