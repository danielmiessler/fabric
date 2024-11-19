package fsdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type PatternsEntity struct {
	*StorageEntity
	SystemPatternFile      string
	UniquePatternsFilePath string
}

// Pattern represents a single pattern with its metadata
type Pattern struct {
	Name        string
	Description string
	Pattern     string
}

// main entry point for getting patterns from any source
func (o *PatternsEntity) GetApplyVariables(source string, variables map[string]string) (*Pattern, error) {
    var pattern *Pattern
    var err error

    // Determine if this is a file path
    isFilePath := strings.HasPrefix(source, "\\") ||
                  strings.HasPrefix(source, "/") ||
                  strings.HasPrefix(source, "~") ||
                  strings.HasPrefix(source, ".")
    
    if isFilePath {
        pattern, err = o.getFromFile(source)
    } else {
        pattern, err = o.getFromDB(source)
    }

    if err != nil {
        return nil, err
    }

    return o.applyVariables(pattern, variables), nil
}

// handles all variable substitution
func (o *PatternsEntity) applyVariables(pattern *Pattern, variables map[string]string) *Pattern {
	if variables != nil && len(variables) > 0 {
			for variableName, value := range variables {
					pattern.Pattern = strings.ReplaceAll(pattern.Pattern, variableName, value)
			}
	}
	return pattern
}

// retrieves a pattern from the database by name
func (o *PatternsEntity) getFromDB(name string) (ret *Pattern, err error) {
	patternPath := filepath.Join(o.Dir, name, o.SystemPatternFile)

	var pattern []byte
	if pattern, err = os.ReadFile(patternPath); err != nil {
		return
	}

	patternStr := string(pattern)
	ret = &Pattern{
		Name:    name,
		Pattern: patternStr,
	}
	return
}

func (o *PatternsEntity) PrintLatestPatterns(latestNumber int) (err error) {
	var contents []byte
	if contents, err = os.ReadFile(o.UniquePatternsFilePath); err != nil {
		err = fmt.Errorf("could not read unique patterns file. Pleas run --updatepatterns (%s)", err)
		return
	}
	uniquePatterns := strings.Split(string(contents), "\n")
	if latestNumber > len(uniquePatterns) {
		latestNumber = len(uniquePatterns)
	}

	for i := len(uniquePatterns) - 1; i > len(uniquePatterns)-latestNumber-1; i-- {
		fmt.Println(uniquePatterns[i])
	}
	return
}

// reads a pattern from a file path and returns it
func (o *PatternsEntity) getFromFile(pathStr string) (*Pattern, error) {
	// Handle home directory expansion
	if strings.HasPrefix(pathStr, "~/") {
			homedir, err := os.UserHomeDir()
			if err != nil {
					return nil, fmt.Errorf("could not get home directory: %v", err)
			}
			pathStr = filepath.Join(homedir, pathStr[2:])
	}

	content, err := os.ReadFile(pathStr)
	if err != nil {
			return nil, fmt.Errorf("could not read pattern file %s: %v", pathStr, err)
	}

	return &Pattern{
			Name:    pathStr,
			Pattern: string(content),
	}, nil
}

// Get required for Storage interface
func (o *PatternsEntity) Get(name string) (*Pattern, error) {
	// Use GetPattern with no variables
	return o.GetApplyVariables(name, nil)
}