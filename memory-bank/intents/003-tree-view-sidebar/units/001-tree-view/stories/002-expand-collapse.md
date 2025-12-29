---
id: 002-expand-collapse
unit: 001-tree-view
intent: 003-tree-view-sidebar
status: ready
priority: must
created: 2025-12-29T16:00:00Z
assigned_bolt: 008-tree-view
---

# Story: 002-expand-collapse

## User Story

**As a** user with many projects
**I want** to collapse projects I'm not interested in
**So that** I can focus on the projects I'm working with

## Acceptance Criteria

- [ ] **Given** a project tree node, **When** clicking toggle button, **Then** children collapse/expand
- [ ] **Given** a project tree node, **When** clicking project name, **Then** only navigation occurs (not toggle)
- [ ] **Given** chevron icon, **When** expanded, **Then** rotates 90 degrees pointing down
- [ ] **Given** chevron icon, **When** collapsed, **Then** rotates back to pointing right
- [ ] **Given** "Expand All" button, **When** clicked, **Then** all projects expand
- [ ] **Given** "Collapse All" button, **When** clicked, **Then** all projects collapse
- [ ] **Given** aria-expanded attribute, **When** toggled, **Then** updates to match state

## Technical Notes

- Use event delegation on `.project-list` for performance
- Toggle `.collapsed` class on `.tree-node`
- CSS handles visibility via `max-height` transition
- Update `aria-expanded` for accessibility
- Add `data-project-id` attribute for persistence (unit 002)

## JavaScript Pseudocode

```javascript
document.querySelectorAll('.tree-toggle').forEach(btn => {
  btn.addEventListener('click', (e) => {
    e.preventDefault();
    const node = btn.closest('.tree-node');
    const isExpanded = node.classList.toggle('collapsed');
    btn.setAttribute('aria-expanded', !isExpanded);
  });
});

document.getElementById('expandAll')?.addEventListener('click', () => {
  document.querySelectorAll('.tree-node').forEach(n => n.classList.remove('collapsed'));
});

document.getElementById('collapseAll')?.addEventListener('click', () => {
  document.querySelectorAll('.tree-node').forEach(n => n.classList.add('collapsed'));
});
```

## Dependencies

### Requires
- 001-tree-markup (HTML structure must exist)

### Enables
- 003-visual-polish (animations depend on toggle working)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Rapid clicking | Debounce or handle gracefully |
| Keyboard Enter on toggle | Same as click |
| Project with no sessions | Toggle button hidden or disabled |

## Out of Scope

- Animation timing (story 003)
- State persistence (unit 002)
