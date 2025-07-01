package openai

import (
	"testing"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
	"github.com/stretchr/testify/assert"
)

func TestBuildResponseRequestWithMaxTokens(t *testing.T) {

	var msgs []*chat.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &chat.ChatCompletionMessage{
			Role:    "User",
			Content: "My msg",
		})
	}

	opts := &common.ChatOptions{
		Temperature: 0.8,
		TopP:        0.9,
		Raw:         false,
		MaxTokens:   50,
	}

	var client = NewClient()
	request := client.buildResponseParams(msgs, opts)
	assert.Equal(t, shared.ResponsesModel(opts.Model), request.Model)
	assert.Equal(t, openai.Float(opts.Temperature), request.Temperature)
	assert.Equal(t, openai.Float(opts.TopP), request.TopP)
	assert.Equal(t, openai.Int(int64(opts.MaxTokens)), request.MaxOutputTokens)
}

func TestBuildResponseRequestNoMaxTokens(t *testing.T) {

	var msgs []*chat.ChatCompletionMessage

	for i := 0; i < 2; i++ {
		msgs = append(msgs, &chat.ChatCompletionMessage{
			Role:    "User",
			Content: "My msg",
		})
	}

	opts := &common.ChatOptions{
		Temperature: 0.8,
		TopP:        0.9,
		Raw:         false,
	}

	var client = NewClient()
	request := client.buildResponseParams(msgs, opts)
	assert.Equal(t, shared.ResponsesModel(opts.Model), request.Model)
	assert.Equal(t, openai.Float(opts.Temperature), request.Temperature)
	assert.Equal(t, openai.Float(opts.TopP), request.TopP)
	assert.False(t, request.MaxOutputTokens.Valid())
}
