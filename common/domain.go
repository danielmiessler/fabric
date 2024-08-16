package common

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	ContextName string
	SessionName string
	PatternName string
	Message     string
}

type ChatOptions struct {
	Model            string
	Temperature      float64
	TopP             float64
	PresencePenalty  float64
	FrequencyPenalty float64
}
