package core

import (
	"context"
	"fmt"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/db"
	"github.com/danielmiessler/fabric/vendors"
)

type Chatter struct {
	db *db.Db

	Stream bool
	DryRun bool

	model  string
	vendor vendors.Vendor
}

func (o *Chatter) Send(request *common.ChatRequest, opts *common.ChatOptions) (message string, err error) {
	var chatRequest *Chat
	if chatRequest, err = o.NewChat(request); err != nil {
		return
	}

	var session *db.Session
	if session, err = chatRequest.BuildChatSession(); err != nil {
		return
	}

	if opts.Model == "" {
		opts.Model = o.model
	}

	if o.Stream {
		channel := make(chan string)
		go func() {
			if streamErr := o.vendor.SendStream(session.Messages, opts, channel); streamErr != nil {
				channel <- streamErr.Error()
			}
		}()

		for response := range channel {
			message += response
			fmt.Print(response)
		}
	} else {
		if message, err = o.vendor.Send(context.Background(), session.Messages, opts); err != nil {
			return
		}
	}

	if chatRequest.Session != nil && message != "" {
		chatRequest.Session.Append(&common.Message{Role: "system", Content: message})
		err = o.db.Sessions.SaveSession(chatRequest.Session)
	}
	return
}

func (o *Chatter) NewChat(request *common.ChatRequest) (ret *Chat, err error) {
	ret = &Chat{}

	if request.ContextName != "" {
		var ctx *db.Context
		if ctx, err = o.db.Contexts.GetContext(request.ContextName); err != nil {
			err = fmt.Errorf("could not find context %s: %v", request.ContextName, err)
			return
		}
		ret.Context = ctx.Content
	}

	if request.SessionName != "" {
		var sess *db.Session
		if sess, err = o.db.Sessions.GetOrCreateSession(request.SessionName); err != nil {
			err = fmt.Errorf("could not find session %s: %v", request.SessionName, err)
			return
		}
		ret.Session = sess
	}

	if request.PatternName != "" {
		var pattern *db.Pattern
		if pattern, err = o.db.Patterns.GetPattern(request.PatternName, request.PatternVariables); err != nil {
			err = fmt.Errorf("could not find pattern %s: %v", request.PatternName, err)
			return
		}

		if pattern.Pattern != "" {
			ret.Pattern = pattern.Pattern
		}
	}

	ret.Message = request.Message
	return
}

type Chat struct {
	Context string
	Pattern string
	Message string
	Session *db.Session
}
