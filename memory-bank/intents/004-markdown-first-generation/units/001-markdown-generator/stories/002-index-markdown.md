---
story: 002-index-markdown
unit: 001-markdown-generator
intent: 004-markdown-first-generation
priority: Must
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Generate Index Markdown Files

## User Story

As a user, I want index files generated in Markdown format listing all projects and sessions, so that navigation data is also in a portable format.

## Acceptance Criteria

- [ ] Main `index.md` lists all projects with:
  - Project path
  - Session count
  - Link to project index
- [ ] Per-project `{project}/index.md` lists all sessions with:
  - Session title/summary
  - Creation date
  - Link to session MD
- [ ] Output is deterministic (sorted by last update date)
- [ ] No timestamps in content that would cause git churn

## Technical Notes

- Replace `templates_index.go` and `templates_project.go` with MD generation
- Sort projects by last update (existing logic)
- Sort sessions by creation date (existing logic)
- Use relative links for navigation

## Dependencies

- Parser provides Project and Session structs
- Existing sorting logic in parser.go

## Estimation

Complexity: S
