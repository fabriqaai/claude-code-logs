# Coding Standards

## Project Structure

- Flat structure: all Go files in root directory
- Test files alongside source (`*_test.go`)
- Single `main` package

## Naming Conventions

- Files: `snake_case.go`
- Types: `PascalCase` (e.g., `Session`, `SearchIndex`)
- Functions: `PascalCase` for exported, `camelCase` for unexported
- Variables: `camelCase`

## Error Handling

- Wrap errors with context: `fmt.Errorf("doing X: %w", err)`
- Return errors to caller, don't panic
- Log warnings to stderr for non-fatal issues

## Testing

- Table-driven tests preferred
- Test files named `*_test.go`
- Run with `go test -v ./...`

## Code Patterns

### Atomic File Writes
```go
tmpFile, _ := os.CreateTemp(dir, "tmp-*.html")
// write to tmpFile
os.Rename(tmpPath, outputPath)
```

### Embedded Assets
```go
//go:embed asset.png
var assetData []byte
```

### Graceful Shutdown
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
```

### HTTP Server
- Use stdlib `net/http`
- CORS middleware for API
- JSON request/response for APIs

## Documentation

- Comments for exported types and functions
- README.md for user-facing docs
- No excessive inline comments
