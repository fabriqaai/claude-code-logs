---
id: 03-stats-page-charts
title: Stats Page with Charts
intent: session-stats-dashboard
complexity: medium
mode: confirm
status: completed
depends_on: [02-stats-api-endpoint]
created: 2025-01-20
---

# Work Item: Stats Page with Charts

## Description

Create a new `/stats` page with HTML template and client-side JavaScript charting to visualize usage analytics.

## Acceptance Criteria

- [ ] `/stats` route serves stats page
- [ ] Page displays 4 charts: messages over time, token usage, project activity, session lengths
- [ ] Charts render correctly with data from `/api/stats`
- [ ] Page follows existing UI style (consistent with index/session pages)
- [ ] Charts are responsive (work on mobile)
- [ ] Loading state shown while fetching data

## Technical Notes

- Create `templates_stats.go` with stats page template
- Use Chart.js via CDN (lightweight, no build step)
- Add `renderStatsPage` handler in `server.go`
- Chart types: line chart for time series, bar chart for project comparison
- Match existing CSS variables and color scheme
