---
id: 05-remove-overlay-update-shortcuts
title: Remove Overlay & Update Shortcuts
intent: dedicated-search-page
complexity: medium
mode: confirm
status: pending
depends_on:
  - 04-search-page-route
created: 2026-01-20
---

# Work Item: Remove Overlay & Update Shortcuts

## Description

Remove the old overlay-based search UI from existing templates and update keyboard shortcuts to navigate to the new search page.

## Acceptance Criteria

- [ ] Search overlay HTML/CSS/JS removed from `templates_index.go`
- [ ] Search overlay removed from `templates_shell.go` (session pages)
- [ ] Search overlay removed from `templates_project.go`
- [ ] Search overlay removed from `templates_stats.go`
- [ ] `/` keyboard shortcut navigates to `/search` page (preserving current query if any)
- [ ] Search icon in header links to `/search` instead of opening overlay
- [ ] No JavaScript errors after removal

## Technical Notes

- Search for `search-overlay`, `searchOverlay`, `openSearch`, `closeSearch` in templates
- Keep the search icon in header but change onclick to navigation
- Update keyboard listener: `if (e.key === '/') window.location.href = '/search'`
- Test on index, project, session, and stats pages

## Files to Modify

- `templates_index.go`
- `templates_shell.go`
- `templates_project.go`
- `templates_stats.go`
