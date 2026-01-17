package main

import (
	"fmt"
	"os"
	"time"

	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
)

var (
	servePort  int
	serveWatch bool
	serveList  bool
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

With --list flag, you can interactively select which projects to serve.

Example:
  claude-code-logs serve
  claude-code-logs serve --port 3000
  claude-code-logs serve --dir /custom/path
  claude-code-logs serve --watch         (regenerates on changes)
  claude-code-logs serve --list          (select projects interactively)
  claude-code-logs serve --list --watch  (select projects + watch mode)`,
	RunE: runServe,
}

func init() {
	serveCmd.Flags().IntVarP(&servePort, "port", "p", 8080, "Server port")
	serveCmd.Flags().BoolVarP(&serveWatch, "watch", "w", false, "Enable watch mode (regenerate on changes)")
	serveCmd.Flags().BoolVarP(&serveList, "list", "l", false, "Interactively select projects to serve")
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

	// Load all projects
	fmt.Println("Discovering projects...")
	start := time.Now()
	projects, err := LoadAllProjects(projectsPath)
	if err != nil {
		return fmt.Errorf("loading projects: %w", err)
	}

	if len(projects) == 0 {
		fmt.Println("No Claude projects found.")
		fmt.Println("Claude projects are typically stored in ~/.claude/projects")
		fmt.Println("Server will start but will have no content to display.")
	}

	// Interactive project selection if --list flag is set
	var selectedFolders []string
	if serveList && len(projects) > 0 {
		projects, selectedFolders, err = selectProjects(projects)
		if err != nil {
			return fmt.Errorf("selecting projects: %w", err)
		}
		if len(projects) == 0 {
			return fmt.Errorf("no projects selected")
		}
		fmt.Printf("Selected %d projects\n", len(projects))
	}

	// Generate HTML
	if len(projects) > 0 {
		fmt.Println("Generating HTML...")
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
			SourceDir:        projectsPath,
			OutputDir:        outDir,
			PollInterval:     30 * time.Second,
			DebounceDelay:    2 * time.Second,
			SelectedProjects: selectedFolders, // nil means all projects
		}

		cancelWatch, err := WatchInBackground(config)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to start watcher: %v\n", err)
			fmt.Println("Server will start without automatic regeneration")
		} else {
			if len(selectedFolders) > 0 {
				fmt.Printf("Watch mode enabled for %d selected projects\n", len(selectedFolders))
			} else {
				fmt.Println("Watch mode enabled - HTML will regenerate on changes")
			}
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

// selectProjects displays an interactive multi-select UI for choosing projects
// Returns the filtered projects and their folder names
func selectProjects(projects []Project) ([]Project, []string, error) {
	// Build options with useful info
	options := make([]string, len(projects))
	for i, p := range projects {
		sessionCount := len(p.Sessions)
		lastActivity := "no sessions"
		if sessionCount > 0 {
			lastUpdate := getProjectLastUpdate(&p)
			if !lastUpdate.IsZero() {
				lastActivity = formatRelativeTime(lastUpdate)
			}
		}
		options[i] = fmt.Sprintf("%s (%d sessions, %s)", p.Path, sessionCount, lastActivity)
	}

	// Create multi-select prompt
	var selectedIndices []int
	prompt := &survey.MultiSelect{
		Message:  "Select projects to serve:",
		Options:  options,
		PageSize: 15,
	}

	if err := survey.AskOne(prompt, &selectedIndices); err != nil {
		return nil, nil, err
	}

	// Build filtered results
	selected := make([]Project, 0, len(selectedIndices))
	folders := make([]string, 0, len(selectedIndices))
	for _, idx := range selectedIndices {
		selected = append(selected, projects[idx])
		folders = append(folders, projects[idx].FolderName)
	}

	return selected, folders, nil
}

// formatRelativeTime formats a time as a human-readable relative string
func formatRelativeTime(t time.Time) string {
	diff := time.Since(t)

	switch {
	case diff < time.Minute:
		return "just now"
	case diff < time.Hour:
		mins := int(diff.Minutes())
		if mins == 1 {
			return "1 min ago"
		}
		return fmt.Sprintf("%d mins ago", mins)
	case diff < 24*time.Hour:
		hours := int(diff.Hours())
		if hours == 1 {
			return "1 hour ago"
		}
		return fmt.Sprintf("%d hours ago", hours)
	case diff < 7*24*time.Hour:
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	default:
		weeks := int(diff.Hours() / 24 / 7)
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}
}
