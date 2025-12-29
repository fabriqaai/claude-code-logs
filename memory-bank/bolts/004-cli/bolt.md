---
id: 004-cli
unit: 005-cli
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: complete
stories:
  - 001-cli-commands
created: 2025-12-29T12:45:00.000Z
started: 2025-12-29T16:30:00.000Z
completed: "2025-12-29T14:58:55Z"
current_stage: null
stages_completed:
  - name: plan
    completed: 2025-12-29T16:40:00.000Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T16:55:00.000Z
    artifact: implementation-walkthrough.md
requires_bolts:
  - 001-parser
  - 002-generator
  - 003-server
enables_bolts:
  - 005-watcher
  - 006-homebrew-tap
requires_units: []
blocks: false
complexity:
  avg_complexity: 2
  avg_uncertainty: 1
  max_dependencies: 3
  testing_scope: 2
---

# Bolt: 004-cli

## Overview

Command-line interface that orchestrates all components (parser, generator, server) and provides user-facing commands.

## Objective

Implement CLI commands using Cobra framework: generate, serve, watch, and version commands with appropriate flags.

## Stories Included

- **001-cli-commands**: Implement CLI commands - Must

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 001-parser (for parsing logs)
- 002-generator (for generating HTML)
- 003-server (for serving content)

### Enables
- 005-watcher (CLI exposes watch command)
- 006-homebrew-tap (distributes CLI binary)

## Success Criteria

- [ ] `claude-logs generate` works
- [ ] `claude-logs serve` starts server on port 8080
- [ ] `claude-logs serve --watch` combines serve and watch
- [ ] `claude-logs watch` monitors for changes
- [ ] `claude-logs version` shows version
- [ ] `--output-dir` flag works on all commands
- [ ] `--port` flag works with serve
- [ ] Graceful shutdown on Ctrl+C
- [ ] Clear error messages
- [ ] Help text for all commands
- [ ] Tests passing

## Notes

- Use Cobra CLI framework
- Default output directory: ~/.claude-logs
- Version embedded at build time via ldflags
- Follow Go CLI conventions
