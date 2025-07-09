package cli

import (
	"os"
	"path/filepath"

	"github.com/danielmiessler/fabric/internal/core"
	"github.com/danielmiessler/fabric/internal/plugins/db/fsdb"
)

// initializeFabric initializes the fabric database and plugin registry
func initializeFabric() (registry *core.PluginRegistry, err error) {
	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return
	}

	fabricDb := fsdb.NewDb(filepath.Join(homedir, ".config/fabric"))
	if err = fabricDb.Configure(); err != nil {
		return
	}

	if registry, err = core.NewPluginRegistry(fabricDb); err != nil {
		return
	}

	return
}
