---
intent: 004-markdown-first-generation
phase: inception
status: units-decomposed
updated: 2026-01-17T16:30:00Z
---

# Markdown-First Generation - Unit Decomposition

## Units Overview

This intent decomposes into 2 units of work:

### Unit 1: 001-markdown-generator

**Description**: Modify Go backend to generate Markdown files instead of HTML, with YAML frontmatter, mtime-based skip logic, and orphan preservation.

**Assigned FRs**:
- FR-1: Markdown Generation with Frontmatter
- FR-2: mtime-Based Regeneration Skip
- FR-3: Preserve Orphaned Markdown Files
- FR-7: Index Files Generation

**Deliverables**:
- Modified generator.go (or new md_generator.go)
- Markdown templates for sessions and indexes
- mtime comparison logic
- Orphan detection and preservation

**Dependencies**:
- Depends on: None (foundational)
- Depended by: 002-client-renderer

**Estimated Complexity**: M

---

### Unit 2: 002-client-renderer

**Description**: Create HTML shell with JavaScript that loads MD files via fetch and renders them using marked.js + highlight.js, plus download/copy buttons.

**Assigned FRs**:
- FR-4: Client-Side Markdown Rendering
- FR-5: Download as Markdown Button
- FR-6: Copy as Markdown Button

**Deliverables**:
- HTML shell template (replaces session template)
- JavaScript for MD fetching and rendering
- marked.js + highlight.js integration
- Download/Copy button implementation

**Dependencies**:
- Depends on: 001-markdown-generator (needs MD files to render)
- Depended by: None

**Estimated Complexity**: M

---

## Unit Dependency Graph

```text
[001-markdown-generator] ──────▶ [002-client-renderer]
     (generates MD)                  (renders MD)
```

## Execution Order

Based on dependencies:

1. **Unit 1**: 001-markdown-generator (foundation - creates MD files)
2. **Unit 2**: 002-client-renderer (consumes MD files)

Sequential execution required - Unit 2 needs MD files to test against.
