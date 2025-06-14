package bedrock

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"
)

// NewClient creates a new AWS Bedrock client
func NewClient() *Client {
	ret := &Client{}
	ret.PluginBase = &plugins.PluginBase{
		Name:            "Bedrock",
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix("Bedrock"),
		ConfigureCustom: ret.configure,
	}

	// Add setup questions
	ret.AwsAccessKeyId = ret.AddSetupQuestion("AWS Access Key ID", true)
	ret.AwsSecretAccessKey = ret.AddSetupQuestionCustom("AWS Secret Access Key", true, "Enter your AWS Secret Access Key")
	ret.AwsRegion = ret.AddSetupQuestionCustom("AWS Region", true, "Enter your AWS Region (e.g., us-east-1)")
	ret.AwsSessionToken = ret.AddSetupQuestionCustom("AWS Session Token", false, "Enter your AWS Session Token (optional, for temporary credentials)")
	ret.AwsProfile = ret.AddSetupQuestionCustom("AWS Profile", false, "Enter your AWS Profile name (optional, alternative to explicit credentials)")
	ret.PreferredModels = ret.AddSetupQuestionCustom("Preferred Models", false, "Enter preferred model IDs (comma separated, optional)")

	return ret
}

// Client represents the AWS Bedrock client
type Client struct {
	*plugins.PluginBase

	AwsAccessKeyId     *plugins.SetupQuestion
	AwsSecretAccessKey *plugins.SetupQuestion
	AwsRegion          *plugins.SetupQuestion
	AwsSessionToken    *plugins.SetupQuestion
	AwsProfile         *plugins.SetupQuestion
	PreferredModels    *plugins.SetupQuestion

	bedrockClient   *bedrockruntime.Client
	preferredModels []string
}

// configure initializes the AWS Bedrock client
func (c *Client) configure() error {
	var cfg aws.Config
	var err error

	ctx := context.Background()

	// Configure based on available credentials
	if c.AwsProfile.Value != "" {
		// Use AWS profile
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(c.AwsRegion.Value),
			config.WithSharedConfigProfile(c.AwsProfile.Value),
		)
	} else if c.AwsAccessKeyId.Value != "" && c.AwsSecretAccessKey.Value != "" {
		// Use explicit credentials
		creds := credentials.NewStaticCredentialsProvider(
			c.AwsAccessKeyId.Value,
			c.AwsSecretAccessKey.Value,
			c.AwsSessionToken.Value,
		)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(c.AwsRegion.Value),
			config.WithCredentialsProvider(creds),
		)
	} else {
		// Use default credential chain (IAM role, etc.)
		cfg, err = config.LoadDefaultConfig(ctx,
			config.WithRegion(c.AwsRegion.Value),
		)
	}

	if err != nil {
		return fmt.Errorf("failed to load AWS configuration: %w", err)
	}

	// Create Bedrock Runtime client
	c.bedrockClient = bedrockruntime.NewFromConfig(cfg)

	// Parse preferred models
	if c.PreferredModels.Value != "" {
		models := strings.Split(c.PreferredModels.Value, ",")
		for _, model := range models {
			c.preferredModels = append(c.preferredModels, strings.TrimSpace(model))
		}
	}

	return nil
}

// ListModels returns available Bedrock models
func (c *Client) ListModels() ([]string, error) {
	// If user specified preferred models, return those
	if len(c.preferredModels) > 0 {
		return c.preferredModels, nil
	}

	// Return a default list of commonly available Bedrock models
	// Note: Bedrock doesn't have a ListModels API, so we provide a curated list
	return []string{
		"anthropic.claude-3-opus-20240229-v1:0",
		"anthropic.claude-3-sonnet-20240229-v1:0",
		"us.anthropic.claude-3-5-sonnet-20240620-v1:0",
		"us.anthropic.claude-3-5-sonnet-20241022-v2:0",
		"eu.anthropic.claude-3-5-sonnet-20240620-v1:0",
		"anthropic.claude-3-haiku-20240307-v1:0",
		"anthropic.claude-v2:1",
		"anthropic.claude-v2",
		"anthropic.claude-instant-v1",
		"meta.llama3-70b-instruct-v1:0",
		"meta.llama3-8b-instruct-v1:0",
		"meta.llama2-70b-chat-v1",
		"meta.llama2-13b-chat-v1",
		"mistral.mistral-7b-instruct-v0:2",
		"mistral.mixtral-8x7b-instruct-v0:1",
		"cohere.command-text-v14",
		"cohere.command-light-text-v14",
		"ai21.j2-ultra-v1",
		"ai21.j2-mid-v1",
	}, nil
}

// Send sends a non-streaming request to Bedrock
func (c *Client) Send(ctx context.Context, messages []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (string, error) {
	if c.bedrockClient == nil {
		return "", fmt.Errorf("bedrock client not configured")
	}

	// Prepare request parameters
	maxTokens := 4096
	// Note: ChatOptions doesn't have MaxTokens field, using default

	temperature := float32(0.7)
	if opts != nil && opts.Temperature > 0 {
		temperature = float32(opts.Temperature)
	}

	topP := float32(0.9)
	if opts != nil && opts.TopP > 0 {
		topP = float32(opts.TopP)
	}

	// Convert messages to Bedrock format
	body, err := ConvertToBedrockFormat(messages, opts.Model, maxTokens, temperature, topP)
	if err != nil {
		return "", fmt.Errorf("failed to convert messages to Bedrock format: %w", err)
	}

	// Invoke model
	input := &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(opts.Model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	}

	output, err := c.bedrockClient.InvokeModel(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to invoke Bedrock model: %w", err)
	}

	// Parse response based on provider
	provider := GetModelProvider(opts.Model)
	response, err := ParseNonStreamingResponse(output.Body, provider)
	if err != nil {
		return "", fmt.Errorf("failed to parse Bedrock response: %w", err)
	}

	return response, nil
}

// SendStream sends a streaming request to Bedrock
func (c *Client) SendStream(messages []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string) error {
	defer func() {
		// Ensure channel is always closed, even if there's a panic
		if r := recover(); r != nil {
			close(channel)
			panic(r)
		}
	}()

	if c.bedrockClient == nil {
		close(channel)
		return fmt.Errorf("bedrock client not configured")
	}

	ctx := context.Background()

	// Check if model supports streaming
	modelInfo, exists := GetModelInfo(opts.Model)
	if exists && !modelInfo.SupportsStreaming {
		// Fall back to non-streaming for models that don't support it
		response, err := c.Send(ctx, messages, opts)
		if err != nil {
			close(channel)
			return err
		}
		channel <- response
		close(channel)
		return nil
	}

	// Prepare request parameters
	maxTokens := 4096
	// Note: ChatOptions doesn't have MaxTokens field, using default

	temperature := float32(0.7)
	if opts != nil && opts.Temperature > 0 {
		temperature = float32(opts.Temperature)
	}

	topP := float32(0.9)
	if opts != nil && opts.TopP > 0 {
		topP = float32(opts.TopP)
	}

	// Convert messages to Bedrock format
	body, err := ConvertToBedrockFormat(messages, opts.Model, maxTokens, temperature, topP)
	if err != nil {
		close(channel)
		return fmt.Errorf("failed to convert messages to Bedrock format: %w", err)
	}

	// Invoke model with streaming
	input := &bedrockruntime.InvokeModelWithResponseStreamInput{
		ModelId:     aws.String(opts.Model),
		ContentType: aws.String("application/json"),
		Accept:      aws.String("application/json"),
		Body:        body,
	}

	output, err := c.bedrockClient.InvokeModelWithResponseStream(ctx, input)
	if err != nil {
		close(channel)
		return fmt.Errorf("failed to invoke Bedrock model with streaming: %w", err)
	}

	// Handle streaming response - this function will close the channel
	return HandleStreamingResponse(ctx, output, opts.Model, channel)
}

// NeedsRawMode indicates if the model needs raw mode
func (c *Client) NeedsRawMode(modelName string) bool {
	// Bedrock models generally don't need raw mode
	return false
}