package cli

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/danielmiessler/fabric/common"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	args := []string{"--copy"}
	expectedFlags := &Flags{Copy: true}
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()
	os.Args = append([]string{"cmd"}, args...)

	flags, err := Init()
	assert.NoError(t, err)
	assert.Equal(t, expectedFlags.Copy, flags.Copy)
}

func TestReadStdin(t *testing.T) {
	input := "test input"
	stdin := io.NopCloser(strings.NewReader(input))
	// No need to cast stdin to *os.File, pass it as io.ReadCloser directly
	content, err := ReadStdin(stdin)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if content != input {
		t.Fatalf("expected %q, got %q", input, content)
	}
}

// ReadStdin function assuming it's part of `cli` package
func ReadStdin(reader io.ReadCloser) (string, error) {
	defer reader.Close()
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func TestBuildChatOptions(t *testing.T) {
	flags := &Flags{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
		Seed:             1,
	}

	expectedOptions := &common.ChatOptions{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
		Raw:              false,
		Seed:             1,
	}
	options := flags.BuildChatOptions()
	assert.Equal(t, expectedOptions, options)
}

func TestBuildChatOptionsDefaultSeed(t *testing.T) {
	flags := &Flags{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
	}

	expectedOptions := &common.ChatOptions{
		Temperature:      0.8,
		TopP:             0.9,
		PresencePenalty:  0.1,
		FrequencyPenalty: 0.2,
		Raw:              false,
		Seed:             0,
	}
	options := flags.BuildChatOptions()
	assert.Equal(t, expectedOptions, options)
}

func TestBuildChatRequest(t *testing.T) {
	flags := &Flags{
		Context: "test-context",
		Session: "test-session",
		Pattern: "test-pattern",
		Message: "test-message",
	}

	expectedRequest := &common.ChatRequest{
		ContextName: "test-context",
		SessionName: "test-session",
		PatternName: "test-pattern",
		Message:     "test-message",
		Meta:        "test",
	}
	request := flags.BuildChatRequest("test")
	assert.Equal(t, expectedRequest, request)
}
