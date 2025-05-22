package anthropic

import (
	"context"
	"fmt"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"
)

const defaultBaseUrl = "https://api.anthropic.com/"

func NewClient() (ret *Client) {
	vendorName := "Anthropic"
	ret = &Client{}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", false)
	ret.ApiBaseURL.Value = defaultBaseUrl
	ret.ApiKey = ret.PluginBase.AddSetupQuestion("API key", true)

	ret.maxTokens = 4096
	ret.defaultRequiredUserMessage = "Hi"
	ret.models = []string{
		string(anthropic.ModelClaude3_7SonnetLatest), string(anthropic.ModelClaude3_7Sonnet20250219),
		string(anthropic.ModelClaude3_5HaikuLatest), string(anthropic.ModelClaude3_5Haiku20241022),
		string(anthropic.ModelClaude3_5SonnetLatest), string(anthropic.ModelClaude3_5Sonnet20241022),
		string(anthropic.ModelClaude_3_5_Sonnet_20240620), string(anthropic.ModelClaude3OpusLatest),
		string(anthropic.ModelClaude_3_Opus_20240229), string(anthropic.ModelClaude_3_Haiku_20240307),
		string(anthropic.ModelClaudeOpus4_20250514), string(anthropic.ModelClaudeSonnet4_20250514),
	}

	return
}

type Client struct {
	*plugins.PluginBase
	ApiBaseURL *plugins.SetupQuestion
	ApiKey     *plugins.SetupQuestion

	maxTokens                  int
	defaultRequiredUserMessage string
	models                     []string

	client anthropic.Client
}

func (an *Client) configure() (err error) {
	if an.ApiBaseURL.Value != "" {
		baseURL := an.ApiBaseURL.Value

		// As of 2.0beta1, using v2 API endpoint.
		// https://github.com/anthropics/anthropic-sdk-go/blob/main/CHANGELOG.md#020-beta1-2025-03-25
		if strings.Contains(baseURL, "-") && !strings.HasSuffix(baseURL, "/v2") {
			baseURL = strings.TrimSuffix(baseURL, "/")
			baseURL = baseURL + "/v2"
		}

		an.client = anthropic.NewClient(
			option.WithAPIKey(an.ApiKey.Value),
			option.WithBaseURL(baseURL),
		)
	} else {
		an.client = anthropic.NewClient(option.WithAPIKey(an.ApiKey.Value))
	}
	return
}

func (an *Client) ListModels() (ret []string, err error) {
	return an.models, nil
}

func (an *Client) SendStream(
	msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string,
) (err error) {
	messages := an.toMessages(msgs)
	if len(messages) == 0 {
		close(channel)
		// No messages to send after normalization, consider this a non-error condition for streaming.
		return nil
	}

	ctx := context.Background()
	stream := an.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:       anthropic.Model(opts.Model),
		MaxTokens:   int64(an.maxTokens),
		TopP:        anthropic.Opt(opts.TopP),
		Temperature: anthropic.Opt(opts.Temperature),
		Messages:    messages,
	})

	for stream.Next() {
		event := stream.Current()

		// directly send any non-empty delta text
		if event.Delta.Text != "" {
			channel <- event.Delta.Text
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
	if len(messages) == 0 {
		// No messages to send after normalization, return empty string and no error.
		return "", nil
	}

	var message *anthropic.Message
	if message, err = an.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:       anthropic.Model(opts.Model),
		MaxTokens:   int64(an.maxTokens),
		TopP:        anthropic.Opt(opts.TopP),
		Temperature: anthropic.Opt(opts.Temperature),
		Messages:    messages,
	}); err != nil {
		return
	}

	if len(message.Content) == 0 {
		// Model returned no content blocks.
		return "", nil
	}
	ret = message.Content[0].Text
	return
}

func (an *Client) toMessages(msgs []*goopenai.ChatCompletionMessage) (ret []anthropic.MessageParam) {
	// Custom normalization for Anthropic:
	// - System messages become the first part of the first user message.
	// - Messages must alternate user/assistant.
	// - Skip empty messages.

	var anthropicMessages []anthropic.MessageParam
	var systemContent string
	isFirstUserMessage := true
	lastRoleWasUser := false

	for _, msg := range msgs {
		if msg.Content == "" {
			continue // Skip empty messages
		}

		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			// Accumulate system content. It will be prepended to the first user message.
			if systemContent != "" {
				systemContent += "\\n" + msg.Content
			} else {
				systemContent = msg.Content
			}
		case goopenai.ChatMessageRoleUser:
			userContent := msg.Content
			if isFirstUserMessage && systemContent != "" {
				userContent = systemContent + "\\n\\n" + userContent
				isFirstUserMessage = false // System content now consumed
			}
			if lastRoleWasUser {
				// Enforce alternation: add a minimal assistant message if two user messages are consecutive.
				// This shouldn't happen with current chatter.go logic but is a safeguard.
				anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock("Okay.")))
			}
			anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(userContent)))
			lastRoleWasUser = true
		case goopenai.ChatMessageRoleAssistant:
			// If the first message is an assistant message, and we have system content,
			// prepend a user message with the system content.
			if isFirstUserMessage && systemContent != "" {
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(systemContent)))
				lastRoleWasUser = true
				isFirstUserMessage = false // System content now consumed
			} else if !lastRoleWasUser && len(anthropicMessages) > 0 {
				// Enforce alternation: add a minimal user message if two assistant messages are consecutive
				// or if an assistant message is first without prior system prompt handling.
				anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(an.defaultRequiredUserMessage)))
				lastRoleWasUser = true
			}
			anthropicMessages = append(anthropicMessages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
			lastRoleWasUser = false
		default:
			// Other roles (like 'meta') are ignored for Anthropic's message structure.
			continue
		}
	}

	// If only system content was provided, create a user message with it.
	if len(anthropicMessages) == 0 && systemContent != "" {
		anthropicMessages = append(anthropicMessages, anthropic.NewUserMessage(anthropic.NewTextBlock(systemContent)))
	}

	return anthropicMessages
}

func (an *Client) NeedsRawMode(modelName string) bool {
	return false
}
