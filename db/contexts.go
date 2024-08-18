package db

type Contexts struct {
	*Storage
}

// GetContext Load a context from file
func (o *Contexts) GetContext(name string) (ret *Context, err error) {
	var content []byte
	if content, err = o.Load(name); err != nil {
		return
	}

	ret = &Context{Name: name, Content: string(content)}
	return
}

type Context struct {
	Name    string
	Content string
}
