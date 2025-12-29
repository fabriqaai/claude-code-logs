# Global Story Index

## Overview

- **Total stories**: 16
- **Generated**: 10
- **Assigned to bolts**: 16
- **Last updated**: 2025-12-29T16:30:00Z

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

### 002-simplify-cli

#### Unit: 001-cli-refactor (3 stories)

- [x] **001-remove-commands** (cli-refactor): Remove generate and watch commands - Must - ✅ GENERATED - Bolt: 007-cli-refactor
- [x] **002-update-serve** (cli-refactor): Update serve with auto-generate and new flags - Must - ✅ GENERATED - Bolt: 007-cli-refactor
- [x] **003-update-docs** (cli-refactor): Update README and help text - Must - ✅ GENERATED - Bolt: 007-cli-refactor

---

### 003-tree-view-sidebar

#### Unit: 001-tree-view (3 stories)

- [ ] **001-tree-markup** (tree-view): Tree view HTML and CSS structure - Must - Bolt: 008-tree-view
- [ ] **002-expand-collapse** (tree-view): JavaScript toggle functionality - Must - Bolt: 008-tree-view
- [ ] **003-visual-polish** (tree-view): Animations and visual enhancements - Should - Bolt: 008-tree-view

#### Unit: 002-resize-persist (3 stories)

- [ ] **001-resize-handle** (resize-persist): Draggable resize functionality - Must - Bolt: 009-resize-persist
- [ ] **002-localstorage** (resize-persist): Persist width and collapse state - Should - Bolt: 009-resize-persist
- [ ] **003-mobile-menu** (resize-persist): Hamburger menu and drawer overlay - Should - Bolt: 009-resize-persist

---

## Stories by Priority

### Must Have (11 stories)

| Story | Unit | Bolt | Description |
|-------|------|------|-------------|
| 001-discover-projects | 001-parser | 001-parser | Discover Claude Code projects |
| 002-parse-sessions | 001-parser | 001-parser | Parse JSONL session files |
| 001-generate-html | 002-generator | 002-generator | Generate HTML pages |
| 001-serve-and-search | 003-server | 003-server | Serve HTML and search API |
| 001-cli-commands | 005-cli | 004-cli | CLI commands |
| 001-remove-commands | 001-cli-refactor | 007-cli-refactor | Remove generate/watch commands |
| 002-update-serve | 001-cli-refactor | 007-cli-refactor | Update serve command |
| 003-update-docs | 001-cli-refactor | 007-cli-refactor | Update documentation |
| 001-tree-markup | 001-tree-view | 008-tree-view | Tree view HTML/CSS |
| 002-expand-collapse | 001-tree-view | 008-tree-view | Expand/collapse toggle |
| 001-resize-handle | 002-resize-persist | 009-resize-persist | Resizable sidebar |

### Should Have (5 stories)

| Story | Unit | Bolt | Description |
|-------|------|------|-------------|
| 001-watch-changes | 004-watcher | 005-watcher | Watch for changes |
| 001-homebrew-distribution | 006-homebrew-tap | 006-homebrew-tap | Homebrew distribution |
| 003-visual-polish | 001-tree-view | 008-tree-view | Visual polish |
| 002-localstorage | 002-resize-persist | 009-resize-persist | State persistence |
| 003-mobile-menu | 002-resize-persist | 009-resize-persist | Mobile hamburger menu |

### Could Have (0 stories)

None.

---

## Stories by Status

- **Planned**: 6 (Intent 003)
- **Generated**: 10 (Intents 001, 002)
- **In Progress**: 0
- **Completed**: 0

---

## Bolt Plan

### Execution Order

```
Intent 001 (completed):
001-parser ──► 002-generator ──► 003-server ──► 004-cli ──► 005-watcher
                                                    │
                                                    └──► 006-homebrew-tap

Intent 002 (completed):
007-cli-refactor (standalone, builds on Intent 001)

Intent 003 (planned):
008-tree-view ──► 009-resize-persist
```

### Bolts Created

| Bolt ID | Intent | Unit | Stories | Status |
|---------|--------|------|---------|--------|
| 001-parser | 001 | 001-parser | 001-discover-projects, 002-parse-sessions | completed |
| 002-generator | 001 | 002-generator | 001-generate-html | completed |
| 003-server | 001 | 003-server | 001-serve-and-search | completed |
| 004-cli | 001 | 005-cli | 001-cli-commands | completed |
| 005-watcher | 001 | 004-watcher | 001-watch-changes | completed |
| 006-homebrew-tap | 001 | 006-homebrew-tap | 001-homebrew-distribution | completed |
| 007-cli-refactor | 002 | 001-cli-refactor | 001-remove-commands, 002-update-serve, 003-update-docs | completed |
| 008-tree-view | 003 | 001-tree-view | 001-tree-markup, 002-expand-collapse, 003-visual-polish | planned |
| 009-resize-persist | 003 | 002-resize-persist | 001-resize-handle, 002-localstorage, 003-mobile-menu | planned |

---

## Dependency Analysis

### Intent 001 Dependencies (completed)

- **002-generator** requires 001-parser (needs parsed session data)
- **003-server** requires 001-parser, 002-generator (needs data and HTML)
- **004-cli** requires 001-parser, 002-generator, 003-server (orchestrates all)
- **005-watcher** requires 001-parser, 002-generator, 004-cli (triggered by CLI)
- **006-homebrew-tap** requires 004-cli (distributes the binary)

### Intent 002 Dependencies (completed)

- **007-cli-refactor** requires Intent 001 complete (modifies existing CLI code)

### Intent 003 Dependencies

- **008-tree-view** requires none (modifies templates.go)
- **009-resize-persist** requires 008-tree-view (needs tree structure for collapse state)

### Cross-Intent Dependencies

- Intent 002 builds on Intent 001's completed CLI
- Intent 003 builds on Intent 001's templates (no direct dependency on Intent 002)

### Dependency Warnings

None - dependency chain is clear and acyclic.

---

## Directories Created

### Intent 001
✅ `memory-bank/bolts/001-parser/bolt.md`
✅ `memory-bank/bolts/002-generator/bolt.md`
✅ `memory-bank/bolts/003-server/bolt.md`
✅ `memory-bank/bolts/004-cli/bolt.md`
✅ `memory-bank/bolts/005-watcher/bolt.md`
✅ `memory-bank/bolts/006-homebrew-tap/bolt.md`

### Intent 002
✅ `memory-bank/bolts/007-cli-refactor/bolt.md`

### Intent 003
✅ `memory-bank/bolts/008-tree-view/bolt.md`
✅ `memory-bank/bolts/009-resize-persist/bolt.md`

---

## Summary

- **9 bolts** created (6 from Intent 001, 1 from Intent 002, 2 from Intent 003)
- **16 stories** covered
- All stories assigned to bolts
- Clear dependency chain established
- Intent 001 complete, Intent 002 complete, Intent 003 ready for construction
