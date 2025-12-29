---
stage: plan
bolt: 006-homebrew-tap
created: 2025-12-29T16:35:00Z
---

## Implementation Plan: homebrew-tap

### Objective

Set up Homebrew distribution infrastructure for claude-code-logs CLI, including:
- goreleaser configuration for cross-platform binary builds
- GitHub Actions for auto-versioning on main branch commits
- GitHub Actions for creating releases on new tags
- Homebrew tap repository with formula

### Deliverables

**In this repository (claude-code-logs):**

1. **`.goreleaser.yaml`** - Cross-platform build configuration
   - Builds for darwin/amd64, darwin/arm64, linux/amd64
   - Creates tarballs with checksums
   - Updates Homebrew formula automatically

2. **`.github/workflows/version-bump.yml`** - Auto-versioning workflow
   - Triggers on push to main
   - Uses anothrNick/github-tag-action
   - Bumps patch version automatically (v0.1.0 → v0.1.1)

3. **`.github/workflows/release.yml`** - Release workflow
   - Triggers on new version tags (v*)
   - Runs goreleaser to build binaries
   - Publishes GitHub release with artifacts

**In homebrew-tap repository (fabriqaai/homebrew-tap):**

4. **`Formula/claude-code-logs.rb`** - Homebrew formula
   - Uses pre-built binaries from releases
   - Supports macOS (arm64, amd64) and Linux (amd64)
   - Auto-updated by goreleaser

5. **`README.md`** - Installation instructions

### Dependencies

| Dependency | Purpose |
|------------|---------|
| `goreleaser/goreleaser-action` | Build and release Go binaries |
| `anothrNick/github-tag-action` | Auto-increment version tags |
| `fabriqaai/homebrew-tap` | Separate tap repository |
| `HOMEBREW_TAP_TOKEN` | GitHub PAT for cross-repo formula updates |

### Technical Approach

**1. goreleaser Configuration**

goreleaser will:
- Build Go binaries for 3 platforms (darwin/amd64, darwin/arm64, linux/amd64)
- Create tarballs with checksums
- Generate changelog from commits
- Update homebrew-tap formula with new version and SHA256

**2. Version Bumping Strategy**

Using `anothrNick/github-tag-action`:
- Every push to `main` triggers patch bump
- Tag format: `v0.1.X`
- Commit message can override: `#major`, `#minor`, `#patch`, `#none`

**3. Release Workflow**

On new tag:
1. Checkout code
2. Set up Go
3. Run goreleaser
4. goreleaser pushes formula update to homebrew-tap

**4. Homebrew Formula**

Formula will:
- Download pre-built binary (not build from source)
- Support `on_macos` and `on_linux` blocks for platform detection
- Use `hardware.cpu.arm?` for Apple Silicon detection

### Acceptance Criteria

- [ ] goreleaser.yaml configured correctly
- [ ] version-bump.yml triggers on main push
- [ ] release.yml triggers on v* tags
- [ ] goreleaser builds all 3 platforms
- [ ] goreleaser updates homebrew-tap formula
- [ ] Formula works: `brew tap fabriqa/tap`
- [ ] Formula works: `brew install fabriqa/tap/claude-code-logs`
- [ ] Binary executes: `claude-code-logs version`

### Prerequisites (Manual Steps Required)

**IMPORTANT**: These steps require human action outside this codebase:

1. **Create fabriqaai/homebrew-tap repository on GitHub**
   - Must be named exactly `homebrew-tap` for Homebrew to recognize it
   - Initialize with README

2. **Create GitHub Personal Access Token (PAT)**
   - Token needs `repo` scope to push to homebrew-tap
   - Add as secret `HOMEBREW_TAP_TOKEN` in claude-code-logs repo

3. **Push initial formula to homebrew-tap**
   - Formula/claude-code-logs.rb will be created by first goreleaser run

### Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| PAT expires | Document renewal process |
| goreleaser fails | Test locally with `goreleaser build --snapshot` |
| Formula syntax error | Validate with `brew audit` |
| Cross-repo push fails | Check PAT permissions |

### Estimated Effort

- goreleaser.yaml: 15-30 minutes
- GitHub workflows: 30-45 minutes
- Formula creation: 15 minutes
- Testing: 30-60 minutes
- **Total**: ~2 hours
