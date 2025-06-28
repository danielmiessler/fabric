package openai

import (
	"context"
	"log/slog"
	"slices"
	"strings"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	openai "github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/openai/openai-go/packages/pagination"
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
	ApiClient  *openai.Client
}

func (o *Client) configure() (ret error) {
	opts := []option.RequestOption{option.WithAPIKey(o.ApiKey.Value)}
	if o.ApiBaseURL.Value != "" {
		opts = append(opts, option.WithBaseURL(o.ApiBaseURL.Value))
	}
	client := openai.NewClient(opts...)
	o.ApiClient = &client
	return
}

func (o *Client) ListModels() (ret []string, err error) {
	var page *pagination.Page[openai.Model]
	if page, err = o.ApiClient.Models.List(context.Background()); err != nil {
		return
	}
	for _, mod := range page.Data {
		ret = append(ret, mod.ID)
	}
	return
}

func (o *Client) SendStream(
	msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {
	req := o.buildChatCompletionParams(msgs, opts)
	stream := o.ApiClient.Chat.Completions.NewStreaming(context.Background(), req)
	for stream.Next() {
		chunk := stream.Current()
		if len(chunk.Choices) > 0 {
			channel <- chunk.Choices[0].Delta.Content
		}
	}
	if stream.Err() == nil {
		channel <- "\n"
	}
	close(channel)
	return stream.Err()
}

func (o *Client) Send(ctx context.Context, msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {
	req := o.buildChatCompletionParams(msgs, opts)

	var resp *openai.ChatCompletion
	if resp, err = o.ApiClient.Chat.Completions.New(ctx, req); err != nil {
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
	openAIModelsNeedingRaw := []string{
		"gpt-4o-mini-search-preview",
		"gpt-4o-mini-search-preview-2025-03-11",
		"gpt-4o-search-preview",
		"gpt-4o-search-preview-2025-03-11",
		"o4-mini-deep-research",
		"o4-mini-deep-research-2025-06-26",
	}
	for _, prefix := range openaiModelsPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			return true
		}
	}
	return slices.Contains(openAIModelsNeedingRaw, modelName)
}

func (o *Client) buildChatCompletionParams(
	inputMsgs []*chat.ChatCompletionMessage, opts *common.ChatOptions,
) (ret openai.ChatCompletionNewParams) {

	// Create a new slice for messages to be sent, converting from []*Msg to []Msg.
	// This also serves as a mutable copy for provider-specific modifications.
	messagesForRequest := make([]openai.ChatCompletionMessageParamUnion, len(inputMsgs))
	for i, msgPtr := range inputMsgs {
		msg := *msgPtr // copy
		// Provider-specific modification for DeepSeek:
		if strings.Contains(opts.Model, "deepseek") && len(inputMsgs) == 1 && msg.Role == chat.ChatMessageRoleSystem {
			msg.Role = chat.ChatMessageRoleUser
		}
		messagesForRequest[i] = convertMessage(msg)
	}
	ret = openai.ChatCompletionNewParams{
		Model:    openai.ChatModel(opts.Model),
		Messages: messagesForRequest,
	}
	if !opts.Raw {
		ret.Temperature = openai.Float(opts.Temperature)
		ret.TopP = openai.Float(opts.TopP)
		ret.PresencePenalty = openai.Float(opts.PresencePenalty)
		ret.FrequencyPenalty = openai.Float(opts.FrequencyPenalty)
		if opts.Seed != 0 {
			ret.Seed = openai.Int(int64(opts.Seed))
		}
	}
	return
}

func convertMessage(msg chat.ChatCompletionMessage) openai.ChatCompletionMessageParamUnion {
	switch msg.Role {
	case chat.ChatMessageRoleSystem:
		return openai.SystemMessage(msg.Content)
	case chat.ChatMessageRoleUser:
		if len(msg.MultiContent) > 0 {
			var parts []openai.ChatCompletionContentPartUnionParam
			for _, p := range msg.MultiContent {
				switch p.Type {
				case chat.ChatMessagePartTypeText:
					parts = append(parts, openai.TextContentPart(p.Text))
				case chat.ChatMessagePartTypeImageURL:
					parts = append(parts, openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{URL: p.ImageURL.URL}))
				}
			}
			return openai.UserMessage(parts)
		}
		return openai.UserMessage(msg.Content)
	default:
		return openai.AssistantMessage(msg.Content)
	}
}
