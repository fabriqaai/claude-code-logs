---
unit: 001-parser
intent: 001-chat-log-viewer
phase: inception
status: draft
---

# Unit Brief: parser

## Purpose

Parse Claude Code JSONL session files from `~/.claude/projects/` into structured Go data types for downstream processing.

## Responsibility

- Discover all project folders in `~/.claude/projects/`
- Decode folder names to human-readable project paths
- List and parse `.jsonl` session files
- Extract messages, tool calls, and metadata
- Handle malformed lines gracefully

## Assigned Requirements

- **FR-1**: Project Discovery
- **FR-2**: JSONL Parsing

## Key Entities

### Types

```go
type Project struct {
    Path        string     // Original path (decoded from folder name)
    FolderName  string     // Encoded folder name
    Sessions    []Session
}

type Session struct {
    ID        string     // UUID from filename
    Summary   string     // From summary entry or first message
    Messages  []Message
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Message struct {
    Type      string     // "user", "assistant", "system"
    UUID      string
    ParentID  string
    Timestamp time.Time
    Content   []ContentBlock
    ToolCalls []ToolCall
}

type ContentBlock struct {
    Type string  // "text", "tool_use", "tool_result"
    Text string
    // ... tool-specific fields
}
```

## Key Operations

1. **DiscoverProjects()** - Scan `~/.claude/projects/`, return list of Project
2. **ParseSession(path)** - Parse single JSONL file into Session
3. **DecodeProjectPath(folderName)** - Convert `-Users-name-project` to `/Users/name/project`

## Dependencies

None - this is the foundational unit.

## Interface

Other units consume:
- `[]Project` - list of all discovered projects with sessions
- `Session` - parsed session with messages

## Technical Constraints

- Stream JSONL line-by-line (don't load entire file)
- Skip malformed lines with warning log
- Handle UTF-8 encoding
- Cross-platform path handling

## Success Criteria

- [ ] Discovers all projects in `~/.claude/projects/`
- [ ] Correctly decodes folder names to paths
- [ ] Parses all message types (user, assistant, system, summary)
- [ ] Extracts tool calls and results
- [ ] Handles 1000+ message sessions without OOM
- [ ] Graceful handling of malformed JSONL

---

## Story Summary

- **Total Stories**: 2
- **Must Have**: 2
- **Should Have**: 0
- **Could Have**: 0

### Stories

- [ ] **001-discover-projects**: Discover projects - Must - Planned
- [ ] **002-parse-sessions**: Parse session files - Must - Planned
