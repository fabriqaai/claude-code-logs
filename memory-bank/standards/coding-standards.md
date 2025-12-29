# Coding Standards

## Overview

Idiomatic Go standards emphasizing simplicity, readability, and the standard library. Flat package structure appropriate for a focused CLI tool.

## Code Formatting

**Tool**: goimports

- Automatically formats code with gofmt
- Organizes and groups imports automatically
- Run on save or pre-commit

**Enforcement**: Pre-commit hook or editor on-save

## Linting

**Tool**: golangci-lint
**Config**: Default linters

**Key Linters**:
- `errcheck` - Ensure errors are handled
- `govet` - Report suspicious constructs
- `staticcheck` - Static analysis checks
- `unused` - Find unused code

## Naming Conventions

| Element | Convention | Example |
|---------|------------|---------|
| Exported functions | PascalCase | `ParseChatLog`, `GenerateHTML` |
| Unexported functions | camelCase | `readFile`, `formatMessage` |
| Exported types | PascalCase | `ChatSession`, `Message` |
| Unexported types | camelCase | `logEntry`, `config` |
| Constants | PascalCase or camelCase | `DefaultPath`, `maxRetries` |
| Files | snake_case | `chat_parser.go`, `html_generator.go` |
| Test files | snake_case + `_test` | `chat_parser_test.go` |

**Conventions**:
- Receivers: Short, 1-2 letter abbreviation of type (`c` for `ChatSession`)
- Interfaces: Verb + "er" suffix when possible (`Reader`, `Parser`)
- Booleans: Use `is`, `has`, `can` prefixes for clarity

## File Organization

**Pattern**: Flat (single package)

**Structure**:
```text
claude-code-logs/
  main.go           # Entry point, CLI setup
  parser.go         # Chat log parsing
  generator.go      # HTML generation
  templates.go      # HTML templates
  *_test.go         # Tests alongside source
  go.mod
  go.sum
```

**Conventions**:
- Tests: Co-located with source files (`foo_test.go`)
- One primary type per file when it makes sense
- Keep files focused and reasonably sized

## Testing Strategy

**Framework**: Standard `testing` package
**Coverage Target**: Critical paths (parsing, generation)

**Conventions**:
- Test naming: `TestFunctionName_Scenario`
- Table-driven tests for multiple cases
- Use `t.Helper()` for test utilities
- Use `t.Parallel()` where safe

**Example**:
```go
func TestParseMessage_ValidJSON(t *testing.T) {
    // test implementation
}

func TestParseMessage_InvalidJSON(t *testing.T) {
    // test implementation
}
```

## Error Handling

**Pattern**: Return errors, wrap with context

**Conventions**:
- Always check returned errors
- Wrap errors with `fmt.Errorf("context: %w", err)`
- Use sentinel errors for known conditions
- Exit with non-zero code on CLI errors

**Example**:
```go
data, err := os.ReadFile(path)
if err != nil {
    return fmt.Errorf("reading chat log %s: %w", path, err)
}
```

## Logging

**Tool**: Standard `log` package or `fmt` for CLI output
**Format**: Text (human-readable CLI output)

**Conventions**:
- Use `fmt.Println` for normal output
- Use `fmt.Fprintf(os.Stderr, ...)` for errors
- Verbose mode flag for debug output
- No logging of sensitive data
