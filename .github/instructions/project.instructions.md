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

## Skills

Specialised agent skills live in `.github/skills/`.
Each skill directory contains a `SKILL.md` (≤ 100 lines) covering one focused topic.

| Skill directory | Purpose |
|---|---|
| `go-patterns` | Architecture patterns, error handling, critical logic |
| `testing-strategy` | Testing philosophy, value-driven approach, naming |
| `testing-patterns` | Table-driven tests, mocks, setup functions, assertions |
| `testing-integration` | Integration and E2E testing design, fixture patterns |
| `common-tasks` | Build commands, adding modules/flags, change checklist |

## Documentation

Update docs only when:

- CLI interface changes.
- Setup process changes.
- New configuration options are added.
- Major architecture changes.

Do **not** update docs for internal refactoring, bug fixes, or performance improvements.

Per-module READMEs: `cmd/glens/`, `cmd/tools/demo/`, `cmd/tools/accuracy/`, `pkg/logging/`.
Root `README.md` links to every module README. `docs/` holds user guides.

## Mandatory Linting — Before Every Commit

**Always** run linters on every changed file before committing or pushing. This is
non-negotiable — CI will reject unlinted changes.

### Go modules

For each Go module you modified, run `make lint` (or `make all` for the full suite):

```bash
cd cmd/glens && make lint        # golangci-lint
cd cmd/tools/demo && make lint
cd cmd/tools/accuracy && make lint
cd pkg/logging && make lint
```

### Markdown files

If you created or edited any `.md` file, run the markdown linter from the repo root:

```bash
npx markdownlint-cli '**/*.md'
```

### YAML / workflow files

If you created or edited any `.yml` or `.yaml` file, validate syntax:

```bash
python3 -c "import yaml; yaml.safe_load(open('path/to/file.yml'))"
```

### Pre-commit hooks

If `.pre-commit-config.yaml` is set up, run all hooks on changed files:

```bash
pre-commit run --files <changed-files>
```

**Never skip linting.** If a linter reports errors, fix them before pushing.

## Conventional Commits

All commits **must** follow the [Conventional Commits](https://www.conventionalcommits.org/)
specification. This is enforced by the `conventional-pre-commit` hook and required by
Release Please for automatic version bumping.

### Format

```text
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

### Allowed types

| Type | Version bump | Use when |
|------|-------------|----------|
| `feat` | minor | Adding new functionality |
| `feat!` | **major** | Breaking change in a feature |
| `fix` | patch | Bug fix |
| `docs` | — | Documentation only |
| `style` | — | Formatting, no code change |
| `refactor` | patch | Code restructuring, no feature/fix |
| `perf` | patch | Performance improvement |
| `test` | — | Adding or correcting tests |
| `build` | patch | Build system or dependency changes |
| `ci` | — | CI configuration changes |
| `chore` | — | Maintenance tasks |
| `revert` | patch | Reverting a previous commit |

A `BREAKING CHANGE:` footer in any commit type triggers a **major** bump.

### Scoping to modules

Use the module directory as scope:

```text
feat(cmd/glens): add --output flag for JSON reports
fix(pkg/logging): handle nil logger gracefully
feat(cmd/api)!: redesign authentication endpoint
```

### Examples

```text
feat: add streaming support for AI responses
fix(cmd/glens): correct exit code on test failure
docs: update quickstart with new flags
chore(ci): upgrade Go version to 1.25
refactor(pkg/logging): simplify level configuration
feat(cmd/api)!: change /analyze endpoint response format

BREAKING CHANGE: response field "result" renamed to "analysis"
```

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
