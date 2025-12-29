---
id: 005-watcher
unit: 004-watcher
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: in-progress
stories:
  - 001-watch-changes
created: 2025-12-29T12:45:00Z
started: 2025-12-29T13:00:00Z
completed: null
current_stage: test
stages_completed:
  - name: plan
    completed: 2025-12-29T13:05:00Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T13:20:00Z
    artifact: implementation-walkthrough.md

# Bolt Dependencies
requires_bolts:
  - 001-parser
  - 002-generator
  - 004-cli
enables_bolts: []
requires_units: []
blocks: false

# Complexity Assessment
complexity:
  avg_complexity: 2
  avg_uncertainty: 1
  max_dependencies: 2
  testing_scope: 2
---

# Bolt: 005-watcher

## Overview

File watcher that monitors ~/.claude/projects/ for changes and triggers HTML regeneration when new sessions appear or existing ones are modified.

## Objective

Implement file watching with fsnotify, debouncing, and incremental regeneration of affected projects.

## Stories Included

- **001-watch-changes**: Watch for file changes and regenerate HTML - Should

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 001-parser (for re-parsing changed sessions)
- 002-generator (for regenerating HTML)
- 004-cli (CLI exposes watch command and --watch flag)

### Enables
- None (final enhancement bolt)

## Success Criteria

- [ ] Detects new .jsonl files within poll interval
- [ ] Detects modified .jsonl files
- [ ] Debounces rapid changes (2 second delay)
- [ ] Only regenerates affected project's HTML
- [ ] Configurable poll interval via --interval flag
- [ ] Graceful shutdown on Ctrl+C
- [ ] Works alongside server mode (serve --watch)
- [ ] Tests passing

## Notes

- Use fsnotify for cross-platform file watching
- Fall back to polling if fsnotify unavailable
- Default poll interval: 30 seconds
- Default debounce delay: 2 seconds
- Run in goroutine for background operation
