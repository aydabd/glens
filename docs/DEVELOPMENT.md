# Glens Development Guide

> For contributors and developers working on Glens

## Workspace structure

```text
go.work                          # Go workspace root (no go.mod here)
├── pkg/logging/                 # module glens/pkg/logging
│   ├── logging.go
│   ├── Makefile
│   └── README.md
├── cmd/glens/                   # module glens/tools/glens
│   ├── main.go
│   ├── cmd/                     # CLI commands (root, analyze, cleanup, models)
│   ├── internal/                # ai, generator, github, parser, reporter
│   ├── Makefile
│   └── README.md
├── cmd/tools/demo/              # module glens/tools/demo
│   ├── internal/loader/ render/
│   ├── Makefile
│   └── README.md
└── cmd/tools/accuracy/          # module glens/tools/accuracy
    ├── internal/analyze/ report/
    ├── Makefile
    └── README.md
```

**`pkg/`** — generic reusable libraries; no `internal/` imports; independently versioned.  
**`cmd/*/internal/`** — module-private code; never imported across module boundaries.

## Setup

```bash
# Clone
git clone https://github.com/aydabd/glens
cd glens

# Create micromamba environment (optional — plain go also works)
make env
```

## Development workflow

Work inside the module you are changing:

```bash
cd cmd/glens          # or cmd/tools/demo, cmd/tools/accuracy, pkg/logging

make fmt              # format
make vet              # go vet
make lint             # golangci-lint
make test             # run tests with race detector
make build            # build binary
make all              # fmt-check + vet + lint + test  (same as CI)
```

`make all` in any module directory is identical to what CI runs — if it passes locally it passes in CI.

## Adding a new module

1. Create the directory: `cmd/tools/<name>/`
2. Run `go mod init glens/tools/<name>`
3. Add `use ./cmd/tools/<name>` to `go.work`
4. Copy `Makefile` from an existing tool module
5. Add a `.pre-commit-config.yaml` covering `go-fmt`, `go-vet`, `golangci-lint`
6. Add `.github/workflows/tool-<name>.yml` triggered on `cmd/tools/<name>/**`
7. Add the binary to `.github/workflows/release.yml`

## Package conventions

- `pkg/` — generic libraries; zero `internal/` imports; can be used by any Go project
- `cmd/*/internal/` — module-private; never import from another module
- All packages follow `gofmt` + `golangci-lint` standards
- Functions ≤ 50 lines; explicit error handling; no global state

## Key components (cmd/glens)

### Issue creation (`cmd/analyze.go`)

Issues are created **only** when tests compile and produce assertion failures.
Connection errors, compilation errors, and passing tests never produce issues.
`isRealTestFailure()` implements this distinction.

### AI layer (`internal/ai/`)

Every provider implements the `AIClient` interface:

```go
type AIClient interface {
    GenerateTest(ctx context.Context, endpoint *parser.Endpoint) (string, string, error)
}
```

Add a new provider by implementing the interface and registering it in `NewManager()`.

### Report generation (`internal/reporter/`)

`SuccessRate` is computed from actual pass/total counts — never hardcoded.

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

Integration tests that need a real spec use `test_specs/sample_api.json` with no network dependency.

## CI workflows

| Workflow | Trigger | What it runs |
|----------|---------|--------------|
| `pkg-logging.yml` | `pkg/logging/**` | `make all` + `go test` |
| `glens.yml` | `cmd/glens/**` | `make all` + `go test` |
| `api.yml` | `cmd/api/**` | `make all` + `go test` |
| `tool-demo.yml` | `cmd/tools/demo/**` | `make all` + `go test` |
| `tool-accuracy.yml` | `cmd/tools/accuracy/**` | `make all` + `go test` |
| `release-please.yml` | push to `main` | Release Please per-module versioning |

Each workflow is fully independent — a change in one module only triggers that module's CI.

## Release

Releases are automated by [Release Please](https://github.com/googleapis/release-please).
Merging to `main` triggers release PR creation. Merging a release PR creates a tag and
GitHub Release. See [docs/RELEASE_PLAN.md](RELEASE_PLAN.md) for the full strategy.

Individual module release workflows also support direct tag-based releases:

```bash
git tag pkg/logging/v0.2.0
git push origin pkg/logging/v0.2.0
```

## Conventional commits

All commits **must** follow [Conventional Commits](https://www.conventionalcommits.org/).
This is enforced by the `conventional-pre-commit` hook (install with `pre-commit install --hook-type commit-msg`).

```text
<type>[optional scope]: <description>
```

Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `revert`.

Use module path as scope: `feat(cmd/glens): add JSON output`.

Breaking changes: add `!` after type (`feat!:`) or add `BREAKING CHANGE:` footer.

## Code review checklist

- [ ] Functions ≤ 50 lines
- [ ] Explicit error handling (no `_` for errors)
- [ ] Table-driven tests for new functions
- [ ] `make all` passes in the changed module
- [ ] Docs updated only if CLI/config/setup changed
