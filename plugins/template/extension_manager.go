package template

import (
	"fmt"
	"path/filepath"
)

// ExtensionManager handles the high-level operations of the extension system
type ExtensionManager struct {
    registry *ExtensionRegistry
    executor *ExtensionExecutor
    configDir string
}

// NewExtensionManager creates a new extension manager instance
func NewExtensionManager(configDir string) *ExtensionManager {
    registry := NewExtensionRegistry(configDir)
    return &ExtensionManager{
        registry: registry,
        executor: NewExtensionExecutor(registry),
        configDir: configDir,
    }
}

// ListExtensions handles the listextensions flag action
func (em *ExtensionManager) ListExtensions() error {
	extensions, err := em.registry.ListExtensions()
	if err != nil {
			return fmt.Errorf("failed to list extensions: %w", err)
	}

	for _, ext := range extensions {
			fmt.Printf("Name: %s\n", ext.Name)
			fmt.Printf("  Executable: %s\n", ext.Executable)
			fmt.Printf("  Type: %s\n", ext.Type)
			fmt.Printf("  Timeout: %s\n", ext.Timeout)
			fmt.Printf("  Description: %s\n", ext.Description)
			fmt.Printf("  Version: %s\n", ext.Version)
			
			fmt.Printf("  Operations:\n")
			for opName, opConfig := range ext.Operations {
					fmt.Printf("    %s:\n", opName)
					fmt.Printf("      Command Template: %s\n", opConfig.CmdTemplate)
			}
			
			if fileConfig := ext.GetFileConfig(); fileConfig != nil {
					fmt.Printf("  File Configuration:\n")
					for k, v := range fileConfig {
							fmt.Printf("    %s: %v\n", k, v)
					}
			}
			fmt.Printf("\n")
	}

	return nil
}

// RegisterExtension handles the addextension flag action
func (em *ExtensionManager) RegisterExtension(configPath string) error {
    absPath, err := filepath.Abs(configPath)
    if err != nil {
        return fmt.Errorf("invalid config path: %w", err)
    }

    if err := em.registry.Register(absPath); err != nil {
        return fmt.Errorf("failed to register extension: %w", err)
    }

    return nil
}

// RemoveExtension handles the rmextension flag action
func (em *ExtensionManager) RemoveExtension(name string) error {
    if err := em.registry.Remove(name); err != nil {
        return fmt.Errorf("failed to remove extension: %w", err)
    }

    return nil
}

// ProcessExtension handles template processing for extension directives
func (em *ExtensionManager) ProcessExtension(name, operation, value string) (string, error) {
	return em.executor.Execute(name, operation, value)
}