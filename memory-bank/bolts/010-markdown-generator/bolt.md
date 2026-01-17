---
id: 010-markdown-generator
unit: 001-markdown-generator
intent: 004-markdown-first-generation
type: ddd-construction-bolt
status: completed
stories:
  - 001-session-markdown
  - 002-index-markdown
  - 003-mtime-skip
  - 004-orphan-preserve
created: 2026-01-17T17:00:00Z
started: 2026-01-17T17:05:00Z
completed: 2026-01-17T18:00:00Z
current_stage: complete
stages_completed:
  - name: model
    completed: 2026-01-17T17:15:00Z
    artifact: ddd-01-domain-model.md
  - name: design
    completed: 2026-01-17T17:25:00Z
    artifact: ddd-02-technical-design.md
  - name: adr
    completed: 2026-01-17T17:30:00Z
    artifact: null
  - name: implement
    completed: 2026-01-17T17:45:00Z
    artifact: markdown_generator.go, frontmatter.go
  - name: test
    completed: 2026-01-17T18:00:00Z
    artifact: ddd-03-test-report.md

requires_bolts: []
enables_bolts:
  - 011-client-renderer
requires_units: []
blocks: false

complexity:
  avg_complexity: 2
  avg_uncertainty: 1
  max_dependencies: 1
  testing_scope: 2
---

# Bolt: 010-markdown-generator

## Overview

Transform the HTML generation system to produce Markdown files with YAML frontmatter, enabling git archival, LLM-friendly output, and preservation of sessions beyond Claude's 30-day cleanup.

## Objective

Generate Markdown files for all sessions and index files, with mtime-based skip logic and orphan preservation.

## Stories Included

- **001-session-markdown**: Generate Session Markdown with Frontmatter (Must)
- **002-index-markdown**: Generate Index Markdown Files (Must)
- **003-mtime-skip**: Implement mtime-Based Skip Logic (Must)
- **004-orphan-preserve**: Preserve Orphaned Markdown Files (Must)

## Bolt Type

**Type**: DDD Construction Bolt
**Definition**: `.specsmd/aidlc/templates/construction/bolt-types/ddd-construction-bolt.md`

## Stages

- ✅ **1. Domain Model**: Complete → ddd-01-domain-model.md
- ✅ **2. Technical Design**: Complete → ddd-02-technical-design.md
- ✅ **3. ADR Analysis**: Complete → (no ADRs needed)
- ✅ **4. Implement**: Complete → Source code
- ✅ **5. Test**: Complete → ddd-03-test-report.md

## Dependencies

### Requires
- None (first bolt in chain)

### Enables
- 011-client-renderer (needs MD files to render)

## Success Criteria

- [ ] All sessions generate valid MD files with frontmatter
- [ ] Frontmatter contains: source, source_hash, project, title, created
- [ ] Index files list all projects/sessions with relative links
- [ ] mtime check prevents unnecessary regeneration
- [ ] `--force` flag bypasses mtime check
- [ ] Orphaned MD files are preserved when JSONL is deleted
- [ ] All unit tests pass

## Notes

- Replaces templates_session.go, templates_index.go, templates_project.go
- Uses crypto/sha256 for source hash
- Output must be deterministic for git (no timestamps in content)
