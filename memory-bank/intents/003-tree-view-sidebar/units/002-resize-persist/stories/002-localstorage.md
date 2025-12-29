---
id: 002-localstorage
unit: 002-resize-persist
intent: 003-tree-view-sidebar
status: ready
priority: should
created: 2025-12-29T16:00:00Z
assigned_bolt: 009-resize-persist
---

# Story: 002-localstorage

## User Story

**As a** user
**I want** my sidebar settings to persist across page loads
**So that** I don't have to reconfigure the sidebar every time

## Acceptance Criteria

- [ ] **Given** sidebar width changed, **When** mouseup, **Then** width saved to localStorage
- [ ] **Given** project collapsed, **When** toggled, **Then** collapsed state saved to localStorage
- [ ] **Given** page load, **When** localStorage has width, **Then** sidebar width restored
- [ ] **Given** page load, **When** localStorage has collapsed IDs, **Then** projects collapsed
- [ ] **Given** localStorage unavailable, **When** page loads, **Then** default values used (no error)
- [ ] **Given** invalid localStorage data, **When** parsing, **Then** reset to defaults

## Technical Notes

- localStorage key: `claude-logs-sidebar`
- Data structure: `{ width: number, collapsed: string[] }`
- Project IDs: Use `data-project-id` attribute (slugified path)
- Load state early to prevent flash of expanded/collapsed

## JavaScript Implementation

```javascript
const STORAGE_KEY = 'claude-logs-sidebar';

function loadSidebarState() {
  try {
    const data = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');

    // Restore width
    if (data.width && data.width >= 200 && data.width <= 500) {
      document.documentElement.style.setProperty('--sidebar-width', data.width + 'px');
    }

    // Restore collapsed state
    if (Array.isArray(data.collapsed)) {
      data.collapsed.forEach(id => {
        const node = document.querySelector(`[data-project-id="${id}"]`);
        if (node) node.classList.add('collapsed');
      });
    }
  } catch (e) {
    console.warn('Failed to load sidebar state:', e);
  }
}

function saveSidebarState() {
  try {
    const width = parseInt(getComputedStyle(document.documentElement)
      .getPropertyValue('--sidebar-width'));
    const collapsed = Array.from(document.querySelectorAll('.tree-node.collapsed'))
      .map(n => n.dataset.projectId)
      .filter(Boolean);

    localStorage.setItem(STORAGE_KEY, JSON.stringify({ width, collapsed }));
  } catch (e) {
    console.warn('Failed to save sidebar state:', e);
  }
}

// Load on DOMContentLoaded
document.addEventListener('DOMContentLoaded', loadSidebarState);

// Save on width change (debounced)
// Save on collapse toggle
```

## Dependencies

### Requires
- 001-resize-handle (width changes to persist)
- Unit 001 tree-view (collapse state to persist)

### Enables
- 003-mobile-menu (state management in place)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| localStorage full | Fail silently, state not saved |
| Project renamed/removed | Orphan IDs ignored on load |
| Multiple tabs | Last write wins (acceptable) |
| Private browsing | May not persist, handle gracefully |

## Out of Scope

- Cross-device sync
- Versioned storage migration
