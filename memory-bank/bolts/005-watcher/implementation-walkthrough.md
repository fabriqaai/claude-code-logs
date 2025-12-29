---
stage: implement
bolt: 005-watcher
created: 2025-12-29T13:15:00Z
---

## Implementation Walkthrough: watcher

### Summary

Implemented a file watcher using fsnotify that monitors `~/.claude/projects/` for changes to `.jsonl` session files. The watcher detects creates, modifications, and deletions, debounces rapid changes using per-project timers, and triggers incremental HTML regeneration for only the affected project.

### Structure Overview

The watcher follows an event-driven architecture with three main components: the fsnotify listener for file system events, a debouncer that coalesces rapid changes using timers, and a regenerator callback that incrementally updates HTML. The implementation supports both standalone operation (`claude-code-logs watch`) and background operation with the server (`claude-code-logs serve --watch`).

### Completed Work

- [x] `watcher.go` - Core watcher implementation with fsnotify integration, debouncing, and regeneration callbacks
- [x] `cmd_watch.go` - Full implementation of standalone watch command with signal handling and graceful shutdown
- [x] `cmd_serve.go` - Integration of watch mode to run alongside the HTTP server
- [x] `go.mod` / `go.sum` - Added fsnotify v1.9.0 dependency

### Key Decisions

- **fsnotify over polling**: Chose event-driven fsnotify for efficiency rather than polling, as it provides immediate detection of changes
- **Per-project debouncing**: Each project has its own debounce timer, allowing independent regeneration without blocking other projects
- **Timer-based debouncing**: Used `time.AfterFunc` for clean debounce implementation that automatically fires after quiet period
- **Callback pattern for regeneration**: Watcher accepts a callback function, making it testable and allowing different regeneration strategies
- **Context-based cancellation**: Used Go contexts for clean shutdown propagation across goroutines
- **Incremental regeneration**: Only regenerates the affected project's pages plus the main index, not the entire site

### Deviations from Plan

- **No polling fallback**: Decided against implementing a polling fallback for v1 since fsnotify is well-supported on all major platforms (macOS, Linux, Windows)
- **File deletion handling deferred**: Full deletion cleanup (removing orphaned HTML files) deferred to future enhancement

### Dependencies Added

- [x] `github.com/fsnotify/fsnotify v1.9.0` - Cross-platform file system notification library
- [x] `golang.org/x/sys v0.13.0` - System call support (transitive dependency of fsnotify)

### Developer Notes

- The watcher watches both the root `~/.claude/projects/` directory AND each project subdirectory to catch all events
- New project directories are detected via periodic scanning (poll interval) since fsnotify can miss rapid directory creations
- The debounce delay should be at least 2 seconds since Claude Code can make multiple rapid writes to session files
- When running with `serve --watch`, the watcher runs in a background goroutine and automatically stops when the server shuts down
