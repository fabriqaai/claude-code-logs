---
stage: plan
bolt: 004-cli
created: 2025-12-29T16:35:00Z
---

## Implementation Plan: CLI

### Objective

Implement a Cobra-based CLI that orchestrates all existing components (parser, generator, server) and provides user-facing commands for generating, serving, and watching Claude Code chat logs.

### Deliverables

- `cmd/` directory with Cobra command structure (or flat structure in main package)
- `generate` command - one-time HTML generation
- `serve` command - HTTP server with optional watch mode
- `watch` command - standalone watcher (for future bolt 005)
- `version` command - display version info
- Updated `main.go` as CLI entry point
- All commands with proper flags, help text, and error handling

### Dependencies

- **cobra** (github.com/spf13/cobra): Industry-standard Go CLI framework
- **parser.go**: `DefaultClaudeProjectsPath()`, `LoadAllProjects()`
- **generator.go**: `GenerateAll()`
- **server.go**: `NewServer()`, `Start()`
- **Note**: Watcher (bolt 005) not yet implemented - `watch` command will be a placeholder

### Technical Approach

**Command Structure**:
```
claude-logs
├── generate   # One-time HTML generation
├── serve      # Start HTTP server (with optional --watch)
├── watch      # Standalone watcher (placeholder for bolt 005)
└── version    # Show version
```

**Shared Flags**:
- `--output-dir` (default: `~/.claude-logs`) - Output directory for HTML
- `--verbose` (default: false) - Verbose output

**Command-Specific Flags**:
- `serve --port` (default: 8080) - Server port
- `serve --watch` (default: false) - Enable watch mode
- `watch --interval` (default: 30) - Poll interval in seconds

**Implementation Notes**:

1. **Flat Package Structure**: Keep commands in main package per coding standards
2. **Signal Handling**: Graceful shutdown on SIGINT/SIGTERM (already in server.go)
3. **Error Messages**: Clear, actionable error messages with exit code 1
4. **Version Injection**: Via ldflags at build time (existing `Version` variable)
5. **Default Output Dir**: Expand `~` to user home directory
6. **Watch Placeholder**: `watch` command prints "not yet implemented" until bolt 005

**File Changes**:

| File | Action | Description |
|------|--------|-------------|
| `main.go` | Rewrite | Replace placeholder with Cobra root command |
| `go.mod` | Update | Add cobra dependency |
| `cmd_generate.go` | Create | Generate command implementation |
| `cmd_serve.go` | Create | Serve command implementation |
| `cmd_watch.go` | Create | Watch command placeholder |
| `cmd_version.go` | Create | Version command implementation |

### Acceptance Criteria

- [ ] `claude-logs generate` generates HTML in output directory
- [ ] `claude-logs serve` starts server on port 8080
- [ ] `claude-logs serve --port 3000` uses custom port
- [ ] `claude-logs serve --watch` enables watch mode (placeholder message for now)
- [ ] `claude-logs watch` shows "not yet implemented" message
- [ ] `claude-logs version` displays version
- [ ] `--output-dir` flag works on all commands
- [ ] Ctrl+C triggers graceful shutdown
- [ ] Invalid flags show help message
- [ ] Missing dependencies show clear error
- [ ] All commands have help text (`--help`)
- [ ] Non-zero exit code on errors

### Risk Assessment

| Risk | Mitigation |
|------|------------|
| Watch not implemented | Placeholder command, defer to bolt 005 |
| Port conflicts | Server already handles this gracefully |
| Output dir permissions | Check writability before generation |

### Out of Scope

- Watcher implementation (bolt 005)
- Config file support
- Shell completion
- Interactive mode
