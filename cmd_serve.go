package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	servePort  int
	serveWatch bool
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the web server",
	Long: `Start a local web server to browse chat logs.

The server provides:
- Static HTML pages for browsing projects and sessions
- Search API for full-text search across all messages
- Real-time stats endpoint

With --watch flag, the server will also:
- Monitor ~/.claude/projects for changes
- Automatically regenerate HTML when sessions are modified
- Debounce rapid changes for efficiency

Example:
  claude-logs serve
  claude-logs serve --port 3000
  claude-logs serve --watch  (regenerates on changes)`,
	RunE: runServe,
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Server port")
	serveCmd.Flags().BoolVarP(&serveWatch, "watch", "w", false, "Enable watch mode (regenerate on changes)")
}

func runServe(cmd *cobra.Command, args []string) error {
	// Validate port
	if servePort < 1 || servePort > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1-65535)", servePort)
	}

	// Get output directory
	outDir, err := getOutputDir()
	if err != nil {
		return err
	}

	logVerbose("Output directory: %s", outDir)
	logVerbose("Port: %d", servePort)
	logVerbose("Watch mode: %v", serveWatch)

	// Check if output directory exists
	if _, err := os.Stat(outDir); os.IsNotExist(err) {
		fmt.Println("Output directory does not exist. Running generate first...")
		if err := runGenerate(cmd, args); err != nil {
			return err
		}
	}

	// Get Claude projects path
	projectsPath, err := DefaultClaudeProjectsPath()
	if err != nil {
		return err
	}

	// Load projects for search index
	fmt.Println("Loading projects for search index...")
	projects, err := LoadAllProjects(projectsPath)
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("Warning: No projects found. Server will start but search will be empty.")
	}

	// Handle watch mode
	if serveWatch {
		config := WatchConfig{
			SourceDir:     projectsPath,
			OutputDir:     outDir,
			PollInterval:  30 * time.Second,
			DebounceDelay: 2 * time.Second,
		}

		cancelWatch, err := WatchInBackground(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to start watcher: %v\n", err)
			fmt.Println("Server will start without automatic regeneration")
		} else {
			fmt.Println("Watch mode enabled - HTML will regenerate on changes")
			defer cancelWatch()
		}
	}

	// Start server
	fmt.Printf("Starting server on http://127.0.0.1:%d\n", servePort)
	return StartServer(servePort, outDir, projects)
}
