---
id: 001-watch-changes
unit: 004-watcher
intent: 001-chat-log-viewer
status: draft
priority: should
created: 2025-12-29T12:00:00Z
assigned_bolt: 005-watcher
implemented: false
---

# Story: 001-watch-changes

## User Story

**As a** developer
**I want** HTML auto-regenerated when I have new sessions
**So that** my viewer stays up to date

## Acceptance Criteria

- [ ] **Given** watcher started, **When** new `.jsonl` file created, **Then** detected within poll interval
- [ ] **Given** watcher started, **When** existing `.jsonl` modified, **Then** detected and HTML regenerated
- [ ] **Given** rapid file changes (< 2 seconds apart), **When** detected, **Then** debounced (single regeneration)
- [ ] **Given** file change detected, **When** regenerating, **Then** only affected project's HTML updated
- [ ] **Given** watcher running, **When** Ctrl+C pressed, **Then** graceful shutdown
- [ ] **Given** `--interval` flag, **When** starting watcher, **Then** poll interval is configurable

## Technical Notes

- Use `fsnotify` for cross-platform file watching
- Fall back to polling if fsnotify unavailable
- Default poll interval: 30 seconds
- Default debounce delay: 2 seconds
- Only watch `.jsonl` files in `~/.claude/projects/`
- Use goroutine for background watching

## Dependencies

### Requires
- 002-generator (for regenerating HTML)
- 001-parser (for re-parsing changed sessions)

### Enables
- 005-cli (exposes watch command and --watch flag)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| File deleted | Remove corresponding HTML |
| New project folder created | Scan and add to index |
| fsnotify not available | Fall back to polling |
| File renamed | Treat as delete + create |
| Permission denied mid-watch | Log warning, continue watching others |
| Disk full during regeneration | Log error, retry next interval |

## Out of Scope

- Notifying browser of changes (no WebSocket)
- Watching custom directories (only ~/.claude/projects)
- Real-time streaming of active sessions
