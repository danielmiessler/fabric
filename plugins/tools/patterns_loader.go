package tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielmiessler/fabric/plugins"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/tools/githelper"

	"github.com/otiai10/copy"
)

const DefaultPatternsGitRepoUrl = "https://github.com/danielmiessler/fabric.git"
const DefaultPatternsGitRepoFolder = "patterns"

func NewPatternsLoader(patterns *fsdb.PatternsEntity) (ret *PatternsLoader) {
	label := "Patterns Loader"
	ret = &PatternsLoader{
		Patterns:       patterns,
		loadedFilePath: patterns.BuildFilePath("loaded"),
	}

	ret.PluginBase = &plugins.PluginBase{
		Name:             label,
		SetupDescription: "Patterns - Downloads patterns [required]",
		EnvNamePrefix:    plugins.BuildEnvVariablePrefix(label),
		ConfigureCustom:  ret.configure,
	}

	ret.DefaultGitRepoUrl = ret.AddSetupQuestionCustom("Git Repo Url", true,
		"Enter the default Git repository URL for the patterns")
	ret.DefaultGitRepoUrl.Value = DefaultPatternsGitRepoUrl

	ret.DefaultFolder = ret.AddSetupQuestionCustom("Git Repo Patterns Folder", true,
		"Enter the default folder in the Git repository where patterns are stored")
	ret.DefaultFolder.Value = DefaultPatternsGitRepoFolder

	return
}

type PatternsLoader struct {
	*plugins.PluginBase
	Patterns *fsdb.PatternsEntity

	DefaultGitRepoUrl *plugins.SetupQuestion
	DefaultFolder     *plugins.SetupQuestion

	loadedFilePath string

	pathPatternsPrefix string
	tempPatternsFolder string
}

func (o *PatternsLoader) configure() (err error) {
	o.pathPatternsPrefix = fmt.Sprintf("%v/", o.DefaultFolder.Value)
	o.tempPatternsFolder = filepath.Join(os.TempDir(), o.DefaultFolder.Value)

	return
}

func (o *PatternsLoader) IsConfigured() (ret bool) {
	ret = o.PluginBase.IsConfigured()
	if ret {
		if _, err := os.Stat(o.loadedFilePath); os.IsNotExist(err) {
			ret = false
		}
	}
	return
}

func (o *PatternsLoader) Setup() (err error) {
	if err = o.PluginBase.Setup(); err != nil {
		return
	}

	if err = o.PopulateDB(); err != nil {
		return
	}
	return
}

// PopulateDB downloads patterns from the internet and populates the patterns folder
func (o *PatternsLoader) PopulateDB() (err error) {
	fmt.Printf("Downloading patterns and Populating %s...\n", o.Patterns.Dir)
	fmt.Println()
	if err = o.gitCloneAndCopy(); err != nil {
		return
	}

	if err = o.movePatterns(); err != nil {
		return
	}
	return
}

// PersistPatterns copies custom patterns to the updated patterns directory
func (o *PatternsLoader) PersistPatterns() (err error) {
	var currentPatterns []os.DirEntry
	if currentPatterns, err = os.ReadDir(o.Patterns.Dir); err != nil {
		return
	}

	newPatternsFolder := o.tempPatternsFolder
	var newPatterns []os.DirEntry
	if newPatterns, err = os.ReadDir(newPatternsFolder); err != nil {
		return
	}

	for _, currentPattern := range currentPatterns {
		for _, newPattern := range newPatterns {
			if currentPattern.Name() == newPattern.Name() {
				break
			}
			err = copy.Copy(filepath.Join(o.Patterns.Dir, newPattern.Name()), filepath.Join(newPatternsFolder, newPattern.Name()))
		}
	}
	return
}

// movePatterns copies the new patterns into the config directory
func (o *PatternsLoader) movePatterns() (err error) {
	if err = os.MkdirAll(o.Patterns.Dir, os.ModePerm); err != nil {
		return
	}

	patternsDir := o.tempPatternsFolder
	if err = o.PersistPatterns(); err != nil {
		return
	}

	if err = copy.Copy(patternsDir, o.Patterns.Dir); err != nil { // copies the patterns to the config directory
		return
	}

	//create an empty file to indicate that the patterns have been updated if not exists
	_, _ = os.Create(o.loadedFilePath)

	err = os.RemoveAll(patternsDir)
	return
}

func (o *PatternsLoader) gitCloneAndCopy() (err error) {
	// Create temp folder if it doesn't exist
	if err = os.MkdirAll(filepath.Dir(o.tempPatternsFolder), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Use the helper to fetch files
	err = githelper.FetchFilesFromRepo(githelper.FetchOptions{
		RepoURL:    o.DefaultGitRepoUrl.Value,
		PathPrefix: o.DefaultFolder.Value,
		DestDir:    o.tempPatternsFolder,
	})
	if err != nil {
		return fmt.Errorf("failed to download patterns: %w", err)
	}

	return nil
}
