package core

import (
	"bytes"
	"context"
	"testing"

	"github.com/danielmiessler/fabric/common"
)

func TestNewVendorsManager(t *testing.T) {
	vendorsManager := NewVendorsManager()
	if vendorsManager == nil {
		t.Fatalf("NewVendorsManager() returned nil")
	}
}

func TestAddVendors(t *testing.T) {
	vendorsManager := NewVendorsManager()
	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)

	if _, exists := vendorsManager.Vendors[mockVendor.GetName()]; !exists {
		t.Fatalf("AddVendors() did not add vendor")
	}
}

func TestGetModels(t *testing.T) {
	vendorsManager := NewVendorsManager()
	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)

	models := vendorsManager.GetModels()
	if models == nil {
		t.Fatalf("GetModels() returned nil")
	}
}

func TestHasVendors(t *testing.T) {
	vendorsManager := NewVendorsManager()
	if vendorsManager.HasVendors() {
		t.Fatalf("HasVendors() should return false for an empty manager")
	}

	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)
	if !vendorsManager.HasVendors() {
		t.Fatalf("HasVendors() should return true after adding a vendor")
	}
}

func TestFindByName(t *testing.T) {
	vendorsManager := NewVendorsManager()
	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)

	foundVendor := vendorsManager.FindByName("testVendor")
	if foundVendor == nil {
		t.Fatalf("FindByName() did not find added vendor")
	}
}

func TestReadModels(t *testing.T) {
	vendorsManager := NewVendorsManager()
	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)

	vendorsManager.readModels()
	if vendorsManager.Models == nil || len(vendorsManager.Models.Vendors) == 0 {
		t.Fatalf("readModels() did not read models correctly")
	}
}

func TestSetup(t *testing.T) {
	vendorsManager := NewVendorsManager()
	mockVendor := &MockVendor{name: "testVendor"}
	vendorsManager.AddVendors(mockVendor)

	vendors, err := vendorsManager.Setup()
	if err != nil {
		t.Fatalf("Setup() error = %v", err)
	}
	if len(vendors) == 0 {
		t.Fatalf("Setup() did not setup any vendors")
	}
}

// MockVendor is a mock implementation of the Vendor interface for testing purposes.
type MockVendor struct {
	*common.Settings
	name string
}

func (o *MockVendor) SendStream(messages []*common.Message, options *common.ChatOptions, strings chan string) error {
	// TODO implement me
	panic("implement me")
}

func (o *MockVendor) Send(ctx context.Context, messages []*common.Message, options *common.ChatOptions) (string, error) {
	// TODO implement me
	panic("implement me")
}

func (o *MockVendor) SetupFillEnvFileContent(buffer *bytes.Buffer) {
	// TODO implement me
	panic("implement me")
}

func (o *MockVendor) IsConfigured() bool {
	return false
}

func (o *MockVendor) GetSettings() *common.Settings {
	return o.Settings
}

func (o *MockVendor) GetName() string {
	return o.name
}

func (o *MockVendor) Configure() error {
	return nil
}

func (o *MockVendor) Setup() error {
	return nil
}

func (o *MockVendor) ListModels() ([]string, error) {
	return []string{"model1", "model2"}, nil
}
