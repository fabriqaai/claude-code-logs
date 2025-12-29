---
id: 001-remove-commands
unit: 001-cli-refactor
intent: 002-simplify-cli
status: ready
priority: must
created: 2025-12-29T14:15:00Z
assigned_bolt: 007-cli-refactor
implemented: false
---

# Story: 001-remove-commands

## User Story

**As a** CLI user
**I want** a simpler command structure
**So that** I have fewer commands to remember and the tool is easier to use

## Acceptance Criteria

- [ ] **Given** I run `claude-logs generate`, **When** the command executes, **Then** I see a helpful message: "Command removed. Use 'claude-logs serve' instead (generates automatically on startup)"
- [ ] **Given** I run `claude-logs watch`, **When** the command executes, **Then** I see a helpful message: "Command removed. Use 'claude-logs serve --watch' instead"
- [ ] **Given** `cmd/generate.go` exists, **When** refactoring is complete, **Then** the file is deleted
- [ ] **Given** `cmd/watch.go` exists, **When** refactoring is complete, **Then** the file is deleted
- [ ] **Given** I run `claude-logs --help`, **When** viewing help, **Then** I only see `serve` and `version` commands listed

## Technical Notes

- Delete `cmd/generate.go` entirely
- Delete `cmd/watch.go` entirely
- Update `cmd/root.go` to remove command registrations
- Consider adding hidden "legacy" commands that show migration messages

## Dependencies

### Requires
- None (first story)

### Enables
- 002-update-serve (serve command needs to absorb generate functionality)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| User has scripts using `generate` | Clear error message guides migration |
| User has scripts using `watch` | Clear error message guides migration |
| Tab completion for old commands | Should not appear |

## Out of Scope

- Actual generation logic (moves to serve, story 002)
- Watch logic (already in serve --watch)
