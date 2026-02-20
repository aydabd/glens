# 🤖 Glens

> OpenAPI Integration Test Generator with AI

Analyzes OpenAPI specs, generates integration tests using AI, and creates GitHub issues **only when tests fail**.

## ✨ Features

- **AI Test Generation**: GPT-4, Claude, Gemini, or local Ollama
- **Smart Issue Creation**: Only for real test failures, not infrastructure issues
- **Multiple AI Models**: Compare outputs from different models
- **Local LLM Support**: Free, private testing with Ollama
- **Detailed Reports**: Markdown, HTML, JSON formats

## 🚀 Quick Start

```bash
# Build
make build

# Run with Ollama (free, local)
make ollama-serve &
make run-ollama

# Or with OpenAI
export OPENAI_API_KEY="sk_xxx"
export GITHUB_TOKEN="ghp_xxx"
export GITHUB_REPOSITORY="owner/repo"
./build/glens analyze https://api.example.com/openapi.json --ai-models=gpt4
```

**Full guide**: See [docs/QUICKSTART.md](docs/QUICKSTART.md)

## �� Documentation

- **[QUICKSTART.md](docs/QUICKSTART.md)** - Installation, usage, examples, troubleshooting
- **[DEVELOPMENT.md](docs/DEVELOPMENT.md)** - Contributing, development setup, testing
- **[Architecture Diagrams](docs/diagrams/ARCHITECTURE.md)** - Visual system design and flows
- **[CLEANUP.md](docs/CLEANUP.md)** - Managing and cleaning up test issues

## 🎯 How It Works

```txt
OpenAPI Spec → AI Generation → Test Execution → Issue (if fail) → Report
```

See [architecture diagrams](docs/diagrams/ARCHITECTURE.md) for detailed flows.

## 🔑 Environment Variables

```bash
GITHUB_TOKEN        # Required for issue creation
GITHUB_REPOSITORY   # Required for issue creation (owner/repo)
OPENAI_API_KEY      # For GPT-4
ANTHROPIC_API_KEY   # For Claude (optional)
GOOGLE_API_KEY      # For Gemini (optional)
```

## 🔧 Common Commands

```bash
make build                     # Build binary
make test                      # Run tests with coverage
make test-integration          # Run integration tests
make run                       # Run with default spec
make run-ollama                # Run with local Ollama
make run-ollama-issues         # Run with Ollama and create issues
make cleanup-test-issues       # Preview issue cleanup (dry-run)
make cleanup-test-issues-confirm # Actually close test issues
make help                      # Show all commands
```

## 📝 Examples

```bash
# Basic usage
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json

# With issue creation
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4 \
  --github-repo=owner/repo \
  --create-issues

# Multiple AI models comparison
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4,ollama:codellama:7b-instruct

# Test specific endpoint
./build/glens analyze https://api.example.com/openapi.json \
  --op-id=getUserById \
  --create-issues=false
```

More examples in [docs/QUICKSTART.md](docs/QUICKSTART.md).

## 🧪 Testing & Accuracy

Test Glens' accuracy in generating integration tests:

```bash
# Run accuracy tests with mock AI (no API keys needed)
./scripts/test_accuracy.sh
```

This evaluates:

- ✅ OpenAPI parsing accuracy
- ✅ Endpoint coverage
- ✅ Test code generation quality
- ✅ Comprehensive reporting

See [ACCURACY_REPORT.md](ACCURACY_REPORT.md) for detailed analysis and
[accuracy_tests/ACCURACY_TESTING.md](accuracy_tests/ACCURACY_TESTING.md) for testing documentation.

## 🤝 Contributing

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for setup and guidelines.

**Philosophy**: Simple code, minimal docs, long-term maintainability.

## 📄 License

MIT License - See LICENSE file.
