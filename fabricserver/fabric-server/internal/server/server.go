package server

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"regexp"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/danielmiessler/fabric/plugins/tools/youtube"
)

// YouTubeProcessor handles YouTube content processing
type YouTubeProcessor struct {
	youtube *youtube.YouTube
}

func NewYouTubeProcessor() (*YouTubeProcessor, error) {
	yt := youtube.NewYouTube()
	
	// Get API key from environment variable
	apiKey := os.Getenv("YOUTUBE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("YOUTUBE_API_KEY environment variable not set")
	}
	
	yt.ApiKey.Value = apiKey
	
	return &YouTubeProcessor{
		youtube: yt,
	}, nil
}

// isYouTubeURL checks if the input is a YouTube URL
func (p *YouTubeProcessor) isYouTubeURL(input string) bool {
	youtubePattern := regexp.MustCompile(`(?i)(?:https?://)?(?:www\.)?(?:youtube\.com|youtu\.be)`)
	return youtubePattern.MatchString(input)
}

// processYouTubeContent processes a YouTube URL to get transcript and other content
func (p *YouTubeProcessor) processYouTubeContent(url string) (string, error) {
	// Extract video ID
	videoId, playlistId, err := p.youtube.GetVideoOrPlaylistId(url)
	if err != nil {
		return "", fmt.Errorf("failed to extract video ID: %v", err)
	}
	
	if videoId == "" && playlistId != "" {
		return "", fmt.Errorf("playlist URLs are not supported, please provide a video URL")
	}

	// Get video transcript
	transcript, err := p.youtube.GrabTranscript(videoId, "en")
	if err != nil {
		return "", fmt.Errorf("failed to get transcript: %v", err)
	}

	return transcript, nil
}

type Server struct {
	router *mux.Router
	youtube *YouTubeProcessor
}

func NewServer() (*Server, error) {
	yt, err := NewYouTubeProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize YouTube processor: %v", err)
	}

	s := &Server{
		router: mux.NewRouter(),
		youtube: yt,
	}
	s.routes()
	return s, nil
}

func (s *Server) routes() {
	// API routes
	s.router.HandleFunc("/api/fabric/options", s.getOptions).Methods("GET")
	s.router.HandleFunc("/chat", s.postChat).Methods("POST")
	
	// Serve static files from the root directory
	s.router.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))
}

func (s *Server) Start() error {
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},                                   // Allow all origins for development
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: false,                                           // Keep this false to allow wildcard origin
		Debug:            true,
	})

	handler := c.Handler(s.router)

	fmt.Println("Fabric server listening at http://localhost:3001")
	return http.ListenAndServe(":3001", handler)
}

func (s *Server) getOptions(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Fetching options from fabric server...")
	
	resp, err := http.Get("http://localhost:8080/patterns/names")
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to fetch patterns: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	w.Header().Set("Content-Type", "application/json")
	io.Copy(w, resp.Body)
}

type Prompt struct {
	UserInput   string `json:"userInput"`
	PatternName string `json:"patternName"`
}

type ChatRequest struct {
	Prompts []Prompt `json:"prompts"`
}

func (s *Server) postChat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if len(req.Prompts) == 0 {
		http.Error(w, "No prompts provided", http.StatusBadRequest)
		return
	}

	// Use first prompt for now since the fabric server only handles one at a time
	prompt := req.Prompts[0]
	
	// Check if input is a YouTube URL and process it first if needed
	userInput := prompt.UserInput
	if s.youtube.isYouTubeURL(userInput) {
		processedContent, err := s.youtube.processYouTubeContent(userInput)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to process YouTube content: %v", err), http.StatusInternalServerError)
			return
		}
		userInput = processedContent
	}
	
	// Forward the request to the fabric server with the required format
	fabricReq := struct {
		Prompts []struct {
			UserInput    string `json:"userInput"`
			Vendor      string `json:"vendor"`
			Model       string `json:"model"`
			ContextName string `json:"contextName,omitempty"`
			PatternName string `json:"patternName"`
		} `json:"prompts"`
		Temperature      *float64 `json:"temperature,omitempty"`
		TopP            *float64 `json:"topP,omitempty"`
		FrequencyPenalty *float64 `json:"frequencyPenalty,omitempty"`
		PresencePenalty  *float64 `json:"presencePenalty,omitempty"`
	}{
		Prompts: []struct {
			UserInput    string `json:"userInput"`
			Vendor      string `json:"vendor"`
			Model       string `json:"model"`
			ContextName string `json:"contextName,omitempty"`
			PatternName string `json:"patternName"`
		}{{
			UserInput:   userInput,
			PatternName: prompt.PatternName,
			Vendor:     "Gemini", // Default vendor
			Model:      "gemini-2.0-flash-exp",  // Default model
		}},
	}

	jsonData, err := json.Marshal(fabricReq)
	if err != nil {
		http.Error(w, "Failed to marshal request", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("http://localhost:8080/chat", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to forward request: %v", err), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Set SSE headers first
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // Disable buffering
	
	// Ensure headers are sent before streaming
	w.WriteHeader(http.StatusOK)
	
	// Create a buffer to read the response line by line
	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Printf("Received from fabric server: %s\n", line)
		// Skip empty lines
		if line == "" {
			continue
		}

		// Extract JSON content from SSE data line
		jsonStr := line
		if strings.HasPrefix(line, "data: ") {
			jsonStr = line[6:]
		}

		// Parse and re-serialize the JSON to ensure proper formatting
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
			fmt.Printf("Error parsing JSON from fabric server: %v\nJSON: %s\n", err, jsonStr)
			continue
		}

		// Re-serialize with proper formatting
		formattedJSON, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("Error re-serializing JSON: %v\n", err)
			continue
		}

		// Send as SSE event
		_, err = fmt.Fprintf(w, "data: %s\n\n", formattedJSON)
		if err != nil {
			fmt.Printf("Error writing to response: %v\n", err)
			return
		}
		// Flush the response writer after each line
		if f, ok := w.(http.Flusher); ok {
			f.Flush()
		}
	}
	
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
	}
}
