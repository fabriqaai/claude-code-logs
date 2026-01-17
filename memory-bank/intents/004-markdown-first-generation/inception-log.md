---
intent: 004-markdown-first-generation
created: 2026-01-17T16:00:00Z
completed: null
status: approved
---

# Inception Log: markdown-first-generation

## Overview

**Intent**: Transform output to Markdown-first generation with client-side rendering, enabling LLM-friendly logs, git archival, and history preservation
**Type**: green-field
**Created**: 2026-01-17

## Artifacts Created

| Artifact | Status | File |
|----------|--------|------|
| Requirements | ✅ | requirements.md |
| System Context | ✅ | system-context.md |
| Units | ✅ | units.md, units/*/unit-brief.md |
| Stories | ✅ | units/*/stories/*.md |
| Bolt Plan | ✅ | bolt-plan.md |

## Summary

| Metric | Count |
|--------|-------|
| Functional Requirements | 7 |
| Non-Functional Requirements | 3 categories |
| Units | 2 |
| Stories | 7 |
| Bolts Planned | 2 |

## Units Breakdown

| Unit | Stories | Bolts | Priority |
|------|---------|-------|----------|
| 001-markdown-generator | 4 | bolt-001-generator | Must |
| 002-client-renderer | 3 | bolt-002-renderer | Must |

## Decision Log

| Date | Decision | Rationale | Approved |
|------|----------|-----------|----------|
| 2026-01-17 | Use marked.js for MD rendering | Fastest, simplest API, native hljs integration | Yes |
| 2026-01-17 | Use highlight.js for syntax highlighting | Best marked.js integration, auto-language detection | Yes |
| 2026-01-17 | Use mtime for regeneration check | Faster than hash, sufficient for local use | Yes |
| 2026-01-17 | Store hash in frontmatter anyway | Future flexibility for `--use-hash` flag | Yes |
| 2026-01-17 | Preserve orphaned MD files | Archive sessions beyond Claude's 30-day retention | Yes |
| 2026-01-17 | Include frontmatter in clipboard copy | Full context for LLM consumption | Yes |
| 2026-01-17 | No indicator for archived sessions | Keep UI simple, archives treated same as active | Yes |
| 2026-01-17 | MD-only output (no dual HTML) | Clean transition, simpler codebase | Yes |

## Scope Changes

| Date | Change | Reason | Impact |
|------|--------|--------|--------|
| - | - | - | - |

## Ready for Construction

**Checklist**:
- [x] All requirements documented
- [x] System context defined
- [x] Units decomposed
- [x] Stories created for all units
- [x] Bolts planned
- [x] Human review complete

## Next Steps

1. ~~Define system context~~ ✅
2. ~~Decompose into units~~ ✅
3. ~~Create user stories~~ ✅
4. ~~Plan bolts~~ ✅
5. ~~Review and approve for construction~~ ✅ Approved

## Dependencies

| From | To | Reason |
|------|----|--------|
| bolt-001-generator | bolt-002-renderer | Renderer needs MD files to exist |
| marked.js (external) | bolt-002-renderer | MD parsing library |
| highlight.js (external) | bolt-002-renderer | Code syntax highlighting |
