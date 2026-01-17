package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// MarkdownGenerator handles Markdown file generation from parsed session data
type MarkdownGenerator struct {
	outputDir string
	sourceDir string
	force     bool
}

// GenerationResult contains statistics about the generation process
type GenerationResult struct {
	Generated int
	Skipped   int
	Errors    []error
}

// NewMarkdownGenerator creates a new MarkdownGenerator instance
func NewMarkdownGenerator(outputDir, sourceDir string, force bool) *MarkdownGenerator {
	return &MarkdownGenerator{
		outputDir: outputDir,
		sourceDir: sourceDir,
		force:     force,
	}
}

// GenerateAll generates Markdown files for all projects and sessions
func (g *MarkdownGenerator) GenerateAll(projects []Project) (*GenerationResult, error) {
	result := &GenerationResult{}

	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return nil, fmt.Errorf("creating output directory: %w", err)
	}

	// Generate session files for each project
	for i := range projects {
		project := &projects[i]
		projectSlug := ProjectSlug(project.Path)
		projectDir := filepath.Join(g.outputDir, projectSlug)

		// Create project directory
		if err := os.MkdirAll(projectDir, 0755); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("creating project dir %s: %w", projectSlug, err))
			continue
		}

		// Generate session files
		for j := range project.Sessions {
			session := &project.Sessions[j]
			mdPath := filepath.Join(projectDir, session.ID+".md")

			// Check if regeneration is needed
			if !g.ShouldRegenerate(session.SourcePath, mdPath) {
				result.Skipped++
				continue
			}

			// Generate session markdown
			if err := g.GenerateSession(session, projectSlug); err != nil {
				result.Errors = append(result.Errors, fmt.Errorf("session %s: %w", session.ID, err))
				continue
			}
			result.Generated++
		}

		// Generate project index (MD)
		if err := g.GenerateProjectIndex(project, projectSlug); err != nil {
			result.Errors = append(result.Errors, fmt.Errorf("project index %s: %w", projectSlug, err))
		}
	}

	// Generate main index (MD)
	if err := g.GenerateMainIndex(projects); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("main index: %w", err))
	}

	return result, nil
}

// GenerateSession generates a Markdown file for a single session
func (g *MarkdownGenerator) GenerateSession(session *Session, projectSlug string) error {
	// Compute source hash
	sourceHash, err := ComputeFileHash(session.SourcePath)
	if err != nil {
		// Log warning but continue with "unknown" hash
		sourceHash = "unknown"
	}

	// Create frontmatter
	fm := NewFrontmatter(session, sourceHash)
	fmBytes, err := fm.Marshal()
	if err != nil {
		return fmt.Errorf("marshaling frontmatter: %w", err)
	}

	// Build markdown content
	var content strings.Builder
	content.Write(fmBytes)
	content.WriteString("\n")

	// Format messages
	for _, msg := range session.Messages {
		content.WriteString(g.formatMessage(&msg))
		content.WriteString("\n")
	}

	// Write MD file
	mdPath := filepath.Join(g.outputDir, projectSlug, session.ID+".md")
	return g.writeFile(mdPath, []byte(content.String()))
}

// GenerateMainIndex generates the main index.md listing all projects
func (g *MarkdownGenerator) GenerateMainIndex(projects []Project) error {
	var content strings.Builder

	content.WriteString("# Claude Code Logs\n\n")
	content.WriteString("| Project | Sessions | Last Activity |\n")
	content.WriteString("|---------|----------|---------------|\n")

	// Sort projects by last update (most recent first) - already sorted by parser
	for _, project := range projects {
		projectSlug := ProjectSlug(project.Path)
		sessionCount := len(project.Sessions)

		lastActivity := ""
		if sessionCount > 0 {
			lastUpdate := getProjectLastUpdate(&project)
			if !lastUpdate.IsZero() {
				lastActivity = lastUpdate.Format("2006-01-02")
			}
		}

		content.WriteString(fmt.Sprintf("| [%s](%s/index.md) | %d | %s |\n",
			escapeMarkdownTableCell(project.Path),
			projectSlug,
			sessionCount,
			lastActivity,
		))
	}

	outputPath := filepath.Join(g.outputDir, "index.md")
	return g.writeFile(outputPath, []byte(content.String()))
}

// GenerateProjectIndex generates the project index.md listing all sessions
func (g *MarkdownGenerator) GenerateProjectIndex(project *Project, projectSlug string) error {
	var content strings.Builder

	content.WriteString(fmt.Sprintf("# Project: %s\n\n", project.Path))
	content.WriteString("| Session | Title | Created |\n")
	content.WriteString("|---------|-------|--------|\n")

	// Sort sessions by creation date (most recent first)
	sessions := make([]Session, len(project.Sessions))
	copy(sessions, project.Sessions)
	sort.Slice(sessions, func(i, j int) bool {
		return sessions[i].CreatedAt.After(sessions[j].CreatedAt)
	})

	for _, session := range sessions {
		created := session.CreatedAt.Format("2006-01-02")
		content.WriteString(fmt.Sprintf("| [%s](%s.md) | %s | %s |\n",
			session.ID,
			session.ID,
			escapeMarkdownTableCell(session.Summary),
			created,
		))
	}

	outputPath := filepath.Join(g.outputDir, projectSlug, "index.md")
	return g.writeFile(outputPath, []byte(content.String()))
}

// ShouldRegenerate determines if a session needs regeneration based on mtime
func (g *MarkdownGenerator) ShouldRegenerate(jsonlPath, mdPath string) bool {
	if g.force {
		return true
	}

	mdInfo, err := os.Stat(mdPath)
	if os.IsNotExist(err) {
		return true // MD doesn't exist, generate it
	}
	if err != nil {
		return true // Error checking, regenerate to be safe
	}

	jsonlInfo, err := os.Stat(jsonlPath)
	if err != nil {
		return false // JSONL doesn't exist (orphan case)
	}

	// Regenerate if JSONL is newer than MD
	return jsonlInfo.ModTime().After(mdInfo.ModTime())
}

// formatMessage formats a single message as Markdown
func (g *MarkdownGenerator) formatMessage(msg *Message) string {
	var content strings.Builder

	// Role header (capitalize first letter)
	role := capitalizeFirst(msg.Role)
	content.WriteString(fmt.Sprintf("## %s\n\n", role))

	// Format content blocks
	for _, block := range msg.Content {
		switch block.Type {
		case "text":
			content.WriteString(block.Text)
			content.WriteString("\n\n")

		case "tool_use":
			content.WriteString(fmt.Sprintf("<details>\n<summary>Tool: %s</summary>\n\n", block.ToolName))
			content.WriteString("```json\n")
			content.WriteString(block.ToolInput)
			content.WriteString("\n```\n\n")
			content.WriteString("</details>\n\n")

		case "tool_result":
			content.WriteString("<details>\n<summary>Result</summary>\n\n")
			content.WriteString("```\n")
			// Truncate very long tool outputs
			output := block.ToolOutput
			if len(output) > 10000 {
				output = output[:10000] + "\n... (truncated)"
			}
			content.WriteString(output)
			content.WriteString("\n```\n\n")
			content.WriteString("</details>\n\n")
		}
	}

	return content.String()
}

// writeFile writes content to a file using atomic write
func (g *MarkdownGenerator) writeFile(outputPath string, content []byte) error {
	dir := filepath.Dir(outputPath)

	// Create temp file in the same directory for atomic rename
	tmpFile, err := os.CreateTemp(dir, "tmp-*.md")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Clean up temp file on error
	defer func() {
		if tmpPath != "" {
			os.Remove(tmpPath)
		}
	}()

	// Write content
	if _, err := tmpFile.Write(content); err != nil {
		tmpFile.Close()
		return fmt.Errorf("writing content: %w", err)
	}

	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Atomic rename
	if err := os.Rename(tmpPath, outputPath); err != nil {
		return fmt.Errorf("renaming temp file: %w", err)
	}

	// Clear tmpPath so defer doesn't try to remove it
	tmpPath = ""
	return nil
}

// escapeMarkdownTableCell escapes characters that would break markdown tables
func escapeMarkdownTableCell(s string) string {
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "\n", " ")
	return s
}

// capitalizeFirst capitalizes the first letter of a string
func capitalizeFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// GenerateAllMarkdown is the main entry point for generating all Markdown files
func GenerateAllMarkdown(projects []Project, outputDir, sourceDir string, force bool) (*GenerationResult, error) {
	gen := NewMarkdownGenerator(outputDir, sourceDir, force)
	return gen.GenerateAll(projects)
}
