---
stage: plan
bolt: 008-tree-view
created: 2025-12-29T17:00:00Z
---

## Implementation Plan: Tree View

### Objective

Transform the flat sidebar list into an interactive hierarchical tree view with expandable/collapsible project nodes. Projects will display as parent nodes with their sessions nested underneath, allowing users to collapse projects they're not interested in.

### Deliverables

- Tree view CSS classes and styles in `baseCSS`
- Updated HTML structure in all 3 templates (indexTemplate, projectIndexTemplate, sessionTemplate)
- JavaScript for expand/collapse toggle functionality
- Smooth 200ms animations for expand/collapse
- "Expand All" / "Collapse All" control buttons
- Chevron icon rotation on toggle
- Session count badge visible on collapsed projects

### Dependencies

- **None**: This is a frontend-only change to `templates.go`
- **No external libraries**: Uses vanilla CSS and JavaScript

### Technical Approach

1. **CSS First**: Add tree-specific CSS classes to `baseCSS`:
   - `.tree-node` - Container for project + children
   - `.tree-node-header` - Clickable header row with toggle button
   - `.tree-toggle` - Button containing chevron
   - `.tree-chevron` - SVG chevron icon
   - `.tree-children` - Container for nested sessions
   - Transitions for smooth height animation and chevron rotation

2. **HTML Structure**: Update sidebar markup in all 3 templates:
   - Wrap each project in `.tree-node`
   - Add `.tree-node-header` with toggle button and project link
   - Move session list into `.tree-children` container
   - Add "Expand All" / "Collapse All" buttons above project list

3. **JavaScript**: Add toggle logic to each template:
   - Event listeners on `.tree-toggle` buttons
   - Toggle `.collapsed` class on `.tree-node`
   - Update `aria-expanded` for accessibility
   - Expand All / Collapse All handlers

4. **Visual Polish**:
   - Chevron rotation (90 degrees) on collapse
   - 200ms ease-out transitions
   - Session count badge styling
   - Hover states on tree nodes
   - Focus ring for accessibility

### File Changes

| File | Changes |
|------|---------|
| `templates.go` | Add CSS to baseCSS, update HTML structure in all 3 templates, add JS |

### Acceptance Criteria

- [ ] Projects render as expandable tree nodes with toggle buttons
- [ ] Sessions nested under parent projects with proper indentation
- [ ] Clicking toggle button expands/collapses children (clicking project name navigates)
- [ ] Chevron rotates 90 degrees on expand/collapse
- [ ] 200ms smooth animation for expand/collapse
- [ ] "Expand All" / "Collapse All" buttons work correctly
- [ ] `aria-expanded` attribute updates correctly
- [ ] Projects with no sessions show without toggle (or toggle disabled)
- [ ] Focus ring visible on toggle button for keyboard users
- [ ] Works in all 3 templates (index, project, session)
- [ ] `go build` succeeds
- [ ] `go test` passes

### Edge Cases

| Scenario | Handling |
|----------|----------|
| Project with 0 sessions | Hide toggle button, show project as leaf node |
| Very long project name | Truncate with ellipsis (existing behavior) |
| 50+ projects | All render without performance issues |
| Rapid clicking | CSS transitions handle gracefully |
| Keyboard navigation | Enter key triggers toggle, focus states visible |
| Reduced motion | Respect `prefers-reduced-motion` media query |
