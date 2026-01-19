package main

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
                <a href="../index.html" class="sidebar-title">
                    <img src="../claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .AllProjects}} projects</div>
                <a href="../stats.html" class="stats-nav-link">
                    <svg viewBox="0 0 20 20" fill="currentColor"><path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zm6-4a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zm6-3a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/></svg>
                    Stats
                </a>
            </div>
            <div class="search-container">
                <div class="search-input-wrapper">
                    <svg class="search-icon" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zZ2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
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
            <div class="session-grid">
                {{range .Project.Sessions}}
                <a href="{{.ID}}.html" class="session-card">
                    <div class="session-card-title">{{.Summary}}</div>
                    <div class="session-card-meta">
                        <span class="session-card-date">{{.CreatedAt.Format "Jan 2, 2006 3:04 PM"}}</span>
                        <span class="session-card-messages">{{len .Messages}} messages</span>
                    </div>
                </a>
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
                <a href="https://github.com/fabriqaai/claude-code-logs">claude-code-logs</a> by <a href="https://fabriqa.ai">fabriqa.ai</a>
            </footer>
        </main>
    </div>
    <script>
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
