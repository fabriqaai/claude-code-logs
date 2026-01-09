---
id: 001-homebrew-distribution
unit: 006-homebrew-tap
intent: 001-chat-log-viewer
status: complete
priority: should
created: 2025-12-29T12:00:00Z
assigned_bolt: 006-homebrew-tap
implemented: true
---

# Story: 001-homebrew-distribution

## User Story

**As a** developer
**I want** to install via Homebrew
**So that** I can easily install and update the tool

## Acceptance Criteria

- [ ] **Given** tap repository exists, **When** `brew tap fabriqa/tap` run, **Then** tap is added successfully
- [ ] **Given** tap added, **When** `brew install fabriqa/tap/claude-code-logs` run, **Then** CLI is installed
- [ ] **Given** any commit to main branch, **When** pushed, **Then** version is auto-bumped (patch)
- [ ] **Given** new version tag created, **When** release workflow runs, **Then** formula is auto-updated
- [ ] **Given** macOS arm64 (Apple Silicon), **When** installing, **Then** correct binary downloaded
- [ ] **Given** macOS amd64 (Intel), **When** installing, **Then** correct binary downloaded
- [ ] **Given** Linux amd64, **When** installing, **Then** correct binary downloaded

## Technical Notes

- Create separate repository: `fabriqaai/homebrew-tap`
- Use goreleaser for cross-platform builds
- GitHub Actions workflow for:
  - Version bump on main commits (anothrNick/github-tag-action)
  - Release creation on new tags
  - Formula update on release
- Formula path: `Formula/claude-code-logs.rb`
- Version format: `v0.1.0`, `v0.1.1`, etc.

## Dependencies

### Requires
- 005-cli (produces the binary to distribute)

### Enables
- None (final distribution unit)

## Edge Cases

| Scenario | Expected Behavior |
|----------|-------------------|
| goreleaser build fails | Release aborted, no broken formula |
| SHA256 mismatch | Homebrew install fails with clear error |
| Previous version installed | Upgrade works correctly |
| Network error during download | Homebrew retries |
| Old formula version cached | `brew update` before install |
| M1 Mac running Rosetta | Serve arm64 binary preferentially |

## Out of Scope

- Linux package managers (apt, yum)
- Windows distribution (chocolatey, scoop)
- Docker image distribution
- Manual binary download page
