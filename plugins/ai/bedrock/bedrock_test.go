package bedrock

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	goopenai "github.com/sashabaranov/go-openai"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	assert.NotNil(t, client)
	assert.Equal(t, "Bedrock", client.Name)
	assert.NotNil(t, client.AwsAccessKeyId)
	assert.NotNil(t, client.AwsSecretAccessKey)
	assert.NotNil(t, client.AwsRegion)
	assert.NotNil(t, client.AwsSessionToken)
	assert.NotNil(t, client.AwsProfile)
	assert.NotNil(t, client.PreferredModels)
}

func TestGetModelProvider(t *testing.T) {
	tests := []struct {
		modelID  string
		expected ModelProvider
	}{
		{"anthropic.claude-3-opus-20240229-v1:0", ProviderAnthropic},
		{"meta.llama3-70b-instruct-v1:0", ProviderMeta},
		{"mistral.mistral-7b-instruct-v0:2", ProviderMistral},
		{"cohere.command-text-v14", ProviderCohere},
		{"ai21.j2-ultra-v1", ProviderAI21},
		{"unknown.model", "unknown"},
	}

	for _, test := range tests {
		t.Run(test.modelID, func(t *testing.T) {
			provider := GetModelProvider(test.modelID)
			assert.Equal(t, test.expected, provider)
		})
	}
}

func TestGetModelInfo(t *testing.T) {
	// Test existing model
	info, exists := GetModelInfo("anthropic.claude-3-opus-20240229-v1:0")
	assert.True(t, exists)
	assert.Equal(t, "Claude 3 Opus", info.DisplayName)
	assert.Equal(t, 4096, info.MaxTokens)
	assert.True(t, info.SupportsStreaming)
	assert.True(t, info.SupportsSystemMessage)

	// Test non-existing model
	_, exists = GetModelInfo("non.existing.model")
	assert.False(t, exists)
}

func TestEnsureAlternatingMessages(t *testing.T) {
	messages := []AnthropicMessage{
		{Role: "user", Content: "Hello"},
		{Role: "user", Content: "How are you?"},
		{Role: "assistant", Content: "I'm fine"},
		{Role: "assistant", Content: "Thanks for asking"},
		{Role: "user", Content: "Great!"},
	}

	result := ensureAlternatingMessages(messages)

	assert.Len(t, result, 3)
	assert.Equal(t, "user", result[0].Role)
	assert.Equal(t, "Hello\nHow are you?", result[0].Content)
	assert.Equal(t, "assistant", result[1].Role)
	assert.Equal(t, "I'm fine\nThanks for asking", result[1].Content)
	assert.Equal(t, "user", result[2].Role)
	assert.Equal(t, "Great!", result[2].Content)
}

func TestConvertToAnthropicFormat(t *testing.T) {
	messages := []*goopenai.ChatCompletionMessage{
		{Role: goopenai.ChatMessageRoleSystem, Content: "You are a helpful assistant"},
		{Role: goopenai.ChatMessageRoleUser, Content: "Hello"},
		{Role: goopenai.ChatMessageRoleAssistant, Content: "Hi there!"},
	}

	// Test with Claude 3 (supports system messages)
	data, err := convertToAnthropicFormat(messages, "anthropic.claude-3-opus-20240229-v1:0", 1000, 0.7, 0.9)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var request AnthropicRequest
	err = json.Unmarshal(data, &request)
	assert.NoError(t, err)
	assert.Equal(t, "You are a helpful assistant", request.System)
	assert.Len(t, request.Messages, 2)

	// Test with Claude 2 (doesn't support system messages)
	data, err = convertToAnthropicFormat(messages, "anthropic.claude-v2", 1000, 0.7, 0.9)
	assert.NoError(t, err)
	assert.NotNil(t, data)

	var request2 AnthropicRequest
	err = json.Unmarshal(data, &request2)
	assert.NoError(t, err)
	assert.Empty(t, request2.System)
	assert.Len(t, request2.Messages, 2)
	assert.Contains(t, request2.Messages[0].Content, "System:")
}