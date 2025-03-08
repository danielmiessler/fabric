package core

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/template"
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

type FileChangeItem struct {
	Type     string           `json:"type"`
	Name     string           `json:"name"`
	Content  string           `json:"content,omitempty"`
	Contents []FileChangeItem `json:"contents,omitempty"`
}

// validateFileChanges performs security checks on proposed file changes
func validateFileChanges(currentDir string, fileChanges []FileChangeItem) error {
	// Normalize the current directory to an absolute path
	currentDir, err := filepath.Abs(currentDir)
	if err != nil {
		return fmt.Errorf("could not resolve current directory: %v", err)
	}

	// Validate each file change
	for _, change := range fileChanges {
		// Skip non-file changes
		if change.Type != "file" {
			continue
		}

		// Validate file path
		filename := change.Name
		if filename == "" {
			return fmt.Errorf("file change with empty filename is not allowed")
		}

		// Convert to absolute path
		absPath, err := filepath.Abs(filepath.Join(currentDir, filename))
		if err != nil {
			return fmt.Errorf("invalid file path: %v", err)
		}

		// Security check: Ensure the file is within the current directory
		if !strings.HasPrefix(absPath, currentDir) {
			return fmt.Errorf("file path %s is outside the current directory", filename)
		}

		// Additional security checks
		// 1. Prevent absolute paths
		if filepath.IsAbs(filename) {
			return fmt.Errorf("absolute paths are not allowed: %s", filename)
		}

		// 2. Block path traversal attempts
		if strings.Contains(filename, "..") {
			return fmt.Errorf("path traversal attempt detected: %s", filename)
		}

	}

	return nil
}

// applyFileChanges applies validated file changes
func applyFileChanges(currentDir string, fileChanges []FileChangeItem) error {
	// Validate changes first
	if err := validateFileChanges(currentDir, fileChanges); err != nil {
		return fmt.Errorf("file change validation failed: %v", err)
	}

	// Apply changes
	for _, change := range fileChanges {
		if change.Type != "file" {
			continue
		}

		// Construct full file path
		filePath := filepath.Join(currentDir, change.Name)

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			return fmt.Errorf("could not create directory for %s: %v", filePath, err)
		}

		// Write file content
		err := os.WriteFile(filePath, []byte(change.Content), 0644)
		if err != nil {
			return fmt.Errorf("could not write file %s: %v", filePath, err)
		}

		// Log the file that was written
		log.Printf("INFO: File written: %s (size: %d bytes)", filePath, len(change.Content))
	}

	return nil
}

// extractFileChangesRecursive recursively extracts file changes from a nested directory structure
func extractFileChangesRecursive(item FileChangeItem) []FileChangeItem {
	var fileChanges []FileChangeItem

	// If this is a file, add it directly
	if item.Type == "file" {
		fileChanges = append(fileChanges, item)
	}

	// If this is a directory, recursively process its contents
	if item.Type == "directory" && len(item.Contents) > 0 {
		for _, subItem := range item.Contents {
			fileChanges = append(fileChanges, extractFileChangesRecursive(subItem)...)
		}
	}

	return fileChanges
}

func extractFileChanges(message string) ([]FileChangeItem, error) {
	const fileChangesMarker = "FILE_MANAGER_INTERFACE_CHANGES:"

	startIndex := strings.Index(message, fileChangesMarker)
	if startIndex == -1 {
		return nil, nil
	}
	log.Printf("DEBUG: Extracting file changes (marker: %s)", fileChangesMarker)
	log.Printf("DEBUG: Message length: %d characters", len(message))

	// Extract the file changes section
	changesText := message[startIndex+len(fileChangesMarker):]

	// Trim whitespace and remove code block markers
	changesText = strings.TrimSpace(changesText)
	changesText = strings.TrimPrefix(changesText, "```json")
	changesText = strings.TrimPrefix(changesText, "```text")
	changesText = strings.TrimPrefix(changesText, "```")
	changesText = strings.TrimSuffix(changesText, "```")
	changesText = strings.TrimSpace(changesText)

	// Define a struct that matches the exact JSON structure
	var topLevelChanges []FileChangeItem

	// Parse the JSON
	err := json.Unmarshal([]byte(changesText), &topLevelChanges)
	if err != nil {
		log.Printf("ERROR: Problematic JSON text: %s", changesText)
		return nil, fmt.Errorf("error parsing file changes JSON: %v", err)
	}

	// Collect all file change items recursively
	var fileChanges []FileChangeItem
	for _, topLevel := range topLevelChanges {
		fileChanges = append(fileChanges, extractFileChangesRecursive(topLevel)...)
	}

	// Ensure we have file changes
	if len(fileChanges) == 0 {
		return nil, fmt.Errorf("no file changes found in JSON")
	}

	return fileChanges, nil
}

func (o *Chatter) Send(request *common.ChatRequest, opts *common.ChatOptions) (session *fsdb.Session, err error) {
	if session, err = o.BuildSession(request, opts.Raw); err != nil {
		return
	}

	vendorMessages := session.GetVendorMessages()
	if len(vendorMessages) == 0 {
		if session.Name != "" {
			err = o.db.Sessions.SaveSession(session)
		}
		err = fmt.Errorf("no messages provided")
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

	// Extract file changes from the message
	fileChanges, err := extractFileChanges(message)
	if err != nil {
		// Log the error but don't stop processing
		fmt.Printf("Error extracting file changes: %v\n", err)
	}

	// Apply file changes
	if len(fileChanges) > 0 {
		// Use the current working directory as the base path
		currentDir, dirErr := os.Getwd()
		if dirErr != nil {
			fmt.Printf("Error getting current directory: %v\n", dirErr)
		} else {
			if applyErr := applyFileChanges(currentDir, fileChanges); applyErr != nil {
				fmt.Printf("Error applying file changes: %v\n", applyErr)
			}
		}
	}

	session.Append(&goopenai.ChatCompletionMessage{Role: goopenai.ChatMessageRoleAssistant, Content: message})

	if session.Name != "" {
		err = o.db.Sessions.SaveSession(session)
	}
	return
}

func (o *Chatter) BuildSession(request *common.ChatRequest, raw bool) (session *fsdb.Session, err error) {
	// If a session name is provided, retrieve it from the database
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

	// if a context name is provided, retrieve it from the database
	var contextContent string
	if request.ContextName != "" {
		var ctx *fsdb.Context
		if ctx, err = o.db.Contexts.Get(request.ContextName); err != nil {
			err = fmt.Errorf("could not find context %s: %v", request.ContextName, err)
			return
		}
		contextContent = ctx.Content
	}

	// Process any template variables in the message content (user input)
	// Double curly braces {{variable}} indicate template substitution
	// Ensure we have a message before processing, other wise we'll get an error when we pass to pattern.go
	if request.Message == nil {
		request.Message = &goopenai.ChatCompletionMessage{
			Role:    goopenai.ChatMessageRoleUser,
			Content: " ",
		}
	}

	// Now we know request.Message is not nil, process template variables
	if request.InputHasVars {
		request.Message.Content, err = template.ApplyTemplate(request.Message.Content, request.PatternVariables, "")
		if err != nil {
			return nil, err
		}
	}

	var patternContent string
	if request.PatternName != "" {
		pattern, err := o.db.Patterns.GetApplyVariables(request.PatternName, request.PatternVariables, request.Message.Content)
		// pattern will now contain user input, and all variables will be resolved, or errored

		if err != nil {
			return nil, fmt.Errorf("could not get pattern %s: %v", request.PatternName, err)
		}
		patternContent = pattern.Pattern
	}

	systemMessage := strings.TrimSpace(contextContent) + strings.TrimSpace(patternContent)
	if request.Language != "" {
		systemMessage = fmt.Sprintf("%s. Please use the language '%s' for the output.", systemMessage, request.Language)
	}

	if raw {
		if request.Message != nil {
			if systemMessage != "" {
				request.Message.Content = systemMessage
				// system contains pattern which contains user input
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
