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
	Short: "Generate HTML and start the web server",
	Long: `Generate HTML from chat logs and start a local web server.

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
  claude-logs serve --dir /custom/path
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

	// Check if output directory is writable (creates if needed)
	if err := ensureWritableDir(outDir); err != nil {
		return fmt.Errorf("output directory not writable: %w", err)
	}

	// Get Claude projects path
	projectsPath, err := DefaultClaudeProjectsPath()
	if err != nil {
		return err
	}

	// Always generate HTML on startup
	fmt.Println("Generating HTML...")
	start := time.Now()
	projects, err := LoadAllProjects(projectsPath)
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("No Claude projects found.")
		fmt.Println("Claude projects are typically stored in ~/.claude/projects")
		fmt.Println("Server will start but will have no content to display.")
	} else {
		// Count total sessions
		totalSessions := 0
		for _, p := range projects {
			totalSessions += len(p.Sessions)
		}
		fmt.Printf("Found %d projects with %d sessions\n", len(projects), totalSessions)

		if err := GenerateAll(projects, outDir); err != nil {
			return fmt.Errorf("generating HTML: %w", err)
		}
		fmt.Printf("Generated HTML in %v\n", time.Since(start).Round(time.Millisecond))
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
