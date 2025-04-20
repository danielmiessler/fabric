package openai_compatible

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

// ProviderConfig defines the configuration for an OpenAI-compatible API provider
type ProviderConfig struct {
	Name    string
	BaseURL string
}

// Client is the common structure for all OpenAI-compatible providers
type Client struct {
	*openai.Client
}

// NewClient creates a new OpenAI-compatible client for the specified provider
func NewClient(providerConfig ProviderConfig) *Client {
	client := &Client{}
	client.Client = openai.NewClientCompatible(providerConfig.Name, providerConfig.BaseURL, nil)
	return client
}

// ProviderMap is a map of provider name to ProviderConfig for O(1) lookup
var ProviderMap = map[string]ProviderConfig{
	"Mistral": {
		Name:    "Mistral",
		BaseURL: "https://api.mistral.ai/v1",
	},
	"LiteLLM": {
		Name:    "LiteLLM",
		BaseURL: "http://localhost:4000",
	},
	"Groq": {
		Name:    "Groq",
		BaseURL: "https://api.groq.com/openai/v1",
	},
	"GrokAI": {
		Name:    "GrokAI",
		BaseURL: "https://api.x.ai/v1",
	},
	"DeepSeek": {
		Name:    "DeepSeek",
		BaseURL: "https://api.deepseek.com",
	},
	"Cerebras": {
		Name:    "Cerebras",
		BaseURL: "https://api.cerebras.ai/v1",
	},
	"OpenRouter": {
		Name:    "OpenRouter",
		BaseURL: "https://openrouter.ai/api/v1",
	},
	"SiliconCloud": {
		Name:    "SiliconCloud",
		BaseURL: "https://api.siliconflow.cn/v1",
	},
}

// GetProviderByName returns the provider configuration for a given name with O(1) lookup
func GetProviderByName(name string) (ProviderConfig, bool) {
	provider, found := ProviderMap[name]
	return provider, found
}

// CreateClient creates a new client for a provider by name
func CreateClient(providerName string) (*Client, bool) {
	providerConfig, found := GetProviderByName(providerName)
	if !found {
		return nil, false
	}
	return NewClient(providerConfig), true
}
