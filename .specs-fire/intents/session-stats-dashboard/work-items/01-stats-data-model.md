---
id: 01-stats-data-model
title: Stats Data Model & Computation
intent: session-stats-dashboard
complexity: medium
mode: confirm
status: completed
depends_on: []
created: 2025-01-20
---

# Work Item: Stats Data Model & Computation

## Description

Create Go structs for analytics data and implement computation logic that processes loaded projects to generate stats. This is the foundation for all stats features.

## Acceptance Criteria

- [ ] `StatsData` struct defined with fields for all metrics
- [ ] `ComputeStats(projects []Project) *StatsData` function implemented
- [ ] Messages per day/week computed correctly
- [ ] Token usage estimated from message content length (rough approximation: chars/4)
- [ ] Active projects per time period tracked
- [ ] Session duration/length calculated (first to last message)
- [ ] Unit tests for computation logic

## Technical Notes

- Add to new file `stats.go`
- Token estimation: use character count / 4 as rough approximation
- Time buckets: daily for last 30 days, weekly for last 12 weeks
- Consider caching computed stats in Server struct
- Reuse existing `Project`, `Session`, `Message` types
