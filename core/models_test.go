package core

import (
	"errors"
	"testing"
)

func TestNewVendorsModels(t *testing.T) {
	vendors := NewVendorsModels()
	if vendors == nil {
		t.Fatalf("NewVendorsModels() returned nil")
	}
	if len(vendors.VendorsModels) != 0 {
		t.Fatalf("NewVendorsModels() returned non-empty VendorsModels map")
	}
}

func TestFindVendorsByModelFirst(t *testing.T) {
	vendors := NewVendorsModels()
	vendors.AddVendorModels("vendor1", []string{"model1", "model2"})
	vendor := vendors.FindVendorsByModelFirst("model1")
	if vendor != "vendor1" {
		t.Fatalf("FindVendorsByModelFirst() = %v, want %v", vendor, "vendor1")
	}
}

func TestFindVendorsByModel(t *testing.T) {
	vendors := NewVendorsModels()
	vendors.AddVendorModels("vendor1", []string{"model1", "model2"})
	foundVendors := vendors.FindVendorsByModel("model1")
	if len(foundVendors) != 1 || foundVendors[0] != "vendor1" {
		t.Fatalf("FindVendorsByModel() = %v, want %v", foundVendors, []string{"vendor1"})
	}
}

func TestAddVendorModels(t *testing.T) {
	vendors := NewVendorsModels()
	vendors.AddVendorModels("vendor1", []string{"model1", "model2"})
	models := vendors.GetVendorModels("vendor1")
	if len(models) != 2 {
		t.Fatalf("AddVendorModels() failed to add models")
	}
}

func TestAddError(t *testing.T) {
	vendors := NewVendorsModels()
	err := errors.New("sample error")
	vendors.AddError(err)
	if len(vendors.Errs) != 1 {
		t.Fatalf("AddError() failed to add error")
	}
}
