---
unit: 002-client-renderer
bolt: 011-client-renderer
stage: design
status: complete
updated: 2026-01-17T18:15:00Z
---

# Technical Design - Client Renderer

## Architecture Pattern

**Pattern**: Static Shell + Dynamic Content Loading

The server generates an HTML shell (sidebar, CSS, JS) while session content is loaded dynamically via fetch() and rendered client-side with marked.js. This preserves the existing UI/UX while enabling MD-first output.

```text
┌─────────────────────────────────────────────────────────────┐
│                      Go Server                               │
│  - Serves static HTML shells                                 │
│  - Serves MD files                                           │
│  - Serves index.html, project/index.html, session.html       │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                    Browser (Client)                          │
│  - Loads HTML shell (sidebar, CSS, empty content area)       │
│  - Fetches session.md via JavaScript                         │
│  - Renders MD → HTML using marked.js                         │
│  - Highlights code using highlight.js                        │
└─────────────────────────────────────────────────────────────┘
```

## File Structure

### New Files
```text
claude-code-logs/
├── templates_shell.go         # HTML shell templates (replaces session template content)
```

### Modified Files
```text
├── markdown_generator.go      # Add shell.html generation
├── templates_css.go           # Add styles for action buttons, MD content
```

### Library Dependencies (CDN)
```text
- marked.js (https://cdn.jsdelivr.net/npm/marked/marked.min.js)
- highlight.js (https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/)
```

## HTML Shell Structure

### Session Shell (`{session-id}.html`)

```html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>{{.Session.Summary}} - Claude Code Logs</title>
    <style>/* Existing CSS */</style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
</head>
<body>
    <!-- Mobile menu button -->
    <button class="mobile-menu-btn">...</button>
    <div class="sidebar-backdrop"></div>

    <div class="container">
        <!-- Sidebar (same as current) -->
        <aside class="sidebar">
            <!-- Navigation tree -->
        </aside>

        <!-- Main content area -->
        <main class="main">
            <header class="page-header">
                <div class="page-breadcrumb">...</div>
                <h1 class="page-title" id="session-title">Loading...</h1>
                <p class="page-subtitle" id="session-meta"></p>
                <div class="page-actions" id="action-bar" style="display:none;">
                    <button class="action-btn" id="download-btn" title="Download Markdown">
                        <svg>...</svg> Download MD
                    </button>
                    <button class="action-btn" id="copy-btn" title="Copy Markdown">
                        <svg>...</svg> Copy MD
                    </button>
                </div>
            </header>

            <!-- Content loaded dynamically -->
            <div id="content-area" class="message-list">
                <div class="loading-indicator">Loading session...</div>
            </div>

            <footer class="footer">...</footer>
        </main>
    </div>

    <!-- Libraries -->
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>

    <!-- Client-side rendering logic -->
    <script>
        // MD URL derived from current page
        const mdUrl = '{{.Session.ID}}.md';
        // ... rendering logic
    </script>
</body>
</html>
```

## JavaScript Implementation

### Main Rendering Script

```javascript
(function() {
    const mdUrl = '{{.Session.ID}}.md';
    const contentArea = document.getElementById('content-area');
    const titleEl = document.getElementById('session-title');
    const metaEl = document.getElementById('session-meta');
    const actionBar = document.getElementById('action-bar');
    const downloadBtn = document.getElementById('download-btn');
    const copyBtn = document.getElementById('copy-btn');

    let rawMarkdown = '';  // Store for download/copy

    // Configure marked
    marked.setOptions({
        gfm: true,
        breaks: true,
        headerIds: true,
        mangle: false,
        highlight: function(code, lang) {
            if (lang && hljs.getLanguage(lang)) {
                return hljs.highlight(code, { language: lang }).value;
            }
            return hljs.highlightAuto(code).value;
        }
    });

    // Fetch and render
    fetch(mdUrl)
        .then(response => {
            if (!response.ok) throw new Error('Failed to load');
            return response.text();
        })
        .then(md => {
            rawMarkdown = md;
            const { frontmatter, body } = parseFrontmatter(md);

            // Update header
            if (frontmatter.title) {
                titleEl.textContent = frontmatter.title;
                document.title = frontmatter.title + ' - Claude Code Logs';
            }
            if (frontmatter.created) {
                metaEl.textContent = formatDate(frontmatter.created);
            }

            // Render markdown
            contentArea.innerHTML = marked.parse(body);

            // Highlight code blocks
            contentArea.querySelectorAll('pre code').forEach(block => {
                hljs.highlightElement(block);
            });

            // Show action buttons
            actionBar.style.display = 'flex';
        })
        .catch(err => {
            contentArea.innerHTML = '<div class="error-state">Failed to load session</div>';
        });

    // Parse frontmatter
    function parseFrontmatter(md) {
        if (!md.startsWith('---\n')) {
            return { frontmatter: {}, body: md };
        }
        const end = md.indexOf('\n---\n', 4);
        if (end === -1) {
            return { frontmatter: {}, body: md };
        }
        const yaml = md.slice(4, end);
        const body = md.slice(end + 5);

        // Simple YAML parsing for our known fields
        const frontmatter = {};
        yaml.split('\n').forEach(line => {
            const match = line.match(/^(\w+):\s*"?([^"]*)"?$/);
            if (match) {
                frontmatter[match[1]] = match[2];
            }
        });

        return { frontmatter, body };
    }

    // Format date
    function formatDate(isoDate) {
        const date = new Date(isoDate);
        return date.toLocaleDateString('en-US', {
            year: 'numeric', month: 'long', day: 'numeric',
            hour: 'numeric', minute: '2-digit'
        });
    }

    // Download button
    downloadBtn.addEventListener('click', function() {
        const blob = new Blob([rawMarkdown], { type: 'text/markdown' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = mdUrl;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        URL.revokeObjectURL(url);
    });

    // Copy button
    copyBtn.addEventListener('click', function() {
        navigator.clipboard.writeText(rawMarkdown).then(() => {
            const originalText = copyBtn.innerHTML;
            copyBtn.innerHTML = '<svg>...</svg> Copied!';
            copyBtn.classList.add('copied');
            setTimeout(() => {
                copyBtn.innerHTML = originalText;
                copyBtn.classList.remove('copied');
            }, 2000);
        }).catch(err => {
            console.error('Copy failed:', err);
        });
    });
})();
```

## CSS Additions

### Action Bar Styles

```css
.page-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
}

.action-btn {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    padding: 8px 12px;
    background: var(--bg-tertiary);
    border: 1px solid var(--border-subtle);
    border-radius: 6px;
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.875rem;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.action-btn:hover {
    background: var(--bg-secondary);
    border-color: var(--border-medium);
    color: var(--text-primary);
}

.action-btn.copied {
    background: var(--accent-subtle);
    border-color: var(--accent-primary);
    color: var(--accent-primary);
}

.action-btn svg {
    width: 16px;
    height: 16px;
}
```

### Loading and Error States

```css
.loading-indicator {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 48px;
    color: var(--text-muted);
    font-size: 1rem;
}

.error-state {
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 48px;
    color: var(--accent-primary);
    background: var(--accent-subtle);
    border-radius: 8px;
    margin: 24px;
}
```

### Markdown Content Styles

```css
/* Rendered markdown in content area */
#content-area h2 {
    font-family: var(--font-display);
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 24px 0 12px;
    padding-bottom: 8px;
    border-bottom: 1px solid var(--border-subtle);
}

#content-area p {
    margin: 12px 0;
    line-height: 1.7;
}

#content-area pre {
    background: var(--bg-code);
    border-radius: 8px;
    padding: 16px;
    overflow-x: auto;
    margin: 16px 0;
}

#content-area code {
    font-family: var(--font-mono);
    font-size: 0.875rem;
}

#content-area pre code {
    background: none;
    padding: 0;
    color: var(--text-inverse);
}

#content-area details {
    background: var(--bg-tertiary);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    margin: 12px 0;
}

#content-area details summary {
    padding: 12px 16px;
    cursor: pointer;
    font-weight: 500;
    color: var(--text-secondary);
}

#content-area details[open] summary {
    border-bottom: 1px solid var(--border-subtle);
}

#content-area details > *:not(summary) {
    padding: 0 16px 16px;
}
```

## Index Pages

Index pages (`index.html`, `{project}/index.html`) remain **server-rendered** since they contain project/session lists that don't need client-side MD rendering. They link to session HTML shells.

## Data Flow

```text
User clicks session link
        │
        ▼
Browser loads {session-id}.html (shell)
        │
        ├── Sidebar renders immediately (from shell)
        ├── Content area shows "Loading..."
        │
        ▼
JavaScript fetches {session-id}.md
        │
        ▼
Parse frontmatter → Update title/meta
        │
        ▼
marked.parse(body) → HTML
        │
        ▼
Insert HTML into #content-area
        │
        ▼
hljs.highlightAll() on code blocks
        │
        ▼
Show action buttons (Download, Copy)
```

## Error Handling

| Error | User Feedback |
|-------|---------------|
| Network failure | "Failed to load session" error state |
| Invalid MD | Render raw text without parsing |
| Clipboard denied | Console warning, no UI feedback |

## Performance Targets

| Metric | Target |
|--------|--------|
| Shell load | < 200ms |
| MD fetch | < 100ms (local) |
| Render time | < 50ms |
| Highlight time | < 30ms |

## Browser Compatibility

| Browser | Support |
|---------|---------|
| Chrome 90+ | ✅ Full |
| Firefox 88+ | ✅ Full |
| Safari 14+ | ✅ Full |
| Edge 90+ | ✅ Full |

Uses standard APIs: fetch, Blob, clipboard API (with fallback).

## Security Considerations

| Concern | Mitigation |
|---------|------------|
| XSS in MD | marked.js sanitizes by default |
| Script injection | No eval, content is text-only |
| Clipboard access | Uses async clipboard API with permissions |
