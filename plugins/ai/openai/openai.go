package openai

import (
	"context"
	"errors"
	"fmt"
	"github.com/danielmiessler/fabric/plugins"
	"io"
	"log/slog"

	"github.com/danielmiessler/fabric/common"
	"github.com/samber/lo"
	"github.com/sashabaranov/go-openai"
	goopenai "github.com/sashabaranov/go-openai"
)

func NewClient() (ret *Client) {
	return NewClientCompatible("OpenAI", "https://api.openai.com/v1", nil)
}

func NewClientCompatible(vendorName string, defaultBaseUrl string, configureCustom func() error) (ret *Client) {
	ret = &Client{}

	if configureCustom == nil {
		configureCustom = ret.configure
	}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: configureCustom,
	}

	ret.ApiKey = ret.AddSetupQuestion("API Key", true)
	ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", false)
	ret.ApiBaseURL.Value = defaultBaseUrl

	return
}

type Client struct {
	*plugins.PluginBase
	ApiKey     *plugins.SetupQuestion
	ApiBaseURL *plugins.SetupQuestion
	ApiClient  *openai.Client
}

func (o *Client) configure() (ret error) {
	config := openai.DefaultConfig(o.ApiKey.Value)
	if o.ApiBaseURL.Value != "" {
		config.BaseURL = o.ApiBaseURL.Value
	}
	o.ApiClient = openai.NewClientWithConfig(config)
	return
}

func (o *Client) ListModels() (ret []string, err error) {
	var models openai.ModelsList
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
	msgs []*common.Message, opts *common.ChatOptions, channel chan string,
) (err error) {
	req := o.buildChatCompletionRequest(msgs, opts)
	req.Stream = true

	var stream *openai.ChatCompletionStream
	if stream, err = o.ApiClient.CreateChatCompletionStream(context.Background(), req); err != nil {
		fmt.Printf("ChatCompletionStream error: %v\n", err)
		return
	}

	defer stream.Close()

	for {
		var response openai.ChatCompletionStreamResponse
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

func (o *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (ret string, err error) {
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

func (o *Client) buildChatCompletionRequest(
	msgs []*common.Message, opts *common.ChatOptions,
) (ret goopenai.ChatCompletionRequest) {
	messages := lo.Map(msgs, func(message *common.Message, _ int) goopenai.ChatCompletionMessage {
		return goopenai.ChatCompletionMessage{Role: message.Role, Content: message.Content}
	})

	if opts.Raw {
		ret = goopenai.ChatCompletionRequest{
			Model:    opts.Model,
			Messages: messages,
		}
	} else {
		if opts.Seed == 0 {
			ret = goopenai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messages,
			}
		} else {
			ret = goopenai.ChatCompletionRequest{
				Model:            opts.Model,
				Temperature:      float32(opts.Temperature),
				TopP:             float32(opts.TopP),
				PresencePenalty:  float32(opts.PresencePenalty),
				FrequencyPenalty: float32(opts.FrequencyPenalty),
				Messages:         messages,
				Seed:             &opts.Seed,
			}
		}
	}
	return
}
