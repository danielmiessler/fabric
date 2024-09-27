package cli

import (
	"github.com/danielmiessler/fabric/core"
	"os"
	"testing"

	"github.com/danielmiessler/fabric/db"
	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	os.Args = os.Args[:1]
	message, err := Cli("test")
	assert.Error(t, err)
	assert.Equal(t, core.NoSessionPatternUserMessages, err.Error())
	assert.Empty(t, message)
}

func TestSetup(t *testing.T) {
	mockDB := db.NewDb(os.TempDir())

	fabric, err := Setup(mockDB, false)
	assert.Error(t, err)
	assert.Nil(t, fabric)
}
