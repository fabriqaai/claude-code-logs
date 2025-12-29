---
intent: 001-chat-log-viewer
phase: inception
status: draft
created: 2025-12-29T10:30:00Z
updated: 2025-12-29T10:30:00Z
---

# Requirements: Chat Log Viewer

## Intent Overview

CLI tool that scans `~/.claude/projects`, parses JSONL chat logs, generates browsable HTML pages, and serves them with full-text search capability. Users can view static HTML files directly (without search) or run the server for the complete experience.

## Business Goals

| Goal | Success Metric | Priority |
|------|----------------|----------|
| Browse Claude Code chat history visually | Can navigate all projects/sessions in browser | Must |
| Search across conversations | Find messages by content within seconds | Must |
| Keep output updated automatically | Watch mode detects new sessions in < 1 min | Should |
| Zero-config quick start | `claude-logs serve` works immediately | Must |

---

## Functional Requirements

### FR-1: Project Discovery
- **Description**: Scan `~/.claude/projects/` and discover all project folders
- **Acceptance Criteria**:
  - Finds all folders matching pattern `-{encoded-path}`
  - Decodes folder names to human-readable project paths
  - Lists all `.jsonl` session files per project
- **Priority**: Must
- **Related Stories**: TBD

### FR-2: JSONL Parsing
- **Description**: Parse Claude Code session files (JSONL format)
- **Acceptance Criteria**:
  - Handles `user`, `assistant`, `system`, `summary` message types
  - Extracts message content, timestamps, tool calls
  - Preserves conversation threading via `parentUuid`
  - Handles malformed lines gracefully (skip with warning)
- **Priority**: Must
- **Related Stories**: TBD

### FR-3: HTML Generation
- **Description**: Generate static HTML pages for each project/session
- **Acceptance Criteria**:
  - Left navigation listing all projects (fixed sidebar, 288px width)
  - Session list per project with recent chats section
  - Conversation view with user/assistant messages styled differently
  - Tool calls rendered with syntax highlighting
  - Responsive design (works on mobile)
  - **Claude.ai visual style**: Match the look and feel of claude.ai web chat (see `ui-standards.md`)
    - Warm cream/beige backgrounds (not pure white)
    - Serif font for assistant responses, sans-serif for user messages
    - Subtle shadows and 0.5px borders
    - Smooth transitions (300ms cubic-bezier)
  - **Footer branding**: "claude-code-logs by fabriqa.ai" linking to https://fabriqa.ai
- **Priority**: Must
- **Related Stories**: TBD
- **Reference**: See `ui-standards.md` for complete CSS specifications

### FR-4: Server Mode
- **Description**: Serve generated HTML with search API
- **Acceptance Criteria**:
  - `claude-logs serve` starts HTTP server on configurable port
  - Serves static HTML files
  - Provides `/api/search` endpoint for full-text search
  - Displays server URL on startup
- **Priority**: Must
- **Related Stories**: TBD

### FR-5: Search Functionality
- **Description**: Full-text search across chat sessions
- **Acceptance Criteria**:
  - Search all sessions globally
  - Filter by project
  - Filter by session
  - Highlights matching text in results
  - Returns context (surrounding messages)
- **Priority**: Must
- **Related Stories**: TBD

### FR-6: Watch Mode
- **Description**: Monitor for new sessions and regenerate HTML
- **Acceptance Criteria**:
  - `claude-logs watch` monitors `~/.claude/projects/` for changes
  - Regenerates affected HTML within 60 seconds of change
  - Can run alongside server mode
  - Configurable poll interval
- **Priority**: Should
- **Related Stories**: TBD

### FR-7: Static Fallback UI
- **Description**: HTML files work when opened directly (without server)
- **Acceptance Criteria**:
  - Navigation and viewing works without server
  - Search UI shows "Start server for search" message
  - Displays the command to start server: `claude-logs serve`
  - No JavaScript errors when opened as file://
- **Priority**: Must
- **Related Stories**: TBD

### FR-8: CLI Commands
- **Description**: Provide intuitive CLI interface
- **Acceptance Criteria**:
  - `claude-logs generate` - one-time HTML generation
  - `claude-logs serve` - start server (default port 8080)
  - `claude-logs watch` - monitor and regenerate
  - `claude-logs serve --watch` - combined mode
  - `--output-dir` flag to specify output location
  - `--port` flag for server port
- **Priority**: Must
- **Related Stories**: TBD

---

## Non-Functional Requirements

### Performance
| Requirement | Metric | Target |
|-------------|--------|--------|
| Initial generation | Time for 100 sessions | < 10 seconds |
| Search response | p95 latency | < 500ms |
| Memory usage | Peak RAM | < 200MB |

### Scalability
| Requirement | Metric | Target |
|-------------|--------|--------|
| Session count | Total sessions | 1000+ |
| Message count | Messages per session | 10,000+ |

### Reliability
| Requirement | Metric | Target |
|-------------|--------|--------|
| Crash recovery | Watch mode auto-restart | Yes |
| Corrupt file handling | Skip and continue | Yes |

---

## Constraints

### Technical Constraints

**Project-wide standards**: Go, flat structure, standard testing (from standards/)

**Intent-specific constraints**:
- Must work offline (no external API calls)
- Single binary distribution (no runtime dependencies)
- Cross-platform (macOS, Linux)

### Business Constraints
- Personal/open-source project
- No external service dependencies

---

## Assumptions

| Assumption | Risk if Invalid | Mitigation |
|------------|-----------------|------------|
| Claude Code log format is stable | Parsing breaks on updates | Version detection, graceful degradation |
| Users have read access to ~/.claude | Tool fails silently | Clear error message with permission fix |
| JSONL files fit in memory line-by-line | Large sessions fail | Stream processing, never load full file |

---

## Open Questions

| Question | Owner | Due Date | Resolution |
|----------|-------|----------|------------|
| Index format for search (in-memory vs file-based)? | TBD | Before construction | Pending |
| Should we support syntax themes (light/dark)? | TBD | Before construction | Pending |
