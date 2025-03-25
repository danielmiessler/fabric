package common

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	// MaxFileSize is the maximum size of a file that can be created (10MB)
	MaxFileSize = 10 * 1024 * 1024
)

// FileChange represents a single file change operation to be performed
type FileChange struct {
	Operation string `json:"operation"` // "create" or "update"
	Path      string `json:"path"`      // Relative path from project root
	Content   string `json:"content"`   // New file content
}

// ParseFileChanges extracts and parses the FILE_CHANGES section from LLM output
func ParseFileChanges(output string) ([]FileChange, error) {
	// Find the FILE_CHANGES: section marker
	fileChangesStart := strings.Index(output, "FILE_CHANGES:")
	if fileChangesStart == -1 {
		return nil, nil // No file changes section found
	}

	// Extract the JSON part
	jsonStart := fileChangesStart + len("FILE_CHANGES:")
	// Find the first [ after the FILE_CHANGES: marker
	jsonArrayStart := strings.Index(output[jsonStart:], "[")
	if jsonArrayStart == -1 {
		return nil, fmt.Errorf("invalid FILE_CHANGES format: no JSON array found")
	}
	jsonStart += jsonArrayStart

	// Find the matching closing bracket for the array with proper bracket counting
	bracketCount := 0
	jsonEnd := jsonStart
	for i := jsonStart; i < len(output); i++ {
		if output[i] == '[' {
			bracketCount++
		} else if output[i] == ']' {
			bracketCount--
			if bracketCount == 0 {
				jsonEnd = i + 1
				break
			}
		}
	}

	if bracketCount != 0 {
		return nil, fmt.Errorf("invalid FILE_CHANGES format: unbalanced brackets")
	}

	// Parse the JSON
	var fileChanges []FileChange
	err := json.Unmarshal([]byte(output[jsonStart:jsonEnd]), &fileChanges)
	if err != nil {
		return nil, fmt.Errorf("failed to parse FILE_CHANGES JSON: %w", err)
	}

	// Validate file changes
	for i, change := range fileChanges {
		// Validate operation
		if change.Operation != "create" && change.Operation != "update" {
			return nil, fmt.Errorf("invalid operation for file change %d: %s", i, change.Operation)
		}

		// Validate path
		if change.Path == "" {
			return nil, fmt.Errorf("empty path for file change %d", i)
		}

		// Check for suspicious paths (directory traversal)
		if strings.Contains(change.Path, "..") {
			return nil, fmt.Errorf("suspicious path for file change %d: %s", i, change.Path)
		}

		// Check file size
		if len(change.Content) > MaxFileSize {
			return nil, fmt.Errorf("file content too large for file change %d: %d bytes", i, len(change.Content))
		}
	}

	return fileChanges, nil
}

// ApplyFileChanges applies the parsed file changes to the file system
func ApplyFileChanges(projectRoot string, changes []FileChange) error {
	for i, change := range changes {
		// Get the absolute path
		absPath := filepath.Join(projectRoot, change.Path)

		// Create directories if necessary
		dir := filepath.Dir(absPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s for file change %d: %w", dir, i, err)
		}

		// Write the file
		if err := os.WriteFile(absPath, []byte(change.Content), 0644); err != nil {
			return fmt.Errorf("failed to write file %s for file change %d: %w", absPath, i, err)
		}

		fmt.Printf("Applied %s operation to %s\n", change.Operation, change.Path)
	}

	return nil
}
