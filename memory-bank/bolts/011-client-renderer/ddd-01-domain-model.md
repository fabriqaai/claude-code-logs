---
unit: 002-client-renderer
bolt: 011-client-renderer
stage: model
status: complete
updated: 2026-01-17T18:10:00Z
---

# Static Model - Client Renderer

## Bounded Context

The Client Renderer is responsible for transforming the output from a static file server into a dynamic browsing experience. It replaces server-side HTML generation with client-side Markdown rendering while preserving the existing UI/UX. This context is purely client-side (JavaScript/HTML/CSS) and consumes the MD files produced by the Markdown Generator.

## Domain Entities

| Entity | Properties | Business Rules |
|--------|------------|----------------|
| HTMLShell | sidebar, header, contentArea, styles, scripts | Generated once per page type; contains navigation but no message content |
| MarkdownContent | rawText, frontmatter, bodyContent | Fetched via HTTP; parsed client-side |
| RenderedView | htmlContent, highlightedCode | Result of marked.js + highlight.js processing |

## Value Objects

| Value Object | Properties | Constraints |
|--------------|------------|-------------|
| FetchedMarkdown | url, content, loadTime | Immutable once fetched |
| ParsedFrontmatter | source, sourceHash, project, title, created | Extracted from MD content |
| HighlightedCodeBlock | language, code, styledHTML | Immutable after highlighting |
| ClipboardContent | text, format | MD with frontmatter for copy action |
| DownloadBlob | filename, content, mimeType | For download action |

## Aggregates

| Aggregate Root | Members | Invariants |
|----------------|---------|------------|
| PageView | HTMLShell, MarkdownContent, RenderedView | Shell loaded first; content fetched async |
| ActionBar | DownloadButton, CopyButton | Visible only when content loaded |

## Domain Events

| Event | Trigger | Payload |
|-------|---------|---------|
| MarkdownFetched | fetch() complete | url, content, loadTimeMs |
| MarkdownRendered | marked.parse() complete | renderTimeMs |
| CodeHighlighted | hljs.highlightAll() complete | blockCount |
| ContentCopied | Copy button click | success/failure |
| DownloadTriggered | Download button click | filename |

## Domain Services

| Service | Operations | Dependencies |
|---------|------------|--------------|
| MarkdownFetcher | fetch(url) → FetchedMarkdown | Browser fetch API |
| MarkdownParser | parse(md) → {frontmatter, body} | String manipulation |
| MarkdownRenderer | render(md) → HTML | marked.js library |
| CodeHighlighter | highlight(codeBlocks) | highlight.js library |
| ClipboardService | copy(text) → Promise<void> | navigator.clipboard API |
| DownloadService | download(content, filename) | Blob API, anchor element |

## UI Components

| Component | Purpose | Interactions |
|-----------|---------|--------------|
| ContentArea | Displays rendered MD | Receives rendered HTML |
| ActionBar | Contains action buttons | Download, Copy buttons |
| DownloadButton | Triggers MD download | Click → download file |
| CopyButton | Copies MD to clipboard | Click → copy + feedback |
| LoadingIndicator | Shows during fetch | Visible while loading |

## Ubiquitous Language

| Term | Definition |
|------|------------|
| Shell | HTML template with sidebar/navigation but empty content area |
| Fetch | HTTP request to load MD file from server |
| Render | Convert Markdown to HTML using marked.js |
| Highlight | Apply syntax coloring to code blocks using highlight.js |
| Frontmatter | YAML metadata at top of MD file (delimited by `---`) |
| Content Area | Main region where rendered session content displays |
| Action Bar | UI region containing Download and Copy buttons |

## Story Coverage Matrix

| Story | Components | Services | Notes |
|-------|------------|----------|-------|
| 001-md-rendering | HTMLShell, ContentArea, LoadingIndicator | MarkdownFetcher, MarkdownRenderer, CodeHighlighter | Core rendering flow |
| 002-download-button | ActionBar, DownloadButton | DownloadService | File download |
| 003-copy-button | ActionBar, CopyButton | ClipboardService | Clipboard with feedback |

## Client-Side Flow

```
┌─────────────────────────────────────────────────────────────┐
│                     Browser Loads Page                       │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│              HTMLShell Renders (sidebar, CSS, JS)            │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Show Loading Indicator                    │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│               fetch(sessionId.md) → FetchedMarkdown          │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│    ParseFrontmatter → Extract metadata + body content        │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│           marked.parse(body) → HTML content                  │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│        Insert HTML into ContentArea + hljs.highlightAll()    │
└─────────────────────────────┬───────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                 Show ActionBar (Download, Copy)              │
└─────────────────────────────────────────────────────────────┘
```

## Library Configuration

### marked.js Options
```javascript
{
  gfm: true,        // GitHub Flavored Markdown
  breaks: true,     // Convert \n to <br>
  headerIds: true,  // Add IDs to headers
  mangle: false     // Don't mangle email addresses
}
```

### highlight.js Options
```javascript
{
  languages: ['javascript', 'typescript', 'python', 'go', 'bash', 'json', 'yaml', 'html', 'css'],
  ignoreUnescapedHTML: true
}
```

## Button Specifications

### Download Button
- **Icon**: Download arrow SVG
- **Label**: "Download MD"
- **Tooltip**: "Download as Markdown file"
- **Filename**: `{session-id}.md`
- **Content**: Full MD including frontmatter

### Copy Button
- **Icon**: Clipboard SVG
- **Label**: "Copy MD"
- **Tooltip**: "Copy Markdown to clipboard"
- **Feedback**: "Copied!" toast (2 seconds)
- **Content**: Full MD including frontmatter
