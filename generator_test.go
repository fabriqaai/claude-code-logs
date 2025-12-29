package main

import (
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestProjectSlug(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
	}{
		{
			name:     "simple path",
			path:     "/Users/john/project",
			expected: "users-john-project",
		},
		{
			name:     "path with spaces",
			path:     "/Users/john/my project",
			expected: "users-john-my-project",
		},
		{
			name:     "path with special chars",
			path:     "/Users/john/project@v2.0",
			expected: "users-john-project-v2-0",
		},
		{
			name:     "empty path",
			path:     "",
			expected: "unknown",
		},
		{
			name:     "root path",
			path:     "/",
			expected: "project",
		},
		{
			name:     "path with uppercase",
			path:     "/Users/John/MyProject",
			expected: "users-john-myproject",
		},
		{
			name:     "windows-like path",
			path:     "C:/Users/john/project",
			expected: "c-users-john-project",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ProjectSlug(tt.path)
			if result != tt.expected {
				t.Errorf("ProjectSlug(%q) = %q, want %q", tt.path, result, tt.expected)
			}
		})
	}
}

func TestRenderText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		contains []string
		excludes []string
	}{
		{
			name:     "plain text",
			input:    "Hello world",
			contains: []string{"<p>Hello world</p>"},
		},
		{
			name:     "inline code",
			input:    "Use `go run` to execute",
			contains: []string{"<code>go run</code>"},
		},
		{
			name:     "code block",
			input:    "```go\nfmt.Println(\"hello\")\n```",
			contains: []string{"<pre><code>", "fmt.Println", "</code></pre>"},
		},
		{
			name:     "multiple paragraphs",
			input:    "First paragraph\n\nSecond paragraph",
			contains: []string{"<p>First paragraph</p>", "<p>Second paragraph</p>"},
		},
		{
			name:     "empty text",
			input:    "",
			contains: []string{},
		},
		{
			name:     "HTML escaping",
			input:    "Use <script> tags",
			contains: []string{"&lt;script&gt;"},
			excludes: []string{"<script>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := string(RenderText(tt.input))

			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("RenderText(%q) should contain %q, got %q", tt.input, s, result)
				}
			}

			for _, s := range tt.excludes {
				if strings.Contains(result, s) {
					t.Errorf("RenderText(%q) should not contain %q, got %q", tt.input, s, result)
				}
			}
		})
	}
}

func TestNewGenerator(t *testing.T) {
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions:   []Session{},
		},
	}

	gen, err := NewGenerator("/tmp/output", projects)
	if err != nil {
		t.Fatalf("NewGenerator failed: %v", err)
	}

	if gen.outputDir != "/tmp/output" {
		t.Errorf("outputDir = %q, want %q", gen.outputDir, "/tmp/output")
	}

	if len(gen.projects) != 1 {
		t.Errorf("projects count = %d, want 1", len(gen.projects))
	}

	if gen.indexTmpl == nil {
		t.Error("indexTmpl should not be nil")
	}

	if gen.projectTmpl == nil {
		t.Error("projectTmpl should not be nil")
	}

	if gen.sessionTmpl == nil {
		t.Error("sessionTmpl should not be nil")
	}
}

func TestGenerateAll(t *testing.T) {
	// Create temp directory for output
	tmpDir, err := os.MkdirTemp("", "generator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test data
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:        "test-session-1",
					Summary:   "Test Session One",
					CreatedAt: now,
					UpdatedAt: now,
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Hello Claude"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "Hello! How can I help you?"},
							},
						},
					},
				},
			},
		},
	}

	// Generate all pages
	err = GenerateAll(projects, tmpDir)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}

	// Verify index.html exists
	indexPath := filepath.Join(tmpDir, "index.html")
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		t.Error("index.html was not created")
	}

	// Verify project directory exists
	projectDir := filepath.Join(tmpDir, "users-test-project1")
	if _, err := os.Stat(projectDir); os.IsNotExist(err) {
		t.Error("project directory was not created")
	}

	// Verify project index exists
	projectIndex := filepath.Join(projectDir, "index.html")
	if _, err := os.Stat(projectIndex); os.IsNotExist(err) {
		t.Error("project index.html was not created")
	}

	// Verify session page exists
	sessionPage := filepath.Join(projectDir, "test-session-1.html")
	if _, err := os.Stat(sessionPage); os.IsNotExist(err) {
		t.Error("session page was not created")
	}

	// Verify index.html content
	indexContent, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index.html: %v", err)
	}

	// Check for key elements
	checks := []string{
		"<!DOCTYPE html>",
		"Claude Code Logs",
		"/Users/test/project1",
		"fabriqa.ai",
		"Static Mode",
	}

	for _, check := range checks {
		if !strings.Contains(string(indexContent), check) {
			t.Errorf("index.html should contain %q", check)
		}
	}

	// Verify session page content
	sessionContent, err := os.ReadFile(sessionPage)
	if err != nil {
		t.Fatalf("Failed to read session page: %v", err)
	}

	sessionChecks := []string{
		"Test Session One",
		"Hello Claude",
		"Hello! How can I help you?",
		"user",
		"assistant",
	}

	for _, check := range sessionChecks {
		if !strings.Contains(string(sessionContent), check) {
			t.Errorf("session page should contain %q", check)
		}
	}
}

func TestGenerateEmptyProject(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Empty project (no sessions)
	projects := []Project{
		{
			Path:       "/Users/test/empty",
			FolderName: "-Users-test-empty",
			Sessions:   []Session{},
		},
	}

	err = GenerateAll(projects, tmpDir)
	if err != nil {
		t.Fatalf("GenerateAll failed for empty project: %v", err)
	}

	// Verify project index shows empty state
	projectIndex := filepath.Join(tmpDir, "users-test-empty", "index.html")
	content, err := os.ReadFile(projectIndex)
	if err != nil {
		t.Fatalf("Failed to read project index: %v", err)
	}

	if !strings.Contains(string(content), "No sessions found") {
		t.Error("Empty project should show 'No sessions found' message")
	}
}

func TestGenerateNoProjects(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-no-projects-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// No projects at all
	projects := []Project{}

	err = GenerateAll(projects, tmpDir)
	if err != nil {
		t.Fatalf("GenerateAll failed with no projects: %v", err)
	}

	// Verify index shows empty state
	indexPath := filepath.Join(tmpDir, "index.html")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		t.Fatalf("Failed to read index: %v", err)
	}

	if !strings.Contains(string(content), "No projects found") {
		t.Error("Empty index should show 'No projects found' message")
	}
}

func TestGenerateToolBlocks(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "generator-tools-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/tools",
			FolderName: "-Users-test-tools",
			Sessions: []Session{
				{
					ID:        "tool-session",
					Summary:   "Tool Test",
					CreatedAt: now,
					UpdatedAt: now,
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "assistant",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Let me check that file."},
								{Type: "tool_use", ToolName: "Read", ToolInput: `{"path": "/test/file.txt"}`},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "user",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "tool_result", ToolOutput: "File contents here"},
							},
						},
					},
				},
			},
		},
	}

	err = GenerateAll(projects, tmpDir)
	if err != nil {
		t.Fatalf("GenerateAll failed: %v", err)
	}

	sessionPage := filepath.Join(tmpDir, "users-test-tools", "tool-session.html")
	content, err := os.ReadFile(sessionPage)
	if err != nil {
		t.Fatalf("Failed to read session page: %v", err)
	}

	toolChecks := []string{
		"tool-block",
		"Read",
		"toggleTool",
	}

	for _, check := range toolChecks {
		if !strings.Contains(string(content), check) {
			t.Errorf("Session with tools should contain %q", check)
		}
	}
}

func TestRenderTextHTMLSafety(t *testing.T) {
	// Test that HTML is properly escaped
	dangerous := "<script>alert('xss')</script>"
	result := RenderText(dangerous)

	if strings.Contains(string(result), "<script>") {
		t.Error("RenderText should escape HTML script tags")
	}

	if !strings.Contains(string(result), "&lt;script&gt;") {
		t.Error("RenderText should convert < to &lt;")
	}
}

func TestTemplatesAreValid(t *testing.T) {
	// Verify all templates parse correctly
	funcMap := template.FuncMap{
		"ProjectSlug": ProjectSlug,
		"RenderText":  RenderText,
	}

	templates := map[string]string{
		"index":   indexTemplate,
		"project": projectIndexTemplate,
		"session": sessionTemplate,
	}

	for name, tmpl := range templates {
		t.Run(name, func(t *testing.T) {
			_, err := template.New(name).Funcs(funcMap).Parse(tmpl)
			if err != nil {
				t.Errorf("Template %q failed to parse: %v", name, err)
			}
		})
	}
}

func TestCSSEmbedded(t *testing.T) {
	// Verify CSS is embedded in templates
	checks := []struct {
		name     string
		template string
	}{
		{"index", indexTemplate},
		{"project", projectIndexTemplate},
		{"session", sessionTemplate},
	}

	cssElements := []string{
		"--bg-cream",
		"--accent-orange",
		"--sidebar-width",
		"--font-serif",
		"288px",
	}

	for _, tc := range checks {
		t.Run(tc.name, func(t *testing.T) {
			for _, css := range cssElements {
				if !strings.Contains(tc.template, css) {
					t.Errorf("Template %q should contain CSS element %q", tc.name, css)
				}
			}
		})
	}
}

func TestFooterBranding(t *testing.T) {
	// Verify footer branding in all templates
	templates := []string{indexTemplate, projectIndexTemplate, sessionTemplate}

	for i, tmpl := range templates {
		if !strings.Contains(tmpl, "fabriqa.ai") {
			t.Errorf("Template %d should contain fabriqa.ai branding", i)
		}
		if !strings.Contains(tmpl, "https://fabriqa.ai") {
			t.Errorf("Template %d should contain link to https://fabriqa.ai", i)
		}
	}
}

func TestStaticModeBanner(t *testing.T) {
	// Verify static mode banner in all templates
	templates := []string{indexTemplate, projectIndexTemplate, sessionTemplate}

	for i, tmpl := range templates {
		if !strings.Contains(tmpl, "Static Mode") {
			t.Errorf("Template %d should contain Static Mode banner", i)
		}
		if !strings.Contains(tmpl, "claude-code-logs serve") {
			t.Errorf("Template %d should mention serve command", i)
		}
	}
}
