# Walkthrough: Time Range Filtering & Navigation Integration

## Summary

Added time range filter buttons (Today, This Week, This Month, All Time) to the stats page, and integrated a Stats navigation link into the sidebar across all pages.

## Work Item 04: Time Range Filtering

### Changes to templates_stats.go

**1. Added filter buttons HTML:**
```html
<div class="time-filter">
    <button type="button" class="time-filter-btn" data-range="today">Today</button>
    <button type="button" class="time-filter-btn" data-range="week">This Week</button>
    <button type="button" class="time-filter-btn active" data-range="month">This Month</button>
    <button type="button" class="time-filter-btn" data-range="all">All Time</button>
</div>
```

**2. Added filtering JavaScript:**
- `filterDataByRange(data, range)` - Filters time series data by cutoff date
- Click handlers for filter buttons
- `updateMessagesChart()` and `updateTokensChart()` - Update charts with new data
- Full data stored in `fullData` variable for instant filtering

**3. Added filter button CSS:**
- Pill-shaped buttons with hover and active states
- Active button uses accent color (terracotta)
- Responsive wrapping on mobile

### How It Works

1. On page load, full stats data is fetched and stored
2. Default filter "This Month" is applied
3. Clicking a filter button:
   - Updates active button state
   - Filters `messagesPerDay` and `tokensPerDay` arrays
   - Recalculates totals from filtered data
   - Updates summary cards and charts

## Work Item 05: Navigation Integration

### Changes to Templates

Added Stats navigation link to sidebar header in:
- `templates_index.go` - Index page
- `templates_project.go` - Project pages  
- `templates_shell.go` - Session pages
- `templates_stats.go` - Stats page (with `.active` class)

**Link HTML:**
```html
<a href="stats.html" class="stats-nav-link">
    <svg viewBox="0 0 20 20" fill="currentColor">
        <path d="M2 11a1 1 0 011-1h2..."/>
    </svg>
    Stats
</a>
```

### Changes to templates_css.go

Added `.stats-nav-link` styles:
- Flexbox layout with icon and text
- Border and hover states
- Active state with accent color
- Positioned below sidebar subtitle

## Verification

```bash
# Start the server
go run . serve

# Open stats page
open http://localhost:8080/stats

# Test:
# 1. Click filter buttons - charts should update instantly
# 2. Navigate to index page - Stats link should appear in sidebar
# 3. Navigate to any project - Stats link should appear
# 4. Navigate to any session - Stats link should appear
# 5. Click Stats link - should navigate to stats page with .active state
```

## Acceptance Criteria

### Time Range Filtering
- [x] Four toggle buttons displayed: Today, This Week, This Month, All Time
- [x] Clicking a button filters charts to that time range
- [x] Active button is visually highlighted
- [x] Default selection is "This Month"
- [x] Filter state persists during page session
- [x] Charts update smoothly when filter changes

### Navigation Integration
- [x] "Stats" link visible in navigation/header
- [x] Link appears on index page, project pages, and session pages
- [x] Link styling consistent with existing navigation
- [x] Current page indicator when on stats page (`.active` class)
