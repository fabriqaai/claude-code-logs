---
run: run-006
work_item: 01-stats-data-model-extension
intent: project-stats-filtering
mode: confirm
checkpoint: plan-approval
approved_at:
---

# Implementation Plan: Extend Stats Data Model

## Approach

Add per-project time series data to enable client-side project filtering on the stats page. The existing aggregated data at the root level will remain unchanged for "All Projects" view, while each `ProjectStat` will gain its own `MessagesPerDay` and `TokensPerDay` arrays.

**Key design decisions:**
1. Track per-project time series during the existing `ComputeStats()` loop (single pass)
2. Use maps keyed by project slug for O(1) lookups during aggregation
3. Apply the same 30-day window as root-level time series
4. Keep backward compatibility - all existing fields remain unchanged

## Files to Create

| File | Purpose |
|------|---------|
| (none) | |

## Files to Modify

| File | Changes |
|------|---------|
| `stats.go` | Extend `ProjectStat` struct with `MessagesPerDay` and `TokensPerDay` fields; update `ComputeStats()` to populate per-project time series |
| `stats_test.go` | Add tests for per-project time series computation |

## Tests

| Test File | Coverage |
|-----------|----------|
| `stats_test.go` | Per-project time series population, correct aggregation by day per project |

## Technical Details

### Struct Changes (stats.go:37-45)

```go
type ProjectStat struct {
    Path           string      `json:"path"`
    Slug           string      `json:"slug"`
    Sessions       int         `json:"sessions"`
    Messages       int         `json:"messages"`
    Tokens         int         `json:"tokens"`
    LastUsed       time.Time   `json:"lastUsed"`
    MessagesPerDay []TimePoint `json:"messagesPerDay"` // NEW
    TokensPerDay   []TimePoint `json:"tokensPerDay"`   // NEW
}
```

### Algorithm Changes (ComputeStats)

1. Create maps to track per-project daily data:
   ```go
   // Map: projectSlug -> date -> value
   projectMessagesByDay := make(map[string]map[string]int)
   projectTokensByDay := make(map[string]map[string]int)
   ```

2. During message iteration (inside the project loop), aggregate per-project:
   ```go
   slug := ProjectSlug(project.Path)
   if projectMessagesByDay[slug] == nil {
       projectMessagesByDay[slug] = make(map[string]int)
       projectTokensByDay[slug] = make(map[string]int)
   }
   projectMessagesByDay[slug][day]++
   projectTokensByDay[slug][day] += tokens
   ```

3. After the loop, convert each project's maps to time series:
   ```go
   projectStat.MessagesPerDay = buildTimeSeries(projectMessagesByDay[slug], 30)
   projectStat.TokensPerDay = buildTimeSeries(projectTokensByDay[slug], 30)
   ```

### Performance Consideration

The additional per-project time series adds memory overhead proportional to `O(projects × 30 days)`. For typical usage (10-50 projects), this is negligible.

---
*Plan generated for checkpoint approval.*

---

# Implementation Plan: Add Project Filter UI

## Work Item: 02-project-filter-ui

**Mode**: Autopilot (no checkpoint)

## Approach

Add a project dropdown selector to the stats page filter area, positioned next to the existing time-range filter buttons. The dropdown will be populated dynamically from `projectStats` data and will default to "All Projects".

## Files to Create

| File | Purpose |
|------|---------|
| (none) | |

## Files to Modify

| File | Changes |
|------|---------|
| `templates_stats.go` | Add project filter dropdown HTML and JavaScript to handle selection |

## Technical Details

1. **HTML** - Add a `<select>` element with class `project-filter` after the time-filter div
2. **JavaScript** - Populate options from `projectStats`, store selection in `currentProject` variable
3. **CSS** - Style the dropdown to match existing filter button aesthetics

---
*Plan saved (Autopilot mode - continuing without checkpoint)*

---

# Implementation Plan: Implement Client-Side Project Filtering

## Work Item: 03-client-side-filtering

**Mode**: Confirm (checkpoint for plan approval)

## Approach

Wire up the project selector to filter all stats displays by modifying the JavaScript filtering logic. When a project is selected, the summary cards and time series charts will use per-project data. When "All Projects" is selected, the aggregated data is used (existing behavior).

**Key implementation points:**
1. Create a unified `getFilteredData()` function that applies both time-range AND project filters
2. When project selected → use `projectStats[slug].messagesPerDay/tokensPerDay` for charts
3. When project selected → use `projectStats[slug].messages/tokens` for summary totals
4. Apply time-range filter to per-project time series data
5. Project Activity chart highlights selected project (or shows all if "All Projects")

## Files to Create

| File | Purpose |
|------|---------|
| (none) | |

## Files to Modify

| File | Changes |
|------|---------|
| `templates_stats.go` | Modify `filterDataByRange()` to handle project filtering; update `updateDisplay()` to use project-specific data; update Project Activity chart to highlight selected project |

## Technical Details

### 1. Create `getProjectData()` Helper

```javascript
function getProjectData(projectSlug) {
    if (!projectSlug) return null;
    return fullData.projectStats.find(function(p) { return p.slug === projectSlug; });
}
```

### 2. Modify `filterDataByRange()` to Handle Both Filters

The function will:
1. Get base data from selected project (or use root-level if "All Projects")
2. Apply time-range filter to the time series
3. Return filtered result with correct totals

### 3. Update `updateDisplay()`

- Update summary cards (messages, tokens) from filtered data
- Update sessions/projects counts appropriately when project filtered
- Redraw charts with filtered time series

### 4. Project Activity Chart Behavior

When a project is selected:
- Highlight the selected project bar with accent color
- Dim other project bars (lower opacity)
- Or optionally filter to show only the selected project

---
*Plan requires checkpoint approval.*
