package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/danielmiessler/fabric/common"
	"github.com/jessevdk/go-flags"
	"golang.org/x/text/language"
)

// Flags create flags struct. the users flags go into this, this will be passed to the chat struct in cli
type Flags struct {
	Pattern            string            `short:"p" long:"pattern" description:"Choose a pattern from the available patterns" default:""`
	PatternVariables   map[string]string `short:"v" long:"variable" description:"Values for pattern variables, e.g. -v=#role:expert -v=#points:30"`
	Context            string            `short:"C" long:"context" description:"Choose a context from the available contexts" default:""`
	Session            string            `long:"session" description:"Choose a session from the available sessions"`
	Setup              bool              `short:"S" long:"setup" description:"Run setup for all reconfigurable parts of fabric"`
	Temperature        float64           `short:"t" long:"temperature" description:"Set temperature" default:"0.7"`
	TopP               float64           `short:"T" long:"topp" description:"Set top P" default:"0.9"`
	Stream             bool              `short:"s" long:"stream" description:"Stream"`
	PresencePenalty    float64           `short:"P" long:"presencepenalty" description:"Set presence penalty" default:"0.0"`
	Raw                bool              `short:"r" long:"raw" description:"Use the defaults of the model without sending chat options (like temperature etc.) and use the user role instead of the system role for patterns."`
	FrequencyPenalty   float64           `short:"F" long:"frequencypenalty" description:"Set frequency penalty" default:"0.0"`
	ListPatterns       bool              `short:"l" long:"listpatterns" description:"List all patterns"`
	ListAllModels      bool              `short:"L" long:"listmodels" description:"List all available models"`
	ListAllContexts    bool              `short:"x" long:"listcontexts" description:"List all contexts"`
	ListAllSessions    bool              `short:"X" long:"listsessions" description:"List all sessions"`
	UpdatePatterns     bool              `short:"U" long:"updatepatterns" description:"Update patterns"`
	Message            string            `hidden:"true" description:"Message to send to chat"`
	Copy               bool              `short:"c" long:"copy" description:"Copy to clipboard"`
	Model              string            `short:"m" long:"model" description:"Choose model"`
	Output             string            `short:"o" long:"output" description:"Output to file" default:""`
	OutputSession      bool              `long:"output-session" description:"Output the entire session (also a temporary one) to the output file"`
	LatestPatterns     string            `short:"n" long:"latest" description:"Number of latest patterns to list" default:"0"`
	ChangeDefaultModel bool              `short:"d" long:"changeDefaultModel" description:"Change default model"`
	YouTube            string            `short:"y" long:"youtube" description:"YouTube video \"URL\" to grab transcript, comments from it and send to chat"`
	YouTubeTranscript  bool              `long:"transcript" description:"Grab transcript from YouTube video and send to chat (it used per default)."`
	YouTubeComments    bool              `long:"comments" description:"Grab comments from YouTube video and send to chat"`
	Language           string            `short:"g" long:"language" description:"Specify the Language Code for the chat, e.g. -g=en -g=zh" default:""`
	ScrapeURL          string            `short:"u" long:"scrape_url" description:"Scrape website URL to markdown using Jina AI"`
	ScrapeQuestion     string            `short:"q" long:"scrape_question" description:"Search question using Jina AI"`
	Seed               int               `short:"e" long:"seed" description:"Seed to be used for LMM generation"`
	WipeContext        string            `short:"w" long:"wipecontext" description:"Wipe context"`
	WipeSession        string            `short:"W" long:"wipesession" description:"Wipe session"`
	PrintContext       string            `long:"printcontext" description:"Print context"`
	PrintSession       string            `long:"printsession" description:"Print session"`
	HtmlReadability    bool              `long:"readability" description:"Convert HTML input into a clean, readable view"`
	DryRun             bool              `long:"dry-run" description:"Show what would be sent to the model without actually sending it"`
	Serve              bool              `long:"serve" description:"Serve the Fabric Rest API"`
	ServeAddress       string            `long:"address" description:"The address to bind the REST API" default:":8080"`
	Version            bool              `long:"version" description:"Print current version"`
}

// Init Initialize flags. returns a Flags struct and an error
func Init() (ret *Flags, err error) {
	var message string

	ret = &Flags{}
	parser := flags.NewParser(ret, flags.Default)
	var args []string
	if args, err = parser.Parse(); err != nil {
		return
	}

	info, _ := os.Stdin.Stat()
	hasStdin := (info.Mode() & os.ModeCharDevice) == 0

	// takes input from stdin if it exists, otherwise takes input from args (the last argument)
	if hasStdin {
		if message, err = readStdin(); err != nil {
			return
		}
	} else if len(args) > 0 {
		message = args[len(args)-1]
	} else {
		message = ""
	}
	ret.Message = message

	return
}

// readStdin reads from stdin and returns the input as a string or an error
func readStdin() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	var input string
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return "", fmt.Errorf("error reading from stdin: %w", err)
		}
		input += line
	}
	return input, nil
}

func (o *Flags) BuildChatOptions() (ret *common.ChatOptions) {
	ret = &common.ChatOptions{
		Temperature:      o.Temperature,
		TopP:             o.TopP,
		PresencePenalty:  o.PresencePenalty,
		FrequencyPenalty: o.FrequencyPenalty,
		Raw:              o.Raw,
		Seed:             o.Seed,
	}
	return
}

func (o *Flags) BuildChatRequest(Meta string) (ret *common.ChatRequest) {
	ret = &common.ChatRequest{
		ContextName:      o.Context,
		SessionName:      o.Session,
		PatternName:      o.Pattern,
		PatternVariables: o.PatternVariables,
		Message:          o.Message,
		Meta:             Meta,
	}
	if o.Language != "" {
		langTag, err := language.Parse(o.Language)
		if err == nil {
			ret.Language = langTag.String()
		}
	}
	return
}

func (o *Flags) AppendMessage(message string) {
	if o.Message != "" {
		o.Message = o.Message + "\n" + message
	} else {
		o.Message = message
	}
	return
}

func (o *Flags) IsChatRequest() (ret bool) {
	ret = (o.Message != "" || o.Context != "") && (o.Session != "" || o.Pattern != "")
	return
}
