---
id: 002-parse-sessions
unit: 001-parser
intent: 001-chat-log-viewer
status: draft
priority: must
created: 2025-12-29T12:00:00Z
assigned_bolt: 001-parser
implemented: false
---

# Story: 002-parse-sessions

## User Story

**As a** developer
**I want** my Claude Code sessions parsed correctly
**So that** I can view all messages and tool calls

## Acceptance Criteria

- [ ] **Given** a valid `.jsonl` session file, **When** parsing, **Then** all messages are extracted with correct roles (user/assistant/system)
- [ ] **Given** a session with tool calls, **When** parsing, **Then** tool_use and tool_result blocks are captured
- [ ] **Given** a session with summary entry, **When** parsing, **Then** summary is extracted as session title
- [ ] **Given** a malformed JSON line, **When** parsing, **Then** line is skipped with warning, parsing continues
- [ ] **Given** a 1000+ message session, **When** parsing, **Then** completes without memory issues (streaming)
- [ ] **Given** multiple sessions in a project, **When** listing, **Then** sessions sorted by creation date descending
- [ ] **Given** message with parentUuid, **When** parsing, **Then** threading relationship is preserved

## Technical Notes

- JSONL format: One JSON object per line
- Message types: `human`, `assistant`, `summary`
- Content blocks: `text`, `tool_use`, `tool_result`
- Stream line-by-line using `bufio.Scanner`
- Extract timestamps from message UUIDs or use file mtime
- Handle nested content blocks (array of blocks)

## Dependencies

### Requires
- 001-discover-projects (provides project paths)

### Enables
- Generator unit (002) - consumes parsed sessions
- Server unit (003) - indexes parsed content for search

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| Empty `.jsonl` file | Return Session with empty Messages |
| File locked by Claude Code | Retry once, then skip with warning |
| Unicode content (emoji, CJK) | Parse correctly, preserve encoding |
| Very long message (100KB+) | Parse fully, no truncation |
| Binary content in text field | Preserve as-is (base64 encoded) |
| Corrupted file mid-parse | Return partial results, log error |

## Out of Scope

- Rendering message content (generator responsibility)
- Message filtering or search (server responsibility)
- Real-time watching of sessions (watcher responsibility)
