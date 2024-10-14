package core

import (
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"os"
	"testing"
)

func TestSaveEnvFile(t *testing.T) {
	registry := NewPluginRegistry(fsdb.NewDb(os.TempDir()))

	err := registry.SaveEnvFile()
	if err != nil {
		t.Fatalf("SaveEnvFile() error = %v", err)
	}
}
