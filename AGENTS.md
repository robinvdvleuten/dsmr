# AGENTS.md

## Project Overview

Go library for parsing Dutch Smart Meter Requirements (DSMR) telegram data. Transforms raw telegram strings into strongly typed Go structures for smart meter measurements (energy consumption, production, gas readings, meter metadata).

## Commands

```bash
# Run tests
go test ./...

# Run linter
golangci-lint run

# Format code
go fmt ./...

# Build/verify compilation
go build ./...
```

## Project Structure

- `parser.go` - Main parser using [participle](https://github.com/alecthomas/participle) to build AST from DSMR telegrams
- `checksum.go` / `checksum_test.go` - CRC checksum verification
- `options.go` - Parser configuration options
- `error.go` - Custom error types
- `_examples/` - Usage examples

## Code Style

- Go 1.19+ compatible
- Uses `github.com/alecthomas/participle/v2` for parsing grammar
- Uses `github.com/alecthomas/assert/v2` for test assertions
- Run `go fmt` before committing
- Linting configured via `.golangci.yml` (excludes participle struct tag issues)

## Testing

- Use `github.com/alecthomas/assert/v2` for assertions
- Test files follow `*_test.go` naming convention

## Additional Resources

- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
