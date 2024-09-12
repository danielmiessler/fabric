package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	var input io.Reader
	var outputFile string
	if len(os.Args) > 1 {
		// File input mode
		file, err := os.Open(os.Args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
			os.Exit(1)
		}
		defer file.Close()
		input = file
		outputFile = strings.TrimSuffix(os.Args[1], filepath.Ext(os.Args[1])) + ".pdf"
	} else {
		// Stdin mode
		input = os.Stdin
		outputFile = "output.pdf"
	}

	// Check if pdflatex is installed
	if _, err := exec.LookPath("pdflatex"); err != nil {
		fmt.Fprintln(os.Stderr, "Error: pdflatex is not installed or not in your PATH.")
		fmt.Fprintln(os.Stderr, "Please install a LaTeX distribution (e.g., TeX Live or MiKTeX) and ensure pdflatex is in your PATH.")
		os.Exit(1)
	}

	// Create a temporary directory
	tmpDir, err := os.MkdirTemp("", "latex_")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary .tex file
	tmpFile, err := os.Create(filepath.Join(tmpDir, "input.tex"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating temporary file: %v\n", err)
		os.Exit(1)
	}

	// Copy input to the temporary file
	_, err = io.Copy(tmpFile, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to temporary file: %v\n", err)
		os.Exit(1)
	}
	tmpFile.Close()

	// Run pdflatex with nonstopmode
	cmd := exec.Command("pdflatex", "-interaction=nonstopmode", "-output-directory", tmpDir, tmpFile.Name())
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running pdflatex: %v\n", err)
		fmt.Fprintf(os.Stderr, "pdflatex output:\n%s\n", output)
		os.Exit(1)
	}

	// Check if PDF was actually created
	pdfPath := filepath.Join(tmpDir, "input.pdf")
	if _, err := os.Stat(pdfPath); os.IsNotExist(err) {
		fmt.Fprintln(os.Stderr, "Error: PDF file was not created. There might be an issue with your LaTeX source.")
		fmt.Fprintf(os.Stderr, "pdflatex output:\n%s\n", output)
		os.Exit(1)
	}

	// Move the output PDF to the current directory
	err = os.Rename(pdfPath, outputFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error moving output file: %v\n", err)
		os.Exit(1)
	}

	// Clean up temporary files
	cleanupTempFiles(tmpDir)

	fmt.Printf("PDF created: %s\n", outputFile)
}

func cleanupTempFiles(dir string) {
	extensions := []string{".aux", ".log", ".out", ".toc", ".lof", ".lot", ".bbl", ".blg"}
	for _, ext := range extensions {
		files, err := filepath.Glob(filepath.Join(dir, "*"+ext))
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error finding %s files: %v\n", ext, err)
			continue
		}
		for _, file := range files {
			if err := os.Remove(file); err != nil {
				fmt.Fprintf(os.Stderr, "Error removing file %s: %v\n", file, err)
			}
		}
	}
}
