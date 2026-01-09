---
id: 003-mobile-menu
unit: 002-resize-persist
intent: 003-tree-view-sidebar
status: complete
priority: should
created: 2025-12-29T16:00:00Z
assigned_bolt: 009-resize-persist
---

# Story: 003-mobile-menu

## User Story

**As a** mobile user
**I want** a hamburger menu that reveals the sidebar
**So that** I can navigate on small screens without the sidebar taking all the space

## Acceptance Criteria

- [ ] **Given** screen width < 768px, **When** viewing page, **Then** sidebar hidden by default
- [ ] **Given** mobile view, **When** rendered, **Then** hamburger button visible in header
- [ ] **Given** hamburger clicked, **When** sidebar closed, **Then** sidebar slides in as overlay
- [ ] **Given** sidebar open as overlay, **When** clicking outside, **Then** sidebar closes
- [ ] **Given** sidebar open as overlay, **When** pressing Escape, **Then** sidebar closes
- [ ] **Given** sidebar open, **When** clicking a link, **Then** sidebar closes and navigates
- [ ] **Given** overlay visible, **When** viewing, **Then** semi-transparent backdrop shown

## Technical Notes

- Use CSS media query `@media (max-width: 768px)`
- Hamburger button: 3 horizontal lines icon
- Overlay: fixed position, z-index above content
- Backdrop: rgba(0,0,0,0.5) overlay
- Transition: slide-in from left, 200ms

## HTML Changes

```html
<!-- Add to header area -->
<button class="mobile-menu-btn" aria-label="Open navigation" aria-expanded="false">
  <svg class="hamburger-icon" viewBox="0 0 24 24">
    <path d="M3 6h18M3 12h18M3 18h18" stroke="currentColor" stroke-width="2"/>
  </svg>
</button>

<!-- Backdrop -->
<div class="sidebar-backdrop" aria-hidden="true"></div>
```

## CSS Additions

```css
.mobile-menu-btn {
  display: none;
  position: fixed;
  top: 16px;
  left: 16px;
  z-index: 200;
  padding: 8px;
  background: var(--bg-secondary);
  border: 1px solid var(--border-medium);
  border-radius: 8px;
  cursor: pointer;
}

.sidebar-backdrop {
  display: none;
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.5);
  z-index: 99;
}

@media (max-width: 768px) {
  .mobile-menu-btn {
    display: flex;
  }

  .sidebar {
    position: fixed;
    left: -100%;
    z-index: 100;
    transition: left 200ms ease-out;
  }

  .sidebar.open {
    left: 0;
  }

  .sidebar-backdrop.active {
    display: block;
  }

  .main {
    margin-left: 0;
  }
}
```

## JavaScript

```javascript
const menuBtn = document.querySelector('.mobile-menu-btn');
const sidebar = document.querySelector('.sidebar');
const backdrop = document.querySelector('.sidebar-backdrop');

function openSidebar() {
  sidebar.classList.add('open');
  backdrop.classList.add('active');
  menuBtn.setAttribute('aria-expanded', 'true');
}

function closeSidebar() {
  sidebar.classList.remove('open');
  backdrop.classList.remove('active');
  menuBtn.setAttribute('aria-expanded', 'false');
}

menuBtn?.addEventListener('click', () => {
  sidebar.classList.contains('open') ? closeSidebar() : openSidebar();
});

backdrop?.addEventListener('click', closeSidebar);

document.addEventListener('keydown', (e) => {
  if (e.key === 'Escape' && sidebar.classList.contains('open')) {
    closeSidebar();
  }
});

// Close on link click
sidebar?.querySelectorAll('a').forEach(link => {
  link.addEventListener('click', closeSidebar);
});
```

## Dependencies

### Requires
- 001-resize-handle (sidebar structure)
- 002-localstorage (state management)

### Enables
- None (final story)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Rotate device while open | Recalculate, close if now desktop |
| Resize browser while open | Same as above |
| Focus trap in overlay | Tab should cycle within sidebar |

## Out of Scope

- Swipe gestures (nice-to-have, not MVP)
- Tablet-specific layout
