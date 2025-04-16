package grokai

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("GrokAI", "https://api.x.ai/v1", nil)
	return
}

type Client struct {
	*openai.Client
}
