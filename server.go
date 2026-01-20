package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"
)

// cacheEntry represents a cached rendered HTML page
type cacheEntry struct {
	content   []byte
	timestamp time.Time
}

// Server represents the HTTP server for serving HTML and search API
type Server struct {
	port         int
	outputDir    string
	index        *SearchIndex
	projects     []Project
	server       *http.Server
	shellTmpl    *template.Template
	indexTmpl    *template.Template
	projectTmpl  *template.Template
	statsTmpl    *template.Template
	searchTmpl   *template.Template
	// Cache for rendered HTML pages
	cache      map[string]*cacheEntry
	cacheMu    sync.RWMutex
	cacheTTL   time.Duration
	// Precomputed stats for the stats API
	stats      *StatsData
}

// NewServer creates a new server instance
func NewServer(port int, outputDir string, projects []Project) (*Server, error) {
	funcMap := template.FuncMap{
		"ProjectSlug": ProjectSlug,
	}

	shellTmpl, err := template.New("shell").Funcs(funcMap).Parse(sessionShellTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing shell template: %w", err)
	}

	indexTmpl, err := template.New("index").Funcs(funcMap).Parse(indexTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing index template: %w", err)
	}

	projectTmpl, err := template.New("project").Funcs(funcMap).Parse(projectIndexTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing project template: %w", err)
	}

	statsTmpl, err := template.New("stats").Funcs(funcMap).Parse(statsTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing stats template: %w", err)
	}

	searchTmpl, err := template.New("search").Funcs(funcMap).Parse(searchTemplate)
	if err != nil {
		return nil, fmt.Errorf("parsing search template: %w", err)
	}

	return &Server{
		port:        port,
		outputDir:   outputDir,
		projects:    projects,
		index:       NewSearchIndex(projects),
		shellTmpl:   shellTmpl,
		indexTmpl:   indexTmpl,
		projectTmpl: projectTmpl,
		statsTmpl:   statsTmpl,
		searchTmpl:  searchTmpl,
		cache:       make(map[string]*cacheEntry),
		cacheTTL:    30 * time.Second, // Cache HTML for 30 seconds
		stats:       ComputeStats(projects),
	}, nil
}

// Start starts the HTTP server and blocks until shutdown
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API routes
	mux.HandleFunc("/api/search", s.handleSearch)
	mux.HandleFunc("/api/stats", s.handleStats)

	// Static file serving
	fileServer := http.FileServer(http.Dir(s.outputDir))
	mux.HandleFunc("/", s.handleStatic(fileServer))

	// Create server
	addr := fmt.Sprintf("127.0.0.1:%d", s.port)
	s.server = &http.Server{
		Addr:         addr,
		Handler:      corsMiddleware(mux),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Check if port is available
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		if strings.Contains(err.Error(), "address already in use") {
			return fmt.Errorf("port %d is already in use. Try a different port with --port flag", s.port)
		}
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	// Setup graceful shutdown
	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		fmt.Println("\nShutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.server.Shutdown(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Server shutdown error: %v\n", err)
		}
		close(done)
	}()

	// Print startup message
	fmt.Printf("Server starting on http://%s\n", addr)
	fmt.Printf("Search index: %d messages, %d terms\n", s.index.MessageCount(), s.index.TermCount())
	fmt.Println("Press Ctrl+C to stop")

	// Start server
	if err := s.server.Serve(listener); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("server error: %w", err)
	}

	<-done
	fmt.Println("Server stopped")
	return nil
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error {
	if s.server == nil {
		return nil
	}
	return s.server.Shutdown(ctx)
}

// handleStatic serves static files and renders HTML dynamically
func (s *Server) handleStatic(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET and HEAD
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Clean the path
		path := filepath.Clean(r.URL.Path)

		// Serve embedded logo
		if path == "/claude-code-icon.png" {
			w.Header().Set("Content-Type", "image/png")
			w.Header().Set("Cache-Control", "public, max-age=86400")
			w.Write(logoData)
			return
		}

		// Handle root - render main index
		if path == "/" || path == "" || path == "/index.html" {
			s.renderMainIndex(w, r)
			return
		}

		// Handle stats page
		if path == "/stats" || path == "/stats.html" {
			s.renderStatsPage(w, r)
			return
		}

		// Handle search page
		if path == "/search" || path == "/search.html" {
			s.renderSearchPage(w, r)
			return
		}

		// Handle project and session paths without .html extension
		trimmedPath := strings.TrimPrefix(path, "/")
		trimmedPath = strings.TrimSuffix(trimmedPath, "/")
		pathParts := strings.Split(trimmedPath, "/")

		// Single path segment without extension - could be a project slug
		if len(pathParts) == 1 && !strings.Contains(trimmedPath, ".") {
			for i := range s.projects {
				if ProjectSlug(s.projects[i].Path) == trimmedPath {
					s.renderProjectIndex(w, r, trimmedPath)
					return
				}
			}
		}

		// Two path segments without extension - could be project/session
		if len(pathParts) == 2 && !strings.Contains(pathParts[1], ".") {
			projectSlug := pathParts[0]
			sessionID := pathParts[1]
			// Verify this is a valid project/session combination
			for i := range s.projects {
				if ProjectSlug(s.projects[i].Path) == projectSlug {
					for j := range s.projects[i].Sessions {
						if s.projects[i].Sessions[j].ID == sessionID {
							s.renderSessionShell(w, r, projectSlug, sessionID)
							return
						}
					}
					break
				}
			}
		}

		// Check if this is a session HTML request (e.g., /project-slug/session-id.html)
		if strings.HasSuffix(path, ".html") && path != "/index.html" {
			parts := strings.Split(strings.TrimPrefix(path, "/"), "/")
			if len(parts) == 2 {
				projectSlug := parts[0]
				filename := parts[1]

				// Check if it's a project index
				if filename == "index.html" {
					s.renderProjectIndex(w, r, projectSlug)
					return
				}

				// Otherwise it's a session page
				sessionID := strings.TrimSuffix(filename, ".html")
				s.renderSessionShell(w, r, projectSlug, sessionID)
				return
			}
		}

		// For static files (.md, .png, etc.), check if file exists
		fullPath := filepath.Join(s.outputDir, path)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			http.NotFound(w, r)
			return
		}

		// Set content type based on extension
		ext := filepath.Ext(path)
		switch ext {
		case ".html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
		case ".css":
			w.Header().Set("Content-Type", "text/css; charset=utf-8")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript; charset=utf-8")
		case ".json":
			w.Header().Set("Content-Type", "application/json; charset=utf-8")
		case ".md":
			w.Header().Set("Content-Type", "text/markdown; charset=utf-8")
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	}
}

// getFromCache returns cached content if available and not expired
func (s *Server) getFromCache(key string) ([]byte, bool) {
	s.cacheMu.RLock()
	defer s.cacheMu.RUnlock()

	entry, ok := s.cache[key]
	if !ok {
		return nil, false
	}

	if time.Since(entry.timestamp) > s.cacheTTL {
		return nil, false
	}

	return entry.content, true
}

// setInCache stores content in the cache
func (s *Server) setInCache(key string, content []byte) {
	s.cacheMu.Lock()
	defer s.cacheMu.Unlock()

	s.cache[key] = &cacheEntry{
		content:   content,
		timestamp: time.Now(),
	}
}

// renderMainIndex renders the main index page dynamically
func (s *Server) renderMainIndex(w http.ResponseWriter, r *http.Request) {
	cacheKey := "index"

	// Check cache first
	if content, ok := s.getFromCache(cacheKey); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
		return
	}

	data := struct {
		Projects []Project
	}{
		Projects: s.projects,
	}

	// Render to buffer for caching
	var buf bytes.Buffer
	if err := s.indexTmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Index template error: %v\n", err)
		return
	}

	content := buf.Bytes()
	s.setInCache(cacheKey, content)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// renderProjectIndex renders a project index page dynamically
func (s *Server) renderProjectIndex(w http.ResponseWriter, r *http.Request, projectSlug string) {
	cacheKey := "project:" + projectSlug

	// Check cache first
	if content, ok := s.getFromCache(cacheKey); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
		return
	}

	// Find the project
	var project *Project
	for i := range s.projects {
		if ProjectSlug(s.projects[i].Path) == projectSlug {
			project = &s.projects[i]
			break
		}
	}

	if project == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Project     *Project
		AllProjects []Project
	}{
		Project:     project,
		AllProjects: s.projects,
	}

	// Render to buffer for caching
	var buf bytes.Buffer
	if err := s.projectTmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Project template error: %v\n", err)
		return
	}

	content := buf.Bytes()
	s.setInCache(cacheKey, content)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// renderSessionShell renders a session HTML shell dynamically
func (s *Server) renderSessionShell(w http.ResponseWriter, r *http.Request, projectSlug, sessionID string) {
	cacheKey := "session:" + projectSlug + "/" + sessionID

	// Check cache first
	if content, ok := s.getFromCache(cacheKey); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
		return
	}

	// Find the project and session
	var project *Project
	var session *Session

	for i := range s.projects {
		if ProjectSlug(s.projects[i].Path) == projectSlug {
			project = &s.projects[i]
			for j := range project.Sessions {
				if project.Sessions[j].ID == sessionID {
					session = &project.Sessions[j]
					break
				}
			}
			break
		}
	}

	if project == nil || session == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Session     *Session
		Project     *Project
		AllProjects []Project
	}{
		Session:     session,
		Project:     project,
		AllProjects: s.projects,
	}

	// Render to buffer for caching
	var buf bytes.Buffer
	if err := s.shellTmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Shell template error: %v\n", err)
		return
	}

	content := buf.Bytes()
	s.setInCache(cacheKey, content)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// renderStatsPage renders the stats dashboard page
func (s *Server) renderStatsPage(w http.ResponseWriter, r *http.Request) {
	cacheKey := "stats"

	// Check cache first
	if content, ok := s.getFromCache(cacheKey); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(content)
		return
	}

	data := struct {
		Projects []Project
	}{
		Projects: s.projects,
	}

	// Render to buffer for caching
	var buf bytes.Buffer
	if err := s.statsTmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Stats template error: %v\n", err)
		return
	}

	content := buf.Bytes()
	s.setInCache(cacheKey, content)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// renderSearchPage renders the search page
func (s *Server) renderSearchPage(w http.ResponseWriter, r *http.Request) {
	cacheKey := "search"

	// Check cache first
	if content, ok := s.getFromCache(cacheKey); ok {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(content) // Error ignored: client may have disconnected
		return
	}

	data := struct {
		Projects []Project
	}{
		Projects: s.projects,
	}

	// Render to buffer for caching
	var buf bytes.Buffer
	if err := s.searchTmpl.Execute(&buf, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		fmt.Fprintf(os.Stderr, "Search template error: %v\n", err)
		return
	}

	content := buf.Bytes()
	s.setInCache(cacheKey, content)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write(content) // Error ignored: client may have disconnected
}

// handleSearch handles POST /api/search requests
func (s *Server) handleSearch(w http.ResponseWriter, r *http.Request) {
	// Only allow POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request
	var req SearchRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}

	// Validate query length
	const maxQueryLength = 1000
	if len(req.Query) > maxQueryLength {
		req.Query = req.Query[:maxQueryLength]
	}

	// Execute search with options
	start := time.Now()
	opts := SearchOptions{
		Offset: req.Offset,
		Limit:  req.Limit,
		Sort:   req.Sort,
	}
	searchResult := s.index.SearchWithOptions(req.Query, req.Project, req.Session, opts)
	duration := time.Since(start)

	// Count total matches across returned results
	totalMatches := 0
	for _, result := range searchResult.Results {
		totalMatches += len(result.Matches)
	}

	// Build response
	response := SearchResponse{
		Results: searchResult.Results,
		Total:   searchResult.Total,
		Query:   req.Query,
		HasMore: searchResult.HasMore,
		Offset:  searchResult.Offset,
	}

	// Log search (for debugging)
	fmt.Printf("Search: %q -> %d/%d results (offset=%d, limit=%d) in %v\n",
		req.Query, len(searchResult.Results), searchResult.Total, req.Offset, req.Limit, duration)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding response: %v\n", err)
	}
}

// handleStats returns full analytics data for the stats dashboard
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check for time range filter
	rangeType := r.URL.Query().Get("range")
	stats := s.stats
	if rangeType != "" && rangeType != "all" {
		stats = FilterStatsByTimeRange(s.stats, rangeType)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(stats)
}

// corsMiddleware adds CORS headers for local development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Allow requests from any origin (localhost only anyway)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		// Handle preflight
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// StartServer is a convenience function to start a server
func StartServer(port int, outputDir string, projects []Project) error {
	server, err := NewServer(port, outputDir, projects)
	if err != nil {
		return fmt.Errorf("creating server: %w", err)
	}
	return server.Start()
}
