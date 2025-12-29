package main

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Generator handles HTML generation from parsed session data
type Generator struct {
	outputDir   string
	projects    []Project
	indexTmpl   *template.Template
	projectTmpl *template.Template
	sessionTmpl *template.Template
}

// NewGenerator creates a new Generator instance
func NewGenerator(outputDir string, projects []Project) (*Generator, error) {
	g := &Generator{
		outputDir: outputDir,
		projects:  projects,
	}

	// Create template functions
	funcMap := template.FuncMap{
		"ProjectSlug": ProjectSlug,
		"RenderText":  RenderText,
	}

	// Parse templates
	var err error
	g.indexTmpl, err = template.New("index").Funcs(funcMap).Parse(indexTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing index template: %w", err)
	}

	g.projectTmpl, err = template.New("project").Funcs(funcMap).Parse(projectIndexTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing project template: %w", err)
	}

	g.sessionTmpl, err = template.New("session").Funcs(funcMap).Parse(sessionTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing session template: %w", err)
	}

	return g, nil
}

// GenerateAll generates the complete static site
func (g *Generator) GenerateAll() error {
	// Create output directory
	if err := os.MkdirAll(g.outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	// Generate main index
	if err := g.GenerateIndex(); err != nil {
		return fmt.Errorf("generating index: %w", err)
	}

	// Generate project pages and session pages
	for i := range g.projects {
		project := &g.projects[i]
		if err := g.GenerateProjectPages(project); err != nil {
			return fmt.Errorf("generating project %s: %w", project.Path, err)
		}
	}

	return nil
}

// GenerateIndex generates the main index.html
func (g *Generator) GenerateIndex() error {
	data := struct {
		Projects    []Project
		ProjectSlug func(string) string
	}{
		Projects:    g.projects,
		ProjectSlug: ProjectSlug,
	}

	outputPath := filepath.Join(g.outputDir, "index.html")
	return g.writeTemplate(g.indexTmpl, data, outputPath)
}

// GenerateProjectPages generates the project index and all session pages
func (g *Generator) GenerateProjectPages(project *Project) error {
	// Create project directory
	projectDir := filepath.Join(g.outputDir, ProjectSlug(project.Path))
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("creating project directory: %w", err)
	}

	// Generate project index
	if err := g.GenerateProjectIndex(project); err != nil {
		return err
	}

	// Generate session pages
	for i := range project.Sessions {
		session := &project.Sessions[i]
		if err := g.GenerateSession(session, project); err != nil {
			return fmt.Errorf("generating session %s: %w", session.ID, err)
		}
	}

	return nil
}

// GenerateProjectIndex generates the project's index.html
func (g *Generator) GenerateProjectIndex(project *Project) error {
	data := struct {
		Project     *Project
		AllProjects []Project
		ProjectSlug func(string) string
	}{
		Project:     project,
		AllProjects: g.projects,
		ProjectSlug: ProjectSlug,
	}

	outputPath := filepath.Join(g.outputDir, ProjectSlug(project.Path), "index.html")
	return g.writeTemplate(g.projectTmpl, data, outputPath)
}

// GenerateSession generates a session HTML file
func (g *Generator) GenerateSession(session *Session, project *Project) error {
	data := struct {
		Session     *Session
		Project     *Project
		AllProjects []Project
		ProjectSlug func(string) string
		RenderText  func(string) template.HTML
	}{
		Session:     session,
		Project:     project,
		AllProjects: g.projects,
		ProjectSlug: ProjectSlug,
		RenderText:  RenderText,
	}

	outputPath := filepath.Join(g.outputDir, ProjectSlug(project.Path), session.ID+".html")
	return g.writeTemplate(g.sessionTmpl, data, outputPath)
}

// writeTemplate writes a template to a file using atomic write
func (g *Generator) writeTemplate(tmpl *template.Template, data interface{}, outputPath string) error {
	// Create temp file in the same directory for atomic rename
	dir := filepath.Dir(outputPath)
	tmpFile, err := os.CreateTemp(dir, "tmp-*.html")
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

	// Execute template to temp file
	if err := tmpl.Execute(tmpFile, data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("executing template: %w", err)
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

// ProjectSlug converts a project path to a URL-safe slug
// Example: "/Users/john/project" -> "users-john-project"
func ProjectSlug(path string) string {
	if path == "" {
		return "unknown"
	}

	// Remove leading slash
	slug := strings.TrimPrefix(path, "/")

	// Convert to lowercase
	slug = strings.ToLower(slug)

	// Replace non-alphanumeric with dashes
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading/trailing dashes
	slug = strings.Trim(slug, "-")

	if slug == "" {
		return "project"
	}

	return slug
}

// RenderText converts plain text to HTML with basic formatting
// Handles code blocks, inline code, and paragraphs
func RenderText(text string) template.HTML {
	if text == "" {
		return ""
	}

	// Process code blocks first (``` ... ```)
	codeBlockRe := regexp.MustCompile("(?s)```([a-zA-Z]*)\n?(.*?)```")
	text = codeBlockRe.ReplaceAllStringFunc(text, func(match string) string {
		parts := codeBlockRe.FindStringSubmatch(match)
		if len(parts) >= 3 {
			code := template.HTMLEscapeString(parts[2])
			return "<pre><code>" + code + "</code></pre>"
		}
		return match
	})

	// Process inline code (` ... `)
	inlineCodeRe := regexp.MustCompile("`([^`]+)`")
	text = inlineCodeRe.ReplaceAllStringFunc(text, func(match string) string {
		parts := inlineCodeRe.FindStringSubmatch(match)
		if len(parts) >= 2 {
			code := template.HTMLEscapeString(parts[1])
			return "<code>" + code + "</code>"
		}
		return match
	})

	// Split into paragraphs and wrap non-code blocks
	var result strings.Builder
	lines := strings.Split(text, "\n\n")
	for _, para := range lines {
		para = strings.TrimSpace(para)
		if para == "" {
			continue
		}
		// Don't wrap if it's already a pre/code block
		if strings.HasPrefix(para, "<pre>") {
			result.WriteString(para)
		} else {
			// Escape HTML in regular paragraphs (but preserve our code tags)
			if !strings.Contains(para, "<code>") && !strings.Contains(para, "</code>") {
				para = template.HTMLEscapeString(para)
			}
			result.WriteString("<p>")
			result.WriteString(para)
			result.WriteString("</p>")
		}
	}

	return template.HTML(result.String())
}

// GenerateAll is the main entry point for generating all HTML
func GenerateAll(projects []Project, outputDir string) error {
	gen, err := NewGenerator(outputDir, projects)
	if err != nil {
		return err
	}
	return gen.GenerateAll()
}
