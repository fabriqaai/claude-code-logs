---
unit: 006-homebrew-tap
intent: 001-chat-log-viewer
phase: complete
status: complete
---

# Unit Brief: homebrew-tap

## Purpose

Create and maintain `fabriqa/tap` Homebrew tap for easy installation on macOS and Linux.

## Responsibility

- Create `homebrew-tap` repository at `fabriqaai/homebrew-tap`
- Define formula for `claude-code-logs`
- Automate formula updates on release
- Document installation instructions

## Assigned Requirements

- Distribution requirement (Homebrew)

## Key Entities

### Repository Structure

```
fabriqaai/homebrew-tap/
├── Formula/
│   └── claude-code-logs.rb
└── README.md
```

### Formula Template

```ruby
class ClaudeCodeLogs < Formula
  desc "Browse and search Claude Code chat logs"
  homepage "https://github.com/fabriqaai/claude-code-logs"
  url "https://github.com/fabriqaai/claude-code-logs/releases/download/v#{version}/claude-code-logs_#{version}_darwin_amd64.tar.gz"
  sha256 "PLACEHOLDER"
  license "MIT"

  depends_on "go" => :build

  def install
    bin.install "claude-code-logs"
  end

  test do
    system "#{bin}/claude-code-logs", "version"
  end
end
```

### Installation

```bash
brew tap fabriqa/tap
brew install claude-code-logs
```

## Key Operations

1. **Create tap repository** - `fabriqaai/homebrew-tap`
2. **Create formula** - `Formula/claude-code-logs.rb`
3. **Setup goreleaser** - For automated releases
4. **Update formula on release** - GitHub Actions
5. **Auto-version on commit** - GitHub Actions bumps patch version on every main branch commit

## Dependencies

- **005-cli**: Produces binary to distribute

## Interface

- Users install via: `brew install fabriqa/tap/claude-code-logs`
- GitHub Actions updates formula on new release

## Versioning Strategy

**Auto-increment patch version on every commit to main:**

```yaml
# .github/workflows/version-bump.yml
name: Version Bump

on:
  push:
    branches: [main]

jobs:
  bump:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - name: Bump version and push tag
        uses: anothrNick/github-tag-action@v1
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          DEFAULT_BUMP: patch
          WITH_V: true
```

**Version format:** `v0.1.0`, `v0.1.1`, `v0.1.2`, etc.

- Every commit to `main` → patch bump (0.0.X)
- Manual tag for minor/major when needed

## Technical Constraints

- Separate repository (`fabriqaai/homebrew-tap`)
- Formula must follow Homebrew conventions
- Support both macOS (arm64, amd64) and Linux (amd64)
- Use goreleaser for cross-platform builds
- Auto-versioning via GitHub Actions on main branch commits

## Success Criteria

- [ ] Tap repository created at `fabriqaai/homebrew-tap`
- [ ] Formula installs working binary
- [ ] `brew install fabriqa/tap/claude-code-logs` works
- [ ] Formula auto-updates on new releases
- [ ] Version auto-increments (patch) on every main commit
- [ ] Supports macOS arm64 (Apple Silicon)
- [ ] Supports macOS amd64 (Intel)
- [ ] Supports Linux amd64

---

## Story Summary

- **Total Stories**: 1
- **Must Have**: 0
- **Should Have**: 1
- **Could Have**: 0

### Stories

- [ ] **001-homebrew-distribution**: Homebrew distribution - Should - Planned
