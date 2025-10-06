# GitHub Copilot Instructions

## Project: Glens - OpenAPI Test Generator

### Quick Context

Go-based CLI tool that generates integration tests from OpenAPI specs using AI models.
Creates GitHub issues ONLY when tests fail against the spec.

### Golden Rules

1. **Simplicity**: Simplest working solution wins
2. **Go Idioms**: Follow standard Go patterns
3. **Minimal Docs**: Update only what users/contributors need
4. **Long-term**: Code must be maintainable for years

### Code Style

- Use `gofmt` and `golint` standards
- Explicit error handling (no silent failures)
- Small functions (<50 lines)
- Clear variable names (no cryptic abbreviations)
- Comments for "why", not "what"

### Architecture Patterns

```go
// Prefer this
func ProcessEndpoint(ep *Endpoint) error {
    if err := validate(ep); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    return process(ep)
}

// Not this
func ProcessEndpoint(ep *Endpoint) error {
    if ep == nil || ep.Path == "" || ep.Method == "" {
        if ep == nil {
            return errors.New("endpoint nil")
        }
        if ep.Path == "" {
            return errors.New("path empty")
        }
        return errors.New("method empty")
    }
    // ... complex nested logic
}
```

### Critical Logic: GitHub Issue Creation

Issues created ONLY when:

- Tests execute successfully (no compilation errors)
- Tests fail with assertion errors (spec violations)
- Failures are from OpenAPI spec mismatches

Issues NOT created when:

- Connection to server fails
- Test compilation fails
- Infrastructure/setup issues

Function: `isRealTestFailure()` in `cmd/analyze.go`

### Testing

```go
// Use table-driven tests
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

### Documentation

**Only two docs needed:**

- `docs/QUICKSTART.md` - User getting started
- `docs/DEVELOPMENT.md` - Contributor guide
- `docs/diagrams/` - Architecture diagrams (Mermaid)

Update only when:

- CLI interface changes
- Setup process changes
- New configuration options
- Major architecture changes (diagrams)

### Dependencies

- Prefer standard library
- Justify any new dependency
- Keep `go.mod` clean

### File Structure

```txt
cmd/           # CLI commands
pkg/           # Core packages
  ├── ai/      # AI model clients
  ├── github/  # GitHub integration
  ├── parser/  # OpenAPI parsing
  └── generator/ # Test generation
configs/       # Example configs only
docs/          # Minimal docs
```

### Common Tasks

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

**Error handling:**

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

**Logging:**

```go
log.Info().Str("key", val).Msg("What happened")
log.Error().Err(err).Msg("What failed")
```

### What to Avoid

- Complex abstractions
- Premature optimization
- Extensive documentation
- Duplicate code/docs
- Magic numbers (use constants)
- Global state

### When Suggesting Changes

Ask:

1. Simplest solution?
2. Maintainable long-term?
3. Follows Go idioms?
4. Necessary?

If all yes → implement
If any no → reconsider

### Current Focus

Feature complete but under active development.
Prioritize code quality over feature additions.
