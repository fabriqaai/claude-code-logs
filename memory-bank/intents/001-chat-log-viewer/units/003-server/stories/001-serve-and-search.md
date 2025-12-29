---
id: 001-serve-and-search
unit: 003-server
intent: 001-chat-log-viewer
status: complete
priority: must
created: 2025-12-29T12:00:00.000Z
assigned_bolt: 003-server
implemented: true
---

# Story: 001-serve-and-search

## User Story

**As a** developer
**I want** a local server with search
**So that** I can browse and search my chat logs

## Acceptance Criteria

- [ ] **Given** `claude-code-logs serve` run, **When** started, **Then** HTTP server listens on configurable port (default 8080)
- [ ] **Given** server running, **When** accessing `/`, **Then** index.html is served
- [ ] **Given** server running, **When** POST `/api/search` with query, **Then** matching results returned with context
- [ ] **Given** search request with project filter, **When** executing, **Then** results limited to that project
- [ ] **Given** any search query, **When** responding, **Then** response time < 500ms
- [ ] **Given** server running, **When** Ctrl+C pressed, **Then** graceful shutdown within 5 seconds
- [ ] **Given** search results, **When** returned, **Then** matching terms are highlighted

## Technical Notes

- Use Go `net/http` standard library
- Build in-memory search index at startup from parsed data
- Consider using bleve for full-text search or custom inverted index
- Bind to localhost only (127.0.0.1) for security
- CORS headers for local development
- JSON request/response format for search API

## Dependencies

### Requires
- 002-generator (provides HTML files to serve)
- 001-parser (provides data for search index)

### Enables
- 005-cli (exposes serve command)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Port already in use | Clear error message, suggest different port |
| No search results | Return empty array, not error |
| Query with special regex chars | Escape properly, treat as literal |
| Very long query (1000+ chars) | Truncate, proceed with search |
| Concurrent search requests | Handle safely (mutex or immutable index) |
| Search during index rebuild | Queue request or use old index |

## Out of Scope

- Authentication/authorization (localhost only)
- HTTPS (local use)
- WebSocket for live updates (future)
- Pagination (return all results, handle in UI)
