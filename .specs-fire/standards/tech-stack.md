# Tech Stack

## Language

- **Go 1.21** - Primary language

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML parsing |
| `github.com/fsnotify/fsnotify` | File watching |
| `github.com/AlecAivazis/survey/v2` | Interactive prompts |

## Standard Library

- `net/http` - HTTP server
- `html/template` - HTML templating
- `encoding/json` - JSON parsing
- `bufio` - JSONL line scanning

## Build & Release

- **GoReleaser** - Cross-platform builds
- **Homebrew** - macOS/Linux distribution

## Data Sources

- Input: `~/.claude/projects/**/*.jsonl` (Claude Code chat logs)
- Output: `~/claude-code-logs/*.md` (Markdown with YAML frontmatter)
