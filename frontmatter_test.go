package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFrontmatterMarshal(t *testing.T) {
	fm := Frontmatter{
		Source:     "test-session.jsonl",
		SourceHash: "abc123def456",
		Project:    "/Users/test/project",
		Title:      "Test Session",
		Created:    "2026-01-17T10:30:00Z",
	}

	data, err := fm.Marshal()
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	output := string(data)

	// Check frontmatter delimiters
	if output[:4] != "---\n" {
		t.Errorf("Expected frontmatter to start with '---\\n', got: %q", output[:4])
	}
	if output[len(output)-4:] != "---\n" {
		t.Errorf("Expected frontmatter to end with '---\\n', got: %q", output[len(output)-4:])
	}

	// Check content contains expected fields
	expectedFields := []string{
		"source: test-session.jsonl",
		"source_hash: abc123def456",
		"project: /Users/test/project",
		"title: Test Session",
		"created: \"2026-01-17T10:30:00Z\"",
	}

	for _, field := range expectedFields {
		if !contains(output, field) {
			t.Errorf("Expected frontmatter to contain %q, got:\n%s", field, output)
		}
	}
}

func TestParseFrontmatter(t *testing.T) {
	content := `---
source: test-session.jsonl
source_hash: abc123def456
project: /Users/test/project
title: Test Session
created: "2026-01-17T10:30:00Z"
---

## User

Hello world
`

	fm, remaining, err := ParseFrontmatter([]byte(content))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	if fm.Source != "test-session.jsonl" {
		t.Errorf("Expected source 'test-session.jsonl', got %q", fm.Source)
	}
	if fm.SourceHash != "abc123def456" {
		t.Errorf("Expected source_hash 'abc123def456', got %q", fm.SourceHash)
	}
	if fm.Project != "/Users/test/project" {
		t.Errorf("Expected project '/Users/test/project', got %q", fm.Project)
	}
	if fm.Title != "Test Session" {
		t.Errorf("Expected title 'Test Session', got %q", fm.Title)
	}

	// Check remaining content
	expectedRemaining := "\n## User\n\nHello world\n"
	if string(remaining) != expectedRemaining {
		t.Errorf("Expected remaining content %q, got %q", expectedRemaining, string(remaining))
	}
}

func TestParseFrontmatterNoFrontmatter(t *testing.T) {
	content := `# Just Markdown

No frontmatter here.
`

	fm, remaining, err := ParseFrontmatter([]byte(content))
	if err != nil {
		t.Fatalf("ParseFrontmatter failed: %v", err)
	}

	// Should return empty frontmatter
	if fm.Source != "" {
		t.Errorf("Expected empty source, got %q", fm.Source)
	}

	// Should return original content
	if string(remaining) != content {
		t.Errorf("Expected remaining to equal original content")
	}
}

func TestComputeFileHash(t *testing.T) {
	// Create a temp file with known content
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "test.txt")

	content := "Hello, World!"
	if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write temp file: %v", err)
	}

	hash, err := ComputeFileHash(tmpFile)
	if err != nil {
		t.Fatalf("ComputeFileHash failed: %v", err)
	}

	// SHA256 of "Hello, World!" is known
	expectedHash := "dffd6021bb2bd5b0af676290809ec3a53191dd81c7f70a4b28688a362182986f"
	if hash != expectedHash {
		t.Errorf("Expected hash %q, got %q", expectedHash, hash)
	}
}

func TestNewFrontmatter(t *testing.T) {
	session := &Session{
		ID:         "abc123",
		Summary:    "Test Session Summary",
		SourcePath: "/path/to/session-abc123.jsonl",
		CWD:        "/Users/test/myproject",
		CreatedAt:  time.Date(2026, 1, 17, 10, 30, 0, 0, time.UTC),
	}

	fm := NewFrontmatter(session, "hashvalue123")

	if fm.Source != "session-abc123.jsonl" {
		t.Errorf("Expected source 'session-abc123.jsonl', got %q", fm.Source)
	}
	if fm.SourceHash != "hashvalue123" {
		t.Errorf("Expected source_hash 'hashvalue123', got %q", fm.SourceHash)
	}
	if fm.Project != "/Users/test/myproject" {
		t.Errorf("Expected project '/Users/test/myproject', got %q", fm.Project)
	}
	if fm.Title != "Test Session Summary" {
		t.Errorf("Expected title 'Test Session Summary', got %q", fm.Title)
	}
	if fm.Created != "2026-01-17T10:30:00Z" {
		t.Errorf("Expected created '2026-01-17T10:30:00Z', got %q", fm.Created)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
