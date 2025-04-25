package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielmiessler/fabric/plugins/tools/youtube"

	"github.com/danielmiessler/fabric/common"
	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/plugins/ai"
	"github.com/danielmiessler/fabric/plugins/db/fsdb"
	"github.com/danielmiessler/fabric/plugins/tools/converter"
	"github.com/danielmiessler/fabric/restapi"
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

	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return
	}

	fabricDb := fsdb.NewDb(filepath.Join(homedir, ".config/fabric"))

	if err = fabricDb.Configure(); err != nil {
		if !currentFlags.Setup {
			println(err.Error())
			currentFlags.Setup = true
		}
	}

	var registry *core.PluginRegistry
	if registry, err = core.NewPluginRegistry(fabricDb); err != nil {
		return
	}

	// if the setup flag is set, run the setup function
	if currentFlags.Setup {
		err = registry.Setup()
		return
	}

	if currentFlags.Serve {
		registry.ConfigureVendors()
		err = restapi.Serve(registry, currentFlags.ServeAddress, currentFlags.ServeAPIKey)
		return
	}

	if currentFlags.ServeOllama {
		registry.ConfigureVendors()
		err = restapi.ServeOllama(registry, currentFlags.ServeAddress, version)
		return
	}

	if currentFlags.UpdatePatterns {
		err = registry.PatternsLoader.PopulateDB()
		return
	}

	if currentFlags.ChangeDefaultModel {
		if err = registry.Defaults.Setup(); err != nil {
			return
		}
		err = registry.SaveEnvFile()
		return
	}

	if currentFlags.LatestPatterns != "0" {
		var parsedToInt int
		if parsedToInt, err = strconv.Atoi(currentFlags.LatestPatterns); err != nil {
			return
		}

		if err = fabricDb.Patterns.PrintLatestPatterns(parsedToInt); err != nil {
			return
		}
		return
	}

	if currentFlags.ListPatterns {
		err = fabricDb.Patterns.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllModels {
		var models *ai.VendorsModels
		if models, err = registry.VendorManager.GetModels(); err != nil {
			return
		}
		models.Print(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllContexts {
		err = fabricDb.Contexts.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListAllSessions {
		err = fabricDb.Sessions.ListNames(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.WipeContext != "" {
		err = fabricDb.Contexts.Delete(currentFlags.WipeContext)
		return
	}

	if currentFlags.WipeSession != "" {
		err = fabricDb.Sessions.Delete(currentFlags.WipeSession)
		return
	}

	if currentFlags.PrintSession != "" {
		err = fabricDb.Sessions.PrintSession(currentFlags.PrintSession)
		return
	}

	if currentFlags.PrintContext != "" {
		err = fabricDb.Contexts.PrintContext(currentFlags.PrintContext)
		return
	}

	if currentFlags.HtmlReadability {
		if msg, cleanErr := converter.HtmlReadability(currentFlags.Message); cleanErr != nil {
			fmt.Println("use original input, because can't apply html readability", err)
		} else {
			currentFlags.Message = msg
		}
	}

	if currentFlags.ListExtensions {
		err = registry.TemplateExtensions.ListExtensions()
		return
	}

	if currentFlags.AddExtension != "" {
		err = registry.TemplateExtensions.RegisterExtension(currentFlags.AddExtension)
		return
	}

	if currentFlags.RemoveExtension != "" {
		err = registry.TemplateExtensions.RemoveExtension(currentFlags.RemoveExtension)
		return
	}

	if currentFlags.ListStrategies {
		err = registry.Strategies.ListStrategies(currentFlags.ShellCompleteOutput)
		return
	}

	if currentFlags.ListVendors {
		err = registry.ListVendors(os.Stdout)
		return
	}

	// if the interactive flag is set, run the interactive function
	// if currentFlags.Interactive {
	// 	interactive.Interactive()
	// }

	// if none of the above currentFlags are set, run the initiate chat function

	var messageTools string

	if currentFlags.YouTube != "" {
		if !registry.YouTube.IsConfigured() {
			err = fmt.Errorf("YouTube is not configured, please run the setup procedure")
			return
		}

		var videoId string
		var playlistId string
		if videoId, playlistId, err = registry.YouTube.GetVideoOrPlaylistId(currentFlags.YouTube); err != nil {
			return
		} else if (videoId == "" || currentFlags.YouTubePlaylist) && playlistId != "" {
			if currentFlags.Output != "" {
				err = registry.YouTube.FetchAndSavePlaylist(playlistId, currentFlags.Output)
			} else {
				var videos []*youtube.VideoMeta
				if videos, err = registry.YouTube.FetchPlaylistVideos(playlistId); err != nil {
					err = fmt.Errorf("error fetching playlist videos: %v", err)
					return
				}

				for _, video := range videos {
					var message string
					if message, err = processYoutubeVideo(currentFlags, registry, video.Id); err != nil {
						return
					}

					if !currentFlags.IsChatRequest() {
						if err = WriteOutput(message, fmt.Sprintf("%v.md", video.TitleNormalized)); err != nil {
							return
						}
					} else {
						messageTools = AppendMessage(messageTools, message)
					}
				}
			}
			return
		}

		if messageTools, err = processYoutubeVideo(currentFlags, registry, videoId); err != nil {
			return
		}
		if !currentFlags.IsChatRequest() {
			err = currentFlags.WriteOutput(messageTools)
			return
		}
	}

	if (currentFlags.ScrapeURL != "" || currentFlags.ScrapeQuestion != "") && registry.Jina.IsConfigured() {
		// Check if the scrape_url flag is set and call ScrapeURL
		if currentFlags.ScrapeURL != "" {
			var website string
			if website, err = registry.Jina.ScrapeURL(currentFlags.ScrapeURL); err != nil {
				return
			}
			messageTools = AppendMessage(messageTools, website)
		}

		// Check if the scrape_question flag is set and call ScrapeQuestion
		if currentFlags.ScrapeQuestion != "" {
			var website string
			if website, err = registry.Jina.ScrapeQuestion(currentFlags.ScrapeQuestion); err != nil {
				return
			}

			messageTools = AppendMessage(messageTools, website)
		}

		if !currentFlags.IsChatRequest() {
			err = currentFlags.WriteOutput(messageTools)
			return
		}
	}

	if messageTools != "" {
		currentFlags.AppendMessage(messageTools)
	}

	var chatter *core.Chatter
	if chatter, err = registry.GetChatter(currentFlags.Model, currentFlags.ModelContextLength, currentFlags.Strategy, currentFlags.Stream, currentFlags.DryRun); err != nil {
		return
	}

	var session *fsdb.Session
	var chatReq *common.ChatRequest
	if chatReq, err = currentFlags.BuildChatRequest(strings.Join(os.Args[1:], " ")); err != nil {
		return
	}

	if chatReq.Language == "" {
		chatReq.Language = registry.Language.DefaultLanguage.Value
	}
	if session, err = chatter.Send(chatReq, currentFlags.BuildChatOptions()); err != nil {
		return
	}

	result := session.GetLastMessage().Content

	if !currentFlags.Stream {
		// print the result if it was not streamed already
		fmt.Println(result)
	}

	// if the copy flag is set, copy the message to the clipboard
	if currentFlags.Copy {
		if err = CopyToClipboard(result); err != nil {
			return
		}
	}

	// if the output flag is set, create an output file
	if currentFlags.Output != "" {
		if currentFlags.OutputSession {
			sessionAsString := session.String()
			err = CreateOutputFile(sessionAsString, currentFlags.Output)
		} else {
			err = CreateOutputFile(result, currentFlags.Output)
		}
	}
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
