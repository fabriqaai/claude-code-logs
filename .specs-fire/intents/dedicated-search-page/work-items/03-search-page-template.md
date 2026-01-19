---
id: 03-search-page-template
title: Create Search Page Template
intent: dedicated-search-page
complexity: medium
mode: confirm
status: pending
depends_on:
  - 01-search-api-pagination
  - 02-relevance-scoring
created: 2026-01-20
---

# Work Item: Create Search Page Template

## Description

Create a new Go template for the dedicated search page. The page should provide a clean, focused search experience with expandable results and "Load more" pagination.

## Acceptance Criteria

- [ ] New `templates_search.go` file with search page HTML template
- [ ] Search input field (auto-focused, placeholder text)
- [ ] Results list showing: session title, project name, match count, first excerpt
- [ ] Each result card is expandable to show all matching excerpts with highlighting
- [ ] "Load more" button appears when `hasMore` is true
- [ ] Empty state when no results
- [ ] Loading state during search
- [ ] Consistent styling with existing pages (stats, sessions)

## Technical Notes

- Use same CSS variables and layout patterns as `templates_stats.go`
- Client-side JS handles: search submission, result rendering, expand/collapse, load more
- Debounce search input (300ms) for live search feel
- Preserve query in URL (`/search?q=term`) for shareability
- Mobile responsive: stack results vertically, full-width search input

## Files to Create

- `templates_search.go` - Template with HTML, CSS, and JS
