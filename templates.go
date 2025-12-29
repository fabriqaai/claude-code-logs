package main

// HTML templates for generating static pages
// All CSS is embedded inline for offline viewing

// baseCSS contains the Claude.ai-inspired styling
const baseCSS = `
/* Claude.ai-inspired styling */
:root {
    --bg-cream: hsl(48, 33.3%, 97.1%);
    --bg-white: #ffffff;
    --accent-orange: hsl(15, 54.2%, 51.2%);
    --text-primary: #1a1a1a;
    --text-secondary: #666666;
    --text-muted: #999999;
    --border-subtle: rgba(0, 0, 0, 0.1);
    --sidebar-width: 288px;
    --font-serif: Georgia, "Times New Roman", Times, serif;
    --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif;
    --font-mono: "SF Mono", Monaco, "Cascadia Code", "Roboto Mono", Consolas, "Courier New", monospace;
}

* {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html, body {
    height: 100%;
    font-family: var(--font-sans);
    font-size: 16px;
    line-height: 1.6;
    color: var(--text-primary);
    background-color: var(--bg-cream);
}

/* Layout */
.container {
    display: flex;
    min-height: 100vh;
}

/* Sidebar */
.sidebar {
    width: var(--sidebar-width);
    min-width: var(--sidebar-width);
    background-color: var(--bg-white);
    border-right: 0.5px solid var(--border-subtle);
    padding: 24px 16px;
    overflow-y: auto;
    position: fixed;
    height: 100vh;
}

.sidebar-header {
    margin-bottom: 24px;
    padding-bottom: 16px;
    border-bottom: 0.5px solid var(--border-subtle);
}

.sidebar-title {
    font-size: 18px;
    font-weight: 600;
    color: var(--text-primary);
    text-decoration: none;
    display: block;
}

.sidebar-title:hover {
    color: var(--accent-orange);
}

.sidebar-subtitle {
    font-size: 12px;
    color: var(--text-muted);
    margin-top: 4px;
}

.project-list {
    list-style: none;
}

.project-item {
    margin-bottom: 8px;
}

.project-link {
    display: block;
    padding: 8px 12px;
    border-radius: 6px;
    text-decoration: none;
    color: var(--text-primary);
    font-size: 14px;
    transition: background-color 0.15s ease;
}

.project-link:hover {
    background-color: var(--bg-cream);
}

.project-link.active {
    background-color: var(--bg-cream);
    font-weight: 500;
}

.project-name {
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.project-meta {
    font-size: 12px;
    color: var(--text-muted);
    margin-top: 2px;
}

.session-list {
    list-style: none;
    margin-top: 8px;
    padding-left: 12px;
    border-left: 2px solid var(--border-subtle);
}

.session-item {
    margin-bottom: 4px;
}

.session-link {
    display: block;
    padding: 6px 10px;
    border-radius: 4px;
    text-decoration: none;
    color: var(--text-secondary);
    font-size: 13px;
    transition: background-color 0.15s ease;
}

.session-link:hover {
    background-color: var(--bg-cream);
    color: var(--text-primary);
}

.session-link.active {
    background-color: var(--bg-cream);
    color: var(--text-primary);
    font-weight: 500;
}

.session-title {
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.session-date {
    font-size: 11px;
    color: var(--text-muted);
}

/* Main content */
.main {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 32px 48px;
    max-width: 900px;
}

.page-header {
    margin-bottom: 32px;
}

.page-title {
    font-size: 28px;
    font-weight: 600;
    color: var(--text-primary);
    margin-bottom: 8px;
}

.page-subtitle {
    font-size: 14px;
    color: var(--text-secondary);
}

/* Static mode banner */
.static-banner {
    background-color: #fff3cd;
    border: 1px solid #ffc107;
    border-radius: 8px;
    padding: 12px 16px;
    margin-bottom: 24px;
    font-size: 14px;
    color: #856404;
}

.static-banner strong {
    font-weight: 600;
}

/* Messages */
.message-list {
    display: flex;
    flex-direction: column;
    gap: 24px;
}

.message {
    padding: 20px 24px;
    border-radius: 12px;
    background-color: var(--bg-white);
    border: 0.5px solid var(--border-subtle);
}

.message-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 12px;
}

.message-role {
    font-size: 13px;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.5px;
}

.message-role.user {
    color: var(--accent-orange);
}

.message-role.assistant {
    color: #6b7280;
}

.message-time {
    font-size: 12px;
    color: var(--text-muted);
}

.message-content {
    font-size: 15px;
    line-height: 1.7;
}

.message.user .message-content {
    font-family: var(--font-sans);
}

.message.assistant .message-content {
    font-family: var(--font-serif);
}

.message-content p {
    margin-bottom: 12px;
}

.message-content p:last-child {
    margin-bottom: 0;
}

/* Code blocks */
.message-content pre {
    background-color: #f6f8fa;
    border: 1px solid var(--border-subtle);
    border-radius: 6px;
    padding: 16px;
    overflow-x: auto;
    margin: 12px 0;
    font-family: var(--font-mono);
    font-size: 13px;
    line-height: 1.5;
}

.message-content code {
    font-family: var(--font-mono);
    font-size: 13px;
    background-color: #f6f8fa;
    padding: 2px 6px;
    border-radius: 4px;
}

.message-content pre code {
    background: none;
    padding: 0;
    border-radius: 0;
}

/* Tool blocks */
.tool-block {
    margin: 12px 0;
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    overflow: hidden;
}

.tool-header {
    background-color: #f6f8fa;
    padding: 10px 14px;
    font-size: 13px;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 8px;
}

.tool-header:hover {
    background-color: #eef1f4;
}

.tool-icon {
    font-size: 14px;
}

.tool-name {
    font-family: var(--font-mono);
    color: var(--text-secondary);
}

.tool-content {
    padding: 12px 14px;
    background-color: var(--bg-white);
    font-family: var(--font-mono);
    font-size: 12px;
    line-height: 1.5;
    max-height: 300px;
    overflow: auto;
    white-space: pre-wrap;
    word-break: break-word;
}

.tool-content.collapsed {
    display: none;
}

/* Empty state */
.empty-state {
    text-align: center;
    padding: 48px 24px;
    color: var(--text-secondary);
}

.empty-state-icon {
    font-size: 48px;
    margin-bottom: 16px;
    opacity: 0.5;
}

.empty-state-title {
    font-size: 18px;
    font-weight: 500;
    margin-bottom: 8px;
    color: var(--text-primary);
}

.empty-state-text {
    font-size: 14px;
}

/* Footer */
.footer {
    margin-top: 48px;
    padding-top: 24px;
    border-top: 0.5px solid var(--border-subtle);
    text-align: center;
    font-size: 13px;
    color: var(--text-muted);
}

.footer a {
    color: var(--accent-orange);
    text-decoration: none;
}

.footer a:hover {
    text-decoration: underline;
}

/* Responsive */
@media (max-width: 768px) {
    .sidebar {
        position: relative;
        width: 100%;
        height: auto;
        border-right: none;
        border-bottom: 0.5px solid var(--border-subtle);
    }

    .main {
        margin-left: 0;
        padding: 24px 16px;
    }

    .container {
        flex-direction: column;
    }
}
`

// indexTemplate is the main page template showing all projects
const indexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Claude Code Logs</title>
    <style>` + baseCSS + `</style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-header">
                <a href="index.html" class="sidebar-title">Claude Code Logs</a>
                <div class="sidebar-subtitle">{{len .Projects}} projects</div>
            </div>
            <ul class="project-list">
                {{range .Projects}}
                <li class="project-item">
                    <a href="{{ProjectSlug .Path}}/index.html" class="project-link">
                        <span class="project-name">{{.Path}}</span>
                        <span class="project-meta">{{len .Sessions}} sessions</span>
                    </a>
                </li>
                {{end}}
            </ul>
        </aside>
        <main class="main">
            <div class="static-banner">
                <strong>Static Mode:</strong> Start the server for search functionality. Run <code>claude-code-logs serve</code>
            </div>
            <header class="page-header">
                <h1 class="page-title">Claude Code Logs</h1>
                <p class="page-subtitle">Browse your Claude Code chat sessions</p>
            </header>
            {{if .Projects}}
            <div class="message-list">
                {{range .Projects}}
                <div class="message">
                    <div class="message-header">
                        <span class="message-role user">{{.Path}}</span>
                        <span class="message-time">{{len .Sessions}} sessions</span>
                    </div>
                    <div class="message-content">
                        {{if .Sessions}}
                        <p>Latest session: <strong>{{(index .Sessions 0).Summary}}</strong></p>
                        <p><a href="{{ProjectSlug .Path}}/index.html">View all sessions →</a></p>
                        {{else}}
                        <p>No sessions found</p>
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">📁</div>
                <h2 class="empty-state-title">No projects found</h2>
                <p class="empty-state-text">Claude Code projects will appear here once you start chatting.</p>
            </div>
            {{end}}
            <footer class="footer">
                <a href="https://fabriqa.ai">claude-code-logs</a> by fabriqa.ai
            </footer>
        </main>
    </div>
</body>
</html>`

// projectIndexTemplate shows all sessions for a project
const projectIndexTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Project.Path}} - Claude Code Logs</title>
    <style>` + baseCSS + `</style>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-header">
                <a href="../index.html" class="sidebar-title">Claude Code Logs</a>
                <div class="sidebar-subtitle">{{len .AllProjects}} projects</div>
            </div>
            <ul class="project-list">
                {{range .AllProjects}}
                <li class="project-item">
                    <a href="../{{ProjectSlug .Path}}/index.html" class="project-link{{if eq .Path $.Project.Path}} active{{end}}">
                        <span class="project-name">{{.Path}}</span>
                        <span class="project-meta">{{len .Sessions}} sessions</span>
                    </a>
                    {{if eq .Path $.Project.Path}}
                    <ul class="session-list">
                        {{range .Sessions}}
                        <li class="session-item">
                            <a href="{{.ID}}.html" class="session-link">
                                <span class="session-title">{{.Summary}}</span>
                                <span class="session-date">{{.CreatedAt.Format "Jan 2, 2006"}}</span>
                            </a>
                        </li>
                        {{end}}
                    </ul>
                    {{end}}
                </li>
                {{end}}
            </ul>
        </aside>
        <main class="main">
            <div class="static-banner">
                <strong>Static Mode:</strong> Start the server for search functionality. Run <code>claude-code-logs serve</code>
            </div>
            <header class="page-header">
                <h1 class="page-title">{{.Project.Path}}</h1>
                <p class="page-subtitle">{{len .Project.Sessions}} sessions</p>
            </header>
            {{if .Project.Sessions}}
            <div class="message-list">
                {{range .Project.Sessions}}
                <div class="message">
                    <div class="message-header">
                        <span class="message-role user">{{.Summary}}</span>
                        <span class="message-time">{{.CreatedAt.Format "Jan 2, 2006 3:04 PM"}}</span>
                    </div>
                    <div class="message-content">
                        <p>{{len .Messages}} messages</p>
                        <p><a href="{{.ID}}.html">View session →</a></p>
                    </div>
                </div>
                {{end}}
            </div>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">💬</div>
                <h2 class="empty-state-title">No sessions found</h2>
                <p class="empty-state-text">Chat sessions for this project will appear here.</p>
            </div>
            {{end}}
            <footer class="footer">
                <a href="https://fabriqa.ai">claude-code-logs</a> by fabriqa.ai
            </footer>
        </main>
    </div>
</body>
</html>`

// sessionTemplate shows a single session's conversation
const sessionTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Session.Summary}} - Claude Code Logs</title>
    <style>` + baseCSS + `</style>
    <script>
    function toggleTool(id) {
        var content = document.getElementById(id);
        if (content) {
            content.classList.toggle('collapsed');
        }
    }
    </script>
</head>
<body>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-header">
                <a href="../../index.html" class="sidebar-title">Claude Code Logs</a>
                <div class="sidebar-subtitle">{{len .AllProjects}} projects</div>
            </div>
            <ul class="project-list">
                {{range .AllProjects}}
                <li class="project-item">
                    <a href="../../{{ProjectSlug .Path}}/index.html" class="project-link{{if eq .Path $.Project.Path}} active{{end}}">
                        <span class="project-name">{{.Path}}</span>
                        <span class="project-meta">{{len .Sessions}} sessions</span>
                    </a>
                    {{if eq .Path $.Project.Path}}
                    <ul class="session-list">
                        {{range .Sessions}}
                        <li class="session-item">
                            <a href="{{.ID}}.html" class="session-link{{if eq .ID $.Session.ID}} active{{end}}">
                                <span class="session-title">{{.Summary}}</span>
                                <span class="session-date">{{.CreatedAt.Format "Jan 2, 2006"}}</span>
                            </a>
                        </li>
                        {{end}}
                    </ul>
                    {{end}}
                </li>
                {{end}}
            </ul>
        </aside>
        <main class="main">
            <div class="static-banner">
                <strong>Static Mode:</strong> Start the server for search functionality. Run <code>claude-code-logs serve</code>
            </div>
            <header class="page-header">
                <h1 class="page-title">{{.Session.Summary}}</h1>
                <p class="page-subtitle">{{.Session.CreatedAt.Format "January 2, 2006 at 3:04 PM"}} · {{len .Session.Messages}} messages</p>
            </header>
            {{if .Session.Messages}}
            <div class="message-list">
                {{range $idx, $msg := .Session.Messages}}
                <div class="message {{$msg.Role}}">
                    <div class="message-header">
                        <span class="message-role {{$msg.Role}}">{{$msg.Role}}</span>
                        <span class="message-time">{{$msg.Timestamp.Format "3:04 PM"}}</span>
                    </div>
                    <div class="message-content">
                        {{range $bidx, $block := $msg.Content}}
                        {{if eq $block.Type "text"}}
                        {{RenderText $block.Text}}
                        {{else if eq $block.Type "tool_use"}}
                        <div class="tool-block">
                            <div class="tool-header" onclick="toggleTool('tool-{{$idx}}-{{$bidx}}')">
                                <span class="tool-icon">🔧</span>
                                <span class="tool-name">{{$block.ToolName}}</span>
                            </div>
                            <div id="tool-{{$idx}}-{{$bidx}}" class="tool-content collapsed">{{$block.ToolInput}}</div>
                        </div>
                        {{else if eq $block.Type "tool_result"}}
                        <div class="tool-block">
                            <div class="tool-header" onclick="toggleTool('result-{{$idx}}-{{$bidx}}')">
                                <span class="tool-icon">📋</span>
                                <span class="tool-name">Result</span>
                            </div>
                            <div id="result-{{$idx}}-{{$bidx}}" class="tool-content collapsed">{{$block.ToolOutput}}</div>
                        </div>
                        {{end}}
                        {{end}}
                    </div>
                </div>
                {{end}}
            </div>
            {{else}}
            <div class="empty-state">
                <div class="empty-state-icon">💬</div>
                <h2 class="empty-state-title">No messages</h2>
                <p class="empty-state-text">This session appears to be empty.</p>
            </div>
            {{end}}
            <footer class="footer">
                <a href="https://fabriqa.ai">claude-code-logs</a> by fabriqa.ai
            </footer>
        </main>
    </div>
</body>
</html>`
