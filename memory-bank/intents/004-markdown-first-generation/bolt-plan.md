---
intent: 004-markdown-first-generation
phase: inception
status: planned
created: 2026-01-17T16:30:00Z
updated: 2026-01-17T16:30:00Z
---

# Bolt Plan: Markdown-First Generation

## Overview

This plan defines the execution order for implementing markdown-first generation. The work is divided into two bolts that must be executed sequentially due to dependencies.

## Execution Order

```
bolt-001-generator (Unit 001)
         │
         ▼
bolt-002-renderer (Unit 002)
```

---

## Bolt 001: Markdown Generator

| Field | Value |
|-------|-------|
| **Bolt ID** | bolt-001-generator |
| **Unit** | 001-markdown-generator |
| **Type** | DDD |
| **Priority** | P1 - Critical Path |
| **Status** | Planned |

### Objective

Generate Markdown files with YAML frontmatter for all sessions and index files, replacing HTML generation.

### Stories Included

| Story | Title | Priority | Complexity |
|-------|-------|----------|------------|
| 001-session-markdown | Generate Session Markdown with Frontmatter | Must | M |
| 002-index-markdown | Generate Index Markdown Files | Must | S |
| 003-mtime-skip | Implement mtime-Based Skip Logic | Must | S |
| 004-orphan-preserve | Preserve Orphaned Markdown Files | Must | M |

### DDD Stages

#### Model Stage
- Define MD file structure and frontmatter schema
- Define index file structure
- Model orphan detection logic

#### Test Stage
- Unit tests for MD generation
- Unit tests for frontmatter formatting
- Unit tests for mtime comparison
- Unit tests for orphan detection

#### Implement Stage
- Implement session MD generator
- Implement index MD generator
- Implement mtime-based skip logic
- Implement orphan preservation

### Dependencies

| Type | Dependency | Reason |
|------|------------|--------|
| Internal | None | First bolt in chain |
| External | crypto/sha256 | For source hash in frontmatter |

### Success Criteria

- [ ] All sessions generate valid MD files
- [ ] Frontmatter contains all required fields
- [ ] Index files list all projects/sessions
- [ ] mtime check prevents unnecessary regeneration
- [ ] Orphaned MD files are preserved
- [ ] All unit tests pass

### Estimated Complexity

**Total: M** (4 stories: 2S + 2M)

---

## Bolt 002: Client Renderer

| Field | Value |
|-------|-------|
| **Bolt ID** | bolt-002-renderer |
| **Unit** | 002-client-renderer |
| **Type** | DDD |
| **Priority** | P1 - Critical Path |
| **Status** | Planned |

### Objective

Create HTML shell with JavaScript that fetches and renders Markdown files client-side, preserving existing UI/UX.

### Stories Included

| Story | Title | Priority | Complexity |
|-------|-------|----------|------------|
| 001-md-rendering | Client-Side Markdown Rendering | Must | M |
| 002-download-button | Download as Markdown Button | Should | S |
| 003-copy-button | Copy as Markdown Button | Should | S |

### DDD Stages

#### Model Stage
- Define HTML shell template structure
- Define JS rendering pipeline
- Define button placement and behavior

#### Test Stage
- Manual testing of rendering quality
- Cross-browser testing
- Performance testing (< 100ms render)

#### Implement Stage
- Create HTML shell template
- Implement JS fetch and render logic
- Integrate marked.js and highlight.js
- Implement download button
- Implement copy button

### Dependencies

| Type | Dependency | Reason |
|------|------------|--------|
| Internal | bolt-001-generator | Needs MD files to render |
| External | marked.js | MD to HTML parsing |
| External | highlight.js | Code syntax highlighting |

### Success Criteria

- [ ] MD files render with identical styling to current HTML
- [ ] Code blocks have syntax highlighting
- [ ] Download button works
- [ ] Copy button works with clipboard feedback
- [ ] Render time < 100ms
- [ ] Works offline (bundled or cached libs)

### Estimated Complexity

**Total: S-M** (3 stories: 2S + 1M)

---

## Summary

| Bolt | Stories | Must | Should | Complexity |
|------|---------|------|--------|------------|
| bolt-001-generator | 4 | 4 | 0 | M |
| bolt-002-renderer | 3 | 1 | 2 | S-M |
| **Total** | **7** | **5** | **2** | **M** |

## Risk Assessment

| Risk | Impact | Likelihood | Mitigation |
|------|--------|------------|------------|
| marked.js rendering differences | Medium | Low | Test early, configure options |
| Performance on large sessions | Medium | Low | Lazy load, pagination if needed |
| Browser compatibility | Low | Low | Use standard APIs with fallbacks |

## Notes

- Both bolts are on the critical path
- Unit 002 cannot start until Unit 001 produces MD files
- Consider running bolt-001 tests with sample MD files for bolt-002 development
