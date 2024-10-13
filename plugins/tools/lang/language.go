package lang

import (
	"github.com/danielmiessler/fabric/plugins"
	"golang.org/x/text/language"
)

func NewLanguage() (ret *Language) {

	label := "Language"
	ret = &Language{}

	ret.Plugin = &plugins.Plugin{
		Label:           label,
		EnvNamePrefix:   plugins.BuildEnvVariablePrefix(label),
		ConfigureCustom: ret.configure,
	}

	ret.DefaultLanguage = ret.Plugin.AddSetupQuestionCustom("Output", false,
		"Enter your default want output lang (for example: zh_CN)")

	return
}

type Language struct {
	*plugins.Plugin
	DefaultLanguage *plugins.SetupQuestion
}

func (o *Language) configure() error {
	if o.DefaultLanguage.Value != "" {
		langTag, err := language.Parse(o.DefaultLanguage.Value)
		if err == nil {
			o.DefaultLanguage.Value = langTag.String()
		} else {
			o.DefaultLanguage.Value = ""
		}
	}

	return nil
}
