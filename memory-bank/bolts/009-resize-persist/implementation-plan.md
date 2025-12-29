---
stage: plan
bolt: 009-resize-persist
created: 2025-12-29T18:00:00Z
---

## Implementation Plan: resize-persist

### Objective

Add resizable sidebar width with drag handle, persist UI state (width, collapsed projects) to localStorage, and provide mobile-responsive hamburger menu for small screens.

### Deliverables

- Resize handle element on sidebar right edge with drag functionality
- localStorage persistence for sidebar width and collapsed project IDs
- Mobile hamburger menu button with slide-out drawer overlay
- Backdrop overlay with click-to-close functionality

### Dependencies

| Dependency | Type | Why Needed |
|------------|------|------------|
| Bolt 008 tree structure | Internal | Collapse state requires tree nodes with `data-project-id` |
| CSS variable `--sidebar-width` | Internal | Already used in layout, will extend for dynamic resize |
| No external libraries | N/A | All features use vanilla JS and CSS |

### Technical Approach

**1. Resize Handle (Story 001)**
- Add `<div class="sidebar-resize-handle">` as last child of sidebar
- Position absolute, right edge, full height, 6px wide
- Use `mousedown` → `mousemove` → `mouseup` event chain
- Update CSS variable `--sidebar-width` during drag
- Clamp between 200px and 500px
- Add touch event support for mobile

**2. localStorage Persistence (Story 002)**
- Storage key: `claude-logs-sidebar`
- Data shape: `{ width: number, collapsed: string[] }`
- Load on `DOMContentLoaded` (before paint ideally)
- Save width on resize `mouseup`
- Save collapsed state on toggle (from Bolt 008)
- Graceful fallback if localStorage unavailable

**3. Mobile Menu (Story 003)**
- CSS media query: `@media (max-width: 768px)`
- Hamburger button: fixed position, z-index above content
- Sidebar: fixed position, `left: -100%` by default, `left: 0` when open
- Backdrop: fixed, `rgba(0,0,0,0.5)`, z-index between content and sidebar
- Close triggers: backdrop click, Escape key, link click

### File Changes

| File | Changes |
|------|---------|
| `templates.go` | Add resize handle HTML, mobile menu HTML, all CSS, all JavaScript |

### Acceptance Criteria

- [ ] Resize handle visible on sidebar right edge
- [ ] Cursor changes to `col-resize` on hover
- [ ] Sidebar width follows mouse during drag
- [ ] Width clamps to 200-500px range
- [ ] Width persists across page loads via localStorage
- [ ] Collapsed project IDs persist across page loads
- [ ] On mobile (< 768px): hamburger button visible
- [ ] Hamburger click opens sidebar as overlay
- [ ] Backdrop visible behind open sidebar
- [ ] Backdrop click closes sidebar
- [ ] Escape key closes sidebar
- [ ] Link click in sidebar closes overlay
- [ ] Tests pass
- [ ] Build succeeds

### Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| localStorage quota exceeded | Catch error, continue without persistence |
| Touch events differ from mouse | Support both touchstart/touchmove/touchend |
| CSS variable not reactive | Test in all template sections |

### Estimated Effort

Stage 2 (Implement): ~60-90 minutes
Stage 3 (Test): ~30-45 minutes
