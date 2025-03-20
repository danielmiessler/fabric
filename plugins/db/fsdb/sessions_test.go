package fsdb

import (
	"testing"

	goopenai "github.com/sashabaranov/go-openai"
)

func TestSessions_GetOrCreateSession(t *testing.T) {
	dir := t.TempDir()
	sessions := &SessionsEntity{
		StorageEntity: &StorageEntity{Dir: dir, FileExtension: ".json"},
	}
	sessionName := "testSession"
	session, err := sessions.Get(sessionName)
	if err != nil {
		t.Fatalf("failed to get or create session: %v", err)
	}
	if session.Name != sessionName {
		t.Errorf("expected session name %v, got %v", sessionName, session.Name)
	}
}

func TestSessions_SaveSession(t *testing.T) {
	dir := t.TempDir()
	sessions := &SessionsEntity{
		StorageEntity: &StorageEntity{Dir: dir, FileExtension: ".json"},
	}
	sessionName := "testSession"
	session := &Session{Name: sessionName, Messages: []*goopenai.ChatCompletionMessage{{Content: "message1"}}}
	err := sessions.SaveSession(session)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}
	if !sessions.Exists(sessionName) {
		t.Errorf("expected session to be saved")
	}
}
