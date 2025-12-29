---
stage: plan
bolt: 003-server
created: 2025-12-29T17:20:00Z
---

## Implementation Plan: server

### Objective

Implement a local HTTP server that serves the generated HTML files and provides a full-text search API. The server enables interactive browsing and searching of Claude Code chat logs.

### Deliverables

- `server.go` - HTTP server implementation with static file serving and search API
- `search.go` - Search index and query logic

### Dependencies

- **001-parser**: Provides `[]Project`, `Session`, `Message` types for building search index
- **002-generator**: Provides generated HTML files to serve
- **Go stdlib**: `net/http` for server, `encoding/json` for API, `os/signal` for graceful shutdown

### Technical Approach

#### 1. Server Architecture

```
Server
├── Static File Handler (GET /*)
│   └── Serves files from outputDir
├── Search API Handler (POST /api/search)
│   └── Queries in-memory search index
└── Search Index
    └── Built at startup from parsed projects
```

#### 2. Static File Serving

Use `http.FileServer` with the output directory:
- Map `/` to `index.html`
- Map `/projects/*` to project directories
- Set appropriate MIME types
- Return 404 for missing files

#### 3. Search Index Design

Simple inverted index (no external dependencies):

```go
type SearchIndex struct {
    // Map term -> list of occurrences
    index map[string][]Match

    // All messages for context retrieval
    messages []IndexedMessage
}

type IndexedMessage struct {
    Project     string
    SessionID   string
    SessionTitle string
    MessageID   string
    Role        string
    Content     string
    Timestamp   time.Time
}

type Match struct {
    MessageIndex int
    Position     int
}
```

**Indexing Strategy**:
- Tokenize message content into lowercase words
- Strip punctuation
- Build inverted index for fast lookup
- Store original content for highlighting

#### 4. Search API

**Endpoint**: `POST /api/search`

**Request**:
```json
{
    "query": "search terms",
    "project": "optional-project-path",
    "session": "optional-session-id"
}
```

**Response**:
```json
{
    "results": [
        {
            "project": "/Users/test/project",
            "sessionId": "abc-123",
            "sessionTitle": "Session summary",
            "matches": [
                {
                    "messageId": "msg-uuid",
                    "role": "assistant",
                    "content": "...matched <mark>text</mark> with context...",
                    "timestamp": "2025-12-29T10:00:00Z"
                }
            ]
        }
    ],
    "total": 5,
    "query": "search terms"
}
```

**Search Algorithm**:
1. Tokenize query into terms
2. Find messages containing ALL terms (AND logic)
3. Apply project/session filters if provided
4. Extract context around matches (50 chars before/after)
5. Wrap matched terms in `<mark>` tags
6. Sort by relevance (number of term occurrences)

#### 5. Server Configuration

| Setting | Default | Description |
|---------|---------|-------------|
| Port | 8080 | HTTP listening port |
| Host | 127.0.0.1 | Localhost only for security |
| OutputDir | ~/.claude-code-logs | Directory with generated HTML |

#### 6. Graceful Shutdown

- Listen for SIGINT/SIGTERM
- Stop accepting new connections
- Complete in-flight requests (5s timeout)
- Clean up resources

### Acceptance Criteria

- [ ] HTTP server listens on configurable port (default 8080)
- [ ] Binds to localhost only (127.0.0.1)
- [ ] Serves index.html at `/`
- [ ] Serves project/session HTML files
- [ ] POST `/api/search` returns JSON results
- [ ] Search results include context (50 chars around match)
- [ ] Matching terms wrapped in `<mark>` tags
- [ ] Filter by project works
- [ ] Filter by session works
- [ ] Response time < 500ms for typical searches
- [ ] Graceful shutdown on Ctrl+C
- [ ] Clear error message if port in use
- [ ] Empty query returns no results (not error)
- [ ] Special characters escaped in search
- [ ] CORS headers for local development

### API Design

```go
// Server represents the HTTP server
type Server struct {
    port      int
    outputDir string
    index     *SearchIndex
    projects  []Project
}

// NewServer creates a new server instance
func NewServer(port int, outputDir string, projects []Project) *Server

// Start starts the HTTP server (blocking)
func (s *Server) Start() error

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown(ctx context.Context) error

// SearchIndex provides full-text search
type SearchIndex struct {...}

// NewSearchIndex builds index from projects
func NewSearchIndex(projects []Project) *SearchIndex

// Search executes a search query
func (idx *SearchIndex) Search(query string, projectFilter, sessionFilter string) []SearchResult
```

### File Structure

```
claude-code-logs/
├── server.go         # HTTP server implementation
├── server_test.go    # Server tests
├── search.go         # Search index and query logic
├── search_test.go    # Search tests
└── ... (existing files)
```

### Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Large index memory usage | Index only text content, not tool blocks |
| Slow search on many sessions | Use inverted index for O(1) term lookup |
| Concurrent access issues | Index is immutable after build |
| Port conflicts | Clear error message with suggestion |

### Out of Scope

- HTTPS (localhost only)
- Authentication
- WebSocket for live updates
- Pagination (return all results)
- Index persistence (rebuilt on each start)
