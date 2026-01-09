---
id: 003-visual-polish
unit: 001-tree-view
intent: 003-tree-view-sidebar
status: complete
priority: should
created: 2025-12-29T16:00:00Z
assigned_bolt: 008-tree-view
---

# Story: 003-visual-polish

## User Story

**As a** user
**I want** smooth animations and polished visuals
**So that** the tree view feels responsive and professional

## Acceptance Criteria

- [ ] **Given** expand/collapse action, **When** animating, **Then** transition takes 200ms ease-out
- [ ] **Given** chevron icon, **When** rotating, **Then** smooth 200ms rotation transition
- [ ] **Given** collapsed project, **When** viewing, **Then** session count badge visible on project row
- [ ] **Given** active session, **When** viewing, **Then** parent project auto-expanded
- [ ] **Given** tree node hover, **When** hovering, **Then** subtle background highlight
- [ ] **Given** focus on toggle, **When** focused, **Then** visible focus ring for accessibility

## Technical Notes

- Use CSS `transition: max-height 200ms ease-out, transform 200ms ease-out`
- Session count badge: `<span class="session-count">5</span>` styled inline
- Auto-expand: Check if current page's project matches, add expanded state
- Focus styles: `:focus-visible` for keyboard users only

## CSS Additions

```css
.tree-children {
  max-height: 1000px;
  overflow: hidden;
  transition: max-height 200ms ease-out;
}

.tree-node.collapsed .tree-children {
  max-height: 0;
}

.tree-chevron {
  transition: transform 200ms ease-out;
}

.tree-node.collapsed .tree-chevron {
  transform: rotate(-90deg);
}

.session-count {
  background: var(--accent-subtle);
  color: var(--accent-primary);
  font-size: 0.7rem;
  padding: 2px 6px;
  border-radius: 10px;
  margin-left: auto;
}
```

## Dependencies

### Requires
- 001-tree-markup (structure)
- 002-expand-collapse (toggle logic)

### Enables
- Unit 002 (polished UI ready for persistence)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Very long session list | max-height large enough or use JS calculation |
| Reduced motion preference | Respect `prefers-reduced-motion` |
| Many rapid toggles | Animations queue gracefully |

## Out of Scope

- Resize functionality (unit 002)
- State persistence (unit 002)
