package db

import (
	"testing"

	"github.com/danielmiessler/fabric/common"
)

func TestSessions_GetOrCreateSession(t *testing.T) {
	dir := t.TempDir()
	sessions := &Sessions{
		Storage: &Storage{Dir: dir, FileExtension: ".json"},
	}
	sessionName := "testSession"
	session, err := sessions.GetOrCreateSession(sessionName)
	if err != nil {
		t.Fatalf("failed to get or create session: %v", err)
	}
	if session.Name != sessionName {
		t.Errorf("expected session name %v, got %v", sessionName, session.Name)
	}
}

func TestSessions_SaveSession(t *testing.T) {
	dir := t.TempDir()
	sessions := &Sessions{
		Storage: &Storage{Dir: dir, FileExtension: ".json"},
	}
	sessionName := "testSession"
	session := &Session{Name: sessionName, Messages: []*common.Message{{Content: "message1"}}}
	err := sessions.SaveSession(session)
	if err != nil {
		t.Fatalf("failed to save session: %v", err)
	}
	if !sessions.Exists(sessionName) {
		t.Errorf("expected session to be saved")
	}
}
