package grocq

import (
	"github.com/danielmiessler/fabric/vendors/openai"
	goopenai "github.com/sashabaranov/go-openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Grocq", ret.configure)
	return
}

type Client struct {
	*openai.Client
}

func (oi *Client) configure() (err error) {
	config := goopenai.DefaultConfig(oi.ApiKey.Value)
	config.BaseURL = "https://api.groq.com/openai/v1"
	oi.ApiClient = goopenai.NewClientWithConfig(config)
	return
}
