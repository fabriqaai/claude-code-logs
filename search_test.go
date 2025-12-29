package main

import (
	"strings"
	"testing"
	"time"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple words",
			input:    "hello world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "mixed case",
			input:    "Hello World",
			expected: []string{"hello", "world"},
		},
		{
			name:     "with punctuation",
			input:    "hello, world! how are you?",
			expected: []string{"hello", "world", "how", "are", "you"},
		},
		{
			name:     "with numbers",
			input:    "version 2.0 released",
			expected: []string{"version", "released"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single char words filtered",
			input:    "a b c hello",
			expected: []string{"hello"},
		},
		{
			name:     "duplicate words",
			input:    "hello hello world world",
			expected: []string{"hello", "world"},
		},
		{
			name:     "code snippet",
			input:    "func main() { fmt.Println(\"hello\") }",
			expected: []string{"func", "main", "fmt", "println", "hello"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tokenize(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("tokenize(%q) = %v, want %v", tt.input, result, tt.expected)
				return
			}
			for i, term := range result {
				if term != tt.expected[i] {
					t.Errorf("tokenize(%q)[%d] = %q, want %q", tt.input, i, term, tt.expected[i])
				}
			}
		})
	}
}

func TestNewSearchIndex(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test Session",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Hello Claude, how are you?"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "I'm doing well, thank you for asking!"},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	if idx.MessageCount() != 2 {
		t.Errorf("MessageCount() = %d, want 2", idx.MessageCount())
	}

	if idx.TermCount() == 0 {
		t.Error("TermCount() should be > 0")
	}
}

func TestSearchIndex_Search_BasicQuery(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test Session",
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
								{Type: "text", Text: "Hello there!"},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for "hello"
	results := idx.Search("hello", "", "")

	if len(results) != 1 {
		t.Fatalf("Search(\"hello\") returned %d results, want 1", len(results))
	}

	if len(results[0].Matches) != 2 {
		t.Errorf("Search(\"hello\") returned %d matches, want 2", len(results[0].Matches))
	}
}

func TestSearchIndex_Search_MultiTermQuery(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test Session",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Hello Claude, help me with Go programming"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "Hello! I'd be happy to help with Go."},
							},
						},
						{
							UUID:      "msg-3",
							Role:      "user",
							Timestamp: now.Add(2 * time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "Just programming in general"},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for "hello go" - should match messages with BOTH terms
	results := idx.Search("hello go", "", "")

	if len(results) != 1 {
		t.Fatalf("Search returned %d session results, want 1", len(results))
	}

	// Only msg-1 and msg-2 have both "hello" and "go"
	if len(results[0].Matches) != 2 {
		t.Errorf("Search(\"hello go\") returned %d matches, want 2", len(results[0].Matches))
	}
}

func TestSearchIndex_Search_ProjectFilter(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Project 1 Session",
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello from project 1"}}},
					},
				},
			},
		},
		{
			Path:       "/Users/test/project2",
			FolderName: "-Users-test-project2",
			Sessions: []Session{
				{
					ID:      "session-2",
					Summary: "Project 2 Session",
					Messages: []Message{
						{UUID: "msg-2", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello from project 2"}}},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search without filter - should find both
	results := idx.Search("hello", "", "")
	if len(results) != 2 {
		t.Errorf("Search without filter returned %d results, want 2", len(results))
	}

	// Search with project filter
	results = idx.Search("hello", "/Users/test/project1", "")
	if len(results) != 1 {
		t.Errorf("Search with project filter returned %d results, want 1", len(results))
	}

	if results[0].Project != "/Users/test/project1" {
		t.Errorf("Filtered result project = %q, want /Users/test/project1", results[0].Project)
	}
}

func TestSearchIndex_Search_SessionFilter(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Session One",
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello session one"}}},
					},
				},
				{
					ID:      "session-2",
					Summary: "Session Two",
					Messages: []Message{
						{UUID: "msg-2", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello session two"}}},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search with session filter
	results := idx.Search("hello", "", "session-1")
	if len(results) != 1 {
		t.Errorf("Search with session filter returned %d results, want 1", len(results))
	}

	if results[0].SessionID != "session-1" {
		t.Errorf("Filtered result sessionId = %q, want session-1", results[0].SessionID)
	}
}

func TestSearchIndex_Search_EmptyQuery(t *testing.T) {
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test",
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: time.Now(),
							Content: []ContentBlock{{Type: "text", Text: "Hello"}}},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	results := idx.Search("", "", "")
	if len(results) != 0 {
		t.Errorf("Empty query should return 0 results, got %d", len(results))
	}
}

func TestSearchIndex_Search_NoResults(t *testing.T) {
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test",
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: time.Now(),
							Content: []ContentBlock{{Type: "text", Text: "Hello world"}}},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	results := idx.Search("nonexistent", "", "")
	if len(results) != 0 {
		t.Errorf("Search for nonexistent term should return 0 results, got %d", len(results))
	}
}

func TestHighlightMatches(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		terms    []string
		contains []string
	}{
		{
			name:     "single term",
			content:  "Hello world",
			terms:    []string{"hello"},
			contains: []string{"<mark>Hello</mark>"},
		},
		{
			name:     "multiple terms",
			content:  "Hello world, hello again",
			terms:    []string{"hello", "world"},
			contains: []string{"<mark>Hello</mark>", "<mark>world</mark>"},
		},
		{
			name:     "case insensitive",
			content:  "HELLO World",
			terms:    []string{"hello"},
			contains: []string{"<mark>HELLO</mark>"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := highlightMatches(tt.content, tt.terms)
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("highlightMatches(%q, %v) = %q, should contain %q",
						tt.content, tt.terms, result, expected)
				}
			}
		})
	}
}

func TestHighlightMatches_LongContent(t *testing.T) {
	// Create long content
	longContent := strings.Repeat("Lorem ipsum dolor sit amet. ", 100) + "TARGET WORD here " + strings.Repeat("More text. ", 100)

	result := highlightMatches(longContent, []string{"target"})

	// Should contain highlighted match
	if !strings.Contains(result, "<mark>TARGET</mark>") {
		t.Error("Long content should contain highlighted match")
	}

	// Should be truncated (not full length)
	if len(result) > 600 {
		t.Errorf("Long content should be truncated, got length %d", len(result))
	}

	// Should contain ellipsis
	if !strings.Contains(result, "...") {
		t.Error("Truncated content should contain ellipsis")
	}
}

func TestExtractTextContent(t *testing.T) {
	msg := Message{
		Content: []ContentBlock{
			{Type: "text", Text: "Hello"},
			{Type: "tool_use", ToolName: "Read", ToolInput: "{}"},
			{Type: "text", Text: "World"},
			{Type: "tool_result", ToolOutput: "result"},
		},
	}

	result := extractTextContent(msg)

	if result != "Hello World" {
		t.Errorf("extractTextContent() = %q, want \"Hello World\"", result)
	}
}

func TestSearchIndex_SkipsToolBlocks(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Test",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "assistant",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Let me read that file"},
								{Type: "tool_use", ToolName: "Read", ToolInput: `{"path": "/secret/file"}`},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for text content should work
	results := idx.Search("read", "", "")
	if len(results) != 1 {
		t.Error("Should find 'read' in text content")
	}

	// Search for tool input should NOT work (tool blocks not indexed)
	results = idx.Search("secret", "", "")
	if len(results) != 0 {
		t.Error("Should NOT find 'secret' from tool input (tool blocks not indexed)")
	}
}
