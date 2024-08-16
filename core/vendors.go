package core

import (
	"fmt"
	"sync"

	"github.com/danielmiessler/fabric/common"
)

func NewVendors() (ret *VendorsController) {
	ret = &VendorsController{
		All:        map[string]common.Vendor{},
		Configured: map[string]common.Vendor{},
	}
	return
}

type VendorsController struct {
	All        map[string]common.Vendor
	Configured map[string]common.Vendor

	Models *VendorsModels
}

func (o *VendorsController) AddVendors(vendors ...common.Vendor) {
	for _, vendor := range vendors {
		o.All[vendor.GetName()] = vendor
	}
}

func (o *VendorsController) AddVendorConfigured(vendor common.Vendor) {
	o.Configured[vendor.GetName()] = vendor
}

func (o *VendorsController) ResetConfigured() {
	o.Configured = map[string]common.Vendor{}
	o.Models = nil
	return
}

func (o *VendorsController) GetModels() (ret *VendorsModels) {
	if o.Models == nil {
		o.readModels()
	}
	ret = o.Models
	return
}

func (o *VendorsController) HasConfiguredVendors() bool {
	return len(o.Configured) > 0
}

func (o *VendorsController) readModels() {
	o.Models = NewVendorsModels()

	var wg sync.WaitGroup
	var channels []ChannelName

	errorsChan := make(chan error, 3)

	for _, vendor := range o.Configured {
		// For each vendor:
		//  - Create a channel to collect output from the vendor model's list
		//  - Create a goroutine to query the vendor on its model
		cn := ChannelName{channel: make(chan []string, 1), name: vendor.GetName()}
		channels = append(channels, cn)
		o.createGoroutine(&wg, vendor, cn, errorsChan)
	}

	// Let's wait for completion
	wg.Wait() // Wait for all goroutines to finish
	close(errorsChan)

	for err := range errorsChan {
		fmt.Println(err)
		o.Models.AddError(err)
	}

	// And collect output
	for _, cn := range channels {
		models := <-cn.channel
		if models != nil {
			o.Models.AddVendorModels(cn.name, models)
		}
	}
	return
}

func (o *VendorsController) FindByName(name string) (ret common.Vendor) {
	ret = o.Configured[name]
	return
}

// Create a goroutine to list models for the given vendor
func (o *VendorsController) createGoroutine(wg *sync.WaitGroup, vendor common.Vendor, cn ChannelName, errorsChan chan error) {
	wg.Add(1)

	go func() {
		defer wg.Done()
		models, err := vendor.ListModels()
		if err != nil {
			errorsChan <- err
			cn.channel <- nil
		} else {
			cn.channel <- models
		}
	}()
}
