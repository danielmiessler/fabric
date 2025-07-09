package ai

import (
	"context"

	"github.com/danielmiessler/fabric/internal/chat"
	"github.com/danielmiessler/fabric/internal/plugins"

	"github.com/danielmiessler/fabric/internal/common"
)

type Vendor interface {
	plugins.Plugin
	ListModels() ([]string, error)
	SendStream([]*chat.ChatCompletionMessage, *common.ChatOptions, chan string) error
	Send(context.Context, []*chat.ChatCompletionMessage, *common.ChatOptions) (string, error)
	NeedsRawMode(modelName string) bool
}
