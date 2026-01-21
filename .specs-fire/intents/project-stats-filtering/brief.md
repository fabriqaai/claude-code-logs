---
id: project-stats-filtering
title: Project-Based Stats Filtering
status: in_progress
created: 2026-01-21
---

# Intent: Project-Based Stats Filtering

## Goal

Add a project filter to the stats page so users can view statistics for specific project(s) instead of all projects combined.

## Users

Developer viewing Claude Code session logs and wanting to analyze usage patterns for individual projects.

## Problem

Currently the stats page shows aggregated data across all projects. There's no way to drill down into individual project statistics to understand per-project usage patterns.

## Success Criteria

- Project dropdown/selector on stats page
- All stats (messages, tokens, charts) filter by selected project
- Works alongside existing time-range filtering (project AND time range)
- "All Projects" option to return to aggregated view

## Constraints

- Integrates with existing stats infrastructure (`StatsData`, `/api/stats`)
- Follows current UI patterns (similar to time-range filter buttons or dropdown)
- Client-side filtering preferred (matches existing time-range pattern)

## Notes

- Existing infrastructure already has `ProjectStats` with per-project data
- Time-range filtering is done client-side for responsiveness
- Project Activity chart already shows top 10 projects by message count
