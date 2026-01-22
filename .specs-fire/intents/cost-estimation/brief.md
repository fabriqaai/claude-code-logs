---
id: cost-estimation
title: Cost Estimation on Stats Dashboard
status: in_progress
created: 2026-01-22
---

# Intent: Cost Estimation on Stats Dashboard

## Goal

Add estimated API cost tracking to the stats dashboard, showing cost breakdowns by project and over time based on token usage.

## Users

Developers using Claude Code who want to understand and monitor their API spending patterns.

## Problem

Currently the stats page shows token counts but doesn't translate these into dollar costs. Users have to manually calculate costs using Anthropic's pricing, which is tedious and error-prone.

## Success Criteria

- Total estimated cost displayed on stats page
- Cost breakdown by project (matches existing project filtering)
- Cost over time chart (daily costs)
- Costs update when time-range filter changes
- Currency formatted as USD ($X.XX)

## Constraints

- Stats page only (no per-session cost display)
- Hardcoded pricing defaults (Sonnet 4 as default model)
- Combined input+output cost (not separated)
- Integrates with existing stats infrastructure

## Pricing Reference

| Model | Input | Output |
|-------|-------|--------|
| Claude Sonnet 4/4.5 | $3/MTok | $15/MTok |
| Claude Opus 4.5 | $5/MTok | $25/MTok |
| Claude Opus 4/4.1 | $15/MTok | $75/MTok |
| Claude Haiku 4.5 | $1/MTok | $5/MTok |
| Claude Haiku 3.5 | $0.80/MTok | $4/MTok |

Default: Sonnet 4 pricing ($3 input, $15 output per million tokens)

## Notes

- Token counts already tracked in existing stats infrastructure
- Client-side cost calculation matches existing time-range filtering pattern
- Future enhancement: detect model from JSONL if available
