package vendors

import (
	"bytes"
	"context"

	"github.com/danielmiessler/fabric/common"
)

type Vendor interface {
	GetName() string
	IsConfigured() bool
	Configure() error
	ListModels() ([]string, error)
	SendStream([]*common.Message, *common.ChatOptions, chan string) error
	Send(context.Context, []*common.Message, *common.ChatOptions) (string, error)
	Setup() error
	SetupFillEnvFileContent(*bytes.Buffer)
}
