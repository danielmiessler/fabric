package main

import (
	"fmt"
	"os"

	"github.com/danielmiessler/fabric/cli"
)

// use to get latest tag, go install -ldflags "-X main.version=$(git describe --tags --always)" github.com/danielmiessler/fabric@latest
var version = "dev" // Default version

func main() {
	_, err := cli.Cli(version)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}
