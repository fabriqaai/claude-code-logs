---
id: 02-stats-api-endpoint
title: Stats API Endpoint
intent: session-stats-dashboard
complexity: low
mode: autopilot
status: completed
depends_on: [01-stats-data-model]
created: 2025-01-20
---

# Work Item: Stats API Endpoint

## Description

Extend the existing `/api/stats` endpoint to return full analytics data including time-series information for charting.

## Acceptance Criteria

- [ ] `/api/stats` returns full `StatsData` as JSON
- [ ] Response includes time-series arrays for charts
- [ ] Response includes summary totals
- [ ] Endpoint remains fast (<100ms for typical dataset)
- [ ] Existing stats fields preserved for backward compatibility

## Technical Notes

- Modify `handleStats` in `server.go`
- Compute stats once at server startup, store in `Server` struct
- Return cached stats (already computed in work item 1)
- JSON structure should be chart-friendly (arrays of {date, value} objects)
