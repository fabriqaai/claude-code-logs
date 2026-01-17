---
story: 004-orphan-preserve
unit: 001-markdown-generator
intent: 004-markdown-first-generation
priority: Must
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Preserve Orphaned Markdown Files

## User Story

As a user, I want MD files to be preserved when the source JSONL is deleted (e.g., Claude's 30-day cleanup), so that my archive remains complete.

## Acceptance Criteria

- [ ] Scan output directory for existing MD files
- [ ] For each MD file, check if corresponding JSONL exists
- [ ] If JSONL missing, keep MD file (do NOT delete)
- [ ] Include orphaned sessions in index listings
- [ ] Orphaned sessions are navigable in UI

## Technical Notes

- During generation, build list of "expected" MD files from JSONL
- Compare with actual MD files in output directory
- Orphans = actual - expected (keep these)
- Read frontmatter from orphaned MD to get metadata for index
- May need to parse frontmatter YAML to extract title, project, date

## Dependencies

- Story 001 (frontmatter format must support reading back)
- YAML parsing for frontmatter

## Estimation

Complexity: M
