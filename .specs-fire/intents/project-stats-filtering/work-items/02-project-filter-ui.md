---
id: 02-project-filter-ui
title: Add Project Filter UI
intent: project-stats-filtering
complexity: low
mode: autopilot
status: pending
depends_on: [01-stats-data-model-extension]
created: 2026-01-21
---

# Work Item: Add Project Filter UI

## Description

Add a project dropdown selector to the stats page, allowing users to select which project's stats to view. Follow existing UI patterns from the time-range filter buttons.

## Acceptance Criteria

- [ ] Dropdown/select element added to stats page filter area
- [ ] Populated dynamically from `projectStats` array
- [ ] "All Projects" as default selected option
- [ ] Styled consistently with existing filter controls
- [ ] Dropdown updates `currentProject` state variable on change

## Technical Notes

- Add to `templates_stats.go` in the filter controls section
- Use native `<select>` element or styled dropdown matching existing UI
- Store selected project slug in JavaScript variable (e.g., `currentProject`)
- Project list available from `data.projectStats` returned by API

## Dependencies

- 01-stats-data-model-extension (needs per-project data structure)
