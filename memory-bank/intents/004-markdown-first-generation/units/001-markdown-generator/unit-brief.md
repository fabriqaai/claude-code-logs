---
unit: 001-markdown-generator
intent: 004-markdown-first-generation
phase: inception
status: ready
created: 2026-01-17T16:30:00Z
updated: 2026-01-17T16:30:00Z
---

# Unit Brief: Markdown Generator

## Purpose

Modify the Go backend to generate Markdown files with YAML frontmatter instead of HTML. Implement mtime-based regeneration skip to avoid git churn, and preserve orphaned MD files when source JSONL is deleted.

## Scope

### In Scope
- Generate session Markdown files with frontmatter (source, hash, project, title, date)
- Generate index Markdown files (main index, per-project index)
- mtime comparison to skip unchanged files
- Preserve MD files when JSONL source is deleted
- `--force` flag to bypass mtime check
- Remove HTML generation (MD-only output)

### Out of Scope
- Client-side rendering (Unit 2)
- Download/Copy buttons (Unit 2)
- UI changes (Unit 2)

---

## Assigned Requirements

| FR | Requirement | Priority |
|----|-------------|----------|
| FR-1 | Markdown Generation with Frontmatter | Must |
| FR-2 | mtime-Based Regeneration Skip | Must |
| FR-3 | Preserve Orphaned Markdown Files | Must |
| FR-7 | Index Files Generation | Must |

---

## Domain Concepts

### Key Entities

| Entity | Description | Attributes |
|--------|-------------|------------|
| SessionMarkdown | MD file for a session | frontmatter, messages, code blocks |
| IndexMarkdown | MD file listing projects/sessions | frontmatter, project list |
| Frontmatter | YAML metadata header | source, source_hash, project, title, created |

### Key Operations

| Operation | Description | Inputs | Outputs |
|-----------|-------------|--------|---------|
| GenerateSessionMD | Convert session to Markdown | Session struct | .md file |
| GenerateIndexMD | Create index listing | []Project | index.md |
| CheckMtime | Compare JSONL vs MD mtime | file paths | bool (skip?) |
| ComputeHash | SHA256 of JSONL content | file path | hash string |
| DetectOrphans | Find MD without JSONL | directories | []orphaned paths |

---

## Story Summary

| Metric | Count |
|--------|-------|
| Total Stories | 4 |
| Must Have | 4 |
| Should Have | 0 |
| Could Have | 0 |

### Stories

| Story ID | Title | Priority | Status |
|----------|-------|----------|--------|
| 001-session-markdown | Generate Session Markdown with Frontmatter | Must | Planned |
| 002-index-markdown | Generate Index Markdown Files | Must | Planned |
| 003-mtime-skip | Implement mtime-Based Skip Logic | Must | Planned |
| 004-orphan-preserve | Preserve Orphaned Markdown Files | Must | Planned |

---

## Dependencies

### Depends On
| Unit | Reason |
|------|--------|
| None | Foundational unit |

### Depended By
| Unit | Reason |
|------|--------|
| 002-client-renderer | Needs MD files to render |

### External Dependencies
| System | Purpose | Risk |
|--------|---------|------|
| Filesystem | Read JSONL, write MD | Low |

---

## Technical Context

### Suggested Technology
- Go (existing codebase)
- text/template or strings for MD generation
- crypto/sha256 for hashing
- os.Stat for mtime comparison

### Integration Points
| Integration | Type | Protocol |
|-------------|------|----------|
| Parser | Internal | Function call |
| Filesystem | OS | File I/O |

### Data Storage
| Data | Type | Volume | Retention |
|------|------|--------|-----------|
| Markdown files | File | 1 per session | Permanent (archive) |

---

## Constraints

- Frontmatter must be valid YAML
- MD must be valid CommonMark + GFM
- Hash algorithm: SHA256
- Deterministic output (sorted, no timestamps in content)

---

## Success Criteria

### Functional
- [ ] Session MD files generated with correct frontmatter
- [ ] Index MD files generated with project/session lists
- [ ] Unchanged files are skipped (mtime check works)
- [ ] Orphaned MD files are preserved

### Non-Functional
- [ ] Full regeneration < 30s for 100 sessions
- [ ] mtime check < 1ms per file

### Quality
- [ ] All existing tests pass
- [ ] New tests for MD generation
- [ ] Code reviewed

---

## Bolt Suggestions

| Bolt | Type | Stories | Objective |
|------|------|---------|-----------|
| bolt-001-md-gen | DDD | S1, S2, S3, S4 | Complete MD generator |

---

## Notes

This unit replaces the existing HTML generation entirely. The Go templates in templates_*.go will be replaced with MD generation logic.
