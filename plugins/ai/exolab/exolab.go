package exolab

import (
	"strings"

	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai/openai"

	goopenai "github.com/sashabaranov/go-openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatibleNoSetupQuestions("Exolab", ret.configure)

	ret.ApiBaseURL = ret.AddSetupQuestion("API Base URL", true)
	ret.ApiBaseURL.Value = "http://localhost:52415"

	ret.ApiModels = ret.AddSetupQuestionCustom("models", true,
		"Enter your deployed Exolab models (comma separated)")

	return
}

type Client struct {
	*openai.Client
	ApiModels *plugins.SetupQuestion

	apiModels []string
}

func (oi *Client) configure() (err error) {
	oi.apiModels = strings.Split(oi.ApiModels.Value, ",")

	config := goopenai.DefaultConfig("")
	config.BaseURL = oi.ApiBaseURL.Value

	oi.ApiClient = goopenai.NewClientWithConfig(config)
	return
}

func (oi *Client) ListModels() (ret []string, err error) {
	ret = oi.apiModels
	return
}

func (oi *Client) NeedsRawMode(modelName string) bool {
	return false
}
