---
id: 001-generate-html
unit: 002-generator
intent: 001-chat-log-viewer
status: complete
priority: must
created: 2025-12-29T12:00:00.000Z
assigned_bolt: 002-generator
implemented: true
---

# Story: 001-generate-html

## User Story

**As a** developer
**I want** HTML pages generated from my sessions
**So that** I can browse my chat history in a browser

## Acceptance Criteria

- [ ] **Given** parsed project data, **When** generating, **Then** index.html with project navigation is created
- [ ] **Given** a session with messages, **When** generating, **Then** session page shows conversation in Claude.ai style
- [ ] **Given** any generated page, **When** viewed offline, **Then** styling works (CSS embedded)
- [ ] **Given** any page, **When** viewing footer, **Then** shows "claude-code-logs by fabriqa.ai" linking to https://fabriqa.ai
- [ ] **Given** page opened as file://, **When** search attempted, **Then** shows "Start server for search" message
- [ ] **Given** multiple projects, **When** viewing index, **Then** 288px sidebar shows project list
- [ ] **Given** code blocks in messages, **When** rendered, **Then** syntax highlighting is applied

## Technical Notes

- Use Go `html/template` for security (auto-escaping)
- Embed CSS inline for offline viewing
- Use Claude.ai CSS variables from `ui-standards.md`:
  - Background: `hsl(48, 33.3%, 97.1%)` (cream)
  - Accent: `hsl(15, 54.2%, 51.2%)` (Claude orange)
  - Serif font for Claude, sans-serif for user
- Atomic file writes (write to temp, then rename)
- No external CDN dependencies

## Dependencies

### Requires
- 001-parser (provides `[]Project` and `Session` data)

### Enables
- 003-server (serves generated HTML)
- 004-watcher (triggers regeneration)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Message with HTML characters | Properly escaped |
| Very long message | Rendered in full, scrollable |
| Empty project (no sessions) | Show "No sessions found" |
| Code block without language hint | Default syntax highlighting |
| Message with images | Show placeholder (unsupported) |
| Unicode in project name | Display correctly |

## Out of Scope

- Live search functionality (server responsibility)
- Real-time updates (watcher responsibility)
- Dark mode (future enhancement)
