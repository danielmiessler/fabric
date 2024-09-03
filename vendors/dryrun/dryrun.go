package dryrun

import (
	"bytes"
	"context"
	"fmt"

	"github.com/danielmiessler/fabric/common"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) GetName() string {
	return "DryRun"
}

func (c *Client) IsConfigured() bool {
	return true
}

func (c *Client) Configure() error {
	return nil
}

func (c *Client) ListModels() ([]string, error) {
	return []string{"dry-run-model"}, nil
}

func (c *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) error {
	output := "Dry run: Would send the following request:\n\n"

	for _, msg := range msgs {
		switch msg.Role {
		case "system":
			output += fmt.Sprintf("System:\n%s\n\n", msg.Content)
		case "user":
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

func (c *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (string, error) {
	fmt.Println("Dry run: Would send the following request:")

	for _, msg := range msgs {
		switch msg.Role {
		case "system":
			fmt.Printf("System:\n%s\n\n", msg.Content)
		case "user":
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

func (c *Client) SetupFillEnvFileContent(buffer *bytes.Buffer) {
	// No environment variables needed for dry run
}
