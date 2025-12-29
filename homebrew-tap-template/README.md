# Homebrew Tap for FabriqaAI Tools

This repository contains Homebrew formulae for FabriqaAI tools.

## Installation

```bash
# Add the tap
brew tap fabriqa/tap

# Install claude-code-logs
brew install claude-code-logs
```

## Available Formulae

### claude-code-logs

Browse and search Claude Code chat logs.

```bash
brew install fabriqa/tap/claude-code-logs
```

**Usage:**

```bash
# Generate HTML from chat logs
claude-logs generate

# Start the local server
claude-logs serve

# Watch for changes and auto-regenerate
claude-logs watch
```

## Updating

```bash
brew update
brew upgrade claude-code-logs
```

## Troubleshooting

### Formula not found

```bash
brew tap fabriqa/tap
brew update
```

### SHA256 mismatch

```bash
brew update
brew reinstall claude-code-logs
```

## License

MIT
