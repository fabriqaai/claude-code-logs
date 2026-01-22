---
id: 02-stats-cost-aggregation
title: Stats Cost Aggregation
intent: cost-estimation
complexity: low
mode: autopilot
status: pending
depends_on:
  - 01-cost-calculation-logic
created: 2026-01-22
---

# Work Item: Stats Cost Aggregation

## Description

Extend the stats data model to include cost fields and update the stats computation to aggregate costs alongside existing token metrics.

## Acceptance Criteria

- [ ] `StatsData` struct extended with `TotalCost float64`
- [ ] `DataPoint` struct extended with `Cost float64` field (for daily cost)
- [ ] `ProjectStat` struct extended with `Cost float64` field
- [ ] `ComputeStats()` calculates costs during aggregation
- [ ] `/api/stats` returns cost data in response
- [ ] Existing stats functionality unchanged (backward compatible)

## Technical Notes

Extend structs in `stats.go`:
```go
type StatsData struct {
    // ... existing fields
    TotalCost float64 `json:"totalCost"`
}

type DataPoint struct {
    // ... existing fields
    Cost float64 `json:"cost"`
}

type ProjectStat struct {
    // ... existing fields
    Cost float64 `json:"cost"`
}
```

In `ComputeStats()`:
- Calculate cost for each day using `CalculateCost()`
- Sum costs per project
- Aggregate total cost

## Dependencies

- 01-cost-calculation-logic (needs `CalculateCost` function)
