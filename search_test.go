package main

import (
	"fmt"
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

func TestParseQuery(t *testing.T) {
	tests := []struct {
		name            string
		query           string
		expectedTerms   []string
		expectedPhrases []string
	}{
		{
			name:            "simple words",
			query:           "hello world",
			expectedTerms:   []string{"hello", "world"},
			expectedPhrases: nil,
		},
		{
			name:            "single quoted phrase",
			query:           `"hello world"`,
			expectedTerms:   []string{"hello", "world"},
			expectedPhrases: []string{"hello world"},
		},
		{
			name:            "phrase with extra terms",
			query:           `"hello world" foo bar`,
			expectedTerms:   []string{"foo", "bar", "hello", "world"},
			expectedPhrases: []string{"hello world"},
		},
		{
			name:            "multiple phrases",
			query:           `"hello world" "foo bar"`,
			expectedTerms:   []string{"hello", "world", "foo", "bar"},
			expectedPhrases: []string{"hello world", "foo bar"},
		},
		{
			name:            "empty quotes",
			query:           `"" hello`,
			expectedTerms:   []string{"hello"},
			expectedPhrases: nil,
		},
		{
			name:            "phrase with special chars",
			query:           `"func main()"`,
			expectedTerms:   []string{"func", "main"},
			expectedPhrases: []string{"func main()"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseQuery(tt.query)

			// Check phrases
			if len(result.Phrases) != len(tt.expectedPhrases) {
				t.Errorf("parseQuery(%q).Phrases = %v, want %v", tt.query, result.Phrases, tt.expectedPhrases)
			} else {
				for i, phrase := range result.Phrases {
					if phrase != tt.expectedPhrases[i] {
						t.Errorf("parseQuery(%q).Phrases[%d] = %q, want %q", tt.query, i, phrase, tt.expectedPhrases[i])
					}
				}
			}

			// Check that all expected terms are present (order may vary)
			termSet := make(map[string]bool)
			for _, term := range result.Terms {
				termSet[term] = true
			}
			for _, expected := range tt.expectedTerms {
				if !termSet[expected] {
					t.Errorf("parseQuery(%q).Terms missing %q, got %v", tt.query, expected, result.Terms)
				}
			}
		})
	}
}

func TestSearchIndex_Search_PhraseQuery(t *testing.T) {
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
								{Type: "text", Text: "Hello world, how are you today?"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "Hello there! The world is great."},
							},
						},
						{
							UUID:      "msg-3",
							Role:      "user",
							Timestamp: now.Add(2 * time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "world hello reversed"},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for exact phrase "hello world" - should only match msg-1
	results := idx.Search(`"hello world"`, "", "")

	if len(results) != 1 {
		t.Fatalf("Search for phrase returned %d session results, want 1", len(results))
	}

	// Only msg-1 has "hello world" as an exact phrase
	if len(results[0].Matches) != 1 {
		t.Errorf("Search for phrase returned %d matches, want 1 (only msg-1 has exact phrase)", len(results[0].Matches))
	}

	if results[0].Matches[0].MessageID != "msg-1" {
		t.Errorf("Expected match on msg-1, got %s", results[0].Matches[0].MessageID)
	}
}

func TestSearchIndex_Search_PhraseWithTerms(t *testing.T) {
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
								{Type: "text", Text: "Hello world, let me help you with Go programming"},
							},
						},
						{
							UUID:      "msg-2",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content: []ContentBlock{
								{Type: "text", Text: "Hello world! I love Python programming."},
							},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for phrase "hello world" AND term "go"
	results := idx.Search(`"hello world" go`, "", "")

	if len(results) != 1 {
		t.Fatalf("Search returned %d session results, want 1", len(results))
	}

	// Only msg-1 has both "hello world" phrase AND "go" term
	if len(results[0].Matches) != 1 {
		t.Errorf("Search returned %d matches, want 1", len(results[0].Matches))
	}

	if results[0].Matches[0].MessageID != "msg-1" {
		t.Errorf("Expected match on msg-1, got %s", results[0].Matches[0].MessageID)
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

func TestSearchWithOptions_Pagination(t *testing.T) {
	now := time.Now()
	// Create multiple sessions to test pagination
	sessions := make([]Session, 25)
	for i := 0; i < 25; i++ {
		sessions[i] = Session{
			ID:      fmt.Sprintf("session-%d", i),
			Summary: fmt.Sprintf("Session %d", i),
			Messages: []Message{
				{
					UUID:      fmt.Sprintf("msg-%d", i),
					Role:      "user",
					Timestamp: now.Add(time.Duration(i) * time.Hour),
					Content:   []ContentBlock{{Type: "text", Text: "Hello world test message"}},
				},
			},
		}
	}

	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions:   sessions,
		},
	}

	idx := NewSearchIndex(projects)

	// Test default pagination (limit 20)
	result := idx.SearchWithOptions("hello", "", "", SearchOptions{})
	if len(result.Results) != 20 {
		t.Errorf("Default limit should return 20 results, got %d", len(result.Results))
	}
	if result.Total != 25 {
		t.Errorf("Total should be 25, got %d", result.Total)
	}
	if !result.HasMore {
		t.Error("HasMore should be true when more results exist")
	}

	// Test with offset
	result = idx.SearchWithOptions("hello", "", "", SearchOptions{Offset: 20})
	if len(result.Results) != 5 {
		t.Errorf("Offset 20 should return 5 results, got %d", len(result.Results))
	}
	if result.HasMore {
		t.Error("HasMore should be false when no more results")
	}
	if result.Offset != 20 {
		t.Errorf("Offset should be 20, got %d", result.Offset)
	}

	// Test with custom limit
	result = idx.SearchWithOptions("hello", "", "", SearchOptions{Limit: 10})
	if len(result.Results) != 10 {
		t.Errorf("Limit 10 should return 10 results, got %d", len(result.Results))
	}

	// Test max limit cap at 100
	result = idx.SearchWithOptions("hello", "", "", SearchOptions{Limit: 200})
	if len(result.Results) != 25 { // only 25 total results
		t.Errorf("Max limit should cap at 100, but we only have 25 results, got %d", len(result.Results))
	}
}

func TestSearchWithOptions_PaginationDefaults(t *testing.T) {
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
							Content:   []ContentBlock{{Type: "text", Text: "Hello world"}},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Test that defaults are applied (offset=0, limit=20)
	result := idx.SearchWithOptions("hello", "", "", SearchOptions{})
	if result.Offset != 0 {
		t.Errorf("Default offset should be 0, got %d", result.Offset)
	}
	// With only 1 result, hasMore should be false
	if result.HasMore {
		t.Error("HasMore should be false when fewer results than limit")
	}
}

func TestSearchWithOptions_RelevanceScoring(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Low relevance",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now.Add(-30 * 24 * time.Hour), // 30 days ago
							Content:   []ContentBlock{{Type: "text", Text: "hello at the end of a long message that has lots of other content before it finally says hello"}},
						},
					},
				},
				{
					ID:      "session-2",
					Summary: "High relevance",
					Messages: []Message{
						{
							UUID:      "msg-2",
							Role:      "user",
							Timestamp: now, // Recent
							Content:   []ContentBlock{{Type: "text", Text: "hello hello hello world"}}, // Multiple occurrences, early position
						},
						{
							UUID:      "msg-3",
							Role:      "assistant",
							Timestamp: now.Add(time.Second),
							Content:   []ContentBlock{{Type: "text", Text: "hello there friend"}},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search with relevance sorting (default)
	result := idx.SearchWithOptions("hello", "", "", SearchOptions{Sort: "relevance"})

	if len(result.Results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result.Results))
	}

	// Session-2 should rank higher (more matches, recent, early position)
	if result.Results[0].SessionID != "session-2" {
		t.Errorf("Session-2 should rank first due to higher relevance, got %s", result.Results[0].SessionID)
	}

	// Verify scores are populated and ordered
	if result.Results[0].Score <= result.Results[1].Score {
		t.Errorf("First result score (%f) should be higher than second (%f)",
			result.Results[0].Score, result.Results[1].Score)
	}
}

func TestSearchWithOptions_SortByRecent(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-old",
					Summary: "Old session",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now.Add(-48 * time.Hour), // 2 days ago
							Content:   []ContentBlock{{Type: "text", Text: "hello hello hello"}}, // More matches
						},
					},
				},
				{
					ID:      "session-new",
					Summary: "New session",
					Messages: []Message{
						{
							UUID:      "msg-2",
							Role:      "user",
							Timestamp: now, // Now
							Content:   []ContentBlock{{Type: "text", Text: "hello"}}, // Fewer matches
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search with recent sorting
	result := idx.SearchWithOptions("hello", "", "", SearchOptions{Sort: "recent"})

	if len(result.Results) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result.Results))
	}

	// session-new should be first despite lower relevance
	if result.Results[0].SessionID != "session-new" {
		t.Errorf("session-new should rank first with sort=recent, got %s", result.Results[0].SessionID)
	}
}

func TestSearchWithOptions_PhraseScoring(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:      "session-1",
					Summary: "Has phrase",
					Messages: []Message{
						{
							UUID:      "msg-1",
							Role:      "user",
							Timestamp: now,
							Content:   []ContentBlock{{Type: "text", Text: "hello world is a common phrase"}},
						},
					},
				},
				{
					ID:      "session-2",
					Summary: "No phrase",
					Messages: []Message{
						{
							UUID:      "msg-2",
							Role:      "user",
							Timestamp: now,
							Content:   []ContentBlock{{Type: "text", Text: "world hello reversed order"}},
						},
					},
				},
			},
		},
	}

	idx := NewSearchIndex(projects)

	// Search for phrase
	result := idx.SearchWithOptions(`"hello world"`, "", "", SearchOptions{Sort: "relevance"})

	// Only session-1 should match the exact phrase
	if len(result.Results) != 1 {
		t.Fatalf("Expected 1 result for phrase search, got %d", len(result.Results))
	}

	if result.Results[0].SessionID != "session-1" {
		t.Errorf("Expected session-1 to match phrase, got %s", result.Results[0].SessionID)
	}

	// Score should include phrase bonus
	if result.Results[0].Score < 2.0 {
		t.Errorf("Phrase match should have score >= 2.0 (phrase bonus), got %f", result.Results[0].Score)
	}
}
