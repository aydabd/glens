# Glens Quick Start Guide

> OpenAPI Integration Test Generator with AI - Simple, Fast, Effective

## What It Does

Analyzes OpenAPI specs, generates tests using AI, runs them, and creates GitHub issues **only when tests fail**.

**Visual Guide**: See [Architecture Diagrams](diagrams/ARCHITECTURE.md) for flow diagrams and system design.

## Install

```bash
git clone <repo>
cd glens
make build
```

## Basic Usage

```bash
# Setup (required once)
export GITHUB_TOKEN="ghp_your_token"
export GITHUB_REPOSITORY="owner/repo"

# Run with local Ollama (free, private)
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=ollama

# Run with OpenAI GPT-4
export OPENAI_API_KEY="sk_your_key"
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4
```

## Common Commands

```bash
# Test specific endpoint
make test-endpoint OP_ID=getUserById

# Test without creating issues (dry run)
./build/glens analyze URL --create-issues=false

# Use specific AI model
./build/glens analyze URL --ai-models=ollama:codellama:7b-instruct

# Multiple models
./build/glens analyze URL --ai-models=gpt4,ollama
```

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

## GitHub Issues

Issues are created **ONLY** when:

- ✅ Tests fail with assertion errors
- ✅ Response doesn't match OpenAPI spec
- ✅ Status codes don't match

Issues are **NOT** created when:

- ❌ Connection fails (infrastructure issue)
- ❌ Test compilation fails (setup issue)
- ❌ All tests pass

## Local AI (Ollama)

Free, private, no API costs:

```bash
# Setup Ollama (in micromamba environment)
make ollama-pull-codellama

# Start server
make ollama-serve

# Run tests
make run-ollama
```

## Makefile Commands

```bash
make build              # Build binary
make test               # Run tests with coverage
make run                # Run with default spec
make run-ollama         # Run with Ollama
make test-endpoint      # Test specific endpoint (set OP_ID=...)
```

## Environment Variables

| Variable            | Required   | Purpose                  |
| ------------------- | ---------- | ------------------------ |
| `GITHUB_TOKEN`      | For issues | GitHub authentication    |
| `GITHUB_REPOSITORY` | For issues | Target repo (owner/repo) |
| `OPENAI_API_KEY`    | For GPT    | OpenAI API access        |

## Flags

```bash
--ai-models strings       # AI models to use (default: gpt4)
--github-repo string      # Target repository
--create-issues           # Create issues on failures (default: true)
--run-tests              # Execute tests (default: true)
--op-id string           # Target specific endpoint
--output string          # Report file (default: reports/report.md)
--debug                  # Enable debug logging
```

## Examples

### Example 1: Quick Test

```bash
make run
```

### Example 2: Production Run

```bash
export GITHUB_TOKEN=$(gh auth token)
export GITHUB_REPOSITORY="myorg/myapi"
export OPENAI_API_KEY="sk_xxx"

./build/glens analyze https://api.production.com/openapi.json \
  --ai-models=gpt4 \
  --create-issues
```

### Example 3: Local Development

```bash
make ollama-serve &
./build/glens analyze ./specs/local-api.yaml \
  --ai-models=ollama \
  --create-issues=false
```

## Troubleshooting

### GitHub token required

```bash
export GITHUB_TOKEN=$(gh auth token)
```

### Ollama not responding

```bash
make ollama-serve
```

### No tests generated

- Check OpenAPI spec is valid
- Verify AI model is accessible
- Use `--debug` flag for details

## Report Output

Check `reports/report.md` for:

- Endpoint analysis
- Test results
- AI model comparison
- Issue numbers (if created)

**Note:** All reports are generated in the `reports/` directory, which is gitignored.

## Get Help

```bash
./build/glens analyze --help
make help
```

## Next Steps

- Review generated report: `cat reports/report.md`
- Check created issues: `gh issue list --repo <owner>/<repo>`
- Customize: Edit `configs/config.yaml`
- For development: See `docs/DEVELOPMENT.md`
