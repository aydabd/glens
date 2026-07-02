---
name: common-tasks
description: >
  Common development tasks for Glens. Use this when building, adding modules,
  CLI flags, config options, or running the change checklist.
---

# Skill: Common Tasks

## Working in a Module

```bash
cd cmd/glens       # or cmd/tools/demo, cmd/tools/accuracy, pkg/logging
make all           # fmt-check + vet + lint + test (identical to CI)
make build
```

## Adding a CLI Flag (cmd/glens)

```go
cmd.Flags().String("flag-name", "default", "Clear description")
_ = viper.BindPFlag("config.key", cmd.Flags().Lookup("flag-name"))
```

## Adding a New Tool Module

1. `mkdir cmd/tools/<name> && cd cmd/tools/<name>`
2. `go mod init glens/tools/<name>`
3. Add `use ./cmd/tools/<name>` to `go.work`
4. Copy `Makefile` and `.pre-commit-config.yaml` from an existing tool.
5. Add `.github/workflows/tool-<name>.yml`.
6. Add binary to `.github/workflows/release.yml`.

## Adding a Config Option

```yaml
# configs/config.example.yaml
new_option: "${ENV_VAR}"
```

## File Structure

```text
go.work                          # workspace root
pkg/logging/                     # generic zerolog wrapper
cmd/glens/                       # main CLI (cobra commands, internal/ packages)
cmd/tools/demo/                  # OpenAPI spec visualiser
cmd/tools/accuracy/              # endpoint accuracy reporter
configs/                         # example configuration files
docs/                            # user guides and architecture diagrams
test_specs/                      # OpenAPI specs used in integration tests
```

## Before Any Change — Checklist

1. Is this the simplest solution?
2. Will this be easy to maintain in two years?
3. Can someone new understand this in five minutes?
4. Does this follow Go idioms?
5. Does this respect the module isolation rules (`pkg/` vs `internal/`)?

If any answer is "no", reconsider the approach.
