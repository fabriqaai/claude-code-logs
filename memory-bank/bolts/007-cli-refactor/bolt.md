---
id: 007-cli-refactor
unit: 001-cli-refactor
intent: 002-simplify-cli
type: simple-construction-bolt
status: completed
stories:
  - 001-remove-commands
  - 002-update-serve
  - 003-update-docs
created: 2025-12-29T14:20:00Z
started: 2025-12-29T15:00:00Z
completed: 2025-12-29T15:25:00Z
current_stage: test
stages_completed:
  - name: plan
    completed: 2025-12-29T15:00:00Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T15:15:00Z
    artifact: implementation-walkthrough.md
  - name: test
    completed: 2025-12-29T15:25:00Z
    artifact: test-walkthrough.md

# Bolt Dependencies
requires_bolts: []
enables_bolts: []
requires_units: []
blocks: false

# Complexity Assessment
complexity:
  avg_complexity: 1        # Low - straightforward refactoring
  avg_uncertainty: 1       # Low - clear requirements
  max_dependencies: 0      # No external dependencies
  testing_scope: 2         # Integration - CLI behavior testing
---

# Bolt: 007-cli-refactor

## Overview

Refactor CLI from 4 commands to 2, update default output directory, and rename flags.

## Objective

Simplify the CLI interface by:
1. Removing `generate` and `watch` commands
2. Making `serve` auto-generate on startup
3. Changing default output to `~/claude-code-logs/`
4. Renaming `--output-dir` to `--dir`
5. Updating all documentation

## Stories Included

- **001-remove-commands**: Remove generate and watch commands (Must)
- **002-update-serve**: Update serve with auto-generate and new flags (Must)
- **003-update-docs**: Update README and help text (Must)

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [x] **1. plan**: Review existing code and plan changes
- [x] **2. implement**: Make code changes
- [x] **3. test**: Verify all acceptance criteria
- [x] **4. document**: Update README and help text (done as part of implement)

## Files to Modify

| File | Action | Description |
|------|--------|-------------|
| `cmd/generate.go` | DELETE | Remove generate command |
| `cmd/watch.go` | DELETE | Remove watch command |
| `cmd/root.go` | MODIFY | Remove command registrations, add legacy handlers |
| `cmd/serve.go` | MODIFY | Add auto-generate, rename flags, update defaults |
| `README.md` | MODIFY | Update documentation |

## Dependencies

### Requires
- None (builds on completed Intent 001)

### Enables
- None (final refactoring bolt)

## Success Criteria

- [x] `claude-logs serve` generates HTML then starts server
- [x] `claude-logs serve --watch` enables file watching
- [x] `claude-logs serve --dir /path` uses custom directory
- [x] Default output is `~/claude-code-logs/`
- [x] Old commands show helpful migration messages
- [x] All existing tests pass
- [x] README accurately reflects new CLI structure

## Notes

This is a focused refactoring bolt. The underlying parser, generator, server, and watcher logic remain unchanged. Only the CLI interface is being simplified.

Breaking change: Users with scripts using `generate` or `watch` commands will need to update them.
