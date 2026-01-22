# Implementation Plan — Run 007

---

## Work Item 3: Stats Page Cost UI

**Mode**: confirm
**Intent**: cost-estimation

### Approach

Add a cost summary card to the stats page and update the JavaScript to display and filter costs. Keep it minimal - one card showing total estimated cost.

### Files to Modify

| File | Changes |
|------|---------|
| `templates_stats.go` | Add cost card HTML, add formatCurrency JS function, update display logic |

### Implementation Details

1. **Add Cost Card** (after Tokens card in stats-summary):
   ```html
   <div class="stat-card">
       <div class="stat-value" id="stat-cost">$0.00</div>
       <div class="stat-label">Est. Cost</div>
   </div>
   ```

2. **Add formatCurrency Function**:
   ```javascript
   function formatCurrency(amount) {
       return '$' + amount.toFixed(2);
   }
   ```

3. **Update Display Logic**:
   - Read `totalCost` from API response
   - Display in cost card with currency formatting
   - Update when time-range filter changes
   - Update when project filter changes
   - Aggregate costs from project data when filtering

4. **Update aggregateProjectsData**:
   - Sum `cost` field from filtered projects
   - Include cost in tokensPerDay aggregation

### Acceptance Criteria Validation

- [x] Cost summary card displaying total estimated cost
- [x] Cost formatted as currency ($X.XX)
- [x] Cost updates when time-range filter changes
- [x] Cost updates when project filter changes
- [x] Cost by project visible (via project filtering)

---

## Work Item 2: Stats Cost Aggregation

**Mode**: autopilot
**Intent**: cost-estimation

### Approach

Extend existing stats structs with cost fields and update ComputeStats() to calculate costs using the CalculateCost function from work item 1.

### Files to Modify

| File | Changes |
|------|---------|
| `stats.go` | Add cost fields to structs, update ComputeStats() |

### Implementation Details

1. **Struct Extensions**:
   - `StatsData`: Add `TotalCost float64`
   - `TimePoint`: Add `Cost float64` (for daily cost tracking)
   - `ProjectStat`: Add `Cost float64`

2. **ComputeStats() Updates**:
   - Track daily costs alongside daily tokens
   - Use estimateTokens() result split as rough 80/20 input/output ratio
   - Sum costs per project
   - Aggregate total cost

3. **Backward Compatibility**:
   - All existing fields unchanged
   - New fields default to 0 if not populated

---

## Work Item 1: Cost Calculation Logic

**Mode**: autopilot
**Intent**: cost-estimation

### Approach

Create a dedicated `cost.go` file with Anthropic pricing constants and a simple cost calculation function. Follow existing codebase patterns (package main, similar file structure to stats.go).

### Files to Create

| File | Purpose |
|------|---------|
| `cost.go` | Pricing constants and CalculateCost function |
| `cost_test.go` | Unit tests for cost calculation |

### Files to Modify

None for this work item.

### Implementation Details

1. **Pricing Constants** (per million tokens):
   - Sonnet 4 Input: $3.00/MTok
   - Sonnet 4 Output: $15.00/MTok

2. **CalculateCost Function**:
   - Takes inputTokens and outputTokens as int64
   - Returns cost in USD as float64
   - Formula: `(input/1M * $3) + (output/1M * $15)`

3. **Test Cases**:
   - Zero tokens → $0.00
   - 1 million input tokens → $3.00
   - 1 million output tokens → $15.00
   - Mixed tokens calculation
   - Small token counts (realistic usage)

### Acceptance Criteria Validation

- [x] New `cost.go` file with pricing constants
- [x] `CalculateCost(inputTokens, outputTokens int64) float64` function
- [x] Default to Sonnet 4 pricing
- [x] Returns cost in USD as float64
- [x] Unit tests for cost calculation
