package cli

import (
	"os"
	"strconv"

	"github.com/danielmiessler/fabric/internal/core"
	"github.com/danielmiessler/fabric/internal/plugins/ai"
	"github.com/danielmiessler/fabric/internal/plugins/db/fsdb"
)

// handleListingCommands handles listing-related commands
func handleListingCommands(currentFlags *Flags, fabricDb *fsdb.Db, registry *core.PluginRegistry) (err error) {
	if currentFlags.LatestPatterns != "0" {
		var parsedToInt int
		if parsedToInt, err = strconv.Atoi(currentFlags.LatestPatterns); err != nil {
			return
		}

		if err = fabricDb.Patterns.PrintLatestPatterns(parsedToInt); err != nil {
			return
		}
		return
	}

	if currentFlags.ListPatterns {
		err = fabricDb.Patterns.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllModels {
		var models *ai.VendorsModels
		if models, err = registry.VendorManager.GetModels(); err != nil {
			return
		}
		models.Print(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllContexts {
		err = fabricDb.Contexts.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllSessions {
		err = fabricDb.Sessions.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListStrategies {
		err = registry.Strategies.ListStrategies(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListVendors {
		err = registry.ListVendors(os.Stdout)
		return
	}

	return nil
}
