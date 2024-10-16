package ai

import (
	"context"
	"github.com/danielmiessler/fabric/plugins"

	"github.com/danielmiessler/fabric/common"
)

type Vendor interface {
	plugins.Plugin
	ListModels() ([]string, error)
	SendStream([]*common.Message, *common.ChatOptions, chan string) error
	Send(context.Context, []*common.Message, *common.ChatOptions) (string, error)
}
