---
id: session-stats-dashboard
title: Session Stats Dashboard
status: in_progress
created: 2025-01-20
---

# Intent: Session Stats Dashboard

## Goal

Add a dedicated `/stats` page with usage analytics, giving users insight into their Claude Code usage patterns over time.

## Users

claude-code-logs users who want to understand their Claude Code usage patterns, compare project activity, and track trends.

## Problem

Currently no way to see usage trends, compare projects, or understand how you're using Claude Code over time. Users have data but no visibility into patterns.

## Success Criteria

- Dedicated `/stats` page accessible from main navigation
- Charts displaying:
  - Messages per day/week
  - Token usage (estimated from message content length)
  - Active projects over time
  - Session duration/length
- Time range toggles: today, this week, this month, all time
- Fast load time (pre-computed or cached stats)

## Constraints

- Token counts must be estimated (not stored in JSONL source files)
- Use stdlib + existing patterns where possible
- Client-side charting (no server-side image generation)
- No new Go dependencies if avoidable

## Notes

- Consider using lightweight JS charting library (Chart.js or similar)
- Stats can be computed at server startup and cached
- API endpoint for stats data (`/api/stats` already exists, can be extended)
