---
intent: 001-chat-log-viewer
type: units-index
created: 2025-12-29T11:20:00Z
---

# Units: chat-log-viewer

## Overview

This intent is decomposed into 6 units for a CLI tool architecture.

## Units

| # | Unit | Purpose | Dependencies | Stories |
|---|------|---------|--------------|---------|
| 001 | parser | Parse JSONL logs into structured data | None | TBD |
| 002 | generator | Generate HTML from parsed data | 001-parser | TBD |
| 003 | server | HTTP server for HTML + Search API | 002-generator | TBD |
| 004 | watcher | File system watcher for regeneration | 002-generator | TBD |
| 005 | cli | CLI commands (generate, serve, watch) | All above | TBD |
| 006 | homebrew-tap | Homebrew tap repository (fabriqa/tap) | 005-cli | TBD |

## Requirement Mapping

| FR | Description | Unit |
|----|-------------|------|
| FR-1 | Project Discovery | 001-parser |
| FR-2 | JSONL Parsing | 001-parser |
| FR-3 | HTML Generation | 002-generator |
| FR-4 | Server Mode | 003-server |
| FR-5 | Search Functionality | 003-server |
| FR-6 | Watch Mode | 004-watcher |
| FR-7 | Static Fallback UI | 002-generator |
| FR-8 | CLI Commands | 005-cli |
| FR-9 | Homebrew Distribution | 006-homebrew-tap |

## Dependency Graph

```
001-parser ──► 002-generator ──► 003-server
                    │
                    ▼
              004-watcher
                    │
                    ▼
               005-cli
                    │
                    ▼
          006-homebrew-tap
```

## Build Order

1. **001-parser** - No dependencies, foundational
2. **002-generator** - Depends on parser output
3. **003-server** - Serves generated HTML
4. **004-watcher** - Triggers generator on changes
5. **005-cli** - Orchestrates all components
6. **006-homebrew-tap** - Distribution (separate repo)
