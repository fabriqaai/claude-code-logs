---
id: 01-stats-data-model-extension
title: Extend Stats Data Model
intent: project-stats-filtering
complexity: medium
mode: confirm
status: completed
depends_on: []
created: 2026-01-21
---

# Work Item: Extend Stats Data Model

## Description

Add per-project time series data to `StatsData` so charts can be filtered by project client-side. Currently `MessagesPerDay` and `TokensPerDay` are aggregated across all projects, making it impossible to filter chart data by individual project.

## Acceptance Criteria

- [ ] `ProjectStat` struct extended with `MessagesPerDay []DataPoint` and `TokensPerDay []DataPoint`
- [ ] `ComputeStats()` computes per-project time series during aggregation
- [ ] `/api/stats` returns per-project time series in response
- [ ] Existing aggregated totals still work (backward compatible)
- [ ] Unit tests pass

## Technical Notes

- Modify `types.go` or `stats.go` to extend `ProjectStat` struct
- In `ComputeStats()`, track time series per project using a map keyed by project slug
- Performance consideration: only compute for projects with activity in the time window
- Keep aggregated `MessagesPerDay`/`TokensPerDay` at root level for "All Projects" view

## Dependencies

None - this is the foundation work item.
