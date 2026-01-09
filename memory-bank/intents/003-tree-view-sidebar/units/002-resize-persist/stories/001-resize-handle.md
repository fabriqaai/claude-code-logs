---
id: 001-resize-handle
unit: 002-resize-persist
intent: 003-tree-view-sidebar
status: complete
priority: must
created: 2025-12-29T16:00:00Z
assigned_bolt: 009-resize-persist
---

# Story: 001-resize-handle

## User Story

**As a** user
**I want** to drag the sidebar edge to resize it
**So that** I can adjust the sidebar width to my preference

## Acceptance Criteria

- [ ] **Given** sidebar rendered, **When** viewing, **Then** resize handle visible on right edge
- [ ] **Given** mouse over handle, **When** hovering, **Then** cursor changes to `col-resize`
- [ ] **Given** mousedown on handle, **When** dragging, **Then** sidebar width follows mouse
- [ ] **Given** dragging, **When** width < 200px, **Then** width clamps to 200px
- [ ] **Given** dragging, **When** width > 500px, **Then** width clamps to 500px
- [ ] **Given** mouseup, **When** releasing, **Then** resize stops, width stays
- [ ] **Given** resize active, **When** dragging, **Then** main content adjusts margin

## Technical Notes

- Resize handle: 6px wide, transparent with visible indicator on hover
- Use mousedown/mousemove/mouseup events (touch events for mobile)
- Update CSS variable `--sidebar-width` for reactive layout
- Prevent text selection during drag with `user-select: none`

## HTML Structure

```html
<aside class="sidebar">
  <!-- existing content -->
  <div class="sidebar-resize-handle" role="separator" aria-orientation="vertical"></div>
</aside>
```

## JavaScript Pseudocode

```javascript
const handle = document.querySelector('.sidebar-resize-handle');
const sidebar = document.querySelector('.sidebar');
let isResizing = false;

handle.addEventListener('mousedown', (e) => {
  isResizing = true;
  document.body.style.cursor = 'col-resize';
  document.body.style.userSelect = 'none';
});

document.addEventListener('mousemove', (e) => {
  if (!isResizing) return;
  const width = Math.min(500, Math.max(200, e.clientX));
  document.documentElement.style.setProperty('--sidebar-width', width + 'px');
});

document.addEventListener('mouseup', () => {
  isResizing = false;
  document.body.style.cursor = '';
  document.body.style.userSelect = '';
});
```

## CSS Additions

```css
.sidebar-resize-handle {
  position: absolute;
  top: 0;
  right: 0;
  width: 6px;
  height: 100%;
  cursor: col-resize;
  background: transparent;
  transition: background-color 150ms;
}

.sidebar-resize-handle:hover {
  background: var(--border-medium);
}
```

## Dependencies

### Requires
- 001-tree-view unit (sidebar structure in place)

### Enables
- 002-localstorage (width to persist)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Mouse leaves window during drag | Continue tracking, stop on mouseup anywhere |
| Touch device | Support touchstart/touchmove/touchend |
| Iframe present | May need to handle mouseup outside |

## Out of Scope

- Persisting width (story 002)
- Mobile behavior (story 003)
