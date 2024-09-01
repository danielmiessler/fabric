package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/danielmiessler/fabric/common"
	"github.com/samber/lo"

	ollamaapi "github.com/ollama/ollama/api"
)

func NewClient() (ret *Client) {
	vendorName := "Ollama"
	ret = &Client{}

	ret.Configurable = &common.Configurable{
		Label:           vendorName,
		EnvNamePrefix:   common.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.ApiUrl = ret.Configurable.AddSetupQuestionCustom("API URL", true,
		"Enter your Ollama URL (as a reminder, it is usually http://localhost:11434)")

	return
}

type Client struct {
	*common.Configurable
	ApiUrl *common.SetupQuestion

	apiUrl *url.URL
	client *ollamaapi.Client
}

func (o *Client) configure() (err error) {
	if o.apiUrl, err = url.Parse(o.ApiUrl.Value); err != nil {
		fmt.Printf("cannot parse URL: %s: %v\n", o.ApiUrl.Value, err)
		return
	}

	o.client = ollamaapi.NewClient(o.apiUrl, &http.Client{Timeout: 1200000 * time.Millisecond})
	return
}

func (o *Client) ListModels() (ret []string, err error) {
	ctx := context.Background()

	var listResp *ollamaapi.ListResponse
	if listResp, err = o.client.List(ctx); err != nil {
		return
	}

	for _, mod := range listResp.Models {
		ret = append(ret, mod.Model)
	}
	return
}

func (o *Client) SendStream(msgs []*common.Message, opts *common.ChatOptions, channel chan string) (err error) {
	req := o.createChatRequest(msgs, opts)

	respFunc := func(resp ollamaapi.ChatResponse) (streamErr error) {
		channel <- resp.Message.Content
		return
	}

	ctx := context.Background()

	if err = o.client.Chat(ctx, &req, respFunc); err != nil {
		return
	}

	close(channel)
	return
}

func (o *Client) Send(ctx context.Context, msgs []*common.Message, opts *common.ChatOptions) (ret string, err error) {
	bf := false

	req := o.createChatRequest(msgs, opts)
	req.Stream = &bf

	respFunc := func(resp ollamaapi.ChatResponse) (streamErr error) {
		ret = resp.Message.Content
		return
	}

	if err = o.client.Chat(ctx, &req, respFunc); err != nil {
		fmt.Printf("FRED --> %s\n", err)
	}
	return
}

func (o *Client) createChatRequest(msgs []*common.Message, opts *common.ChatOptions) (ret ollamaapi.ChatRequest) {
	messages := lo.Map(msgs, func(message *common.Message, _ int) (ret ollamaapi.Message) {
		return ollamaapi.Message{Role: message.Role, Content: message.Content}
	})

	options := map[string]interface{}{
		"temperature":       opts.Temperature,
		"presence_penalty":  opts.PresencePenalty,
		"frequency_penalty": opts.FrequencyPenalty,
		"top_p":             opts.TopP,
	}

	ret = ollamaapi.ChatRequest{
		Model:    opts.Model,
		Messages: messages,
		Options:  options,
	}
	return
}
