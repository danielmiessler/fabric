package ai

import (
	"context"

	"github.com/danielmiessler/fabric/plugins"
	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
)

type Vendor interface {
	plugins.Plugin
	ListModels() ([]string, error)
	SendStream([]*goopenai.ChatCompletionMessage, *common.ChatOptions, chan string) error
	Send(context.Context, []*goopenai.ChatCompletionMessage, *common.ChatOptions) (string, error)
	NeedsRawMode(modelName string) bool
}
