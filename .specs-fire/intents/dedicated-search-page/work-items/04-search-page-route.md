---
id: 04-search-page-route
title: Add Search Page Route
intent: dedicated-search-page
complexity: low
mode: autopilot
status: pending
depends_on:
  - 03-search-page-template
created: 2026-01-20
---

# Work Item: Add Search Page Route

## Description

Register the `/search` route in the server and wire up template rendering with caching.

## Acceptance Criteria

- [ ] `/search` route serves the search page template
- [ ] Template is parsed in `NewServer()` like other templates
- [ ] Page rendering uses cache (same pattern as stats page)
- [ ] Clean URLs work: `/search` and `/search?q=term`

## Technical Notes

- Follow exact pattern of `renderStatsPage` in server.go
- Add `searchTmpl` field to Server struct
- Parse template in NewServer
- Add route handler in handleStatic or as explicit route

## Files to Modify

- `server.go` - Add searchTmpl, parse in NewServer, add route handler
