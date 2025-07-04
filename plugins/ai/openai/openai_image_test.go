package openai

import (
	"testing"

	"github.com/danielmiessler/fabric/chat"
	"github.com/danielmiessler/fabric/common"
	"github.com/openai/openai-go/responses"
	"github.com/stretchr/testify/assert"
)

func TestShouldUseImageGeneration(t *testing.T) {
	client := NewClient()

	// Test with image file specified
	opts := &common.ChatOptions{
		ImageFile: "output.png",
	}
	assert.True(t, client.shouldUseImageGeneration(opts), "Should use image generation when image file is specified")

	// Test without image file
	opts = &common.ChatOptions{
		ImageFile: "",
	}
	assert.False(t, client.shouldUseImageGeneration(opts), "Should not use image generation when no image file is specified")
}

func TestAddImageGenerationTool(t *testing.T) {
	client := NewClient()

	// Test with image generation enabled
	opts := &common.ChatOptions{
		ImageFile: "output.png",
	}
	tools := []responses.ToolUnionParam{}
	result := client.addImageGenerationTool(opts, tools)

	assert.Len(t, result, 1, "Should add one image generation tool")
	assert.NotNil(t, result[0].OfImageGeneration, "Should have image generation tool")
	assert.Equal(t, "image_generation", string(result[0].OfImageGeneration.Type))
	assert.Equal(t, "gpt-image-1", result[0].OfImageGeneration.Model)
	assert.Equal(t, "png", result[0].OfImageGeneration.OutputFormat)

	// Test without image generation
	opts = &common.ChatOptions{
		ImageFile: "",
	}
	tools = []responses.ToolUnionParam{}
	result = client.addImageGenerationTool(opts, tools)

	assert.Len(t, result, 0, "Should not add image generation tool when not needed")
}

func TestBuildResponseParams_WithImageGeneration(t *testing.T) {
	client := NewClient()
	opts := &common.ChatOptions{
		Model:     "gpt-image-1",
		ImageFile: "output.png",
	}

	msgs := []*chat.ChatCompletionMessage{
		{Role: "user", Content: "Generate an image of a cat"},
	}

	params := client.buildResponseParams(msgs, opts)

	assert.NotNil(t, params.Tools, "Expected tools when image generation is enabled")

	// Should have image generation tool
	hasImageTool := false
	for _, tool := range params.Tools {
		if tool.OfImageGeneration != nil {
			hasImageTool = true
			assert.Equal(t, "image_generation", string(tool.OfImageGeneration.Type))
			assert.Equal(t, "gpt-image-1", tool.OfImageGeneration.Model)
			break
		}
	}
	assert.True(t, hasImageTool, "Should have image generation tool")
}

func TestBuildResponseParams_WithBothSearchAndImage(t *testing.T) {
	client := NewClient()
	opts := &common.ChatOptions{
		Model:          "gpt-image-1",
		Search:         true,
		SearchLocation: "America/Los_Angeles",
		ImageFile:      "output.png",
	}

	msgs := []*chat.ChatCompletionMessage{
		{Role: "user", Content: "Search for cat images and generate one"},
	}

	params := client.buildResponseParams(msgs, opts)

	assert.NotNil(t, params.Tools, "Expected tools when both search and image generation are enabled")
	assert.Len(t, params.Tools, 2, "Should have both search and image generation tools")

	hasSearchTool := false
	hasImageTool := false

	for _, tool := range params.Tools {
		if tool.OfWebSearchPreview != nil {
			hasSearchTool = true
		}
		if tool.OfImageGeneration != nil {
			hasImageTool = true
		}
	}

	assert.True(t, hasSearchTool, "Should have web search tool")
	assert.True(t, hasImageTool, "Should have image generation tool")
}
