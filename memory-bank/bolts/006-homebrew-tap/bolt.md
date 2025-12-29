---
id: 006-homebrew-tap
unit: 006-homebrew-tap
intent: 001-chat-log-viewer
type: simple-construction-bolt
status: planned
stories:
  - 001-homebrew-distribution
created: 2025-12-29T12:45:00Z
started: null
completed: null
current_stage: null
stages_completed: []

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

- [ ] **1. Plan**: Pending → `implementation-plan.md`
- [ ] **2. Implement**: Pending → Source code + `implementation-walkthrough.md`
- [ ] **3. Test**: Pending → Tests + `test-walkthrough.md`

## Dependencies

### Requires
- 004-cli (produces the binary to distribute)

### Enables
- None (final distribution bolt)

## Success Criteria

- [ ] fabriqaai/homebrew-tap repository created
- [ ] Formula/claude-code-logs.rb works
- [ ] `brew tap fabriqa/tap` succeeds
- [ ] `brew install fabriqa/tap/claude-code-logs` installs correctly
- [ ] goreleaser.yaml configured for cross-platform builds
- [ ] GitHub Actions version-bump.yml auto-increments patch on main commits
- [ ] GitHub Actions release.yml creates releases on new tags
- [ ] Supports macOS arm64 (Apple Silicon)
- [ ] Supports macOS amd64 (Intel)
- [ ] Supports Linux amd64
- [ ] Formula auto-updates on new releases

## Notes

- Separate repository: fabriqaai/homebrew-tap
- Use anothrNick/github-tag-action for version bumping
- goreleaser handles cross-compilation and release artifacts
- Version format: v0.1.0, v0.1.1, etc.
