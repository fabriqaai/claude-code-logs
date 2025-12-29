---
stage: plan
bolt: 005-watcher
created: 2025-12-29T13:00:00Z
---

## Implementation Plan: watcher

### Objective

Implement a file watcher that monitors `~/.claude/projects/` for changes to `.jsonl` files and triggers incremental HTML regeneration. The watcher will use fsnotify for event-driven watching with debouncing to handle rapid file changes efficiently.

### Deliverables

- `watcher.go` - Core watcher implementation with fsnotify and debouncing
- Updated `cmd_watch.go` - Full implementation of standalone watch command
- Updated `cmd_serve.go` - Integration of watch mode with server
- Updated `go.mod` - fsnotify dependency

### Dependencies

- `github.com/fsnotify/fsnotify` - Cross-platform file system notifications
- `001-parser` - Re-parse changed sessions via `ParseSession()`
- `002-generator` - Regenerate HTML via `GenerateProjectPages()`, `GenerateIndex()`

### Technical Approach

#### Architecture

```
Watcher
├── WatchConfig (poll interval, debounce delay, directories)
├── fsnotify.Watcher (file system events)
├── Debouncer (coalesce rapid changes)
├── Regenerator (selective HTML updates)
└── Signal Handler (graceful shutdown)
```

#### Key Components

1. **Watcher struct** - Main orchestrator
   - Holds fsnotify watcher instance
   - Manages goroutines for event processing
   - Handles graceful shutdown via context

2. **Debouncer** - Coalesces rapid file changes
   - Uses timer-based approach (2 second default)
   - Tracks pending changes per project
   - Fires regeneration after quiet period

3. **Event Handler** - Processes fsnotify events
   - Filters for `.jsonl` files only
   - Maps file path to project
   - Queues project for regeneration

4. **Regenerator** - Selective HTML updates
   - Re-parses only changed sessions
   - Regenerates only affected project HTML
   - Updates main index.html

#### Event Flow

```
File Change → fsnotify Event → Filter (.jsonl) → Debouncer → Regenerate Project
```

#### Concurrency Model

- Main goroutine: Handles signals, coordinates shutdown
- Event goroutine: Processes fsnotify events, feeds debouncer
- Debounce goroutine: Manages timers, triggers regeneration
- Use `context.Context` for cancellation propagation

### Acceptance Criteria

- [ ] Detects new `.jsonl` files within poll interval
- [ ] Detects modified `.jsonl` files and regenerates HTML
- [ ] Debounces rapid changes (< 2 seconds apart) into single regeneration
- [ ] Only regenerates affected project's HTML (not full site)
- [ ] Graceful shutdown on Ctrl+C (SIGINT/SIGTERM)
- [ ] Configurable poll interval via `--interval` flag
- [ ] Works standalone (`claude-code-logs watch`)
- [ ] Works with server (`claude-code-logs serve --watch`)

### File Structure

```
.
├── watcher.go          # Core watcher implementation
├── cmd_watch.go        # Updated watch command (full impl)
├── cmd_serve.go        # Updated serve command (watch integration)
└── go.mod              # Add fsnotify dependency
```

### API Design

```go
// WatchConfig configures the file watcher
type WatchConfig struct {
    SourceDir     string        // ~/.claude/projects/
    OutputDir     string        // ~/.claude-code-logs/
    PollInterval  time.Duration // Default: 30s (for new directory scanning)
    DebounceDelay time.Duration // Default: 2s
}

// Watcher monitors for file changes
type Watcher struct {
    config    WatchConfig
    fsWatcher *fsnotify.Watcher
    // ...
}

// StartWatcher starts watching for changes (blocking)
func StartWatcher(ctx context.Context, config WatchConfig) error

// WatchInBackground starts watcher in goroutine (for serve --watch)
func WatchInBackground(ctx context.Context, config WatchConfig) error
```

### Edge Cases Handling

| Scenario | Approach |
|----------|----------|
| File deleted | Remove corresponding session HTML, regenerate project index |
| New project folder | Add recursive watch, generate project pages |
| fsnotify unavailable | Log error, exit (no polling fallback for v1) |
| Permission denied | Log warning, continue watching other files |
| Disk full | Log error, retry on next change |
| Concurrent changes | Debouncer coalesces into single regeneration |

### Implementation Notes

1. **fsnotify recursive watching** - fsnotify doesn't support recursive watching natively. We need to:
   - Watch the root projects directory for new project folders
   - Add watches for each project folder
   - Handle `Create` events for new directories

2. **Debouncing strategy** - Use a map of project → timer:
   - On file change, reset/create timer for that project
   - When timer fires (after 2s of quiet), regenerate that project
   - This ensures rapid saves don't cause multiple regenerations

3. **Graceful shutdown** - Use `context.WithCancel`:
   - Catch SIGINT/SIGTERM
   - Cancel context to stop all goroutines
   - Wait for goroutines to finish
   - Close fsnotify watcher

4. **Server integration** - For `serve --watch`:
   - Start watcher in background goroutine
   - Share context with server
   - When server stops, watcher stops too

### Out of Scope (v1)

- Polling fallback (fsnotify only)
- Browser live reload (no WebSocket)
- Custom watch directories
- File deletion cleanup (deferred)
