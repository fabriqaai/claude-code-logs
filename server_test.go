package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestNewServer(t *testing.T) {
	projects := []Project{
		{Path: "/test", FolderName: "-test", Sessions: []Session{}},
	}

	server, err := NewServer(8080, "/tmp/output", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	if server.port != 8080 {
		t.Errorf("port = %d, want 8080", server.port)
	}

	if server.outputDir != "/tmp/output" {
		t.Errorf("outputDir = %q, want /tmp/output", server.outputDir)
	}

	if server.index == nil {
		t.Error("index should not be nil")
	}
}

func TestHandleSearch_ValidRequest(t *testing.T) {
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
					},
				},
			},
		},
	}

	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// Create request
	reqBody := SearchRequest{Query: "hello"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	// Record response
	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	// Check status
	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	// Check response
	var response SearchResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if response.Query != "hello" {
		t.Errorf("response.Query = %q, want \"hello\"", response.Query)
	}

	if response.Total == 0 {
		t.Error("response.Total should be > 0")
	}
}

func TestHandleSearch_EmptyQuery(t *testing.T) {
	projects := []Project{}
	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	reqBody := SearchRequest{Query: ""}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Empty query should return OK, got %v", rr.Code)
	}

	var response SearchResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if response.Total != 0 {
		t.Errorf("Empty query should return 0 results, got %d", response.Total)
	}
}

func TestHandleSearch_InvalidMethod(t *testing.T) {
	projects := []Project{}
	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/search", nil)
	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("GET should return 405, got %v", rr.Code)
	}
}

func TestHandleSearch_InvalidJSON(t *testing.T) {
	projects := []Project{}
	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader([]byte("not json")))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("Invalid JSON should return 400, got %v", rr.Code)
	}
}

func TestHandleSearch_WithFilters(t *testing.T) {
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
						{UUID: "msg-1", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello world"}}},
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
					Summary: "Other Session",
					Messages: []Message{
						{UUID: "msg-2", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello there"}}},
					},
				},
			},
		},
	}

	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// Search with project filter
	reqBody := SearchRequest{Query: "hello", Project: "/Users/test/project1"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	var response SearchResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	if len(response.Results) != 1 {
		t.Errorf("Filtered search should return 1 result, got %d", len(response.Results))
	}

	if len(response.Results) > 0 && response.Results[0].Project != "/Users/test/project1" {
		t.Error("Filter should only return project1 results")
	}
}

func TestHandleSearch_LongQuery(t *testing.T) {
	projects := []Project{}
	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	// Create very long query
	longQuery := ""
	for i := 0; i < 2000; i++ {
		longQuery += "a"
	}

	reqBody := SearchRequest{Query: longQuery}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	// Should not error, just truncate
	if rr.Code != http.StatusOK {
		t.Errorf("Long query should return OK, got %v", rr.Code)
	}
}

func TestHandleStats(t *testing.T) {
	now := time.Now()
	projects := []Project{
		{
			Path:       "/Users/test/project1",
			FolderName: "-Users-test-project1",
			Sessions: []Session{
				{
					ID:        "session-1",
					Summary:   "Test",
					CreatedAt: now.Add(-time.Hour),
					UpdatedAt: now,
					Messages: []Message{
						{UUID: "msg-1", Role: "user", Timestamp: now,
							Content: []ContentBlock{{Type: "text", Text: "Hello"}}},
					},
				},
			},
		},
	}

	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodGet, "/api/stats", nil)
	rr := httptest.NewRecorder()
	server.handleStats(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Stats should return OK, got %v", rr.Code)
	}

	var stats StatsData
	if err := json.Unmarshal(rr.Body.Bytes(), &stats); err != nil {
		t.Fatalf("Failed to parse stats: %v", err)
	}

	if stats.TotalProjects != 1 {
		t.Errorf("Stats TotalProjects = %v, want 1", stats.TotalProjects)
	}

	if stats.TotalMessages != 1 {
		t.Errorf("Stats TotalMessages = %v, want 1", stats.TotalMessages)
	}

	if stats.TotalSessions != 1 {
		t.Errorf("Stats TotalSessions = %v, want 1", stats.TotalSessions)
	}

	if len(stats.MessagesPerDay) == 0 {
		t.Error("Stats should include MessagesPerDay time series")
	}

	if len(stats.ProjectStats) != 1 {
		t.Errorf("Stats should have 1 project stat, got %d", len(stats.ProjectStats))
	}
}

func TestHandleStats_InvalidMethod(t *testing.T) {
	projects := []Project{}
	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/stats", nil)
	rr := httptest.NewRecorder()
	server.handleStats(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST to stats should return 405, got %v", rr.Code)
	}
}

func TestHandleStatic(t *testing.T) {
	// Create temp directory with test files
	tmpDir, err := os.MkdirTemp("", "server-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projects := []Project{}
	server, err := NewServer(8080, tmpDir, projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	fileServer := http.FileServer(http.Dir(tmpDir))
	handler := server.handleStatic(fileServer)

	// Test root path - should render dynamic main index
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("GET / should return OK, got %v", rr.Code)
	}

	// Root path now renders dynamic HTML with the index template
	if !bytes.Contains(rr.Body.Bytes(), []byte("Claude Code Logs")) {
		t.Error("Response should contain dynamically rendered index content")
	}
}

func TestHandleStatic_NotFound(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "server-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projects := []Project{}
	server, err := NewServer(8080, tmpDir, projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	fileServer := http.FileServer(http.Dir(tmpDir))
	handler := server.handleStatic(fileServer)

	req := httptest.NewRequest(http.MethodGet, "/nonexistent.html", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("GET nonexistent should return 404, got %v", rr.Code)
	}
}

func TestHandleStatic_InvalidMethod(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "server-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	projects := []Project{}
	server, err := NewServer(8080, tmpDir, projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	fileServer := http.FileServer(http.Dir(tmpDir))
	handler := server.handleStatic(fileServer)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	rr := httptest.NewRecorder()
	handler(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Errorf("POST to static should return 405, got %v", rr.Code)
	}
}

func TestCorsMiddleware(t *testing.T) {
	handler := corsMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test regular request
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Error("CORS header should be set")
	}

	// Test preflight
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("OPTIONS should return OK, got %v", rr.Code)
	}
}

func TestSearchResponseFormat(t *testing.T) {
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
							Role:      "assistant",
							Timestamp: now,
							Content: []ContentBlock{
								{Type: "text", Text: "Hello, I can help you with that!"},
							},
						},
					},
				},
			},
		},
	}

	server, err := NewServer(8080, "/tmp", projects)
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	reqBody := SearchRequest{Query: "help"}
	body, _ := json.Marshal(reqBody)
	req := httptest.NewRequest(http.MethodPost, "/api/search", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	server.handleSearch(rr, req)

	var response SearchResponse
	json.Unmarshal(rr.Body.Bytes(), &response)

	// Verify response structure
	if len(response.Results) == 0 {
		t.Fatal("Expected at least one result")
	}

	result := response.Results[0]

	// Check all required fields
	if result.Project == "" {
		t.Error("Result should have project")
	}
	if result.SessionID == "" {
		t.Error("Result should have sessionId")
	}
	if result.SessionTitle == "" {
		t.Error("Result should have sessionTitle")
	}
	if len(result.Matches) == 0 {
		t.Error("Result should have matches")
	}

	match := result.Matches[0]
	if match.MessageID == "" {
		t.Error("Match should have messageId")
	}
	if match.Role == "" {
		t.Error("Match should have role")
	}
	if match.Content == "" {
		t.Error("Match should have content")
	}

	// Check highlighting
	if !bytes.Contains([]byte(match.Content), []byte("<mark>")) {
		t.Error("Match content should contain highlighted terms")
	}
}
