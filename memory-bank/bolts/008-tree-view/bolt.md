---
bolt-id: "008-tree-view"
bolt_type: "simple-construction-bolt"
intent: "003-tree-view-sidebar"
unit: "001-tree-view"
status: completed
created: 2025-12-29T16:00:00Z
started: 2025-12-29T17:00:00Z
completed: 2025-12-29T17:30:00Z
current_stage: complete
stages_completed:
  - name: plan
    completed: 2025-12-29T17:05:00Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T17:20:00Z
    artifact: implementation-walkthrough.md
  - name: test
    completed: 2025-12-29T17:30:00Z
    artifact: test-walkthrough.md
stories:
  - "001-tree-markup"
  - "002-expand-collapse"
  - "003-visual-polish"
---

# Bolt 008: Tree View

## Overview

Implement hierarchical tree view for sidebar with expandable/collapsible project nodes.

## Assigned Stories

| Story | Priority | Description |
|-------|----------|-------------|
| 001-tree-markup | Must | Tree view HTML and CSS structure |
| 002-expand-collapse | Must | JavaScript toggle functionality |
| 003-visual-polish | Should | Animations and visual enhancements |

## Scope

### In Scope
- Tree view HTML structure with toggle buttons
- CSS for tree nodes, indentation, chevrons
- JavaScript for expand/collapse
- Smooth animations
- Expand All / Collapse All buttons
- Apply to all 3 templates

### Out of Scope
- Sidebar resizing (Bolt 009)
- State persistence (Bolt 009)
- Mobile hamburger menu (Bolt 009)

## Implementation Approach

1. **CSS First**: Add tree view styles to `baseCSS` in templates.go
2. **HTML Structure**: Update sidebar markup in all 3 templates
3. **JavaScript**: Add toggle logic at end of each template
4. **Polish**: Add animations, chevron rotation, hover states

## Files to Modify

| File | Changes |
|------|---------|
| `templates.go` | Add CSS classes, update HTML structure, add JS |

## Acceptance Criteria

- [x] Projects render as expandable tree nodes
- [x] Sessions nested under parent projects
- [x] Chevron rotates on expand/collapse
- [x] 200ms smooth animation
- [x] Expand All / Collapse All work
- [x] Tests pass
- [x] Build succeeds

## Dependencies

- **Requires**: None
- **Blocks**: Bolt 009 (needs tree structure for collapse state persistence)
