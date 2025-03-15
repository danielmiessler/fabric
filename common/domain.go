package common

import goopenai "github.com/sashabaranov/go-openai"

const ChatMessageRoleMeta = "meta"

type ChatRequest struct {
	ContextName      string
	SessionName      string
	PatternName      string
	PatternVariables map[string]string
	Message          *goopenai.ChatCompletionMessage
	Language         string
	Meta             string
	InputHasVars     bool
	StrategyName     string
}

type ChatOptions struct {
	Model              string
	Temperature        float64
	TopP               float64
	PresencePenalty    float64
	FrequencyPenalty   float64
	Raw                bool
	Seed               int
	ModelContextLength int
}

// NormalizeMessages remove empty messages and ensure messages order user-assist-user
func NormalizeMessages(msgs []*goopenai.ChatCompletionMessage, defaultUserMessage string) (ret []*goopenai.ChatCompletionMessage) {
	// Iterate over messages to enforce the odd position rule for user messages
	fullMessageIndex := 0
	for _, message := range msgs {
		if message.Content == "" {
			// Skip empty messages as the anthropic API doesn't accept them
			continue
		}

		// Ensure, that each odd position shall be a user message
		if fullMessageIndex%2 == 0 && message.Role != goopenai.ChatMessageRoleUser {
			ret = append(ret, &goopenai.ChatCompletionMessage{Role: goopenai.ChatMessageRoleUser, Content: defaultUserMessage})
			fullMessageIndex++
		}
		ret = append(ret, message)
		fullMessageIndex++
	}
	return
}
