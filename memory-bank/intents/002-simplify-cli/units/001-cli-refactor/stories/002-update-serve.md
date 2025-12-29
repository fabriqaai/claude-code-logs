---
id: 002-update-serve
unit: 001-cli-refactor
intent: 002-simplify-cli
status: done
priority: must
created: 2025-12-29T14:15:00Z
assigned_bolt: 007-cli-refactor
implemented: true
---

# Story: 002-update-serve

## User Story

**As a** CLI user
**I want** the serve command to generate HTML automatically on startup
**So that** I don't need to run a separate generate command first

## Acceptance Criteria

- [ ] **Given** I run `claude-code-logs serve`, **When** the server starts, **Then** HTML is generated before the server begins listening
- [ ] **Given** I run `claude-code-logs serve --dir /custom/path`, **When** the server starts, **Then** HTML is generated to `/custom/path` and served from there
- [ ] **Given** I run `claude-code-logs serve` without --dir, **When** the server starts, **Then** default output directory is `~/claude-code-logs/` (not `~/.claude-code-logs/`)
- [ ] **Given** the output directory doesn't exist, **When** serve runs, **Then** the directory is created automatically
- [ ] **Given** I run `claude-code-logs serve -d /path`, **When** the server starts, **Then** `-d` works as short form of `--dir`
- [ ] **Given** I run `claude-code-logs serve --output-dir`, **When** parsing flags, **Then** flag is not recognized (old flag removed)
- [ ] **Given** I run `claude-code-logs serve -o`, **When** parsing flags, **Then** flag is not recognized (old flag removed)

## Technical Notes

- Move generation call to start of serve command's Run function
- Rename flag from `--output-dir`/`-o` to `--dir`/`-d`
- Update default value for output directory
- Preserve all existing serve functionality (port, watch, verbose)

### Flag Changes

| Old Flag | New Flag | Default |
|----------|----------|---------|
| `--output-dir` | `--dir` | `~/claude-code-logs/` |
| `-o` | `-d` | `~/claude-code-logs/` |
| `--port` | `--port` | `8080` (unchanged) |
| `-p` | `-p` | `8080` (unchanged) |
| `--watch` | `--watch` | `false` (unchanged) |
| `-w` | `-w` | `false` (unchanged) |
| `--verbose` | `--verbose` | `false` (unchanged) |
| `-v` | `-v` | `false` (unchanged) |

## Dependencies

### Requires
- 001-remove-commands (generate command removed)

### Enables
- 003-update-docs (documentation needs accurate commands)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Generation fails | Show error, don't start server |
| Output dir not writable | Show permission error |
| ~/ expansion on different OS | Works on macOS and Linux |
| $HOME not set | Fallback to `/tmp/claude-code-logs/` or error |

## Out of Scope

- Changes to generation logic itself
- Changes to server logic itself
- Changes to watch logic itself
