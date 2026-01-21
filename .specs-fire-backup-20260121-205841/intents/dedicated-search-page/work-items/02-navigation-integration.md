---
id: 02-navigation-integration
title: Integrate Search Link in Navigation
intent: dedicated-search-page
complexity: low
mode: autopilot
status: pending
depends_on: [01-hook-search-page]
created: 2026-01-20
---

# Work Item: Integrate Search Link in Navigation

## Description

Add a search page link to the sidebar navigation across all page templates. Users should be able to navigate to search from any page via the sidebar.

## Acceptance Criteria

- [ ] Search link visible in sidebar header on all pages (index, stats, project, session, search)
- [ ] Search icon (magnifying glass) displayed with "Search" label
- [ ] Active styling (`.stats-nav-link.active`) applied when on `/search` page
- [ ] Link position consistent - after "Stats" link in nav section

## Technical Notes

Add to sidebar header section in each template file:

```html
<a href="/search" class="stats-nav-link{{if eq .CurrentPage "search"}} active{{end}}">
    <svg><!-- search icon --></svg>
    Search
</a>
```

Templates to update:
- `templates_index.go`
- `templates_stats.go`
- `templates_project.go`
- `templates_session.go`
- `templates_search.go`

Use existing search icon SVG (magnifying glass) from the codebase.

## Dependencies

- 01-hook-search-page (search page must exist to link to)
