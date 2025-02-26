package common

import (
	"testing"

	goopenai "github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
)

func TestNormalizeMessages(t *testing.T) {
	msgs := []*goopenai.ChatCompletionMessage{
		{Role: goopenai.ChatMessageRoleUser, Content: "Hello"},
		{Role: goopenai.ChatMessageRoleAssistant, Content: "Hi there!"},
		{Role: goopenai.ChatMessageRoleUser, Content: ""},
		{Role: goopenai.ChatMessageRoleUser, Content: ""},
		{Role: goopenai.ChatMessageRoleUser, Content: "How are you?"},
	}

	expected := []*goopenai.ChatCompletionMessage{
		{Role: goopenai.ChatMessageRoleUser, Content: "Hello"},
		{Role: goopenai.ChatMessageRoleAssistant, Content: "Hi there!"},
		{Role: goopenai.ChatMessageRoleUser, Content: "How are you?"},
	}

	actual := NormalizeMessages(msgs, "default")
	assert.Equal(t, expected, actual)
}
