package core

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/db"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/otiai10/copy"
)

func NewPatternsLoader(patterns *db.Patterns) (ret *PatternsLoader) {
	label := "Patterns Loader"
	ret = &PatternsLoader{
		Patterns: patterns,
	}

	ret.Configurable = &common.Configurable{
		Label:           label,
		EnvNamePrefix:   common.BuildEnvVariablePrefix(label),
		ConfigureCustom: ret.configure,
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
	*common.Configurable
	Patterns *db.Patterns

	DefaultGitRepoUrl *common.SetupQuestion
	DefaultFolder     *common.SetupQuestion

	pathPatternsPrefix string
	tempPatternsFolder string
}

func (o *PatternsLoader) configure() (err error) {
	o.pathPatternsPrefix = fmt.Sprintf("%v/", o.DefaultFolder.Value)
	o.tempPatternsFolder = filepath.Join(os.TempDir(), o.DefaultFolder.Value)

	return
}

// PopulateDB downloads patterns from the internet and populates the patterns folder
func (o *PatternsLoader) PopulateDB() (err error) {
	fmt.Printf("Downloading patterns and Populating %s..\n", o.Patterns.Dir)
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
			copy.Copy(filepath.Join(o.Patterns.Dir, newPattern.Name()), filepath.Join(newPatternsFolder, newPattern.Name()))
		}
	}
	return
}

// movePatterns copies the new patterns into the config directory
func (o *PatternsLoader) movePatterns() (err error) {
	os.MkdirAll(o.Patterns.Dir, os.ModePerm)

	patternsDir := o.tempPatternsFolder
	if err = o.PersistPatterns(); err != nil {
		return
	}

	copy.Copy(patternsDir, o.Patterns.Dir) // copies the patterns to the config directory
	err = os.RemoveAll(patternsDir)
	return
}

// checks if a pattern already exists in the directory
// func DoesPatternExistAlready(name string) (bool, error) {
// 	entry := db.Entry{
// 		Label: name,
// 	}
// 	_, err := entry.GetByName()
// 	if err != nil {
// 		return false, err
// 	}
// 	return true, nil
// }

func (o *PatternsLoader) gitCloneAndCopy() (err error) {
	// Clones the given repository, creating the remote, the local branches
	// and fetching the objects, everything in memory:
	var r *git.Repository
	if r, err = git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: o.DefaultGitRepoUrl.Value,
	}); err != nil {
		fmt.Println(err)
		return
	}

	// ... retrieves the branch pointed by HEAD
	var ref *plumbing.Reference
	if ref, err = r.Head(); err != nil {
		fmt.Println(err)
		return
	}

	// ... retrieves the commit history for /patterns folder
	var cIter object.CommitIter
	if cIter, err = r.Log(&git.LogOptions{
		From: ref.Hash(),
		PathFilter: func(path string) bool {
			return path == o.DefaultFolder.Value || strings.HasPrefix(path, o.pathPatternsPrefix)
		},
	}); err != nil {
		fmt.Println(err)
		return err
	}

	var changes []db.DirectoryChange
	// ... iterates over the commits
	if err = cIter.ForEach(func(c *object.Commit) (err error) {
		// Get the files changed in this commit by comparing with its parents
		parentIter := c.Parents()
		if err = parentIter.ForEach(func(parent *object.Commit) (err error) {
			var patch *object.Patch
			if patch, err = parent.Patch(c); err != nil {
				fmt.Println(err)
				return
			}

			for _, fileStat := range patch.Stats() {
				if strings.HasPrefix(fileStat.Name, o.pathPatternsPrefix) {
					dir := filepath.Dir(fileStat.Name)
					changes = append(changes, db.DirectoryChange{Dir: dir, Timestamp: c.Committer.When})
				}
			}
			return
		}); err != nil {
			fmt.Println(err)
			return
		}
		return
	}); err != nil {
		fmt.Println(err)
		return
	}

	// Sort changes by timestamp
	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Timestamp.Before(changes[j].Timestamp)
	})

	o.makeUniqueList(changes)

	var commit *object.Commit
	if commit, err = r.CommitObject(ref.Hash()); err != nil {
		fmt.Println(err)
		return
	}

	var tree *object.Tree
	if tree, err = commit.Tree(); err != nil {
		fmt.Println(err)
		return
	}

	if err = tree.Files().ForEach(func(f *object.File) (err error) {
		if strings.HasPrefix(f.Name, o.pathPatternsPrefix) {
			// Create the local file path
			localPath := filepath.Join(os.TempDir(), f.Name)

			// Create the directories if they don't exist
			if err = os.MkdirAll(filepath.Dir(localPath), os.ModePerm); err != nil {
				fmt.Println(err)
				return
			}

			// Write the file to the local filesystem
			var blob *object.Blob
			if blob, err = r.BlobObject(f.Hash); err != nil {
				fmt.Println(err)
				return
			}
			err = o.writeBlobToFile(blob, localPath)
			return
		}

		return
	}); err != nil {
		fmt.Println(err)
	}

	return
}

func (o *PatternsLoader) writeBlobToFile(blob *object.Blob, path string) (err error) {
	var reader io.ReadCloser
	if reader, err = blob.Reader(); err != nil {
		return
	}
	defer reader.Close()

	// Create the file
	var file *os.File
	if file, err = os.Create(path); err != nil {
		return
	}
	defer file.Close()

	// Copy the contents of the blob to the file
	if _, err = io.Copy(file, reader); err != nil {
		return
	}
	return
}

func (o *PatternsLoader) makeUniqueList(changes []db.DirectoryChange) {
	uniqueItems := make(map[string]bool)
	for _, change := range changes {
		if strings.TrimSpace(change.Dir) != "" && !strings.Contains(change.Dir, "=>") {
			pattern := strings.ReplaceAll(change.Dir, o.pathPatternsPrefix, "")
			pattern = strings.TrimSpace(pattern)
			uniqueItems[pattern] = true
		}
	}

	finalList := make([]string, 0, len(uniqueItems))
	for _, change := range changes {
		pattern := strings.ReplaceAll(change.Dir, o.pathPatternsPrefix, "")
		pattern = strings.TrimSpace(pattern)
		if _, exists := uniqueItems[pattern]; exists {
			finalList = append(finalList, pattern)
			delete(uniqueItems, pattern) // Remove to avoid duplicates in the final list
		}
	}

	joined := strings.Join(finalList, "\n")
	os.WriteFile(o.Patterns.UniquePatternsFilePath, []byte(joined), 0o644)
}
