---
intent: 003-tree-view-sidebar
phase: inception
status: context-defined
updated: 2025-12-29T16:00:00Z
---

# System Context: Tree View Sidebar

## System Boundary

```
┌─────────────────────────────────────────────────────────────────┐
│                        Generated HTML Pages                      │
│  ┌─────────────────────────────────────────────────────────────┐ │
│  │                        Browser                               │ │
│  │  ┌─────────────────┐  ┌────────────────────────────────────┐ │ │
│  │  │    Sidebar      │  │          Main Content              │ │ │
│  │  │  ┌───────────┐  │  │                                    │ │ │
│  │  │  │ Tree View │◄─┼──┼── This Intent                      │ │ │
│  │  │  │ - Projects│  │  │                                    │ │ │
│  │  │  │   └Sessions│ │  │  ┌─────────────────────────────┐   │ │ │
│  │  │  ├───────────┤  │  │  │   Messages / Session View   │   │ │ │
│  │  │  │ Resize    │  │  │  │                             │   │ │ │
│  │  │  │ Handle    │  │  │  └─────────────────────────────┘   │ │ │
│  │  │  └───────────┘  │  │                                    │ │ │
│  │  └─────────────────┘  └────────────────────────────────────┘ │ │
│  │                                                               │ │
│  │  localStorage: sidebar-width, collapsed-projects              │ │
│  └───────────────────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────────────────────┘
```

## Actors

| Actor | Type | Description |
|-------|------|-------------|
| User | Human | Views and navigates chat logs |
| Browser | System | Renders HTML, executes JavaScript, stores localStorage |
| Generator | System | Produces HTML files with embedded CSS/JS (Intent 001) |

## External Systems

| System | Integration | Data Exchanged |
|--------|-------------|----------------|
| localStorage | Browser API | Sidebar width, collapsed project IDs |
| CSS Variables | Browser | Theme colors, transitions |

## Data Flow

```
[Page Load]
     │
     ▼
[Read localStorage] ──► [sidebar-width, collapsed-projects]
     │
     ▼
[Apply Initial State]
     │
     ├──► [Set sidebar width from stored value]
     └──► [Collapse projects from stored list]

[User Interaction]
     │
     ├──► [Click Project] ──► [Toggle expand/collapse] ──► [Save to localStorage]
     │
     ├──► [Drag Resize Handle] ──► [Update width] ──► [Save to localStorage]
     │
     └──► [Click Hamburger (mobile)] ──► [Toggle sidebar overlay]
```

## Component Architecture

```
templates.go
├── baseCSS
│   ├── .sidebar (existing)
│   ├── .sidebar-resize-handle (new)
│   ├── .tree-node (new)
│   ├── .tree-node-toggle (new)
│   └── .mobile-overlay (new)
│
├── indexTemplate
│   └── Sidebar with tree view markup
│
├── projectIndexTemplate
│   └── Sidebar with tree view markup
│
└── sessionTemplate
    └── Sidebar with tree view markup

+ Embedded JavaScript for:
  - Tree expand/collapse
  - Resize handling
  - localStorage persistence
  - Mobile hamburger menu
```

## Integration Points

| Component | File | Change Type |
|-----------|------|-------------|
| CSS Styles | templates.go:baseCSS | Add tree view styles |
| HTML Structure | templates.go:*Template | Modify sidebar markup |
| JavaScript | templates.go:*Template | Add interaction logic |
| Generator | generator.go | No changes needed |

## Constraints from Context

1. **No external dependencies** - All CSS/JS must be embedded
2. **Progressive enhancement** - Basic navigation works without JS
3. **Template duplication** - Changes must be applied to all 3 templates
4. **Offline-first** - Must work when served as static files

## Security Considerations

- localStorage only stores UI preferences, no sensitive data
- No external scripts or resources loaded
- All user input is for UI state only (no injection risk)
