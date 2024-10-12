package gemini

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/danielmiessler/fabric/common"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

const modelsNamePrefix = "models/"

func NewClient() (ret *Client) {
	vendorName := "Gemini"
	ret = &Client{}

	ret.Configurable = &common.Configurable{
		Label:         vendorName,
		EnvNamePrefix: common.BuildEnvVariablePrefix(vendorName),
	}

	ret.ApiKey = ret.Configurable.AddSetupQuestion("API key", true)

	return
}

type Client struct {
	*common.Configurable
	ApiKey *common.SetupQuestion
}

func (o *Client) ListModels() (ret []string, err error) {
	ctx := context.Background()
	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(o.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	iter := client.ListModels(ctx)
	for {
		var resp *genai.ModelInfo
		if resp, err = iter.Next(); err != nil {
			if errors.Is(err, iterator.Done) {
				err = nil
			}
			break
		}

		name := o.buildModelNameSimple(resp.Name)
		ret = append(ret, name)
	}
	return
}

func (o *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (ret string, err error) {
	systemInstruction, messages := toMessages(msgs)

	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(o.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	model := client.GenerativeModel(o.buildModelNameFull(opts.Model))
	model.SetTemperature(float32(opts.Temperature))
	model.SetTopP(float32(opts.TopP))
	model.SystemInstruction = systemInstruction

	var response *genai.GenerateContentResponse
	if response, err = model.GenerateContent(ctx, messages...); err != nil {
		return
	}

	ret = o.extractText(response)
	return
}

func (o *Client) buildModelNameSimple(fullModelName string) string {
	return strings.TrimPrefix(fullModelName, modelsNamePrefix)
}

func (o *Client) buildModelNameFull(modelName string) string {
	return fmt.Sprintf("%v%v", modelsNamePrefix, modelName)
}

func (o *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) (err error) {
	ctx := context.Background()
	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(o.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	systemInstruction, messages := toMessages(msgs)

	model := client.GenerativeModel(o.buildModelNameFull(opts.Model))
	model.SetTemperature(float32(opts.Temperature))
	model.SetTopP(float32(opts.TopP))
	model.SystemInstruction = systemInstruction

	iter := model.GenerateContentStream(ctx, messages...)
	for {
		if resp, iterErr := iter.Next(); iterErr == nil {
			for _, candidate := range resp.Candidates {
				if candidate.Content != nil {
					for _, part := range candidate.Content.Parts {
						if text, ok := part.(genai.Text); ok {
							channel <- string(text)
						}
					}
				}
			}
		} else {
			if !errors.Is(iterErr, iterator.Done) {
				channel <- fmt.Sprintf("%v\n", iterErr)
			}
			close(channel)
			break
		}
	}
	return
}

func (o *Client) extractText(response *genai.GenerateContentResponse) (ret string) {
	for _, candidate := range response.Candidates {
		if candidate.Content == nil {
			break
		}
		for _, part := range candidate.Content.Parts {
			if text, ok := part.(genai.Text); ok {
				ret += string(text)
			}
		}
	}
	return
}

func toMessages(msgs []*common.Message) (systemInstruction *genai.Content, messages []genai.Part) {
	if len(msgs) >= 2 {
		systemInstruction = &genai.Content{
			Parts: []genai.Part{
				genai.Text(msgs[0].Content),
			},
		}
		for _, msg := range msgs[1:] {
			messages = append(messages, genai.Text(msg.Content))
		}
	} else {
		messages = append(messages, genai.Text(msgs[0].Content))
	}
	return
}
