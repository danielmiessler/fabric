package db

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/danielmiessler/fabric/common"
)

type Sessions struct {
	*Storage
}

func (o *Sessions) LoadOrCreateSession(name string) (ret *Session, err error) {
	if name == "" {
		return &Session{}, nil
	}

	path := o.BuildFilePath(name)
	if _, statErr := os.Stat(path); errors.Is(statErr, os.ErrNotExist) {
		fmt.Printf("Creating new session: %s\n", name)
		ret = &Session{Name: name, sessions: o}
	} else {
		ret, err = o.loadSession(name)
	}
	return
}

// LoadSession Load a session from file
func (o *Sessions) LoadSession(name string) (ret *Session, err error) {
	if name == "" {
		return &Session{}, nil
	}
	ret, err = o.loadSession(name)
	return
}

func (o *Sessions) loadSession(name string) (ret *Session, err error) {
	ret = &Session{Name: name, sessions: o}
	if err = o.LoadAsJson(name, &ret.Messages); err != nil {
		return
	}
	return
}

type Session struct {
	Name     string
	Messages []*common.Message

	sessions *Sessions
}

func (o *Session) Append(messages ...*common.Message) {
	o.Messages = append(o.Messages, messages...)
}

// Save the session on disk
func (o *Session) Save() (err error) {
	var jsonBytes []byte
	if jsonBytes, err = json.Marshal(o.Messages); err == nil {
		err = o.sessions.Save(o.Name, jsonBytes)
	} else {
		err = fmt.Errorf("could not marshal session %o: %o", o.Name, err)
	}
	return
}
