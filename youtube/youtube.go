package youtube

import (
	"github.com/danielmiessler/fabric/common"
)

func NewYouTube() (ret *YouTube) {

	label := "YouTube"
	ret = &YouTube{}

	ret.Configurable = &common.Configurable{
		Label:         label,
		EnvNamePrefix: common.BuildEnvVariablePrefix(label),
	}

	ret.ApiKey = ret.AddSetupQuestion("API key", true)

	return
}

type YouTube struct {
	*common.Configurable
	ApiKey *common.SetupQuestion
}
