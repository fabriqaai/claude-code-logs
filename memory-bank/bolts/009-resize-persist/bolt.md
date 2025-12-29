---
bolt-id: "009-resize-persist"
bolt_type: "simple-construction-bolt"
intent: "003-tree-view-sidebar"
unit: "002-resize-persist"
status: completed
created: 2025-12-29T16:00:00Z
started: 2025-12-29T18:00:00Z
completed: 2025-12-29T19:00:00Z
current_stage: test
stages_completed:
  - name: plan
    completed: 2025-12-29T18:05:00Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T18:45:00Z
    artifact: implementation-walkthrough.md
  - name: test
    completed: 2025-12-29T19:00:00Z
stories:
  - "001-resize-handle"
  - "002-localstorage"
  - "003-mobile-menu"
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

## Files Modified

| File | Changes |
|------|---------|
| `templates_css.go` | CSS for resize handle, mobile menu, backdrop |
| `templates_index.go` | Index page with full sidebar functionality |
| `templates_project.go` | Project page with full sidebar functionality |
| `templates_session.go` | Session page with full sidebar functionality |

Note: Original `templates.go` was split into 4 files for maintainability.

## Acceptance Criteria

- [x] Resize handle visible and functional
- [x] Width clamps between 200-500px
- [x] Width persists across page loads
- [x] Collapsed projects persist across page loads
- [x] Mobile: hamburger button appears < 768px
- [x] Mobile: sidebar slides in as overlay
- [x] Backdrop click closes overlay
- [x] Tests pass
- [x] Build succeeds

## Dependencies

- **Requires**: Bolt 008 (tree structure must exist for collapse state)
- **Blocks**: None (final bolt for this intent)
