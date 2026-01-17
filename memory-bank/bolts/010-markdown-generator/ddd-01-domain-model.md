---
unit: 001-markdown-generator
bolt: 010-markdown-generator
stage: model
status: complete
updated: 2026-01-17T17:10:00Z
---

# Static Model - Markdown Generator

## Bounded Context

The Markdown Generator is responsible for transforming Claude Code session logs (JSONL format) into portable Markdown files with YAML frontmatter. This context handles file generation, change detection, and orphan preservation—distinct from the existing Parser context (which reads JSONL) and the future Client Renderer context (which displays MD).

## Domain Entities

| Entity | Properties | Business Rules |
|--------|------------|----------------|
| Session | ID, Title, Messages[], CreatedAt, UpdatedAt, CWD, SourcePath | Already exists in Parser context; extended with SourceHash |
| Project | Path, FolderName, Sessions[] | Already exists; groups sessions by working directory |
| MarkdownFile | FilePath, Content, Frontmatter, ModifiedTime | Represents generated .md output; content derived from Session |
| IndexFile | FilePath, Entries[], Type (main\|project) | Represents index.md listing projects or sessions |

## Value Objects

| Value Object | Properties | Constraints |
|--------------|------------|-------------|
| Frontmatter | Source, SourceHash, Project, Title, Created | Immutable; computed from Session; YAML-serializable |
| SourceHash | Value (string) | SHA256 hex digest of JSONL file content; 64 chars |
| FileModTime | Time (timestamp) | From os.Stat(); used for skip logic |
| RelativeLink | Path (string) | Relative path for navigation; no absolute paths |

## Aggregates

| Aggregate Root | Members | Invariants |
|----------------|---------|------------|
| GenerationContext | Sessions[], Projects[], OutputDir, SourceDir, ForceFlag | Sessions must have unique IDs; OutputDir must exist |
| OrphanSet | OrphanedMDFiles[], ExpectedMDFiles[] | Orphans = Actual - Expected; orphans are never deleted |

## Domain Events

| Event | Trigger | Payload |
|-------|---------|---------|
| SessionMarkdownGenerated | Session converted to MD | SessionID, OutputPath, SourceHash |
| SessionMarkdownSkipped | mtime check passed (no change) | SessionID, Reason |
| IndexMarkdownGenerated | Index file created/updated | IndexType, OutputPath, EntryCount |
| OrphanDetected | MD exists without JSONL | OrphanPath, FrontmatterMetadata |

## Domain Services

| Service | Operations | Dependencies |
|---------|------------|--------------|
| SessionMarkdownGenerator | GenerateSessionMD(session) → MarkdownFile | Frontmatter builder, Message formatter |
| IndexMarkdownGenerator | GenerateMainIndex(projects), GenerateProjectIndex(project) | RelativeLink builder |
| ModificationTimeChecker | ShouldRegenerate(jsonlPath, mdPath, force) → bool | FileSystem for mtime |
| SourceHasher | ComputeHash(jsonlPath) → SourceHash | crypto/sha256 |
| OrphanDetector | FindOrphans(outputDir, expectedFiles) → OrphanSet | FileSystem glob, Frontmatter parser |
| FrontmatterParser | ParseFrontmatter(mdContent) → Frontmatter | YAML parser |
| MessageFormatter | FormatMessage(message) → MarkdownString | Code fence handling, tool call formatting |

## Repository Interfaces

| Repository | Entity | Methods |
|------------|--------|---------|
| MarkdownFileWriter | MarkdownFile | Write(file), Exists(path), GetMTime(path) |
| MarkdownFileReader | MarkdownFile | Read(path), ListAll(dir), GetFrontmatter(path) |

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Session | A single Claude Code conversation with a unique ID |
| Project | A directory that was used as CWD for one or more sessions |
| Frontmatter | YAML metadata block at the top of a Markdown file, delimited by `---` |
| Source Hash | SHA256 digest of the original JSONL file, stored in frontmatter for verification |
| mtime | File modification timestamp from the filesystem |
| Skip Logic | Decision to not regenerate an MD file because JSONL hasn't changed |
| Orphan | An MD file whose source JSONL no longer exists (Claude deleted it after 30 days) |
| Force Flag | CLI option `--force` that bypasses mtime check and regenerates all files |
| Index File | A Markdown file listing projects (main index) or sessions (project index) |

## Story Coverage Matrix

| Story | Entities | Services | Notes |
|-------|----------|----------|-------|
| 001-session-markdown | Session, MarkdownFile, Frontmatter | SessionMarkdownGenerator, SourceHasher, MessageFormatter | Core generation |
| 002-index-markdown | Project, IndexFile, RelativeLink | IndexMarkdownGenerator | Navigation files |
| 003-mtime-skip | FileModTime | ModificationTimeChecker | Performance optimization |
| 004-orphan-preserve | OrphanSet, Frontmatter | OrphanDetector, FrontmatterParser | Archive preservation |

## Message Formatting Rules

Messages are formatted as Markdown with these rules:

1. **User messages**: `## User` header followed by message content
2. **Assistant messages**: `## Assistant` header followed by message content
3. **Code blocks**: Preserve language fencing (```language)
4. **Tool calls**: Format as collapsible details or code blocks
5. **System messages**: Omit or format as blockquote (TBD in Technical Design)

## Frontmatter Schema

```yaml
---
source: "session-abc123.jsonl"
source_hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
project: "/Users/user/code/my-project"
title: "Session summary from first message"
created: "2026-01-17T10:30:00Z"
---
```

## Index File Schemas

### Main Index (index.md)
```markdown
# Claude Code Logs

| Project | Sessions | Last Activity |
|---------|----------|---------------|
| [project-name](project-name/index.md) | 5 | 2026-01-17 |
```

### Project Index ({project}/index.md)
```markdown
# Project: /Users/user/code/my-project

| Session | Title | Created |
|---------|-------|---------|
| [abc123](abc123.md) | Fix authentication bug | 2026-01-17 |
```
