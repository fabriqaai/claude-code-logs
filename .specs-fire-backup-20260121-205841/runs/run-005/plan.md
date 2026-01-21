# Implementation Plan: run-005

## Work Item 04: Keyboard Shortcut for Search Focus

**Mode**: confirm
**Intent**: dedicated-search-page

### Approach

Add a global `/` keyboard shortcut to all non-search pages that navigates to `/search`.

### Files to Modify

1. **templates_index.go** - Add keydown listener
2. **templates_stats.go** - Add keydown listener
3. **templates_project.go** - Add keydown listener
4. **templates_session.go** - Add keydown listener
5. **templates_shell.go** - Add keydown listener

### JavaScript Added

```javascript
document.addEventListener('keydown', function(e) {
    if (e.target.matches('input, textarea, [contenteditable]')) return;
    if (e.key === '/') {
        e.preventDefault();
        window.location.href = '/search';
    }
});
```

---

## Work Item 01: Hook Search Page into Server

**Mode**: autopilot
**Intent**: dedicated-search-page

### Approach

Integrate the existing `templates_search.go` template into the server by following the established pattern used for the stats page.

### Files to Modify

1. **server.go**
   - Add `searchTmpl *template.Template` to Server struct (line ~36)
   - Parse template in `NewServer()` (after stats template parsing)
   - Add route handling for `/search` in `handleStatic()` (after stats handling)
   - Create `renderSearchPage()` function (after `renderStatsPage()`)

### Implementation Steps

1. Add `searchTmpl` field to Server struct
2. Add template parsing in NewServer() with error handling
3. Add searchTmpl to the return struct
4. Add route handling for `/search` and `/search.html` paths
5. Create `renderSearchPage()` function following stats page pattern

### Tests

- Verify `/search` URL returns 200
- Verify page renders with correct HTML structure
- Verify search API integration works

### Acceptance Criteria

- [ ] `/search` URL serves the search page
- [ ] Page renders correctly with sidebar and project tree
- [ ] Search functionality works via existing `/api/search` API calls
- [ ] Cache strategy implemented (cache key: "search")
