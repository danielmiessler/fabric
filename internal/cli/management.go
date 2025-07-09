package cli

import (
	"github.com/danielmiessler/fabric/internal/plugins/db/fsdb"
)

// handleManagementCommands handles management-related commands (delete, print, etc.)
func handleManagementCommands(currentFlags *Flags, fabricDb *fsdb.Db) (err error) {
	if currentFlags.WipeContext != "" {
		err = fabricDb.Contexts.Delete(currentFlags.WipeContext)
		return
	}

	if currentFlags.WipeSession != "" {
		err = fabricDb.Sessions.Delete(currentFlags.WipeSession)
		return
	}

	if currentFlags.PrintSession != "" {
		err = fabricDb.Sessions.PrintSession(currentFlags.PrintSession)
		return
	}

	if currentFlags.PrintContext != "" {
		err = fabricDb.Contexts.PrintContext(currentFlags.PrintContext)
		return
	}

	return nil
}
