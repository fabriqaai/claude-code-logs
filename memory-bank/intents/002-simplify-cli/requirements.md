---
intent: 002-simplify-cli
phase: complete
status: complete
created: 2025-12-29T14:00:00Z
updated: 2026-01-09T12:00:00Z
---

# Requirements: Simplify CLI

## Intent Overview

Simplify the CLI interface by reducing redundant commands from 4 to 2. The `serve` command will auto-generate HTML on startup and support watch mode via flag. Default output directory changes to `~/claude-code-logs/` with a `--dir` flag for customization.

## Business Goals

| Goal | Success Metric | Priority |
|------|----------------|----------|
| Reduce CLI complexity | Only 2 commands: `serve` and `version` | Must |
| Sensible defaults | Works with zero flags for common use case | Must |
| User-configurable output | `--dir` flag allows custom working directory | Must |

---

## Functional Requirements

### FR-1: Remove `generate` Command
- **Description**: Eliminate standalone `generate` command; auto-generate on `serve` startup
- **Acceptance Criteria**:
  - `claude-code-logs generate` no longer exists
  - `claude-code-logs serve` generates HTML before starting server
  - Generation happens automatically on every `serve` invocation
- **Priority**: Must
- **Related Stories**: TBD

### FR-2: Remove `watch` Command
- **Description**: Eliminate standalone `watch` command; use `serve --watch` instead
- **Acceptance Criteria**:
  - `claude-code-logs watch` no longer exists
  - `claude-code-logs serve --watch` enables file watching
  - Watch mode regenerates HTML when source files change
- **Priority**: Must
- **Related Stories**: TBD

### FR-3: Default Output Directory
- **Description**: Change default output directory to `~/claude-code-logs/`
- **Acceptance Criteria**:
  - Default output is `$HOME/claude-code-logs/` (not `~/.claude-code-logs/`)
  - Directory is created if it doesn't exist
  - Works on macOS and Linux
- **Priority**: Must
- **Related Stories**: TBD

### FR-4: Working Directory Flag
- **Description**: Add `--dir` flag to specify custom working directory
- **Acceptance Criteria**:
  - `claude-code-logs serve --dir /custom/path` uses specified directory
  - Short form `-d` available
  - Replaces `--output-dir` / `-o` flags
  - Directory is created if it doesn't exist
- **Priority**: Must
- **Related Stories**: TBD

### FR-5: Update Help Text
- **Description**: Update CLI help to reflect simplified commands
- **Acceptance Criteria**:
  - `claude-code-logs --help` shows only `serve` and `version`
  - `claude-code-logs serve --help` shows all available flags
  - Examples in help text are accurate
- **Priority**: Must
- **Related Stories**: TBD

### FR-6: Update README
- **Description**: Update documentation to reflect new CLI structure
- **Acceptance Criteria**:
  - README shows simplified command structure
  - Examples use new flag names (`--dir` instead of `--output-dir`)
  - Quick start section is updated
- **Priority**: Must
- **Related Stories**: TBD

---

## Non-Functional Requirements

### Backwards Compatibility
| Requirement | Metric | Target |
|-------------|--------|--------|
| Breaking change | Intentional | Yes - old commands will fail with helpful error |

### User Experience
| Requirement | Metric | Target |
|-------------|--------|--------|
| Time to first use | Commands to type | 1 (`claude-code-logs serve`) |
| Learning curve | Commands to remember | 2 |

---

## Constraints

### Technical Constraints

**Project-wide standards**: Go, Cobra CLI framework

**Intent-specific constraints**:
- Must maintain all existing functionality (generation, serving, watching)
- Only the CLI interface changes, not the underlying implementation

### Business Constraints
- Breaking change is acceptable (major version bump if needed)

---

## Assumptions

| Assumption | Risk if Invalid | Mitigation |
|------------|-----------------|------------|
| Users prefer simpler CLI | Some users may want standalone `generate` | Document workaround or add back if requested |
| `~/claude-code-logs/` is acceptable default | Users may prefer hidden directory | `--dir` flag provides override |

---

## Open Questions

| Question | Owner | Due Date | Resolution |
|----------|-------|----------|------------|
| Should old commands show deprecation warning or just fail? | TBD | Before construction | Pending |
