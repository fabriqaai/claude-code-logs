# Walkthrough: Stats Page with Charts

## Summary

Created a new `/stats` page that displays analytics charts for Claude Code usage, including messages over time, token usage, and project activity.

## What Was Created

### templates_stats.go

New file containing:
- `statsTemplate` - Full HTML template for the stats page
- `statsCSS` - Stats-specific CSS styling

**Template Structure:**
- Same sidebar navigation as other pages (index, project, session)
- Summary cards grid: Total Messages, Sessions, Projects, Est. Tokens
- Charts grid:
  - Messages per Day (line chart)
  - Tokens per Day (line chart)
  - Project Activity (horizontal bar chart)
- Session stats: Avg Session Length, Avg Messages/Session

**Chart.js Integration:**
```html
<script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
```

Charts fetch data from `/api/stats` on page load and render using Chart.js:
- Line charts with smooth curves and fill
- Horizontal bar chart for projects (top 10)
- Consistent color scheme using CSS variables

### server.go Changes

1. Added `statsTmpl` field to `Server` struct
2. Added template parsing in `NewServer`
3. Added `/stats` and `/stats.html` route handling in `handleStatic`
4. Added `renderStatsPage` handler with caching

## Page URL

The stats page is accessible at:
- `/stats`
- `/stats.html`

## Verification

```bash
# Start the server
go run . serve

# Open in browser
open http://localhost:8080/stats
```

**Expected:**
- Loading spinner while fetching `/api/stats`
- Summary cards with formatted numbers (1.2K, 1.5M, etc.)
- Messages chart showing last 30 days
- Tokens chart showing estimated usage
- Project activity bar chart (top 10 projects by messages)
- Responsive layout on mobile

## Acceptance Criteria

- [x] `/stats` route serves stats page
- [x] Page displays 4 charts (3 chart canvases + summary cards)
- [x] Charts render correctly with data from `/api/stats`
- [x] Page follows existing UI style (sidebar, fonts, colors)
- [x] Charts are responsive (flexbox grid)
- [x] Loading state shown while fetching data
