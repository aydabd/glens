---
applyTo: "**"
---

# Glens Project – AI Agent Instructions

> **Single source of truth** for all AI coding agents (GitHub Copilot, Claude Code,
> OpenAI Codex, Cursor, and others).
> Edit **only this file** to update instructions across every agent.

## Project Overview

Glens is a Go-based CLI tool that generates integration tests from OpenAPI specs
using AI models. It creates GitHub issues **only** when tests fail against the spec.

## Golden Rules

1. **Simplicity** – Simplest working solution wins.
2. **Go Idioms** – Follow standard Go patterns (`gofmt`, `golint`).
3. **Minimal Docs** – Update only what users/contributors need.
4. **Long-term** – Code must be maintainable for years.

## Code Style

- Use `gofmt` and `golint` standards.
- Keep functions small and focused (max 50 lines).
- Use meaningful variable names; avoid abbreviations.
- Add comments only when code is not self-explanatory (explain *why*, not *what*).
- Handle errors explicitly — never ignore them.
- Write tests for new functionality.

## Architecture Patterns

```go
// Prefer this – small, composable functions
func ProcessEndpoint(ep *Endpoint) error {
    if err := validate(ep); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    return process(ep)
}

// Avoid this – deep nesting / repeated nil-checks
func ProcessEndpoint(ep *Endpoint) error {
    if ep == nil || ep.Path == "" || ep.Method == "" {
        if ep == nil {
            return errors.New("endpoint nil")
        }
        // ...
    }
}
```

### Error handling

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

### Configuration struct

```go
type Config struct {
    // Required fields first
    Token string `yaml:"token"`
    // Optional fields after
    Debug bool `yaml:"debug,omitempty"`
}
```

### Logging

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

**Key function**: `isRealTestFailure()` in `cmd/analyze.go` — preserve this distinction
when modifying test-failure detection logic.

## Testing

Use table-driven tests with `testify`:

```go
func TestSomething(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {"valid", "input", "output", false},
        {"invalid", "", "", true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Function(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.Equal(t, tt.want, got)
        })
    }
}
```

- Mock external dependencies.
- Test error cases, not just happy paths.
- Keep test code as simple as production code.

## Documentation

Two docs are needed for users/contributors:

- `docs/QUICKSTART.md` – Getting started guide.
- `docs/DEVELOPMENT.md` – Contributor guide.
- `docs/diagrams/` – Architecture diagrams (Mermaid).

Update docs **only** when:

- CLI interface changes.
- Setup process changes.
- New configuration options are added.
- Major architecture changes (diagrams).

Do **not** update docs for internal refactoring, bug fixes, code reorganization,
or performance improvements.

## File Structure

```txt
cmd/           # CLI commands (one per file)
pkg/           # Core packages
  ├── ai/      # AI model clients
  ├── github/  # GitHub integration
  ├── parser/  # OpenAPI parsing
  └── generator/ # Test generation
configs/       # Example configurations only
docs/          # Minimal docs
```

## Common Tasks

**Adding a CLI flag:**

```go
cmd.Flags().String("flag-name", "default", "Clear description")
_ = viper.BindPFlag("config.key", cmd.Flags().Lookup("flag-name"))
```

**Adding a config option:**

```yaml
# configs/config.example.yaml
new_option: "${ENV_VAR}" # default_value
```

## Before Any Change – Checklist

Ask yourself:

1. Is this the simplest solution?
2. Will this be easy to maintain in two years?
3. Can someone new understand this in five minutes?
4. Does this follow Go idioms?
5. Is this change necessary?

If any answer is "no", reconsider the approach.

## What to Avoid

- ❌ Features added "just in case"
- ❌ Complex abstractions or inheritance hierarchies
- ❌ Reflection unless absolutely necessary
- ❌ Documentation novels / duplicate docs
- ❌ Magic numbers (use named constants)
- ❌ Global state
- ❌ New dependencies without strong justification
- ❌ Premature optimization
