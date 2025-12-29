---
id: 003-update-docs
unit: 001-cli-refactor
intent: 002-simplify-cli
status: ready
priority: must
created: 2025-12-29T14:15:00Z
assigned_bolt: 007-cli-refactor
implemented: false
---

# Story: 003-update-docs

## User Story

**As a** new user
**I want** accurate documentation
**So that** I can learn how to use the simplified CLI correctly

## Acceptance Criteria

- [ ] **Given** I view README.md, **When** reading Quick Start, **Then** I see only `claude-logs serve` (not generate + serve)
- [ ] **Given** I view README.md, **When** reading Commands Overview, **Then** I see only 2 commands: `serve` and `version`
- [ ] **Given** I view README.md, **When** reading Configuration table, **Then** `--dir`/`-d` is documented (not `--output-dir`/`-o`)
- [ ] **Given** I view README.md, **When** reading default output, **Then** it shows `~/claude-code-logs/` (not `~/.claude-logs/`)
- [ ] **Given** I run `claude-logs serve --help`, **When** viewing help, **Then** all flags are accurately described
- [ ] **Given** I run `claude-logs --help`, **When** viewing root help, **Then** only serve and version are listed

## Technical Notes

### README Changes

**Quick Start section** (before):
```bash
claude-logs generate
claude-logs serve
```

**Quick Start section** (after):
```bash
claude-logs serve
# Open http://localhost:8080 in your browser
```

**Commands Overview** (after):
| Command | Description |
|---------|-------------|
| `serve` | Generate HTML and start web server |
| `version` | Display version information |

**Configuration table** - update flags:
| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dir` | `-d` | Output directory for HTML | `~/claude-code-logs` |
| `--port` | `-p` | Server port | `8080` |
| `--watch` | `-w` | Auto-regenerate on changes | `false` |
| `--verbose` | `-v` | Verbose output | `false` |

## Dependencies

### Requires
- 002-update-serve (need final flag names)

### Enables
- None (final story)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Old documentation cached | Users see version in README |
| Help text line length | Keep under 80 chars for terminal |

## Out of Scope

- Creating new documentation files
- Updating external documentation (Homebrew tap README)
