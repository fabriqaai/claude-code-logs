# Global Story Index

## Overview

- **Total stories**: 7
- **Generated**: 7
- **Assigned to bolts**: 7
- **Last updated**: 2025-12-29T12:50:00Z

---

## Stories by Intent

### 001-chat-log-viewer

#### Unit: 001-parser (2 stories)

- [x] **001-discover-projects** (parser): Discover Claude Code projects - Must - ✅ GENERATED - Bolt: 001-parser
- [x] **002-parse-sessions** (parser): Parse JSONL session files - Must - ✅ GENERATED - Bolt: 001-parser

#### Unit: 002-generator (1 story)

- [x] **001-generate-html** (generator): Generate HTML pages with Claude.ai styling - Must - ✅ GENERATED - Bolt: 002-generator

#### Unit: 003-server (1 story)

- [x] **001-serve-and-search** (server): Serve HTML and provide search API - Must - ✅ GENERATED - Bolt: 003-server

#### Unit: 004-watcher (1 story)

- [x] **001-watch-changes** (watcher): Watch for file changes and regenerate - Should - ✅ GENERATED - Bolt: 005-watcher

#### Unit: 005-cli (1 story)

- [x] **001-cli-commands** (cli): Implement CLI commands - Must - ✅ GENERATED - Bolt: 004-cli

#### Unit: 006-homebrew-tap (1 story)

- [x] **001-homebrew-distribution** (homebrew-tap): Create Homebrew tap and formula - Should - ✅ GENERATED - Bolt: 006-homebrew-tap

---

## Stories by Priority

### Must Have (5 stories)

| Story | Unit | Bolt | Description |
|-------|------|------|-------------|
| 001-discover-projects | 001-parser | 001-parser | Discover Claude Code projects |
| 002-parse-sessions | 001-parser | 001-parser | Parse JSONL session files |
| 001-generate-html | 002-generator | 002-generator | Generate HTML pages |
| 001-serve-and-search | 003-server | 003-server | Serve HTML and search API |
| 001-cli-commands | 005-cli | 004-cli | CLI commands |

### Should Have (2 stories)

| Story | Unit | Bolt | Description |
|-------|------|------|-------------|
| 001-watch-changes | 004-watcher | 005-watcher | Watch for changes |
| 001-homebrew-distribution | 006-homebrew-tap | 006-homebrew-tap | Homebrew distribution |

### Could Have (0 stories)

None.

---

## Stories by Status

- **Planned**: 0
- **Generated**: 7
- **In Progress**: 0
- **Completed**: 0

---

## Bolt Plan

### Execution Order

```
001-parser ──► 002-generator ──► 003-server ──► 004-cli ──► 005-watcher
                                                    │
                                                    └──► 006-homebrew-tap
```

### Bolts Created

| Bolt ID | Unit | Stories | Status |
|---------|------|---------|--------|
| 001-parser | 001-parser | 001-discover-projects, 002-parse-sessions | planned |
| 002-generator | 002-generator | 001-generate-html | planned |
| 003-server | 003-server | 001-serve-and-search | planned |
| 004-cli | 005-cli | 001-cli-commands | planned |
| 005-watcher | 004-watcher | 001-watch-changes | planned |
| 006-homebrew-tap | 006-homebrew-tap | 001-homebrew-distribution | planned |

---

## Dependency Analysis

### Within-Intent Dependencies

- **002-generator** requires 001-parser (needs parsed session data)
- **003-server** requires 001-parser, 002-generator (needs data and HTML)
- **004-cli** requires 001-parser, 002-generator, 003-server (orchestrates all)
- **005-watcher** requires 001-parser, 002-generator, 004-cli (triggered by CLI)
- **006-homebrew-tap** requires 004-cli (distributes the binary)

### Cross-Unit Dependencies

None (single intent, all units within same intent)

### Dependency Warnings

None - dependency chain is clear and acyclic.

---

## Directories Created

✅ `memory-bank/bolts/001-parser/bolt.md`
✅ `memory-bank/bolts/002-generator/bolt.md`
✅ `memory-bank/bolts/003-server/bolt.md`
✅ `memory-bank/bolts/004-cli/bolt.md`
✅ `memory-bank/bolts/005-watcher/bolt.md`
✅ `memory-bank/bolts/006-homebrew-tap/bolt.md`

---

## Summary

- **6 bolts** created
- **7 stories** covered
- All stories assigned to bolts
- Clear dependency chain established
