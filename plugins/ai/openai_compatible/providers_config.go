package openai_compatible

import (
	"os"
	"strings"

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
	"AIML": {
		Name:    "AIML",
		BaseURL: "https://api.aimlapi.com/v1",
	},
	"Cerebras": {
		Name:    "Cerebras",
		BaseURL: "https://api.cerebras.ai/v1",
	},
	"DeepSeek": {
		Name:    "DeepSeek",
		BaseURL: "https://api.deepseek.com",
	},
	"GrokAI": {
		Name:    "GrokAI",
		BaseURL: "https://api.x.ai/v1",
	},
	"Groq": {
		Name:    "Groq",
		BaseURL: "https://api.groq.com/openai/v1",
	},
	"Langdock": {
		Name:    "Langdock",
		BaseURL: "https://api.langdock.com/openai/{{REGION=us}}/v1",
	},
	"LiteLLM": {
		Name:    "LiteLLM",
		BaseURL: "http://localhost:4000",
	},
	"Mistral": {
		Name:    "Mistral",
		BaseURL: "https://api.mistral.ai/v1",
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
	if strings.Contains(provider.BaseURL, "{{") && strings.Contains(provider.BaseURL, "}}") {
		// Extract the template variable and default value
		start := strings.Index(provider.BaseURL, "{{")
		end := strings.Index(provider.BaseURL, "}}") + 2
		template := provider.BaseURL[start:end]

		// Parse the template to get variable name and default value
		inner := template[2 : len(template)-2] // Remove {{ and }}
		parts := strings.Split(inner, "=")
		if len(parts) == 2 {
			varName := strings.TrimSpace(parts[0])
			defaultValue := strings.TrimSpace(parts[1])

			// Create environment variable name
			envVarName := strings.ToUpper(provider.Name) + "_" + varName

			// Get value from environment or use default
			envValue := os.Getenv(envVarName)
			if envValue == "" {
				envValue = defaultValue
			}

			// Replace the template with the actual value
			provider.BaseURL = strings.Replace(provider.BaseURL, template, envValue, 1)
		}
	}
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
