package main

import (
	"fmt"
	"os"
)

// Version is set at build time via ldflags
var Version = "dev"

func main() {
	// Placeholder - CLI will be implemented in bolt 004-cli
	fmt.Println("claude-code-logs", Version)
	fmt.Println("CLI not yet implemented. Use parser functions programmatically.")

	// Quick test of parser
	projectsPath, err := DefaultClaudeProjectsPath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	projects, err := DiscoverProjects(projectsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFound %d projects:\n", len(projects))
	for _, p := range projects {
		fmt.Printf("  %s\n", p.Path)
	}
}
