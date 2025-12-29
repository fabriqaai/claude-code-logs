---
id: 001-tree-markup
unit: 001-tree-view
intent: 003-tree-view-sidebar
status: ready
priority: must
created: 2025-12-29T16:00:00Z
assigned_bolt: 008-tree-view
---

# Story: 001-tree-markup

## User Story

**As a** user browsing chat logs
**I want** projects displayed as a hierarchical tree
**So that** I can see the structure and relationship between projects and sessions

## Acceptance Criteria

- [ ] **Given** the sidebar renders, **When** viewing projects, **Then** each project has a tree-node wrapper with toggle button
- [ ] **Given** a project with sessions, **When** rendered, **Then** sessions are nested children with proper indentation
- [ ] **Given** CSS loaded, **When** viewing tree, **Then** indentation is 16px per level
- [ ] **Given** a collapsed project, **When** rendered, **Then** sessions are hidden via CSS (max-height: 0)
- [ ] **Given** expand/collapse header, **When** rendered, **Then** "Expand All" and "Collapse All" buttons appear

## Technical Notes

- Add `.tree-node`, `.tree-toggle`, `.tree-children` CSS classes
- Use CSS `max-height` transition for smooth expand/collapse
- Chevron icon using CSS transform or SVG
- Update all 3 templates: indexTemplate, projectIndexTemplate, sessionTemplate

## HTML Structure

```html
<ul class="project-list tree">
  <li class="tree-node">
    <div class="tree-node-header">
      <button class="tree-toggle" aria-expanded="true">
        <svg class="tree-chevron">...</svg>
      </button>
      <a href="..." class="project-link">
        <span class="project-name">Project Name</span>
        <span class="project-meta">5 sessions</span>
      </a>
    </div>
    <ul class="tree-children session-list">
      <li class="session-item">...</li>
    </ul>
  </li>
</ul>
```

## Dependencies

### Requires
- None (first story)

### Enables
- 002-expand-collapse (needs markup to attach events)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Project with 0 sessions | Show project without expand toggle |
| Very long project name | Truncate with ellipsis |
| 50+ projects | All render without performance issues |

## Out of Scope

- JavaScript toggle logic (story 002)
- Animation polish (story 003)
- State persistence (unit 002)
