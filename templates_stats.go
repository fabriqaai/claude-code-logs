package main

// statsTemplate is the stats page template showing analytics charts
const statsTemplate = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Stats - Claude Code Logs</title>
    <style>` + baseCSS + statsCSS + `</style>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
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
                <a href="index.html" class="sidebar-title">
                    <img src="claude-code-icon.png" alt="Claude Code" class="sidebar-logo">
                    Claude Code Logs
                </a>
                <div class="sidebar-subtitle">{{len .Projects}} projects</div>
                <a href="stats.html" class="stats-nav-link active">
                    <svg viewBox="0 0 20 20" fill="currentColor"><path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zm6-4a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zm6-3a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z"/></svg>
                    Stats
                </a>
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
                <li class="tree-node collapsed" data-project="{{ProjectSlug $project.Path}}">
                    <div class="tree-node-header">
                        <button type="button" class="tree-toggle{{if eq (len $project.Sessions) 0}} hidden{{end}}" aria-expanded="false" aria-label="Toggle {{$project.Path}}">
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
        <main class="main stats-main">
            <header class="page-header">
                <div class="page-breadcrumb">
                    <a href="index.html">Home</a>
                    <span class="page-breadcrumb-separator">/</span>
                    <span>Stats</span>
                </div>
                <h1 class="page-title">Usage Statistics</h1>
                <p class="page-subtitle">Analytics for your Claude Code sessions</p>
            </header>

            <div id="stats-loading" class="stats-loading">
                <div class="loading-spinner"></div>
                <span>Loading statistics...</span>
            </div>

            <div id="stats-content" class="stats-content" style="display: none;">
                <!-- Time Range Filter -->
                <div class="time-filter">
                    <button type="button" class="time-filter-btn" data-range="today">Today</button>
                    <button type="button" class="time-filter-btn" data-range="week">This Week</button>
                    <button type="button" class="time-filter-btn active" data-range="month">This Month</button>
                    <button type="button" class="time-filter-btn" data-range="all">All Time</button>
                </div>

                <!-- Summary Cards -->
                <div class="stats-summary">
                    <div class="stat-card">
                        <div class="stat-value" id="stat-messages">-</div>
                        <div class="stat-label">Total Messages</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value" id="stat-sessions">-</div>
                        <div class="stat-label">Total Sessions</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value" id="stat-projects">-</div>
                        <div class="stat-label">Projects</div>
                    </div>
                    <div class="stat-card">
                        <div class="stat-value" id="stat-tokens">-</div>
                        <div class="stat-label">Est. Tokens</div>
                    </div>
                </div>

                <!-- Charts Grid -->
                <div class="charts-grid">
                    <div class="chart-card">
                        <div class="chart-header">
                            <h3 class="chart-title">Messages per Day</h3>
                            <span class="chart-subtitle">Last 30 days</span>
                        </div>
                        <div class="chart-container">
                            <canvas id="messagesChart"></canvas>
                        </div>
                    </div>

                    <div class="chart-card">
                        <div class="chart-header">
                            <h3 class="chart-title">Tokens per Day</h3>
                            <span class="chart-subtitle">Estimated usage</span>
                        </div>
                        <div class="chart-container">
                            <canvas id="tokensChart"></canvas>
                        </div>
                    </div>

                    <div class="chart-card chart-wide">
                        <div class="chart-header">
                            <h3 class="chart-title">Project Activity</h3>
                            <span class="chart-subtitle">Messages by project</span>
                        </div>
                        <div class="chart-container chart-container-tall">
                            <canvas id="projectsChart"></canvas>
                        </div>
                    </div>
                </div>

                <!-- Session Stats -->
                <div class="session-stats">
                    <div class="stat-card stat-card-wide">
                        <div class="stat-inline">
                            <div>
                                <div class="stat-value" id="stat-avg-length">-</div>
                                <div class="stat-label">Avg Session Length</div>
                            </div>
                            <div>
                                <div class="stat-value" id="stat-avg-messages">-</div>
                                <div class="stat-label">Avg Messages/Session</div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>

            <footer class="footer">
                <a href="https://github.com/fabriqaai/claude-code-logs">claude-code-logs</a> by <a href="https://fabriqa.ai">fabriqa.ai</a>
            </footer>
        </main>
    </div>

    <script>
    (function() {
        var loading = document.getElementById('stats-loading');
        var content = document.getElementById('stats-content');

        // Chart.js default options
        Chart.defaults.font.family = "'DM Sans', -apple-system, sans-serif";
        Chart.defaults.color = '#57534E';

        // Color palette
        var accentColor = '#C2410C';
        var accentColorLight = 'rgba(194, 65, 12, 0.1)';
        var greenColor = '#059669';

        function formatNumber(num) {
            if (num >= 1000000) return (num / 1000000).toFixed(1) + 'M';
            if (num >= 1000) return (num / 1000).toFixed(1) + 'K';
            return num.toString();
        }

        function formatMinutes(mins) {
            if (mins < 60) return Math.round(mins) + ' min';
            var hours = Math.floor(mins / 60);
            var remainMins = Math.round(mins % 60);
            return hours + 'h ' + remainMins + 'm';
        }

        // Store full data for filtering
        var fullData = null;
        var messagesChart = null;
        var tokensChart = null;
        var projectsChart = null;
        var currentRange = 'month';

        // Time range filter setup
        var filterBtns = document.querySelectorAll('.time-filter-btn');
        filterBtns.forEach(function(btn) {
            btn.addEventListener('click', function() {
                var range = btn.dataset.range;
                if (range === currentRange) return;

                // Update active state
                filterBtns.forEach(function(b) { b.classList.remove('active'); });
                btn.classList.add('active');
                currentRange = range;

                // Re-render with filtered data
                if (fullData) {
                    var filtered = filterDataByRange(fullData, range);
                    updateDisplay(filtered);
                }
            });
        });

        function filterDataByRange(data, range) {
            if (range === 'all') return data;

            var now = new Date();
            var cutoff;
            switch (range) {
                case 'today':
                    cutoff = new Date(now.getFullYear(), now.getMonth(), now.getDate());
                    break;
                case 'week':
                    cutoff = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
                    break;
                case 'month':
                default:
                    cutoff = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
            }
            var cutoffStr = cutoff.toISOString().split('T')[0];

            // Filter time series
            var filteredMessages = data.messagesPerDay.filter(function(d) { return d.date >= cutoffStr; });
            var filteredTokens = data.tokensPerDay.filter(function(d) { return d.date >= cutoffStr; });

            // Recalculate totals from filtered data
            var totalMessages = filteredMessages.reduce(function(sum, d) { return sum + d.value; }, 0);
            var totalTokens = filteredTokens.reduce(function(sum, d) { return sum + d.value; }, 0);

            return {
                totalProjects: data.totalProjects,
                totalSessions: data.totalSessions,
                totalMessages: totalMessages,
                totalTokens: totalTokens,
                messagesPerDay: filteredMessages,
                tokensPerDay: filteredTokens,
                projectStats: data.projectStats,
                avgSessionLengthMins: data.avgSessionLengthMins,
                avgMessagesPerSession: data.avgMessagesPerSession
            };
        }

        function updateDisplay(data) {
            // Update summary cards
            document.getElementById('stat-messages').textContent = formatNumber(data.totalMessages);
            document.getElementById('stat-tokens').textContent = formatNumber(data.totalTokens);

            // Update charts
            updateMessagesChart(data.messagesPerDay);
            updateTokensChart(data.tokensPerDay);
        }

        function updateMessagesChart(data) {
            if (messagesChart) {
                messagesChart.data.labels = data.map(function(d) {
                    return new Date(d.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                });
                messagesChart.data.datasets[0].data = data.map(function(d) { return d.value; });
                messagesChart.update();
            }
        }

        function updateTokensChart(data) {
            if (tokensChart) {
                tokensChart.data.labels = data.map(function(d) {
                    return new Date(d.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                });
                tokensChart.data.datasets[0].data = data.map(function(d) { return d.value; });
                tokensChart.update();
            }
        }

        // Fetch and render stats
        fetch('/api/stats')
            .then(function(r) { return r.json(); })
            .then(function(data) {
                fullData = data;

                // Apply default filter (month)
                var filtered = filterDataByRange(data, currentRange);

                // Update summary cards
                document.getElementById('stat-messages').textContent = formatNumber(filtered.totalMessages);
                document.getElementById('stat-sessions').textContent = formatNumber(data.totalSessions);
                document.getElementById('stat-projects').textContent = formatNumber(data.totalProjects);
                document.getElementById('stat-tokens').textContent = formatNumber(filtered.totalTokens);
                document.getElementById('stat-avg-length').textContent = formatMinutes(data.avgSessionLengthMins);
                document.getElementById('stat-avg-messages').textContent = data.avgMessagesPerSession.toFixed(1);

                // Render charts with filtered data
                messagesChart = renderMessagesChart(filtered.messagesPerDay);
                tokensChart = renderTokensChart(filtered.tokensPerDay);
                projectsChart = renderProjectsChart(data.projectStats);

                // Show content
                loading.style.display = 'none';
                content.style.display = 'block';
            })
            .catch(function(err) {
                console.error('Failed to load stats:', err);
                loading.innerHTML = '<div class="error-state"><div class="error-icon">!</div><p>Failed to load statistics</p></div>';
            });

        function renderMessagesChart(data) {
            var ctx = document.getElementById('messagesChart').getContext('2d');
            return new Chart(ctx, {
                type: 'line',
                data: {
                    labels: data.map(function(d) {
                        return new Date(d.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                    }),
                    datasets: [{
                        label: 'Messages',
                        data: data.map(function(d) { return d.value; }),
                        borderColor: accentColor,
                        backgroundColor: accentColorLight,
                        fill: true,
                        tension: 0.3,
                        pointRadius: 2,
                        pointHoverRadius: 5
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false }
                    },
                    scales: {
                        x: {
                            grid: { display: false },
                            ticks: { maxTicksLimit: 7 }
                        },
                        y: {
                            beginAtZero: true,
                            grid: { color: 'rgba(0,0,0,0.05)' }
                        }
                    }
                }
            });
        }

        function renderTokensChart(data) {
            var ctx = document.getElementById('tokensChart').getContext('2d');
            return new Chart(ctx, {
                type: 'line',
                data: {
                    labels: data.map(function(d) {
                        return new Date(d.date).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
                    }),
                    datasets: [{
                        label: 'Tokens',
                        data: data.map(function(d) { return d.value; }),
                        borderColor: greenColor,
                        backgroundColor: 'rgba(5, 150, 105, 0.1)',
                        fill: true,
                        tension: 0.3,
                        pointRadius: 2,
                        pointHoverRadius: 5
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    plugins: {
                        legend: { display: false }
                    },
                    scales: {
                        x: {
                            grid: { display: false },
                            ticks: { maxTicksLimit: 7 }
                        },
                        y: {
                            beginAtZero: true,
                            grid: { color: 'rgba(0,0,0,0.05)' },
                            ticks: {
                                callback: function(value) {
                                    return formatNumber(value);
                                }
                            }
                        }
                    }
                }
            });
        }

        function renderProjectsChart(data) {
            // Sort by messages and take top 10
            var sorted = data.slice().sort(function(a, b) { return b.messages - a.messages; }).slice(0, 10);

            var ctx = document.getElementById('projectsChart').getContext('2d');
            return new Chart(ctx, {
                type: 'bar',
                data: {
                    labels: sorted.map(function(d) {
                        var name = d.path.split('/').pop() || d.slug;
                        return name.length > 25 ? name.substring(0, 22) + '...' : name;
                    }),
                    datasets: [{
                        label: 'Messages',
                        data: sorted.map(function(d) { return d.messages; }),
                        backgroundColor: accentColor,
                        borderRadius: 4
                    }]
                },
                options: {
                    responsive: true,
                    maintainAspectRatio: false,
                    indexAxis: 'y',
                    plugins: {
                        legend: { display: false }
                    },
                    scales: {
                        x: {
                            beginAtZero: true,
                            grid: { color: 'rgba(0,0,0,0.05)' }
                        },
                        y: {
                            grid: { display: false }
                        }
                    }
                }
            });
        }
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

        // Search
        var searchInput = document.getElementById('searchInput');
        var searchResults = document.getElementById('searchResults');
        var debounceTimer;

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
                            html += '<a href="' + r.projectSlug + '/' + r.sessionId + '.html" class="search-result-item">';
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

// statsCSS contains styles specific to the stats page
const statsCSS = `
/* Stats Page Specific Styles */
.stats-main {
    max-width: calc(1200px + 96px);
}

.stats-loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 80px;
    color: var(--text-muted);
    gap: 16px;
}

.loading-spinner {
    width: 32px;
    height: 32px;
    border: 3px solid var(--border-subtle);
    border-top-color: var(--accent-primary);
    border-radius: 50%;
    animation: spin 1s linear infinite;
}

@keyframes spin {
    to { transform: rotate(360deg); }
}

.stats-content {
    animation: fadeIn 0.4s ease-out;
}

/* Time Range Filter */
.time-filter {
    display: flex;
    gap: 8px;
    margin-bottom: 24px;
    flex-wrap: wrap;
}

.time-filter-btn {
    padding: 8px 16px;
    border: 1px solid var(--border-subtle);
    border-radius: 20px;
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.time-filter-btn:hover {
    border-color: var(--border-medium);
    background: var(--bg-tertiary);
}

.time-filter-btn.active {
    background: var(--accent-primary);
    border-color: var(--accent-primary);
    color: white;
}

.time-filter-btn.active:hover {
    background: var(--accent-hover);
    border-color: var(--accent-hover);
}

/* Summary Cards */
.stats-summary {
    display: grid;
    grid-template-columns: repeat(4, 1fr);
    gap: 16px;
    margin-bottom: 24px;
}

.stat-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    text-align: center;
    transition: all var(--transition-fast);
    box-shadow: var(--shadow-sm);
}

.stat-card:hover {
    box-shadow: var(--shadow-md);
    transform: translateY(-2px);
}

.stat-value {
    font-family: var(--font-display);
    font-size: 2rem;
    font-weight: 600;
    color: var(--accent-primary);
    margin-bottom: 4px;
}

.stat-label {
    font-size: 0.85rem;
    color: var(--text-muted);
    text-transform: uppercase;
    letter-spacing: 0.05em;
}

.stat-card-wide {
    grid-column: span 2;
}

.stat-inline {
    display: flex;
    justify-content: space-around;
    gap: 32px;
}

/* Charts Grid */
.charts-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 20px;
    margin-bottom: 24px;
}

.chart-card {
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    padding: 20px;
    box-shadow: var(--shadow-sm);
}

.chart-card.chart-wide {
    grid-column: span 2;
}

.chart-header {
    display: flex;
    align-items: baseline;
    justify-content: space-between;
    margin-bottom: 16px;
}

.chart-title {
    font-family: var(--font-display);
    font-size: 1rem;
    font-weight: 500;
    color: var(--text-primary);
    margin: 0;
}

.chart-subtitle {
    font-size: 0.8rem;
    color: var(--text-muted);
}

.chart-container {
    position: relative;
    height: 200px;
}

.chart-container-tall {
    height: 300px;
}

/* Session Stats */
.session-stats {
    display: grid;
    grid-template-columns: 1fr;
    gap: 16px;
}

/* Error State */
.error-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 12px;
}

.error-icon {
    width: 48px;
    height: 48px;
    background: var(--accent-subtle);
    color: var(--accent-primary);
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    font-weight: 600;
}

/* Responsive */
@media (max-width: 1024px) {
    .stats-summary {
        grid-template-columns: repeat(2, 1fr);
    }
    
    .stat-card-wide {
        grid-column: span 2;
    }
}

@media (max-width: 768px) {
    .stats-summary {
        grid-template-columns: 1fr 1fr;
    }
    
    .charts-grid {
        grid-template-columns: 1fr;
    }
    
    .chart-card.chart-wide {
        grid-column: span 1;
    }
    
    .stat-card-wide {
        grid-column: span 2;
    }
    
    .stat-inline {
        flex-direction: column;
        gap: 16px;
    }
}

@media (max-width: 480px) {
    .stats-summary {
        grid-template-columns: 1fr;
    }
    
    .stat-card-wide {
        grid-column: span 1;
    }
    
    .stat-value {
        font-size: 1.5rem;
    }
}
`
