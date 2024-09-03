package common

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	ContextName      string
	SessionName      string
	PatternName      string
	PatternVariables map[string]string
	Message          string
}

type ChatOptions struct {
	Model            string
	Temperature      float64
	TopP             float64
	PresencePenalty  float64
	FrequencyPenalty float64
}

// NormalizeMessages remove empty messages and ensure messages order user-assist-user
func NormalizeMessages(msgs []*Message, defaultUserMessage string) (ret []*Message) {
	// Iterate over messages to enforce the odd position rule for user messages
	fullMessageIndex := 0
	for _, message := range msgs {
		if message.Content == "" {
			// Skip empty messages as the anthropic API doesn't accept them
			continue
		}

		// Ensure, that each odd position shall be a user message
		if fullMessageIndex%2 == 0 && message.Role != "user" {
			ret = append(ret, &Message{Role: "user", Content: defaultUserMessage})
			fullMessageIndex++
		}
		ret = append(ret, message)
		fullMessageIndex++
	}
	return
}
