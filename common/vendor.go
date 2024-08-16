package common

type Vendor interface {
	GetName() string
	IsConfigured() bool
	Configure() error
	ListModels() ([]string, error)
	SendStream([]*Message, *ChatOptions, chan string) error
	Send([]*Message, *ChatOptions) (string, error)
	GetSettings() Settings
	Setup() error
}
