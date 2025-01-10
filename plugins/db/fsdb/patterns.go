package fsdb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins/template"
)

const inputSentinel = "__FABRIC_INPUT_SENTINEL_TOKEN__"

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

// GetApplyVariables main entry point for getting patterns from any source
func (o *PatternsEntity) GetApplyVariables(
	source string, variables map[string]string, input string) (pattern *Pattern, err error) {

	// Determine if this is a file path
	isFilePath := strings.HasPrefix(source, "\\") ||
		strings.HasPrefix(source, "/") ||
		strings.HasPrefix(source, "~") ||
		strings.HasPrefix(source, ".")

	if isFilePath {
		// Resolve the file path using GetAbsolutePath
		absPath, err := common.GetAbsolutePath(source)
		if err != nil {
			return nil, fmt.Errorf("could not resolve file path: %v", err)
		}

		// Use the resolved absolute path to get the pattern
		pattern, err = o.getFromFile(absPath)
	} else {
		// Otherwise, get the pattern from the database
		pattern, err = o.getFromDB(source)
	}

	if err != nil {
		return
	}

	// Apply variables to the pattern
	err = o.applyVariables(pattern, variables, input)
	return
}

func (o *PatternsEntity) applyVariables(
	pattern *Pattern, variables map[string]string, input string) (err error) {

	// Ensure pattern has an {{input}} placeholder
	// If not present, append it on a new line
	if !strings.Contains(pattern.Pattern, "{{input}}") {
		if !strings.HasSuffix(pattern.Pattern, "\n") {
			pattern.Pattern += "\n"
		}
		pattern.Pattern += "{{input}}"
	}

	// Temporarily replace {{input}} with a sentinel token to protect it
	// from recursive variable resolution
	withSentinel := strings.ReplaceAll(pattern.Pattern, "{{input}}", inputSentinel)

	// Process all other template variables in the pattern
	// At this point, our sentinel ensures {{input}} won't be affected
	var processed string
	if processed, err = template.ApplyTemplate(withSentinel, variables, ""); err != nil {
		return
	}

	// Finally, replace our sentinel with the actual user input
	// The input has already been processed for variables if InputHasVars was true
	pattern.Pattern = strings.ReplaceAll(processed, inputSentinel, input)
	return
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
		err = fmt.Errorf("could not read unique patterns file. Please run --updatepatterns (%s)", err)
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
func (o *PatternsEntity) getFromFile(pathStr string) (pattern *Pattern, err error) {
	// Handle home directory expansion
	if strings.HasPrefix(pathStr, "~/") {
		var homedir string
		if homedir, err = os.UserHomeDir(); err != nil {
			err = fmt.Errorf("could not get home directory: %v", err)
			return
		}
		pathStr = filepath.Join(homedir, pathStr[2:])
	}

	var content []byte
	if content, err = os.ReadFile(pathStr); err != nil {
		err = fmt.Errorf("could not read pattern file %s: %v", pathStr, err)
		return
	}
	pattern = &Pattern{
		Name:    pathStr,
		Pattern: string(content),
	}
	return
}

// Get required for Storage interface
func (o *PatternsEntity) Get(name string) (*Pattern, error) {
	// Use GetPattern with no variables
	return o.GetApplyVariables(name, nil, "")
}
