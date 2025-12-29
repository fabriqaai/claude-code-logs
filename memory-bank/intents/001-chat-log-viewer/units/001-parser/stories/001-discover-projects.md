---
id: 001-discover-projects
unit: 001-parser
intent: 001-chat-log-viewer
status: complete
priority: must
created: 2025-12-29T12:00:00.000Z
assigned_bolt: 001-parser
implemented: true
---

# Story: 001-discover-projects

## User Story

**As a** developer using Claude Code
**I want** my projects automatically discovered
**So that** I can browse all my chat history without configuration

## Acceptance Criteria

- [ ] **Given** `~/.claude/projects/` exists with project folders, **When** I run discovery, **Then** all project folders are listed
- [ ] **Given** a folder name like `-Users-name-myproject`, **When** decoding, **Then** it returns `/Users/name/myproject`
- [ ] **Given** empty `~/.claude/projects/`, **When** I run discovery, **Then** an empty list is returned without error
- [ ] **Given** `~/.claude/` does not exist, **When** I run discovery, **Then** a clear error message is shown
- [ ] **Given** mixed valid/invalid folders, **When** I run discovery, **Then** valid projects are returned, invalid ones skipped with warning

## Technical Notes

- Path decoding: Replace leading dash with `/`, replace `-` with `/` in path segments
- Use `os.UserHomeDir()` to locate home directory cross-platform
- Return `[]Project` slice, not error on empty
- Consider Windows path encoding for future compatibility

## Dependencies

### Requires
- None (first story - foundational)

### Enables
- 002-parse-sessions (needs project paths to find sessions)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Folder with no `.jsonl` files | Include project with empty Sessions slice |
| Symlinked project folder | Follow symlink, treat as regular folder |
| Permission denied on folder | Skip folder, log warning |
| Very long path (200+ chars) | Handle without truncation |
| Non-UTF8 folder name | Skip with warning |

## Out of Scope

- Parsing session contents (story 002)
- Custom source directory (handled by CLI)
- Recursive project discovery beyond immediate children
