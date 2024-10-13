package openrouter

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("OpenRouter", "https://openrouter.ai/api/v1", nil)

	return
}

type Client struct {
	*openai.Client
}
