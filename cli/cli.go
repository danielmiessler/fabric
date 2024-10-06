package cli

import (
	"fmt"
	"github.com/danielmiessler/fabric/converter"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/db"
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

	fabricDb := db.NewDb(filepath.Join(homedir, ".config/fabric"))

	// if the setup flag is set, run the setup function
	if currentFlags.Setup || currentFlags.SetupSkipPatterns || currentFlags.SetupVendor != "" {
		_ = fabricDb.Configure()
		if currentFlags.SetupVendor != "" {
			_, err = SetupVendor(fabricDb, currentFlags.SetupVendor)
		} else {
			_, err = Setup(fabricDb, currentFlags.SetupSkipPatterns)
		}
		return
	}

	var fabric *core.Fabric
	if err = fabricDb.Configure(); err != nil {
		fmt.Println("init is failed, run start the setup procedure", err)
		if fabric, err = Setup(fabricDb, currentFlags.SetupSkipPatterns); err != nil {
			return
		}
	} else {
		if fabric, err = core.NewFabric(fabricDb); err != nil {
			fmt.Println("fabric can't initialize, please run the --setup procedure", err)
			return
		}
	}

	if currentFlags.UpdatePatterns {
		err = fabric.PopulateDB()
		return
	}

	if currentFlags.ChangeDefaultModel {
		err = fabric.SetupDefaultModel()
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
		err = fabricDb.Patterns.ListNames()
		return
	}

	if currentFlags.ListAllModels {
		fabric.GetModels().Print()
		return
	}

	if currentFlags.ListAllContexts {
		err = fabricDb.Contexts.ListNames()
		return
	}

	if currentFlags.ListAllSessions {
		err = fabricDb.Sessions.ListNames()
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
			var language = "en"
			if currentFlags.Language != "" {
				language = currentFlags.Language
			}
			if transcript, err = fabric.YouTube.GrabTranscript(videoId, language); err != nil {
				return
			}

			currentFlags.AppendMessage(transcript)
		}

		if currentFlags.YouTubeComments {
			var comments []string
			if comments, err = fabric.YouTube.GrabComments(videoId); err != nil {
				return
			}

			commentsString := strings.Join(comments, "\n")

			currentFlags.AppendMessage(commentsString)
		}

		if !currentFlags.IsChatRequest() {
			// if the pattern flag is not set, we wanted only to grab the transcript or comments
			fmt.Println(currentFlags.Message)
			return
		}
	}

	if (currentFlags.ScrapeURL != "" || currentFlags.ScrapeQuestion != "") && fabric.Jina.IsConfigured() {
		// Check if the scrape_url flag is set and call ScrapeURL
		if currentFlags.ScrapeURL != "" {
			var website string
			if website, err = fabric.Jina.ScrapeURL(currentFlags.ScrapeURL); err != nil {
				return
			}

			currentFlags.AppendMessage(website)
		}

		// Check if the scrape_question flag is set and call ScrapeQuestion
		if currentFlags.ScrapeQuestion != "" {
			var website string
			if website, err = fabric.Jina.ScrapeQuestion(currentFlags.ScrapeQuestion); err != nil {
				return
			}

			currentFlags.AppendMessage(website)
		}

		if !currentFlags.IsChatRequest() {
			// if the pattern flag is not set, we wanted only to grab the url or get the answer to the question
			fmt.Println(currentFlags.Message)
			return
		}
	}

	var chatter *core.Chatter
	if chatter, err = fabric.GetChatter(currentFlags.Model, currentFlags.Stream, currentFlags.DryRun); err != nil {
		return
	}

	var session *db.Session
	if session, err = chatter.Send(currentFlags.BuildChatRequest(strings.Join(os.Args[1:], " ")), currentFlags.BuildChatOptions()); err != nil {
		return
	}

	result := session.GetLastMessage().Content

	if !currentFlags.Stream {
		// print the result if it was not streamed already
		fmt.Println(result)
	}

	// if the copy flag is set, copy the message to the clipboard
	if currentFlags.Copy {
		if err = fabric.CopyToClipboard(result); err != nil {
			return
		}
	}

	// if the output flag is set, create an output file
	if currentFlags.Output != "" {
		if currentFlags.OutputSession {
			sessionAsString := session.String()
			err = fabric.CreateOutputFile(sessionAsString, currentFlags.Output)
		} else {
			err = fabric.CreateOutputFile(result, currentFlags.Output)
		}
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

func SetupVendor(db *db.Db, vendorName string) (ret *core.Fabric, err error) {
	ret = core.NewFabricForSetup(db)
	if err = ret.SetupVendor(vendorName); err != nil {
		return
	}
	return
}
