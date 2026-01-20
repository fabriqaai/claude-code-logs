package main

// searchTemplate is the dedicated search page template
const searchTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Search - Claude Code Logs</title>
    <style>` + baseCSS + searchCSS + `</style>
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
                <a href="/" class="sidebar-title">
                    <img src="/claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .Projects}} projects</div>
                <a href="/search" class="stats-nav-link active">
                    <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    Search
                </a>
                <a href="/stats" class="stats-nav-link">
                    <svg viewBox="0 0 20 20" fill="currentColor"><path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zm6-4a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zm6-3a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/></svg>
                    Stats
                </a>
            </div>
            <div class="tree-controls">
                <button type="button" class="tree-control-btn" id="expandAll">Expand All</button>
                <button type="button" class="tree-control-btn" id="collapseAll">Collapse All</button>
            </div>
            <ul class="project-list">
                {{range $project := .Projects}}
                <li class="tree-node collapsed" data-project="{{ProjectSlug $project.Path}}">
                    <div class="tree-node-header">
                        <button type="button" class="tree-toggle{{if eq (len $project.Sessions) 0}} hidden{{end}}" aria-expanded="false" aria-label="Toggle {{$project.Path}}">
                            <svg class="tree-chevron" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
                            </svg>
                        </button>
                        <a href="/{{ProjectSlug $project.Path}}/" class="tree-node-link">
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
                            <a href="/{{ProjectSlug $project.Path}}/{{.ID}}" class="session-link">
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
        <main class="main search-main">
            <header class="page-header">
                <div class="page-breadcrumb">
                    <a href="/">Home</a>
                    <span class="page-breadcrumb-separator">/</span>
                    <span>Search</span>
                </div>
                <h1 class="page-title">Search Sessions</h1>
                <p class="page-subtitle">Find messages across all your Claude Code sessions</p>
            </header>

            <div class="search-box-container">
                <div class="search-box">
                    <svg class="search-box-icon" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    <input type="text"
                           id="searchInput"
                           class="search-box-input"
                           placeholder="Search messages, code, conversations..."
                           autocomplete="off"
                           autofocus>
                    <div class="search-box-hint">
                        <span class="search-hint-text">Use quotes for exact phrases: "hello world"</span>
                    </div>
                </div>
            </div>

            <div class="search-meta" id="searchMeta" style="display: none;">
                <span id="searchMetaText"></span>
                <div class="search-sort">
                    <label>Sort by:</label>
                    <select id="searchSort">
                        <option value="relevance">Relevance</option>
                        <option value="recent">Most Recent</option>
                    </select>
                </div>
            </div>

            <div class="search-results-container">
                <div id="searchInitial" class="search-initial">
                    <div class="search-initial-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <path d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                        </svg>
                    </div>
                    <p>Enter a search term to find messages across your sessions</p>
                </div>

                <div id="searchLoading" class="search-loading" style="display: none;">
                    <div class="loading-spinner"></div>
                    <span>Searching...</span>
                </div>

                <div id="searchEmpty" class="search-empty" style="display: none;">
                    <div class="search-empty-icon">
                        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                            <path d="M9.172 9.172a4 4 0 015.656 5.656M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                        </svg>
                    </div>
                    <p>No results found for "<span id="searchEmptyQuery"></span>"</p>
                    <span class="search-empty-hint">Try different keywords or check your spelling</span>
                </div>

                <div id="searchResults" class="search-results-list"></div>
            </div>

            <div id="loadMoreContainer" class="load-more-container" style="display: none;">
                <button type="button" id="loadMoreBtn" class="load-more-btn">
                    Load more results
                </button>
                <div id="loadMoreSpinner" class="load-more-spinner" style="display: none;">
                    <div class="loading-spinner loading-spinner-sm"></div>
                </div>
            </div>

            <footer class="footer">
                <a href="https://github.com/fabriqaai/claude-code-logs">claude-code-logs</a> by <a href="https://fabriqa.ai">fabriqa.ai</a>
            </footer>
        </main>
    </div>

    <script>
    (function() {
        // Elements
        var searchInput = document.getElementById('searchInput');
        var searchMeta = document.getElementById('searchMeta');
        var searchMetaText = document.getElementById('searchMetaText');
        var searchSort = document.getElementById('searchSort');
        var searchInitial = document.getElementById('searchInitial');
        var searchLoading = document.getElementById('searchLoading');
        var searchEmpty = document.getElementById('searchEmpty');
        var searchEmptyQuery = document.getElementById('searchEmptyQuery');
        var searchResults = document.getElementById('searchResults');
        var loadMoreContainer = document.getElementById('loadMoreContainer');
        var loadMoreBtn = document.getElementById('loadMoreBtn');
        var loadMoreSpinner = document.getElementById('loadMoreSpinner');

        // State
        var currentQuery = '';
        var currentOffset = 0;
        var currentTotal = 0;
        var hasMore = false;
        var debounceTimer = null;
        var isLoading = false;

        // Read initial query from URL
        var urlParams = new URLSearchParams(window.location.search);
        var initialQuery = urlParams.get('q') || '';
        if (initialQuery) {
            searchInput.value = initialQuery;
            performSearch(initialQuery, false);
        }

        // Debounced search on input
        searchInput.addEventListener('input', function() {
            clearTimeout(debounceTimer);
            var query = this.value.trim();

            if (query.length < 2) {
                showInitial();
                updateURL('');
                return;
            }

            debounceTimer = setTimeout(function() {
                performSearch(query, false);
            }, 300);
        });

        // Sort change
        searchSort.addEventListener('change', function() {
            if (currentQuery) {
                performSearch(currentQuery, false);
            }
        });

        // Load more
        loadMoreBtn.addEventListener('click', function() {
            if (!isLoading && hasMore) {
                performSearch(currentQuery, true);
            }
        });

        function updateURL(query) {
            var url = new URL(window.location);
            if (query) {
                url.searchParams.set('q', query);
            } else {
                url.searchParams.delete('q');
            }
            window.history.replaceState({}, '', url);
        }

        function showInitial() {
            searchInitial.style.display = 'flex';
            searchLoading.style.display = 'none';
            searchEmpty.style.display = 'none';
            searchResults.innerHTML = '';
            searchMeta.style.display = 'none';
            loadMoreContainer.style.display = 'none';
            currentQuery = '';
            currentOffset = 0;
        }

        function showLoading(append) {
            if (!append) {
                searchInitial.style.display = 'none';
                searchLoading.style.display = 'flex';
                searchEmpty.style.display = 'none';
                searchResults.innerHTML = '';
                loadMoreContainer.style.display = 'none';
            } else {
                loadMoreBtn.style.display = 'none';
                loadMoreSpinner.style.display = 'flex';
            }
        }

        function showEmpty(query) {
            searchInitial.style.display = 'none';
            searchLoading.style.display = 'none';
            searchEmpty.style.display = 'flex';
            searchEmptyQuery.textContent = query;
            searchMeta.style.display = 'none';
            loadMoreContainer.style.display = 'none';
        }

        function showResults(data, append) {
            searchInitial.style.display = 'none';
            searchLoading.style.display = 'none';
            searchEmpty.style.display = 'none';

            // Update meta
            searchMeta.style.display = 'flex';
            searchMetaText.textContent = data.total + ' session' + (data.total !== 1 ? 's' : '') + ' found';

            // Render results
            if (!append) {
                searchResults.innerHTML = '';
            }

            data.results.forEach(function(result) {
                searchResults.appendChild(createResultCard(result));
            });

            // Update load more
            hasMore = data.hasMore;
            currentOffset = data.offset + data.results.length;
            currentTotal = data.total;

            if (hasMore) {
                loadMoreContainer.style.display = 'flex';
                loadMoreBtn.style.display = 'block';
                loadMoreSpinner.style.display = 'none';
            } else {
                loadMoreContainer.style.display = 'none';
            }
        }

        function createResultCard(result) {
            var card = document.createElement('div');
            card.className = 'search-result-card';

            var matchCount = result.matches.length;
            var firstMatch = result.matches[0];

            card.innerHTML =
                '<div class="search-result-header" role="button" tabindex="0">' +
                    '<div class="search-result-info">' +
                        '<div class="search-result-title">' + escapeHtml(result.sessionTitle) + '</div>' +
                        '<div class="search-result-project">' + escapeHtml(result.project) + '</div>' +
                    '</div>' +
                    '<div class="search-result-meta">' +
                        '<span class="search-result-count">' + matchCount + ' match' + (matchCount !== 1 ? 'es' : '') + '</span>' +
                        '<span class="search-result-score" title="Relevance score">' + result.score.toFixed(1) + '</span>' +
                        '<svg class="search-result-chevron" viewBox="0 0 20 20" fill="currentColor">' +
                            '<path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"/>' +
                        '</svg>' +
                    '</div>' +
                '</div>' +
                '<div class="search-result-preview">' +
                    '<div class="search-result-excerpt">' + firstMatch.content + '</div>' +
                '</div>' +
                '<div class="search-result-expanded">' +
                    '<div class="search-result-matches">' +
                        result.matches.map(function(m, i) {
                            return '<div class="search-match-item">' +
                                '<div class="search-match-role">' + m.role + '</div>' +
                                '<div class="search-match-content">' + m.content + '</div>' +
                            '</div>';
                        }).join('') +
                    '</div>' +
                    '<a href="/' + result.projectSlug + '/' + result.sessionId + '" class="search-result-link">' +
                        'View full session →' +
                    '</a>' +
                '</div>';

            // Toggle expand on header click
            var header = card.querySelector('.search-result-header');
            header.addEventListener('click', function() {
                card.classList.toggle('expanded');
            });
            header.addEventListener('keydown', function(e) {
                if (e.key === 'Enter' || e.key === ' ') {
                    e.preventDefault();
                    card.classList.toggle('expanded');
                }
            });

            return card;
        }

        function escapeHtml(text) {
            var div = document.createElement('div');
            div.textContent = text;
            return div.innerHTML;
        }

        function performSearch(query, append) {
            if (isLoading) return;

            isLoading = true;
            currentQuery = query;
            updateURL(query);

            var offset = append ? currentOffset : 0;
            var sort = searchSort.value;

            showLoading(append);

            fetch('/api/search', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    query: query,
                    offset: offset,
                    limit: 20,
                    sort: sort
                })
            })
            .then(function(r) { return r.json(); })
            .then(function(data) {
                isLoading = false;

                if (!data.results || data.results.length === 0) {
                    if (append) {
                        // No more results
                        hasMore = false;
                        loadMoreContainer.style.display = 'none';
                    } else {
                        showEmpty(query);
                    }
                    return;
                }

                showResults(data, append);
            })
            .catch(function(err) {
                isLoading = false;
                console.error('Search error:', err);
                if (!append) {
                    showEmpty(query);
                }
            });
        }

        // Keyboard shortcut: / to focus search
        document.addEventListener('keydown', function(e) {
            if (e.key === '/' && document.activeElement !== searchInput) {
                e.preventDefault();
                searchInput.focus();
                searchInput.select();
            }
            if (e.key === 'Escape' && document.activeElement === searchInput) {
                searchInput.blur();
            }
        });
    })();
    </script>

    <!-- Sidebar functionality -->
    <script>
    (function() {
        var STORAGE_KEY = 'claude-code-logs-sidebar';
        var sidebar = document.querySelector('.sidebar');
        var resizeHandle = document.querySelector('.sidebar-resize-handle');
        var menuBtn = document.querySelector('.mobile-menu-btn');
        var backdrop = document.querySelector('.sidebar-backdrop');

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

        loadSidebarState();

        // Resize handle
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
        }

        // Mobile menu
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

        sidebar.querySelectorAll('a').forEach(function(link) {
            link.addEventListener('click', function() {
                if (window.innerWidth <= 768) closeSidebar();
            });
        });

        // Tree toggle
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
    })();
    </script>
</body>
</html>`

// searchCSS contains styles specific to the search page
const searchCSS = `
/* Search Page Specific Styles */
.search-main {
    max-width: calc(900px + 96px);
}

/* Search Box */
.search-box-container {
    margin-bottom: 32px;
}

.search-box {
    position: relative;
    max-width: 640px;
}

.search-box-icon {
    position: absolute;
    left: 16px;
    top: 50%;
    transform: translateY(-50%);
    width: 20px;
    height: 20px;
    color: var(--text-muted);
    pointer-events: none;
}

.search-box-input {
    width: 100%;
    padding: 16px 16px 16px 48px;
    font-family: var(--font-body);
    font-size: 1.1rem;
    border: 2px solid var(--border-medium);
    border-radius: 12px;
    background: var(--bg-secondary);
    color: var(--text-primary);
    transition: all var(--transition-fast);
}

.search-box-input:focus {
    outline: none;
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 3px var(--accent-subtle);
}

.search-box-input::placeholder {
    color: var(--text-muted);
}

.search-box-hint {
    margin-top: 8px;
    padding-left: 4px;
}

.search-hint-text {
    font-size: 0.8rem;
    color: var(--text-muted);
}

/* Search Meta */
.search-meta {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 12px 0;
    margin-bottom: 16px;
    border-bottom: 1px solid var(--border-subtle);
    font-size: 0.9rem;
    color: var(--text-secondary);
}

.search-sort {
    display: flex;
    align-items: center;
    gap: 8px;
}

.search-sort label {
    color: var(--text-muted);
}

.search-sort select {
    padding: 6px 12px;
    border: 1px solid var(--border-medium);
    border-radius: 6px;
    background: var(--bg-secondary);
    color: var(--text-primary);
    font-family: var(--font-body);
    font-size: 0.85rem;
    cursor: pointer;
}

/* Search States */
.search-initial,
.search-loading,
.search-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 60px 20px;
    color: var(--text-muted);
    text-align: center;
}

.search-initial-icon,
.search-empty-icon {
    width: 64px;
    height: 64px;
    margin-bottom: 16px;
    color: var(--border-medium);
}

.search-initial-icon svg,
.search-empty-icon svg {
    width: 100%;
    height: 100%;
}

.search-empty-hint {
    margin-top: 8px;
    font-size: 0.85rem;
}

.loading-spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--border-subtle);
    border-top-color: var(--accent-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
    margin-bottom: 12px;
}

.loading-spinner-sm {
    width: 20px;
    height: 20px;
    border-width: 2px;
    margin-bottom: 0;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

/* Search Results */
.search-results-list {
    display: flex;
    flex-direction: column;
    gap: 12px;
}

.search-result-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    overflow: hidden;
    transition: all var(--transition-fast);
}

.search-result-card:hover {
    border-color: var(--border-medium);
    box-shadow: var(--shadow-sm);
}

.search-result-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 16px;
    cursor: pointer;
    user-select: none;
}

.search-result-header:focus {
    outline: none;
    background: var(--bg-tertiary);
}

.search-result-info {
    flex: 1;
    min-width: 0;
}

.search-result-title {
    font-family: var(--font-display);
    font-size: 1rem;
    font-weight: 500;
    color: var(--text-primary);
    margin-bottom: 4px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.search-result-project {
    font-size: 0.8rem;
    color: var(--text-muted);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.search-result-meta {
    display: flex;
    align-items: center;
    gap: 12px;
    flex-shrink: 0;
    margin-left: 16px;
}

.search-result-count {
    font-size: 0.8rem;
    color: var(--text-muted);
    background: var(--bg-tertiary);
    padding: 4px 8px;
    border-radius: 4px;
}

.search-result-score {
    font-size: 0.75rem;
    color: var(--accent-primary);
    font-weight: 500;
}

.search-result-chevron {
    width: 20px;
    height: 20px;
    color: var(--text-muted);
    transition: transform var(--transition-fast);
}

.search-result-card.expanded .search-result-chevron {
    transform: rotate(180deg);
}

/* Preview (always visible) */
.search-result-preview {
    padding: 0 16px 16px;
    border-top: 1px solid var(--border-subtle);
}

.search-result-excerpt {
    font-size: 0.9rem;
    color: var(--text-secondary);
    line-height: 1.6;
    padding-top: 12px;
}

.search-result-excerpt mark {
    background: var(--accent-subtle);
    color: var(--accent-primary);
    padding: 1px 3px;
    border-radius: 3px;
}

/* Expanded State */
.search-result-expanded {
    display: none;
    padding: 16px;
    background: var(--bg-tertiary);
    border-top: 1px solid var(--border-subtle);
}

.search-result-card.expanded .search-result-expanded {
    display: block;
}

.search-result-card.expanded .search-result-preview {
    display: none;
}

.search-result-matches {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 16px;
    max-height: 400px;
    overflow-y: auto;
}

.search-match-item {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    padding: 12px;
}

.search-match-role {
    font-size: 0.75rem;
    font-weight: 500;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
    margin-bottom: 6px;
}

.search-match-content {
    font-size: 0.9rem;
    color: var(--text-secondary);
    line-height: 1.6;
}

.search-match-content mark {
    background: var(--accent-subtle);
    color: var(--accent-primary);
    padding: 1px 3px;
    border-radius: 3px;
}

.search-result-link {
    display: inline-flex;
    align-items: center;
    gap: 4px;
    font-size: 0.9rem;
    color: var(--accent-primary);
    text-decoration: none;
    font-weight: 500;
    transition: color var(--transition-fast);
}

.search-result-link:hover {
    color: var(--accent-hover);
    text-decoration: underline;
}

/* Load More */
.load-more-container {
    display: flex;
    justify-content: center;
    padding: 24px 0;
}

.load-more-btn {
    padding: 12px 24px;
    font-family: var(--font-body);
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--accent-primary);
    background: var(--bg-secondary);
    border: 2px solid var(--accent-primary);
    border-radius: 8px;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.load-more-btn:hover {
    background: var(--accent-primary);
    color: white;
}

.load-more-spinner {
    display: flex;
    align-items: center;
    justify-content: center;
}

/* Responsive */
@media (max-width: 768px) {
    .search-box-input {
        font-size: 1rem;
        padding: 14px 14px 14px 44px;
    }

    .search-meta {
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
    }

    .search-result-header {
        flex-direction: column;
        align-items: flex-start;
        gap: 12px;
    }

    .search-result-meta {
        margin-left: 0;
        width: 100%;
        justify-content: space-between;
    }
}
`
