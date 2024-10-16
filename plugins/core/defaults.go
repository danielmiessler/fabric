package core

import (
	"fmt"
	"strconv"

	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/pkg/errors"
)

func NeeDefaults() (ret *Defaults) {
	vendorName := "Default"
	ret = &Defaults{
		PluginBase: &plugins.PluginBase{
			Name:             vendorName,
			SetupDescription: "Configure the default AI Vendor and Model",
			EnvNamePrefix:    plugins.BuildEnvVariablePrefix(vendorName),
		},
	}

	ret.Vendor = ret.AddSetting("Vendor", true)
	ret.Model = ret.AddSetupQuestionCustom("Model", true,
		"Enter the index the name of your default model")

	return
}

type Defaults struct {
	*plugins.PluginBase

	Vendor *plugins.Setting
	Model  *plugins.SetupQuestion
}

func (o *Defaults) Setup(vendorsModels *ai.VendorsModels) (err error) {
	vendorsModels.Print()

	if err = o.Ask(o.Name); err != nil {
		return
	}

	index, parseErr := strconv.Atoi(o.Model.Value)
	if parseErr == nil {
		o.Vendor.Value, o.Model.Value = vendorsModels.GetGroupAndItemByItemNumber(index)
	} else {
		o.Vendor.Value = vendorsModels.FindGroupsByItemFirst(o.Model.Value)
	}

	//verify
	vendorNames := vendorsModels.FindGroupsByItem(o.Model.Value)
	if len(vendorNames) == 0 {
		err = errors.Errorf("You need to chose an available default model.")
		return
	}

	fmt.Println()
	o.Vendor.Print()
	o.Model.Print()

	return
}
