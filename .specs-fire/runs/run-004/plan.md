# Run 004 - Implementation Plan

**Status**: Approved
**Approved**: 2026-01-20
**Mode**: Confirm (batch)

---

## Work Items

1. **01-search-api-pagination** - Extend Search API with Pagination
2. **02-relevance-scoring** - Add Relevance Scoring

---

## Approach

Both work items modify the same search API infrastructure. Implementing together allows a clean, unified approach:

1. **Pagination**: Add `Offset`/`Limit` to `SearchRequest`, add `Total`/`HasMore`/`Offset` to `SearchResponse`, slice results server-side after scoring/sorting

2. **Relevance Scoring**: Add `Score` field to `SearchResult`, implement scoring formula based on term frequency/position/phrase bonus/recency, add `Sort` parameter (`relevance`/`recent`)

3. **Server Integration**: Update `handleSearch` to parse new parameters from JSON body

---

## Files to Modify

| File | Changes |
|------|---------|
| `search.go` | Update `SearchRequest` (add Offset, Limit, Sort), update `SearchResponse` (add HasMore, Offset), add `Score` to `SearchResult`, implement `scoreResult()` function, modify `Search()` to accept pagination/sort params |
| `server.go` | Update `handleSearch` to pass pagination/sort params to search |

---

## Implementation Details

### SearchRequest changes

```go
type SearchRequest struct {
    Query   string `json:"query"`
    Project string `json:"project,omitempty"`
    Session string `json:"session,omitempty"`
    Offset  int    `json:"offset,omitempty"`  // default 0
    Limit   int    `json:"limit,omitempty"`   // default 20, max 100
    Sort    string `json:"sort,omitempty"`    // "relevance" (default) or "recent"
}
```

### SearchResponse changes

```go
type SearchResponse struct {
    Results []SearchResult `json:"results"`
    Total   int            `json:"total"`   // Total matching sessions
    Query   string         `json:"query"`
    HasMore bool           `json:"hasMore"`
    Offset  int            `json:"offset"`
}
```

### SearchResult changes

```go
type SearchResult struct {
    // ... existing fields ...
    Score float64 `json:"score"` // relevance score
}
```

### Scoring Formula

- Base: 1.0 per matching message
- Term frequency: +0.1 per additional term occurrence (capped at 1.0)
- Position bonus: +0.5 if match in first 200 chars of any message
- Exact phrase bonus: +2.0 per phrase match
- Recency bonus: +0.2 if any match within last 7 days

---

## Tests

| Test | Purpose |
|------|---------|
| `TestSearch_Pagination` | Verify offset/limit slicing works correctly |
| `TestSearch_PaginationDefaults` | Verify default values (offset=0, limit=20) |
| `TestSearch_PaginationMaxLimit` | Verify limit caps at 100 |
| `TestSearch_RelevanceScoring` | Verify scoring increases with term frequency |
| `TestSearch_SortByRelevance` | Verify results sorted by score descending |
| `TestSearch_SortByRecent` | Verify sort=recent uses timestamp ordering |

---

## Acceptance Criteria

### Pagination
- [ ] `SearchRequest` accepts `offset` (default 0) and `limit` (default 20, max 100)
- [ ] `SearchResponse` includes `total`, `hasMore`, `offset`
- [ ] Results sliced server-side after scoring/sorting
- [ ] Backwards compatible (existing clients work)

### Relevance Scoring
- [ ] Each `SearchResult` includes `score` field
- [ ] Scoring factors: term frequency, position, phrase, recency
- [ ] Results sorted by score descending by default
- [ ] API accepts `sort` param: `relevance` or `recent`
