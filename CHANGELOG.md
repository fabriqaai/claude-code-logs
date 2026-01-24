# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.1.27] - 2026-01-22

### Added
- Cost estimation on stats dashboard - estimated costs based on token usage and model pricing
- Cost breakdown by model type (Opus, Sonnet, Haiku) with input/output token differentiation
- Total estimated cost displayed in stats summary

## [0.1.26] - 2026-01-21

### Added
- Project-based filtering on stats dashboard - filter messages, tokens, and charts by individual project
- Project filter dropdown selector with "All Projects" default option
- Per-project time series data (`MessagesPerDay`, `TokensPerDay`) in `ProjectStat` struct
- Project Activity chart now highlights selected project

### Changed
- Stats filtering now combines project and time-range filters client-side for responsive UX

## [0.1.24] - 2026-01-21

### Added
- Dedicated search page at `/search` route
- Keyboard shortcut `/` to navigate to search from any page
- Navigation links to search page in sidebar

### Removed
- Overlay-based search modal (replaced with dedicated page)

## [0.1.23] - 2026-01-20

### Added
- Stats dashboard with usage analytics
- Messages per day and tokens per day charts
- Project activity breakdown chart
- Time range filtering (Today, This Week, This Month, All Time)
- Summary cards for total messages, sessions, projects, and estimated tokens

## [0.1.22] - 2026-01-19

### Added
- Initial release with core features
- Markdown generation from Claude Code JSONL logs
- Full-text search across all messages
- Inline session search with highlighting
- Tree view sidebar with collapsible projects
- File watching for auto-regeneration
- Mobile responsive design

[0.1.27]: https://github.com/fabriqaai/claude-code-logs/compare/v0.1.26...v0.1.27
[0.1.26]: https://github.com/fabriqaai/claude-code-logs/compare/v0.1.24...v0.1.26
[0.1.24]: https://github.com/fabriqaai/claude-code-logs/compare/v0.1.23...v0.1.24
[0.1.23]: https://github.com/fabriqaai/claude-code-logs/compare/v0.1.22...v0.1.23
[0.1.22]: https://github.com/fabriqaai/claude-code-logs/releases/tag/v0.1.22
