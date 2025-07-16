package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/danielmiessler/fabric/internal/core"
	"github.com/danielmiessler/fabric/internal/domain"
	"github.com/danielmiessler/fabric/internal/plugins/db/fsdb"
)

// handleChatProcessing handles the main chat processing logic
func handleChatProcessing(currentFlags *Flags, registry *core.PluginRegistry, messageTools string) (err error) {
	if messageTools != "" {
		currentFlags.AppendMessage(messageTools)
	}

	var chatter *core.Chatter
	if chatter, err = registry.GetChatter(currentFlags.Model, currentFlags.ModelContextLength,
		currentFlags.Strategy, currentFlags.Stream, currentFlags.DryRun); err != nil {
		return
	}

	var session *fsdb.Session
	var chatReq *domain.ChatRequest
	if chatReq, err = currentFlags.BuildChatRequest(strings.Join(os.Args[1:], " ")); err != nil {
		return
	}

	if chatReq.Language == "" {
		chatReq.Language = registry.Language.DefaultLanguage.Value
	}
	var chatOptions *domain.ChatOptions
	if chatOptions, err = currentFlags.BuildChatOptions(); err != nil {
		return
	}
	if session, err = chatter.Send(chatReq, chatOptions); err != nil {
		return
	}

	result := session.GetLastMessage().Content

	if !currentFlags.Stream || currentFlags.SuppressThink {
		// print the result if it was not streamed already or suppress-think disabled streaming output
		fmt.Println(result)
	}

	// if the copy flag is set, copy the message to the clipboard
	if currentFlags.Copy {
		if err = CopyToClipboard(result); err != nil {
			return
		}
	}

	// if the output flag is set, create an output file
	if currentFlags.Output != "" {
		if currentFlags.OutputSession {
			sessionAsString := session.String()
			err = CreateOutputFile(sessionAsString, currentFlags.Output)
		} else {
			err = CreateOutputFile(result, currentFlags.Output)
		}
	}
	return
}
