# glens

> OpenAPI Integration Test Generator with AI — Go workspace monorepo

Analyzes OpenAPI specs, generates integration tests using AI, and creates GitHub
issues **only when tests fail**. All binaries are cross-platform and available
as standalone downloads — no Go toolchain required.

## Workspace layout

```text
go.work
├── pkg/logging           # module glens/pkg/logging    — generic zerolog wrapper
├── cmd/glens             # module glens/tools/glens    — main CLI
├── cmd/tools/demo        # module glens/tools/demo     — OpenAPI spec visualiser
└── cmd/tools/accuracy    # module glens/tools/accuracy — endpoint accuracy reporter
```

Each module is independently buildable, testable, and releasable. Any module can be moved to its own repository without changes.

## Module documentation

| Module | README | Description |
|--------|--------|-------------|
| `cmd/glens` | [cmd/glens/README.md](cmd/glens/README.md) | Main CLI — generate tests, create issues |
| `cmd/tools/demo` | [cmd/tools/demo/README.md](cmd/tools/demo/README.md) | Render an OpenAPI spec summary |
| `cmd/tools/accuracy` | [cmd/tools/accuracy/README.md](cmd/tools/accuracy/README.md) | Endpoint accuracy report |
| `pkg/logging` | [pkg/logging/README.md](pkg/logging/README.md) | Generic zerolog setup wrapper |

## Download binaries

Pre-built binaries for `linux/amd64`, `linux/arm64`, `darwin/amd64`,
`darwin/arm64`, and `windows/amd64` are on the
[releases page](https://github.com/aydabd/glens/releases).
Download and run — no dependencies.

## Quick start (glens CLI)

```bash
# Build from source
cd cmd/glens && make build

# Run (binary lands in ../../build/glens)
export OPENAI_API_KEY="sk_xxx"
export GITHUB_TOKEN="ghp_xxx"
../../build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4 --github-repo=owner/repo
```

Full usage: [cmd/glens/README.md](cmd/glens/README.md)

## Development

```bash
# Set up micromamba environment (optional — plain go also works)
make env

# Work inside a single module
cd cmd/glens
make all        # fmt-check + vet + lint + test (identical to CI)
make build
make test
```

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for full contributor guide.

## Documentation

- [docs/QUICKSTART.md](docs/QUICKSTART.md) — Getting started
- [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) — Contributor guide
- [docs/diagrams/architecture.md](docs/diagrams/architecture.md) — Architecture diagrams
- [docs/saas/SAAS_PLAN.md](docs/saas/SAAS_PLAN.md) — SaaS transformation plan (GCP)

## License

MIT License — see [LICENSE](LICENSE).
