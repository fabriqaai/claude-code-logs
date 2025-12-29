package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	watchInterval int
	watchDebounce int
)

var watchCmd = &cobra.Command{
	Use:   "watch",
	Short: "Watch for changes and regenerate HTML",
	Long: `Watch for changes in Claude Code chat logs and regenerate HTML automatically.

This command monitors ~/.claude/projects for new or modified sessions
and regenerates the affected HTML pages using fsnotify for file system events.

Features:
- Detects new and modified .jsonl files
- Debounces rapid changes (default: 2 seconds)
- Only regenerates affected projects (incremental)
- Graceful shutdown with Ctrl+C

Example:
  claude-logs watch
  claude-logs watch --interval 60
  claude-logs watch --debounce 5`,
	RunE: runWatch,
}

func init() {
	watchCmd.Flags().IntVarP(&watchInterval, "interval", "i", 30, "Poll interval for new directories (seconds)")
	watchCmd.Flags().IntVarP(&watchDebounce, "debounce", "d", 2, "Debounce delay before regenerating (seconds)")
}

func runWatch(cmd *cobra.Command, args []string) error {
	// Get directories
	outDir, err := getOutputDir()
	if err != nil {
		return err
	}

	sourceDir, err := DefaultClaudeProjectsPath()
	if err != nil {
		return err
	}

	// Validate source directory exists
	if _, err := os.Stat(sourceDir); os.IsNotExist(err) {
		return fmt.Errorf("claude projects directory not found: %s", sourceDir)
	}

	// Run initial generation if output doesn't exist
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		fmt.Println("Output directory does not exist. Running initial generation...")
		if err := runGenerate(cmd, args); err != nil {
			return fmt.Errorf("initial generation failed: %w", err)
		}
		fmt.Println()
	}

	// Configure watcher
	config := WatchConfig{
		SourceDir:     sourceDir,
		OutputDir:     outDir,
		PollInterval:  time.Duration(watchInterval) * time.Second,
		DebounceDelay: time.Duration(watchDebounce) * time.Second,
	}

	logVerbose("Source directory: %s", sourceDir)
	logVerbose("Output directory: %s", outDir)
	logVerbose("Poll interval: %d seconds", watchInterval)
	logVerbose("Debounce delay: %d seconds", watchDebounce)

	// Set up context with signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle Ctrl+C and termination signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		fmt.Println("\nReceived shutdown signal...")
		cancel()
	}()

	fmt.Println("Starting file watcher...")
	fmt.Println("Press Ctrl+C to stop")
	fmt.Println()

	// Start watching
	return StartWatcher(ctx, config)
}
