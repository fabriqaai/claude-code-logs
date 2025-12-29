package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestDecodeProjectPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard path",
			input:    "-Users-name-project",
			expected: "/Users/name/project",
		},
		{
			name:     "deep nested path",
			input:    "-Users-cengiz-han-workspace-code-claude-code-logs",
			expected: "/Users/cengiz/han/workspace/code/claude/code/logs",
		},
		{
			name:     "root marker",
			input:    "-",
			expected: "/",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "no leading dash",
			input:    "some-folder",
			expected: "some-folder",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DecodeProjectPath(tt.input)
			if result != tt.expected {
				t.Errorf("DecodeProjectPath(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestDiscoverProjects(t *testing.T) {
	// Create temporary directory structure
	tmpDir := t.TempDir()

	// Create test project directories
	testProjects := []string{
		"-Users-test-project1",
		"-Users-test-project2",
		".hidden-folder", // Should be skipped
	}

	for _, p := range testProjects {
		if err := os.MkdirAll(filepath.Join(tmpDir, p), 0755); err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}
	}

	// Create a non-directory file (should be skipped)
	if err := os.WriteFile(filepath.Join(tmpDir, "some-file.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test discovery
	projects, err := DiscoverProjects(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverProjects failed: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(projects))
	}

	// Check project paths are decoded
	for _, p := range projects {
		if p.Path == "" {
			t.Errorf("Project path should not be empty")
		}
		if p.FolderName == "" {
			t.Errorf("Project folder name should not be empty")
		}
	}
}

func TestDiscoverProjects_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	projects, err := DiscoverProjects(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverProjects failed on empty dir: %v", err)
	}

	if len(projects) != 0 {
		t.Errorf("Expected 0 projects for empty directory, got %d", len(projects))
	}
}

func TestDiscoverProjects_NonExistentDirectory(t *testing.T) {
	_, err := DiscoverProjects("/non/existent/path")
	if err == nil {
		t.Error("Expected error for non-existent directory, got nil")
	}
}

func TestParseSession(t *testing.T) {
	// Create temporary JSONL file
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	jsonlContent := `{"type":"summary","summary":"Test Session","leafUuid":"abc123"}
{"type":"user","uuid":"msg1","parentUuid":null,"timestamp":"2025-12-29T10:00:00.000Z","message":{"role":"user","content":"Hello, Claude!"}}
{"type":"assistant","uuid":"msg2","parentUuid":"msg1","timestamp":"2025-12-29T10:00:05.000Z","message":{"role":"assistant","content":[{"type":"text","text":"Hello! How can I help you today?"}]}}
{"type":"file-history-snapshot","messageId":"msg1","snapshot":{}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed: %v", err)
	}

	// Check summary
	if session.Summary != "Test Session" {
		t.Errorf("Expected summary 'Test Session', got %q", session.Summary)
	}

	// Check messages
	if len(session.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(session.Messages))
	}

	// Check first message (user)
	if session.Messages[0].Role != "user" {
		t.Errorf("Expected first message role 'user', got %q", session.Messages[0].Role)
	}

	// Check second message (assistant)
	if session.Messages[1].Role != "assistant" {
		t.Errorf("Expected second message role 'assistant', got %q", session.Messages[1].Role)
	}

	// Check threading
	if session.Messages[1].ParentUUID != "msg1" {
		t.Errorf("Expected parent UUID 'msg1', got %q", session.Messages[1].ParentUUID)
	}
}

func TestParseSession_WithToolCalls(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	jsonlContent := `{"type":"summary","summary":"Tool Test"}
{"type":"assistant","uuid":"msg1","parentUuid":null,"timestamp":"2025-12-29T10:00:00.000Z","message":{"role":"assistant","content":[{"type":"tool_use","id":"tool1","name":"Read","input":{"file_path":"/test/file.txt"}}]}}
{"type":"assistant","uuid":"msg2","parentUuid":"msg1","timestamp":"2025-12-29T10:00:01.000Z","message":{"role":"assistant","content":[{"type":"tool_result","tool_use_id":"tool1","content":"file contents here"}]}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed: %v", err)
	}

	if len(session.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(session.Messages))
	}

	// Check tool_use block
	toolUseBlock := session.Messages[0].Content[0]
	if toolUseBlock.Type != "tool_use" {
		t.Errorf("Expected type 'tool_use', got %q", toolUseBlock.Type)
	}
	if toolUseBlock.ToolName != "Read" {
		t.Errorf("Expected tool name 'Read', got %q", toolUseBlock.ToolName)
	}

	// Check tool_result block
	toolResultBlock := session.Messages[1].Content[0]
	if toolResultBlock.Type != "tool_result" {
		t.Errorf("Expected type 'tool_result', got %q", toolResultBlock.Type)
	}
	if toolResultBlock.ToolUseID != "tool1" {
		t.Errorf("Expected tool_use_id 'tool1', got %q", toolResultBlock.ToolUseID)
	}
}

func TestParseSession_MalformedJSON(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	jsonlContent := `{"type":"summary","summary":"Test"}
{this is not valid json}
{"type":"user","uuid":"msg1","timestamp":"2025-12-29T10:00:00.000Z","message":{"role":"user","content":"Hello"}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession should not fail on malformed lines: %v", err)
	}

	// Should have summary and one message (malformed line skipped)
	if session.Summary != "Test" {
		t.Errorf("Expected summary 'Test', got %q", session.Summary)
	}
	if len(session.Messages) != 1 {
		t.Errorf("Expected 1 message (malformed line skipped), got %d", len(session.Messages))
	}
}

func TestParseSession_EmptyFile(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	if err := os.WriteFile(sessionPath, []byte(""), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed on empty file: %v", err)
	}

	if len(session.Messages) != 0 {
		t.Errorf("Expected 0 messages for empty file, got %d", len(session.Messages))
	}
}

func TestListSessions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create project directory
	projectDir := filepath.Join(tmpDir, "-Users-test-project")
	if err := os.MkdirAll(projectDir, 0755); err != nil {
		t.Fatalf("Failed to create project directory: %v", err)
	}

	// Create test sessions with different timestamps
	sessions := []struct {
		id        string
		timestamp time.Time
	}{
		{"session1", time.Date(2025, 12, 29, 10, 0, 0, 0, time.UTC)},
		{"session2", time.Date(2025, 12, 29, 11, 0, 0, 0, time.UTC)},
		{"session3", time.Date(2025, 12, 29, 9, 0, 0, 0, time.UTC)},
	}

	for _, s := range sessions {
		content := `{"type":"summary","summary":"Session ` + s.id + `"}
{"type":"user","uuid":"msg1","timestamp":"` + s.timestamp.Format(time.RFC3339) + `","message":{"role":"user","content":"test"}}
`
		path := filepath.Join(projectDir, s.id+".jsonl")
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to create session file: %v", err)
		}
	}

	// Also create a non-jsonl file (should be skipped)
	if err := os.WriteFile(filepath.Join(projectDir, "readme.txt"), []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create non-jsonl file: %v", err)
	}

	project := &Project{
		FolderName: "-Users-test-project",
		Path:       "/Users/test/project",
	}

	sessionList, err := ListSessions(tmpDir, project)
	if err != nil {
		t.Fatalf("ListSessions failed: %v", err)
	}

	if len(sessionList) != 3 {
		t.Errorf("Expected 3 sessions, got %d", len(sessionList))
	}

	// Check sessions are sorted by creation date (newest first)
	if len(sessionList) >= 2 {
		if !sessionList[0].CreatedAt.After(sessionList[1].CreatedAt) {
			t.Errorf("Sessions should be sorted newest first")
		}
	}
}

func TestDefaultClaudeProjectsPath(t *testing.T) {
	path, err := DefaultClaudeProjectsPath()
	if err != nil {
		t.Fatalf("DefaultClaudeProjectsPath failed: %v", err)
	}

	if path == "" {
		t.Error("Expected non-empty path")
	}

	// Should end with .claude/projects
	if !filepath.IsAbs(path) {
		t.Errorf("Expected absolute path, got %q", path)
	}
}
