# Run 004 - Implementation Walkthrough

## Summary

This run implemented **pagination** and **relevance scoring** for the search API, preparing the foundation for the dedicated search page.

---

## What Changed

### 1. New Types in `search.go`

```go
// SearchOptions - pagination and sort parameters
type SearchOptions struct {
    Offset int
    Limit  int
    Sort   string // "relevance" (default) or "recent"
}

// SearchResultWithPagination - includes pagination metadata
type SearchResultWithPagination struct {
    Results []SearchResult
    Total   int
    HasMore bool
    Offset  int
}
```

### 2. Updated Structs

**SearchRequest** now accepts:
- `offset` - starting position (default: 0)
- `limit` - max results to return (default: 20, max: 100)
- `sort` - "relevance" or "recent"

**SearchResponse** now includes:
- `hasMore` - boolean indicating more results exist
- `offset` - current offset for pagination state

**SearchResult** now includes:
- `score` - relevance score (float64)

### 3. New Methods

**`SearchWithOptions()`** - Main search method with pagination/sorting:
1. Parses query for terms and phrases
2. Finds matching messages via inverted index
3. Groups matches by session
4. Calculates relevance score per session
5. Sorts by score (default) or recency
6. Applies pagination (offset/limit)
7. Returns results with pagination metadata

**`calculateScore()`** - Computes relevance score:
- Base: 1.0 per matching message in session
- Term frequency: +0.1 per extra occurrence (capped at 1.0)
- Position bonus: +0.5 if match in first 200 chars
- Phrase bonus: +2.0 per exact phrase match
- Recency bonus: +0.2 if within last 7 days

### 4. Server Integration

`handleSearch()` in `server.go` now:
1. Extracts `offset`, `limit`, `sort` from request JSON
2. Passes options to `SearchWithOptions()`
3. Returns pagination metadata in response

---

## How to Verify

### Test the API

```bash
# Start server
go run . serve

# Search with defaults (relevance sort, limit 20)
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"query": "hello"}'

# Search with pagination
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"query": "hello", "offset": 0, "limit": 10}'

# Search sorted by recency
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"query": "hello", "sort": "recent"}'

# Load more (page 2)
curl -X POST http://localhost:8080/api/search \
  -H "Content-Type: application/json" \
  -d '{"query": "hello", "offset": 20, "limit": 20}'
```

### Expected Response Format

```json
{
  "results": [
    {
      "project": "/path/to/project",
      "projectSlug": "project-name",
      "sessionId": "abc123",
      "sessionTitle": "Session Title",
      "matches": [...],
      "score": 3.7
    }
  ],
  "total": 45,
  "query": "hello",
  "hasMore": true,
  "offset": 0
}
```

### Run Tests

```bash
# Run all search tests
go test -v -run "TestSearch" ./...

# Run specific pagination tests
go test -v -run "TestSearchWithOptions" ./...
```

---

## Backwards Compatibility

The original `Search()` method is preserved and works exactly as before:
- Returns `[]SearchResult` (not the new wrapper type)
- Results still sorted by recency (as before)
- No score field was breaking change - just added

Existing overlay search will continue to work without changes.

---

## Next Steps

The API is now ready for the dedicated search page (work items 03-05):
1. **03-search-page-template** - Create the search page UI
2. **04-search-page-route** - Wire up `/search` route
3. **05-remove-overlay** - Remove old overlay search
