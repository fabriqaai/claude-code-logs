---
story: 003-copy-button
unit: 002-client-renderer
intent: 004-markdown-first-generation
priority: Should
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Copy as Markdown Button

## User Story

As a user, I want to copy the current session as Markdown to my clipboard, so that I can quickly paste it into other tools (LLMs, editors, etc.).

## Acceptance Criteria

- [ ] Copy button visible in session view UI
- [ ] Clicking button copies MD content to clipboard
- [ ] Copied content includes YAML frontmatter
- [ ] Visual feedback shown after copy (toast/icon change)
- [ ] Button has appropriate icon and tooltip
- [ ] Works across browsers (Chrome, Firefox, Safari)

## Technical Notes

- Use navigator.clipboard.writeText() API
- Fall back to document.execCommand('copy') for older browsers
- Show brief "Copied!" feedback (CSS animation or toast)
- MD content is already fetched for rendering, reuse it

## Dependencies

- Story 001 (MD content must be fetched first)
- Clipboard API browser support

## Estimation

Complexity: S
