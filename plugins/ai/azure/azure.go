package azure

import (
	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"
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
	// Azure Open Client has models and deployments. We need to get the deployments to use them for chat
	// There is no easy way to get the deployments from the API, so we need to ask the user to provide them
	ret = oi.apiDeployments
	return
}
