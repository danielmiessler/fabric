package openai

import (
	"testing"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	openai "github.com/openai/openai-go"
	"github.com/stretchr/testify/assert"
)

func TestBuildChatCompletionRequestPinSeed(t *testing.T) {

	var msgs []*chat.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &chat.ChatCompletionMessage{
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

	var client = NewClient()
	request := client.buildChatCompletionParams(msgs, opts)
	assert.Equal(t, openai.ChatModel(opts.Model), request.Model)
	assert.Equal(t, openai.Float(opts.Temperature), request.Temperature)
	assert.Equal(t, openai.Float(opts.TopP), request.TopP)
	assert.Equal(t, openai.Float(opts.PresencePenalty), request.PresencePenalty)
	assert.Equal(t, openai.Float(opts.FrequencyPenalty), request.FrequencyPenalty)
	assert.Equal(t, openai.Int(int64(opts.Seed)), request.Seed)
}

func TestBuildChatCompletionRequestNilSeed(t *testing.T) {

	var msgs []*chat.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &chat.ChatCompletionMessage{
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

	var client = NewClient()
	request := client.buildChatCompletionParams(msgs, opts)
	assert.Equal(t, openai.ChatModel(opts.Model), request.Model)
	assert.Equal(t, openai.Float(opts.Temperature), request.Temperature)
	assert.Equal(t, openai.Float(opts.TopP), request.TopP)
	assert.Equal(t, openai.Float(opts.PresencePenalty), request.PresencePenalty)
	assert.Equal(t, openai.Float(opts.FrequencyPenalty), request.FrequencyPenalty)
	assert.False(t, request.Seed.Valid())
}
