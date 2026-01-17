---
intent: 004-markdown-first-generation
phase: inception
status: draft
created: 2026-01-17T16:00:00Z
updated: 2026-01-17T16:00:00Z
---

# Requirements: Markdown-First Generation

## Intent Overview

Transform the output format from HTML-only to Markdown-first generation. Markdown files with YAML frontmatter become the source of truth, rendered to HTML client-side using marked.js + highlight.js. This enables LLM-friendly logs, git-based archival, and preservation of session history beyond Claude's 30-day retention.

## Business Goals

| Goal | Success Metric | Priority |
|------|----------------|----------|
| LLM-friendly output | Markdown files can be fed to AI tools | Must |
| Git-friendly archival | No unnecessary regeneration/churn | Must |
| Preserve history | Orphaned MD files kept when JSONL deleted | Must |
| Better code highlighting | Syntax highlighting for all languages | Should |
| Download/copy capability | Users can export sessions as Markdown | Should |

---

## Functional Requirements

### FR-1: Markdown Generation with Frontmatter
- **Description**: Generate `.md` files instead of (or alongside) HTML for each session
- **Acceptance Criteria**:
  - Each session produces a `{session-id}.md` file
  - Frontmatter contains: source filename, source_hash (SHA256), project path, title, created date
  - Content is valid Markdown with proper code block fencing
- **Priority**: Must
- **Related Stories**: TBD

### FR-2: mtime-Based Regeneration Skip
- **Description**: Skip regeneration if source JSONL hasn't changed since last MD generation
- **Acceptance Criteria**:
  - Compare JSONL mtime with MD mtime
  - If MD is newer or equal, skip regeneration
  - If JSONL is newer, regenerate MD
  - `--force` flag bypasses mtime check
- **Priority**: Must
- **Related Stories**: TBD

### FR-3: Preserve Orphaned Markdown Files
- **Description**: When source JSONL is deleted (Claude cleanup), keep the MD file as archive
- **Acceptance Criteria**:
  - MD files without corresponding JSONL are NOT deleted
  - Orphaned MDs still appear in index/navigation
  - Clear indicator that source is archived (optional)
- **Priority**: Must
- **Related Stories**: TBD

### FR-4: Client-Side Markdown Rendering
- **Description**: Render Markdown to HTML in browser using marked.js + highlight.js
- **Acceptance Criteria**:
  - Session pages load MD file via fetch and render client-side
  - Code blocks have syntax highlighting with auto-language detection
  - GFM features work (tables, task lists, strikethrough)
  - Render time < 100ms for typical session
- **Priority**: Must
- **Related Stories**: TBD

### FR-5: Download as Markdown Button
- **Description**: Add button to download session as raw Markdown file
- **Acceptance Criteria**:
  - Download button visible on session page
  - Downloads the raw `.md` file with frontmatter
  - Filename matches session ID or title
- **Priority**: Should
- **Related Stories**: TBD

### FR-6: Copy as Markdown Button
- **Description**: Add button to copy session content as Markdown to clipboard
- **Acceptance Criteria**:
  - Copy button visible on session page
  - Copies full Markdown content INCLUDING frontmatter
  - Visual feedback on successful copy
- **Priority**: Should
- **Related Stories**: TBD

### FR-7: Index Files Generation
- **Description**: Generate index.md files for main page and per-project pages
- **Acceptance Criteria**:
  - `index.md` at root lists all projects
  - `{project}/index.md` lists all sessions for that project
  - Index files are regenerated when session list changes
  - Deterministic output (sorted, no random elements)
- **Priority**: Must
- **Related Stories**: TBD

---

## Non-Functional Requirements

### Performance
| Requirement | Metric | Target |
|-------------|--------|--------|
| MD Render Time | Client-side render | < 100ms |
| Regeneration Check | mtime comparison | < 1ms per file |
| Full Regeneration | All sessions | < 30s for 100 sessions |

### Compatibility
| Requirement | Standard | Notes |
|-------------|----------|-------|
| Markdown | CommonMark + GFM | Tables, task lists, code fences |
| Frontmatter | YAML | Standard Jekyll-style frontmatter |
| Browsers | ES6+ | Modern browsers only (local tool) |

### Maintainability
| Requirement | Metric | Target |
|-------------|--------|--------|
| Library Dependencies | External JS libs | marked.js + highlight.js only |
| Bundle Size | Total JS | < 200KB |

---

## Constraints

### Technical Constraints

**Intent-specific constraints**:
- Must use marked.js for Markdown rendering (decision made based on performance analysis)
- Must use highlight.js for syntax highlighting (best integration with marked.js)
- Frontmatter must be valid YAML parseable by common tools
- Hash algorithm: SHA256 for source_hash field

### Business Constraints
- Backward compatible: existing HTML output should still work during transition
- No breaking changes to CLI interface

---

## Assumptions

| Assumption | Risk if Invalid | Mitigation |
|------------|-----------------|------------|
| marked.js handles all GFM features needed | Missing feature | Fall back to markdown-it |
| mtime is reliable on user's filesystem | Unnecessary regeneration | Add `--use-hash` flag as fallback |
| Users have modern browsers | Rendering fails | Document browser requirements |

---

## Open Questions

| Question | Owner | Due Date | Resolution |
|----------|-------|----------|------------|
| Include frontmatter in clipboard copy? | User | 2026-01-17 | **Yes** - include frontmatter |
| Indicator for archived sessions? | User | 2026-01-17 | **No** - no special indicator needed |
| Transition strategy: MD-only or dual output? | User | 2026-01-17 | **MD-only** - replace HTML generation entirely |
