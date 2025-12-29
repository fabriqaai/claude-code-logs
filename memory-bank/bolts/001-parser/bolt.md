---
id: 001-parser
unit: 001-parser
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: planned
stories:
  - 001-discover-projects
  - 002-parse-sessions
created: 2025-12-29T12:45:00Z
started: null
completed: null
current_stage: null
stages_completed: []

# Bolt Dependencies
requires_bolts: []
enables_bolts:
  - 002-generator
  - 003-server
  - 004-cli
  - 005-watcher
requires_units: []
blocks: false

# Complexity Assessment
complexity:
  avg_complexity: 2
  avg_uncertainty: 1
  max_dependencies: 1
  testing_scope: 2
---

# Bolt: 001-parser

## Overview

First bolt establishing the foundation for parsing Claude Code JSONL session files. This is the core data layer that all other bolts depend on.

## Objective

Implement project discovery and JSONL session parsing functionality to extract structured data from Claude Code log files.

## Stories Included

- **001-discover-projects**: Discover Claude Code projects from ~/.claude/projects/ - Must
- **002-parse-sessions**: Parse JSONL session files into structured Go types - Must

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- None (first bolt - foundational)

### Enables
- 002-generator (needs parsed session data)
- 003-server (needs parsed data for search index)
- 004-cli (orchestrates parser)
- 005-watcher (triggers parser on file changes)

## Success Criteria

- [ ] Discovers all projects in ~/.claude/projects/
- [ ] Correctly decodes folder names to human-readable paths
- [ ] Parses JSONL files line-by-line (streaming)
- [ ] Extracts all message types (human, assistant, summary)
- [ ] Handles tool_use and tool_result blocks
- [ ] Graceful handling of malformed JSONL lines
- [ ] Tests passing for all acceptance criteria

## Notes

- This bolt uses the simple-construction-bolt type as it's CLI/utility work
- JSONL parsing should be streaming to handle large files
- Path decoding logic based on Claude Code folder naming convention
