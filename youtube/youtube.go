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

	ret.ApiKey = ret.AddSetting("ApiKey", false)

	return
}

type YouTube struct {
	*common.Configurable
	ApiKey *common.Setting
}
