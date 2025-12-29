---
id: 001-cli-commands
unit: 005-cli
intent: 001-chat-log-viewer
status: draft
priority: must
created: 2025-12-29T12:00:00Z
assigned_bolt: 004-cli
implemented: false
---

# Story: 001-cli-commands

## User Story

**As a** developer
**I want** simple CLI commands
**So that** I can generate, serve, and watch my logs

## Acceptance Criteria

- [ ] **Given** user runs `claude-logs generate`, **When** executed, **Then** HTML files generated in output directory
- [ ] **Given** user runs `claude-logs serve`, **When** executed, **Then** HTTP server starts on default port 8080
- [ ] **Given** user runs `claude-logs watch`, **When** executed, **Then** watcher monitors for changes
- [ ] **Given** user runs `claude-logs serve --watch`, **When** executed, **Then** server starts with background watcher
- [ ] **Given** user provides `--output-dir PATH`, **When** any command runs, **Then** uses specified output directory
- [ ] **Given** user provides `--port PORT`, **When** serve runs, **Then** uses specified port
- [ ] **Given** user runs `claude-logs version`, **When** executed, **Then** shows version number
- [ ] **Given** any command running, **When** Ctrl+C pressed, **Then** graceful shutdown

## Technical Notes

- Use Cobra for CLI framework
- Default output directory: `~/.claude-logs`
- Default port: 8080
- Default watch interval: 30 seconds
- Follow Go CLI conventions (flags, help text)
- Non-zero exit code on errors
- Version embedded at build time via ldflags

## Dependencies

### Requires
- 001-parser (for parsing logs)
- 002-generator (for generating HTML)
- 003-server (for serving content)
- 004-watcher (for watching changes)

### Enables
- 006-homebrew-tap (distributes the CLI binary)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Invalid port number | Clear error message |
| Output dir not writable | Error message, exit 1 |
| No Claude projects found | Generate empty index, info message |
| Running multiple instances | Second instance fails with port conflict |
| Very long output path | Handle correctly |
| Config file present | Merge with CLI flags (flags take precedence) |

## Out of Scope

- GUI interface
- Interactive mode
- Shell completion (future enhancement)
- Config file creation command
