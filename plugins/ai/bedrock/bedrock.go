// Package bedrock provides a plugin to use Amazon Bedrock models.
// Supported models are defined in the MODELS variable.
// To add additional models, append them to the MODELS array. Models must support the Converse and ConverseStream operations
// Authentication uses the  AWS credential provider chain, similar.to the AWS CLI and SDKs
// https://docs.aws.amazon.com/sdkref/latest/guide/standardized-credentials.html
package bedrock

import (
	"context"
	"fmt"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	goopenai "github.com/sashabaranov/go-openai"
)

// BedrockClient is a plugin to add support for Amazon Bedrock
type BedrockClient struct {
	*plugins.PluginBase
	runtimeClient      *bedrockruntime.Client
	controlPlaneClient *bedrock.Client

	bedrockRegion *plugins.SetupQuestion
}

// NewClient returns a new Bedrock plugin client
func NewClient() (ret *BedrockClient) {
	vendorName := "Bedrock"
	ret = &BedrockClient{}

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.APIOptions = append(cfg.APIOptions, middleware.AddUserAgentKeyValue("aiosc", "fabric"))

	if err != nil {
		fmt.Printf("Unable to load AWS Config: %s\n", err)
	}

	runtimeClient := bedrockruntime.NewFromConfig(cfg)
	controlPlaneClient := bedrock.NewFromConfig(cfg)

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.runtimeClient = runtimeClient
	ret.controlPlaneClient = controlPlaneClient

	if cfg.Region != "" {
		ret.bedrockRegion.Value = cfg.Region
	}

	ret.bedrockRegion = ret.PluginBase.AddSetupQuestion("AWS Region", true)

	return
}

func (c *BedrockClient) configure() (err error) {

	if c.bedrockRegion.Value != "" {
		// Load the AWS config with the specified region
		ctx := context.TODO()

		cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(c.bedrockRegion.Value))
		if err != nil {
			return fmt.Errorf("unable to load AWS Config: %w", err)
		}

		cfg.APIOptions = append(cfg.APIOptions, middleware.AddUserAgentKeyValue("aiosc", "fabric"))

		c.runtimeClient = bedrockruntime.NewFromConfig(cfg)
		c.controlPlaneClient = bedrock.NewFromConfig(cfg)
	}

	return
}

// ListModels lists the models available for use with the Bedrock plugin
func (c *BedrockClient) ListModels() ([]string, error) {
	models := []string{}
	ctx := context.TODO()

	foundationModels, err := c.controlPlaneClient.ListFoundationModels(ctx, &bedrock.ListFoundationModelsInput{})
	if err != nil {
		return nil, err
	}

	for _, model := range foundationModels.ModelSummaries {
		models = append(models, *model.ModelId)
	}

	inferenceProfilesPaginator := bedrock.NewListInferenceProfilesPaginator(c.controlPlaneClient, &bedrock.ListInferenceProfilesInput{})

	for inferenceProfilesPaginator.HasMorePages() {
		inferenceProfiles, err := inferenceProfilesPaginator.NextPage(context.TODO())
		if err != nil {
			return nil, err
		}

		for _, profile := range inferenceProfiles.InferenceProfileSummaries {
			models = append(models, *profile.InferenceProfileId)
		}
	}

	return models, nil
}

// SendStream sends the messages to the the Bedrock ConverseStream API
func (c *BedrockClient) SendStream(msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string) (err error) {

	messages := c.toMessages(msgs)

	var converseInput = bedrockruntime.ConverseStreamInput{
		ModelId:  aws.String(opts.Model),
		Messages: messages,
		InferenceConfig: &types.InferenceConfiguration{
			Temperature: aws.Float32(float32(opts.Temperature)),
			TopP:        aws.Float32(float32(opts.TopP))},
	}

	response, err := c.runtimeClient.ConverseStream(context.TODO(), &converseInput)
	if err != nil {
		fmt.Printf("Error conversing with Bedrock: %s\n", err)
		return
	}

	for event := range response.GetStream().Events() {
		// Possible ConverseStream event types
		// https://docs.aws.amazon.com/bedrock/latest/userguide/conversation-inference-call.html#conversation-inference-call-response-converse-stream
		switch v := event.(type) {

		case *types.ConverseStreamOutputMemberContentBlockDelta:
			text, ok := v.Value.Delta.(*types.ContentBlockDeltaMemberText)
			if ok {
				channel <- text.Value
			}

		case *types.ConverseStreamOutputMemberMessageStop:
			channel <- "\n"
			close(channel)

		// Unused Events
		case *types.ConverseStreamOutputMemberMessageStart,
			*types.ConverseStreamOutputMemberContentBlockStart,
			*types.ConverseStreamOutputMemberContentBlockStop,
			*types.ConverseStreamOutputMemberMetadata:

		default:
			fmt.Printf("Error: Unknown stream event type: %T\n", v)
		}
	}

	return nil
}

// Send sends the messages the Bedrock Converse API
func (c *BedrockClient) Send(ctx context.Context, msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {

	messages := c.toMessages(msgs)

	var converseInput = bedrockruntime.ConverseInput{
		ModelId:  aws.String(opts.Model),
		Messages: messages,
	}
	response, err := c.runtimeClient.Converse(ctx, &converseInput)
	if err != nil {
		fmt.Printf("Error conversing with Bedrock: %s\n", err)
		return "", err
	}

	responseText, _ := response.Output.(*types.ConverseOutputMemberMessage)
	responseContentBlock := responseText.Value.Content[0]
	text, _ := responseContentBlock.(*types.ContentBlockMemberText)
	return text.Value, nil
}

func (c *BedrockClient) NeedsRawMode(modelName string) bool {
	return false
}

// toMessages converts the array of input messages from the ChatCompletionMessageType to the
// Bedrock Converse Message type
// The system role messages are mapped to the user role as they contain a mix of system messages,
// pattern content and user input.
func (c *BedrockClient) toMessages(inputMessages []*goopenai.ChatCompletionMessage) (messages []types.Message) {
	for _, msg := range inputMessages {
		roles := map[string]types.ConversationRole{
			goopenai.ChatMessageRoleUser:      types.ConversationRoleUser,
			goopenai.ChatMessageRoleAssistant: types.ConversationRoleAssistant,
			goopenai.ChatMessageRoleSystem:    types.ConversationRoleUser,
		}

		role, ok := roles[msg.Role]
		if !ok {
			continue
		}

		message := types.Message{
			Role:    role,
			Content: []types.ContentBlock{&types.ContentBlockMemberText{Value: msg.Content}},
		}
		messages = append(messages, message)

	}

	return
}
