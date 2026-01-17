---
story: 002-download-button
unit: 002-client-renderer
intent: 004-markdown-first-generation
priority: Should
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Download as Markdown Button

## User Story

As a user, I want to download the current session as a Markdown file, so that I can save it locally or share it with others.

## Acceptance Criteria

- [ ] Download button visible in session view UI
- [ ] Clicking button downloads the MD file
- [ ] Downloaded file includes YAML frontmatter
- [ ] Filename matches session ID (e.g., `abc123.md`)
- [ ] Button has appropriate icon and tooltip
- [ ] Button placement matches existing UI patterns

## Technical Notes

- Use Blob API to create downloadable file
- Use anchor element with download attribute
- MD content is already fetched for rendering, reuse it
- Consider placement: toolbar area near session title

## Dependencies

- Story 001 (MD content must be fetched first)

## Estimation

Complexity: S
