---
name: go-patterns
description: >
  Architecture patterns for Glens Go code. Use this when writing or reviewing
  Go functions — composable design, error handling, and logging patterns.
---

# Skill: Go Architecture Patterns

## Composable Functions

Prefer small, composable functions with clear responsibilities:

```go
func ProcessEndpoint(ep *Endpoint) error {
    if err := validate(ep); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    return process(ep)
}
```

## Error Handling

Always wrap errors with context using `fmt.Errorf` and `%w`:

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

## Logging (via pkg/logging)

```go
log.Info().Str("key", val).Msg("What happened")
log.Error().Err(err).Msg("What failed")
```

## Critical Logic: GitHub Issue Creation

Issues are created **only** when:

- Tests execute successfully (no compilation errors).
- Tests fail with assertion errors (spec violations).
- Failures are genuine OpenAPI spec mismatches.

Issues are **not** created when:

- Connection to the server fails.
- Test compilation fails.
- Infrastructure / setup issues occur.

**Key function**: `isRealTestFailure()` in `cmd/glens/cmd/analyze.go` — preserve
this distinction when modifying test-failure detection logic.
