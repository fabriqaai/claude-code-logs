package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestMarkdownGeneratorShouldRegenerate(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a JSONL file
	jsonlPath := filepath.Join(tmpDir, "session.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(`{"type":"summary"}`), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	// Create MD file that's older than JSONL
	mdPath := filepath.Join(tmpDir, "session.md")
	if err := os.WriteFile(mdPath, []byte("# Old MD"), 0644); err != nil {
		t.Fatalf("Failed to write MD: %v", err)
	}

	// Set MD file to be older
	oldTime := time.Now().Add(-1 * time.Hour)
	if err := os.Chtimes(mdPath, oldTime, oldTime); err != nil {
		t.Fatalf("Failed to set mtime: %v", err)
	}

	gen, err := NewMarkdownGenerator(tmpDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	// MD is older than JSONL - should regenerate
	if !gen.ShouldRegenerate(jsonlPath, mdPath) {
		t.Error("Expected ShouldRegenerate=true when MD is older than JSONL")
	}

	// Update MD to be newer
	if err := os.WriteFile(mdPath, []byte("# New MD"), 0644); err != nil {
		t.Fatalf("Failed to update MD: %v", err)
	}

	// MD is newer than JSONL - should not regenerate
	if gen.ShouldRegenerate(jsonlPath, mdPath) {
		t.Error("Expected ShouldRegenerate=false when MD is newer than JSONL")
	}

	// Force mode - should always regenerate
	genForce, err := NewMarkdownGenerator(tmpDir, tmpDir, true, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}
	if !genForce.ShouldRegenerate(jsonlPath, mdPath) {
		t.Error("Expected ShouldRegenerate=true in force mode")
	}
}

func TestMarkdownGeneratorShouldRegenerateNonExistent(t *testing.T) {
	tmpDir := t.TempDir()

	jsonlPath := filepath.Join(tmpDir, "session.jsonl")
	if err := os.WriteFile(jsonlPath, []byte(`{"type":"summary"}`), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	mdPath := filepath.Join(tmpDir, "nonexistent.md")

	gen, err := NewMarkdownGenerator(tmpDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	// MD doesn't exist - should regenerate
	if !gen.ShouldRegenerate(jsonlPath, mdPath) {
		t.Error("Expected ShouldRegenerate=true when MD doesn't exist")
	}
}

func TestMarkdownGeneratorGenerateSession(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake JSONL file for hashing
	jsonlPath := filepath.Join(tmpDir, "src", "test-session.jsonl")
	if err := os.MkdirAll(filepath.Dir(jsonlPath), 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.WriteFile(jsonlPath, []byte(`{"type":"summary"}`), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	session := &Session{
		ID:         "abc123",
		Summary:    "Test Session",
		SourcePath: jsonlPath,
		CWD:        "/Users/test/project",
		CreatedAt:  time.Date(2026, 1, 17, 10, 30, 0, 0, time.UTC),
		Messages: []Message{
			{
				Role: "user",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello, assistant!"},
				},
				Timestamp: time.Date(2026, 1, 17, 10, 30, 0, 0, time.UTC),
			},
			{
				Role: "assistant",
				Content: []ContentBlock{
					{Type: "text", Text: "Hello! How can I help?"},
				},
				Timestamp: time.Date(2026, 1, 17, 10, 31, 0, 0, time.UTC),
			},
		},
	}

	outDir := filepath.Join(tmpDir, "output")

	// Create project directory
	projectSlug := "test-project"
	projectDir := filepath.Join(outDir, projectSlug)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	project := &Project{
		Path:       "/Users/test/project",
		FolderName: "test-project",
		Sessions:   []Session{*session},
	}

	gen, err := NewMarkdownGenerator(outDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	if err := gen.GenerateSession(session, project, projectSlug); err != nil {
		t.Fatalf("GenerateSession failed: %v", err)
	}

	// Read generated file
	mdPath := filepath.Join(projectDir, "abc123.md")
	content, err := os.ReadFile(mdPath)
	if err != nil {
		t.Fatalf("Failed to read generated MD: %v", err)
	}

	mdContent := string(content)

	// Check frontmatter exists
	if !strings.HasPrefix(mdContent, "---\n") {
		t.Error("Expected MD to start with frontmatter")
	}

	// Check content
	expectedParts := []string{
		"source: test-session.jsonl",
		"project: /Users/test/project",
		"title: Test Session",
		"## User",
		"Hello, assistant!",
		"## Assistant",
		"Hello! How can I help?",
	}

	for _, part := range expectedParts {
		if !strings.Contains(mdContent, part) {
			t.Errorf("Expected MD to contain %q", part)
		}
	}
}

func TestMarkdownGeneratorGenerateMainIndex(t *testing.T) {
	tmpDir := t.TempDir()

	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{ID: "s1", CreatedAt: time.Date(2026, 1, 17, 10, 0, 0, 0, time.UTC)},
				{ID: "s2", CreatedAt: time.Date(2026, 1, 16, 10, 0, 0, 0, time.UTC)},
			},
		},
		{
			Path:       "/Users/test/project2",
			FolderName: "-Users-test-project2",
			Sessions: []Session{
				{ID: "s3", CreatedAt: time.Date(2026, 1, 15, 10, 0, 0, 0, time.UTC)},
			},
		},
	}

	gen, err := NewMarkdownGenerator(tmpDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	if err := gen.GenerateMainIndex(projects); err != nil {
		t.Fatalf("GenerateMainIndex failed: %v", err)
	}

	// Read generated file
	content, err := os.ReadFile(filepath.Join(tmpDir, "index.md"))
	if err != nil {
		t.Fatalf("Failed to read index.md: %v", err)
	}

	mdContent := string(content)

	// Check structure
	expectedParts := []string{
		"# Claude Code Logs",
		"| Project | Sessions | Last Activity |",
		"/Users/test/project1",
		"/Users/test/project2",
		"| 2 |", // project1 has 2 sessions
		"| 1 |", // project2 has 1 session
	}

	for _, part := range expectedParts {
		if !strings.Contains(mdContent, part) {
			t.Errorf("Expected index to contain %q, got:\n%s", part, mdContent)
		}
	}
}

func TestMarkdownGeneratorGenerateProjectIndex(t *testing.T) {
	tmpDir := t.TempDir()

	project := &Project{
		Path:       "/Users/test/myproject",
		FolderName: "-Users-test-myproject",
		Sessions: []Session{
			{
				ID:        "session1",
				Summary:   "First Session",
				CreatedAt: time.Date(2026, 1, 17, 10, 0, 0, 0, time.UTC),
			},
			{
				ID:        "session2",
				Summary:   "Second Session",
				CreatedAt: time.Date(2026, 1, 16, 10, 0, 0, 0, time.UTC),
			},
		},
	}

	projectSlug := "users-test-myproject"
	projectDir := filepath.Join(tmpDir, projectSlug)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	gen, err := NewMarkdownGenerator(tmpDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	if err := gen.GenerateProjectIndex(project, projectSlug); err != nil {
		t.Fatalf("GenerateProjectIndex failed: %v", err)
	}

	// Read generated file
	content, err := os.ReadFile(filepath.Join(projectDir, "index.md"))
	if err != nil {
		t.Fatalf("Failed to read project index.md: %v", err)
	}

	mdContent := string(content)

	// Check structure
	expectedParts := []string{
		"# Project: /Users/test/myproject",
		"| Session | Title | Created |",
		"[session1](session1.md)",
		"[session2](session2.md)",
		"First Session",
		"Second Session",
	}

	for _, part := range expectedParts {
		if !strings.Contains(mdContent, part) {
			t.Errorf("Expected project index to contain %q", part)
		}
	}
}

func TestEscapeMarkdownTableCell(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"simple text", "simple text"},
		{"text|with|pipes", "text\\|with\\|pipes"},
		{"text\nwith\nnewlines", "text with newlines"},
		{"pipes|and\nnewlines", "pipes\\|and newlines"},
	}

	for _, tt := range tests {
		result := escapeMarkdownTableCell(tt.input)
		if result != tt.expected {
			t.Errorf("escapeMarkdownTableCell(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"user", "User"},
		{"assistant", "Assistant"},
		{"", ""},
		{"A", "A"},
		{"already Capitalized", "Already Capitalized"},
	}

	for _, tt := range tests {
		result := capitalizeFirst(tt.input)
		if result != tt.expected {
			t.Errorf("capitalizeFirst(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToolCallFormatting(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a fake JSONL file for hashing
	jsonlPath := filepath.Join(tmpDir, "src", "test-session.jsonl")
	if err := os.MkdirAll(filepath.Dir(jsonlPath), 0755); err != nil {
		t.Fatalf("Failed to create src dir: %v", err)
	}
	if err := os.WriteFile(jsonlPath, []byte(`{"type":"summary"}`), 0644); err != nil {
		t.Fatalf("Failed to write JSONL: %v", err)
	}

	session := &Session{
		ID:         "tooltest",
		Summary:    "Tool Test",
		SourcePath: jsonlPath,
		CWD:        "/Users/test/project",
		CreatedAt:  time.Date(2026, 1, 17, 10, 30, 0, 0, time.UTC),
		Messages: []Message{
			{
				Role: "assistant",
				Content: []ContentBlock{
					{
						Type:      "tool_use",
						ToolName:  "Read",
						ToolInput: `{"file_path":"/test/file.txt"}`,
					},
				},
			},
			{
				Role: "user",
				Content: []ContentBlock{
					{
						Type:       "tool_result",
						ToolUseID:  "tool123",
						ToolOutput: "File contents here",
					},
				},
			},
		},
	}

	outDir := filepath.Join(tmpDir, "output")

	projectSlug := "test-project"
	projectDir := filepath.Join(outDir, projectSlug)
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project dir: %v", err)
	}

	project := &Project{
		Path:       "/Users/test/project",
		FolderName: "test-project",
		Sessions:   []Session{*session},
	}

	gen, err := NewMarkdownGenerator(outDir, tmpDir, false, false, nil)
	if err != nil {
		t.Fatalf("NewMarkdownGenerator failed: %v", err)
	}

	if err := gen.GenerateSession(session, project, projectSlug); err != nil {
		t.Fatalf("GenerateSession failed: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(projectDir, "tooltest.md"))
	if err != nil {
		t.Fatalf("Failed to read generated MD: %v", err)
	}

	mdContent := string(content)

	// Check tool call formatting
	expectedParts := []string{
		"<details>",
		"<summary>Tool: Read</summary>",
		"```json",
		`{"file_path":"/test/file.txt"}`,
		"</details>",
		"<summary>Result</summary>",
		"File contents here",
	}

	for _, part := range expectedParts {
		if !strings.Contains(mdContent, part) {
			t.Errorf("Expected MD to contain %q, got:\n%s", part, mdContent)
		}
	}
}
