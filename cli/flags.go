package cli

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/danielmiessler/fabric/common"
	"github.com/jessevdk/go-flags"
)

// Flags create flags struct. the users flags go into this, this will be passed to the chat struct in cli
type Flags struct {
	Pattern                 string            `short:"p" long:"pattern" description:"Choose a pattern" default:""`
	PatternVariables        map[string]string `short:"v" long:"variable" description:"Values for pattern variables, e.g. -v=$name:John -v=$age:30"`
	Context                 string            `short:"C" long:"context" description:"Choose a context" default:""`
	Session                 string            `long:"session" description:"Choose a session"`
	Setup                   bool              `short:"S" long:"setup" description:"Run setup"`
	SetupSkipUpdatePatterns bool              `long:"setup-skip-update-patterns" description:"Skip update patterns at setup"`
	Temperature             float64           `short:"t" long:"temperature" description:"Set temperature" default:"0.7"`
	TopP                    float64           `short:"T" long:"topp" description:"Set top P" default:"0.9"`
	Stream                  bool              `short:"s" long:"stream" description:"Stream"`
	PresencePenalty         float64           `short:"P" long:"presencepenalty" description:"Set presence penalty" default:"0.0"`
	FrequencyPenalty        float64           `short:"F" long:"frequencypenalty" description:"Set frequency penalty" default:"0.0"`
	ListPatterns            bool              `short:"l" long:"listpatterns" description:"List all patterns"`
	ListAllModels           bool              `short:"L" long:"listmodels" description:"List all available models"`
	ListAllContexts         bool              `short:"x" long:"listcontexts" description:"List all contexts"`
	ListAllSessions         bool              `short:"X" long:"listsessions" description:"List all sessions"`
	UpdatePatterns          bool              `short:"U" long:"updatepatterns" description:"Update patterns"`
	Message                 string            `hidden:"true" description:"Message to send to chat"`
	Copy                    bool              `short:"c" long:"copy" description:"Copy to clipboard"`
	Model                   string            `short:"m" long:"model" description:"Choose model"`
	Output                  string            `short:"o" long:"output" description:"Output to file" default:""`
	LatestPatterns          string            `short:"n" long:"latest" description:"Number of latest patterns to list" default:"0"`
	ChangeDefaultModel      bool              `short:"d" long:"changeDefaultModel" description:"Change default pattern"`
	YouTube                 string            `short:"y" long:"youtube" description:"YouTube video url to grab transcript, comments from it and send to chat"`
	YouTubeTranscript       bool              `long:"transcript" description:"Grab transcript from YouTube video and send to chat"`
	YouTubeComments         bool              `long:"comments" description:"Grab comments from YouTube video and send to chat"`
	DryRun                  bool              `long:"dry-run" description:"Show what would be sent to the model without actually sending it"`
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
	}
	return
}

func (o *Flags) BuildChatRequest() (ret *common.ChatRequest) {
	ret = &common.ChatRequest{
		ContextName:      o.Context,
		SessionName:      o.Session,
		PatternName:      o.Pattern,
		PatternVariables: o.PatternVariables,
		Message:          o.Message,
	}
	return
}
