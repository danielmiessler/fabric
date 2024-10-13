package mistral

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Mistral", "https://api.mistral.ai/v1", nil)
	return
}

type Client struct {
	*openai.Client
}
