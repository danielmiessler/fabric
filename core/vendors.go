package core

import (
	"context"
	"fmt"
	"github.com/danielmiessler/fabric/common"
	"sync"
)

func NewVendorsManager() *VendorsManager {
	return &VendorsManager{
		Vendors: map[string]common.Vendor{},
	}
}

type VendorsManager struct {
	Vendors map[string]common.Vendor
	Models  *VendorsModels
}

func (o *VendorsManager) AddVendors(vendors ...common.Vendor) {
	for _, vendor := range vendors {
		o.Vendors[vendor.GetName()] = vendor
	}
}

func (o *VendorsManager) Reset() {
	o.Vendors = map[string]common.Vendor{}
	o.Models = nil
}

func (o *VendorsManager) GetModels() *VendorsModels {
	if o.Models == nil {
		o.readModels()
	}
	return o.Models
}

func (o *VendorsManager) HasVendors() bool {
	return len(o.Vendors) > 0
}

func (o *VendorsManager) FindByName(name string) common.Vendor {
	return o.Vendors[name]
}

func (o *VendorsManager) readModels() {
	o.Models = NewVendorsModels()

	var wg sync.WaitGroup
	resultsChan := make(chan modelResult, len(o.Vendors))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, vendor := range o.Vendors {
		wg.Add(1)
		go o.fetchVendorModels(ctx, &wg, vendor, resultsChan)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Collect results
	for result := range resultsChan {
		if result.err != nil {
			fmt.Println(result.vendorName, result.err)
			o.Models.AddError(result.err)
			cancel() // Cancel remaining goroutines if needed
		} else {
			o.Models.AddVendorModels(result.vendorName, result.models)
		}
	}
}

func (o *VendorsManager) fetchVendorModels(
	ctx context.Context, wg *sync.WaitGroup, vendor common.Vendor, resultsChan chan<- modelResult) {

	defer wg.Done()

	models, err := vendor.ListModels()
	select {
	case <-ctx.Done():
		// Context canceled, don't send the result
		return
	case resultsChan <- modelResult{vendorName: vendor.GetName(), models: models, err: err}:
		// Result sent
	}
}

type modelResult struct {
	vendorName string
	models     []string
	err        error
}
