package db

import (
	"os"
	"path/filepath"
	"testing"
)

func TestContexts_GetContext(t *testing.T) {
	dir := t.TempDir()
	contexts := &Contexts{
		Storage: &Storage{Dir: dir},
	}
	contextName := "testContext"
	contextPath := filepath.Join(dir, contextName)
	contextContent := "test content"
	err := os.WriteFile(contextPath, []byte(contextContent), 0644)
	if err != nil {
		t.Fatalf("failed to write context file: %v", err)
	}
	context, err := contexts.GetContext(contextName)
	if err != nil {
		t.Fatalf("failed to get context: %v", err)
	}
	expectedContext := &Context{Name: contextName, Content: contextContent}
	if *context != *expectedContext {
		t.Errorf("expected %v, got %v", expectedContext, context)
	}
}
