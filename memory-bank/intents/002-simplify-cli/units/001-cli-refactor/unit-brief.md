---
unit: 001-cli-refactor
intent: 002-simplify-cli
phase: inception
status: ready
created: 2025-12-29T14:10:00Z
updated: 2025-12-29T14:10:00Z
---

# Unit Brief: CLI Refactor

## Purpose

Simplify the CLI interface by reducing commands from 4 to 2, changing default output directory, and updating flag names.

## Scope

### In Scope
- Remove `generate` command (functionality moves to `serve` startup)
- Remove `watch` command (functionality stays in `serve --watch`)
- Change default output directory to `~/claude-code-logs/`
- Rename `--output-dir` to `--dir`
- Update all help text
- Update README.md

### Out of Scope
- Parser logic (unchanged)
- Generator logic (unchanged)
- Server logic (unchanged)
- Watcher logic (unchanged)
- HTML templates (unchanged)

---

## Assigned Requirements

| FR | Requirement | Priority |
|----|-------------|----------|
| FR-1 | Remove `generate` command | Must |
| FR-2 | Remove `watch` command | Must |
| FR-3 | Default output directory to ~/claude-code-logs/ | Must |
| FR-4 | Add --dir flag for custom working directory | Must |
| FR-5 | Update help text | Must |
| FR-6 | Update README | Must |

---

## Domain Concepts

### Key Entities
| Entity | Description | Attributes |
|--------|-------------|------------|
| ServeCommand | Main CLI command | port, dir, watch, verbose |
| VersionCommand | Version display | version, commit, date |
| Config | Runtime configuration | outputDir, port, watchEnabled |

### Key Operations
| Operation | Description | Inputs | Outputs |
|-----------|-------------|--------|---------|
| serve | Generate HTML + start server | flags | HTTP server |
| version | Display version info | none | stdout |

---

## Story Summary

| Metric | Count |
|--------|-------|
| Total Stories | 3 |
| Must Have | 3 |
| Should Have | 0 |
| Could Have | 0 |

### Stories

| Story ID | Title | Priority | Status |
|----------|-------|----------|--------|
| 001-remove-commands | Remove generate and watch commands | Must | Planned |
| 002-update-serve | Update serve with auto-generate and new flags | Must | Planned |
| 003-update-docs | Update README and help text | Must | Planned |

---

## Dependencies

### Depends On
| Unit | Reason |
|------|--------|
| None | Modifies existing CLI code from Intent 001 |

### Depended By
| Unit | Reason |
|------|--------|
| None | This is a leaf refactoring |

### External Dependencies
| System | Purpose | Risk |
|--------|---------|------|
| None | N/A | N/A |

---

## Technical Context

### Suggested Technology
- Go (existing)
- Cobra CLI framework (existing)

### Files to Modify

| File | Change |
|------|--------|
| `cmd/root.go` | Remove generate/watch command registration |
| `cmd/serve.go` | Add auto-generate on startup, rename flags |
| `cmd/generate.go` | DELETE |
| `cmd/watch.go` | DELETE |
| `README.md` | Update documentation |

### Integration Points
| Integration | Type | Protocol |
|-------------|------|----------|
| None | N/A | N/A |

---

## Constraints

- Must maintain backwards compatibility for `serve` command
- Must provide clear error if user tries old commands
- Single binary distribution must be preserved

---

## Success Criteria

### Functional
- [ ] `claude-code-logs serve` generates HTML then starts server
- [ ] `claude-code-logs serve --watch` enables file watching
- [ ] `claude-code-logs serve --dir /path` uses custom directory
- [ ] Default output is `~/claude-code-logs/`
- [ ] `claude-code-logs generate` shows helpful error
- [ ] `claude-code-logs watch` shows helpful error
- [ ] `claude-code-logs version` works unchanged

### Non-Functional
- [ ] No regression in existing functionality
- [ ] Help text is accurate and clear

### Quality
- [ ] All existing tests pass
- [ ] New flags are tested
- [ ] README matches actual CLI behavior

---

## Bolt Suggestions

| Bolt | Type | Stories | Objective |
|------|------|---------|-----------|
| 007-cli-refactor | simple-construction-bolt | All 3 | Complete CLI refactoring |

---

## Notes

This is a small, focused refactoring. All changes are localized to the CLI layer. The underlying parser, generator, server, and watcher logic remain unchanged.

Consider: Should old commands (`generate`, `watch`) fail silently, show deprecation warning, or show helpful migration message? Recommendation: Show helpful message like "Command removed. Use 'claude-code-logs serve' instead."
