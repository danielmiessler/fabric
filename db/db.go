package db

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"path/filepath"
	"time"
)

func NewDb(dir string) (db *Db) {

	db = &Db{Dir: dir}

	db.EnvFilePath = db.FilePath(".env")

	db.Patterns = &Patterns{
		Storage:                &Storage{Label: "Patterns", Dir: db.FilePath("patterns"), ItemIsDir: true},
		SystemPatternFile:      "system.md",
		UniquePatternsFilePath: db.FilePath("unique_patterns.txt"),
	}

	db.Sessions = &Sessions{
		&Storage{Label: "Sessions", Dir: db.FilePath("sessions"), FileExtension: ".json"}}

	db.Contexts = &Contexts{
		&Storage{Label: "Contexts", Dir: db.FilePath("contexts")}}

	return
}

type Db struct {
	Dir string

	Patterns *Patterns
	Sessions *Sessions
	Contexts *Contexts

	EnvFilePath string
}

func (o *Db) Configure() (err error) {
	if err = os.MkdirAll(o.Dir, os.ModePerm); err != nil {
		return
	}

	if err = o.LoadEnvFile(); err != nil {
		return
	}

	if err = o.Patterns.Configure(); err != nil {
		return
	}

	if err = o.Sessions.Configure(); err != nil {
		return
	}

	if err = o.Contexts.Configure(); err != nil {
		return
	}

	return
}

func (o *Db) LoadEnvFile() (err error) {
	if err = godotenv.Load(o.EnvFilePath); err != nil {
		err = fmt.Errorf("error loading .env file: %s", err)
	}
	return
}

func (o *Db) IsEnvFileExists() (ret bool) {
	_, err := os.Stat(o.EnvFilePath)
	ret = !os.IsNotExist(err)
	return
}

func (o *Db) SaveEnv(content string) (err error) {
	err = os.WriteFile(o.EnvFilePath, []byte(content), 0644)
	return
}

func (o *Db) FilePath(fileName string) (ret string) {
	return filepath.Join(o.Dir, fileName)
}

type DirectoryChange struct {
	Dir       string
	Timestamp time.Time
}
