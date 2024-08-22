package common

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeMessages(t *testing.T) {
	msgs := []*Message{
		{Role: "user", Content: "Hello"},
		{Role: "bot", Content: "Hi there!"},
		{Role: "bot", Content: ""},
		{Role: "user", Content: ""},
		{Role: "user", Content: "How are you?"},
	}

	expected := []*Message{
		{Role: "user", Content: "Hello"},
		{Role: "bot", Content: "Hi there!"},
		{Role: "user", Content: "How are you?"},
	}

	actual := NormalizeMessages(msgs, "default")
	assert.Equal(t, expected, actual)
}
