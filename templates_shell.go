package main

// sessionShellTemplate is the HTML shell for session pages with client-side MD rendering
// The actual content is loaded via JavaScript from the corresponding .md file
const sessionShellTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Session.Summary}} - Claude Code Logs</title>
    <style>` + baseCSS + shellCSS + `</style>
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/styles/github-dark.min.css">
</head>
<body>
    <button class="mobile-menu-btn" aria-label="Open navigation" aria-expanded="false">
        <svg class="hamburger-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M3 6h18M3 12h18M3 18h18"/>
        </svg>
    </button>
    <div class="sidebar-backdrop" aria-hidden="true"></div>
    <div class="container">
        <aside class="sidebar">
            <div class="sidebar-resize-handle" role="separator" aria-orientation="vertical"></div>
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
                <div class="page-breadcrumb">
                    <a href="../../index.html">Home</a>
                    <span class="page-breadcrumb-separator">/</span>
                    <a href="index.html">{{.Project.Path}}</a>
                </div>
                <h1 class="page-title" id="session-title">{{.Session.Summary}}</h1>
                <p class="page-subtitle" id="session-meta">{{.Session.CreatedAt.Format "January 2, 2006 at 3:04 PM"}} · {{len .Session.Messages}} messages</p>
                <div class="page-actions" id="action-bar">
                    <a href="file://{{.Session.SourcePath}}" class="page-action-link" title="{{.Session.SourcePath}}">
                        <svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4z" clip-rule="evenodd"/></svg>
                        View Source JSONL
                    </a>
                    <button class="action-btn" id="download-btn" title="Download Markdown">
                        <svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M3 17a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm3.293-7.707a1 1 0 011.414 0L9 10.586V3a1 1 0 112 0v7.586l1.293-1.293a1 1 0 111.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>
                        Download MD
                    </button>
                    <button class="action-btn" id="copy-btn" title="Copy Markdown to clipboard">
                        <svg viewBox="0 0 20 20" fill="currentColor"><path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z"/><path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z"/></svg>
                        <span class="copy-btn-text">Copy MD</span>
                    </button>
                </div>
            </header>

            <!-- Content loaded dynamically from MD file -->
            <div id="content-area" class="message-list">
                <div class="loading-indicator">
                    <div class="loading-spinner"></div>
                    <span>Loading session...</span>
                </div>
            </div>

            <footer class="footer">
                <a href="https://github.com/fabriqaai/claude-code-logs">claude-code-logs</a> by <a href="https://fabriqa.ai">fabriqa.ai</a>
            </footer>
        </main>
    </div>

    <!-- Libraries -->
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <script src="https://cdnjs.cloudflare.com/ajax/libs/highlight.js/11.9.0/highlight.min.js"></script>

    <!-- Client-side rendering logic -->
    <script>
    (function() {
        var mdUrl = '{{.Session.ID}}.md';
        var contentArea = document.getElementById('content-area');
        var downloadBtn = document.getElementById('download-btn');
        var copyBtn = document.getElementById('copy-btn');
        var copyBtnText = copyBtn.querySelector('.copy-btn-text');

        var rawMarkdown = '';  // Store for download/copy

        // Configure marked
        marked.setOptions({
            gfm: true,
            breaks: true,
            headerIds: true,
            mangle: false
        });

        // Custom renderer for code blocks with highlight.js
        var renderer = new marked.Renderer();
        renderer.code = function(code, lang) {
            var highlighted;
            if (lang && hljs.getLanguage(lang)) {
                highlighted = hljs.highlight(code, { language: lang }).value;
            } else {
                highlighted = hljs.highlightAuto(code).value;
            }
            return '<pre><code class="hljs ' + (lang || '') + '">' + highlighted + '</code></pre>';
        };
        marked.setOptions({ renderer: renderer });

        // Fetch and render
        fetch(mdUrl)
            .then(function(response) {
                if (!response.ok) throw new Error('Failed to load');
                return response.text();
            })
            .then(function(md) {
                rawMarkdown = md;
                var parsed = parseFrontmatter(md);

                // Render markdown body
                contentArea.innerHTML = marked.parse(parsed.body);

                // Style the rendered content for conversation display
                styleContent();
            })
            .catch(function(err) {
                contentArea.innerHTML = '<div class="error-state"><div class="error-icon">!</div><h2>Failed to load session</h2><p>The markdown file could not be loaded.</p></div>';
                console.error('Failed to load markdown:', err);
            });

        // Parse frontmatter
        function parseFrontmatter(md) {
            if (!md.startsWith('---\n')) {
                return { frontmatter: {}, body: md };
            }
            var end = md.indexOf('\n---\n', 4);
            if (end === -1) {
                return { frontmatter: {}, body: md };
            }
            var yaml = md.slice(4, end);
            var body = md.slice(end + 5);

            // Simple YAML parsing for our known fields
            var frontmatter = {};
            yaml.split('\n').forEach(function(line) {
                var match = line.match(/^(\w+):\s*"?([^"]*)"?$/);
                if (match) {
                    frontmatter[match[1]] = match[2];
                }
            });

            return { frontmatter: frontmatter, body: body };
        }

        // Style the content for conversation display
        function styleContent() {
            // Convert ## User and ## Assistant headers to message blocks
            var headers = contentArea.querySelectorAll('h2');
            headers.forEach(function(h2) {
                var role = h2.textContent.toLowerCase();
                if (role === 'user' || role === 'assistant') {
                    var messageDiv = document.createElement('div');
                    messageDiv.className = 'message ' + role;

                    var bodyDiv = document.createElement('div');
                    bodyDiv.className = 'message-body';

                    var headerDiv = document.createElement('div');
                    headerDiv.className = 'message-header';
                    headerDiv.innerHTML = '<span class="message-role ' + role + '">' + role.toUpperCase() + '</span>';
                    bodyDiv.appendChild(headerDiv);

                    var contentDiv = document.createElement('div');
                    contentDiv.className = 'message-content';

                    // Collect all siblings until next h2
                    var sibling = h2.nextElementSibling;
                    while (sibling && sibling.tagName !== 'H2') {
                        var next = sibling.nextElementSibling;
                        contentDiv.appendChild(sibling);
                        sibling = next;
                    }

                    bodyDiv.appendChild(contentDiv);
                    messageDiv.appendChild(bodyDiv);
                    h2.parentNode.replaceChild(messageDiv, h2);
                }
            });

            // Style details elements as tool blocks
            var details = contentArea.querySelectorAll('details');
            details.forEach(function(detail) {
                detail.classList.add('tool-block');
                var summary = detail.querySelector('summary');
                if (summary) {
                    var text = summary.textContent;
                    var isResult = text.toLowerCase().includes('result');
                    summary.innerHTML = '<span class="tool-icon' + (isResult ? ' result' : '') + '">' + (isResult ? 'R' : 'T') + '</span>' +
                        '<span class="tool-name">' + text + '</span>' +
                        '<span class="tool-toggle">▼</span>';
                }
            });
        }

        // Download button
        downloadBtn.addEventListener('click', function() {
            if (!rawMarkdown) return;
            var blob = new Blob([rawMarkdown], { type: 'text/markdown' });
            var url = URL.createObjectURL(blob);
            var a = document.createElement('a');
            a.href = url;
            a.download = mdUrl;
            document.body.appendChild(a);
            a.click();
            document.body.removeChild(a);
            URL.revokeObjectURL(url);
        });

        // Copy button
        copyBtn.addEventListener('click', function() {
            if (!rawMarkdown) return;
            navigator.clipboard.writeText(rawMarkdown).then(function() {
                var originalText = copyBtnText.textContent;
                copyBtnText.textContent = 'Copied!';
                copyBtn.classList.add('copied');
                setTimeout(function() {
                    copyBtnText.textContent = originalText;
                    copyBtn.classList.remove('copied');
                }, 2000);
            }).catch(function(err) {
                console.error('Copy failed:', err);
            });
        });
    })();
    </script>

    <!-- Sidebar and search functionality (same as existing) -->
    <script>` + sidebarJS + `</script>
</body>
</html>`

// sidebarJS contains the sidebar, tree, and search functionality
const sidebarJS = `
(function() {
    var STORAGE_KEY = 'claude-code-logs-sidebar';
    var sidebar = document.querySelector('.sidebar');
    var resizeHandle = document.querySelector('.sidebar-resize-handle');
    var menuBtn = document.querySelector('.mobile-menu-btn');
    var backdrop = document.querySelector('.sidebar-backdrop');

    // localStorage helpers
    function loadSidebarState() {
        try {
            var data = JSON.parse(localStorage.getItem(STORAGE_KEY) || '{}');
            if (data.width && data.width >= 200 && data.width <= 500) {
                document.documentElement.style.setProperty('--sidebar-width', data.width + 'px');
            }
            if (Array.isArray(data.collapsed)) {
                data.collapsed.forEach(function(id) {
                    var node = document.querySelector('[data-project="' + id + '"]');
                    if (node) {
                        node.classList.add('collapsed');
                        var toggle = node.querySelector('.tree-toggle');
                        if (toggle) toggle.setAttribute('aria-expanded', 'false');
                    }
                });
            }
        } catch (e) { console.warn('Failed to load sidebar state:', e); }
    }

    function saveSidebarState() {
        try {
            var width = parseInt(getComputedStyle(document.documentElement).getPropertyValue('--sidebar-width'));
            var collapsed = Array.from(document.querySelectorAll('.tree-node.collapsed'))
                .map(function(n) { return n.dataset.project; })
                .filter(Boolean);
            localStorage.setItem(STORAGE_KEY, JSON.stringify({ width: width, collapsed: collapsed }));
        } catch (e) { console.warn('Failed to save sidebar state:', e); }
    }

    // Load state on init
    loadSidebarState();

    // Resize handle functionality
    var isResizing = false;
    if (resizeHandle) {
        resizeHandle.addEventListener('mousedown', function(e) {
            isResizing = true;
            resizeHandle.classList.add('active');
            document.body.style.cursor = 'col-resize';
            document.body.style.userSelect = 'none';
            e.preventDefault();
        });

        document.addEventListener('mousemove', function(e) {
            if (!isResizing) return;
            var width = Math.min(500, Math.max(200, e.clientX));
            document.documentElement.style.setProperty('--sidebar-width', width + 'px');
        });

        document.addEventListener('mouseup', function() {
            if (isResizing) {
                isResizing = false;
                resizeHandle.classList.remove('active');
                document.body.style.cursor = '';
                document.body.style.userSelect = '';
                saveSidebarState();
            }
        });

        // Touch support
        resizeHandle.addEventListener('touchstart', function(e) {
            isResizing = true;
            resizeHandle.classList.add('active');
        }, { passive: true });

        document.addEventListener('touchmove', function(e) {
            if (!isResizing) return;
            var touch = e.touches[0];
            var width = Math.min(500, Math.max(200, touch.clientX));
            document.documentElement.style.setProperty('--sidebar-width', width + 'px');
        }, { passive: true });

        document.addEventListener('touchend', function() {
            if (isResizing) {
                isResizing = false;
                resizeHandle.classList.remove('active');
                saveSidebarState();
            }
        });
    }

    // Mobile menu functionality
    function openSidebar() {
        sidebar.classList.add('open');
        backdrop.classList.add('active');
        menuBtn.setAttribute('aria-expanded', 'true');
    }

    function closeSidebar() {
        sidebar.classList.remove('open');
        backdrop.classList.remove('active');
        menuBtn.setAttribute('aria-expanded', 'false');
    }

    if (menuBtn) {
        menuBtn.addEventListener('click', function() {
            sidebar.classList.contains('open') ? closeSidebar() : openSidebar();
        });
    }

    if (backdrop) {
        backdrop.addEventListener('click', closeSidebar);
    }

    document.addEventListener('keydown', function(e) {
        if (e.key === 'Escape' && sidebar.classList.contains('open')) {
            closeSidebar();
        }
    });

    // Close sidebar on link click (mobile)
    sidebar.querySelectorAll('a').forEach(function(link) {
        link.addEventListener('click', function() {
            if (window.innerWidth <= 768) closeSidebar();
        });
    });

    // Tree view toggle functionality
    document.querySelectorAll('.tree-toggle').forEach(function(btn) {
        btn.addEventListener('click', function(e) {
            e.preventDefault();
            e.stopPropagation();
            var node = btn.closest('.tree-node');
            var isCollapsed = node.classList.toggle('collapsed');
            btn.setAttribute('aria-expanded', !isCollapsed);
            saveSidebarState();
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
            saveSidebarState();
        });
    }

    if (collapseAllBtn) {
        collapseAllBtn.addEventListener('click', function() {
            document.querySelectorAll('.tree-node').forEach(function(n) {
                n.classList.add('collapsed');
                var toggle = n.querySelector('.tree-toggle');
                if (toggle) toggle.setAttribute('aria-expanded', 'false');
            });
            saveSidebarState();
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
`
