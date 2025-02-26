package groq

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

// NewClient initializes and returns a new Groq Client.
func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Groq", "https://api.groq.com/openai/v1", nil)
	return
}

// Client wraps the openai.Client to provide additional functionality specific to Groq.
type Client struct {
	*openai.Client
}
