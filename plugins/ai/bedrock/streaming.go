package bedrock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

// StreamingResponse structures for different providers

// AnthropicStreamingResponse represents a streaming response chunk from Anthropic
type AnthropicStreamingResponse struct {
	Type         string `json:"type"`
	Index        int    `json:"index,omitempty"`
	Delta        *Delta `json:"delta,omitempty"`
	StopReason   string `json:"stop_reason,omitempty"`
	StopSequence string `json:"stop_sequence,omitempty"`
}

// Delta represents the incremental content in a streaming response
type Delta struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// MetaStreamingResponse represents a streaming response from Meta models
type MetaStreamingResponse struct {
	Generation           string `json:"generation"`
	PromptTokenCount     int    `json:"prompt_token_count,omitempty"`
	GenerationTokenCount int    `json:"generation_token_count,omitempty"`
	StopReason           string `json:"stop_reason,omitempty"`
}

// MistralStreamingResponse represents a streaming response from Mistral models
type MistralStreamingResponse struct {
	Outputs []struct {
		Text       string `json:"text"`
		StopReason string `json:"stop_reason,omitempty"`
	} `json:"outputs"`
}

// CohereStreamingResponse represents a streaming response from Cohere models
type CohereStreamingResponse struct {
	IsFinished   bool     `json:"is_finished"`
	EventType    string   `json:"event_type"`
	Text         string   `json:"text,omitempty"`
	FinishReason string   `json:"finish_reason,omitempty"`
}

// AI21StreamingResponse represents a response from AI21 models (non-streaming)
type AI21StreamingResponse struct {
	Completions []struct {
		Data struct {
			Text string `json:"text"`
		} `json:"data"`
		FinishReason struct {
			Reason string `json:"reason"`
		} `json:"finishReason"`
	} `json:"completions"`
}

// HandleStreamingResponse processes streaming responses from Bedrock
func HandleStreamingResponse(ctx context.Context, output *bedrockruntime.InvokeModelWithResponseStreamOutput, modelID string, channel chan string) error {
	defer close(channel)

	provider := GetModelProvider(modelID)
	reader := output.GetStream().Reader

	eventReader := reader.Events()
	for {
		event := <-eventReader
		if event == nil {
			break
		}

		switch v := event.(type) {
		case *types.ResponseStreamMemberChunk:
			err := processChunk(v.Value.Bytes, provider, channel)
			if err != nil {
				return fmt.Errorf("error processing chunk: %w", err)
			}

		default:
			// Handle other event types if needed
		}
	}

	if err := reader.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	return nil
}

// processChunk processes a single chunk based on the provider
func processChunk(data []byte, provider ModelProvider, channel chan string) error {
	switch provider {
	case ProviderAnthropic:
		return processAnthropicChunk(data, channel)
	case ProviderMeta:
		return processMetaChunk(data, channel)
	case ProviderMistral:
		return processMistralChunk(data, channel)
	case ProviderCohere:
		return processCohereChunk(data, channel)
	case ProviderAI21:
		return processAI21Chunk(data, channel)
	default:
		return fmt.Errorf("unsupported provider for streaming: %s", provider)
	}
}

// processAnthropicChunk processes an Anthropic streaming chunk
func processAnthropicChunk(data []byte, channel chan string) error {
	var response AnthropicStreamingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	if response.Type == "content_block_delta" && response.Delta != nil {
		channel <- response.Delta.Text
	}

	return nil
}

// processMetaChunk processes a Meta streaming chunk
func processMetaChunk(data []byte, channel chan string) error {
	var response MetaStreamingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	if response.Generation != "" {
		channel <- response.Generation
	}

	return nil
}

// processMistralChunk processes a Mistral streaming chunk
func processMistralChunk(data []byte, channel chan string) error {
	var response MistralStreamingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	for _, output := range response.Outputs {
		if output.Text != "" {
			channel <- output.Text
		}
	}

	return nil
}

// processCohereChunk processes a Cohere streaming chunk
func processCohereChunk(data []byte, channel chan string) error {
	// Cohere sends chunks as newline-delimited JSON
	decoder := json.NewDecoder(bytes.NewReader(data))
	for {
		var response CohereStreamingResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		if response.EventType == "text-generation" && response.Text != "" {
			channel <- response.Text
		}
	}

	return nil
}

// processAI21Chunk processes an AI21 response (non-streaming)
func processAI21Chunk(data []byte, channel chan string) error {
	var response AI21StreamingResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return err
	}

	for _, completion := range response.Completions {
		if completion.Data.Text != "" {
			channel <- completion.Data.Text
		}
	}

	return nil
}

// ParseNonStreamingResponse parses non-streaming responses from different providers
func ParseNonStreamingResponse(data []byte, provider ModelProvider) (string, error) {
	switch provider {
	case ProviderAnthropic:
		var response struct {
			Content []struct {
				Text string `json:"text"`
			} `json:"content"`
		}
		if err := json.Unmarshal(data, &response); err != nil {
			return "", err
		}
		if len(response.Content) > 0 {
			return response.Content[0].Text, nil
		}
		return "", nil

	case ProviderMeta:
		var response MetaStreamingResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return "", err
		}
		return response.Generation, nil

	case ProviderMistral:
		var response MistralStreamingResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return "", err
		}
		if len(response.Outputs) > 0 {
			return response.Outputs[0].Text, nil
		}
		return "", nil

	case ProviderCohere:
		var response struct {
			Generations []struct {
				Text string `json:"text"`
			} `json:"generations"`
		}
		if err := json.Unmarshal(data, &response); err != nil {
			return "", err
		}
		if len(response.Generations) > 0 {
			return response.Generations[0].Text, nil
		}
		return "", nil

	case ProviderAI21:
		var response AI21StreamingResponse
		if err := json.Unmarshal(data, &response); err != nil {
			return "", err
		}
		if len(response.Completions) > 0 {
			return response.Completions[0].Data.Text, nil
		}
		return "", nil

	default:
		return "", fmt.Errorf("unsupported provider: %s", provider)
	}
}