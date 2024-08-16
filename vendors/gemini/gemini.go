package gemini

import (
	"context"
	"errors"

	"github.com/danielmiessler/fabric/common"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

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

	client *genai.Client
}

func (ge *Client) ListModels() (ret []string, err error) {
	ctx := context.Background()
	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(ge.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	iter := client.ListModels(ctx)
	for {
		var resp *genai.ModelInfo
		if resp, err = iter.Next(); err != nil {
			break
		}
		ret = append(ret, resp.Name)
	}
	return
}

func (ge *Client) Send(msgs []*common.Message, opts *common.ChatOptions) (ret string, err error) {
	systemInstruction, userText := toContent(msgs)

	ctx := context.Background()
	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(ge.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	model := ge.client.GenerativeModel(opts.Model)
	model.SetTemperature(float32(opts.Temperature))
	model.SetTopP(float32(opts.TopP))
	model.SystemInstruction = systemInstruction

	var response *genai.GenerateContentResponse
	if response, err = model.GenerateContent(ctx, genai.Text(userText)); err != nil {
		return
	}

	ret = ge.extractText(response)
	return
}

func (ge *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) (err error) {
	ctx := context.Background()
	var client *genai.Client
	if client, err = genai.NewClient(ctx, option.WithAPIKey(ge.ApiKey.Value)); err != nil {
		return
	}
	defer client.Close()

	systemInstruction, userText := toContent(msgs)

	model := client.GenerativeModel(opts.Model)
	model.SetTemperature(float32(opts.Temperature))
	model.SetTopP(float32(opts.TopP))
	model.SystemInstruction = systemInstruction

	iter := model.GenerateContentStream(ctx, genai.Text(userText))
	for {
		var resp *genai.GenerateContentResponse
		if resp, err = iter.Next(); err == nil {
			for _, candidate := range resp.Candidates {
				if candidate.Content != nil {
					for _, part := range candidate.Content.Parts {
						if text, ok := part.(genai.Text); ok {
							channel <- string(text)
						}
					}
				}
			}
		} else if errors.Is(err, iterator.Done) {
			channel <- "\n"
			close(channel)
			err = nil
		}
		return
	}
}

func (ge *Client) extractText(response *genai.GenerateContentResponse) (ret string) {
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

// Current implementation does not support session
// We need to retrieve the System instruction and User instruction
// Considering how we've built msgs, it's the last 2 messages
// FIXME: I know it's not clean, but will make it for now
func toContent(msgs []*common.Message) (ret *genai.Content, userText string) {
	sys := msgs[len(msgs)-2]
	usr := msgs[len(msgs)-1]

	ret = &genai.Content{
		Parts: []genai.Part{
			genai.Part(genai.Text(sys.Content)),
		},
	}
	userText = usr.Content

	return
}
