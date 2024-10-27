package core

import (
	"context"
	"fmt"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
)

const NoSessionPatternUserMessages = "no session, pattern or user messages provided"

type Chatter struct {
	db *fsdb.Db

	Stream bool
	DryRun bool

	model              string
	modelContextLength int
	vendor             ai.Vendor
}

func (o *Chatter) Send(request *common.ChatRequest, opts *common.ChatOptions) (session *fsdb.Session, err error) {
	if session, err = o.BuildSession(request, opts.Raw); err != nil {
		return
	}

	if opts.Model == "" {
		opts.Model = o.model
	}

	if opts.ModelContextLength == 0 {
		opts.ModelContextLength = o.modelContextLength
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

	session.Append(&goopenai.ChatCompletionMessage{Role: goopenai.ChatMessageRoleAssistant, Content: message})

	if session.Name != "" {
		err = o.db.Sessions.SaveSession(session)
	}
	return
}

func (o *Chatter) BuildSession(request *common.ChatRequest, raw bool) (session *fsdb.Session, err error) {
	if request.SessionName != "" {
		var sess *fsdb.Session
		if sess, err = o.db.Sessions.Get(request.SessionName); err != nil {
			err = fmt.Errorf("could not find session %s: %v", request.SessionName, err)
			return
		}
		session = sess
	} else {
		session = &fsdb.Session{}
	}

	if request.Meta != "" {
		session.Append(&goopenai.ChatCompletionMessage{Role: common.ChatMessageRoleMeta, Content: request.Meta})
	}

	var contextContent string
	if request.ContextName != "" {
		var ctx *fsdb.Context
		if ctx, err = o.db.Contexts.Get(request.ContextName); err != nil {
			err = fmt.Errorf("could not find context %s: %v", request.ContextName, err)
			return
		}
		contextContent = ctx.Content
	}

	var patternContent string
	if request.PatternName != "" {
		var pattern *fsdb.Pattern
		if pattern, err = o.db.Patterns.GetApplyVariables(request.PatternName, request.PatternVariables); err != nil {
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

	if raw {
		if request.Message != nil {
			if systemMessage != "" {
				request.Message.Content = systemMessage + request.Message.Content
			}
		} else {
			if systemMessage != "" {
				request.Message = &goopenai.ChatCompletionMessage{Role: goopenai.ChatMessageRoleSystem, Content: systemMessage}
			}
		}
	} else {
		if systemMessage != "" {
			session.Append(&goopenai.ChatCompletionMessage{Role: goopenai.ChatMessageRoleSystem, Content: systemMessage})
		}
	}

	if request.Message != nil {
		session.Append(request.Message)
	}

	if session.IsEmpty() {
		session = nil
		err = fmt.Errorf(NoSessionPatternUserMessages)
	}
	return
}
