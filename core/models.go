package core

import (
	"fmt"
	"sort"
)

func NewVendorsModels() *VendorsModels {
	return &VendorsModels{VendorsModels: make(map[string][]string)}
}

type VendorsModels struct {
	Vendors       []string
	VendorsModels map[string][]string
	Errs          []error
}

func (o *VendorsModels) AddVendorModels(vendor string, models []string) {
	o.Vendors = append(o.Vendors, vendor)
	o.VendorsModels[vendor] = models
}

func (o *VendorsModels) GetVendorAndModelByModelIndex(modelIndex int) (vendor string, model string) {
	vendorModelIndexFrom := 0
	vendorModelIndexTo := 0
	for _, currenVendor := range o.Vendors {
		vendorModelIndexFrom = vendorModelIndexTo + 1
		vendorModelIndexTo = vendorModelIndexFrom + len(o.VendorsModels[currenVendor]) - 1

		if modelIndex >= vendorModelIndexFrom && modelIndex <= vendorModelIndexTo {
			vendor = currenVendor
			model = o.VendorsModels[currenVendor][modelIndex-vendorModelIndexFrom]
			break
		}
	}
	return
}

func (o *VendorsModels) AddError(err error) {
	o.Errs = append(o.Errs, err)
}

func (o *VendorsModels) Print() {
	fmt.Printf("\nAvailable vendor models:\n")

	sort.Strings(o.Vendors)

	var currentModelIndex int
	for _, vendor := range o.Vendors {
		fmt.Println()
		fmt.Printf("%s\n", vendor)
		fmt.Println()
		currentModelIndex = o.PrintVendor(vendor, currentModelIndex)
	}
	return
}

func (o *VendorsModels) PrintVendor(vendor string, modelIndex int) (currentModelIndex int) {
	currentModelIndex = modelIndex
	models := o.VendorsModels[vendor]
	for _, model := range models {
		currentModelIndex++
		fmt.Printf("\t[%d]\t%s\n", currentModelIndex, model)
	}
	fmt.Println()
	return
}

func (o *VendorsModels) GetVendorModels(vendor string) (models []string) {
	models = o.VendorsModels[vendor]
	return
}

func (o *VendorsModels) HasVendor(vendor string) (ret bool) {
	ret = o.VendorsModels[vendor] != nil
	return
}

func (o *VendorsModels) FindVendorsByModelFirst(model string) (ret string) {
	vendors := o.FindVendorsByModel(model)
	if len(vendors) > 0 {
		ret = vendors[0]
	}
	return
}

func (o *VendorsModels) FindVendorsByModel(model string) (vendors []string) {
	for vendor, models := range o.VendorsModels {
		for _, m := range models {
			if m == model {
				vendors = append(vendors, vendor)
				continue
			}
		}
	}
	return
}
