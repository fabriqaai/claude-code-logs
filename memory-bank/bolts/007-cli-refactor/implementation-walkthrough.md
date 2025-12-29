---
stage: implement
bolt: 007-cli-refactor
created: 2025-12-29T15:15:00Z
---

## Implementation Walkthrough: CLI Refactor

### Summary

Simplified the CLI from 4 commands to 2 by removing `generate` and `watch` commands, making `serve` the primary entry point with auto-generation. Changed the default output directory and renamed flags for clarity.

### Structure Overview

The CLI now has a streamlined command structure with `serve` as the main command. Legacy commands are kept as hidden commands that display migration messages to help users update their scripts.

### Completed Work

- [x] `main.go` - Updated global flag from `outputDir`/`--output-dir`/`-o` to `dirFlag`/`--dir`/`-d`, changed default from `~/.claude-code-logs` to `~/claude-code-logs`, added hidden legacy command handlers
- [x] `cmd_serve.go` - Added auto-generation on startup (always generates, not just when dir missing), moved `ensureWritableDir` here, updated help text
- [x] `cmd_generate.go` - DELETED (functionality absorbed by serve)
- [x] `cmd_watch.go` - DELETED (functionality available via serve --watch)
- [x] `cmd_test.go` - Updated tests for new flag name, new default directory, and legacy command behavior
- [x] `README.md` - Complete documentation overhaul reflecting new CLI structure

### Key Decisions

- **Hidden legacy commands**: Rather than removing commands entirely, they're kept hidden and show helpful migration messages
- **Always generate on startup**: Unlike before where generation only happened if dir was missing, now serve always regenerates for fresh content
- **Flag rename**: `--dir/-d` is clearer than `--output-dir/-o` for the simplified use case

### Deviations from Plan

None - implementation followed plan exactly.

### Dependencies Added

None - no new dependencies required.

### Developer Notes

Breaking change for users with scripts using `generate` or `watch` commands. The migration messages guide them to the new patterns.
