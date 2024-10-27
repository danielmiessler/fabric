package dryrun

import (
	"bytes"
	"context"
	"fmt"
	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
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

func (c *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) error {
	output := "Dry run: Would send the following request:\n\n"

	for _, msg := range msgs {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			output += fmt.Sprintf("System:\n%s\n\n", msg.Content)
		case goopenai.ChatMessageRoleAssistant:
			output += fmt.Sprintf("Assistant:\n%s\n\n", msg.Content)
		case goopenai.ChatMessageRoleUser:
			output += fmt.Sprintf("User:\n%s\n\n", msg.Content)
		default:
			output += fmt.Sprintf("%s:\n%s\n\n", msg.Role, msg.Content)
		}
	}

	output += "Options:\n"
	output += fmt.Sprintf("Model: %s\n", opts.Model)
	output += fmt.Sprintf("Temperature: %f\n", opts.Temperature)
	output += fmt.Sprintf("TopP: %f\n", opts.TopP)
	output += fmt.Sprintf("PresencePenalty: %f\n", opts.PresencePenalty)
	output += fmt.Sprintf("FrequencyPenalty: %f\n", opts.FrequencyPenalty)

	channel <- output
	close(channel)
	return nil
}

func (c *Client) Send(_ context.Context, msgs []*common.Message, opts *common.ChatOptions) (string, error) {
	fmt.Println("Dry run: Would send the following request:")

	for _, msg := range msgs {
		switch msg.Role {
		case goopenai.ChatMessageRoleSystem:
			fmt.Printf("System:\n%s\n\n", msg.Content)
		case goopenai.ChatMessageRoleAssistant:
			fmt.Printf("Assistant:\n%s\n\n", msg.Content)
		case goopenai.ChatMessageRoleUser:
			fmt.Printf("User:\n%s\n\n", msg.Content)
		default:
			fmt.Printf("%s:\n%s\n\n", msg.Role, msg.Content)
		}
	}

	fmt.Println("Options:")
	fmt.Printf("Model: %s\n", opts.Model)
	fmt.Printf("Temperature: %f\n", opts.Temperature)
	fmt.Printf("TopP: %f\n", opts.TopP)
	fmt.Printf("PresencePenalty: %f\n", opts.PresencePenalty)
	fmt.Printf("FrequencyPenalty: %f\n", opts.FrequencyPenalty)

	return "", nil
}

func (c *Client) Setup() error {
	return nil
}

func (c *Client) SetupFillEnvFileContent(_ *bytes.Buffer) {
	// No environment variables needed for dry run
}
