---
id: 03-stats-page-cost-ui
title: Stats Page Cost UI
intent: cost-estimation
complexity: medium
mode: confirm
status: pending
depends_on:
  - 02-stats-cost-aggregation
created: 2026-01-22
---

# Work Item: Stats Page Cost UI

## Description

Add cost display elements to the stats page: a cost summary card showing total estimated cost, and integrate cost data into the existing charts or add a dedicated cost chart.

## Acceptance Criteria

- [ ] Cost summary card displaying total estimated cost
- [ ] Cost formatted as currency ($X.XX)
- [ ] Cost updates when time-range filter changes
- [ ] Cost updates when project filter changes
- [ ] Cost by project visible (in project activity section or separate)
- [ ] Optional: Cost over time chart line

## Technical Notes

In `templates_stats.go`:

1. Add cost summary card:
```html
<div class="stat-card">
    <div class="stat-value" id="total-cost">$0.00</div>
    <div class="stat-label">Estimated Cost</div>
</div>
```

2. Update JavaScript to:
   - Read `totalCost` from API response
   - Format as currency: `cost.toFixed(2)`
   - Update on filter changes (reuse existing filter logic)

3. Project costs:
   - Show in project activity breakdown
   - Format: "project-name: X messages, Y tokens, $Z.ZZ"

4. Optional enhancement:
   - Add cost line to tokens chart (secondary Y-axis)
   - Or separate "Cost per Day" chart

## Dependencies

- 02-stats-cost-aggregation (needs cost data in API response)
