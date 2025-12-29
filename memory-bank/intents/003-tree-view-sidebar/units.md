---
intent: 003-tree-view-sidebar
phase: inception
status: units-decomposed
updated: 2025-12-29T16:00:00Z
---

# Tree View Sidebar - Unit Decomposition

## Units Overview

This intent decomposes into **2 units** of work. The first handles core tree view functionality, the second handles resize and persistence.

## Requirement-to-Unit Mapping

| FR | Requirement | Unit |
|----|-------------|------|
| FR-1 | Tree View Structure | 001-tree-view |
| FR-2 | Expand/Collapse | 001-tree-view |
| FR-6 | Visual Enhancements | 001-tree-view |
| FR-3 | Resizable Width | 002-resize-persist |
| FR-4 | State Persistence | 002-resize-persist |
| FR-5 | Mobile Responsive | 002-resize-persist |

---

### Unit 1: 001-tree-view

**Description**: Implement hierarchical tree view with expand/collapse functionality

**Stories**:
- Story-1: Tree view markup and CSS
- Story-2: Expand/collapse JavaScript
- Story-3: Visual enhancements (chevrons, animations)

**Deliverables**:
- Modified `templates.go` - new tree view CSS classes
- Modified `templates.go` - updated HTML structure for tree nodes
- Modified `templates.go` - expand/collapse JavaScript

**Dependencies**:
- Depends on: None
- Depended by: 002-resize-persist (needs tree structure in place)

**Estimated Complexity**: M (Medium)

---

### Unit 2: 002-resize-persist

**Description**: Add resizable sidebar and localStorage persistence

**Stories**:
- Story-1: Resizable sidebar with drag handle
- Story-2: localStorage persistence for width and state
- Story-3: Mobile hamburger menu and drawer

**Deliverables**:
- Modified `templates.go` - resize handle CSS/JS
- Modified `templates.go` - localStorage integration
- Modified `templates.go` - mobile responsive overlay

**Dependencies**:
- Depends on: 001-tree-view (needs tree structure for collapse state)
- Depended by: None

**Estimated Complexity**: M (Medium)

---

## Unit Dependency Graph

```text
[001-tree-view] ──► [002-resize-persist]
```

## Execution Order

1. **Bolt 008**: 001-tree-view (tree structure and expand/collapse)
2. **Bolt 009**: 002-resize-persist (resize, persistence, mobile)
