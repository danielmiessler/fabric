// File: plugins/ai/cerebras/cerebras.go
package cerebras

import (
	"github.com/danielmiessler/fabric/plugins/ai/openai"
)

// NewClient initializes and returns a new Cerebras Client.
func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("Cerebras", "https://api.cerebras.ai/v1", nil)
	return
}

// Client wraps the openai.Client to provide additional functionality specific to Cerebras.
type Client struct {
	*openai.Client
}
