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
func Cli() (message string, err error) {
	var currentFlags *Flags
	if currentFlags, err = Init(); err != nil {
		// we need to reset error, because we don't want to show double help messages
		err = nil
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

		if currentFlags.YouTubeTranscript {
			var transcript string
			if transcript, err = fabric.YouTube.GrabTranscript(videoId); err != nil {
				return
			}

			if currentFlags.Message != "" {
				currentFlags.Message = currentFlags.Message + "\n" + transcript
			} else {
				currentFlags.Message = transcript
			}
		}

		if currentFlags.YouTubeComments {
			var comments []string
			if comments, err = fabric.YouTube.GrabComments(videoId); err != nil {
				return
			}

			commentsString := strings.Join(comments, "\n")

			if currentFlags.Message != "" {
				currentFlags.Message = currentFlags.Message + "\n" + commentsString
			} else {
				currentFlags.Message = commentsString
			}
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
