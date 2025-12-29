---
intent: 002-simplify-cli
created: 2025-12-29T14:00:00Z
completed: 2025-12-29T14:30:00Z
status: complete
---

# Inception Log: simplify-cli

## Overview

**Intent**: Simplify CLI from 4 commands to 2, change default output directory
**Type**: refactoring
**Created**: 2025-12-29

## Artifacts Created

| Artifact | Status | File |
|----------|--------|------|
| Requirements | ✅ | requirements.md |
| System Context | ✅ | system-context.md |
| Units | ✅ | units.md |
| Unit Brief | ✅ | units/001-cli-refactor/unit-brief.md |
| Stories | ✅ | units/001-cli-refactor/stories/*.md |
| Bolt Plan | ✅ | memory-bank/bolts/007-cli-refactor/bolt.md |

## Summary

| Metric | Count |
|--------|-------|
| Functional Requirements | 6 |
| Non-Functional Requirements | 2 |
| Units | 1 |
| Stories | 3 |
| Bolts Planned | 1 |

## Units Breakdown

| Unit | Stories | Bolts | Priority |
|------|---------|-------|----------|
| 001-cli-refactor | 3 | 1 | Must |

## Decision Log

| Date | Decision | Rationale | Approved |
|------|----------|-----------|----------|
| 2025-12-29 | Remove `generate` command | Redundant - serve can auto-generate | Yes |
| 2025-12-29 | Remove `watch` command | Redundant - serve --watch covers this | Yes |
| 2025-12-29 | Default to ~/claude-code-logs/ | More visible than hidden ~/.claude-code-logs/ | Yes |
| 2025-12-29 | Use --dir instead of --output-dir | Shorter, clearer | Yes |
| 2025-12-29 | Show helpful migration messages | Better UX for users with existing scripts | Yes |

## Scope Changes

| Date | Change | Reason | Impact |
|------|--------|--------|--------|
| None | - | - | - |

## Ready for Construction

**Checklist**:
- [x] All requirements documented
- [x] System context defined
- [x] Units decomposed
- [x] Stories created for all units
- [x] Bolts planned
- [x] Human review complete

## Next Steps

1. Begin Construction Phase
2. Start with Bolt: 007-cli-refactor
3. Execute: `/specsmd-construction-agent --bolt="007-cli-refactor"`

## Dependencies

- Requires Intent 001 (chat-log-viewer) to be complete
- All Intent 001 bolts are complete
