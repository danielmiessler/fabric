package cli

import (
	"github.com/danielmiessler/fabric/internal/core"
)

// handleConfigurationCommands handles configuration-related commands
func handleConfigurationCommands(currentFlags *Flags, registry *core.PluginRegistry) (err error) {
	if currentFlags.UpdatePatterns {
		if err = registry.PatternsLoader.PopulateDB(); err != nil {
			return
		}
		// Save configuration in case any paths were migrated during pattern loading
		err = registry.SaveEnvFile()
		return
	}

	if currentFlags.ChangeDefaultModel {
		if err = registry.Defaults.Setup(); err != nil {
			return
		}
		err = registry.SaveEnvFile()
		return
	}

	return nil
}
