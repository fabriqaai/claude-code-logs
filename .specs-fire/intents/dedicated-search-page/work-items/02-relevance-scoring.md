---
id: 02-relevance-scoring
title: Add Relevance Scoring
intent: dedicated-search-page
complexity: medium
mode: confirm
status: pending
depends_on: []
created: 2026-01-20
---

# Work Item: Add Relevance Scoring

## Description

Implement a relevance scoring system for search results. Results should be ranked by match quality rather than just recency.

## Acceptance Criteria

- [ ] Each `SearchResult` includes a `score` field (float64)
- [ ] Scoring factors: term frequency, match position (earlier in content = higher), exact phrase bonus
- [ ] Results sorted by score descending by default
- [ ] API accepts optional `sort` param: `relevance` (default) or `recent`

## Technical Notes

Scoring formula (suggested):
- Base: 1.0 per matching message in session
- Term frequency bonus: +0.1 per additional term occurrence (capped)
- Position bonus: +0.5 if match in first 200 chars
- Exact phrase bonus: +2.0 per phrase match
- Recency bonus: +0.2 if within last 7 days

Keep scoring simple and tunable. Can refine later based on usage.

## Files to Modify

- `search.go` - Add scoring logic, update SearchResult struct, add sort param
- `server.go` - Parse sort param from request
