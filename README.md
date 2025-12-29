# claude-code-logs

Browse and search your Claude Code chat logs with a local web interface.

## Features

- **HTML Generation**: Converts JSONL chat logs to browsable HTML pages
- **Full-Text Search**: Search across all messages with highlighted results
- **File Watching**: Auto-regenerate when chat logs change
- **Local Server**: Browse your logs at `http://localhost:8080`
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
# Generate HTML from your chat logs
claude-logs generate

# Start the web server
claude-logs serve

# Open http://localhost:8080 in your browser
```

## Usage

### Generate HTML

Convert Claude Code chat logs to static HTML pages:

```bash
claude-logs generate
claude-logs generate --output-dir /custom/path
claude-logs generate --verbose
```

### Start Server

Launch a local web server to browse your logs:

```bash
claude-logs serve                    # Default port 8080
claude-logs serve --port 3000        # Custom port
claude-logs serve --watch            # Auto-regenerate on changes
```

### Watch Mode

Monitor for changes and regenerate automatically:

```bash
claude-logs watch                    # Watch with defaults
claude-logs watch --debounce 5       # 5 second debounce delay
claude-logs watch --interval 60      # Check for new dirs every 60s
```

### Version Info

```bash
claude-logs version
```

## How It Works

1. **Scans** `~/.claude/projects/` for Claude Code chat sessions
2. **Parses** JSONL files containing conversation history
3. **Generates** HTML pages with syntax highlighting and formatting
4. **Serves** files locally with a search API

## Output Structure

```
~/.claude-logs/
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
| `--output-dir` | `-o` | Output directory for HTML | `~/.claude-logs` |
| `--verbose` | `-v` | Verbose output | `false` |
| `--port` | `-p` | Server port (serve only) | `8080` |
| `--watch` | `-w` | Enable watch mode (serve only) | `false` |
| `--debounce` | `-d` | Debounce delay in seconds (watch only) | `2` |
| `--interval` | `-i` | Poll interval in seconds (watch only) | `30` |

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
go build -o claude-logs .

# Run locally
./claude-logs generate
./claude-logs serve
```

## License

MIT

## Links

- [GitHub Repository](https://github.com/fabriqaai/claude-code-logs)
- [Homebrew Tap](https://github.com/fabriqaai/homebrew-tap)
- [Releases](https://github.com/fabriqaai/claude-code-logs/releases)
