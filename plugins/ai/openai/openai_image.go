package openai

// This file contains helper methods for image generation and processing
// using OpenAI's Responses API and Image API.

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/danielmiessler/fabric/common"
	"github.com/openai/openai-go/responses"
)

// ImageGenerationResponseType is the type used for image generation calls in responses
const ImageGenerationResponseType = "image_generation_call"
const ImageGenerationToolType = "image_generation"

// getOutputFormatFromExtension determines the API output format based on file extension
func getOutputFormatFromExtension(imagePath string) string {
	if imagePath == "" {
		return "png" // Default format
	}

	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".png":
		return "png"
	case ".webp":
		return "webp"
	case ".jpg":
		return "jpeg"
	case ".jpeg":
		return "jpeg"
	default:
		return "png" // Default fallback
	}
}

// addImageGenerationTool adds the image generation tool to the request if needed
func (o *Client) addImageGenerationTool(opts *common.ChatOptions, tools []responses.ToolUnionParam) []responses.ToolUnionParam {
	// Check if the request seems to be asking for image generation
	if o.shouldUseImageGeneration(opts) {
		outputFormat := getOutputFormatFromExtension(opts.ImageFile)
		imageGenTool := responses.ToolUnionParam{
			OfImageGeneration: &responses.ToolImageGenerationParam{
				Type:         ImageGenerationToolType,
				Model:        "gpt-image-1",
				OutputFormat: outputFormat,
				Quality:      "auto",
				Size:         "auto",
			},
		}
		tools = append(tools, imageGenTool)
	}
	return tools
}

// shouldUseImageGeneration determines if image generation should be enabled
// This is a heuristic based on the presence of --image-file flag
func (o *Client) shouldUseImageGeneration(opts *common.ChatOptions) bool {
	return opts.ImageFile != ""
}

// extractAndSaveImages extracts generated images from the response and saves them
func (o *Client) extractAndSaveImages(resp *responses.Response, opts *common.ChatOptions) error {
	if opts.ImageFile == "" {
		return nil // No image file specified, skip saving
	}

	// Extract image data from response
	for _, item := range resp.Output {
		if item.Type == ImageGenerationResponseType {
			imageCall := item.AsImageGenerationCall()
			if imageCall.Status == "completed" && imageCall.Result != "" {
				// Decode base64 image data
				imageData, err := base64.StdEncoding.DecodeString(imageCall.Result)
				if err != nil {
					return fmt.Errorf("failed to decode image data: %w", err)
				}

				// Ensure directory exists
				dir := filepath.Dir(opts.ImageFile)
				if dir != "." {
					if err := os.MkdirAll(dir, 0755); err != nil {
						return fmt.Errorf("failed to create directory %s: %w", dir, err)
					}
				}

				// Save image to file
				if err := os.WriteFile(opts.ImageFile, imageData, 0644); err != nil {
					return fmt.Errorf("failed to save image to %s: %w", opts.ImageFile, err)
				}

				fmt.Printf("Image saved to: %s\n", opts.ImageFile)
				return nil
			}
		}
	}

	return nil
}
