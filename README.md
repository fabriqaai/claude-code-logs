# claude-code-logs

Browse and search your Claude Code chat logs with a local web interface.

## Features

- **HTML Generation**: Converts JSONL chat logs to browsable HTML pages
- **Full-Text Search**: Search across all messages with highlighted results
- **Tree View Sidebar**: Collapsible project/session tree with resizable width
- **Card-Based Layout**: Clean card views for projects and sessions
- **Breadcrumb Navigation**: Easy navigation back to projects from sessions
- **Source File Links**: Quick access to original JSONL files
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
# Start the server (generates HTML automatically)
claude-code-logs serve

# Open http://localhost:8080 in your browser
```

## Usage

### Serve

Generate HTML and start a local web server:

```bash
claude-code-logs serve                    # Default port 8080, outputs to ~/claude-code-logs
claude-code-logs serve --port 3000        # Custom port
claude-code-logs serve --dir /custom/path # Custom output directory
claude-code-logs serve --watch            # Auto-regenerate on changes
claude-code-logs serve --verbose          # Verbose output
```

### Version Info

```bash
claude-code-logs version
```

## Commands

| Command | Description |
|---------|-------------|
| `serve` | Generate HTML and start web server |
| `version` | Display version information |

## How It Works

1. **Scans** `~/.claude/projects/` for Claude Code chat sessions
2. **Parses** JSONL files containing conversation history
3. **Generates** HTML pages with syntax highlighting and formatting
4. **Serves** files locally with a search API

## Output Structure

```
~/claude-code-logs/
├── index.html              # Project listing
├── projects/
│   └── my-project/
│       ├── index.html      # Session listing
│       └── abc123.html     # Individual session
└── static/
    └── styles.css
```

## Configuration

| Flag | Short | Description | Default |
|------|-------|-------------|---------|
| `--dir` | `-d` | Output directory for HTML | `~/claude-code-logs` |
| `--port` | `-p` | Server port | `8080` |
| `--watch` | `-w` | Auto-regenerate on changes | `false` |
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
