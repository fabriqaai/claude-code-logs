# Walkthrough: Stats API Endpoint

## Summary

Extended `/api/stats` endpoint to return full analytics data (`StatsData`) for the stats dashboard, including time-series data for charting.

## What Changed

### server.go

**1. Added `stats` field to `Server` struct:**
```go
type Server struct {
    // ... existing fields ...
    stats *StatsData  // Precomputed stats for the stats API
}
```

**2. Compute stats at server startup:**
```go
return &Server{
    // ... existing fields ...
    stats: ComputeStats(projects),
}, nil
```

**3. Modified `handleStats` handler:**
- Returns full `StatsData` structure as JSON
- Supports optional `?range=` query parameter for time filtering
- Uses `FilterStatsByTimeRange` from stats.go for filtering

### server_test.go

Updated `TestHandleStats` to validate the new response structure:
- Checks `TotalProjects`, `TotalMessages`, `TotalSessions`
- Validates `MessagesPerDay` time series exists
- Verifies `ProjectStats` array

## API Response

The `/api/stats` endpoint now returns:

```json
{
  "totalProjects": 5,
  "totalSessions": 42,
  "totalMessages": 1234,
  "totalTokens": 567890,
  "messagesPerDay": [{"date": "2026-01-20", "value": 15}, ...],
  "tokensPerDay": [{"date": "2026-01-20", "value": 1200}, ...],
  "projectStats": [{"path": "/project", "slug": "project", "sessions": 10, ...}],
  "avgSessionLengthMins": 23.5,
  "avgMessagesPerSession": 29.4,
  "computedAt": "2026-01-20T10:00:00Z"
}
```

**Query Parameters:**
- `?range=today` - Filter to today only
- `?range=week` - Filter to last 7 days
- `?range=month` - Filter to last 30 days
- `?range=all` or no param - Return all data

## Verification

```bash
# Start the server
go run . serve

# Test the endpoint
curl http://localhost:8080/api/stats | jq .

# Test with time range filter
curl "http://localhost:8080/api/stats?range=week" | jq .
```

## Acceptance Criteria

- [x] `/api/stats` returns full `StatsData` as JSON
- [x] Response includes time-series arrays for charts
- [x] Response includes summary totals
- [x] Endpoint remains fast (<100ms for typical dataset) - stats precomputed at startup
- [x] Existing stats fields preserved for backward compatibility
