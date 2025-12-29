---
stage: plan
bolt: 002-generator
created: 2025-12-29T16:35:00Z
---

## Implementation Plan: generator

### Objective

Generate static HTML pages from parsed Claude Code session data with Claude.ai visual styling. Create a complete static site that can be viewed offline with project navigation, conversation views, and embedded CSS.

### Deliverables

- `generator.go` - Main generation logic with public API
- `templates.go` - HTML templates as Go template strings
- Output structure:
  ```
  {outputDir}/
  ├── index.html           # Main page with project sidebar
  └── {project-slug}/
      └── {session-id}.html  # Individual session pages
  ```

### Dependencies

- **001-parser**: Provides `[]Project`, `Session`, `Message`, `ContentBlock` types
- **Go stdlib**: `html/template` for secure HTML generation, `os` for file I/O, `path/filepath` for cross-platform paths

### Technical Approach

#### 1. Template Strategy

Use Go `html/template` for security (auto-escaping HTML). Define templates as string constants in `templates.go` to avoid external file dependencies:

- `baseTemplate` - HTML5 boilerplate with CSS
- `indexTemplate` - Project list with sidebar navigation
- `sessionTemplate` - Conversation view for a single session

#### 2. Styling (Claude.ai Visual Design)

CSS embedded inline for offline viewing. Key styles:

| Element | Style |
|---------|-------|
| Background | `hsl(48, 33.3%, 97.1%)` (warm cream) |
| Accent | `hsl(15, 54.2%, 51.2%)` (Claude orange) |
| User messages | Sans-serif font |
| Assistant messages | Serif font (Georgia, Times) |
| Sidebar | 288px fixed width |
| Borders | 0.5px subtle |
| Footer | "claude-code-logs by fabriqa.ai" |

#### 3. Page Structure

**Index Page**:
- Left sidebar (288px): Project list with session counts
- Main content: Welcome message, instructions
- Static mode banner: "Start server for search"

**Session Page**:
- Left sidebar (288px): Project list, session list for current project
- Main content: Message thread in conversation format
- Each message: Role indicator, timestamp, content blocks

#### 4. Content Block Rendering

| Block Type | Rendering |
|------------|-----------|
| `text` | Markdown-like rendering (code blocks, paragraphs) |
| `tool_use` | Collapsible section with tool name and input |
| `tool_result` | Collapsible section with output |

#### 5. Code Syntax Highlighting

Use inline CSS for code blocks (no external dependencies):
- Monospace font
- Background: slightly darker than page
- No external syntax highlighting library (keep it simple)

#### 6. File Generation Strategy

1. Create output directory structure
2. Generate index.html
3. For each project, create project directory
4. For each session, generate session HTML file
5. Use atomic writes (temp file → rename) for reliability

#### 7. URL Scheme

- `index.html` - Main entry point
- `{project-slug}/` - Directory per project (URL-safe project name)
- `{project-slug}/{session-id}.html` - Session page

Project slug derived from path: `/Users/john/project` → `users-john-project`

### Acceptance Criteria

- [ ] Generates valid HTML5 pages
- [ ] Index.html displays project list with session counts
- [ ] 288px sidebar with project navigation
- [ ] Session pages show conversation in Claude.ai style
- [ ] User messages: sans-serif font
- [ ] Assistant messages: serif font
- [ ] Cream background (#faf8f3 / hsl(48, 33.3%, 97.1%))
- [ ] CSS embedded (no external files)
- [ ] Footer: "claude-code-logs by fabriqa.ai" → https://fabriqa.ai
- [ ] "Start server for search" message in static mode
- [ ] Code blocks: monospace, background styling
- [ ] Tool use/result blocks: collapsible display
- [ ] HTML special characters escaped
- [ ] Empty project handling: "No sessions found"
- [ ] Unicode project names display correctly
- [ ] Atomic file writes (temp → rename)
- [ ] Cross-platform path handling

### API Design

```go
// GenerateAll generates the complete static site
func GenerateAll(projects []Project, outputDir string) error

// GenerateIndex generates the main index.html
func GenerateIndex(projects []Project, outputDir string) error

// GenerateSession generates a session HTML file
func GenerateSession(session *Session, project *Project, outputDir string) error

// Helper: create URL-safe slug from project path
func ProjectSlug(path string) string
```

### File Structure

```
claude-code-logs/
├── generator.go       # Generation logic
├── generator_test.go  # Tests for generator
├── templates.go       # HTML template strings
└── ... (existing files)
```

### Risks & Mitigations

| Risk | Mitigation |
|------|------------|
| Large sessions slow to generate | Generate files in parallel (future optimization) |
| Complex markdown in messages | Basic text formatting only, no full markdown parser |
| External CSS dependency | Embed all CSS inline |
| File path edge cases | Use filepath package, test with special characters |

### Out of Scope

- Live search (server responsibility)
- Real-time updates (watcher responsibility)
- Dark mode (future enhancement)
- Full markdown rendering (keep simple)
- Syntax highlighting with language detection (basic monospace only)
