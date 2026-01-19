# System Architecture

## Overview

claude-code-logs is a CLI tool that parses Claude Code's JSONL chat session files, generates browsable Markdown documentation, and serves them through a local web interface with full-text search.

## System Context

The tool bridges Claude Code's raw chat logs with a human-friendly browsable interface. It reads from Claude Code's data directory and serves content locally.

### Context Diagram

```
┌─────────────────┐     ┌──────────────────────┐     ┌─────────────────┐
│   Claude Code   │────>│   claude-code-logs   │────>│    Browser      │
│  (~/.claude/)   │     │      CLI Tool        │     │  (localhost)    │
└─────────────────┘     └──────────────────────┘     └─────────────────┘
      JSONL                   Go Binary               HTML/Markdown
```

### Users

- **Developer**: Primary user viewing their Claude Code chat history
- **Team Lead**: May review team members' AI interactions (if shared)

### External Systems

- **Claude Code**: Source of JSONL chat session files
- **Browser**: Renders the served HTML/Markdown content

## Architecture Pattern

**Pattern**: Pipeline + HTTP Server
**Rationale**: Data flows through parsing → generation → serving stages. Simple, stateless design fits CLI tool use case.

## Component Architecture

### Components

#### Parser

- **Purpose**: Extract structured data from JSONL files
- **Responsibilities**: Parse messages, handle content blocks, extract metadata
- **Dependencies**: Standard library only

#### Generator

- **Purpose**: Convert parsed sessions to Markdown files
- **Responsibilities**: Generate frontmatter, format content, write files atomically
- **Dependencies**: Parser output

#### Server

- **Purpose**: Serve generated content via HTTP
- **Responsibilities**: Route requests, render templates, handle API
- **Dependencies**: Generator output, templates

#### Search

- **Purpose**: Provide full-text search over sessions
- **Responsibilities**: Build inverted index, execute queries, highlight matches
- **Dependencies**: Parser output

#### Watcher

- **Purpose**: Detect changes in source files
- **Responsibilities**: Monitor file system, trigger regeneration
- **Dependencies**: fsnotify

### Component Diagram

```
┌─────────────────────────────────────────────────────────┐
│                    CLI (cobra)                          │
└─────────────────────────────────────────────────────────┘
          │              │              │
          v              v              v
    ┌──────────┐   ┌──────────┐   ┌──────────┐
    │  Parser  │──>│ Generator│──>│  Server  │
    └──────────┘   └──────────┘   └──────────┘
          │                             │
          v                             v
    ┌──────────┐                 ┌──────────┐
    │  Search  │                 │ Templates│
    └──────────┘                 └──────────┘
          ^
          │
    ┌──────────┐
    │ Watcher  │
    └──────────┘
```

## Data Flow

Sessions are parsed from JSONL, indexed for search, converted to Markdown, and served as HTML.

```
~/.claude/projects/**/*.jsonl
          │
          v
    ┌──────────┐
    │  Parse   │──> []Session
    └──────────┘
          │
    ┌─────┴─────┐
    v           v
┌────────┐ ┌────────┐
│ Index  │ │Generate│──> *.md files
└────────┘ └────────┘
    │           │
    v           v
┌────────────────────┐
│   HTTP Server      │──> Browser
│ /api/search        │
│ /{project}/{id}    │
└────────────────────┘
```

## Technology Stack

| Layer | Technology | Purpose |
|-------|------------|---------|
| CLI | Cobra | Command parsing, flags, subcommands |
| Parsing | encoding/json | JSONL line parsing |
| Templates | html/template | Server-side HTML rendering |
| HTTP | net/http | Local web server |
| Search | Custom inverted index | Full-text search |
| File Watch | fsnotify | Real-time file monitoring |

## Non-Functional Requirements

### Performance

- **Startup time**: < 2s for 1000 sessions
- **Search latency**: < 100ms for typical queries
- **Memory**: < 200MB for 10K indexed messages

### Security

- Server binds to localhost only (127.0.0.1)
- No authentication (local-only access)
- No external network calls
- Atomic file writes prevent corruption

### Scalability

Designed for single-user local use. Performance scales linearly with session count. Search index held in memory limits practical size to ~100K messages.

## Constraints

- Must run without external dependencies (single binary)
- Must work offline (no network required)
- Must support macOS and Linux
- Must preserve original JSONL files (read-only)

## Key Decisions

| Decision | Choice | Rationale |
|----------|--------|-----------|
| Language | Go | Single binary, cross-platform, good stdlib |
| Search | In-memory inverted index | Simple, fast for local use case |
| Rendering | Server-side templates + client-side Markdown | Fast initial load, rich formatting |
| Storage | Filesystem (Markdown) | Human-readable, git-friendly, no DB needed |

---
*Generated by FIRE project-init*
