---
unit: 001-tree-view
intent: 003-tree-view-sidebar
status: complete
created: 2025-12-29T16:00:00Z
---

# Unit Brief: 001-tree-view

## Purpose

Transform the flat sidebar list into an interactive hierarchical tree view with expandable/collapsible project nodes.

## Scope

### In Scope
- Tree view HTML structure with proper nesting
- CSS for tree nodes, indentation, and visual hierarchy
- Chevron icons for expand/collapse indicators
- JavaScript for toggle functionality
- Smooth CSS animations for expand/collapse
- "Expand All" / "Collapse All" controls

### Out of Scope
- Sidebar resizing (Unit 002)
- State persistence (Unit 002)
- Mobile hamburger menu (Unit 002)

## Technical Approach

1. **HTML Structure**: Wrap projects in tree-node containers with toggle buttons
2. **CSS**: Add tree-specific classes for indentation, transitions, icons
3. **JavaScript**: Event listeners for toggle, chevron rotation, height animation

## Dependencies

- **Requires**: None
- **Enables**: 002-resize-persist

## Acceptance Criteria

- [ ] Projects displayed as collapsible tree nodes
- [ ] Sessions nested under their parent project
- [ ] Chevron rotates 90 degrees on expand
- [ ] Smooth 200ms animation on expand/collapse
- [ ] Expand All / Collapse All buttons functional
- [ ] Works in all 3 templates (index, project, session)

## Stories

1. `001-tree-markup.md` - Tree view HTML and CSS structure
2. `002-expand-collapse.md` - JavaScript toggle functionality
3. `003-visual-polish.md` - Animations and visual enhancements
