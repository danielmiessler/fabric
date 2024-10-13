package core

import (
	"fmt"
	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/pkg/errors"
	"strconv"
)

func NeeDefaults() (ret *Defaults) {
	vendorName := "Default"
	ret = &Defaults{
		Plugin: &plugins.Plugin{
			Label:         vendorName,
			EnvNamePrefix: plugins.BuildEnvVariablePrefix(vendorName),
		},
	}

	ret.Vendor = ret.AddSetting("Vendor", true)
	ret.Model = ret.AddSetupQuestionCustom("Model", true,
		"Enter the index the name of your default model")

	return
}

type Defaults struct {
	*plugins.Plugin

	Vendor *plugins.Setting
	Model  *plugins.SetupQuestion
}

func (o *Defaults) Setup(vendorsModels *ai.VendorsModels) (err error) {
	vendorsModels.Print()

	if err = o.Ask(o.Label); err != nil {
		return
	}

	index, parseErr := strconv.Atoi(o.Model.Value)
	if parseErr == nil {
		o.Vendor.Value, o.Model.Value = vendorsModels.GetVendorAndModelByModelIndex(index)
	} else {
		o.Vendor.Value = vendorsModels.FindVendorsByModelFirst(o.Model.Value)
	}

	//verify
	vendorNames := vendorsModels.FindVendorsByModel(o.Model.Value)
	if len(vendorNames) == 0 {
		err = errors.Errorf("You need to chose an available default model.")
		return
	}

	fmt.Println()
	o.Vendor.Print()
	o.Model.Print()

	return
}
