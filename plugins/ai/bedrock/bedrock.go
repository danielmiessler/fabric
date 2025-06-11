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
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"

	goopenai "github.com/sashabaranov/go-openai"
)

// BedrockClient is a plugin to add support for Amazon Bedrock
type BedrockClient struct {
	*plugins.PluginBase
	client *bedrockruntime.Client
}

// NewClient returns a new Bedrock plugin client
func NewClient() (ret *BedrockClient) {
	vendorName := "Bedrock"

	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx)
	cfg.APIOptions = append(cfg.APIOptions, middleware.AddUserAgentKeyValue("aiosc", "fabric"))

	if err != nil {
		fmt.Printf("Unable to load AWS Config: %s\n", err)
	}

	client := bedrockruntime.NewFromConfig(cfg)

	ret = &BedrockClient{
		PluginBase: &plugins.PluginBase{
			Name:          vendorName,
			EnvNamePrefix: plugins.BuildEnvVariablePrefix(vendorName),
		},
		client: client,
	}

	return
}

// ListModels lists the models available for use with the Bedrock plugin
func (c *BedrockClient) ListModels() ([]string, error) {
	return MODELS, nil
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

	response, err := c.client.ConverseStream(context.TODO(), &converseInput)
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
	response, err := c.client.Converse(ctx, &converseInput)
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

var MODELS = []string{
	"amazon.nova-micro-v1:0",
	"amazon.nova-lite-v1:0",
	"amazon.nova-pro-v1:0",
	"amazon.nova-premier-v1:0",

	"amazon.titan-tg1-large",
	"amazon.titan-text-premier-v1:0",

	"amazon.titan-text-lite-v1",
	"amazon.titan-text-express-v1",

	"ai21.jamba-instruct-v1:0",
	"ai21.jamba-1-5-large-v1:0",
	"ai21.jamba-1-5-mini-v1:0",

	"anthropic.claude-instant-v1",
	"anthropic.claude-v2",
	"anthropic.claude-v2:1",
	"anthropic.claude-3-haiku-20240307-v1:0",
	"anthropic.claude-3-sonnet-20240229-v1:0",
	"anthropic.claude-3-opus-20240229-v1:0",
	"anthropic.claude-3-5-haiku-20241022-v1:0",
	"anthropic.claude-3-5-sonnet-20240620-v1:0",
	"anthropic.claude-3-5-sonnet-20241022-v2:0",
	"anthropic.claude-3-7-sonnet-20250219-v1:0",
	"anthropic.claude-sonnet-4-20250514-v1:0",
	"anthropic.claude-opus-4-20250514-v1:0",

	"meta.llama3-8b-instruct-v1:0",
	"meta.llama3-70b-instruct-v1:0",
	"meta.llama3-1-8b-instruct-v1:0",
	"meta.llama3-1-70b-instruct-v1:0",
	"meta.llama3-2-11b-instruct-v1:0",
	"meta.llama3-2-90b-instruct-v1:0",
	"meta.llama3-2-1b-instruct-v1:0",
	"meta.llama3-2-3b-instruct-v1:0",
	"meta.llama3-3-70b-instruct-v1:0",
	"meta.llama4-scout-17b-instruct-v1:0",
	"meta.llama4-maverick-17b-instruct-v1:0",

	"mistral.mistral-7b-instruct-v0:2",
	"mistral.mixtral-8x7b-instruct-v0:1",
	"mistral.mistral-small-2402-v1:0",
	"mistral.mistral-large-2402-v1:0",
	"mistral.pixtral-large-2502-v1:0",

	// Cross Region Inferences Profiles
	// https://docs.aws.amazon.com/bedrock/latest/userguide/inference-profiles-support.html#inference-profiles-support-system
	"us.amazon.nova-lite-v1:0",
	"us.amazon.nova-lite-v1:0",
	"us.amazon.nova-micro-v1:0",
	"us.amazon.nova-micro-v1:0",
	"us.amazon.nova-premier-v1:0",
	"us.amazon.nova-premier-v1:0",
	"us.amazon.nova-pro-v1:0",
	"us.amazon.nova-pro-v1:0",
	"us.anthropic.claude-3-5-haiku-20241022-v1:0",
	"us.anthropic.claude-3-5-haiku-20241022-v1:0",
	"us.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"us.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"us.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"us.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"us.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"us.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"us.anthropic.claude-3-haiku-20240307-v1:0",
	"us.anthropic.claude-3-haiku-20240307-v1:0",
	"us.anthropic.claude-3-opus-20240229-v1:0",
	"us.anthropic.claude-3-opus-20240229-v1:0",
	"us.anthropic.claude-3-sonnet-20240229-v1:0",
	"us.anthropic.claude-3-sonnet-20240229-v1:0",
	"us.anthropic.claude-opus-4-20250514-v1:0",
	"us.anthropic.claude-opus-4-20250514-v1:0",
	"us.anthropic.claude-sonnet-4-20250514-v1:0",
	"us.anthropic.claude-sonnet-4-20250514-v1:0",
	"us.deepseek.r1-v1:0",
	"us.deepseek.r1-v1:0",
	"us.meta.llama3-1-405b-instruct-v1:0",
	"us.meta.llama3-1-405b-instruct-v1:0",
	"us.meta.llama3-1-70b-instruct-v1:0",
	"us.meta.llama3-1-70b-instruct-v1:0",
	"us.meta.llama3-1-8b-instruct-v1:0",
	"us.meta.llama3-1-8b-instruct-v1:0",
	"us.meta.llama3-2-11b-instruct-v1:0",
	"us.meta.llama3-2-11b-instruct-v1:0",
	"us.meta.llama3-2-1b-instruct-v1:0",
	"us.meta.llama3-2-1b-instruct-v1:0",
	"us.meta.llama3-2-3b-instruct-v1:0",
	"us.meta.llama3-2-3b-instruct-v1:0",
	"us.meta.llama3-2-90b-instruct-v1:0",
	"us.meta.llama3-2-90b-instruct-v1:0",
	"us.meta.llama3-3-70b-instruct-v1:0",
	"us.meta.llama3-3-70b-instruct-v1:0",
	"us.meta.llama4-maverick-17b-instruct-v1:0",
	"us.meta.llama4-maverick-17b-instruct-v1:0",
	"us.meta.llama4-scout-17b-instruct-v1:0",
	"us.meta.llama4-scout-17b-instruct-v1:0",
	"us.mistral.pixtral-large-2502-v1:0",
	"us.mistral.pixtral-large-2502-v1:0",
	"us.writer.palmyra-x4-v1:0",
	"us.writer.palmyra-x4-v1:0",
	"us.writer.palmyra-x5-v1:0",
	"us.writer.palmyra-x5-v1:0",
	"us-gov.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"us-gov.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"us-gov.anthropic.claude-3-haiku-20240307-v1:0",
	"us-gov.anthropic.claude-3-haiku-20240307-v1:0",
	"eu.amazon.nova-lite-v1:0",
	"eu.amazon.nova-lite-v1:0",
	"eu.amazon.nova-micro-v1:0",
	"eu.amazon.nova-micro-v1:0",
	"eu.amazon.nova-pro-v1:0",
	"eu.amazon.nova-pro-v1:0",
	"eu.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"eu.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"eu.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"eu.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"eu.anthropic.claude-3-haiku-20240307-v1:0",
	"eu.anthropic.claude-3-haiku-20240307-v1:0",
	"eu.anthropic.claude-3-sonnet-20240229-v1:0",
	"eu.anthropic.claude-3-sonnet-20240229-v1:0",
	"eu.anthropic.claude-sonnet-4-20250514-v1:0",
	"eu.anthropic.claude-sonnet-4-20250514-v1:0",
	"eu.meta.llama3-2-1b-instruct-v1:0",
	"eu.meta.llama3-2-1b-instruct-v1:0",
	"eu.meta.llama3-2-3b-instruct-v1:0",
	"eu.meta.llama3-2-3b-instruct-v1:0",
	"eu.mistral.pixtral-large-2502-v1:0",
	"eu.mistral.pixtral-large-2502-v1:0",
	"apac.amazon.nova-lite-v1:0",
	"apac.amazon.nova-lite-v1:0",
	"apac.amazon.nova-micro-v1:0",
	"apac.amazon.nova-micro-v1:0",
	"apac.amazon.nova-pro-v1:0",
	"apac.amazon.nova-pro-v1:0",
	"apac.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"apac.anthropic.claude-3-5-sonnet-20240620-v1:0",
	"apac.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"apac.anthropic.claude-3-5-sonnet-20241022-v2:0",
	"apac.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"apac.anthropic.claude-3-7-sonnet-20250219-v1:0",
	"apac.anthropic.claude-3-haiku-20240307-v1:0",
	"apac.anthropic.claude-3-haiku-20240307-v1:0",
	"apac.anthropic.claude-3-sonnet-20240229-v1:0",
	"apac.anthropic.claude-3-sonnet-20240229-v1:0",
	"apac.anthropic.claude-sonnet-4-20250514-v1:0",
	"apac.anthropic.claude-sonnet-4-20250514-v1:0",
}
