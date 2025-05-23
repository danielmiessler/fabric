package azure

import (
	"strings"

	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai/openai"
	goopenai "github.com/sashabaranov/go-openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Azure", "", ret.configure)
	ret.ApiDeployments = ret.AddSetupQuestionCustom("deployments", true,
		"Enter your Azure deployments (comma separated)")
	ret.ApiVersion = ret.AddSetupQuestionCustom("API Version", false,
		"Enter the Azure API version (optional)")

	return
}

type Client struct {
	*openai.Client
	ApiDeployments *plugins.SetupQuestion
	ApiVersion     *plugins.SetupQuestion

	apiDeployments []string
}

func (oi *Client) configure() (err error) {
	oi.apiDeployments = strings.Split(oi.ApiDeployments.Value, ",")
	config := goopenai.DefaultAzureConfig(oi.ApiKey.Value, oi.ApiBaseURL.Value)
	if oi.ApiVersion.Value != "" {
		config.APIVersion = oi.ApiVersion.Value
	}
	oi.ApiClient = goopenai.NewClientWithConfig(config)
	return
}

func (oi *Client) ListModels() (ret []string, err error) {
	ret = oi.apiDeployments
	return
}

func (oi *Client) NeedsRawMode(modelName string) bool {
	return false
}
