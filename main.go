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
	dirFlag string
	verbose bool
)

// rootCmd is the base command for the CLI
var rootCmd = &cobra.Command{
	Use:   "claude-code-logs",
	Short: "Browse and search Claude Code chat logs",
	Long: `claude-code-logs generates HTML pages from Claude Code chat logs and serves them locally.

It scans ~/.claude/projects for chat sessions and creates a searchable web interface.`,
}

// Legacy commands that show migration messages
var generateCmd = &cobra.Command{
	Use:    "generate",
	Short:  "Generate HTML pages (DEPRECATED)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command removed. Use 'claude-code-logs serve' instead (generates automatically on startup)")
	},
}

var watchCmd = &cobra.Command{
	Use:    "watch",
	Short:  "Watch for changes (DEPRECATED)",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command removed. Use 'claude-code-logs serve --watch' instead")
	},
}

func init() {
	// Global flags available to all commands
	rootCmd.PersistentFlags().StringVarP(&dirFlag, "dir", "d", "", "Output directory for HTML (default: ~/claude-code-logs)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")

	// Add subcommands
	rootCmd.AddCommand(serveCmd)
	rootCmd.AddCommand(versionCmd)

	// Add hidden legacy commands for migration messages
	rootCmd.AddCommand(generateCmd)
	rootCmd.AddCommand(watchCmd)
}

// getOutputDir returns the output directory, expanding ~ to home directory
func getOutputDir() (string, error) {
	if dirFlag != "" {
		return expandPath(dirFlag)
	}

	// Default to ~/claude-code-logs
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("getting home directory: %w", err)
	}
	return filepath.Join(home, "claude-code-logs"), nil
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

