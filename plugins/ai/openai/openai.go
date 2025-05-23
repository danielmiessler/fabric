package openai

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"strings"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"
)

func NewClient() (ret *Client) {
	return NewClientCompatible("OpenAI", "https://api.openai.com/v1", nil)
}

func NewClientCompatible(vendorName string, defaultBaseUrl string, configureCustom func() error) (ret *Client) {
	ret = NewClientCompatibleNoSetupQuestions(vendorName, configureCustom)

	ret.ApiKey = ret.AddSetupQuestion("API Key", true)
	ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", false)
	ret.ApiBaseURL.Value = defaultBaseUrl

	return
}

func NewClientCompatibleNoSetupQuestions(vendorName string, configureCustom func() error) (ret *Client) {
	ret = &Client{}

	if configureCustom == nil {
		configureCustom = ret.configure
	}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: configureCustom,
	}

	return
}

type Client struct {
	*plugins.PluginBase
	ApiKey     *plugins.SetupQuestion
	ApiBaseURL *plugins.SetupQuestion
	ApiClient  *goopenai.Client
}

func (o *Client) configure() (ret error) {
	config := goopenai.DefaultConfig(o.ApiKey.Value)
	if o.ApiBaseURL.Value != "" {
		config.BaseURL = o.ApiBaseURL.Value
	}
	o.ApiClient = goopenai.NewClientWithConfig(config)
	return
}

func (o *Client) ListModels() (ret []string, err error) {
	var models goopenai.ModelsList
	if models, err = o.ApiClient.ListModels(context.Background()); err != nil {
		return
	}

	model := models.Models
	for _, mod := range model {
		ret = append(ret, mod.ID)
	}
	return
}

func (o *Client) SendStream(
	msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {
	req := o.buildChatCompletionRequest(msgs, opts)
	req.Stream = true

	var stream *goopenai.ChatCompletionStream
	if stream, err = o.ApiClient.CreateChatCompletionStream(context.Background(), req); err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}

	defer stream.Close()

	for {
		var response goopenai.ChatCompletionStreamResponse
		if response, err = stream.Recv(); err == nil {
			if len(response.Choices) > 0 {
				channel <- response.Choices[0].Delta.Content
			} else {
				channel <- "\n"
				close(channel)
				break
			}
		} else if errors.Is(err, io.EOF) {
			channel <- "\n"
			close(channel)
			err = nil
			break
		} else if err != nil {
			fmt.Printf("\nStream error: %v\n", err)
			break
		}
	}
	return
}

func (o *Client) Send(ctx context.Context, msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {
	req := o.buildChatCompletionRequest(msgs, opts)

	var resp goopenai.ChatCompletionResponse
	if resp, err = o.ApiClient.CreateChatCompletion(ctx, req); err != nil {
		return
	}
	if len(resp.Choices) > 0 {
		ret = resp.Choices[0].Message.Content
		slog.Debug("SystemFingerprint: " + resp.SystemFingerprint)
	}
	return
}

func (o *Client) NeedsRawMode(modelName string) bool {
	openaiModelsPrefixes := []string{
		"o1",
		"o3",
		"o4",
	}
	for _, prefix := range openaiModelsPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			return true
		}
	}
	return false
}

func (o *Client) buildChatCompletionRequest(
	inputMsgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions,
) (ret goopenai.ChatCompletionRequest) {

	// Create a new slice for messages to be sent, converting from []*Msg to []Msg.
	// This also serves as a mutable copy for provider-specific modifications.
	messagesForRequest := make([]goopenai.ChatCompletionMessage, len(inputMsgs))
	for i, msgPtr := range inputMsgs {
		messagesForRequest[i] = *msgPtr // Dereference and copy
	}

	// Provider-specific modification for DeepSeek:
	// DeepSeek requires the last message to be a user message.
	// If fabric constructs a single system message (common when a pattern includes user input),
	// we change its role to user for DeepSeek.
	if strings.Contains(opts.Model, "deepseek") { // Heuristic to identify DeepSeek models
		if len(messagesForRequest) == 1 && messagesForRequest[0].Role == goopenai.ChatMessageRoleSystem {
			messagesForRequest[0].Role = goopenai.ChatMessageRoleUser
		}
		// Note: This handles the most common case arising from pattern usage.
		// More complex scenarios where a multi-message sequence ends in 'system'
		// are not currently expected from chatter.go's BuildSession logic for OpenAI providers
		// but might require further rules if they arise.
	}

	if opts.Raw {
		ret = goopenai.ChatCompletionRequest{
			Model:    opts.Model,
			Messages: messagesForRequest,
		}
	} else {
		if opts.Seed == 0 {
			ret = goopenai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messagesForRequest,
			}
		} else {
			ret = goopenai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messagesForRequest,
				Seed:             &opts.Seed,
			}
		}
	}
	return
}
