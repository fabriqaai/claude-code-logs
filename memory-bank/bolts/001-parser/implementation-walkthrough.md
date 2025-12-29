---
stage: implement
bolt: 001-parser
created: 2025-12-29T13:15:00Z
---

## Implementation Walkthrough: parser

### Summary

Implemented the core parser module for Claude Code logs. The parser discovers projects in `~/.claude/projects/`, decodes folder names to human-readable paths, and parses JSONL session files into structured Go types with support for text content, tool calls, and threading.

### Structure Overview

Flat package structure with three source files following Go conventions. All code is in the `main` package for this CLI tool. Standard library only - no external dependencies.

### Completed Work

- [x] `go.mod` - Go module definition for github.com/fabriqaai/claude-code-logs
- [x] `types.go` - Data structures for Project, Session, Message, ContentBlock
- [x] `parser.go` - Core parsing functions (DiscoverProjects, ParseSession, DecodeProjectPath)
- [x] `parser_test.go` - Comprehensive unit tests with table-driven tests
- [x] `main.go` - Placeholder entry point for testing parser

### Key Decisions

- **Standard library only**: No external dependencies for JSONL parsing - uses `encoding/json` with `bufio.Scanner` for streaming
- **Streaming parsing**: Uses `bufio.Scanner` with 10MB buffer to handle large session files line-by-line without loading into memory
- **Graceful error handling**: Malformed JSON lines are skipped with warnings rather than failing the entire parse
- **Flexible content parsing**: Handles both string content and array of content blocks
- **Null-safe**: Uses pointer for `parentUuid` to handle JSON null values

### Deviations from Plan

- Added `LoadProjectWithSessions` and `LoadAllProjects` convenience functions not in original plan
- Used `*string` pointer for ParentUUID to properly handle JSON null values
- Added 10MB buffer for Scanner to handle very long lines in real JSONL files

### Dependencies Added

None - standard library only as planned.

### Developer Notes

- Path decoding algorithm: The folder `-Users-name-project` becomes `/Users/name/project`. Note that underscores in original paths are preserved (not affected by the dash-to-slash conversion).
- JSONL entry types: `summary`, `user`, `assistant` are processed. `file-history-snapshot` and unknown types are silently skipped.
- Timestamp parsing: Tries RFC3339 first, then falls back to `2006-01-02T15:04:05.000Z` format used by Claude.
- Tool calls: `tool_use` blocks capture name, id, and input JSON. `tool_result` blocks capture tool_use_id and output.
