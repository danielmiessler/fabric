package cli

import (
	"os"
	"testing"

	"github.com/danielmiessler/fabric/db"
	"github.com/stretchr/testify/assert"
)

func TestCli(t *testing.T) {
	message, err := Cli()
	assert.NoError(t, err)
	assert.Empty(t, message)
}

func TestSetup(t *testing.T) {
	mockDB := db.NewDb(os.TempDir())

	fabric, err := Setup(mockDB, false)
	assert.Error(t, err)
	assert.Nil(t, fabric)
}
