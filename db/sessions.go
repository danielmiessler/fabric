package db

import (
	"fmt"
	"github.com/danielmiessler/fabric/common"
)

type Sessions struct {
	*Storage
}

func (o *Sessions) GetOrCreateSession(name string) (session *Session, err error) {
	session = &Session{Name: name}

	if o.Exists(name) {
		err = o.LoadAsJson(name, &session.Messages)
	} else {
		fmt.Printf("Creating new session: %s\n", name)
	}
	return
}

func (o *Sessions) PrintSession(name string) (err error) {
	if o.Exists(name) {
		var session Session
		if err = o.LoadAsJson(name, &session.Messages); err == nil {
			fmt.Println(session)
		}
	}
	return
}

func (o *Sessions) SaveSession(session *Session) (err error) {
	return o.SaveAsJson(session.Name, session.Messages)
}

type Session struct {
	Name     string
	Messages []*common.Message
}

func (o *Session) IsEmpty() bool {
	return len(o.Messages) == 0
}

func (o *Session) Append(messages ...*common.Message) {
	o.Messages = append(o.Messages, messages...)
}

func (o *Session) GetLastMessage() (ret *common.Message) {
	if len(o.Messages) > 0 {
		ret = o.Messages[len(o.Messages)-1]
	}
	return
}

func (o *Session) String() (ret string) {
	for _, message := range o.Messages {
		ret += fmt.Sprintf("[%v] >\n\n%v\n\n", message.Role, message.Content)
	}
	return
}
