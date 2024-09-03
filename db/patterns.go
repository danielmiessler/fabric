package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Patterns struct {
	*Storage
	SystemPatternFile      string
	UniquePatternsFilePath string
}

// GetPattern finds a pattern by name and returns the pattern as an entry or an error
func (o *Patterns) GetPattern(name string, variables map[string]string) (ret *Pattern, err error) {
	patternPath := filepath.Join(o.Dir, name, o.SystemPatternFile)

	var pattern []byte
	if pattern, err = os.ReadFile(patternPath); err != nil {
		return
	}

	patternStr := string(pattern)

	if variables != nil && len(variables) > 0 {
		for variableName, value := range variables {
			patternStr = strings.ReplaceAll(patternStr, variableName, value)
		}
	}

	ret = &Pattern{
		Name:    name,
		Pattern: patternStr,
	}
	return
}

func (o *Patterns) PrintLatestPatterns(latestNumber int) (err error) {
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

type Pattern struct {
	Name        string
	Description string
	Pattern     string
}
