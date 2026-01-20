---
id: 01-hook-search-page
title: Hook Search Page into Server
intent: dedicated-search-page
complexity: low
mode: autopilot
status: pending
depends_on: []
created: 2026-01-20
---

# Work Item: Hook Search Page into Server

## Description

Integrate the existing `templates_search.go` template into the server routing layer. The template is complete - this work item connects it to the server so `/search` URL serves the page.

## Acceptance Criteria

- [ ] `/search` URL serves the search page
- [ ] Page renders correctly with sidebar and project tree
- [ ] Search functionality works via existing `/api/search` API calls
- [ ] Cache strategy implemented (cache key: "search")

## Technical Notes

Follow the pattern established by the stats page:

1. Add `searchTmpl *template.Template` to Server struct in `server.go`
2. Parse template in `NewServer()`: `template.New("search").Parse(searchTemplate)`
3. Create `renderSearchPage()` function similar to `renderStatsPage()`
4. Add route handling in `handleStatic()` for `/search` and `/search.html` paths

Reference files:
- `templates_search.go` - existing complete template
- `server.go` lines 184-188 - stats page routing pattern
- `server.go` `renderStatsPage()` - function pattern to follow

## Dependencies

None - this is the foundation work item.
