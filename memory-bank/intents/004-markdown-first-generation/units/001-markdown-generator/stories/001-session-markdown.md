---
story: 001-session-markdown
unit: 001-markdown-generator
intent: 004-markdown-first-generation
priority: Must
status: planned
created: 2026-01-17T16:30:00Z
---

# Story: Generate Session Markdown with Frontmatter

## User Story

As a user, I want each session to be saved as a Markdown file with YAML frontmatter containing metadata, so that I can archive logs in git and feed them to LLM tools.

## Acceptance Criteria

- [ ] Each session generates a `{session-id}.md` file
- [ ] Frontmatter contains:
  - `source`: original JSONL filename
  - `source_hash`: SHA256 of JSONL content
  - `project`: actual project path (from cwd)
  - `title`: session summary
  - `created`: session creation timestamp
- [ ] Messages are formatted as Markdown:
  - User messages with `## User` header
  - Assistant messages with `## Assistant` header
  - Code blocks with proper language fencing
  - Tool calls in collapsible sections or code blocks
- [ ] Output is valid CommonMark + GFM

## Technical Notes

- Replace `templates_session.go` HTML generation with MD generation
- Use text/template or string building for MD output
- Compute SHA256 hash of JSONL file for frontmatter
- Format timestamps consistently (ISO 8601)

## Dependencies

- Parser (existing) provides Session struct
- crypto/sha256 for hashing

## Estimation

Complexity: M
