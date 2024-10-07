package core

import (
	"context"
	"fmt"
	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/db"
	"github.com/danielmiessler/fabric/vendors"
	goopenai "github.com/sashabaranov/go-openai"
	"strings"
)

type Chatter struct {
	db *db.Db

	Stream bool
	DryRun bool

	model  string
	vendor vendors.Vendor
}

func (o *Chatter) Send(request *common.ChatRequest, opts *common.ChatOptions) (session *db.Session, err error) {
	if session, err = o.BuildSession(request, opts.Raw); err != nil {
		return
	}

	if opts.Model == "" {
		opts.Model = o.model
	}

	message := ""

	if o.Stream {
		channel := make(chan string)
		go func() {
			if streamErr := o.vendor.SendStream(session.GetVendorMessages(), opts, channel); streamErr != nil {
				channel <- streamErr.Error()
			}
		}()

		for response := range channel {
			message += response
			fmt.Print(response)
		}
	} else {
		if message, err = o.vendor.Send(context.Background(), session.GetVendorMessages(), opts); err != nil {
			return
		}
	}

	if message == "" {
		session = nil
		err = fmt.Errorf("empty response")
		return
	}

	session.Append(&common.Message{Role: goopenai.ChatMessageRoleAssistant, Content: message})

	if session.Name != "" {
		err = o.db.Sessions.SaveSession(session)
	}
	return
}

func (o *Chatter) BuildSession(request *common.ChatRequest, raw bool) (session *db.Session, err error) {
	if request.SessionName != "" {
		var sess *db.Session
		if sess, err = o.db.Sessions.GetOrCreateSession(request.SessionName); err != nil {
			err = fmt.Errorf("could not find session %s: %v", request.SessionName, err)
			return
		}
		session = sess
	} else {
		session = &db.Session{}
	}

	if request.Meta != "" {
		session.Append(&common.Message{Role: common.ChatMessageRoleMeta, Content: request.Meta})
	}

	var contextContent string
	if request.ContextName != "" {
		var ctx *db.Context
		if ctx, err = o.db.Contexts.GetContext(request.ContextName); err != nil {
			err = fmt.Errorf("could not find context %s: %v", request.ContextName, err)
			return
		}
		contextContent = ctx.Content
	}

	var patternContent string
	if request.PatternName != "" {
		var pattern *db.Pattern
		if pattern, err = o.db.Patterns.GetPattern(request.PatternName, request.PatternVariables); err != nil {
			err = fmt.Errorf("could not find pattern %s: %v", request.PatternName, err)
			return
		}

		if pattern.Pattern != "" {
			patternContent = pattern.Pattern
		}
	}

	systemMessage := strings.TrimSpace(contextContent) + strings.TrimSpace(patternContent)
	if request.Language != "" {
		systemMessage = fmt.Sprintf("%s. Please use the language '%s' for the output.", systemMessage, request.Language)
	}
	userMessage := strings.TrimSpace(request.Message)

	if raw {
		// use the user role instead of the system role in raw mode
		message := systemMessage + userMessage
		if message != "" {
			session.Append(&common.Message{Role: goopenai.ChatMessageRoleUser, Content: message})
		}
	} else {
		if systemMessage != "" {
			session.Append(&common.Message{Role: goopenai.ChatMessageRoleSystem, Content: systemMessage})
		}
		if userMessage != "" {
			session.Append(&common.Message{Role: goopenai.ChatMessageRoleUser, Content: userMessage})
		}
	}

	if session.IsEmpty() {
		session = nil
		err = fmt.Errorf(NoSessionPatternUserMessages)
	}
	return
}
