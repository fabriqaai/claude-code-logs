---
unit: 002-generator
intent: 001-chat-log-viewer
phase: inception
status: draft
---

# Unit Brief: generator

## Purpose

Generate static HTML pages from parsed session data with Claude.ai visual styling.

## Responsibility

- Generate main index page with project navigation
- Generate session pages with conversation view
- Apply Claude.ai visual styling (see `ui-standards.md`)
- Embed CSS for offline viewing
- Include footer branding (fabriqa.ai)
- Show "Start server for search" message in static mode

## Assigned Requirements

- **FR-3**: HTML Generation
- **FR-7**: Static Fallback UI

## Key Entities

### Templates

- `index.html` - Project list with sidebar navigation
- `session.html` - Conversation view for a single session
- `styles.css` - Claude.ai visual styling (embedded)

### Output Structure

```
~/.claude-logs/
├── index.html           # Project navigation
├── projects/
│   └── {project-name}/
│       ├── index.html   # Session list for project
│       └── sessions/
│           └── {session-id}.html
└── assets/
    └── styles.css       # Embedded in HTML
```

## Key Operations

1. **GenerateAll(projects, outputDir)** - Generate complete site
2. **GenerateIndex(projects)** - Generate main index
3. **GenerateProject(project)** - Generate project page with session list
4. **GenerateSession(session)** - Generate session conversation view
5. **RenderMessage(message)** - Render single message with styling

## Dependencies

- **001-parser**: Provides `[]Project` and `Session` data

## Interface

- Input: Parsed project/session data from parser
- Output: HTML files written to `outputDir`

## Technical Constraints

- Use Go `html/template` for security (auto-escaping)
- Embed CSS inline for offline viewing
- Responsive design (mobile-friendly)
- No external CDN dependencies
- Atomic file writes (write to temp, then rename)

## Styling Requirements

Reference: `ui-standards.md`

- Warm cream backgrounds (`hsl(48, 33.3%, 97.1%)`)
- Serif font for Claude responses
- Sans-serif for user messages
- 0.5px subtle borders
- Claude orange accent (`hsl(15, 54.2%, 51.2%)`)
- Footer: "claude-code-logs by fabriqa.ai"

## Success Criteria

- [ ] Generates valid HTML5 pages
- [ ] Matches Claude.ai visual style
- [ ] Works offline (no external dependencies)
- [ ] Responsive on mobile
- [ ] Shows "Start server for search" when opened as file://
- [ ] Footer branding visible on all pages
- [ ] Syntax highlighting for code blocks

---

## Story Summary

- **Total Stories**: 1
- **Must Have**: 1
- **Should Have**: 0
- **Could Have**: 0

### Stories

- [ ] **001-generate-html**: Generate HTML pages - Must - Planned
