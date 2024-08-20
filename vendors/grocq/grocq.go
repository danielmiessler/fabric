package grocq

import (
	"github.com/danielmiessler/fabric/vendors/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Grocq", "https://api.groq.com/openai/v1", nil)
	return
}

type Client struct {
	*openai.Client
}
