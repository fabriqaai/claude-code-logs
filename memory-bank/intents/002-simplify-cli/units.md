---
intent: 002-simplify-cli
phase: inception
status: units-decomposed
updated: 2025-12-29T14:10:00Z
---

# Simplify CLI - Unit Decomposition

## Units Overview

This intent decomposes into **1 unit** of work. This is a focused refactoring effort that touches only the CLI layer.

## Requirement-to-Unit Mapping

| FR | Requirement | Unit |
|----|-------------|------|
| FR-1 | Remove `generate` command | 001-cli-refactor |
| FR-2 | Remove `watch` command | 001-cli-refactor |
| FR-3 | Default output directory | 001-cli-refactor |
| FR-4 | Working directory flag | 001-cli-refactor |
| FR-5 | Update help text | 001-cli-refactor |
| FR-6 | Update README | 001-cli-refactor |

---

### Unit 1: 001-cli-refactor

**Description**: Refactor CLI commands from 4 to 2, update defaults and flags

**Stories**:
- Story-1: Remove generate and watch commands
- Story-2: Update serve command with auto-generate and new flags
- Story-3: Update documentation

**Deliverables**:
- Modified `cmd/root.go` - remove generate/watch commands
- Modified `cmd/serve.go` - add auto-generate, update flags
- Removed `cmd/generate.go` and `cmd/watch.go`
- Updated `README.md`

**Dependencies**:
- Depends on: None (modifies existing code from Intent 001)
- Depended by: None

**Estimated Complexity**: S (Small)

---

## Unit Dependency Graph

```text
[001-cli-refactor] (standalone - no dependencies)
```

## Execution Order

1. Single bolt: cli-refactor (all changes in one focused session)
