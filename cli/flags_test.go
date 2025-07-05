package cli

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
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
	options, err := flags.BuildChatOptions()
	assert.NoError(t, err)
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
	options, err := flags.BuildChatOptions()
	assert.NoError(t, err)
	assert.Equal(t, expectedOptions, options)
}

func TestInitWithYAMLConfig(t *testing.T) {
	// Create a temporary YAML config file
	configContent := `
temperature: 0.9
model: gpt-4
pattern: analyze
stream: true
`
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write([]byte(configContent)); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Test 1: Basic YAML loading
	t.Run("Load YAML config", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "--config", tmpfile.Name()}

		flags, err := Init()
		assert.NoError(t, err)
		assert.Equal(t, 0.9, flags.Temperature)
		assert.Equal(t, "gpt-4", flags.Model)
		assert.Equal(t, "analyze", flags.Pattern)
		assert.True(t, flags.Stream)
	})

	// Test 2: CLI overrides YAML
	t.Run("CLI overrides YAML", func(t *testing.T) {
		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "--config", tmpfile.Name(), "--temperature", "0.7", "--model", "gpt-3.5-turbo"}

		flags, err := Init()
		assert.NoError(t, err)
		assert.Equal(t, 0.7, flags.Temperature)
		assert.Equal(t, "gpt-3.5-turbo", flags.Model)
		assert.Equal(t, "analyze", flags.Pattern) // unchanged from YAML
		assert.True(t, flags.Stream)              // unchanged from YAML
	})

	// Test 3: Invalid YAML config
	t.Run("Invalid YAML config", func(t *testing.T) {
		badConfig := `
temperature: "not a float"
model: 123  # should be string
`
		badfile, err := os.CreateTemp("", "bad-config.*.yaml")
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(badfile.Name())

		if _, err := badfile.Write([]byte(badConfig)); err != nil {
			t.Fatal(err)
		}
		if err := badfile.Close(); err != nil {
			t.Fatal(err)
		}

		oldArgs := os.Args
		defer func() { os.Args = oldArgs }()
		os.Args = []string{"cmd", "--config", badfile.Name()}

		_, err = Init()
		assert.Error(t, err)
	})
}

func TestValidateImageFile(t *testing.T) {
	t.Run("Empty path should be valid", func(t *testing.T) {
		err := validateImageFile("")
		assert.NoError(t, err)
	})

	t.Run("Valid extensions should pass", func(t *testing.T) {
		validExtensions := []string{".png", ".jpeg", ".jpg", ".webp"}
		for _, ext := range validExtensions {
			filename := "/tmp/test" + ext
			err := validateImageFile(filename)
			assert.NoError(t, err, "Extension %s should be valid", ext)
		}
	})

	t.Run("Invalid extensions should fail", func(t *testing.T) {
		invalidExtensions := []string{".gif", ".bmp", ".tiff", ".svg", ".txt", ""}
		for _, ext := range invalidExtensions {
			filename := "/tmp/test" + ext
			err := validateImageFile(filename)
			assert.Error(t, err, "Extension %s should be invalid", ext)
			assert.Contains(t, err.Error(), "invalid image file extension")
		}
	})

	t.Run("Existing file should fail", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "test*.png")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		// Validation should fail because file exists
		err = validateImageFile(tempFile.Name())
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "image file already exists")
	})

	t.Run("Non-existing file with valid extension should pass", func(t *testing.T) {
		nonExistentFile := filepath.Join(os.TempDir(), "non_existent_file.png")
		// Make sure the file doesn't exist
		os.Remove(nonExistentFile)

		err := validateImageFile(nonExistentFile)
		assert.NoError(t, err)
	})
}

func TestBuildChatOptionsWithImageFileValidation(t *testing.T) {
	t.Run("Valid image file should pass", func(t *testing.T) {
		flags := &Flags{
			ImageFile: "/tmp/output.png",
		}

		options, err := flags.BuildChatOptions()
		assert.NoError(t, err)
		assert.Equal(t, "/tmp/output.png", options.ImageFile)
	})

	t.Run("Invalid extension should fail", func(t *testing.T) {
		flags := &Flags{
			ImageFile: "/tmp/output.gif",
		}

		options, err := flags.BuildChatOptions()
		assert.Error(t, err)
		assert.Nil(t, options)
		assert.Contains(t, err.Error(), "invalid image file extension")
	})

	t.Run("Existing file should fail", func(t *testing.T) {
		// Create a temporary file
		tempFile, err := os.CreateTemp("", "existing*.png")
		assert.NoError(t, err)
		defer os.Remove(tempFile.Name())
		tempFile.Close()

		flags := &Flags{
			ImageFile: tempFile.Name(),
		}

		options, err := flags.BuildChatOptions()
		assert.Error(t, err)
		assert.Nil(t, options)
		assert.Contains(t, err.Error(), "image file already exists")
	})
}
