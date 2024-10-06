package cli

import (
	"github.com/danielmiessler/fabric/core"
	"os"
	"testing"

	"github.com/danielmiessler/fabric/db"
	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	t.Skip("Skipping test for now, collision with flag -t")
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	os.Args = []string{os.Args[0]}
	err := Cli("test")
	assert.Error(t, err)
	assert.Equal(t, core.NoSessionPatternUserMessages, err.Error())
}

func TestSetup(t *testing.T) {
	mockDB := db.NewDb(os.TempDir())

	fabric, err := Setup(mockDB, false)
	assert.Error(t, err)
	assert.Nil(t, fabric)
}
