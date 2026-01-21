# claude-code-logs

Browse and search your Claude Code chat logs with a local web interface.

<!-- Last updated: 2026-01-21T17:30:00Z -->

Blog article about this project https://www.cengizhan.com/p/announcing-claude-code-logs-a-searchable

## Features

- **Markdown-First**: Generates Markdown files with YAML frontmatter for easy archival and version control
- **Server-Side Rendering**: HTML pages rendered at runtime with caching for fast consecutive requests
- **Client-Side Rendering**: Markdown content rendered in browser using marked.js + highlight.js
- **Full-Text Search**: Search across all messages with highlighted results and full-page overlay
- **Inline Search**: Filter messages within a session with real-time highlighting
- **Hide Tool Calls**: Toggle to hide tool calls for a compact conversation view
- **Tree View Sidebar**: Collapsible project/session tree with resizable width
- **Card-Based Layout**: Clean card views for projects and sessions
- **Download & Copy**: Download or copy session content as Markdown
- **Copy JSONL Path**: Quick copy of source JSONL file path to clipboard
- **File Watching**: Auto-regenerate when chat logs change
- **Local Server**: Browse your logs at `http://localhost:8080`
- **Mobile Responsive**: Works on desktop and mobile devices
- **Cross-Platform**: macOS (Intel & Apple Silicon) and Linux

## Installation

### Homebrew (macOS/Linux)

```bash
brew tap fabriqaai/tap
brew install claude-code-logs
```

### From Source

```bash
go install github.com/fabriqaai/claude-code-logs@latest
```

### Manual Download

Download the latest release from [GitHub Releases](https://github.com/fabriqaai/claude-code-logs/releases).

## Quick Start

```bash
# Start the server (generates Markdown automatically)
claude-code-logs serve

# Open http://localhost:8080 in your browser
```

## Usage

### Serve

Generate Markdown and start a local web server:

```bash
claude-code-logs serve                      # Default port 8080, outputs to ~/claude-code-logs
claude-code-logs serve --port 3000          # Custom port
claude-code-logs serve --dir /custom/path   # Custom output directory
claude-code-logs serve --watch              # Auto-regenerate on changes
claude-code-logs serve --list               # Interactively select projects
claude-code-logs serve --force              # Force regeneration (ignore mtime)
claude-code-logs serve --verbose            # Verbose output
```

### Version Info

```bash
claude-code-logs version
```

## Commands

| Command | Description |
|---------|-------------|
| `serve` | Generate Markdown and start web server |
| `version` | Display version information |

## How It Works

1. **Scans** `~/.claude/projects/` for Claude Code chat sessions
2. **Parses** JSONL files containing conversation history
3. **Generates** Markdown files with YAML frontmatter (source, hash, project, title, created)
4. **Serves** HTML pages rendered at runtime with client-side Markdown rendering
5. **Provides** search API for full-text search across all messages

## Output Structure

```
~/claude-code-logs/
├── index.md                    # Main project listing
├── my-project/
│   ├── index.md                # Session listing for project
│   ├── abc123.md               # Session Markdown with frontmatter
│   └── def456.md               # Another session
└── another-project/
    └── ...
```

### Markdown Format

Each session is saved as Markdown with YAML frontmatter:

```markdown
---
source: my-session.jsonl
source_hash: sha256:abc123...
project: /Users/me/my-project
title: Session Summary
created: 2024-01-17T10:30:00Z
---

## User

Hello, can you help me?

## Assistant

Of course! How can I help?
```

## Configuration

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dir` | `-d` | Output directory for Markdown | `~/claude-code-logs` |
| `--port` | `-p` | Server port | `8080` |
| `--watch` | `-w` | Auto-regenerate on changes | `false` |
| `--list` | `-l` | Interactively select projects | `false` |
| `--force` | `-f` | Force regeneration (ignore mtime) | `false` |
| `--verbose` | `-v` | Verbose output | `false` |

## Requirements

- Claude Code must be installed and have generated chat logs
- Chat logs are stored in `~/.claude/projects/`

## Development

```bash
# Clone the repository
git clone https://github.com/fabriqaai/claude-code-logs.git
cd claude-code-logs

# Run tests
go test -v ./...

# Build
go build -o claude-code-logs .

# Run locally
./claude-code-logs serve
```

## License

MIT

## Links

- [GitHub Repository](https://github.com/fabriqaai/claude-code-logs)
- [Homebrew Tap](https://github.com/fabriqaai/homebrew-tap)
- [Releases](https://github.com/fabriqaai/claude-code-logs/releases)
