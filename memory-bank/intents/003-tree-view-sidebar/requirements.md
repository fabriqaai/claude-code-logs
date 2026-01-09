---
intent: 003-tree-view-sidebar
phase: complete
status: complete
created: 2025-12-29T16:00:00Z
updated: 2026-01-09T12:00:00Z
---

# Requirements: Tree View Sidebar

## Intent Overview

Transform the left sidebar from a flat list into an interactive tree view with expandable/collapsible nodes and resizable width. This improves navigation for users with many projects and sessions.

## Business Goals

| Goal | Success Metric | Priority |
|------|----------------|----------|
| Improve navigation | Users can collapse projects to reduce visual clutter | Must |
| Flexible layout | Sidebar width adjustable to user preference | Must |
| Persistent state | Collapsed/expanded state remembered across sessions | Should |
| Mobile friendly | Sidebar works on mobile devices | Should |

---

## Functional Requirements

### FR-1: Tree View Structure
- **Description**: Display projects and sessions in a hierarchical tree structure
- **Acceptance Criteria**:
  - Projects shown as parent nodes with expand/collapse icons
  - Sessions shown as child nodes under their project
  - Visual indentation indicates hierarchy
  - Expand/collapse icons (chevron or +/-) clearly visible
- **Priority**: Must
- **Related Stories**: TBD

### FR-2: Expand/Collapse Functionality
- **Description**: Allow users to expand and collapse project nodes
- **Acceptance Criteria**:
  - Clicking project row or chevron toggles expand/collapse
  - Collapsed projects hide their session list
  - Expanded projects show their session list
  - Smooth animation on expand/collapse (200-300ms)
  - "Expand All" / "Collapse All" buttons in sidebar header
- **Priority**: Must
- **Related Stories**: TBD

### FR-3: Resizable Sidebar Width
- **Description**: Allow users to drag to resize sidebar width
- **Acceptance Criteria**:
  - Drag handle visible on sidebar right edge
  - Cursor changes to resize cursor on hover
  - Minimum width: 200px
  - Maximum width: 500px
  - Default width: 280px (current)
  - Smooth resize without content jumping
- **Priority**: Must
- **Related Stories**: TBD

### FR-4: State Persistence
- **Description**: Remember sidebar state across page loads
- **Acceptance Criteria**:
  - Sidebar width saved to localStorage
  - Expand/collapse state of each project saved
  - State restored on page load
  - Graceful fallback if localStorage unavailable
- **Priority**: Should
- **Related Stories**: TBD

### FR-5: Mobile Responsive Behavior
- **Description**: Sidebar adapts to mobile screen sizes
- **Acceptance Criteria**:
  - On screens < 768px: sidebar collapses to hamburger menu
  - Hamburger button visible in header
  - Tapping hamburger shows sidebar as overlay/drawer
  - Tapping outside drawer closes it
  - Swipe gesture to close drawer (nice to have)
- **Priority**: Should
- **Related Stories**: TBD

### FR-6: Visual Enhancements
- **Description**: Improve tree view visual clarity
- **Acceptance Criteria**:
  - Chevron rotates 90 degrees when expanded
  - Active project/session highlighted
  - Hover states on all clickable elements
  - Session count badge on collapsed projects
- **Priority**: Should
- **Related Stories**: TBD

---

## Non-Functional Requirements

### Performance
| Requirement | Metric | Target |
|-------------|--------|--------|
| Expand/collapse animation | Frame rate | 60fps |
| Resize responsiveness | Lag | < 16ms |
| Initial render | Time to interactive | < 100ms |

### Accessibility
| Requirement | Metric | Target |
|-------------|--------|--------|
| Keyboard navigation | Arrow keys work | Yes |
| Screen reader support | ARIA labels | Yes |
| Focus indicators | Visible focus ring | Yes |

### Browser Support
| Browser | Minimum Version |
|---------|-----------------|
| Chrome | 90+ |
| Firefox | 90+ |
| Safari | 14+ |
| Edge | 90+ |

---

## Constraints

### Technical Constraints

**Project-wide standards**:
- Go templates for HTML generation
- Vanilla JavaScript (no frameworks)
- CSS embedded in templates

**Intent-specific constraints**:
- Must work without JavaScript for basic viewing (progressive enhancement)
- No external dependencies (self-contained)
- Must maintain search functionality

### Business Constraints
- Non-breaking change (enhances existing UI)

---

## Assumptions

| Assumption | Risk if Invalid | Mitigation |
|------------|-----------------|------------|
| Users have modern browsers | Older browsers may break | CSS/JS feature detection |
| localStorage available | State won't persist | Fallback to defaults |
| Touch events work consistently | Mobile may have issues | Test on real devices |

---

## Open Questions

| Question | Owner | Due Date | Resolution |
|----------|-------|----------|------------|
| Should collapsed state be per-page or global? | TBD | Before construction | Global across all pages |
| Double-click to expand all children? | TBD | Before construction | Pending |
