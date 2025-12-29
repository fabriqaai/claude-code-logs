---
unit: 005-cli
intent: 001-chat-log-viewer
phase: inception
status: complete
---

# Unit Brief: cli

## Purpose

Command-line interface that orchestrates all components (parser, generator, server, watcher).

## Responsibility

- Define CLI commands and flags
- Wire up components
- Handle configuration
- Display progress and errors

## Assigned Requirements

- **FR-8**: CLI Commands

## Key Entities

### Commands

```bash
# One-time generation
claude-logs generate [--output-dir PATH]

# Start server (primary mode)
claude-logs serve [--port PORT] [--output-dir PATH]

# Start watcher (background regeneration)
claude-logs watch [--output-dir PATH] [--interval SECONDS]

# Combined: serve + watch
claude-logs serve --watch [--port PORT] [--output-dir PATH]

# Show version
claude-logs version
```

### Flags

| Flag | Default | Description |
|------|---------|-------------|
| `--output-dir` | `~/.claude-logs` | Output directory for HTML |
| `--port` | `8080` | Server port |
| `--interval` | `30` | Watch poll interval (seconds) |
| `--watch` | `false` | Enable watch mode with serve |
| `--verbose` | `false` | Verbose output |

### Configuration File (optional)

```yaml
# ~/.claude-logs.yaml
output_dir: ~/.claude-logs
port: 8080
watch_interval: 30
```

## Key Operations

1. **main()** - Entry point, parse commands
2. **cmdGenerate()** - One-time generation
3. **cmdServe()** - Start server
4. **cmdWatch()** - Start watcher
5. **loadConfig()** - Load configuration

## Dependencies

- **001-parser**: For parsing logs
- **002-generator**: For generating HTML
- **003-server**: For serving content
- **004-watcher**: For watching changes

## Interface

- User runs CLI commands
- CLI orchestrates components

## Technical Constraints

- Use Cobra for CLI framework
- Follow Go CLI conventions
- Clear error messages
- Non-zero exit code on errors
- Graceful shutdown (Ctrl+C)

## Success Criteria

- [ ] All commands work as documented
- [ ] Sensible defaults (works with just `claude-logs serve`)
- [ ] Clear error messages
- [ ] Graceful shutdown on Ctrl+C
- [ ] Help text for all commands
- [ ] Version command works

---

## Story Summary

- **Total Stories**: 1
- **Must Have**: 1
- **Should Have**: 0
- **Could Have**: 0

### Stories

- [ ] **001-cli-commands**: CLI commands - Must - Planned
