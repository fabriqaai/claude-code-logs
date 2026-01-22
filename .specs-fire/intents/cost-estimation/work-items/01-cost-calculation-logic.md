---
id: 01-cost-calculation-logic
title: Cost Calculation Logic
intent: cost-estimation
complexity: low
mode: autopilot
status: pending
depends_on: []
created: 2026-01-22
---

# Work Item: Cost Calculation Logic

## Description

Add cost calculation functions with Anthropic pricing constants. Create a dedicated `cost.go` file with pricing data and a function to calculate cost from token counts.

## Acceptance Criteria

- [ ] New `cost.go` file with pricing constants for Claude models
- [ ] `CalculateCost(inputTokens, outputTokens int64) float64` function
- [ ] Default to Sonnet 4 pricing ($3/MTok input, $15/MTok output)
- [ ] Returns cost in USD as float64
- [ ] Unit tests for cost calculation

## Technical Notes

Pricing constants (per million tokens):
```go
const (
    SonnetInputPrice  = 3.0   // $/MTok
    SonnetOutputPrice = 15.0  // $/MTok
)

func CalculateCost(inputTokens, outputTokens int64) float64 {
    inputCost := float64(inputTokens) / 1_000_000 * SonnetInputPrice
    outputCost := float64(outputTokens) / 1_000_000 * SonnetOutputPrice
    return inputCost + outputCost
}
```

## Dependencies

None - this is the foundation work item.
