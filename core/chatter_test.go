package core

import (
	"testing"
)

func TestBuildChatSession(t *testing.T) {
	chat := &Chat{
		Context: "test context",
		Pattern: "test pattern",
		Message: "test message",
	}
	session, err := chat.BuildChatSession()
	if err != nil {
		t.Fatalf("BuildChatSession() error = %v", err)
	}

	if session == nil {
		t.Fatalf("BuildChatSession() returned nil session")
	}
}
