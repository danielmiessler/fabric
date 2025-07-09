package cli

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/danielmiessler/fabric/internal/core"
	"github.com/danielmiessler/fabric/internal/tools/converter"
	"github.com/danielmiessler/fabric/internal/tools/youtube"
)

// Cli Controls the cli. It takes in the flags and runs the appropriate functions
func Cli(version string) (err error) {
	var currentFlags *Flags
	if currentFlags, err = Init(); err != nil {
		return
	}

	if currentFlags.Version {
		fmt.Println(version)
		return
	}

	// Initialize database and registry
	var registry, err2 = initializeFabric()
	if err2 != nil {
		if !currentFlags.Setup {
			println(err2.Error())
			currentFlags.Setup = true
		}
	}

	// Handle setup and server commands
	var handled bool
	if handled, err = handleSetupAndServerCommands(currentFlags, registry, version); err != nil || handled {
		return
	}

	// Handle configuration commands
	if handled, err = handleConfigurationCommands(currentFlags, registry); err != nil || handled {
		return
	}

	// Handle listing commands
	if handled, err = handleListingCommands(currentFlags, registry.Db, registry); err != nil || handled {
		return
	}

	// Handle management commands
	if handled, err = handleManagementCommands(currentFlags, registry.Db); err != nil || handled {
		return
	}

	// Handle extension commands
	if handled, err = handleExtensionCommands(currentFlags, registry); err != nil || handled {
		return
	}

	// Process HTML readability if needed
	if currentFlags.HtmlReadability {
		if msg, cleanErr := converter.HtmlReadability(currentFlags.Message); cleanErr != nil {
			fmt.Println("use original input, because can't apply html readability", err)
		} else {
			currentFlags.Message = msg
		}
	}

	// Handle tool-based message processing
	var messageTools string
	if messageTools, err = handleToolProcessing(currentFlags, registry); err != nil {
		return
	}

	// Handle chat processing
	err = handleChatProcessing(currentFlags, registry, messageTools)
	return
}

func processYoutubeVideo(
	flags *Flags, registry *core.PluginRegistry, videoId string) (message string, err error) {

	if (!flags.YouTubeComments && !flags.YouTubeMetadata) || flags.YouTubeTranscript || flags.YouTubeTranscriptWithTimestamps {
		var transcript string
		var language = "en"
		if flags.Language != "" || registry.Language.DefaultLanguage.Value != "" {
			if flags.Language != "" {
				language = flags.Language
			} else {
				language = registry.Language.DefaultLanguage.Value
			}
		}
		if flags.YouTubeTranscriptWithTimestamps {
			if transcript, err = registry.YouTube.GrabTranscriptWithTimestamps(videoId, language); err != nil {
				return
			}
		} else {
			if transcript, err = registry.YouTube.GrabTranscript(videoId, language); err != nil {
				return
			}
		}
		message = AppendMessage(message, transcript)
	}

	if flags.YouTubeComments {
		var comments []string
		if comments, err = registry.YouTube.GrabComments(videoId); err != nil {
			return
		}

		commentsString := strings.Join(comments, "\n")

		message = AppendMessage(message, commentsString)
	}

	if flags.YouTubeMetadata {
		var metadata *youtube.VideoMetadata
		if metadata, err = registry.YouTube.GrabMetadata(videoId); err != nil {
			return
		}
		metadataJson, _ := json.MarshalIndent(metadata, "", "  ")
		message = AppendMessage(message, string(metadataJson))
	}

	return
}

func WriteOutput(message string, outputFile string) (err error) {
	fmt.Println(message)
	if outputFile != "" {
		err = CreateOutputFile(message, outputFile)
	}
	return
}
