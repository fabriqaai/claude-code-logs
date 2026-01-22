---
id: 03-client-side-filtering
title: Implement Client-Side Project Filtering
intent: project-stats-filtering
complexity: medium
mode: confirm
status: completed
depends_on: [01-stats-data-model-extension, 02-project-filter-ui]
created: 2026-01-21
---

# Work Item: Implement Client-Side Project Filtering

## Description

Wire up the project selector to filter all stats displays. When a project is selected, summary cards, time series charts, and project activity should all reflect the filtered data.

## Acceptance Criteria

- [ ] Summary cards (Total Messages, Total Tokens) update when project selected
- [ ] Messages per Day chart filters to show only selected project's data
- [ ] Tokens per Day chart filters to show only selected project's data
- [ ] Project filter combines with time-range filter (both applied together)
- [ ] "All Projects" shows aggregated data (current behavior)
- [ ] Project Activity chart highlights or filters to selected project

## Technical Notes

- Extend `filterDataByRange()` to also handle project filtering
- Create `filterDataByProject()` or combine into unified filter function
- When project selected:
  - Use `projectStats[selectedProject].messagesPerDay` for charts
  - Use `projectStats[selectedProject].totalMessages` for summary cards
- When "All Projects" selected:
  - Use root-level aggregated data (existing behavior)
- Update chart redraw logic to handle both filters

## Dependencies

- 01-stats-data-model-extension (per-project time series data)
- 02-project-filter-ui (dropdown selector)
