package openai

// This file contains helper methods for the Chat Completions API.
// These methods are used as fallbacks for OpenAI-compatible providers
// that don't support the newer Responses API (e.g., Groq, Mistral, etc.).

import (
	"context"
	"strings"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

// sendChatCompletions sends a request using the Chat Completions API
func (o *Client) sendChatCompletions(ctx context.Context, msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {
	req := o.buildChatCompletionParams(msgs, opts)

	var resp *openai.ChatCompletion
	if resp, err = o.ApiClient.Chat.Completions.New(ctx, req); err != nil {
		return
	}
	if len(resp.Choices) > 0 {
		ret = resp.Choices[0].Message.Content
	}
	return
}

// sendStreamChatCompletions sends a streaming request using the Chat Completions API
func (o *Client) sendStreamChatCompletions(
	msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {
	defer close(channel)

	req := o.buildChatCompletionParams(msgs, opts)
	stream := o.ApiClient.Chat.Completions.NewStreaming(context.Background(), req)
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 && chunk.Choices[0].Delta.Content != "" {
			channel <- chunk.Choices[0].Delta.Content
		}
	}
	if stream.Err() == nil {
		channel <- "\n"
	}
	return stream.Err()
}

// buildChatCompletionParams builds parameters for the Chat Completions API
func (o *Client) buildChatCompletionParams(
	inputMsgs []*chat.ChatCompletionMessage, opts *common.ChatOptions,
) (ret openai.ChatCompletionNewParams) {

	messages := make([]openai.ChatCompletionMessageParamUnion, len(inputMsgs))
	for i, msgPtr := range inputMsgs {
		msg := *msgPtr
		if strings.Contains(opts.Model, "deepseek") && len(inputMsgs) == 1 && msg.Role == chat.ChatMessageRoleSystem {
			msg.Role = chat.ChatMessageRoleUser
		}
		messages[i] = o.convertChatMessage(msg)
	}

	ret = openai.ChatCompletionNewParams{
		Model:    shared.ChatModel(opts.Model),
		Messages: messages,
	}

	if !opts.Raw {
		ret.Temperature = openai.Float(opts.Temperature)
		ret.TopP = openai.Float(opts.TopP)
		if opts.MaxTokens != 0 {
			ret.MaxTokens = openai.Int(int64(opts.MaxTokens))
		}
		if opts.PresencePenalty != 0 {
			ret.PresencePenalty = openai.Float(opts.PresencePenalty)
		}
		if opts.FrequencyPenalty != 0 {
			ret.FrequencyPenalty = openai.Float(opts.FrequencyPenalty)
		}
		if opts.Seed != 0 {
			ret.Seed = openai.Int(int64(opts.Seed))
		}
	}
	return
}

// convertChatMessage converts fabric chat message to OpenAI chat completion message
func (o *Client) convertChatMessage(msg chat.ChatCompletionMessage) openai.ChatCompletionMessageParamUnion {
	// For now, simplify to text-only messages to get the basic functionality working
	// Multi-content support can be added later if needed
	switch msg.Role {
	case chat.ChatMessageRoleSystem:
		return openai.SystemMessage(msg.Content)
	case chat.ChatMessageRoleUser:
		return openai.UserMessage(msg.Content)
	case chat.ChatMessageRoleAssistant:
		return openai.AssistantMessage(msg.Content)
	default:
		return openai.UserMessage(msg.Content)
	}
}
