package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/danielmiessler/fabric/core"
	"github.com/danielmiessler/fabric/db"
)

// Cli Controls the cli. It takes in the flags and runs the appropriate functions
func Cli() (message string, err error) {
	var currentFlags *Flags
	if currentFlags, err = Init(); err != nil {
		// we need to reset error, because we want to show double help messages
		err = nil
		return
	}

	var homedir string
	if homedir, err = os.UserHomeDir(); err != nil {
		return
	}

	db := db.NewDb(filepath.Join(homedir, ".config/fabric"))

	// if the setup flag is set, run the setup function
	if currentFlags.Setup {
		_ = db.Configure()
		_, err = Setup(db, currentFlags.SetupSkipUpdatePatterns)
		return
	}

	var fabric *core.Fabric
	if err = db.Configure(); err != nil {
		fmt.Println("init is failed, run start the setup procedure", err)
		if fabric, err = Setup(db, currentFlags.SetupSkipUpdatePatterns); err != nil {
			return
		}
	} else {
		if fabric, err = core.NewFabric(db); err != nil {
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

		if err = db.Patterns.PrintLatestPatterns(parsedToInt); err != nil {
			return
		}
		return
	}

	// if the list patterns flag is set, run the list all patterns function
	if currentFlags.ListPatterns {
		err = db.Patterns.ListNames()
		return
	}

	// if the list all models flag is set, run the list all models function
	if currentFlags.ListAllModels {
		fabric.GetModels().Print()
		return
	}

	// if the list all contexts flag is set, run the list all contexts function
	if currentFlags.ListAllContexts {
		err = db.Contexts.ListNames()
		return
	}

	// if the list all sessions flag is set, run the list all sessions function
	if currentFlags.ListAllSessions {
		err = db.Sessions.ListNames()
		return
	}

	// Check for ScrapeURL flag first
	if currentFlags.ScrapeURL != "" {
		fmt.Println("ScrapeURL flag is set") // Debug print
		url := currentFlags.ScrapeURL
		curlCommand := fmt.Sprintf("curl https://r.jina.ai/%s", url)
		if err := exec.Command("sh", "-c", curlCommand).Run(); err != nil {
			return "", fmt.Errorf("failed to run curl command: %w", err)
		}
		os.Exit(0)
	}

	// if the interactive flag is set, run the interactive function
	// if currentFlags.Interactive {
	// 	interactive.Interactive()
	// }

	// if none of the above currentFlags are set, run the initiate chat function

	var chatter *core.Chatter
	if chatter, err = fabric.GetChatter(currentFlags.Model, currentFlags.Stream); err != nil {
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
	ret = core.NewFabricForSetup(db)

	if err = ret.Setup(); err != nil {
		return
	}

	if !skipUpdatePatterns {
		if err = ret.PopulateDB(); err != nil {
			return
		}
	}

	return
}
