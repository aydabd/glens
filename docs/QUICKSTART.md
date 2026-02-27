# Glens Quick Start Guide

> OpenAPI Integration Test Generator with AI

## What It Does

Parses OpenAPI specs, generates integration tests using AI, runs them, and creates GitHub issues **only when tests fail**.

## Install

### Option 1 — Download binary (no Go required)

Download the binary for your platform from the [releases page](https://github.com/aydabd/glens/releases) and run it directly.

### Option 2 — Build from source

```bash
git clone https://github.com/aydabd/glens
cd glens/cmd/glens
make build          # binary → ../../build/glens
```

## Basic Usage

```bash
# Required for issue creation
export GITHUB_TOKEN="ghp_your_token"
export GITHUB_REPOSITORY="owner/repo"

# Run with local Ollama (free, private)
./build/glens analyze https://api.example.com/openapi.json --ai-models=ollama

# Run with OpenAI GPT-4
export OPENAI_API_KEY="sk_your_key"
./build/glens analyze https://api.example.com/openapi.json --ai-models=gpt4

# Dry run (no issue creation)
./build/glens analyze https://api.example.com/openapi.json --create-issues=false
```

## Configuration

Create `configs/config.yaml` inside `cmd/glens/`:

```yaml
github:
  token: "${GITHUB_TOKEN}"
  repository: "${GITHUB_REPOSITORY}"

ai_models:
  openai:
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4-turbo"
```

## GitHub Issue Logic

Issues are created **only** when:

- Tests compile and run successfully
- Tests fail with assertion errors (spec violations)

Issues are **not** created when:

- Connection to the server fails
- Test compilation fails
- All tests pass

## Makefile Commands (run from `cmd/glens/`)

```bash
make all        # fmt-check + vet + lint + test (same as CI)
make build      # Build binary
make test       # Run tests with race detector
make fmt        # Format code
make lint       # Run golangci-lint
make clean      # Remove build artifacts
```

## Flags

```text
--ai-models strings    AI models to use (default: gpt4)
--github-repo string   Target repository (owner/repo)
--create-issues        Create issues on failures (default: true)
--run-tests            Execute tests (default: true)
--op-id string         Target a specific endpoint by operationId
--output string        Report file path (default: reports/report.md)
--debug                Enable debug logging
```

## Environment Variables

| Variable | Required | Purpose |
|----------|----------|---------|
| `GITHUB_TOKEN` | For issue creation | GitHub authentication |
| `GITHUB_REPOSITORY` | For issue creation | Target repo (`owner/repo`) |
| `OPENAI_API_KEY` | For GPT-4 | OpenAI API |
| `ANTHROPIC_API_KEY` | For Claude | Anthropic API |
| `GOOGLE_API_KEY` | For Gemini | Google API |

## Tools

Two standalone utility binaries are also available:

- **`glens-demo`** — renders an OpenAPI spec as a human-readable summary. See [cmd/tools/demo/README.md](../cmd/tools/demo/README.md).
- **`glens-accuracy`** — produces an endpoint accuracy report for one or more specs. See [cmd/tools/accuracy/README.md](../cmd/tools/accuracy/README.md).

## Troubleshooting

### GitHub token missing

```bash
export GITHUB_TOKEN=$(gh auth token)
```

### No tests generated

- Confirm the OpenAPI spec URL is reachable
- Verify the AI model is accessible
- Add `--debug` for detailed logs

### Reports

Reports are written to `reports/report.md` (gitignored). Open with any Markdown viewer.
