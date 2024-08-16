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

// GetByName finds a pattern by name and returns the pattern as an entry or an error
func (o *Patterns) GetByName(name string) (ret *Pattern, err error) {
	patternPath := filepath.Join(o.Dir, name, o.SystemPatternFile)

	var pattern []byte
	if pattern, err = os.ReadFile(patternPath); err != nil {
		return
	}
	ret = &Pattern{
		Name:    name,
		Pattern: string(pattern),
	}
	return
}

func (o *Patterns) LatestPatterns(latestNumber int) (err error) {
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
