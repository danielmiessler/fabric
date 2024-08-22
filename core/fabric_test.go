package core

import (
	"os"
	"testing"

	"github.com/danielmiessler/fabric/db"
)

func TestNewFabric(t *testing.T) {
	_, err := NewFabric(db.NewDb(os.TempDir()))
	if err == nil {
		t.Fatal("without setup error expected")
	}
}

func TestSaveEnvFile(t *testing.T) {
	fabric := NewFabricBase(db.NewDb(os.TempDir()))

	err := fabric.SaveEnvFile()
	if err != nil {
		t.Fatalf("SaveEnvFile() error = %v", err)
	}
}

func TestCopyToClipboard(t *testing.T) {
	t.Skip("skipping test, because of docker env. in ci.")
	fabric := NewFabricBase(db.NewDb(os.TempDir()))

	message := "test message"
	err := fabric.CopyToClipboard(message)
	if err != nil {
		t.Fatalf("CopyToClipboard() error = %v", err)
	}
}

func TestCreateOutputFile(t *testing.T) {
	mockDb := &db.Db{}
	fabric := NewFabricBase(mockDb)

	fileName := "test_output.txt"
	message := "test message"
	err := fabric.CreateOutputFile(message, fileName)
	if err != nil {
		t.Fatalf("CreateOutputFile() error = %v", err)
	}

	defer os.Remove(fileName)
}
