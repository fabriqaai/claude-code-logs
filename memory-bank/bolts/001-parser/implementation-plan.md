---
stage: plan
bolt: 001-parser
created: 2025-12-29T13:00:00Z
---

## Implementation Plan: parser

### Objective

Implement project discovery and JSONL session parsing to extract structured data from Claude Code log files in `~/.claude/projects/`.

### Deliverables

1. **types.go** - Data structures for Project, Session, Message, ContentBlock
2. **parser.go** - Core parsing functions:
   - `DiscoverProjects()` - Scan and list all projects
   - `DecodeProjectPath()` - Convert folder names to paths
   - `ParseSession()` - Parse JSONL into Session struct
   - `ListSessions()` - List all sessions for a project
3. **parser_test.go** - Unit tests for all parsing functions

### Dependencies

- Go standard library only:
  - `os` - File system operations
  - `bufio` - Line-by-line streaming
  - `encoding/json` - JSON parsing
  - `path/filepath` - Path manipulation
  - `strings` - String operations
  - `sort` - Sorting sessions by date
  - `time` - Timestamp handling

### Technical Approach

#### JSONL Structure (from real data analysis)

Each line is one of these types:
- `{"type":"summary","summary":"...","leafUuid":"..."}` - Session title
- `{"type":"user","message":{...},"uuid":"...","timestamp":"...","parentUuid":...}` - User message
- `{"type":"assistant","message":{...},"uuid":"...","timestamp":"...","parentUuid":"..."}` - Assistant response
- `{"type":"file-history-snapshot",...}` - Skip these

Message content blocks:
- `{"type":"text","text":"..."}` - Text content
- `{"type":"tool_use","id":"...","name":"...","input":{...}}` - Tool call
- `{"type":"tool_result","tool_use_id":"...","content":"..."}` - Tool result

#### Path Decoding

Folder name format: `-Users-cengiz-han-workspace-code-project`
Decoded path: `/Users/cengiz_han/workspace/code/project`

Algorithm:
1. Remove leading `-`
2. Split by `-`
3. Join with `/`
4. Prepend `/`

Note: Underscores in original path become underscores (preserved).

#### File Structure

```
claude-code-logs/
├── main.go           # Entry point (empty for now)
├── types.go          # Data structures
├── parser.go         # Parsing logic
├── parser_test.go    # Tests
├── go.mod
└── go.sum
```

### Acceptance Criteria

#### Story: 001-discover-projects
- [ ] Scans `~/.claude/projects/` directory
- [ ] Decodes folder names to human-readable paths
- [ ] Returns empty list (not error) for empty directory
- [ ] Returns error for missing `~/.claude/` directory
- [ ] Skips invalid folders with warning

#### Story: 002-parse-sessions
- [ ] Parses JSONL line-by-line (streaming)
- [ ] Extracts summary as session title
- [ ] Extracts user/assistant messages with content
- [ ] Captures tool_use and tool_result blocks
- [ ] Preserves parentUuid for threading
- [ ] Skips malformed lines with warning
- [ ] Sorts sessions by creation date descending

### Type Definitions

```go
type Project struct {
    Path       string    // Decoded path: /Users/name/project
    FolderName string    // Encoded: -Users-name-project
    Sessions   []Session
}

type Session struct {
    ID        string    // UUID from filename
    Summary   string    // From summary entry
    Messages  []Message
    CreatedAt time.Time
    UpdatedAt time.Time
}

type Message struct {
    UUID       string
    ParentUUID string
    Role       string    // "user" or "assistant"
    Content    []ContentBlock
    Timestamp  time.Time
}

type ContentBlock struct {
    Type       string // "text", "tool_use", "tool_result"
    Text       string // For text blocks
    ToolName   string // For tool_use
    ToolInput  string // JSON string of input
    ToolUseID  string // For tool_result
    ToolOutput string // For tool_result
}
```

### Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Large session files | Stream line-by-line, don't load all into memory |
| Malformed JSON | Skip bad lines, log warning, continue parsing |
| File locked by Claude | Retry once, then skip with warning |
| Unknown message types | Skip unknown types gracefully |
