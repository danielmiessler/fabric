package ai

import (
	"context"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/plugins"

	"github.com/danielmiessler/fabric/common"
)

type Vendor interface {
	plugins.Plugin
	ListModels() ([]string, error)
	SendStream([]*chat.ChatCompletionMessage, *common.ChatOptions, chan string) error
	Send(context.Context, []*chat.ChatCompletionMessage, *common.ChatOptions) (string, error)
	NeedsRawMode(modelName string) bool
}
