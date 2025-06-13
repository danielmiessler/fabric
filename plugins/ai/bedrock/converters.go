package bedrock

import (
	"encoding/json"
	"fmt"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"
)

// Message format structures for different providers

// AnthropicMessage represents a message in Anthropic's format
type AnthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AnthropicRequest represents the request format for Anthropic models
type AnthropicRequest struct {
	Messages         []AnthropicMessage `json:"messages"`
	System           string             `json:"system,omitempty"`
	MaxTokens        int                `json:"max_tokens"`
	Temperature      float32            `json:"temperature,omitempty"`
	TopP             float32            `json:"top_p,omitempty"`
	StopSequences    []string           `json:"stop_sequences,omitempty"`
	AnthropicVersion string             `json:"anthropic_version"`
}

// MetaRequest represents the request format for Meta Llama models
type MetaRequest struct {
	Prompt       string  `json:"prompt"`
	MaxGenLen    int     `json:"max_gen_len"`
	Temperature  float32 `json:"temperature,omitempty"`
	TopP         float32 `json:"top_p,omitempty"`
}

// MistralMessage represents a message in Mistral's format
type MistralMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// MistralRequest represents the request format for Mistral models
type MistralRequest struct {
	Prompt       string  `json:"prompt"`
	MaxTokens    int     `json:"max_tokens"`
	Temperature  float32 `json:"temperature,omitempty"`
	TopP         float32 `json:"top_p,omitempty"`
	TopK         int     `json:"top_k,omitempty"`
}

// CohereRequest represents the request format for Cohere models
type CohereRequest struct {
	Prompt       string   `json:"prompt"`
	MaxTokens    int      `json:"max_tokens"`
	Temperature  float32  `json:"temperature,omitempty"`
	P            float32  `json:"p,omitempty"`
	K            int      `json:"k,omitempty"`
	StopSequences []string `json:"stop_sequences,omitempty"`
}

// AI21Request represents the request format for AI21 models
type AI21Request struct {
	Prompt       string   `json:"prompt"`
	MaxTokens    int      `json:"maxTokens"`
	Temperature  float32  `json:"temperature,omitempty"`
	TopP         float32  `json:"topP,omitempty"`
	StopSequences []string `json:"stopSequences,omitempty"`
}

// ConvertToBedrockFormat converts OpenAI messages to provider-specific format
func ConvertToBedrockFormat(messages []*goopenai.ChatCompletionMessage, modelID string, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	provider := GetModelProvider(modelID)

	switch provider {
	case ProviderAnthropic:
		return convertToAnthropicFormat(messages, modelID, maxTokens, temperature, topP)
	case ProviderMeta:
		return convertToMetaFormat(messages, maxTokens, temperature, topP)
	case ProviderMistral:
		return convertToMistralFormat(messages, maxTokens, temperature, topP)
	case ProviderCohere:
		return convertToCohereFormat(messages, maxTokens, temperature, topP)
	case ProviderAI21:
		return convertToAI21Format(messages, maxTokens, temperature, topP)
	default:
		return nil, fmt.Errorf("unsupported model provider: %s", provider)
	}
}

// convertToAnthropicFormat converts messages to Anthropic's format
func convertToAnthropicFormat(messages []*goopenai.ChatCompletionMessage, modelID string, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	var anthropicMessages []AnthropicMessage
	var systemMessage string
	var systemContent string

	modelInfo, exists := GetModelInfo(modelID)
	if !exists {
		// Default settings for unknown models - be conservative and assume no system message support
		modelInfo = ModelInfo{SupportsSystemMessage: false}
	}
	for _, msg := range messages {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			if modelInfo.SupportsSystemMessage {
				systemMessage = msg.Content
			} else {
				// Store system content to embed in first user message
				systemContent = msg.Content
			}
		case goopenai.ChatMessageRoleUser:
			content := msg.Content
			// If we have system content and this is the first user message, prepend it
			if systemContent != "" && len(anthropicMessages) == 0 {
				content = fmt.Sprintf("System: %s\n\nHuman: %s", systemContent, content)
				systemContent = ""
			}
			anthropicMessages = append(anthropicMessages, AnthropicMessage{
				Role:    "user",
				Content: content,
			})
		case goopenai.ChatMessageRoleAssistant:
			anthropicMessages = append(anthropicMessages, AnthropicMessage{
				Role:    "assistant",
				Content: msg.Content,
			})
		}
	}

	// If no messages were created (only system), create a minimal user message
	if len(anthropicMessages) == 0 {
		content := "Please proceed with the request."
		if systemContent != "" {
			content = fmt.Sprintf("System: %s\n\nHuman: %s", systemContent, content)
		}
		anthropicMessages = append(anthropicMessages, AnthropicMessage{
			Role:    "user",
			Content: content,
		})
	}

	// Ensure conversation starts with user message
	if len(anthropicMessages) > 0 && anthropicMessages[0].Role != "user" {
		anthropicMessages = append([]AnthropicMessage{{Role: "user", Content: "Start the conversation"}}, anthropicMessages...)
	}

	// Ensure alternating user/assistant messages
	anthropicMessages = ensureAlternatingMessages(anthropicMessages)

	request := AnthropicRequest{
		Messages:         anthropicMessages,
		MaxTokens:        maxTokens,
		Temperature:      temperature,
		TopP:             topP,
		AnthropicVersion: "bedrock-2023-05-31",
	}

	if modelInfo.SupportsSystemMessage && systemMessage != "" {
		request.System = systemMessage
	}

	return json.Marshal(request)
}

// ensureAlternatingMessages ensures messages alternate between user and assistant
func ensureAlternatingMessages(messages []AnthropicMessage) []AnthropicMessage {
	if len(messages) == 0 {
		return messages
	}

	result := []AnthropicMessage{messages[0]}
	expectedRole := "assistant"
	if messages[0].Role == "assistant" {
		expectedRole = "user"
	}

	for i := 1; i < len(messages); i++ {
		if messages[i].Role == expectedRole {
			result = append(result, messages[i])
			if expectedRole == "user" {
				expectedRole = "assistant"
			} else {
				expectedRole = "user"
			}
		} else {
			// Same role as previous, merge content
			result[len(result)-1].Content += "\n" + messages[i].Content
		}
	}

	return result
}

// convertToMetaFormat converts messages to Meta's format
func convertToMetaFormat(messages []*goopenai.ChatCompletionMessage, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	// Meta models use a simple prompt format
	var promptBuilder strings.Builder

	// Build conversation prompt
	for _, msg := range messages {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			promptBuilder.WriteString(fmt.Sprintf("<s>[INST] <<SYS>>\n%s\n<</SYS>>\n\n", msg.Content))
		case goopenai.ChatMessageRoleUser:
			promptBuilder.WriteString(fmt.Sprintf("[INST] %s [/INST]", msg.Content))
		case goopenai.ChatMessageRoleAssistant:
			promptBuilder.WriteString(fmt.Sprintf(" %s </s><s>", msg.Content))
		}
	}

	request := MetaRequest{
		Prompt:      promptBuilder.String(),
		MaxGenLen:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	}

	return json.Marshal(request)
}

// convertToMistralFormat converts messages to Mistral's format
func convertToMistralFormat(messages []*goopenai.ChatCompletionMessage, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	// Mistral uses a specific prompt format
	var promptBuilder strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			promptBuilder.WriteString(fmt.Sprintf("<s>[INST] System: %s\n", msg.Content))
		case goopenai.ChatMessageRoleUser:
			promptBuilder.WriteString(fmt.Sprintf("[INST] %s [/INST]", msg.Content))
		case goopenai.ChatMessageRoleAssistant:
			promptBuilder.WriteString(fmt.Sprintf("%s</s>", msg.Content))
		}
	}

	request := MistralRequest{
		Prompt:      promptBuilder.String(),
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	}

	return json.Marshal(request)
}

// convertToCohereFormat converts messages to Cohere's format
func convertToCohereFormat(messages []*goopenai.ChatCompletionMessage, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	// Cohere uses a simple prompt format
	var promptBuilder strings.Builder

	for _, msg := range messages {
		promptBuilder.WriteString(fmt.Sprintf("%s: %s\n", strings.Title(msg.Role), msg.Content))
	}

	request := CohereRequest{
		Prompt:      promptBuilder.String(),
		MaxTokens:   maxTokens,
		Temperature: temperature,
		P:           topP,
	}

	return json.Marshal(request)
}

// convertToAI21Format converts messages to AI21's format
func convertToAI21Format(messages []*goopenai.ChatCompletionMessage, maxTokens int, temperature float32, topP float32) ([]byte, error) {
	// AI21 uses a simple prompt format
	var promptBuilder strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			promptBuilder.WriteString(fmt.Sprintf("Instructions: %s\n\n", msg.Content))
		case goopenai.ChatMessageRoleUser:
			promptBuilder.WriteString(fmt.Sprintf("Human: %s\n", msg.Content))
		case goopenai.ChatMessageRoleAssistant:
			promptBuilder.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		}
	}

	request := AI21Request{
		Prompt:      promptBuilder.String(),
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	}

	return json.Marshal(request)
}