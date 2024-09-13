package siliconcloud

import (
	"github.com/danielmiessler/fabric/vendors/openai"
)

func NewClient() (ret *Client) {
	ret = &Client{}
	ret.Client = openai.NewClientCompatible("SiliconCloud", "https://api.siliconflow.cn/v1", nil)
	return
}

type Client struct {
	*openai.Client
}
