package db

import (
	"encoding/json"
	"fmt"
	"github.com/samber/lo"
	"os"
	"path/filepath"
	"strings"
)

type Storage struct {
	Label         string
	Dir           string
	ItemIsDir     bool
	FileExtension string
}

func (o *Storage) Configure() (err error) {
	if err = os.MkdirAll(o.Dir, os.ModePerm); err != nil {
		return
	}
	return
}

// GetNames finds all patterns in the patterns directory and enters the id, name, and pattern into a slice of Entry structs. it returns these entries or an error
func (o *Storage) GetNames() (ret []string, err error) {
	var entries []os.DirEntry
	if entries, err = os.ReadDir(o.Dir); err != nil {
		err = fmt.Errorf("could not read items from directory: %v", err)
		return
	}

	if o.ItemIsDir {
		ret = lo.FilterMap(entries, func(item os.DirEntry, index int) (ret string, ok bool) {
			if ok = item.IsDir(); ok {
				ret = item.Name()
			}
			return
		})
	} else {
		if o.FileExtension == "" {
			ret = lo.FilterMap(entries, func(item os.DirEntry, index int) (ret string, ok bool) {
				if ok = !item.IsDir(); ok {
					ret = item.Name()
				}
				return
			})
		} else {
			ret = lo.FilterMap(entries, func(item os.DirEntry, index int) (ret string, ok bool) {
				if ok = !item.IsDir() && filepath.Ext(item.Name()) == o.FileExtension; ok {
					ret = strings.TrimSuffix(item.Name(), o.FileExtension)
				}
				return
			})
		}
	}
	return
}

func (o *Storage) ListNames() (err error) {
	var names []string
	if names, err = o.GetNames(); err != nil {
		return
	}

	if len(names) == 0 {
		fmt.Printf("\nNo %v\n", o.Label)
		return
	}

	fmt.Printf("\n%v:\n", o.Label)
	for _, item := range names {
		fmt.Printf("\t%s\n", item)
	}
	return
}

func (o *Storage) BuildFilePathByName(name string) (ret string) {
	ret = o.BuildFilePath(o.buildFileName(name))
	return
}

func (o *Storage) BuildFilePath(fileName string) (ret string) {
	ret = filepath.Join(o.Dir, fileName)
	return
}

func (o *Storage) buildFileName(name string) string {
	return fmt.Sprintf("%s%v", name, o.FileExtension)
}

func (o *Storage) Delete(name string) (err error) {
	if err = os.Remove(o.BuildFilePathByName(name)); err != nil {
		err = fmt.Errorf("could not delete %s: %v", name, err)
	}
	return
}

func (o *Storage) Exists(name string) (ret bool) {
	_, err := os.Stat(o.BuildFilePathByName(name))
	ret = !os.IsNotExist(err)
	return
}

func (o *Storage) Rename(oldName, newName string) (err error) {
	if err = os.Rename(o.BuildFilePathByName(oldName), o.BuildFilePathByName(newName)); err != nil {
		err = fmt.Errorf("could not rename %s to %s: %v", oldName, newName, err)
	}
	return
}

func (o *Storage) Save(name string, content []byte) (err error) {
	if err = os.WriteFile(o.BuildFilePathByName(name), content, 0644); err != nil {
		err = fmt.Errorf("could not save %s: %v", name, err)
	}
	return
}

func (o *Storage) Load(name string) (ret []byte, err error) {
	if ret, err = os.ReadFile(o.BuildFilePathByName(name)); err != nil {
		err = fmt.Errorf("could not load %s: %v", name, err)
	}
	return
}

func (o *Storage) SaveAsJson(name string, item interface{}) (err error) {
	var jsonString []byte
	if jsonString, err = json.Marshal(item); err == nil {
		err = o.Save(name, jsonString)
	} else {
		err = fmt.Errorf("could not marshal %s: %s", name, err)
	}

	return err
}

func (o *Storage) LoadAsJson(name string, item interface{}) (err error) {
	var content []byte
	if content, err = o.Load(name); err != nil {
		return
	}

	if err = json.Unmarshal(content, &item); err != nil {
		err = fmt.Errorf("could not unmarshal %s: %s", name, err)
	}
	return
}
