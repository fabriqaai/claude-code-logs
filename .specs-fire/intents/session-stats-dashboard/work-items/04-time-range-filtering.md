---
id: 04-time-range-filtering
title: Time Range Filtering
intent: session-stats-dashboard
complexity: low
mode: autopilot
status: completed
depends_on: [03-stats-page-charts]
created: 2025-01-20
---

# Work Item: Time Range Filtering

## Description

Add toggle buttons to the stats page for filtering data by time range: today, this week, this month, all time.

## Acceptance Criteria

- [ ] Four toggle buttons displayed: Today, This Week, This Month, All Time
- [ ] Clicking a button filters all charts to that time range
- [ ] Active button is visually highlighted
- [ ] Default selection is "This Month"
- [ ] Filter state persists during page session
- [ ] Charts update smoothly when filter changes

## Technical Notes

- Client-side filtering (data already loaded)
- Store full dataset, filter on toggle click
- Use CSS classes for active state styling
- Consider URL param for shareable filtered views (optional)
