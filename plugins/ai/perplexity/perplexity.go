package perplexity

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Perplexity", "https://api.perplexity.ai", nil)
	
	return
}

type Client struct {
	*openai.Client
}

// The endpoint needed to list model from perplexity doesn't exist, so we 
// will return a list of models that were available at the time of writing.
// TODO use models config like it is done in Azure Vendor
func (o *Client) ListModels() (ret []string, err error) {
	ret = []string{"llama-3.1-sonar-small-128k-online", "llama-3.1-sonar-large-128k-online", "llama-3.1-sonar-huge-128k-online", "llama-3.1-sonar-small-128k-chat", "llama-3.1-sonar-large-128k-chat", "llama-3.1-8b-instruct", "llama-3.1-70b-instruct"}
	return
}

