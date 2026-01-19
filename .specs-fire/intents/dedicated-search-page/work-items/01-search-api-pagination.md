---
id: 01-search-api-pagination
title: Extend Search API with Pagination
intent: dedicated-search-page
complexity: medium
mode: confirm
status: pending
depends_on: []
created: 2026-01-20
---

# Work Item: Extend Search API with Pagination

## Description

Add pagination support to the `/api/search` endpoint to handle large result sets efficiently. The API should accept offset/limit parameters and return pagination metadata.

## Acceptance Criteria

- [ ] `SearchRequest` accepts `offset` (default 0) and `limit` (default 20, max 100) params
- [ ] `SearchResponse` includes `total` (total matching sessions), `hasMore` boolean, `offset` value
- [ ] Results are sliced server-side based on offset/limit after scoring/sorting
- [ ] Existing clients (overlay) continue to work (backwards compatible defaults)

## Technical Notes

- Modify `search.go`: update `SearchRequest` and `SearchResponse` structs
- Modify `Search()` method to accept pagination params
- Apply pagination AFTER filtering and sorting (so total count is accurate)
- Consider: paginate by sessions (not individual matches within sessions)

## Files to Modify

- `search.go` - Request/Response structs, Search method
- `server.go` - handleSearch to parse new params
