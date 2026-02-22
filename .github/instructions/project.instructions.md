---
applyTo: "**"
---

# Glens Project — AI Agent Instructions

> **Single source of truth** for all AI coding agents (GitHub Copilot, Claude Code,
> OpenAI Codex, Cursor, and others).
> Edit **only this file** to update instructions across every agent.

## Project Overview

Glens is a Go workspace monorepo. Every binary and library lives in its own isolated module
under `cmd/` or `pkg/`. Each module is independently buildable, testable, and releasable —
and can be moved to a separate repository without changes.

```
go.work
├── pkg/logging           # module glens/pkg/logging    — generic zerolog wrapper
├── cmd/glens             # module glens/tools/glens    — main CLI
├── cmd/tools/demo        # module glens/tools/demo     — OpenAPI spec visualiser
└── cmd/tools/accuracy    # module glens/tools/accuracy — endpoint accuracy reporter
```

## Golden Rules

1. **Simplicity** — simplest working solution wins.
2. **Go Idioms** — follow standard Go patterns (`gofmt`, `golangci-lint`).
3. **Minimal Docs** — update only what users or contributors need.
4. **Long-term** — code must be maintainable for years.
5. **Isolation** — never import across module boundaries; each module owns its `internal/`.

## Package Layout Rules

- `pkg/` — generic, reusable libraries. **Must never import from any `internal/` package.**
  Versioned independently via semver git tags (e.g. `pkg/logging/v0.1.0`).
- `cmd/*/internal/` — module-private implementation. Never imported from another module.
- `cmd/glens` imports `glens/pkg/logging` via a workspace `replace` directive.

## Code Style

- Use `gofmt` and `golangci-lint` standards.
- Keep functions small and focused (max 50 lines).
- Use meaningful variable names; avoid abbreviations.
- Add comments only when code is not self-explanatory (explain *why*, not *what*).
- Handle errors explicitly — never ignore them.
- Write tests for new functionality.

## Architecture Patterns

```go
// Prefer this — small, composable functions
func ProcessEndpoint(ep *Endpoint) error {
    if err := validate(ep); err != nil {
        return fmt.Errorf("validation: %w", err)
    }
    return process(ep)
}
```

### Error handling

```go
if err != nil {
    return fmt.Errorf("operation context: %w", err)
}
```

### Logging (via pkg/logging)

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

**Key function**: `isRealTestFailure()` in `cmd/glens/cmd/analyze.go` — preserve this
distinction when modifying test-failure detection logic.

## Skills

Specialised agent skills live in `.github/skills/` (Copilot CLI / Claude CLI format).
Each skill directory contains a `SKILL.md` with `name` and `description` frontmatter.
Each skill file is ≤ 100 lines and covers one focused topic.

| Skill directory | Purpose |
|---|---|
| `testing-strategy` | Testing philosophy, value-driven approach, naming, corner-case checklist |
| `testing-patterns` | Go-specific patterns: table-driven tests, mocks, helpers, assertions |
| `testing-integration` | Integration and end-to-end testing design and helper patterns |


@.github/instructions/testing-integration.instructions.md

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
- Integration tests use `test_specs/sample_api.json` — no network dependency.

## Documentation

Update docs only when:
- CLI interface changes.
- Setup process changes.
- New configuration options are added.
- Major architecture changes.

Do **not** update docs for internal refactoring, bug fixes, or performance improvements.

Per-module READMEs:
- `cmd/glens/README.md`
- `cmd/tools/demo/README.md`
- `cmd/tools/accuracy/README.md`
- `pkg/logging/README.md`

Root `README.md` links to every module README. `docs/` holds user guides and architecture diagrams.

## File Structure

```
go.work                          # workspace root — no go.mod here
pkg/
  logging/                       # module glens/pkg/logging
    logging.go                   # zerolog wrapper
    go.mod
    Makefile
    README.md
cmd/
  glens/                         # module glens/tools/glens
    main.go
    cmd/                         # cobra commands
    internal/
      ai/                        # AI provider clients
      generator/                 # test generation + execution
      github/                    # GitHub API client
      parser/                    # OpenAPI spec parser
      reporter/                  # report generation
    go.mod
    Makefile
    README.md
  tools/
    demo/                        # module glens/tools/demo
      internal/loader/ render/
      go.mod
      Makefile
      README.md
    accuracy/                    # module glens/tools/accuracy
      internal/analyze/ report/
      go.mod
      Makefile
      README.md
configs/                         # example configuration files
docs/                            # user guides and architecture diagrams
test_specs/                      # OpenAPI specs used in integration tests
```

## Common Tasks

**Working in a module:**

```bash
cd cmd/glens       # or cmd/tools/demo, cmd/tools/accuracy, pkg/logging
make all           # fmt-check + vet + lint + test (identical to CI)
make build
```

**Adding a CLI flag (cmd/glens):**

```go
cmd.Flags().String("flag-name", "default", "Clear description")
_ = viper.BindPFlag("config.key", cmd.Flags().Lookup("flag-name"))
```

**Adding a new tool module:**

1. `mkdir cmd/tools/<name> && cd cmd/tools/<name>`
2. `go mod init glens/tools/<name>`
3. Add `use ./cmd/tools/<name>` to `go.work`
4. Copy `Makefile` and `.pre-commit-config.yaml` from an existing tool
5. Add `.github/workflows/tool-<name>.yml`
6. Add binary to `.github/workflows/release.yml`

**Adding a config option:**

```yaml
# configs/config.example.yaml
new_option: "${ENV_VAR}"
```

## Before Any Change — Checklist

1. Is this the simplest solution?
2. Will this be easy to maintain in two years?
3. Can someone new understand this in five minutes?
4. Does this follow Go idioms?
5. Does this respect the module isolation rules (`pkg/` vs `internal/`)?

If any answer is "no", reconsider the approach.

## What to Avoid

- Features added "just in case"
- Imports that cross module boundaries
- `pkg/` packages importing from `internal/`
- Complex abstractions or inheritance hierarchies
- Reflection unless absolutely necessary
- Magic numbers (use named constants)
- Global state
- New dependencies without strong justification
- Premature optimization
