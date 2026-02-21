# glens

> OpenAPI Integration Test Generator with AI

Analyzes OpenAPI specs, generates integration tests using AI, and creates GitHub issues **only when tests fail**.

Module: `glens/tools/glens`

## Features

- AI-powered test generation (GPT-4, Claude, Gemini, local Ollama)
- Issues created only for real spec violations — never for infrastructure errors
- Multi-model comparison reports
- Markdown, HTML, and JSON report formats

## Install

### Download binary (recommended)

Download the pre-built binary for your platform from the [releases page](https://github.com/aydabd/glens/releases) — no Go toolchain required.

### Build from source

```bash
cd cmd/glens
make build        # builds to ../../build/glens
```

Or with plain Go:

```bash
cd cmd/glens
go build -o ../../build/glens .
```

## Usage

```bash
# Analyze a spec (no issue creation)
./build/glens analyze https://api.example.com/openapi.json --create-issues=false

# With OpenAI and GitHub issue creation
export OPENAI_API_KEY="sk_xxx"
export GITHUB_TOKEN="ghp_xxx"
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4 \
  --github-repo=owner/repo

# Local Ollama (free, private)
./build/glens analyze https://api.example.com/openapi.json --ai-models=ollama

# Target one endpoint
./build/glens analyze https://api.example.com/openapi.json --op-id=getUserById
```

## Makefile targets

Run from this directory (`cmd/glens/`):

| Target | Description |
|--------|-------------|
| `make all` | fmt-check + vet + lint + test (same as CI) |
| `make fmt` | Format source |
| `make fmt-check` | Fail if source is unformatted |
| `make vet` | Run `go vet` |
| `make tidy` | Run `go mod tidy` |
| `make lint` | Run golangci-lint |
| `make test` | Run tests with race detector |
| `make build` | Build binary to `../../build/glens` |
| `make clean` | Remove build artifacts |

Micromamba is used automatically when available; plain `go` is used as fallback.

## Environment variables

| Variable | Required | Purpose |
|----------|----------|---------|
| `GITHUB_TOKEN` | For issue creation | GitHub authentication |
| `GITHUB_REPOSITORY` | For issue creation | Target repo (`owner/repo`) |
| `OPENAI_API_KEY` | For GPT-4 | OpenAI API access |
| `ANTHROPIC_API_KEY` | For Claude | Anthropic API access |
| `GOOGLE_API_KEY` | For Gemini | Google API access |

## Configuration

Create `configs/config.yaml`:

```yaml
github:
  token: "${GITHUB_TOKEN}"
  repository: "${GITHUB_REPOSITORY}"

ai_models:
  openai:
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4-turbo"
```

## Issue creation logic

Issues are created **only** when:
- Tests compile and run successfully
- Tests fail with assertion errors (spec violations)

Issues are **not** created for:
- Connection failures
- Test compilation errors
- Infrastructure problems

Key function: `isRealTestFailure()` in `cmd/analyze.go`.

## Module structure

```
cmd/glens/
├── main.go                 # Entry point
├── cmd/                    # CLI command definitions
│   ├── root.go             # Config, logging, root cobra command
│   ├── analyze.go          # Main analysis pipeline
│   ├── cleanup.go          # Issue cleanup command
│   └── models.go           # AI model management command
├── internal/               # Private implementation (never imported externally)
│   ├── ai/                 # AI provider clients
│   ├── generator/          # Test generation and execution
│   ├── github/             # GitHub API client
│   ├── parser/             # OpenAPI spec parser
│   └── reporter/           # Report generation
├── go.mod                  # Module: glens/tools/glens
├── Makefile
└── README.md
```
