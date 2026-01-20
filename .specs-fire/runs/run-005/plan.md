# Run 005 - Implementation Plan

**Status**: Approved
**Approved**: 2026-01-20
**Mode**: Confirm

---

## Work Item

**03-search-page-template** - Create Search Page Template

---

## Approach

Create a new `templates_search.go` file following the exact patterns from `templates_stats.go`:
- Same sidebar structure with project tree
- Same baseCSS + custom searchCSS
- Client-side JavaScript for search interactions

---

## File to Create

| File | Purpose |
|------|---------|
| `templates_search.go` | Search page template with HTML, CSS, and JS |

---

## Key Features

1. **Search Input** - Large, centered, auto-focused, debounced, URL-preserved
2. **Result Cards** - Session title, project, match count, expandable excerpts
3. **Pagination** - "Load more" button with hasMore state
4. **States** - Initial, loading, no results, results

---

## Acceptance Criteria

- [ ] New `templates_search.go` file with search page HTML template
- [ ] Search input field (auto-focused, placeholder text)
- [ ] Results list showing: session title, project name, match count, first excerpt
- [ ] Each result card is expandable to show all matching excerpts
- [ ] "Load more" button appears when `hasMore` is true
- [ ] Empty state when no results
- [ ] Loading state during search
- [ ] Consistent styling with existing pages
