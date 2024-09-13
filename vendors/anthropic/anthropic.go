package anthropic

import (
	"context"
	"errors"
	"fmt"

	"github.com/danielmiessler/fabric/common"
	"github.com/liushuangls/go-anthropic/v2"
)

func NewClient() (ret *Client) {
	vendorName := "Anthropic"
	ret = &Client{}

	ret.Configurable = &common.Configurable{
		Label:           vendorName,
		EnvNamePrefix:   common.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.ApiKey = ret.Configurable.AddSetupQuestion("API key", true)

	// we could provide a setup question for the following settings
	ret.maxTokens = 4096
	ret.defaultRequiredUserMessage = "Hi"
	ret.models = []string{
		anthropic.ModelClaude3Haiku20240307, anthropic.ModelClaude3Opus20240229,
		anthropic.ModelClaude3Opus20240229, anthropic.ModelClaude2Dot0, anthropic.ModelClaude2Dot1,
		anthropic.ModelClaudeInstant1Dot2, "claude-3-5-sonnet-20240620",
	}

	return
}

type Client struct {
	*common.Configurable
	ApiKey *common.SetupQuestion

	maxTokens                  int
	defaultRequiredUserMessage string
	models                     []string

	client *anthropic.Client
}

func (an *Client) configure() (err error) {
	an.client = anthropic.NewClient(an.ApiKey.Value)
	return
}

func (an *Client) ListModels() (ret []string, err error) {
	return an.models, nil
}

func (an *Client) SendStream(
	msgs []*common.Message, opts *common.ChatOptions, channel chan string,
) (err error) {
	ctx := context.Background()
	req := an.buildMessagesRequest(msgs, opts)
	req.Stream = true

	if _, err = an.client.CreateMessagesStream(ctx, anthropic.MessagesStreamRequest{
		MessagesRequest: req,
		OnContentBlockDelta: func(data anthropic.MessagesEventContentBlockDeltaData) {
			// fmt.Printf("Stream Content: %s\n", data.Delta.Text)
			channel <- *data.Delta.Text
		},
	}); err != nil {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages stream error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages stream error: %v\n", err)
		}
	} else {
		close(channel)
	}
	return
}

func (an *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (ret string, err error) {
	req := an.buildMessagesRequest(msgs, opts)
	req.Stream = false

	var resp anthropic.MessagesResponse
	if resp, err = an.client.CreateMessages(ctx, req); err == nil {
		ret = *resp.Content[0].Text
	} else {
		var e *anthropic.APIError
		if errors.As(err, &e) {
			fmt.Printf("Messages error, type: %s, message: %s", e.Type, e.Message)
		} else {
			fmt.Printf("Messages error: %v\n", err)
		}
	}
	return
}

func (an *Client) buildMessagesRequest(msgs []*common.Message, opts *common.ChatOptions) (ret anthropic.MessagesRequest) {
	temperature := float32(opts.Temperature)
	topP := float32(opts.TopP)

	messages := an.toMessages(msgs)

	ret = anthropic.MessagesRequest{
		Model:       opts.Model,
		Temperature: &temperature,
		TopP:        &topP,
		Messages:    messages,
		MaxTokens:   an.maxTokens,
	}
	return
}

func (an *Client) toMessages(msgs []*common.Message) (ret []anthropic.Message) {
	// we could call the method before calling the specific vendor
	normalizedMessages := common.NormalizeMessages(msgs, an.defaultRequiredUserMessage)

	// Iterate over the incoming session messages and process them
	for _, msg := range normalizedMessages {
		var message anthropic.Message
		switch msg.Role {
		case "user":
			message = anthropic.NewUserTextMessage(msg.Content)
		case "system":
			message = anthropic.NewAssistantTextMessage(msg.Content)
		default:
			message = anthropic.NewAssistantTextMessage(msg.Content)
		}
		ret = append(ret, message)
	}
	return
}
