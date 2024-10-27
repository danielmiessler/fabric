package groq

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Groq", "https://api.groq.com/openai/v1", nil)
	return
}

type Client struct {
	*openai.Client
}
