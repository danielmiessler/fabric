package db

import (
	"os"
)

type Contexts struct {
	*Storage
}

// LoadContext Load a context from file
func (o *Contexts) LoadContext(name string) (ret *Context, err error) {
	path := o.BuildFilePathByName(name)

	var content []byte
	if content, err = os.ReadFile(path); err != nil {
		return
	}

	ret = &Context{Name: name, Content: string(content)}
	return
}

type Context struct {
	Name    string
	Content string

	contexts *Contexts
}

// Save the session on disk
func (o *Context) Save() (err error) {
	err = o.contexts.Save(o.Name, []byte(o.Content))
	return err
}
