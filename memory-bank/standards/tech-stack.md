# Tech Stack

## Overview

Go-based CLI tool for scanning Claude Code chat logs and generating HTML reports. Go provides fast binaries, cross-platform support, and zero runtime dependencies - ideal for developer tooling.

## Languages

**Go**

- Fast compilation and execution
- Single binary distribution (no runtime dependencies)
- Excellent standard library for file I/O and HTML templating
- Strong concurrency primitives for parallel file processing

## Framework

**Standard Library + Cobra**

- `html/template` - Built-in HTML templating
- `encoding/json` - JSON parsing for chat logs
- `cobra` - Industry-standard CLI framework for Go
- `filepath` - Cross-platform path handling

## Authentication

N/A - Local CLI tool, no authentication required.

## Infrastructure & Deployment

**Homebrew**

- Primary distribution via Homebrew tap
- Easy installation: `brew install <tap>/claude-code-logs`
- Automatic updates via `brew upgrade`

## Package Manager

**Go Modules**

- Standard Go dependency management
- `go.mod` and `go.sum` for reproducible builds

## Decision Relationships

- Go's `html/template` provides secure HTML generation with auto-escaping
- Cobra integrates well with Go's flag package and provides shell completions
- Homebrew formula can be auto-generated from goreleaser config
