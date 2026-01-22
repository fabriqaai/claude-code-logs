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
                <a href="/search" class="stats-nav-link">
                    <svg viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd"/>
                    </svg>
                    Search
                </a>
                <a href="/stats" class="stats-nav-link active">
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
                <!-- Filters -->
                <div class="stats-filters">
                    <!-- Time Range Filter -->
                    <div class="time-filter">
                        <button type="button" class="time-filter-btn" data-range="today">Today</button>
                        <button type="button" class="time-filter-btn" data-range="week">This Week</button>
                        <button type="button" class="time-filter-btn active" data-range="month">This Month</button>
                        <button type="button" class="time-filter-btn" data-range="all">All Time</button>
                    </div>

                    <div class="filter-row">
                        <!-- Account Filter -->
                        <div class="account-filter-container">
                            <label class="filter-label">Account:</label>
                            <select id="account-filter" class="account-filter">
                                <option value="">All Accounts</option>
                            </select>
                            <button type="button" class="account-settings-btn" id="account-settings-btn" title="Configure Accounts">
                                <svg viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clip-rule="evenodd"/>
                                </svg>
                            </button>
                        </div>

                        <!-- Project Multi-Select Filter -->
                        <div class="project-filter-container">
                            <label class="filter-label">Projects:</label>
                            <div class="multi-select-wrapper">
                                <div class="multi-select-display" id="project-filter-display">
                                    <span class="multi-select-placeholder">All Projects</span>
                                    <svg class="multi-select-chevron" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd"/>
                                    </svg>
                                </div>
                                <div class="multi-select-dropdown" id="project-filter-dropdown">
                                    <div class="multi-select-search-wrapper">
                                        <input type="text" class="multi-select-search" id="project-filter-search" placeholder="Search projects...">
                                    </div>
                                    <div class="multi-select-options" id="project-filter-options"></div>
                                    <div class="multi-select-actions">
                                        <button type="button" class="multi-select-action" id="project-filter-clear">Clear All</button>
                                        <button type="button" class="multi-select-action primary" id="project-filter-done">Done</button>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>

                <!-- Selected Projects Chips -->
                <div class="selected-chips" id="selected-projects-chips"></div>

                <!-- Account Settings Modal -->
                <div class="modal-backdrop" id="account-modal-backdrop">
                    <div class="modal" id="account-modal">
                        <div class="modal-header">
                            <h2 class="modal-title">Account Settings</h2>
                            <button type="button" class="modal-close" id="account-modal-close">&times;</button>
                        </div>
                        <div class="modal-body">
                            <p class="modal-description">Map path patterns to accounts. Projects matching a pattern will be assigned to that account.</p>
                            <div class="account-mappings" id="account-mappings">
                                <!-- Dynamic account mapping rows -->
                            </div>
                            <button type="button" class="add-mapping-btn" id="add-mapping-btn">
                                <svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M10 3a1 1 0 011 1v5h5a1 1 0 110 2h-5v5a1 1 0 11-2 0v-5H4a1 1 0 110-2h5V4a1 1 0 011-1z" clip-rule="evenodd"/></svg>
                                Add Account Mapping
                            </button>
                        </div>
                        <div class="modal-footer">
                            <button type="button" class="modal-btn" id="account-modal-cancel">Cancel</button>
                            <button type="button" class="modal-btn primary" id="account-modal-save">Save</button>
                        </div>
                    </div>
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
                    <div class="stat-card">
                        <div class="stat-value" id="stat-cost">-</div>
                        <div class="stat-label">Est. Cost</div>
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

        function formatCurrency(amount) {
            if (amount >= 1000) return '$' + (amount / 1000).toFixed(1) + 'K';
            if (amount >= 100) return '$' + amount.toFixed(0);
            return '$' + amount.toFixed(2);
        }

        // Store full data for filtering
        var fullData = null;
        var messagesChart = null;
        var tokensChart = null;
        var projectsChart = null;
        var currentRange = 'month';
        var selectedProjects = []; // Array of selected project slugs (empty = all)

        // Multi-select elements
        var filterDisplay = document.getElementById('project-filter-display');
        var filterDropdown = document.getElementById('project-filter-dropdown');
        var filterOptions = document.getElementById('project-filter-options');
        var filterSearch = document.getElementById('project-filter-search');
        var filterClear = document.getElementById('project-filter-clear');
        var filterDone = document.getElementById('project-filter-done');
        var chipsContainer = document.getElementById('selected-projects-chips');

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
                    var filtered = filterData(fullData, currentRange, selectedProjects);
                    updateDisplay(filtered);
                }
            });
        });

        // Multi-select dropdown toggle
        filterDisplay.addEventListener('click', function(e) {
            e.stopPropagation();
            var isOpen = filterDropdown.classList.contains('open');
            if (isOpen) {
                closeDropdown();
            } else {
                openDropdown();
            }
        });

        function openDropdown() {
            filterDropdown.classList.add('open');
            filterDisplay.classList.add('active');
            filterSearch.focus();
        }

        function closeDropdown() {
            filterDropdown.classList.remove('open');
            filterDisplay.classList.remove('active');
            filterSearch.value = '';
            filterOptionsFromSearch('');
        }

        // Close dropdown when clicking outside
        document.addEventListener('click', function(e) {
            if (!filterDropdown.contains(e.target) && !filterDisplay.contains(e.target)) {
                closeDropdown();
            }
        });

        // Search filtering
        filterSearch.addEventListener('input', function() {
            filterOptionsFromSearch(filterSearch.value.toLowerCase());
        });

        function filterOptionsFromSearch(query) {
            var options = filterOptions.querySelectorAll('.multi-select-option');
            options.forEach(function(opt) {
                var label = opt.querySelector('.multi-select-option-label').textContent.toLowerCase();
                opt.style.display = label.includes(query) ? '' : 'none';
            });
        }

        // Clear all selections
        filterClear.addEventListener('click', function() {
            selectedProjects = [];
            updateProjectSelection();
            if (fullData) {
                var filtered = filterData(fullData, currentRange, selectedProjects);
                updateDisplay(filtered);
            }
        });

        // Done button
        filterDone.addEventListener('click', function() {
            closeDropdown();
        });

        function populateProjectFilter(projectStats) {
            filterOptions.innerHTML = '';
            // Sort projects by messages (most active first)
            var sorted = projectStats.slice().sort(function(a, b) { return b.messages - a.messages; });
            sorted.forEach(function(project) {
                var name = project.path.split('/').pop() || project.slug;
                var displayName = name.length > 35 ? name.substring(0, 32) + '...' : name;

                var option = document.createElement('div');
                option.className = 'multi-select-option';
                option.dataset.slug = project.slug;
                option.innerHTML = '<div class="multi-select-checkbox"><svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/></svg></div>' +
                    '<span class="multi-select-option-label">' + displayName + '</span>' +
                    '<span class="multi-select-option-meta">' + formatNumber(project.messages) + ' msgs</span>';

                option.addEventListener('click', function() {
                    toggleProjectSelection(project.slug);
                });

                filterOptions.appendChild(option);
            });
        }

        function toggleProjectSelection(slug) {
            var idx = selectedProjects.indexOf(slug);
            if (idx === -1) {
                selectedProjects.push(slug);
            } else {
                selectedProjects.splice(idx, 1);
            }
            updateProjectSelection();
            if (fullData) {
                var filtered = filterData(fullData, currentRange, selectedProjects);
                updateDisplay(filtered);
            }
        }

        function updateProjectSelection() {
            // Update dropdown options
            var options = filterOptions.querySelectorAll('.multi-select-option');
            options.forEach(function(opt) {
                var slug = opt.dataset.slug;
                if (selectedProjects.indexOf(slug) !== -1) {
                    opt.classList.add('selected');
                } else {
                    opt.classList.remove('selected');
                }
            });

            // Update display text
            var placeholder = filterDisplay.querySelector('.multi-select-placeholder');
            var countBadge = filterDisplay.querySelector('.multi-select-count');

            if (selectedProjects.length === 0) {
                placeholder.textContent = 'All Projects';
                if (countBadge) countBadge.remove();
            } else {
                placeholder.textContent = selectedProjects.length + ' selected';
                if (!countBadge) {
                    countBadge = document.createElement('span');
                    countBadge.className = 'multi-select-count';
                    filterDisplay.insertBefore(countBadge, filterDisplay.querySelector('.multi-select-chevron'));
                }
                countBadge.textContent = selectedProjects.length;
            }

            // Update chips
            renderChips();
        }

        function renderChips() {
            chipsContainer.innerHTML = '';
            if (!fullData) return;

            selectedProjects.forEach(function(slug) {
                var project = fullData.projectStats.find(function(p) { return p.slug === slug; });
                if (!project) return;

                var name = project.path.split('/').pop() || slug;
                var displayName = name.length > 20 ? name.substring(0, 17) + '...' : name;

                var chip = document.createElement('div');
                chip.className = 'selected-chip';
                chip.innerHTML = '<span>' + displayName + '</span>' +
                    '<button type="button" class="selected-chip-remove" data-slug="' + slug + '">' +
                    '<svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/></svg>' +
                    '</button>';

                chip.querySelector('.selected-chip-remove').addEventListener('click', function() {
                    toggleProjectSelection(slug);
                });

                chipsContainer.appendChild(chip);
            });
        }

        function getSelectedProjectsData() {
            if (!fullData || selectedProjects.length === 0) return null;
            return fullData.projectStats.filter(function(p) {
                return selectedProjects.indexOf(p.slug) !== -1;
            });
        }

        function aggregateProjectsData(projects) {
            if (!projects || projects.length === 0) return null;

            // Aggregate time series by day
            var messagesByDay = {};
            var tokensByDay = {};
            var costByDay = {};
            var totalMessages = 0;
            var totalTokens = 0;
            var totalCost = 0;
            var totalSessions = 0;
            var totalSessionMins = 0;

            projects.forEach(function(p) {
                totalMessages += p.messages;
                totalTokens += p.tokens;
                totalCost += p.cost || 0;
                totalSessions += p.sessions;
                totalSessionMins += (p.avgSessionLengthMins || 0) * p.sessions;

                (p.messagesPerDay || []).forEach(function(d) {
                    messagesByDay[d.date] = (messagesByDay[d.date] || 0) + d.value;
                });
                (p.tokensPerDay || []).forEach(function(d) {
                    tokensByDay[d.date] = (tokensByDay[d.date] || 0) + d.value;
                    costByDay[d.date] = (costByDay[d.date] || 0) + (d.cost || 0);
                });
            });

            // Convert to sorted arrays
            var dates = Object.keys(messagesByDay).sort();
            var messagesPerDay = dates.map(function(d) { return { date: d, value: messagesByDay[d] }; });
            var tokensPerDay = dates.map(function(d) { return { date: d, value: tokensByDay[d] || 0, cost: costByDay[d] || 0 }; });

            return {
                messages: totalMessages,
                tokens: totalTokens,
                cost: totalCost,
                sessions: totalSessions,
                messagesPerDay: messagesPerDay,
                tokensPerDay: tokensPerDay,
                avgSessionLengthMins: totalSessions > 0 ? totalSessionMins / totalSessions : 0,
                avgMessagesPerSession: totalSessions > 0 ? totalMessages / totalSessions : 0
            };
        }

        function filterData(data, range, projects) {
            // Get base data - either from selected projects or all
            var selectedData = getSelectedProjectsData();
            var aggregated = selectedData ? aggregateProjectsData(selectedData) : null;

            var baseMessagesPerDay = aggregated ? aggregated.messagesPerDay : data.messagesPerDay;
            var baseTokensPerDay = aggregated ? aggregated.tokensPerDay : data.tokensPerDay;
            var baseTotalMessages = aggregated ? aggregated.messages : data.totalMessages;
            var baseTotalTokens = aggregated ? aggregated.tokens : data.totalTokens;
            var baseTotalCost = aggregated ? aggregated.cost : (data.totalCost || 0);
            var baseSessions = aggregated ? aggregated.sessions : data.totalSessions;
            var baseAvgLength = aggregated ? aggregated.avgSessionLengthMins : data.avgSessionLengthMins;
            var baseAvgMessages = aggregated ? aggregated.avgMessagesPerSession : data.avgMessagesPerSession;

            // If no time filter, return project-filtered data
            if (range === 'all') {
                return {
                    totalProjects: projects.length > 0 ? projects.length : data.totalProjects,
                    totalSessions: baseSessions,
                    totalMessages: baseTotalMessages,
                    totalTokens: baseTotalTokens,
                    totalCost: baseTotalCost,
                    messagesPerDay: baseMessagesPerDay || [],
                    tokensPerDay: baseTokensPerDay || [],
                    projectStats: data.projectStats,
                    avgSessionLengthMins: baseAvgLength,
                    avgMessagesPerSession: baseAvgMessages,
                    selectedProjects: projects
                };
            }

            // Apply time range filter
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
            var filteredMessages = (baseMessagesPerDay || []).filter(function(d) { return d.date >= cutoffStr; });
            var filteredTokens = (baseTokensPerDay || []).filter(function(d) { return d.date >= cutoffStr; });

            // Recalculate totals from filtered data
            var totalMessages = filteredMessages.reduce(function(sum, d) { return sum + d.value; }, 0);
            var totalTokens = filteredTokens.reduce(function(sum, d) { return sum + d.value; }, 0);
            var totalCost = filteredTokens.reduce(function(sum, d) { return sum + (d.cost || 0); }, 0);

            return {
                totalProjects: projects.length > 0 ? projects.length : data.totalProjects,
                totalSessions: baseSessions,
                totalMessages: totalMessages,
                totalTokens: totalTokens,
                totalCost: totalCost,
                messagesPerDay: filteredMessages,
                tokensPerDay: filteredTokens,
                projectStats: data.projectStats,
                avgSessionLengthMins: baseAvgLength,
                avgMessagesPerSession: baseAvgMessages,
                selectedProjects: projects
            };
        }

        function updateDisplay(data) {
            // Update summary cards
            document.getElementById('stat-messages').textContent = formatNumber(data.totalMessages);
            document.getElementById('stat-tokens').textContent = formatNumber(data.totalTokens);
            document.getElementById('stat-sessions').textContent = formatNumber(data.totalSessions);
            document.getElementById('stat-projects').textContent = formatNumber(data.totalProjects);
            document.getElementById('stat-cost').textContent = formatCurrency(data.totalCost || 0);

            // Update bottom stats (now filtered!)
            document.getElementById('stat-avg-length').textContent = formatMinutes(data.avgSessionLengthMins);
            document.getElementById('stat-avg-messages').textContent = data.avgMessagesPerSession.toFixed(1);

            // Update charts
            updateMessagesChart(data.messagesPerDay);
            updateTokensChart(data.tokensPerDay);
            updateProjectsChart(data.projectStats, data.selectedProjects);
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

        function updateProjectsChart(projectStats, selectedProjectSlugs) {
            if (!projectsChart) return;

            var displayData;
            if (selectedProjectSlugs && selectedProjectSlugs.length > 0) {
                // Show only selected projects, sorted by messages
                displayData = projectStats
                    .filter(function(p) { return selectedProjectSlugs.indexOf(p.slug) !== -1; })
                    .sort(function(a, b) { return b.messages - a.messages; });
            } else {
                // Show top 10 by messages
                displayData = projectStats.slice().sort(function(a, b) { return b.messages - a.messages; }).slice(0, 10);
            }

            // All selected projects get accent color
            var colors = displayData.map(function() { return accentColor; });

            projectsChart.data.labels = displayData.map(function(d) {
                var name = d.path.split('/').pop() || d.slug;
                return name.length > 25 ? name.substring(0, 22) + '...' : name;
            });
            projectsChart.data.datasets[0].data = displayData.map(function(d) { return d.messages; });
            projectsChart.data.datasets[0].backgroundColor = colors;
            projectsChart.update();
        }

        // Fetch and render stats
        fetch('/api/stats')
            .then(function(r) { return r.json(); })
            .then(function(data) {
                fullData = data;

                // Populate project filter dropdown
                populateProjectFilter(data.projectStats);

                // Apply default filter (month)
                var filtered = filterData(data, currentRange, selectedProjects);

                // Update summary cards
                document.getElementById('stat-messages').textContent = formatNumber(filtered.totalMessages);
                document.getElementById('stat-sessions').textContent = formatNumber(data.totalSessions);
                document.getElementById('stat-projects').textContent = formatNumber(data.totalProjects);
                document.getElementById('stat-tokens').textContent = formatNumber(filtered.totalTokens);
                document.getElementById('stat-cost').textContent = formatCurrency(filtered.totalCost || 0);
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

        // ===============================
        // Account Filter Functionality
        // ===============================

        var ACCOUNT_STORAGE_KEY = 'claude-code-logs-accounts';
        var accountFilter = document.getElementById('account-filter');
        var accountSettingsBtn = document.getElementById('account-settings-btn');
        var accountModalBackdrop = document.getElementById('account-modal-backdrop');
        var accountModal = document.getElementById('account-modal');
        var accountModalClose = document.getElementById('account-modal-close');
        var accountModalCancel = document.getElementById('account-modal-cancel');
        var accountModalSave = document.getElementById('account-modal-save');
        var accountMappingsContainer = document.getElementById('account-mappings');
        var addMappingBtn = document.getElementById('add-mapping-btn');
        var selectedAccount = '';
        var accountMappings = []; // [{name: 'Work', pattern: '/work/'}, ...]

        // Load account mappings from localStorage
        function loadAccountMappings() {
            try {
                var saved = localStorage.getItem(ACCOUNT_STORAGE_KEY);
                if (saved) {
                    accountMappings = JSON.parse(saved);
                }
            } catch (e) {
                console.warn('Failed to load account mappings:', e);
            }
            populateAccountFilter();
        }

        // Save account mappings to localStorage
        function saveAccountMappings() {
            try {
                localStorage.setItem(ACCOUNT_STORAGE_KEY, JSON.stringify(accountMappings));
            } catch (e) {
                console.warn('Failed to save account mappings:', e);
            }
        }

        // Populate account filter dropdown
        function populateAccountFilter() {
            // Clear existing options except "All Accounts"
            while (accountFilter.options.length > 1) {
                accountFilter.remove(1);
            }
            // Add account options
            accountMappings.forEach(function(mapping) {
                var option = document.createElement('option');
                option.value = mapping.name;
                option.textContent = mapping.name;
                accountFilter.appendChild(option);
            });
            // Add "Unassigned" option if there are any mappings
            if (accountMappings.length > 0) {
                var unassigned = document.createElement('option');
                unassigned.value = '__unassigned__';
                unassigned.textContent = 'Unassigned';
                accountFilter.appendChild(unassigned);
            }
        }

        // Get account for a project based on path matching
        function getProjectAccount(projectPath) {
            for (var i = 0; i < accountMappings.length; i++) {
                var mapping = accountMappings[i];
                if (projectPath.includes(mapping.pattern)) {
                    return mapping.name;
                }
            }
            return null; // Unassigned
        }

        // Filter projects by account
        function filterProjectsByAccount(projects, account) {
            if (!account) return projects; // All accounts
            if (account === '__unassigned__') {
                return projects.filter(function(p) {
                    return getProjectAccount(p.path) === null;
                });
            }
            return projects.filter(function(p) {
                return getProjectAccount(p.path) === account;
            });
        }

        // Account filter change handler
        accountFilter.addEventListener('change', function() {
            selectedAccount = accountFilter.value;
            // Clear project selection when account changes
            selectedProjects = [];
            updateProjectSelection();
            // Repopulate project filter with filtered projects
            if (fullData) {
                var accountFiltered = filterProjectsByAccount(fullData.projectStats, selectedAccount);
                populateProjectFilterWithProjects(accountFiltered);
                var filtered = filterData(fullData, currentRange, selectedProjects);
                updateDisplay(filtered);
            }
        });

        // Populate project filter with specific projects list
        function populateProjectFilterWithProjects(projects) {
            filterOptions.innerHTML = '';
            var sorted = projects.slice().sort(function(a, b) { return b.messages - a.messages; });
            sorted.forEach(function(project) {
                var name = project.path.split('/').pop() || project.slug;
                var displayName = name.length > 35 ? name.substring(0, 32) + '...' : name;

                var option = document.createElement('div');
                option.className = 'multi-select-option';
                option.dataset.slug = project.slug;
                option.innerHTML = '<div class="multi-select-checkbox"><svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd"/></svg></div>' +
                    '<span class="multi-select-option-label">' + displayName + '</span>' +
                    '<span class="multi-select-option-meta">' + formatNumber(project.messages) + ' msgs</span>';

                option.addEventListener('click', function() {
                    toggleProjectSelection(project.slug);
                });

                filterOptions.appendChild(option);
            });
        }

        // Override filterData to include account filtering
        var originalFilterData = filterData;
        filterData = function(data, range, projects) {
            // First filter by account
            var accountFilteredStats = filterProjectsByAccount(data.projectStats, selectedAccount);

            // Create a modified data object with account-filtered stats
            var accountFilteredData = Object.assign({}, data, {
                projectStats: accountFilteredStats
            });

            // Then apply the original filter logic
            // Get base data - either from selected projects or account-filtered
            var selectedData = null;
            if (projects && projects.length > 0) {
                selectedData = accountFilteredStats.filter(function(p) {
                    return projects.indexOf(p.slug) !== -1;
                });
            }

            var aggregated = selectedData && selectedData.length > 0 ? aggregateProjectsData(selectedData) : null;

            // If no projects selected, aggregate all account-filtered projects
            if (!aggregated && selectedAccount) {
                aggregated = aggregateProjectsData(accountFilteredStats);
            }

            var baseMessagesPerDay = aggregated ? aggregated.messagesPerDay : data.messagesPerDay;
            var baseTokensPerDay = aggregated ? aggregated.tokensPerDay : data.tokensPerDay;
            var baseTotalMessages = aggregated ? aggregated.messages : data.totalMessages;
            var baseTotalTokens = aggregated ? aggregated.tokens : data.totalTokens;
            var baseTotalCost = aggregated ? aggregated.cost : (data.totalCost || 0);
            var baseSessions = aggregated ? aggregated.sessions : data.totalSessions;
            var baseAvgLength = aggregated ? aggregated.avgSessionLengthMins : data.avgSessionLengthMins;
            var baseAvgMessages = aggregated ? aggregated.avgMessagesPerSession : data.avgMessagesPerSession;

            // Determine project count
            var projectCount = projects.length > 0 ? projects.length : (selectedAccount ? accountFilteredStats.length : data.totalProjects);

            // If no time filter, return filtered data
            if (range === 'all') {
                return {
                    totalProjects: projectCount,
                    totalSessions: baseSessions,
                    totalMessages: baseTotalMessages,
                    totalTokens: baseTotalTokens,
                    totalCost: baseTotalCost,
                    messagesPerDay: baseMessagesPerDay || [],
                    tokensPerDay: baseTokensPerDay || [],
                    projectStats: accountFilteredStats,
                    avgSessionLengthMins: baseAvgLength,
                    avgMessagesPerSession: baseAvgMessages,
                    selectedProjects: projects
                };
            }

            // Apply time range filter
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
            var filteredMessages = (baseMessagesPerDay || []).filter(function(d) { return d.date >= cutoffStr; });
            var filteredTokens = (baseTokensPerDay || []).filter(function(d) { return d.date >= cutoffStr; });

            // Recalculate totals from filtered data
            var totalMessages = filteredMessages.reduce(function(sum, d) { return sum + d.value; }, 0);
            var totalTokens = filteredTokens.reduce(function(sum, d) { return sum + d.value; }, 0);
            var totalCost = filteredTokens.reduce(function(sum, d) { return sum + (d.cost || 0); }, 0);

            return {
                totalProjects: projectCount,
                totalSessions: baseSessions,
                totalMessages: totalMessages,
                totalTokens: totalTokens,
                totalCost: totalCost,
                messagesPerDay: filteredMessages,
                tokensPerDay: filteredTokens,
                projectStats: accountFilteredStats,
                avgSessionLengthMins: baseAvgLength,
                avgMessagesPerSession: baseAvgMessages,
                selectedProjects: projects
            };
        };

        // Account settings modal handlers
        accountSettingsBtn.addEventListener('click', function() {
            openAccountModal();
        });

        accountModalClose.addEventListener('click', closeAccountModal);
        accountModalCancel.addEventListener('click', closeAccountModal);

        accountModalBackdrop.addEventListener('click', function(e) {
            if (e.target === accountModalBackdrop) {
                closeAccountModal();
            }
        });

        function openAccountModal() {
            renderMappingRows();
            accountModalBackdrop.classList.add('open');
        }

        function closeAccountModal() {
            accountModalBackdrop.classList.remove('open');
        }

        function renderMappingRows() {
            accountMappingsContainer.innerHTML = '';
            accountMappings.forEach(function(mapping, index) {
                addMappingRow(mapping.name, mapping.pattern, index);
            });
            // Add empty row if no mappings
            if (accountMappings.length === 0) {
                addMappingRow('', '', 0);
            }
        }

        function addMappingRow(name, pattern, index) {
            var row = document.createElement('div');
            row.className = 'account-mapping-row';
            row.dataset.index = index;
            row.innerHTML = '<input type="text" class="name-input" placeholder="Account name (e.g., Work)" value="' + (name || '') + '">' +
                '<input type="text" class="pattern-input" placeholder="Path pattern (e.g., /work/)" value="' + (pattern || '') + '">' +
                '<button type="button" class="remove-mapping-btn"><svg viewBox="0 0 20 20" fill="currentColor"><path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/></svg></button>';

            row.querySelector('.remove-mapping-btn').addEventListener('click', function() {
                row.remove();
            });

            accountMappingsContainer.appendChild(row);
        }

        addMappingBtn.addEventListener('click', function() {
            var rows = accountMappingsContainer.querySelectorAll('.account-mapping-row');
            addMappingRow('', '', rows.length);
        });

        accountModalSave.addEventListener('click', function() {
            // Collect mappings from form
            var newMappings = [];
            var rows = accountMappingsContainer.querySelectorAll('.account-mapping-row');
            rows.forEach(function(row) {
                var name = row.querySelector('.name-input').value.trim();
                var pattern = row.querySelector('.pattern-input').value.trim();
                if (name && pattern) {
                    newMappings.push({ name: name, pattern: pattern });
                }
            });

            accountMappings = newMappings;
            saveAccountMappings();
            populateAccountFilter();
            closeAccountModal();

            // Reset account filter and refresh
            selectedAccount = '';
            accountFilter.value = '';
            selectedProjects = [];
            updateProjectSelection();
            if (fullData) {
                populateProjectFilter(fullData.projectStats);
                var filtered = filterData(fullData, currentRange, selectedProjects);
                updateDisplay(filtered);
            }
        });

        // Initialize account mappings on page load
        loadAccountMappings();

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

        // Global keyboard shortcut: / to search
        document.addEventListener('keydown', function(e) {
            if (e.target.matches('input, textarea, [contenteditable]')) return;
            if (e.key === '/') {
                e.preventDefault();
                window.location.href = '/search';
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

/* Stats Filters Container */
.stats-filters {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 16px;
    margin-bottom: 24px;
    flex-wrap: wrap;
}

/* Time Range Filter */
.time-filter {
    display: flex;
    gap: 8px;
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

/* Project Filter */
.project-filter-container {
    display: flex;
    align-items: center;
    gap: 8px;
}

.filter-label {
    font-size: 0.85rem;
    font-weight: 500;
    color: var(--text-secondary);
}

/* Multi-Select Component */
.multi-select-wrapper {
    position: relative;
}

.multi-select-display {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 8px 12px;
    border: 1px solid var(--border-subtle);
    border-radius: 20px;
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    min-width: 150px;
    max-width: 250px;
    transition: all var(--transition-fast);
}

.multi-select-display:hover {
    border-color: var(--border-medium);
    background-color: var(--bg-tertiary);
}

.multi-select-display.active {
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 2px var(--accent-subtle);
}

.multi-select-placeholder {
    flex: 1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.multi-select-count {
    background: var(--accent-primary);
    color: white;
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 0.75rem;
    font-weight: 600;
}

.multi-select-chevron {
    width: 16px;
    height: 16px;
    flex-shrink: 0;
    transition: transform var(--transition-fast);
}

.multi-select-display.active .multi-select-chevron {
    transform: rotate(180deg);
}

.multi-select-dropdown {
    display: none;
    position: absolute;
    top: calc(100% + 4px);
    right: 0;
    min-width: 280px;
    max-width: 350px;
    background: var(--bg-secondary);
    border: 1px solid var(--border-subtle);
    border-radius: 12px;
    box-shadow: var(--shadow-lg);
    z-index: 100;
    overflow: hidden;
}

.multi-select-dropdown.open {
    display: block;
    animation: dropdownFadeIn 0.15s ease-out;
}

@keyframes dropdownFadeIn {
    from { opacity: 0; transform: translateY(-4px); }
    to { opacity: 1; transform: translateY(0); }
}

.multi-select-search-wrapper {
    padding: 8px;
    border-bottom: 1px solid var(--border-subtle);
}

.multi-select-search {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    background: var(--bg-primary);
    color: var(--text-primary);
    font-family: var(--font-body);
    font-size: 0.85rem;
}

.multi-select-search:focus {
    outline: none;
    border-color: var(--accent-primary);
}

.multi-select-options {
    max-height: 240px;
    overflow-y: auto;
    padding: 4px 0;
}

.multi-select-option {
    display: flex;
    align-items: center;
    gap: 10px;
    padding: 10px 12px;
    cursor: pointer;
    transition: background var(--transition-fast);
}

.multi-select-option:hover {
    background: var(--bg-tertiary);
}

.multi-select-option.selected {
    background: var(--accent-subtle);
}

.multi-select-checkbox {
    width: 18px;
    height: 18px;
    border: 2px solid var(--border-medium);
    border-radius: 4px;
    display: flex;
    align-items: center;
    justify-content: center;
    flex-shrink: 0;
    transition: all var(--transition-fast);
}

.multi-select-option.selected .multi-select-checkbox {
    background: var(--accent-primary);
    border-color: var(--accent-primary);
}

.multi-select-checkbox svg {
    width: 12px;
    height: 12px;
    color: white;
    opacity: 0;
}

.multi-select-option.selected .multi-select-checkbox svg {
    opacity: 1;
}

.multi-select-option-label {
    flex: 1;
    font-size: 0.85rem;
    color: var(--text-primary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
}

.multi-select-option-meta {
    font-size: 0.75rem;
    color: var(--text-muted);
}

.multi-select-actions {
    display: flex;
    justify-content: space-between;
    padding: 8px;
    border-top: 1px solid var(--border-subtle);
    background: var(--bg-tertiary);
}

.multi-select-action {
    padding: 6px 12px;
    border: none;
    border-radius: 6px;
    background: transparent;
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.8rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.multi-select-action:hover {
    background: var(--bg-secondary);
}

.multi-select-action.primary {
    background: var(--accent-primary);
    color: white;
}

.multi-select-action.primary:hover {
    background: var(--accent-hover);
}

/* Selected Chips */
.selected-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    margin-bottom: 16px;
    min-height: 0;
}

.selected-chips:empty {
    display: none;
}

.selected-chip {
    display: flex;
    align-items: center;
    gap: 6px;
    padding: 6px 10px;
    background: var(--accent-subtle);
    border: 1px solid var(--accent-primary);
    border-radius: 16px;
    font-size: 0.8rem;
    color: var(--accent-primary);
    font-weight: 500;
}

.selected-chip-remove {
    width: 16px;
    height: 16px;
    border: none;
    background: transparent;
    color: var(--accent-primary);
    cursor: pointer;
    padding: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 50%;
    transition: all var(--transition-fast);
}

.selected-chip-remove:hover {
    background: var(--accent-primary);
    color: white;
}

.selected-chip-remove svg {
    width: 12px;
    height: 12px;
}

/* Filter Row */
.filter-row {
    display: flex;
    gap: 16px;
    align-items: center;
}

/* Account Filter */
.account-filter-container {
    display: flex;
    align-items: center;
    gap: 8px;
}

.account-filter {
    padding: 8px 32px 8px 12px;
    border: 1px solid var(--border-subtle);
    border-radius: 20px;
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.85rem;
    font-weight: 500;
    cursor: pointer;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%2357534E' d='M3 4.5L6 7.5L9 4.5'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 12px center;
    min-width: 130px;
    transition: all var(--transition-fast);
}

.account-filter:hover {
    border-color: var(--border-medium);
    background-color: var(--bg-tertiary);
}

.account-filter:focus {
    outline: none;
    border-color: var(--accent-primary);
    box-shadow: 0 0 0 2px var(--accent-subtle);
}

.account-settings-btn {
    width: 32px;
    height: 32px;
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    background: var(--bg-secondary);
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all var(--transition-fast);
}

.account-settings-btn:hover {
    border-color: var(--border-medium);
    color: var(--text-secondary);
    background: var(--bg-tertiary);
}

.account-settings-btn svg {
    width: 16px;
    height: 16px;
}

/* Modal */
.modal-backdrop {
    display: none;
    position: fixed;
    inset: 0;
    background: rgba(0, 0, 0, 0.5);
    z-index: 1000;
    align-items: center;
    justify-content: center;
    padding: 20px;
}

.modal-backdrop.open {
    display: flex;
    animation: fadeIn 0.15s ease-out;
}

.modal {
    background: var(--bg-secondary);
    border-radius: 16px;
    box-shadow: var(--shadow-lg);
    width: 100%;
    max-width: 500px;
    max-height: 80vh;
    display: flex;
    flex-direction: column;
    animation: modalSlideIn 0.2s ease-out;
}

@keyframes modalSlideIn {
    from { opacity: 0; transform: translateY(-20px); }
    to { opacity: 1; transform: translateY(0); }
}

.modal-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 20px 24px;
    border-bottom: 1px solid var(--border-subtle);
}

.modal-title {
    font-family: var(--font-display);
    font-size: 1.25rem;
    font-weight: 600;
    color: var(--text-primary);
    margin: 0;
}

.modal-close {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-size: 24px;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    transition: all var(--transition-fast);
}

.modal-close:hover {
    background: var(--bg-tertiary);
    color: var(--text-primary);
}

.modal-body {
    padding: 20px 24px;
    overflow-y: auto;
    flex: 1;
}

.modal-description {
    font-size: 0.9rem;
    color: var(--text-secondary);
    margin: 0 0 16px;
    line-height: 1.5;
}

.account-mappings {
    display: flex;
    flex-direction: column;
    gap: 12px;
    margin-bottom: 16px;
}

.account-mapping-row {
    display: flex;
    gap: 8px;
    align-items: center;
}

.account-mapping-row input {
    flex: 1;
    padding: 10px 12px;
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    background: var(--bg-primary);
    color: var(--text-primary);
    font-family: var(--font-body);
    font-size: 0.9rem;
}

.account-mapping-row input:focus {
    outline: none;
    border-color: var(--accent-primary);
}

.account-mapping-row input::placeholder {
    color: var(--text-muted);
}

.account-mapping-row .pattern-input {
    flex: 1.5;
}

.remove-mapping-btn {
    width: 32px;
    height: 32px;
    border: none;
    background: transparent;
    color: var(--text-muted);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 8px;
    transition: all var(--transition-fast);
    flex-shrink: 0;
}

.remove-mapping-btn:hover {
    background: rgba(220, 38, 38, 0.1);
    color: #DC2626;
}

.remove-mapping-btn svg {
    width: 16px;
    height: 16px;
}

.add-mapping-btn {
    display: flex;
    align-items: center;
    gap: 8px;
    padding: 10px 16px;
    border: 1px dashed var(--border-medium);
    border-radius: 8px;
    background: transparent;
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.9rem;
    cursor: pointer;
    transition: all var(--transition-fast);
    width: 100%;
    justify-content: center;
}

.add-mapping-btn:hover {
    border-color: var(--accent-primary);
    color: var(--accent-primary);
    background: var(--accent-subtle);
}

.add-mapping-btn svg {
    width: 16px;
    height: 16px;
}

.modal-footer {
    display: flex;
    justify-content: flex-end;
    gap: 12px;
    padding: 16px 24px;
    border-top: 1px solid var(--border-subtle);
}

.modal-btn {
    padding: 10px 20px;
    border: 1px solid var(--border-subtle);
    border-radius: 8px;
    background: var(--bg-secondary);
    color: var(--text-secondary);
    font-family: var(--font-body);
    font-size: 0.9rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
}

.modal-btn:hover {
    border-color: var(--border-medium);
    background: var(--bg-tertiary);
}

.modal-btn.primary {
    background: var(--accent-primary);
    border-color: var(--accent-primary);
    color: white;
}

.modal-btn.primary:hover {
    background: var(--accent-hover);
    border-color: var(--accent-hover);
}

/* Summary Cards */
.stats-summary {
    display: grid;
    grid-template-columns: repeat(5, 1fr);
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
        grid-template-columns: repeat(3, 1fr);
    }

    .stat-card-wide {
        grid-column: span 2;
    }
}

@media (max-width: 900px) {
    .stats-summary {
        grid-template-columns: repeat(2, 1fr);
    }
}

@media (max-width: 768px) {
    .stats-filters {
        flex-direction: column;
        align-items: flex-start;
    }

    .filter-row {
        flex-direction: column;
        width: 100%;
        gap: 12px;
    }

    .account-filter-container {
        width: 100%;
    }

    .account-filter {
        flex: 1;
    }

    .project-filter-container {
        width: 100%;
    }

    .multi-select-display {
        max-width: none;
    }

    .multi-select-dropdown {
        left: 0;
        right: 0;
        min-width: auto;
        max-width: none;
    }

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

    .modal {
        max-width: 100%;
        margin: 0;
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
