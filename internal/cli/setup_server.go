package cli

import (
	"github.com/danielmiessler/fabric/internal/core"
	restapi "github.com/danielmiessler/fabric/internal/server"
)

// handleSetupAndServerCommands handles setup and server-related commands
func handleSetupAndServerCommands(currentFlags *Flags, registry *core.PluginRegistry, version string) (err error) {
	// if the setup flag is set, run the setup function
	if currentFlags.Setup {
		err = registry.Setup()
		return
	}

	if currentFlags.Serve {
		registry.ConfigureVendors()
		err = restapi.Serve(registry, currentFlags.ServeAddress, currentFlags.ServeAPIKey)
		return
	}

	if currentFlags.ServeOllama {
		registry.ConfigureVendors()
		err = restapi.ServeOllama(registry, currentFlags.ServeAddress, version)
		return
	}

	return nil
}
