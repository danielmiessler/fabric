package deepseek

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("DeepSeek", "https://api.deepseek.com", nil)
	return
}

type Client struct {
	*openai.Client
}
