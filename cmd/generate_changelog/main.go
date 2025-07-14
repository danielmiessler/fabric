package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/danielmiessler/fabric/cmd/generate_changelog/internal/changelog"
	"github.com/danielmiessler/fabric/cmd/generate_changelog/internal/config"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
	cfg = &config.Config{}
)

var rootCmd = &cobra.Command{
	Use:   "generate_changelog",
	Short: "Generate changelog from git history and GitHub PRs",
	Long: `A high-performance changelog generator that walks git history,
collects version information and pull requests, and generates a
comprehensive changelog in markdown format.`,
	RunE: run,
}

func init() {
	rootCmd.Flags().StringVarP(&cfg.RepoPath, "repo", "r", ".", "Repository path")
	rootCmd.Flags().StringVarP(&cfg.OutputFile, "output", "o", "", "Output file (default: stdout)")
	rootCmd.Flags().IntVarP(&cfg.Limit, "limit", "l", 0, "Limit number of versions (0 = all)")
	rootCmd.Flags().StringVarP(&cfg.Version, "version", "v", "", "Generate changelog for specific version")
	rootCmd.Flags().BoolVar(&cfg.SaveData, "save-data", false, "Save version data to JSON for debugging")
	rootCmd.Flags().StringVar(&cfg.CacheFile, "cache", "./cmd/generate_changelog/changelog.db", "Cache database file")
	rootCmd.Flags().BoolVar(&cfg.NoCache, "no-cache", false, "Disable cache usage")
	rootCmd.Flags().BoolVar(&cfg.RebuildCache, "rebuild-cache", false, "Rebuild cache from scratch")
	rootCmd.Flags().StringVar(&cfg.GitHubToken, "token", "", "GitHub API token (or set GITHUB_TOKEN env var)")
	rootCmd.Flags().BoolVar(&cfg.ForcePRSync, "force-pr-sync", false, "Force a full PR sync from GitHub (ignores cache age)")
	rootCmd.Flags().BoolVar(&cfg.EnableAISummary, "ai-summarize", false, "Generate AI-enhanced summaries using Fabric")
}

func run(cmd *cobra.Command, args []string) error {
	if cfg.GitHubToken == "" {
		cfg.GitHubToken = os.Getenv("GITHUB_TOKEN")
	}

	generator, err := changelog.New(cfg)
	if err != nil {
		return fmt.Errorf("failed to create changelog generator: %w", err)
	}

	output, err := generator.Generate()
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %w", err)
	}

	if cfg.OutputFile != "" {
		if err := os.WriteFile(cfg.OutputFile, []byte(output), 0644); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Printf("Changelog written to %s\n", cfg.OutputFile)
	} else {
		fmt.Print(output)
	}

	return nil
}

func main() {
	// Load .env file from the same directory as the binary
	if exePath, err := os.Executable(); err == nil {
		envPath := filepath.Join(filepath.Dir(exePath), ".env")
		if _, err := os.Stat(envPath); err == nil {
			// .env file exists, load it
			if err := godotenv.Load(envPath); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to load .env file: %v\n", err)
			}
		}
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
