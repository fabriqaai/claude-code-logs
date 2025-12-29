---
stage: implement
bolt: 006-homebrew-tap
created: 2025-12-29T16:50:00Z
---

## Implementation Walkthrough: homebrew-tap

### Summary

Created goreleaser configuration and GitHub Actions workflows for automated versioning and cross-platform binary releases. Added homebrew-tap template files for the external tap repository that goreleaser will update automatically.

### Structure Overview

The implementation adds release automation to the existing CLI project. goreleaser builds binaries for three platforms (darwin/amd64, darwin/arm64, linux/amd64), creates GitHub releases, and pushes formula updates to the external homebrew-tap repository.

### Completed Work

- [x] `.goreleaser.yaml` - Cross-platform build configuration with Homebrew formula generation
- [x] `.github/workflows/version-bump.yml` - Auto-increment patch version on main branch commits
- [x] `.github/workflows/release.yml` - Build and release binaries on new version tags
- [x] `main.go` - Added commit and date build variables for version injection
- [x] `cmd_version.go` - Updated to display commit hash and build date
- [x] `homebrew-tap-template/README.md` - Installation documentation for tap repository
- [x] `homebrew-tap-template/Formula/claude-code-logs.rb.template` - Reference formula structure

### Key Decisions

- **Pre-built binaries over source builds**: Formula downloads pre-built binaries rather than compiling from source, reducing install time and eliminating Go dependency for end users
- **anothrNick/github-tag-action**: Chosen for its simplicity and support for commit message overrides (#major, #minor, #none)
- **Separate homebrew-tap repository**: Following Homebrew best practices for third-party formulae
- **CGO_ENABLED=0**: Ensures fully static binaries without C dependencies

### Deviations from Plan

None

### Dependencies Added

- [x] `goreleaser/goreleaser-action@v6` - GitHub Action for running goreleaser
- [x] `anothrNick/github-tag-action@v1` - GitHub Action for version bumping

### Developer Notes

**Prerequisites for release workflow:**
1. Create `fabriqaai/homebrew-tap` repository on GitHub
2. Create GitHub PAT with `repo` scope
3. Add `HOMEBREW_TAP_TOKEN` secret to this repository's settings

**Local testing:**
```bash
# Verify goreleaser config
goreleaser check

# Test build locally (no release)
goreleaser build --snapshot --clean
```

**Version control:**
- Commit messages with `#major` bump major version
- Commit messages with `#minor` bump minor version
- Commit messages with `#none` skip version bump
- Default behavior bumps patch version
