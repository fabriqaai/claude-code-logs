package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// Build info set at build time via ldflags
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Global flags
var (
	outputDir string
	verbose   bool
)

// rootCmd is the base command for the CLI
var rootCmd = &cobra.Command{
	Use:   "claude-logs",
	Short: "Browse and search Claude Code chat logs",
	Long: `claude-logs generates HTML pages from Claude Code chat logs and serves them locally.

It scans ~/.claude/projects for chat sessions and creates a searchable web interface.`,
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&outputDir, "output-dir", "o", "", "Output directory for HTML (default: ~/.claude-logs)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(watchCmd)
	rootCmd.AddCommand(versionCmd)
}

// getOutputDir returns the output directory, expanding ~ to home directory
func getOutputDir() (string, error) {
	if outputDir != "" {
		return expandPath(outputDir)
	}

	// Default to ~/.claude-logs
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, ".claude-logs"), nil
}

// expandPath expands ~ to the user's home directory
func expandPath(path string) (string, error) {
	if len(path) == 0 {
		return path, nil
	}

	if path[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("expanding home directory: %w", err)
		}
		return filepath.Join(home, path[1:]), nil
	}

	return path, nil
}

// logVerbose prints a message if verbose mode is enabled
func logVerbose(format string, args ...interface{}) {
	if verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

