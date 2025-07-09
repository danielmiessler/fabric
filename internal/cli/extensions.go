package cli

import (
	"github.com/danielmiessler/fabric/internal/core"
)

// handleExtensionCommands handles extension-related commands
func handleExtensionCommands(currentFlags *Flags, registry *core.PluginRegistry) (err error) {
	if currentFlags.ListExtensions {
		err = registry.TemplateExtensions.ListExtensions()
		return
	}

	if currentFlags.AddExtension != "" {
		err = registry.TemplateExtensions.RegisterExtension(currentFlags.AddExtension)
		return
	}

	if currentFlags.RemoveExtension != "" {
		err = registry.TemplateExtensions.RemoveExtension(currentFlags.RemoveExtension)
		return
	}

	return nil
}
