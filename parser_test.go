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

func TestGetProjectLastUpdate(t *testing.T) {
	t1 := time.Date(2025, 12, 25, 10, 0, 0, 0, time.UTC)
	t2 := time.Date(2025, 12, 26, 10, 0, 0, 0, time.UTC)
	t3 := time.Date(2025, 12, 27, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name     string
		project  Project
		expected time.Time
	}{
		{
			name: "single session",
			project: Project{
				Path: "/test/project1",
				Sessions: []Session{
					{ID: "s1", UpdatedAt: t1},
				},
			},
			expected: t1,
		},
		{
			name: "multiple sessions - returns most recent",
			project: Project{
				Path: "/test/project2",
				Sessions: []Session{
					{ID: "s1", UpdatedAt: t1},
					{ID: "s2", UpdatedAt: t3},
					{ID: "s3", UpdatedAt: t2},
				},
			},
			expected: t3,
		},
		{
			name: "no sessions - returns zero time",
			project: Project{
				Path:     "/test/project3",
				Sessions: []Session{},
			},
			expected: time.Time{},
		},
		{
			name: "nil sessions - returns zero time",
			project: Project{
				Path:     "/test/project4",
				Sessions: nil,
			},
			expected: time.Time{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getProjectLastUpdate(&tt.project)
			if !result.Equal(tt.expected) {
				t.Errorf("getProjectLastUpdate() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestParseMessage_EdgeCases(t *testing.T) {
	// Test with empty message content
	entry := &jsonlEntry{
		UUID:      "msg1",
		Timestamp: "2025-12-29T10:00:00.000Z",
	}

	msg, err := parseMessage(entry)
	if err != nil {
		t.Fatalf("parseMessage failed on empty message: %v", err)
	}
	if msg.UUID != "msg1" {
		t.Errorf("Expected UUID 'msg1', got %q", msg.UUID)
	}
	if len(msg.Content) != 0 {
		t.Errorf("Expected empty content, got %d blocks", len(msg.Content))
	}
}

func TestParseMessage_NullParentUUID(t *testing.T) {
	// Test with nil parentUUID
	entry := &jsonlEntry{
		UUID:       "msg1",
		ParentUUID: nil,
		Timestamp:  "2025-12-29T10:00:00.000Z",
	}

	msg, err := parseMessage(entry)
	if err != nil {
		t.Fatalf("parseMessage failed: %v", err)
	}
	if msg.ParentUUID != "" {
		t.Errorf("Expected empty ParentUUID for nil, got %q", msg.ParentUUID)
	}
}

func TestParseMessage_WithParentUUID(t *testing.T) {
	parentID := "parent-msg"
	entry := &jsonlEntry{
		UUID:       "msg1",
		ParentUUID: &parentID,
		Timestamp:  "2025-12-29T10:00:00.000Z",
	}

	msg, err := parseMessage(entry)
	if err != nil {
		t.Fatalf("parseMessage failed: %v", err)
	}
	if msg.ParentUUID != "parent-msg" {
		t.Errorf("Expected ParentUUID 'parent-msg', got %q", msg.ParentUUID)
	}
}

func TestParseMessage_AlternateTimestampFormat(t *testing.T) {
	// Test with different timestamp format
	entry := &jsonlEntry{
		UUID:      "msg1",
		Timestamp: "2025-12-29T10:00:00Z", // RFC3339 format
	}

	msg, err := parseMessage(entry)
	if err != nil {
		t.Fatalf("parseMessage failed: %v", err)
	}

	expected := time.Date(2025, 12, 29, 10, 0, 0, 0, time.UTC)
	if !msg.Timestamp.Equal(expected) {
		t.Errorf("Expected timestamp %v, got %v", expected, msg.Timestamp)
	}
}

func TestParseMessage_InvalidTimestamp(t *testing.T) {
	entry := &jsonlEntry{
		UUID:      "msg1",
		Timestamp: "not-a-valid-timestamp",
	}

	msg, err := parseMessage(entry)
	if err != nil {
		t.Fatalf("parseMessage should not fail on invalid timestamp: %v", err)
	}
	// Should fallback to current time (approximately)
	if msg.Timestamp.IsZero() {
		t.Error("Timestamp should not be zero (should fallback to now)")
	}
}

func TestParseSession_WithCWD(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	jsonlContent := `{"type":"user","uuid":"msg1","timestamp":"2025-12-29T10:00:00.000Z","cwd":"/actual/project/path","message":{"role":"user","content":"Hello!"}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed: %v", err)
	}

	if session.CWD != "/actual/project/path" {
		t.Errorf("Expected CWD '/actual/project/path', got %q", session.CWD)
	}
}

func TestParseSession_SummaryFallback(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	// No summary entry - should use first user message as fallback
	jsonlContent := `{"type":"user","uuid":"msg1","timestamp":"2025-12-29T10:00:00.000Z","message":{"role":"user","content":"This is the first user message that will become the summary"}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed: %v", err)
	}

	if session.Summary != "This is the first user message that will become the summary" {
		t.Errorf("Expected summary to be first user message, got %q", session.Summary)
	}
}

func TestParseSession_SummaryFallbackTruncation(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "test-session.jsonl")

	// Create a long message that should be truncated
	longMessage := "This is a very long message that exceeds one hundred characters and should be truncated to approximately one hundred characters with ellipsis"
	jsonlContent := `{"type":"user","uuid":"msg1","timestamp":"2025-12-29T10:00:00.000Z","message":{"role":"user","content":"` + longMessage + `"}}
`

	if err := os.WriteFile(sessionPath, []byte(jsonlContent), 0644); err != nil {
		t.Fatalf("Failed to create test JSONL file: %v", err)
	}

	session, err := ParseSession(sessionPath, "test-session")
	if err != nil {
		t.Fatalf("ParseSession failed: %v", err)
	}

	if len(session.Summary) > 100 {
		t.Errorf("Summary should be truncated to ~100 chars, got %d chars", len(session.Summary))
	}

	if session.Summary[len(session.Summary)-3:] != "..." {
		t.Error("Truncated summary should end with '...'")
	}
}

func TestDiscoverProjects_SkipsNonEncodedPaths(t *testing.T) {
	tmpDir := t.TempDir()

	// Create directories - some with encoded paths, some without
	dirs := []string{
		"-Users-test-project1", // Valid encoded path
		"-Users-test-project2", // Valid encoded path
		"regular-folder",       // Not an encoded path (no leading dash)
		"subfolder",            // Not an encoded path
	}

	for _, d := range dirs {
		if err := os.MkdirAll(filepath.Join(tmpDir, d), 0755); err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}
	}

	projects, err := DiscoverProjects(tmpDir)
	if err != nil {
		t.Fatalf("DiscoverProjects failed: %v", err)
	}

	if len(projects) != 2 {
		t.Errorf("Expected 2 projects (only encoded paths), got %d", len(projects))
	}
}

func TestDiscoverProjects_NotADirectory(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "not-a-dir")

	if err := os.WriteFile(filePath, []byte("test"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := DiscoverProjects(filePath)
	if err == nil {
		t.Error("Expected error when path is not a directory")
	}
}
