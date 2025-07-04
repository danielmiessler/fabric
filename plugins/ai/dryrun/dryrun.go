package dryrun

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/danielmiessler/fabric/chat"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
)

type Client struct {
	*plugins.PluginBase
}

func NewClient() *Client {
	return &Client{PluginBase: &plugins.PluginBase{Name: "DryRun"}}
}

func (c *Client) ListModels() ([]string, error) {
	return []string{"dry-run-model"}, nil
}

func (c *Client) formatMultiContentMessage(msg *chat.ChatCompletionMessage) string {
	var builder strings.Builder

	if len(msg.MultiContent) > 0 {
		builder.WriteString(fmt.Sprintf("%s:\n", msg.Role))
		for _, part := range msg.MultiContent {
			builder.WriteString(fmt.Sprintf("  - Type: %s\n", part.Type))
			if part.Type == chat.ChatMessagePartTypeImageURL {
				builder.WriteString(fmt.Sprintf("    Image URL: %s\n", part.ImageURL.URL))
			} else {
				builder.WriteString(fmt.Sprintf("    Text: %s\n", part.Text))
			}
		}
		builder.WriteString("\n")
	} else {
		builder.WriteString(fmt.Sprintf("%s:\n%s\n\n", msg.Role, msg.Content))
	}

	return builder.String()
}

func (c *Client) formatMessages(msgs []*chat.ChatCompletionMessage) string {
	var builder strings.Builder

	for _, msg := range msgs {
		switch msg.Role {
		case chat.ChatMessageRoleSystem:
			builder.WriteString(fmt.Sprintf("System:\n%s\n\n", msg.Content))
		case chat.ChatMessageRoleAssistant:
			builder.WriteString(c.formatMultiContentMessage(msg))
		case chat.ChatMessageRoleUser:
			builder.WriteString(c.formatMultiContentMessage(msg))
		default:
			builder.WriteString(fmt.Sprintf("%s:\n%s\n\n", msg.Role, msg.Content))
		}
	}

	return builder.String()
}

func (c *Client) formatOptions(opts *common.ChatOptions) string {
	var builder strings.Builder

	builder.WriteString("Options:\n")
	builder.WriteString(fmt.Sprintf("Model: %s\n", opts.Model))
	builder.WriteString(fmt.Sprintf("Temperature: %f\n", opts.Temperature))
	builder.WriteString(fmt.Sprintf("TopP: %f\n", opts.TopP))
	builder.WriteString(fmt.Sprintf("PresencePenalty: %f\n", opts.PresencePenalty))
	builder.WriteString(fmt.Sprintf("FrequencyPenalty: %f\n", opts.FrequencyPenalty))
	if opts.ModelContextLength != 0 {
		builder.WriteString(fmt.Sprintf("ModelContextLength: %d\n", opts.ModelContextLength))
	}
	if opts.Search {
		builder.WriteString("Search: enabled\n")
		if opts.SearchLocation != "" {
			builder.WriteString(fmt.Sprintf("SearchLocation: %s\n", opts.SearchLocation))
		}
	}

	return builder.String()
}

func (c *Client) SendStream(msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions, channel chan string) error {
	var builder strings.Builder
	builder.WriteString("Dry run: Would send the following request:\n\n")
	builder.WriteString(c.formatMessages(msgs))
	builder.WriteString(c.formatOptions(opts))

	channel <- builder.String()
	close(channel)
	return nil
}

func (c *Client) Send(_ context.Context, msgs []*chat.ChatCompletionMessage, opts *common.ChatOptions) (string, error) {
	fmt.Println("Dry run: Would send the following request:")
	fmt.Print(c.formatMessages(msgs))
	fmt.Print(c.formatOptions(opts))

	return "", nil
}

func (c *Client) Setup() error {
	return nil
}

func (c *Client) SetupFillEnvFileContent(_ *bytes.Buffer) {
	// No environment variables needed for dry run
}

func (c *Client) NeedsRawMode(modelName string) bool {
	return false
}
