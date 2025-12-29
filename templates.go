package main

// HTML templates for generating static pages
// All CSS is embedded inline for offline viewing

// baseCSS contains the Claude.ai-inspired styling
const baseCSS = `
/* Claude.ai Aesthetic - Warm, Refined, Distinctive */
@import url('https://fonts.googleapis.com/css2?family=DM+Sans:ital,opsz,wght@0,9..40,400;0,9..40,500;0,9..40,600;1,9..40,400&family=Fraunces:ital,opsz,wght@0,9..144,400;0,9..144,500;0,9..144,600;1,9..144,400&family=JetBrains+Mono:wght@400;500&display=swap');

:root {
    /* Warm, inviting color palette */
    --bg-primary: #FAF9F6;
    --bg-secondary: #FFFFFF;
    --bg-tertiary: #F5F4F0;
    --bg-code: #1E1E1E;

    /* Terracotta accent - Claude's signature warmth */
    --accent-primary: #C2410C;
    --accent-hover: #9A3412;
    --accent-subtle: #FED7AA;

    /* Refined text hierarchy */
    --text-primary: #1C1917;
    --text-secondary: #57534E;
    --text-muted: #A8A29E;
    --text-inverse: #FAFAF9;

    /* Subtle, warm borders */
    --border-subtle: rgba(28, 25, 23, 0.08);
    --border-medium: rgba(28, 25, 23, 0.12);

    /* Shadows for depth */
    --shadow-sm: 0 1px 2px rgba(28, 25, 23, 0.04);
    --shadow-md: 0 4px 12px rgba(28, 25, 23, 0.08);
    --shadow-lg: 0 8px 24px rgba(28, 25, 23, 0.12);

    /* Layout */
    --sidebar-width: 280px;
    --content-max-width: 768px;

    /* Typography - Distinctive choices */
    --font-display: 'Fraunces', Georgia, serif;
    --font-body: 'DM Sans', -apple-system, sans-serif;
    --font-mono: 'JetBrains Mono', 'SF Mono', monospace;

    /* Timing */
    --transition-fast: 150ms cubic-bezier(0.4, 0, 0.2, 1);
    --transition-smooth: 300ms cubic-bezier(0.4, 0, 0.2, 1);
}

*, *::before, *::after {
    margin: 0;
    padding: 0;
    box-sizing: border-box;
}

html {
    font-size: 16px;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}

body {
    font-family: var(--font-body);
    font-size: 1rem;
    line-height: 1.6;
    color: var(--text-primary);
    background-color: var(--bg-primary);
    min-height: 100vh;
}

/* Layout */
.container {
    display: flex;
    min-height: 100vh;
}

/* Sidebar - Refined navigation */
.sidebar {
    width: var(--sidebar-width);
    min-width: var(--sidebar-width);
    background-color: var(--bg-secondary);
    border-right: 1px solid var(--border-subtle);
    padding: 24px 16px;
    overflow-y: auto;
    position: fixed;
    height: 100vh;
    display: flex;
    flex-direction: column;
}

.sidebar::-webkit-scrollbar {
    width: 6px;
}

.sidebar::-webkit-scrollbar-track {
    background: transparent;
}

.sidebar::-webkit-scrollbar-thumb {
    background: var(--border-medium);
    border-radius: 3px;
}

.sidebar-header {
    margin-bottom: 32px;
    padding-bottom: 20px;
    border-bottom: 1px solid var(--border-subtle);
}

.sidebar-title {
    font-family: var(--font-display);
    font-size: 1.25rem;
    font-weight: 500;
    color: var(--text-primary);
    text-decoration: none;
    display: flex;
    align-items: center;
    gap: 10px;
    transition: color var(--transition-fast);
}

.sidebar-logo {
    width: 28px;
    height: 28px;
    flex-shrink: 0;
    object-fit: contain;
}

.sidebar-title:hover {
    color: var(--accent-primary);
}

.sidebar-subtitle {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin-top: 6px;
    margin-left: 34px;
    letter-spacing: 0.02em;
}

.project-list {
    list-style: none;
    flex: 1;
}

.project-item {
    margin-bottom: 4px;
}

.project-link {
    display: block;
    padding: 10px 12px;
    border-radius: 8px;
    text-decoration: none;
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 400;
    transition: all var(--transition-fast);
    position: relative;
}

.project-link:hover {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
}

.project-link.active {
    background-color: var(--accent-subtle);
    color: var(--accent-primary);
    font-weight: 500;
}

.project-link.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 20px;
    background: var(--accent-primary);
    border-radius: 0 2px 2px 0;
}

.project-name {
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.project-meta {
    font-size: 0.7rem;
    color: var(--text-muted);
    margin-top: 2px;
    font-weight: 400;
}

.session-list {
    list-style: none;
    margin-top: 8px;
    margin-left: 12px;
    padding-left: 12px;
    border-left: 1px solid var(--border-subtle);
}

.session-item {
    margin-bottom: 2px;
}

.session-link {
    display: block;
    padding: 8px 10px;
    border-radius: 6px;
    text-decoration: none;
    color: var(--text-muted);
    font-size: 0.8rem;
    transition: all var(--transition-fast);
}

.session-link:hover {
    background-color: var(--bg-tertiary);
    color: var(--text-secondary);
}

.session-link.active {
    background-color: var(--bg-tertiary);
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
    font-size: 0.7rem;
    color: var(--text-muted);
    margin-top: 1px;
}

/* Main content area */
.main {
    flex: 1;
    margin-left: var(--sidebar-width);
    padding: 48px;
    max-width: calc(var(--content-max-width) + 96px);
    animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
    from { opacity: 0; transform: translateY(8px); }
    to { opacity: 1; transform: translateY(0); }
}

.page-header {
    margin-bottom: 48px;
}

.page-title {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 8px;
    letter-spacing: -0.02em;
    line-height: 1.2;
}

.page-subtitle {
    font-size: 0.9rem;
    color: var(--text-muted);
    font-weight: 400;
}

/* Messages - Conversation flow */
.message-list {
    display: flex;
    flex-direction: column;
    gap: 32px;
}

.message {
    display: flex;
    gap: 16px;
    animation: messageIn 0.3s ease-out backwards;
}

@keyframes messageIn {
    from { opacity: 0; transform: translateY(12px); }
    to { opacity: 1; transform: translateY(0); }
}

.message:nth-child(1) { animation-delay: 0.05s; }
.message:nth-child(2) { animation-delay: 0.1s; }
.message:nth-child(3) { animation-delay: 0.15s; }
.message:nth-child(4) { animation-delay: 0.2s; }
.message:nth-child(5) { animation-delay: 0.25s; }

/* Avatar styling */
.message::before {
    content: '';
    width: 32px;
    height: 32px;
    border-radius: 50%;
    flex-shrink: 0;
    margin-top: 4px;
}

.message.user::before {
    background: linear-gradient(135deg, #78716C 0%, #57534E 100%);
}

.message.assistant::before {
    background: linear-gradient(135deg, var(--accent-primary) 0%, #EA580C 100%);
}

.message-body {
    flex: 1;
    min-width: 0;
}

.message-header {
    display: flex;
    align-items: center;
    gap: 12px;
    margin-bottom: 8px;
}

.message-role {
    font-size: 0.8rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.message-role.user {
    color: var(--text-secondary);
}

.message-role.assistant {
    color: var(--accent-primary);
}

.message-time {
    font-size: 0.75rem;
    color: var(--text-muted);
}

.message-content {
    font-size: 0.95rem;
    line-height: 1.75;
    color: var(--text-primary);
}

.message.assistant .message-content {
    font-family: var(--font-display);
    font-weight: 400;
}

.message-content p {
    margin-bottom: 16px;
}

.message-content p:last-child {
    margin-bottom: 0;
}

.message-content a {
    color: var(--accent-primary);
    text-decoration: none;
    border-bottom: 1px solid var(--accent-subtle);
    transition: border-color var(--transition-fast);
}

.message-content a:hover {
    border-color: var(--accent-primary);
}

/* Code blocks - Dark, refined */
.message-content pre {
    background-color: var(--bg-code);
    border-radius: 12px;
    padding: 20px 24px;
    overflow-x: auto;
    margin: 20px 0;
    font-family: var(--font-mono);
    font-size: 0.85rem;
    line-height: 1.6;
    color: #E5E5E5;
    box-shadow: var(--shadow-md);
}

.message-content code {
    font-family: var(--font-mono);
    font-size: 0.85em;
    background-color: var(--bg-tertiary);
    color: var(--accent-primary);
    padding: 3px 8px;
    border-radius: 6px;
    font-weight: 500;
}

.message-content pre code {
    background: none;
    color: inherit;
    padding: 0;
    border-radius: 0;
    font-weight: 400;
}

/* Tool blocks - Clean, professional */
.tool-block {
    margin: 20px 0;
    border: 1px solid var(--border-medium);
    border-radius: 12px;
    overflow: hidden;
    background: var(--bg-secondary);
    box-shadow: var(--shadow-sm);
    transition: box-shadow var(--transition-fast);
}

.tool-block:hover {
    box-shadow: var(--shadow-md);
}

.tool-header {
    background-color: var(--bg-tertiary);
    padding: 12px 16px;
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 10px;
    transition: background-color var(--transition-fast);
    user-select: none;
}

.tool-header:hover {
    background-color: var(--border-subtle);
}

.tool-icon {
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--accent-primary);
    color: white;
    border-radius: 4px;
    font-size: 10px;
    font-weight: 600;
}

.tool-icon.result {
    background: #059669;
}

.tool-name {
    font-family: var(--font-mono);
    color: var(--text-secondary);
    font-weight: 500;
}

.tool-toggle {
    margin-left: auto;
    color: var(--text-muted);
    transition: transform var(--transition-fast);
}

.tool-content {
    padding: 16px;
    background-color: var(--bg-code);
    font-family: var(--font-mono);
    font-size: 0.8rem;
    line-height: 1.6;
    max-height: 320px;
    overflow: auto;
    white-space: pre-wrap;
    word-break: break-word;
    color: #E5E5E5;
}

.tool-content.collapsed {
    display: none;
}

/* Cards for project/session listings */
.card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 16px;
    padding: 24px;
    transition: all var(--transition-smooth);
    box-shadow: var(--shadow-sm);
}

.card:hover {
    border-color: var(--border-medium);
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
}

.card-title {
    font-family: var(--font-display);
    font-size: 1.1rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 8px;
}

.card-meta {
    font-size: 0.8rem;
    color: var(--text-muted);
    margin-bottom: 12px;
}

.card-content {
    font-size: 0.9rem;
    color: var(--text-secondary);
    line-height: 1.6;
}

.card-link {
    display: inline-flex;
    align-items: center;
    gap: 6px;
    margin-top: 16px;
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--accent-primary);
    text-decoration: none;
    transition: gap var(--transition-fast);
}

.card-link:hover {
    gap: 10px;
}

.card-link::after {
    content: '→';
}

/* Empty state - Elegant placeholder */
.empty-state {
    text-align: center;
    padding: 80px 32px;
    color: var(--text-secondary);
}

.empty-state-icon {
    width: 64px;
    height: 64px;
    margin: 0 auto 24px;
    background: var(--bg-tertiary);
    border-radius: 16px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 28px;
}

.empty-state-title {
    font-family: var(--font-display);
    font-size: 1.25rem;
    font-weight: 500;
    margin-bottom: 8px;
    color: var(--text-primary);
}

.empty-state-text {
    font-size: 0.9rem;
    color: var(--text-muted);
    max-width: 320px;
    margin: 0 auto;
}

/* Footer - Subtle branding */
.footer {
    margin-top: 64px;
    padding-top: 32px;
    border-top: 1px solid var(--border-subtle);
    text-align: center;
    font-size: 0.8rem;
    color: var(--text-muted);
}

.footer a {
    color: var(--text-secondary);
    text-decoration: none;
    font-weight: 500;
    transition: color var(--transition-fast);
}

.footer a:hover {
    color: var(--accent-primary);
}

/* Responsive design */
@media (max-width: 1024px) {
    .main {
        padding: 32px;
    }
}

@media (max-width: 768px) {
    .sidebar {
        position: relative;
        width: 100%;
        height: auto;
        border-right: none;
        border-bottom: 1px solid var(--border-subtle);
        padding: 20px 16px;
    }

    .main {
        margin-left: 0;
        padding: 24px 16px;
    }

    .container {
        flex-direction: column;
    }

    .page-title {
        font-size: 1.5rem;
    }

    .message {
        gap: 12px;
    }

    .message::before {
        width: 28px;
        height: 28px;
    }
}

/* Selection styling */
::selection {
    background: var(--accent-subtle);
    color: var(--accent-primary);
}

/* Focus styles for accessibility */
:focus-visible {
    outline: 2px solid var(--accent-primary);
    outline-offset: 2px;
}

/* Tree View - Collapsible project hierarchy */
.tree-controls {
    display: flex;
    gap: 8px;
    margin-bottom: 12px;
    padding-bottom: 12px;
    border-bottom: 1px solid var(--border-subtle);
}

.tree-control-btn {
    font-family: var(--font-body);
    font-size: 0.7rem;
    color: var(--text-muted);
    background: none;
    border: 1px solid var(--border-subtle);
    border-radius: 4px;
    padding: 4px 8px;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.tree-control-btn:hover {
    color: var(--text-secondary);
    border-color: var(--border-medium);
    background: var(--bg-tertiary);
}

.tree-control-btn:focus-visible {
    outline: 2px solid var(--accent-primary);
    outline-offset: 1px;
}

.tree-node {
    margin-bottom: 2px;
}

.tree-node-header {
    display: flex;
    align-items: center;
    gap: 4px;
}

.tree-toggle {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 20px;
    height: 20px;
    padding: 0;
    border: none;
    background: none;
    cursor: pointer;
    border-radius: 4px;
    flex-shrink: 0;
    transition: background-color var(--transition-fast);
}

.tree-toggle:hover {
    background: var(--bg-tertiary);
}

.tree-toggle:focus-visible {
    outline: 2px solid var(--accent-primary);
    outline-offset: 1px;
}

.tree-toggle.hidden {
    visibility: hidden;
}

.tree-chevron {
    width: 12px;
    height: 12px;
    color: var(--text-muted);
    transition: transform 200ms ease-out;
}

.tree-node.collapsed .tree-chevron {
    transform: rotate(-90deg);
}

.tree-node-link {
    flex: 1;
    min-width: 0;
    display: flex;
    align-items: center;
    padding: 8px 10px;
    border-radius: 8px;
    text-decoration: none;
    color: var(--text-secondary);
    font-size: 0.875rem;
    font-weight: 400;
    transition: all var(--transition-fast);
    position: relative;
}

.tree-node-link:hover {
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
}

.tree-node-link.active {
    background-color: var(--accent-subtle);
    color: var(--accent-primary);
    font-weight: 500;
}

.tree-node-link.active::before {
    content: '';
    position: absolute;
    left: 0;
    top: 50%;
    transform: translateY(-50%);
    width: 3px;
    height: 20px;
    background: var(--accent-primary);
    border-radius: 0 2px 2px 0;
}

.tree-node-content {
    flex: 1;
    min-width: 0;
}

.tree-node-name {
    display: block;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.tree-node-meta {
    font-size: 0.7rem;
    color: var(--text-muted);
    margin-top: 2px;
    font-weight: 400;
}

.session-count {
    background: var(--accent-subtle);
    color: var(--accent-primary);
    font-size: 0.65rem;
    font-weight: 500;
    padding: 2px 6px;
    border-radius: 10px;
    margin-left: 8px;
    flex-shrink: 0;
}

.tree-children {
    max-height: 2000px;
    overflow: hidden;
    transition: max-height 200ms ease-out, opacity 200ms ease-out;
    opacity: 1;
}

.tree-node.collapsed .tree-children {
    max-height: 0;
    opacity: 0;
}

/* Respect reduced motion preference */
@media (prefers-reduced-motion: reduce) {
    .tree-chevron,
    .tree-children {
        transition: none;
    }
}

/* Search UI */
.search-container {
    margin-bottom: 24px;
}

.search-input-wrapper {
    position: relative;
}

.search-input {
    width: 100%;
    padding: 10px 12px 10px 36px;
    font-family: var(--font-body);
    font-size: 0.875rem;
    border: 1px solid var(--border-medium);
    border-radius: 8px;
    background: var(--bg-primary);
    color: var(--text-primary);
    transition: all var(--transition-fast);
}

.search-input:focus {
    outline: none;
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 3px var(--accent-subtle);
}

.search-input::placeholder {
    color: var(--text-muted);
}

.search-icon {
    position: absolute;
    left: 12px;
    top: 50%;
    transform: translateY(-50%);
    width: 16px;
    height: 16px;
    color: var(--text-muted);
    pointer-events: none;
}

.search-results {
    position: absolute;
    top: 100%;
    left: 0;
    right: 0;
    margin-top: 8px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-medium);
    border-radius: 12px;
    box-shadow: var(--shadow-lg);
    max-height: 400px;
    overflow-y: auto;
    z-index: 100;
    display: none;
}

.search-results.active {
    display: block;
}

.search-result-item {
    padding: 12px 16px;
    border-bottom: 1px solid var(--border-subtle);
    cursor: pointer;
    transition: background-color var(--transition-fast);
}

.search-result-item:last-child {
    border-bottom: none;
}

.search-result-item:hover {
    background-color: var(--bg-tertiary);
}

.search-result-project {
    font-size: 0.7rem;
    color: var(--text-muted);
    margin-bottom: 4px;
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.search-result-title {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 4px;
}

.search-result-content {
    font-size: 0.8rem;
    color: var(--text-secondary);
    line-height: 1.5;
}

.search-result-content mark {
    background: var(--accent-subtle);
    color: var(--accent-primary);
    padding: 1px 3px;
    border-radius: 3px;
}

.search-no-results {
    padding: 24px 16px;
    text-align: center;
    color: var(--text-muted);
    font-size: 0.85rem;
}

.search-loading {
    padding: 24px 16px;
    text-align: center;
    color: var(--text-muted);
    font-size: 0.85rem;
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
                <a href="index.html" class="sidebar-title">
                    <img src="claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .Projects}} projects</div>
            </div>
            <div class="search-container">
                <div class="search-input-wrapper">
                    <svg class="search-icon" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    <input type="text" class="search-input" id="searchInput" placeholder="Search messages..." autocomplete="off">
                    <div class="search-results" id="searchResults"></div>
                </div>
            </div>
            <div class="tree-controls">
                <button type="button" class="tree-control-btn" id="expandAll">Expand All</button>
                <button type="button" class="tree-control-btn" id="collapseAll">Collapse All</button>
            </div>
            <ul class="project-list">
                {{range $project := .Projects}}
                <li class="tree-node" data-project="{{ProjectSlug $project.Path}}">
                    <div class="tree-node-header">
                        <button type="button" class="tree-toggle{{if eq (len $project.Sessions) 0}} hidden{{end}}" aria-expanded="true" aria-label="Toggle {{$project.Path}}">
                            <svg class="tree-chevron" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
                            </svg>
                        </button>
                        <a href="{{ProjectSlug $project.Path}}/index.html" class="tree-node-link">
                            <span class="tree-node-content">
                                <span class="tree-node-name">{{$project.Path}}</span>
                                <span class="tree-node-meta">{{len $project.Sessions}} sessions</span>
                            </span>
                            {{if gt (len $project.Sessions) 0}}<span class="session-count">{{len $project.Sessions}}</span>{{end}}
                        </a>
                    </div>
                    {{if gt (len $project.Sessions) 0}}
                    <ul class="tree-children session-list">
                        {{range $project.Sessions}}
                        <li class="session-item">
                            <a href="{{ProjectSlug $project.Path}}/{{.ID}}.html" class="session-link">
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
    <script>
    (function() {
        // Tree view toggle functionality
        document.querySelectorAll('.tree-toggle').forEach(function(btn) {
            btn.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                var node = btn.closest('.tree-node');
                var isCollapsed = node.classList.toggle('collapsed');
                btn.setAttribute('aria-expanded', !isCollapsed);
            });
        });

        // Expand All / Collapse All
        var expandAllBtn = document.getElementById('expandAll');
        var collapseAllBtn = document.getElementById('collapseAll');

        if (expandAllBtn) {
            expandAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.remove('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'true');
                });
            });
        }

        if (collapseAllBtn) {
            collapseAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.add('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'false');
                });
            });
        }

        // Search functionality
        var searchInput = document.getElementById('searchInput');
        var searchResults = document.getElementById('searchResults');
        var debounceTimer;
        var baseUrl = '';

        searchInput.addEventListener('input', function() {
            clearTimeout(debounceTimer);
            var query = this.value.trim();

            if (query.length < 2) {
                searchResults.classList.remove('active');
                return;
            }

            debounceTimer = setTimeout(function() {
                searchResults.innerHTML = '<div class="search-loading">Searching...</div>';
                searchResults.classList.add('active');

                fetch('/api/search', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: query, limit: 20 })
                })
                .then(function(r) { return r.json(); })
                .then(function(data) {
                    if (!data.results || data.results.length === 0) {
                        searchResults.innerHTML = '<div class="search-no-results">No results found</div>';
                        return;
                    }
                    var html = '';
                    data.results.forEach(function(r) {
                        r.matches.slice(0, 3).forEach(function(m) {
                            html += '<a href="' + baseUrl + r.projectSlug + '/' + r.sessionId + '.html" class="search-result-item">';
                            html += '<div class="search-result-project">' + r.project + '</div>';
                            html += '<div class="search-result-title">' + r.sessionTitle + '</div>';
                            html += '<div class="search-result-content">' + m.content.substring(0, 150) + '...</div>';
                            html += '</a>';
                        });
                    });
                    searchResults.innerHTML = html;
                })
                .catch(function() {
                    searchResults.innerHTML = '<div class="search-no-results">Search error</div>';
                });
            }, 300);
        });

        document.addEventListener('click', function(e) {
            if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
                searchResults.classList.remove('active');
            }
        });

        searchInput.addEventListener('keydown', function(e) {
            if (e.key === 'Escape') {
                searchResults.classList.remove('active');
                searchInput.blur();
            }
        });
    })();
    </script>
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
                <a href="../index.html" class="sidebar-title">
                    <img src="../claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .AllProjects}} projects</div>
            </div>
            <div class="search-container">
                <div class="search-input-wrapper">
                    <svg class="search-icon" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    <input type="text" class="search-input" id="searchInput" placeholder="Search messages..." autocomplete="off">
                    <div class="search-results" id="searchResults"></div>
                </div>
            </div>
            <div class="tree-controls">
                <button type="button" class="tree-control-btn" id="expandAll">Expand All</button>
                <button type="button" class="tree-control-btn" id="collapseAll">Collapse All</button>
            </div>
            <ul class="project-list">
                {{range .AllProjects}}
                <li class="tree-node{{if ne .Path $.Project.Path}} collapsed{{end}}" data-project="{{ProjectSlug .Path}}">
                    <div class="tree-node-header">
                        <button type="button" class="tree-toggle{{if eq (len .Sessions) 0}} hidden{{end}}" aria-expanded="{{if eq .Path $.Project.Path}}true{{else}}false{{end}}" aria-label="Toggle {{.Path}}">
                            <svg class="tree-chevron" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
                            </svg>
                        </button>
                        <a href="../{{ProjectSlug .Path}}/index.html" class="tree-node-link{{if eq .Path $.Project.Path}} active{{end}}">
                            <span class="tree-node-content">
                                <span class="tree-node-name">{{.Path}}</span>
                                <span class="tree-node-meta">{{len .Sessions}} sessions</span>
                            </span>
                            {{if gt (len .Sessions) 0}}<span class="session-count">{{len .Sessions}}</span>{{end}}
                        </a>
                    </div>
                    {{if gt (len .Sessions) 0}}
                    <ul class="tree-children session-list">
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
    <script>
    (function() {
        // Tree view toggle functionality
        document.querySelectorAll('.tree-toggle').forEach(function(btn) {
            btn.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                var node = btn.closest('.tree-node');
                var isCollapsed = node.classList.toggle('collapsed');
                btn.setAttribute('aria-expanded', !isCollapsed);
            });
        });

        // Expand All / Collapse All
        var expandAllBtn = document.getElementById('expandAll');
        var collapseAllBtn = document.getElementById('collapseAll');

        if (expandAllBtn) {
            expandAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.remove('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'true');
                });
            });
        }

        if (collapseAllBtn) {
            collapseAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.add('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'false');
                });
            });
        }

        // Search functionality
        var searchInput = document.getElementById('searchInput');
        var searchResults = document.getElementById('searchResults');
        var debounceTimer;
        var baseUrl = '../';

        searchInput.addEventListener('input', function() {
            clearTimeout(debounceTimer);
            var query = this.value.trim();

            if (query.length < 2) {
                searchResults.classList.remove('active');
                return;
            }

            debounceTimer = setTimeout(function() {
                searchResults.innerHTML = '<div class="search-loading">Searching...</div>';
                searchResults.classList.add('active');

                fetch('/api/search', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: query, limit: 20 })
                })
                .then(function(r) { return r.json(); })
                .then(function(data) {
                    if (!data.results || data.results.length === 0) {
                        searchResults.innerHTML = '<div class="search-no-results">No results found</div>';
                        return;
                    }
                    var html = '';
                    data.results.forEach(function(r) {
                        r.matches.slice(0, 3).forEach(function(m) {
                            html += '<a href="' + baseUrl + r.projectSlug + '/' + r.sessionId + '.html" class="search-result-item">';
                            html += '<div class="search-result-project">' + r.project + '</div>';
                            html += '<div class="search-result-title">' + r.sessionTitle + '</div>';
                            html += '<div class="search-result-content">' + m.content.substring(0, 150) + '...</div>';
                            html += '</a>';
                        });
                    });
                    searchResults.innerHTML = html;
                })
                .catch(function() {
                    searchResults.innerHTML = '<div class="search-no-results">Search error</div>';
                });
            }, 300);
        });

        document.addEventListener('click', function(e) {
            if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
                searchResults.classList.remove('active');
            }
        });

        searchInput.addEventListener('keydown', function(e) {
            if (e.key === 'Escape') {
                searchResults.classList.remove('active');
                searchInput.blur();
            }
        });
    })();
    </script>
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
                <a href="../../index.html" class="sidebar-title">
                    <img src="../../claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .AllProjects}} projects</div>
            </div>
            <div class="search-container">
                <div class="search-input-wrapper">
                    <svg class="search-icon" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    <input type="text" class="search-input" id="searchInput" placeholder="Search messages..." autocomplete="off">
                    <div class="search-results" id="searchResults"></div>
                </div>
            </div>
            <div class="tree-controls">
                <button type="button" class="tree-control-btn" id="expandAll">Expand All</button>
                <button type="button" class="tree-control-btn" id="collapseAll">Collapse All</button>
            </div>
            <ul class="project-list">
                {{range .AllProjects}}
                <li class="tree-node{{if ne .Path $.Project.Path}} collapsed{{end}}" data-project="{{ProjectSlug .Path}}">
                    <div class="tree-node-header">
                        <button type="button" class="tree-toggle{{if eq (len .Sessions) 0}} hidden{{end}}" aria-expanded="{{if eq .Path $.Project.Path}}true{{else}}false{{end}}" aria-label="Toggle {{.Path}}">
                            <svg class="tree-chevron" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
                            </svg>
                        </button>
                        <a href="../../{{ProjectSlug .Path}}/index.html" class="tree-node-link{{if eq .Path $.Project.Path}} active{{end}}">
                            <span class="tree-node-content">
                                <span class="tree-node-name">{{.Path}}</span>
                                <span class="tree-node-meta">{{len .Sessions}} sessions</span>
                            </span>
                            {{if gt (len .Sessions) 0}}<span class="session-count">{{len .Sessions}}</span>{{end}}
                        </a>
                    </div>
                    {{if gt (len .Sessions) 0}}
                    <ul class="tree-children session-list">
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
            <header class="page-header">
                <h1 class="page-title">{{.Session.Summary}}</h1>
                <p class="page-subtitle">{{.Session.CreatedAt.Format "January 2, 2006 at 3:04 PM"}} · {{len .Session.Messages}} messages</p>
            </header>
            {{if .Session.Messages}}
            <div class="message-list">
                {{range $idx, $msg := .Session.Messages}}
                <div class="message {{$msg.Role}}">
                    <div class="message-body">
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
                                    <span class="tool-icon">T</span>
                                    <span class="tool-name">{{$block.ToolName}}</span>
                                    <span class="tool-toggle">▼</span>
                                </div>
                                <div id="tool-{{$idx}}-{{$bidx}}" class="tool-content collapsed">{{$block.ToolInput}}</div>
                            </div>
                            {{else if eq $block.Type "tool_result"}}
                            <div class="tool-block">
                                <div class="tool-header" onclick="toggleTool('result-{{$idx}}-{{$bidx}}')">
                                    <span class="tool-icon result">R</span>
                                    <span class="tool-name">Result</span>
                                    <span class="tool-toggle">▼</span>
                                </div>
                                <div id="result-{{$idx}}-{{$bidx}}" class="tool-content collapsed">{{$block.ToolOutput}}</div>
                            </div>
                            {{end}}
                            {{end}}
                        </div>
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
    <script>
    (function() {
        // Tree view toggle functionality
        document.querySelectorAll('.tree-toggle').forEach(function(btn) {
            btn.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                var node = btn.closest('.tree-node');
                var isCollapsed = node.classList.toggle('collapsed');
                btn.setAttribute('aria-expanded', !isCollapsed);
            });
        });

        // Expand All / Collapse All
        var expandAllBtn = document.getElementById('expandAll');
        var collapseAllBtn = document.getElementById('collapseAll');

        if (expandAllBtn) {
            expandAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.remove('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'true');
                });
            });
        }

        if (collapseAllBtn) {
            collapseAllBtn.addEventListener('click', function() {
                document.querySelectorAll('.tree-node').forEach(function(n) {
                    n.classList.add('collapsed');
                    var toggle = n.querySelector('.tree-toggle');
                    if (toggle) toggle.setAttribute('aria-expanded', 'false');
                });
            });
        }

        // Search functionality
        var searchInput = document.getElementById('searchInput');
        var searchResults = document.getElementById('searchResults');
        var debounceTimer;
        var baseUrl = '../../';

        searchInput.addEventListener('input', function() {
            clearTimeout(debounceTimer);
            var query = this.value.trim();

            if (query.length < 2) {
                searchResults.classList.remove('active');
                return;
            }

            debounceTimer = setTimeout(function() {
                searchResults.innerHTML = '<div class="search-loading">Searching...</div>';
                searchResults.classList.add('active');

                fetch('/api/search', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ query: query, limit: 20 })
                })
                .then(function(r) { return r.json(); })
                .then(function(data) {
                    if (!data.results || data.results.length === 0) {
                        searchResults.innerHTML = '<div class="search-no-results">No results found</div>';
                        return;
                    }
                    var html = '';
                    data.results.forEach(function(r) {
                        r.matches.slice(0, 3).forEach(function(m) {
                            html += '<a href="' + baseUrl + r.projectSlug + '/' + r.sessionId + '.html" class="search-result-item">';
                            html += '<div class="search-result-project">' + r.project + '</div>';
                            html += '<div class="search-result-title">' + r.sessionTitle + '</div>';
                            html += '<div class="search-result-content">' + m.content.substring(0, 150) + '...</div>';
                            html += '</a>';
                        });
                    });
                    searchResults.innerHTML = html;
                })
                .catch(function() {
                    searchResults.innerHTML = '<div class="search-no-results">Search error</div>';
                });
            }, 300);
        });

        document.addEventListener('click', function(e) {
            if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
                searchResults.classList.remove('active');
            }
        });

        searchInput.addEventListener('keydown', function(e) {
            if (e.key === 'Escape') {
                searchResults.classList.remove('active');
                searchInput.blur();
            }
        });
    })();
    </script>
</body>
</html>`
