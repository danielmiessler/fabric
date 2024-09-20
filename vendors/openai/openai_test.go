package openai

import (
	"testing"

	"github.com/danielmiessler/fabric/common"
	"github.com/sashabaranov/go-openai"
	goopenai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestBuildChatCompletionRequestPinSeed(t *testing.T) {

	var msgs []*common.Message

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &common.Message{
			Role:    "User",
			Content: "My msg",
		})
	}

	opts := &common.ChatOptions{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
		Raw:              false,
		Seed:             1,
	}

	var expectedMessages []openai.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		expectedMessages = append(expectedMessages,
			openai.ChatCompletionMessage{
				Role:    msgs[i].Role,
				Content: msgs[i].Content,
			},
		)
	}

	var expectedRequest = goopenai.ChatCompletionRequest{
		Model:            opts.Model,
		Temperature:      float32(opts.Temperature),
		TopP:             float32(opts.TopP),
		PresencePenalty:  float32(opts.PresencePenalty),
		FrequencyPenalty: float32(opts.FrequencyPenalty),
		Messages:         expectedMessages,
		Seed:             &opts.Seed,
	}

	var client = NewClient()
	request := client.buildChatCompletionRequest(msgs, opts)
	assert.Equal(t, expectedRequest, request)
}

func TestBuildChatCompletionRequestNilSeed(t *testing.T) {

	var msgs []*common.Message

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &common.Message{
			Role:    "User",
			Content: "My msg",
		})
	}

	opts := &common.ChatOptions{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
		Raw:              false,
		Seed:             0,
	}

	var expectedMessages []openai.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		expectedMessages = append(expectedMessages,
			openai.ChatCompletionMessage{
				Role:    msgs[i].Role,
				Content: msgs[i].Content,
			},
		)
	}

	var expectedRequest = goopenai.ChatCompletionRequest{
		Model:            opts.Model,
		Temperature:      float32(opts.Temperature),
		TopP:             float32(opts.TopP),
		PresencePenalty:  float32(opts.PresencePenalty),
		FrequencyPenalty: float32(opts.FrequencyPenalty),
		Messages:         expectedMessages,
		Seed:             nil,
	}

	var client = NewClient()
	request := client.buildChatCompletionRequest(msgs, opts)
	assert.Equal(t, expectedRequest, request)
}
