package litellm

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("LiteLLM", "http://localhost:4000", nil)
	return
}

type Client struct {
	*openai.Client
}
