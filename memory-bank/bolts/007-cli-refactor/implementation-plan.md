---
stage: plan
bolt: 007-cli-refactor
created: 2025-12-29T15:00:00Z
---

## Implementation Plan: CLI Refactor

### Objective

Simplify CLI from 4 commands (generate, serve, watch, version) to 2 commands (serve, version), with serve auto-generating HTML on startup.

### Deliverables

- Simplified CLI with only `serve` and `version` commands
- Auto-generation on serve startup (always, not just when dir missing)
- Renamed flag `--dir`/`-d` (was `--output-dir`/`-o`)
- New default output directory `~/claude-code-logs/` (was `~/.claude-code-logs`)
- Legacy command handlers showing migration messages
- Updated documentation

### Dependencies

- None (builds on completed Intent 001)

### Technical Approach

1. **Delete cmd_generate.go and cmd_watch.go** - Remove command files entirely
2. **Update main.go**:
   - Change default output from `~/.claude-code-logs` to `~/claude-code-logs/`
   - Rename `--output-dir`/`-o` to `--dir`/`-d`
   - Remove command registrations for generate and watch
   - Add hidden legacy commands that show migration messages
3. **Update cmd_serve.go**:
   - Always run generation on startup (not just when dir doesn't exist)
   - Ensure flag documentation reflects new names
4. **Update README.md**:
   - Update Quick Start section
   - Update Commands Overview
   - Update Configuration table
   - Update examples

### Acceptance Criteria

- [ ] `claude-code-logs serve` generates HTML then starts server
- [ ] `claude-code-logs serve --watch` enables file watching
- [ ] `claude-code-logs serve --dir /path` uses custom directory
- [ ] Default output is `~/claude-code-logs/`
- [ ] Old commands show helpful migration messages
- [ ] All existing tests pass
- [ ] README accurately reflects new CLI structure
