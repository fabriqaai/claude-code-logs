---
id: 05-navigation-integration
title: Navigation Integration
intent: session-stats-dashboard
complexity: low
mode: autopilot
status: completed
depends_on: [03-stats-page-charts]
created: 2025-01-20
---

# Work Item: Navigation Integration

## Description

Add a "Stats" link to the main navigation so users can easily access the stats dashboard from any page.

## Acceptance Criteria

- [ ] "Stats" link visible in navigation/header
- [ ] Link appears on index page, project pages, and session pages
- [ ] Link styling consistent with existing navigation
- [ ] Current page indicator when on stats page

## Technical Notes

- Modify templates: `templates_index.go`, `templates_project.go`, `templates_shell.go`
- Add to sidebar or header navigation
- Use existing nav styling patterns
- Consider icon (chart/graph icon) alongside text
