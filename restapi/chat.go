package restapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	goopenai "github.com/sashabaranov/go-openai"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	registry *core.PluginRegistry
	db       *fsdb.Db
}

type PromptRequest struct {
	UserInput   string `json:"userInput"`
	Vendor      string `json:"vendor"`
	Model       string `json:"model"`
	ContextName string `json:"contextName"`
	PatternName string `json:"patternName"`
}

type ChatRequest struct {
	Prompts            []PromptRequest `json:"prompts"`
	common.ChatOptions                 // Embed the ChatOptions from common package
}

type StreamResponse struct {
	Type    string `json:"type"`    // "content", "error", "complete"
	Format  string `json:"format"`  // "markdown", "mermaid", "plain"
	Content string `json:"content"` // The actual content
}

func NewChatHandler(r *gin.Engine, registry *core.PluginRegistry, db *fsdb.Db) *ChatHandler {
	handler := &ChatHandler{
		registry: registry,
		db:       db,
	}

	r.POST("/chat", handler.HandleChat)

	return handler
}

func (h *ChatHandler) HandleChat(c *gin.Context) {
	var request ChatRequest

	if err := c.BindJSON(&request); err != nil {
		log.Printf("Error binding JSON: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid request format: %v", err)})
		return
	}

	log.Printf("Received chat request with %d prompts", len(request.Prompts))

	// Set headers for SSE
	c.Writer.Header().Set("Content-Type", "text/readystream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
	c.Writer.Header().Set("X-Accel-Buffering", "no")

	clientGone := c.Writer.CloseNotify()

	for i, prompt := range request.Prompts {
		select {
		case <-clientGone:
			log.Printf("Client disconnected")
			return
		default:
			log.Printf("Processing prompt %d: Model=%s Pattern=%s Context=%s",
				i+1, prompt.Model, prompt.PatternName, prompt.ContextName)

			// Create chat channel for streaming
			streamChan := make(chan string)

			// Start chat processing in goroutine
			go func(p PromptRequest) {
				defer close(streamChan)

				chatter, err := h.registry.GetChatter(p.Model, 2048, false, false)
				if err != nil {
					log.Printf("Error creating chatter: %v", err)
					streamChan <- fmt.Sprintf("Error: %v", err)
					return
				}

				chatReq := &common.ChatRequest{
					Message: &goopenai.ChatCompletionMessage{
						Role:    "user",
						Content: p.UserInput,
					},
					PatternName: p.PatternName,
					ContextName: p.ContextName,
				}

				opts := &common.ChatOptions{
					Model:            p.Model,
					Temperature:      request.Temperature,
					TopP:             request.TopP,
					FrequencyPenalty: request.FrequencyPenalty,
					PresencePenalty:  request.PresencePenalty,
				}

				session, err := chatter.Send(chatReq, opts)
				if err != nil {
					log.Printf("Error from chatter.Send: %v", err)
					streamChan <- fmt.Sprintf("Error: %v", err)
					return
				}

				if session == nil {
					log.Printf("No session returned from chatter.Send")
					streamChan <- "Error: No response from model"
					return
				}

				// Get the last message from the session
				lastMsg := session.GetLastMessage()
				if lastMsg != nil {
					streamChan <- lastMsg.Content
				} else {
					log.Printf("No message content in session")
					streamChan <- "Error: No response content"
				}
			}(prompt)

			// Read from streamChan and write to client
			for content := range streamChan {
				select {
				case <-clientGone:
					return
				default:
					if strings.HasPrefix(content, "Error:") {
						response := StreamResponse{
							Type:    "error",
							Format:  "plain",
							Content: content,
						}
						if err := writeSSEResponse(c.Writer, response); err != nil {
							log.Printf("Error writing error response: %v", err)
							return
						}
					} else {
						response := StreamResponse{
							Type:    "content",
							Format:  detectFormat(content),
							Content: content,
						}
						if err := writeSSEResponse(c.Writer, response); err != nil {
							log.Printf("Error writing content response: %v", err)
							return
						}
					}
				}
			}

			// Signal completion of this prompt
			completeResponse := StreamResponse{
				Type:    "complete",
				Format:  "plain",
				Content: "",
			}
			if err := writeSSEResponse(c.Writer, completeResponse); err != nil {
				log.Printf("Error writing completion response: %v", err)
				return
			}
		}
	}
}

func writeSSEResponse(w gin.ResponseWriter, response StreamResponse) error {
	data, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("error marshaling response: %v", err)
	}

	if _, err := fmt.Fprintf(w, "data: %s\n\n", string(data)); err != nil {
		return fmt.Errorf("error writing response: %v", err)
	}

	w.(http.Flusher).Flush()
	return nil
}

func detectFormat(content string) string {
	if strings.HasPrefix(content, "graph TD") ||
		strings.HasPrefix(content, "gantt") ||
		strings.HasPrefix(content, "flowchart") ||
		strings.HasPrefix(content, "sequenceDiagram") ||
		strings.HasPrefix(content, "classDiagram") ||
		strings.HasPrefix(content, "stateDiagram") {
		return "mermaid"
	}
	if strings.Contains(content, "```") ||
		strings.Contains(content, "#") ||
		strings.Contains(content, "*") ||
		strings.Contains(content, "_") ||
		strings.Contains(content, "-") {
		return "markdown"
	}
	return "plain"
}
