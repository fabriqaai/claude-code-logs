---
id: 04-keyboard-shortcut
title: Keyboard Shortcut for Search Focus
intent: dedicated-search-page
complexity: medium
mode: confirm
status: pending
depends_on: [01-hook-search-page]
created: 2026-01-20
---

# Work Item: Keyboard Shortcut for Search Focus

## Description

Implement a global `/` keyboard shortcut that navigates users to the search page from any page in the application. This is a common UX pattern (GitHub, Slack, etc.) for quick access to search.

## Acceptance Criteria

- [ ] Pressing `/` on any page navigates to `/search`
- [ ] If already on search page, `/` focuses the search input
- [ ] Shortcut does NOT trigger when user is typing in input/textarea/contenteditable
- [ ] Works consistently across index, stats, project, session pages
- [ ] No conflicts with existing keyboard shortcuts

## Technical Notes

Add global keydown listener to base template JavaScript:

```javascript
document.addEventListener('keydown', function(e) {
    // Skip if typing in an input
    if (e.target.matches('input, textarea, [contenteditable]')) return;

    if (e.key === '/') {
        e.preventDefault();
        if (window.location.pathname === '/search') {
            document.querySelector('.search-input')?.focus();
        } else {
            window.location.href = '/search';
        }
    }
});
```

Add to:
- `templates_index.go`
- `templates_stats.go`
- `templates_project.go`
- `templates_session.go`

The search page (`templates_search.go`) may already have this - verify and ensure consistency.

**Checkpoint**: Test across all pages before completing.

## Dependencies

- 01-hook-search-page (search page must exist to navigate to)
