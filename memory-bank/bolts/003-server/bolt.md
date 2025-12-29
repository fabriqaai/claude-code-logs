---
id: 003-server
unit: 003-server
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: planned
stories:
  - 001-serve-and-search
created: 2025-12-29T12:45:00Z
started: null
completed: null
current_stage: null
stages_completed: []

# Bolt Dependencies
requires_bolts:
  - 001-parser
  - 002-generator
enables_bolts:
  - 004-cli
requires_units: []
blocks: false

# Complexity Assessment
complexity:
  avg_complexity: 2
  avg_uncertainty: 2
  max_dependencies: 2
  testing_scope: 2
---

# Bolt: 003-server

## Overview

HTTP server that serves generated HTML files and provides a search API for full-text search across Claude Code sessions.

## Objective

Implement local HTTP server with static file serving and POST /api/search endpoint for searching chat logs.

## Stories Included

- **001-serve-and-search**: Serve HTML and provide full-text search API - Must

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 001-parser (provides data for building search index)
- 002-generator (provides HTML files to serve)

### Enables
- 004-cli (exposes serve command)

## Success Criteria

- [ ] HTTP server on configurable port (default 8080)
- [ ] Serves static HTML files correctly
- [ ] POST /api/search endpoint functional
- [ ] Search with project/session filters
- [ ] Returns matches with context
- [ ] Response time < 500ms for search
- [ ] Highlights matching terms
- [ ] Graceful shutdown on Ctrl+C
- [ ] Tests passing

## Notes

- Use Go net/http standard library
- In-memory search index built at startup
- Consider using bleve for full-text search
- Bind to localhost only (127.0.0.1) for security
- CORS headers for local development
