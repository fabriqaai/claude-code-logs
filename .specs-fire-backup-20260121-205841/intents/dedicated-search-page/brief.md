---
id: dedicated-search-page
title: Dedicated Search Page
status: in_progress
created: 2026-01-20
---

# Intent: Dedicated Search Page

## Goal

Replace the current search overlay with a dedicated `/search` page that provides better UX for finding content across sessions - with relevance sorting, pagination, and highlighted results.

## Users

Anyone using claude-code-logs to search across their Claude Code sessions.

## Problem

The current overlay-based search is limited:
- No sorting - results appear in arbitrary order
- No pagination - all results load at once, poor performance at scale
- Hard to scan - excerpts don't highlight matching terms
- Modal UX - can't easily compare results or navigate

## Success Criteria

- `/search` route serves a dedicated search page
- Results sorted by match quality (relevance) by default
- "Load more" pagination for large result sets
- Expandable excerpts with inline term highlighting
- Keyboard shortcut (`/`) navigates to search page
- Search feels fast and scannable even with 1000+ results

## Constraints

- Integrate with existing `/api/search` endpoint (extend if needed)
- Maintain Go templates + vanilla JS approach (no new frameworks)
- Keep consistent styling with existing pages (sessions, stats)
- Remove old overlay search code

## Notes

- May need to extend search API to support pagination params (offset/limit)
- Consider adding sort parameter to API for future sort options
- Relevance scoring may need tuning based on match position, frequency, etc.
