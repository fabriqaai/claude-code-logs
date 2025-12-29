---
bolt-id: "009-resize-persist"
intent: "003-tree-view-sidebar"
unit: "002-resize-persist"
status: planned
created: 2025-12-29T16:00:00Z
---

# Bolt 009: Resize & Persist

## Overview

Add resizable sidebar width, localStorage persistence for UI state, and mobile hamburger menu.

## Assigned Stories

| Story | Priority | Description |
|-------|----------|-------------|
| 001-resize-handle | Must | Draggable resize functionality |
| 002-localstorage | Should | Persist width and collapse state |
| 003-mobile-menu | Should | Hamburger menu and drawer overlay |

## Scope

### In Scope
- Resize handle on sidebar edge
- Drag-to-resize with min/max bounds (200-500px)
- localStorage for width and collapsed project IDs
- Mobile hamburger button
- Slide-out drawer overlay
- Backdrop click to close

### Out of Scope
- Tree view structure (Bolt 008)
- Swipe gestures

## Implementation Approach

1. **Resize Handle**: Add HTML element, CSS positioning, JS drag events
2. **localStorage**: Save/load functions, integrate with resize and collapse
3. **Mobile**: CSS media query, hamburger button, overlay logic

## Files to Modify

| File | Changes |
|------|---------|
| `templates.go` | Add resize handle HTML/CSS/JS, localStorage logic, mobile menu |

## Acceptance Criteria

- [ ] Resize handle visible and functional
- [ ] Width clamps between 200-500px
- [ ] Width persists across page loads
- [ ] Collapsed projects persist across page loads
- [ ] Mobile: hamburger button appears < 768px
- [ ] Mobile: sidebar slides in as overlay
- [ ] Backdrop click closes overlay
- [ ] Tests pass
- [ ] Build succeeds

## Dependencies

- **Requires**: Bolt 008 (tree structure must exist for collapse state)
- **Blocks**: None (final bolt for this intent)
