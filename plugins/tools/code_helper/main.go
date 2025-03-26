package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	// Define command line flags
	maxDepth := flag.Int("depth", 3, "Maximum directory depth to scan")
	ignorePatterns := flag.String("ignore", ".git,node_modules,vendor", "Comma-separated patterns to ignore")
	outputFile := flag.String("out", "", "Output file (default: stdout)")
	showHelp := flag.Bool("help", false, "Show help message")

	// Parse command line flags
	flag.Parse()

	// Show help if requested or no arguments provided
	if *showHelp || flag.NArg() < 1 {
		printUsage()
		os.Exit(0)
	}

	// Get directory and instructions from positional arguments
	directory := flag.Arg(0)
	instructions := ""
	if flag.NArg() > 1 {
		instructions = flag.Arg(1)
	}

	// Validate directory
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: Directory '%s' does not exist\n", directory)
		os.Exit(1)
	}
	// Parse ignore patterns
	ignoreList := ParseIgnorePatterns(*ignorePatterns)

	// Scan directory and generate JSON
	var jsonData []byte
	var err error
	jsonData, err = ScanDirectory(directory, *maxDepth, instructions, ignoreList)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Write output
	if *outputFile != "" {
		if err := os.WriteFile(*outputFile, jsonData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
			os.Exit(1)
		}
	} else {
		fmt.Print(string(jsonData))
	}
}

// ParseIgnorePatterns converts a comma-separated string of patterns to a slice
func ParseIgnorePatterns(patterns string) []string {
	if patterns == "" {
		return nil
	}
	return strings.Split(patterns, ",")
}

func printUsage() {
	fmt.Println(`code_helper - Code project scanner for use with Fabric AI

Usage:
  code_helper [options] <directory> [instructions]

Examples:
  # Scan current directory with instructions
  code_helper . "Add input validation to all user inputs"

  # Scan specific directory with depth limit
  code_helper -depth 4 ./my-project "Implement error handling"

  # Output to file instead of stdout
  code_helper -out project.json ./src "Fix security issues"

Options:`)
	flag.PrintDefaults()
}
