package azure

import (
	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	goopenai "github.com/sashabaranov/go-openai"
	"strings"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Azure", "", ret.configure)
	ret.ApiDeployments = ret.AddSetupQuestionCustom("deployments", true,
		"Enter your Azure deployments (comma separated)")

	return
}

type Client struct {
	*openai.Client
	ApiDeployments *plugins.SetupQuestion

	apiDeployments []string
}

func (oi *Client) configure() (err error) {
	oi.apiDeployments = strings.Split(oi.ApiDeployments.Value, ",")
	oi.ApiClient = goopenai.NewClientWithConfig(goopenai.DefaultAzureConfig(oi.ApiKey.Value, oi.ApiBaseURL.Value))
	return
}

func (oi *Client) ListModels() (ret []string, err error) {
	ret = oi.apiDeployments
	return
}
