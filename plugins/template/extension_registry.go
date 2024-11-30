package template

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	// Add this import
)

// ExtensionDefinition represents a single extension configuration
type ExtensionDefinition struct {
    // Global properties
    Name        string                     `yaml:"name"`
    Executable  string                     `yaml:"executable"`
    Type        string                     `yaml:"type"`
    Timeout     string                     `yaml:"timeout"`
    Description string                     `yaml:"description"`
    Version     string                     `yaml:"version"`
    Env         []string                   `yaml:"env"`
    
    // Operation-specific commands
    Operations  map[string]OperationConfig `yaml:"operations"`
    
    // Additional config
    Config map[string]interface{} `yaml:"config"`
}

type OperationConfig struct {
    CmdTemplate string `yaml:"cmd_template"`
}

type ExtensionRegistry struct {
    configDir string
    registry  struct {
        Extensions   map[string]*ExtensionDefinition
        ConfigHashes map[string]string
        ExecutableHashes map[string]string
    }
}


// Helper methods for Config access
func (e *ExtensionDefinition) GetOutputMethod() string {
	if output, ok := e.Config["output"].(map[string]interface{}); ok {
			if method, ok := output["method"].(string); ok {
					return method
			}
	}
	return "stdout" // default to stdout if not specified
}

func (e *ExtensionDefinition) GetFileConfig() map[string]interface{} {
	if output, ok := e.Config["output"].(map[string]interface{}); ok {
			if fileConfig, ok := output["file_config"].(map[string]interface{}); ok {
					return fileConfig
			}
	}
	return nil
}

func (e *ExtensionDefinition) IsCleanupEnabled() bool {
	if fc := e.GetFileConfig(); fc != nil {
			if cleanup, ok := fc["cleanup"].(bool); ok {
					return cleanup
			}
	}
	return false // default to no cleanup
}


func NewExtensionRegistry(configDir string) *ExtensionRegistry {
	r := &ExtensionRegistry{
			configDir: configDir,
	}
	r.registry.Extensions = make(map[string]*ExtensionDefinition)
	r.registry.ConfigHashes = make(map[string]string)
	r.registry.ExecutableHashes = make(map[string]string)
	
	// Ensure extensions directory exists
	r.ensureConfigDir()
	
	// Load existing registry if it exists
	if err := r.loadRegistry(); err != nil {
			// Since we're in a constructor, we can't return error
			// Log it if we have logging, but continue with empty registry
			if Debug {
					fmt.Printf("Warning: could not load extension registry: %v\n", err)
			}
	}
	
	return r
}

func (r *ExtensionRegistry) ensureConfigDir() error {
    extDir := filepath.Join(r.configDir, "extensions")
    return os.MkdirAll(extDir, 0755)
}

func (r *ExtensionRegistry) Register(configPath string) error {
    // Read and parse the extension definition
    data, err := os.ReadFile(configPath)
    if err != nil {
        return fmt.Errorf("failed to read config file: %w", err)
    }

    var ext ExtensionDefinition
    if err := yaml.Unmarshal(data, &ext); err != nil {
        return fmt.Errorf("failed to parse config file: %w", err)
    }

    // Verify Executable exists
    if _, err := os.Stat(ext.Executable); err != nil {
        return fmt.Errorf("Executable not found: %w", err)
    }

    // Calculate hashes using template package functions
    configHash := ComputeStringHash(string(data))
    ExecutableHash, err := ComputeHash(ext.Executable)
    if err != nil {
        return fmt.Errorf("failed to hash Executable: %w", err)
    }

    // Store extension and hashes
    r.registry.Extensions[ext.Name] = &ext
    r.registry.ConfigHashes[ext.Name] = configHash
    r.registry.ExecutableHashes[ext.Name] = ExecutableHash

    return r.saveRegistry()
}

func (r *ExtensionRegistry) Remove(name string) error {
    if _, exists := r.registry.Extensions[name]; !exists {
        return fmt.Errorf("extension %s not found", name)
    }

    delete(r.registry.Extensions, name)
    delete(r.registry.ConfigHashes, name)
    delete(r.registry.ExecutableHashes, name)

    return r.saveRegistry()
}

func (r *ExtensionRegistry) Verify(name string) error {
    ext, exists := r.registry.Extensions[name]
    if !exists {
        return fmt.Errorf("extension %s not found", name)
    }

    // Verify Executable hash using template package function
    currentExecutableHash, err := ComputeHash(ext.Executable)
    if err != nil {
        return fmt.Errorf("failed to verify Executable: %w", err)
    }

    if currentExecutableHash != r.registry.ExecutableHashes[name] {
        return fmt.Errorf("Executable hash mismatch for %s", name)
    }

    return nil
}

func (r *ExtensionRegistry) GetExtension(name string) (*ExtensionDefinition, error) {
    ext, exists := r.registry.Extensions[name]
    if !exists {
        return nil, fmt.Errorf("extension %s not found", name)
    }
    
    if err := r.Verify(name); err != nil {
        return nil, err
    }
    
    return ext, nil
}

func (r *ExtensionRegistry) ListExtensions() ([]*ExtensionDefinition, error) {
    exts := make([]*ExtensionDefinition, 0, len(r.registry.Extensions))
    for _, ext := range r.registry.Extensions {
        exts = append(exts, ext)
    }
    return exts, nil
}

func (r *ExtensionRegistry) calculateFileHash(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }

    return hex.EncodeToString(h.Sum(nil)), nil
}

func (r *ExtensionRegistry) saveRegistry() error {
    data, err := yaml.Marshal(r.registry)
    if err != nil {
        return fmt.Errorf("failed to marshal extension registry: %w", err)
    }

    registryPath := filepath.Join(r.configDir, "extensions", "extensions.yaml")
    return os.WriteFile(registryPath, data, 0644)
}

func (r *ExtensionRegistry) loadRegistry() error {
	registryPath := filepath.Join(r.configDir, "extensions", "extensions.yaml")
	data, err := os.ReadFile(registryPath)
	if err != nil {
			if os.IsNotExist(err) {
					return nil // New registry
			}
			return fmt.Errorf("failed to read extension registry: %w", err)
	}

	// Need to unmarshal the data into our registry
	if err := yaml.Unmarshal(data, &r.registry); err != nil {
			return fmt.Errorf("failed to parse extension registry: %w", err)
	}

	return nil
}