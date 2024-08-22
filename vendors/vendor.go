package vendors

import "github.com/danielmiessler/fabric/common"

type Vendor interface {
	GetName() string
	IsConfigured() bool
	Configure() error
	ListModels() ([]string, error)
	SendStream([]*common.Message, *common.ChatOptions, chan string) error
	Send([]*common.Message, *common.ChatOptions) (string, error)
	GetSettings() common.Settings
	Setup() error
}
