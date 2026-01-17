---
unit: 001-markdown-generator
bolt: 010-markdown-generator
stage: design
status: complete
updated: 2026-01-17T17:20:00Z
---

# Technical Design - Markdown Generator

## Architecture Pattern

**Pattern**: Service Layer within existing CLI architecture

The markdown generator integrates into the existing `claude-code-logs` Go CLI tool. It replaces the HTML template generation with Markdown file generation while preserving the existing Parser layer.

```text
┌─────────────────────────────────────────────────────────┐
│                     CLI Layer                            │
│  cmd_serve.go (flags: --force, existing --watch/--list)  │
├─────────────────────────────────────────────────────────┤
│                  Generator Layer                         │
│  markdown_generator.go (NEW - replaces HTML templates)   │
├─────────────────────────────────────────────────────────┤
│                   Parser Layer                           │
│  parser.go (existing - provides Session, Project)        │
├─────────────────────────────────────────────────────────┤
│                Infrastructure Layer                      │
│  File I/O, crypto/sha256, os.Stat                        │
└─────────────────────────────────────────────────────────┘
```

## File Structure

### New Files
```text
claude-code-logs/
├── markdown_generator.go      # Core MD generation logic
├── markdown_generator_test.go # Unit tests
├── frontmatter.go             # Frontmatter types and serialization
└── orphan_detector.go         # Orphan detection logic
```

### Modified Files
```text
├── cmd_serve.go               # Add --force flag, call MD generator
├── generator.go               # Refactor to use MD generator (or replace)
├── watcher.go                 # Update to regenerate MD instead of HTML
```

### Removed/Deprecated Files
```text
├── templates_session.go       # REMOVE - replaced by markdown_generator.go
├── templates_index.go         # REMOVE - replaced by markdown_generator.go
├── templates_project.go       # REMOVE - replaced by markdown_generator.go
```

## Data Flow

```text
                    ┌──────────────┐
                    │  JSONL Files │
                    │  (~/.claude) │
                    └──────┬───────┘
                           │
                    ┌──────▼───────┐
                    │    Parser    │
                    │  (existing)  │
                    └──────┬───────┘
                           │
              ┌────────────┼────────────┐
              │            │            │
       ┌──────▼──────┐ ┌───▼───┐ ┌─────▼─────┐
       │   Session   │ │Project│ │  Orphan   │
       │  []Message  │ │  []   │ │ Detector  │
       └──────┬──────┘ └───┬───┘ └─────┬─────┘
              │            │           │
       ┌──────▼──────┐ ┌───▼───┐ ┌─────▼─────┐
       │  Session MD │ │ Index │ │  Preserve │
       │  Generator  │ │  Gen  │ │  Orphans  │
       └──────┬──────┘ └───┬───┘ └─────┬─────┘
              │            │           │
              └────────────┼───────────┘
                           │
                    ┌──────▼───────┐
                    │   Output     │
                    │  Directory   │
                    │  (.md files) │
                    └──────────────┘
```

## Core Type Definitions

### Frontmatter (frontmatter.go)

```go
type Frontmatter struct {
    Source     string `yaml:"source"`      // Original JSONL filename
    SourceHash string `yaml:"source_hash"` // SHA256 of JSONL content
    Project    string `yaml:"project"`     // Actual project path (from CWD)
    Title      string `yaml:"title"`       // Session summary
    Created    string `yaml:"created"`     // ISO 8601 timestamp
}

func (f Frontmatter) Marshal() ([]byte, error)
func ParseFrontmatter(content []byte) (Frontmatter, []byte, error)
```

### MarkdownGenerator (markdown_generator.go)

```go
type MarkdownGenerator struct {
    OutputDir string
    SourceDir string
    Force     bool
}

type GenerationResult struct {
    Generated int
    Skipped   int
    Orphans   int
    Errors    []error
}

func NewMarkdownGenerator(outputDir, sourceDir string, force bool) *MarkdownGenerator
func (g *MarkdownGenerator) GenerateAll(projects []Project) (*GenerationResult, error)
func (g *MarkdownGenerator) GenerateSession(session *Session, projectFolder string) error
func (g *MarkdownGenerator) GenerateMainIndex(projects []Project) error
func (g *MarkdownGenerator) GenerateProjectIndex(project *Project, folder string) error
func (g *MarkdownGenerator) ShouldRegenerate(jsonlPath, mdPath string) bool
```

### OrphanDetector (orphan_detector.go)

```go
type OrphanInfo struct {
    Path       string
    Frontmatter Frontmatter
}

func DetectOrphans(outputDir string, expectedFiles map[string]bool) ([]OrphanInfo, error)
func PreserveOrphansInIndex(orphans []OrphanInfo, indexPath string) error
```

## API Design (CLI Flags)

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--force` | bool | false | Bypass mtime check, regenerate all files |
| `--output` | string | "./output" | Output directory for MD files |
| `--source` | string | "~/.claude/projects" | Source directory for JSONL files |

Existing flags preserved: `--watch`, `--list`, `--port`

## File Output Formats

### Session Markdown ({session-id}.md)

```markdown
---
source: "session-abc123.jsonl"
source_hash: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
project: "/Users/user/code/my-project"
title: "Fix authentication bug in login flow"
created: "2026-01-17T10:30:00Z"
---

## User

Can you help me fix the authentication bug?

## Assistant

I'll help you fix that. Let me look at the code...

```go
func authenticate(user string) error {
    // Fixed implementation
}
```

## User

That worked, thanks!
```

### Main Index (index.md)

```markdown
# Claude Code Logs

| Project | Sessions | Last Activity |
|---------|----------|---------------|
| [my-project](my-project/index.md) | 5 | 2026-01-17 |
| [another-project](another-project/index.md) | 3 | 2026-01-15 |
```

### Project Index ({project}/index.md)

```markdown
# Project: /Users/user/code/my-project

| Session | Title | Created |
|---------|-------|---------|
| [abc123](abc123.md) | Fix authentication bug | 2026-01-17 |
| [def456](def456.md) | Add user registration | 2026-01-16 |
```

## mtime Skip Logic

```go
func (g *MarkdownGenerator) ShouldRegenerate(jsonlPath, mdPath string) bool {
    if g.Force {
        return true
    }

    mdInfo, err := os.Stat(mdPath)
    if os.IsNotExist(err) {
        return true // MD doesn't exist, generate it
    }

    jsonlInfo, err := os.Stat(jsonlPath)
    if err != nil {
        return false // JSONL doesn't exist (orphan case)
    }

    // Regenerate if JSONL is newer than MD
    return jsonlInfo.ModTime().After(mdInfo.ModTime())
}
```

## Orphan Detection Logic

```go
func DetectOrphans(outputDir string, expectedFiles map[string]bool) ([]OrphanInfo, error) {
    var orphans []OrphanInfo

    // Walk output directory for all .md files
    filepath.Walk(outputDir, func(path string, info os.FileInfo, err error) error {
        if !strings.HasSuffix(path, ".md") || info.IsDir() {
            return nil
        }

        // Skip index files
        if strings.HasSuffix(path, "index.md") {
            return nil
        }

        // If not in expected set, it's an orphan
        if !expectedFiles[path] {
            content, _ := os.ReadFile(path)
            fm, _, _ := ParseFrontmatter(content)
            orphans = append(orphans, OrphanInfo{Path: path, Frontmatter: fm})
        }
        return nil
    })

    return orphans, nil
}
```

## Message Formatting Rules

| Message Type | Format |
|--------------|--------|
| User | `## User\n\n{content}` |
| Assistant | `## Assistant\n\n{content}` |
| Code blocks | Preserve original fencing (```lang) |
| Tool calls | `<details><summary>Tool: {name}</summary>\n\n```json\n{args}\n```\n</details>` |
| Tool results | `<details><summary>Result</summary>\n\n```\n{output}\n```\n</details>` |

## Error Handling

| Error Type | Handling |
|------------|----------|
| JSONL parse error | Log warning, skip session, continue with others |
| File write error | Return error, abort generation |
| mtime stat error | Assume regeneration needed |
| Frontmatter parse error | Log warning, use defaults for orphan |
| Hash computation error | Log warning, use "unknown" as hash |

## Security Considerations

| Concern | Approach |
|---------|----------|
| Path traversal | Sanitize project folder names, use filepath.Clean |
| File permissions | Use 0644 for files, 0755 for directories |
| Sensitive data | No filtering (user's own logs); warn in docs about sharing |

## NFR Implementation

| Requirement | Design Approach |
|-------------|-----------------|
| Performance | mtime skip logic avoids unnecessary work |
| Determinism | No timestamps in content body; sorted output |
| Git-friendly | Deterministic output prevents unnecessary diffs |
| Offline | No network dependencies |

## External Dependencies

| Package | Purpose | Import |
|---------|---------|--------|
| crypto/sha256 | Source hash computation | Standard library |
| gopkg.in/yaml.v3 | Frontmatter serialization | Existing dependency |
| os, filepath | File operations | Standard library |

## Integration Points

### With Existing Code

1. **parser.go** - No changes needed; already provides Session/Project structs
2. **cmd_serve.go** - Add `--force` flag; call MarkdownGenerator instead of HTML generator
3. **watcher.go** - Update onChange callback to regenerate MD files
4. **generator.go** - Either refactor or replace entirely with markdown_generator.go

### Output Structure

```text
output/
├── index.md                    # Main index
├── my-project/
│   ├── index.md                # Project index
│   ├── abc123.md               # Session MD
│   └── def456.md               # Session MD
├── another-project/
│   ├── index.md
│   └── xyz789.md
└── shell.html                  # HTML shell for client rendering (Bolt 011)
```
