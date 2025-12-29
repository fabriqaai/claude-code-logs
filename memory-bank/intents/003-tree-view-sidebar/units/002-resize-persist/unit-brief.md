---
unit: 002-resize-persist
intent: 003-tree-view-sidebar
status: ready
created: 2025-12-29T16:00:00Z
---

# Unit Brief: 002-resize-persist

## Purpose

Add resizable sidebar width with drag handle and persist UI state (width, collapsed projects) to localStorage. Include mobile-responsive hamburger menu.

## Scope

### In Scope
- Resize handle on sidebar right edge
- Drag-to-resize functionality with min/max bounds
- localStorage for sidebar width
- localStorage for collapsed project IDs
- Mobile hamburger menu button
- Mobile slide-out drawer overlay

### Out of Scope
- Tree view structure (Unit 001)
- Expand/collapse logic (Unit 001 - we just persist state)

## Technical Approach

1. **Resize Handle**: Absolute-positioned div on sidebar edge, mousedown/mousemove/mouseup events
2. **localStorage**: JSON storage for `{width: number, collapsed: string[]}`
3. **Mobile**: CSS media query triggers hamburger mode, JS toggles overlay

## Dependencies

- **Requires**: 001-tree-view (collapse state needs tree structure)
- **Enables**: None (final unit)

## Acceptance Criteria

- [ ] Drag handle visible on sidebar right edge
- [ ] Cursor changes to `col-resize` on hover
- [ ] Sidebar resizes between 200px and 500px
- [ ] Width persists across page loads
- [ ] Collapsed project IDs persist across page loads
- [ ] On mobile (<768px), sidebar becomes hamburger menu
- [ ] Hamburger click shows sidebar as overlay
- [ ] Clicking outside overlay closes it

## Stories

1. `001-resize-handle.md` - Draggable resize functionality
2. `002-localstorage.md` - Persist width and collapse state
3. `003-mobile-menu.md` - Hamburger menu and drawer overlay
