---
id: 006-homebrew-tap
unit: 006-homebrew-tap
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: completed
stories:
  - 001-homebrew-distribution
created: 2025-12-29T12:45:00Z
started: 2025-12-29T16:30:00Z
completed: 2025-12-29T17:20:00Z
current_stage: completed
stages_completed:
  - name: plan
    completed: 2025-12-29T16:40:00Z
    artifact: implementation-plan.md
  - name: implement
    completed: 2025-12-29T16:55:00Z
    artifact: implementation-walkthrough.md
  - name: test
    completed: 2025-12-29T17:20:00Z
    artifact: test-walkthrough.md

# Bolt Dependencies
requires_bolts:
  - 004-cli
enables_bolts: []
requires_units: []
blocks: false

# Complexity Assessment
complexity:
  avg_complexity: 2
  avg_uncertainty: 2
  max_dependencies: 1
  testing_scope: 1
---

# Bolt: 006-homebrew-tap

## Overview

Create and configure Homebrew tap for distributing claude-code-logs CLI, including GitHub Actions for automated versioning and formula updates.

## Objective

Set up fabriqaai/homebrew-tap repository with formula, goreleaser configuration, and GitHub Actions for auto-versioning.

## Stories Included

- **001-homebrew-distribution**: Create Homebrew tap and formula - Should

## Bolt Type

**Type**: Simple Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/simple-construction-bolt.md`

## Stages

- [x] **1. Plan**: Completed → `implementation-plan.md`
- [x] **2. Implement**: Completed → Source code + `implementation-walkthrough.md`
- [x] **3. Test**: Completed → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 004-cli (produces the binary to distribute)

### Enables
- None (final distribution bolt)

## Success Criteria

- [ ] fabriqaai/homebrew-tap repository created (manual prerequisite)
- [ ] Formula/claude-code-logs.rb works (pending first release)
- [ ] `brew tap fabriqa/tap` succeeds (pending first release)
- [ ] `brew install fabriqa/tap/claude-code-logs` installs correctly (pending first release)
- [x] goreleaser.yaml configured for cross-platform builds
- [x] GitHub Actions version-bump.yml auto-increments patch on main commits
- [x] GitHub Actions release.yml creates releases on new tags
- [x] Supports macOS arm64 (Apple Silicon)
- [x] Supports macOS amd64 (Intel)
- [x] Supports Linux amd64
- [x] Formula auto-updates on new releases (configured in goreleaser)

## Notes

- Separate repository: fabriqaai/homebrew-tap
- Use anothrNick/github-tag-action for version bumping
- goreleaser handles cross-compilation and release artifacts
- Version format: v0.1.0, v0.1.1, etc.
