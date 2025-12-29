---
unit: 003-server
intent: 001-chat-log-viewer
phase: inception
status: draft
---

# Unit Brief: server

## Purpose

HTTP server that serves generated HTML and provides a search API for full-text search across sessions.

## Responsibility

- Serve static HTML files from output directory
- Provide `/api/search` endpoint for full-text search
- Build and maintain search index
- Return search results with context

## Assigned Requirements

- **FR-4**: Server Mode
- **FR-5**: Search Functionality

## Key Entities

### API Endpoints

```
GET  /                    # Serve index.html
GET  /projects/*          # Serve project/session HTML
GET  /assets/*            # Serve static assets
POST /api/search          # Full-text search

# Search Request
{
  "query": "string",
  "project": "optional-project-filter",
  "session": "optional-session-filter"
}

# Search Response
{
  "results": [
    {
      "project": "project-name",
      "session": "session-id",
      "sessionTitle": "Session summary",
      "matches": [
        {
          "messageId": "uuid",
          "role": "user|assistant",
          "content": "matched text with context",
          "timestamp": "2025-12-29T10:00:00Z"
        }
      ]
    }
  ],
  "total": 42,
  "query": "original query"
}
```

### Search Index

In-memory index built at startup:
- Index all message content
- Support filtering by project/session
- Highlight matching terms

## Key Operations

1. **StartServer(port, outputDir)** - Start HTTP server
2. **BuildIndex(projects)** - Build search index from parsed data
3. **Search(query, filters)** - Execute search, return results
4. **ServeStatic(path)** - Serve static files

## Dependencies

- **002-generator**: Provides generated HTML to serve
- **001-parser**: Provides data for building search index

## Interface

- CLI provides: port, output directory
- Server exposes: HTTP endpoints

## Technical Constraints

- Use Go `net/http` standard library
- In-memory search index (bleve or custom)
- Configurable port (default 8080)
- Localhost only (no external access)
- Graceful shutdown on SIGINT/SIGTERM

## Success Criteria

- [ ] Serves generated HTML correctly
- [ ] Search returns relevant results
- [ ] Search results include context (surrounding text)
- [ ] Highlights matching terms
- [ ] Filter by project/session works
- [ ] Response time < 500ms for search
- [ ] Graceful shutdown

---

## Story Summary

- **Total Stories**: 1
- **Must Have**: 1
- **Should Have**: 0
- **Could Have**: 0

### Stories

- [ ] **001-serve-and-search**: Serve HTML and search API - Must - Planned
