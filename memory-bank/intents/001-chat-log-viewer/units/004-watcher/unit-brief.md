---
unit: 004-watcher
intent: 001-chat-log-viewer
phase: inception
status: draft
---

# Unit Brief: watcher

## Purpose

Monitor `~/.claude/projects/` for changes and trigger HTML regeneration when new sessions are created or updated.

## Responsibility

- Watch for new/modified `.jsonl` files
- Debounce rapid changes
- Trigger generator for affected projects
- Configurable poll interval

## Assigned Requirements

- **FR-6**: Watch Mode

## Key Entities

### Watch Events

```go
type WatchEvent struct {
    Type    string    // "created", "modified", "deleted"
    Path    string    // Full path to changed file
    Project string    // Project folder name
    Session string    // Session ID
}
```

### Configuration

```go
type WatchConfig struct {
    PollInterval  time.Duration  // Default: 30 seconds
    DebounceDelay time.Duration  // Default: 2 seconds
    SourceDir     string         // Default: ~/.claude/projects/
    OutputDir     string         // Default: ~/.claude-code-logs/
}
```

## Key Operations

1. **StartWatcher(config)** - Start watching for changes
2. **StopWatcher()** - Graceful shutdown
3. **OnChange(event)** - Handle file change event

## Dependencies

- **002-generator**: Called to regenerate HTML
- **001-parser**: Called to re-parse changed sessions

## Interface

- CLI starts watcher with config
- Watcher calls generator when changes detected

## Technical Constraints

- Use `fsnotify` for cross-platform file watching
- Fall back to polling if fsnotify unavailable
- Debounce rapid changes (Claude Code writes frequently)
- Only watch `.jsonl` files
- Efficient - don't regenerate unchanged files

## Success Criteria

- [ ] Detects new session files within poll interval
- [ ] Detects modified session files
- [ ] Debounces rapid changes (no duplicate regeneration)
- [ ] Only regenerates affected projects
- [ ] Graceful shutdown
- [ ] Works alongside server mode

---

## Story Summary

- **Total Stories**: 1
- **Must Have**: 0
- **Should Have**: 1
- **Could Have**: 0

### Stories

- [ ] **001-watch-changes**: Watch for changes - Should - Planned
