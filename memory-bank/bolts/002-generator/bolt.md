---
id: 002-generator
unit: 002-generator
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: complete
stories:
  - 001-generate-html
created: 2025-12-29T12:45:00.000Z
started: 2025-12-29T16:30:00.000Z
completed: "2025-12-29T14:08:10Z"
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
enables_bolts:
  - 003-server
  - 004-cli
  - 005-watcher
requires_units: []
blocks: false
complexity:
  avg_complexity: 2
  avg_uncertainty: 1
  max_dependencies: 1
  testing_scope: 2
---

# Bolt: 002-generator

## Overview

Generates static HTML pages from parsed Claude Code session data with Claude.ai visual styling.

## Objective

Create HTML templates and generation logic to produce browsable chat log pages with project navigation, conversation view, and Claude.ai styling.

## Stories Included

- **001-generate-html**: Generate HTML pages with Claude.ai styling - Must

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 001-parser (provides []Project and Session data structures)

### Enables
- 003-server (serves generated HTML)
- 004-cli (exposes generate command)
- 005-watcher (triggers regeneration)

## Success Criteria

- [ ] Generates valid HTML5 pages
- [ ] Index.html with project navigation sidebar (288px)
- [ ] Session pages with conversation view
- [ ] Claude.ai visual styling (cream background, serif/sans-serif fonts)
- [ ] CSS embedded for offline viewing
- [ ] Footer with "claude-code-logs by fabriqa.ai" link
- [ ] Shows "Start server for search" in static mode
- [ ] Syntax highlighting for code blocks
- [ ] Tests passing

## Notes

- Use Go html/template for security (auto-escaping)
- Reference ui-standards.md for exact CSS values
- Atomic file writes (temp file then rename)
- No external CDN dependencies
