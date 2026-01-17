---
intent: 004-markdown-first-generation
phase: inception
status: context-defined
updated: 2026-01-17T16:30:00Z
---

# Markdown-First Generation - System Context

## System Overview

Transform claude-code-logs from HTML-first to Markdown-first generation. The Go backend generates Markdown files with YAML frontmatter as the source of truth. An HTML shell with embedded JavaScript loads these MD files and renders them client-side using marked.js + highlight.js, preserving the existing UI/UX.

## Context Diagram

```
┌─────────────────────────────────────────────────────────────────────┐
│                        claude-code-logs                              │
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────────────┐│
│  │   Parser     │────▶│  Markdown    │────▶│   Output Directory   ││
│  │ (existing)   │     │  Generator   │     │  ├── index.md        ││
│  └──────────────┘     │  (modified)  │     │  ├── project/        ││
│         │             └──────────────┘     │  │   ├── index.md    ││
│         │                    │             │  │   └── session.md  ││
│         ▼                    │             │  └── shell.html      ││
│  ┌──────────────┐            │             └──────────────────────┘│
│  │ ~/.claude/   │            │                        │            │
│  │ projects/    │            ▼                        ▼            │
│  │ (JSONL)      │     ┌──────────────┐     ┌──────────────────────┐│
│  └──────────────┘     │   mtime      │     │   HTTP Server        ││
│                       │   Checker    │     │   (existing)         ││
│                       └──────────────┘     └──────────────────────┘│
│                                                       │            │
└───────────────────────────────────────────────────────│────────────┘
                                                        │
                                                        ▼
┌─────────────────────────────────────────────────────────────────────┐
│                           Browser                                    │
│                                                                      │
│  ┌──────────────┐     ┌──────────────┐     ┌──────────────────────┐│
│  │  HTML Shell  │────▶│ fetch(.md)   │────▶│   marked.js +        ││
│  │  (sidebar,   │     │              │     │   highlight.js       ││
│  │   CSS, JS)   │     └──────────────┘     └──────────────────────┘│
│  └──────────────┘                                    │             │
│         │                                            ▼             │
│         │                                  ┌──────────────────────┐│
│         └─────────────────────────────────▶│   Rendered UI        ││
│                                            │   (same as today)    ││
│                                            └──────────────────────┘│
└─────────────────────────────────────────────────────────────────────┘
```

## External Integrations

- **marked.js**: Markdown parsing library (CDN or bundled)
- **highlight.js**: Code syntax highlighting (CDN or bundled)
- **Filesystem**: Read JSONL from ~/.claude/projects, write MD to output dir

## High-Level Constraints

- Must preserve existing UI/UX exactly
- Must work offline (no external API calls for rendering)
- MD files must be valid CommonMark + GFM
- Frontmatter must be valid YAML

## Key NFR Goals

- Client-side render time < 100ms
- mtime check < 1ms per file
- Full regeneration < 30s for 100 sessions
- Zero git churn for unchanged files
