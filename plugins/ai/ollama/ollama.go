package ollama

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	ollamaapi "github.com/ollama/ollama/api"
	"github.com/samber/lo"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins"
)

const defaultBaseUrl = "http://localhost:11434"

func NewClient() (ret *Client) {
	vendorName := "Ollama"
	ret = &Client{}

	ret.PluginBase = &plugins.PluginBase{
		Name:            vendorName,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(vendorName),
		ConfigureCustom: ret.configure,
	}

	ret.ApiUrl = ret.AddSetupQuestionCustom("API URL", true,
		"Enter your Ollama URL (as a reminder, it is usually http://localhost:11434')")
	ret.ApiUrl.Value = defaultBaseUrl
	ret.ApiKey = ret.PluginBase.AddSetupQuestion("API key", false)
	ret.ApiKey.Value = ""

	return
}

type Client struct {
	*plugins.PluginBase
	ApiUrl *plugins.SetupQuestion
	ApiKey *plugins.SetupQuestion
	apiUrl *url.URL
	client *ollamaapi.Client
}

type transport_sec struct {
	underlyingTransport http.RoundTripper
	ApiKey              *plugins.SetupQuestion
}

func (t *transport_sec) RoundTrip(req *http.Request) (*http.Response, error) {
	if t.ApiKey.Value != "" {
		req.Header.Add("Authorization", "Bearer "+t.ApiKey.Value)
	}
	return t.underlyingTransport.RoundTrip(req)
}

func (o *Client) configure() (err error) {
	if o.apiUrl, err = url.Parse(o.ApiUrl.Value); err != nil {
		fmt.Printf("cannot parse URL: %s: %v\n", o.ApiUrl.Value, err)
		return
	}

	o.client = ollamaapi.NewClient(o.apiUrl, &http.Client{Timeout: 1200000 * time.Millisecond, Transport: &transport_sec{underlyingTransport: http.DefaultTransport, ApiKey: o.ApiKey}})
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

func (o *Client) SendStream(msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions, channel chan string) (err error) {
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

func (o *Client) Send(ctx context.Context, msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (ret string, err error) {
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

func (o *Client) createChatRequest(msgs []*goopenai.ChatCompletionMessage, opts *common.ChatOptions) (ret ollamaapi.ChatRequest) {
	messages := lo.Map(msgs, func(message *goopenai.ChatCompletionMessage, _ int) (ret ollamaapi.Message) {
		return ollamaapi.Message{Role: message.Role, Content: message.Content}
	})

	options := map[string]interface{}{
		"temperature":       opts.Temperature,
		"presence_penalty":  opts.PresencePenalty,
		"frequency_penalty": opts.FrequencyPenalty,
		"top_p":             opts.TopP,
	}

	if opts.ModelContextLength != 0 {
		options["num_ctx"] = opts.ModelContextLength
	}

	ret = ollamaapi.ChatRequest{
		Model:    opts.Model,
		Messages: messages,
		Options:  options,
	}
	return
}

func (o *Client) NeedsRawMode(modelName string) bool {
	ollamaPrefixes := []string{
		"llama3",
		"llama2",
	}
	for _, prefix := range ollamaPrefixes {
		if strings.HasPrefix(modelName, prefix) {
			return true
		}
	}
	return false
}
