package anthropic

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"
)

//const baseUrl = "https://api.anthropic.com/"

func NewClient() (ret *Client) {
	vendorName := "Anthropic"
	ret = &Client{}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	//ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", false)
	//ret.ApiBaseURL.Value = baseUrl
	ret.ApiKey = ret.PluginBase.AddSetupQuestion("API key", true)

	// we could provide a setup question for the following settings
	ret.maxTokens = 4096
	ret.defaultRequiredUserMessage = "Hi"
	ret.models = []string{
		anthropic.ModelClaude3_5HaikuLatest, anthropic.ModelClaude3_5Haiku20241022,
		anthropic.ModelClaude3_5SonnetLatest, anthropic.ModelClaude3_5Sonnet20241022,
		anthropic.ModelClaude_3_5_Sonnet_20240620, anthropic.ModelClaude3OpusLatest,
		anthropic.ModelClaude_3_Opus_20240229, anthropic.ModelClaude_3_Sonnet_20240229,
		anthropic.ModelClaude_3_Haiku_20240307, anthropic.ModelClaude_2_1,
		anthropic.ModelClaude_2_0, anthropic.ModelClaude_Instant_1_2,
	}

	return
}

type Client struct {
	*plugins.PluginBase
	//ApiBaseURL *plugins.SetupQuestion
	ApiKey *plugins.SetupQuestion

	maxTokens                  int
	defaultRequiredUserMessage string
	models                     []string

	client *anthropic.Client
}

func (an *Client) configure() (err error) {
	/*if an.ApiBaseURL.Value != "" {
		an.client = anthropic.NewClient(
			option.WithAPIKey(an.ApiKey.Value), option.WithBaseURL(an.ApiBaseURL.Value),
		)
	} else {
	*/
	an.client = anthropic.NewClient(option.WithAPIKey(an.ApiKey.Value))
	//}
	return
}

func (an *Client) ListModels() (ret []string, err error) {
	return an.models, nil
}

func (an *Client) SendStream(
	msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {

	messages := an.toMessages(msgs)

	ctx := context.Background()
	stream := an.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:       anthropic.F(opts.Model),
		MaxTokens:   anthropic.F(int64(an.maxTokens)),
		TopP:        anthropic.F(opts.TopP),
		Temperature: anthropic.F(opts.Temperature),
		Messages:    anthropic.F(messages),
	})

	for stream.Next() {
		event := stream.Current()

		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				channel <- delta.Text
			}
		}
	}

	if stream.Err() != nil {
		fmt.Printf("Messages stream error: %v\n", stream.Err())
	}
	close(channel)
	return
}

func (an *Client) Send(ctx context.Context, msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {
	messages := an.toMessages(msgs)

	var message *anthropic.Message
	if message, err = an.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.F(opts.Model),
		MaxTokens:   anthropic.F(int64(an.maxTokens)),
		TopP:        anthropic.F(opts.TopP),
		Temperature: anthropic.F(opts.Temperature),
		Messages:    anthropic.F(messages),
	}); err != nil {
		return
	}
	ret = message.Content[0].Text
	return
}

func (an *Client) toMessages(msgs []*goopenai.ChatCompletionMessage) (ret []anthropic.MessageParam) {
	// we could call the method before calling the specific vendor
	normalizedMessages := common.NormalizeMessages(msgs, an.defaultRequiredUserMessage)

	// Iterate over the incoming session messages and process them
	for _, msg := range normalizedMessages {
		var message anthropic.MessageParam
		switch msg.Role {
		case goopenai.ChatMessageRoleUser:
			message = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
		default:
			message = anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content))
		}
		ret = append(ret, message)
	}
	return
}
