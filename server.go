package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

// Server represents the HTTP server for serving HTML and search API
type Server struct {
	port      int
	outputDir string
	index     *SearchIndex
	projects  []Project
	server    *http.Server
}

// NewServer creates a new server instance
func NewServer(port int, outputDir string, projects []Project) *Server {
	return &Server{
		port:      port,
		outputDir: outputDir,
		projects:  projects,
		index:     NewSearchIndex(projects),
	}
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

// handleStatic serves static files with proper content types
func (s *Server) handleStatic(fileServer http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow GET and HEAD
		if r.Method != http.MethodGet && r.Method != http.MethodHead {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Clean the path
		path := filepath.Clean(r.URL.Path)

		// Serve index.html for root
		if path == "/" || path == "" {
			path = "/index.html"
		}

		// Check if file exists
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
		}

		// Serve the file
		fileServer.ServeHTTP(w, r)
	}
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

	// Execute search
	start := time.Now()
	results := s.index.Search(req.Query, req.Project, req.Session)
	duration := time.Since(start)

	// Count total matches
	totalMatches := 0
	for _, result := range results {
		totalMatches += len(result.Matches)
	}

	// Build response
	response := SearchResponse{
		Results: results,
		Total:   totalMatches,
		Query:   req.Query,
	}

	// Log search (for debugging)
	fmt.Printf("Search: %q -> %d results in %v\n", req.Query, totalMatches, duration)

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding response: %v\n", err)
	}
}

// handleStats returns index statistics
func (s *Server) handleStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := map[string]interface{}{
		"projects": len(s.projects),
		"messages": s.index.MessageCount(),
		"terms":    s.index.TermCount(),
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
	server := NewServer(port, outputDir, projects)
	return server.Start()
}
