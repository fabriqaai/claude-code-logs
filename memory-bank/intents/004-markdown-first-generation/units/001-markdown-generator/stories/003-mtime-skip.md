---
story: 003-mtime-skip
unit: 001-markdown-generator
intent: 004-markdown-first-generation
priority: Must
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Implement mtime-Based Skip Logic

## User Story

As a user, I want unchanged sessions to be skipped during regeneration, so that my git repo doesn't have unnecessary churn and regeneration is fast.

## Acceptance Criteria

- [ ] Before generating MD, check if MD file exists
- [ ] Compare JSONL mtime with MD mtime
- [ ] If MD mtime >= JSONL mtime, skip regeneration
- [ ] If JSONL is newer, regenerate MD
- [ ] `--force` flag bypasses mtime check and regenerates all
- [ ] Log skipped files in verbose mode

## Technical Notes

- Use os.Stat() to get file modification times
- Add `--force` flag to serve command
- Track skip count for summary output
- Handle edge case: MD exists but JSONL doesn't (orphan - see story 004)

## Dependencies

- None (filesystem operations only)

## Estimation

Complexity: S
