package dryrun

import (
	"bytes"
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

func (c *Client) SendStream(messages []*common.Message, options *common.ChatOptions, channel chan string) error {
	output := "Dry run: Would send the following request:\n\n"

	for _, msg := range messages {
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
	output += fmt.Sprintf("Model: %s\n", options.Model)
	output += fmt.Sprintf("Temperature: %f\n", options.Temperature)
	output += fmt.Sprintf("TopP: %f\n", options.TopP)
	output += fmt.Sprintf("PresencePenalty: %f\n", options.PresencePenalty)
	output += fmt.Sprintf("FrequencyPenalty: %f\n", options.FrequencyPenalty)

	channel <- output
	close(channel)
	return nil
}

func (c *Client) Send(messages []*common.Message, options *common.ChatOptions) (string, error) {
	fmt.Println("Dry run: Would send the following request:")

	for _, msg := range messages {
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
	fmt.Printf("Model: %s\n", options.Model)
	fmt.Printf("Temperature: %f\n", options.Temperature)
	fmt.Printf("TopP: %f\n", options.TopP)
	fmt.Printf("PresencePenalty: %f\n", options.PresencePenalty)
	fmt.Printf("FrequencyPenalty: %f\n", options.FrequencyPenalty)

	return "", nil
}

func (c *Client) Setup() error {
	return nil
}

func (c *Client) SetupFillEnvFileContent(buffer *bytes.Buffer) {
	// No environment variables needed for dry run
}
