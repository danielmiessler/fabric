package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/db"
)

// Cli Controls the cli. It takes in the flags and runs the appropriate functions
func Cli(version string) (message string, err error) {
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

	fabricDb := db.NewDb(filepath.Join(homedir, ".config/fabric"))

	// if the setup flag is set, run the setup function
	if currentFlags.Setup {
		_ = fabricDb.Configure()
		_, err = Setup(fabricDb, currentFlags.SetupSkipUpdatePatterns)
		return
	}

	var fabric *core.Fabric
	if err = fabricDb.Configure(); err != nil {
		fmt.Println("init is failed, run start the setup procedure", err)
		if fabric, err = Setup(fabricDb, currentFlags.SetupSkipUpdatePatterns); err != nil {
			return
		}
	} else {
		if fabric, err = core.NewFabric(fabricDb); err != nil {
			fmt.Println("fabric can't initialize, please run the --setup procedure", err)
			return
		}
	}

	// if the update patterns flag is set, run the update patterns function
	if currentFlags.UpdatePatterns {
		err = fabric.PopulateDB()
		return
	}

	if currentFlags.ChangeDefaultModel {
		err = fabric.SetupDefaultModel()
		return
	}

	// if the latest patterns flag is set, run the latest patterns function
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

	// if the list patterns flag is set, run the list all patterns function
	if currentFlags.ListPatterns {
		err = fabricDb.Patterns.ListNames()
		return
	}

	// if the list all models flag is set, run the list all models function
	if currentFlags.ListAllModels {
		fabric.GetModels().Print()
		return
	}

	// if the list all contexts flag is set, run the list all contexts function
	if currentFlags.ListAllContexts {
		err = fabricDb.Contexts.ListNames()
		return
	}

	// if the list all sessions flag is set, run the list all sessions function
	if currentFlags.ListAllSessions {
		err = fabricDb.Sessions.ListNames()
		return
	}

	// if the wipe context flag is set, run the wipe context function
	if currentFlags.WipeContext != "" {
		err = fabricDb.Contexts.Delete(currentFlags.WipeContext)
		return
	}

	// if the wipe session flag is set, run the wipe session function
	if currentFlags.WipeSession != "" {
		err = fabricDb.Sessions.Delete(currentFlags.WipeSession)
		return
	}

	// if the interactive flag is set, run the interactive function
	// if currentFlags.Interactive {
	// 	interactive.Interactive()
	// }

	// if none of the above currentFlags are set, run the initiate chat function

	if currentFlags.YouTube != "" {
		if fabric.YouTube.IsConfigured() == false {
			err = fmt.Errorf("YouTube is not configured, please run the setup procedure")
			return
		}

		var videoId string
		if videoId, err = fabric.YouTube.GetVideoId(currentFlags.YouTube); err != nil {
			return
		}

		if !currentFlags.YouTubeComments || currentFlags.YouTubeTranscript {
			var transcript string
			if transcript, err = fabric.YouTube.GrabTranscript(videoId); err != nil {
				return
			}

			// fmt.Println(transcript)

			currentFlags.AppendMessage(transcript)
		}

		if currentFlags.YouTubeComments {
			var comments []string
			if comments, err = fabric.YouTube.GrabComments(videoId); err != nil {
				return
			}

			commentsString := strings.Join(comments, "\n")

			// fmt.Println(commentsString)

			currentFlags.AppendMessage(commentsString)
		}

		if currentFlags.Pattern == "" {
			// if the pattern flag is not set, we wanted only to grab the transcript or comments
			fmt.Println(currentFlags.Message)
			return
		}
	}

	if (currentFlags.ScrapeURL != "" || currentFlags.ScrapeQuestion != "") && fabric.Jina.IsConfigured() {
		// Check if the scrape_url flag is set and call ScrapeURL
		if currentFlags.ScrapeURL != "" {
			if message, err = fabric.Jina.ScrapeURL(currentFlags.ScrapeURL); err != nil {
				return
			}

			//fmt.Println(message)

			currentFlags.AppendMessage(message)
		}

		// Check if the scrape_question flag is set and call ScrapeQuestion
		if currentFlags.ScrapeQuestion != "" {
			if message, err = fabric.Jina.ScrapeQuestion(currentFlags.ScrapeQuestion); err != nil {
				return
			}

			//fmt.Println(message)

			currentFlags.AppendMessage(message)
		}

		if currentFlags.Pattern == "" {
			// if the pattern flag is not set, we wanted only to grab the url or get the answer to the question
			fmt.Println(currentFlags.Message)
			return
		}
	}

	if currentFlags.HtmlReadability {
		if msg, err := core.HtmlReadability(currentFlags.Message); err != nil {
			fmt.Println("use readability parser msg err:", err)
		} else {
			currentFlags.Message = msg
		}
	}

	var chatter *core.Chatter
	if chatter, err = fabric.GetChatter(currentFlags.Model, currentFlags.Stream, currentFlags.DryRun); err != nil {
		return
	}

	if message, err = chatter.Send(currentFlags.BuildChatRequest(), currentFlags.BuildChatOptions()); err != nil {
		return
	}

	if !currentFlags.Stream {
		fmt.Println(message)
	}

	// if the copy flag is set, copy the message to the clipboard
	if currentFlags.Copy {
		if err = fabric.CopyToClipboard(message); err != nil {
			return
		}
	}

	// if the output flag is set, create an output file
	if currentFlags.Output != "" {
		err = fabric.CreateOutputFile(message, currentFlags.Output)
	}
	return
}

func Setup(db *db.Db, skipUpdatePatterns bool) (ret *core.Fabric, err error) {
	instance := core.NewFabricForSetup(db)

	if err = instance.Setup(); err != nil {
		return
	}

	if !skipUpdatePatterns {
		if err = instance.PopulateDB(); err != nil {
			return
		}
	}
	ret = instance
	return
}
