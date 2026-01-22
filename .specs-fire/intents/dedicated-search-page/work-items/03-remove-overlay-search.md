---
id: 03-remove-overlay-search
title: Remove Old Overlay Search
intent: dedicated-search-page
complexity: medium
mode: confirm
status: completed
depends_on: [02-navigation-integration]
created: 2026-01-20
---

# Work Item: Remove Old Overlay Search

## Description

Remove the deprecated overlay/dropdown search functionality from other pages. Now that the dedicated search page exists and is linked in navigation, the inline search overlays are redundant.

## Acceptance Criteria

- [ ] No overlay search input remains in sidebar on any page
- [ ] Associated JavaScript for overlay search removed
- [ ] CSS for overlay components removed (if any orphaned)
- [ ] Pages still render and function correctly
- [ ] No console errors or dead code references

## Technical Notes

Components to remove from templates:
- Mini search input in sidebar header
- Search dropdown results container
- JavaScript handling: fetch, debounce, result rendering for overlay
- Any CSS specific to overlay search

Templates to clean:
- `templates_index.go`
- `templates_stats.go`
- `templates_project.go`

**Checkpoint**: Review changes before committing to ensure no regressions.

## Dependencies

- 02-navigation-integration (users need alternative navigation to search before removing overlay)
