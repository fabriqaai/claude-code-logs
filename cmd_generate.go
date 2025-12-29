package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate HTML pages from chat logs",
	Long: `Generate static HTML pages from Claude Code chat logs.

This command scans ~/.claude/projects for chat sessions and generates
HTML files in the output directory (default: ~/.claude-logs).

Example:
  claude-logs generate
  claude-logs generate --output-dir /tmp/my-logs`,
	RunE: runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) error {
	start := time.Now()

	// Get output directory
	outDir, err := getOutputDir()
	if err != nil {
		return err
	}

	logVerbose("Output directory: %s", outDir)

	// Check if output directory is writable
	if err := ensureWritableDir(outDir); err != nil {
		return fmt.Errorf("output directory not writable: %w", err)
	}

	// Get Claude projects path
	projectsPath, err := DefaultClaudeProjectsPath()
	if err != nil {
		return err
	}

	logVerbose("Scanning projects from: %s", projectsPath)

	// Load all projects and sessions
	fmt.Println("Discovering projects...")
	projects, err := LoadAllProjects(projectsPath)
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("No Claude projects found.")
		fmt.Println("Claude projects are typically stored in ~/.claude/projects")
		return nil
	}

	// Count total sessions
	totalSessions := 0
	for _, p := range projects {
		totalSessions += len(p.Sessions)
	}

	fmt.Printf("Found %d projects with %d sessions\n", len(projects), totalSessions)
	logVerbose("Projects:")
	for _, p := range projects {
		logVerbose("  %s (%d sessions)", p.Path, len(p.Sessions))
	}

	// Generate HTML
	fmt.Println("Generating HTML...")
	if err := GenerateAll(projects, outDir); err != nil {
		return fmt.Errorf("generating HTML: %w", err)
	}

	duration := time.Since(start)
	fmt.Printf("Generated HTML in %s\n", outDir)
	fmt.Printf("Completed in %v\n", duration.Round(time.Millisecond))

	return nil
}

// ensureWritableDir ensures the directory exists and is writable
func ensureWritableDir(path string) error {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}

	// Check if writable by creating a temp file
	testFile := path + "/.write-test"
	f, err := os.Create(testFile)
	if err != nil {
		return fmt.Errorf("directory not writable: %w", err)
	}
	f.Close()
	os.Remove(testFile)

	return nil
}
