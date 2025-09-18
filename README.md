# 🤖 Glens

A powerful GitHub Actions workflow system that analyzes OpenAPI specifications and automatically
generates comprehensive integration tests using AI models.
Starts with OpenAI GPT-4 by default, with easy expansion to multiple AI models (Anthropic Claude Sonnet, Google Gemini Flash Pro).

## ✨ Features

- **AI Test Generation**: Generate tests using OpenAI GPT-4 (default), with support for Anthropic Claude Sonnet and Google Gemini Flash Pro
- **GitHub Integration**: Automatically create GitHub issues for each endpoint with detailed test requirements
- **Comprehensive Testing**: Generate tests covering happy path, error handling, boundary testing, and security validation
- **Multiple Frameworks**: Support for testify, Ginkgo, and standard Go testing frameworks
- **Detailed Reporting**: Generate comprehensive reports in Markdown, HTML, and JSON formats
- **Performance Analysis**: Compare AI model performance and get actionable recommendations
- **GitHub Actions Ready**: Complete workflow automation with artifact management and notifications

## 🚀 Quick Start

### 1. Setup Repository

1. Copy the `glens` directory to your repository
2. Copy the GitHub Actions workflow to `.github/workflows/`
3. Configure the required secrets in your GitHub repository

### 2. Configure Secrets

Add these secrets to your GitHub repository (`Settings` → `Secrets and variables` → `Actions`):

```bash
# Required for AI model access (start with OpenAI only)
OPENAI_API_KEY=your_openai_api_key

# Optional: Add these when you want to test multiple AI models
# ANTHROPIC_API_KEY=your_anthropic_api_key
# GOOGLE_API_KEY=your_google_api_key
# GOOGLE_PROJECT_ID=your_google_project_id

# Required for GitHub integration
GITHUB_TOKEN=automatically_provided_by_github

# Optional for notifications
SLACK_WEBHOOK_URL=your_slack_webhook_url
```

### 3. Run the Workflow

#### Manual Trigger

1. Go to `Actions` → `OpenAPI Integration Test Generator`
2. Click `Run workflow`
3. Provide your OpenAPI specification URL
4. Select AI models and configuration options
5. Run the workflow

#### Automatic Triggers

- **Scheduled**: Runs daily at 2 AM UTC (configurable)
- **Push Events**: Triggers when OpenAPI specs in `openapi-specs/` directory change

## 📋 Configuration

### Basic Usage

```bash
# Analyze a public OpenAPI spec
./glens analyze https://petstore3.swagger.io/api/v3/openapi.json

# Use specific AI models (cloud)
./glens analyze --ai-models gpt4,sonnet4 https://api.example.com/openapi.json

# Use local LLM via Ollama
./glens analyze --ai-models ollama https://api.example.com/openapi.json

# Use specific Ollama model
./glens analyze --ai-models ollama:codellama:7b-instruct https://api.example.com/openapi.json

# Mix local and cloud models
./glens analyze --ai-models gpt4,ollama:deepseek-coder:6.7b-instruct https://api.example.com/openapi.json

# Create GitHub issues
./glens analyze --github-repo owner/repo --create-issues https://api.example.com/openapi.json

# Generate HTML report
./glens analyze --output report.html https://api.example.com/openapi.json
```

### Configuration File

Create a `config.yaml` file for advanced configuration:

```yaml
# AI Model Configuration
ai_models:
  openai:
    api_key: "${OPENAI_API_KEY}"
    model: "gpt-4-turbo"
    timeout: "60s"
    max_tokens: 4000

  anthropic:
    api_key: "${ANTHROPIC_API_KEY}"
    model: "claude-3-sonnet-20240229"
    timeout: "60s"
    max_tokens: 4000

  google:
    api_key: "${GOOGLE_API_KEY}"
    project_id: "${GOOGLE_PROJECT_ID}"
    model: "gemini-1.5-flash"
    timeout: "60s"
    max_tokens: 4000

  # Local LLM via Ollama
  ollama:
    base_url: "http://localhost:11434"
    model: "codellama:7b-instruct"
    timeout: "300s"
    temperature: 0.1
    max_tokens: 4000

# GitHub Configuration
github:
  token: "${GITHUB_TOKEN}"
  repository: "owner/repo"
  create_issues: true
  issue_labels:
    - "integration-test"
    - "ai-generated"
    - "openapi"

# Test Configuration
test_generation:
  framework: "testify"
  timeout: "30s"
  include_security_tests: true
  include_boundary_tests: true
  include_error_handling: true

# Reporting Configuration
reporting:
  output_format: "markdown"
  include_prompt_details: true
  include_execution_logs: true
  compare_models: true
```

## 🎯 Workflow Inputs

| Input            | Description                                      | Default                  | Required |
| ---------------- | ------------------------------------------------ | ------------------------ | -------- |
| `openapi_url`    | OpenAPI specification URL or file path           | -                        | ✅       |
| `ai_models`      | AI models to use (comma-separated)               | `gpt4,sonnet4,flash-pro` | ❌       |
| `github_repo`    | Target repository for issues                     | Current repo             | ❌       |
| `test_framework` | Test framework (`testify`, `ginkgo`, `standard`) | `testify`                | ❌       |
| `create_issues`  | Create GitHub issues for endpoints               | `true`                   | ❌       |
| `run_tests`      | Execute generated tests                          | `true`                   | ❌       |
| `output_format`  | Report format (`markdown`, `html`, `json`)       | `markdown`               | ❌       |

## 📊 Generated Reports

The tool generates comprehensive reports including:

### Executive Summary

- Total endpoints analyzed
- Tests generated and executed
- Success rates and performance metrics
- Overall health score

### AI Model Comparison

- Performance rankings by quality, coverage, and reliability
- Strengths and weaknesses analysis
- Token usage and cost analysis
- Recommendations for model selection

### Detailed Results

- Per-endpoint test results
- Generated test code
- Execution logs and errors
- Security and boundary test coverage

### Actionable Recommendations

- Code quality improvements
- Performance optimizations
- Security enhancements
- Best practices guidance

## 🏠 Local LLM Setup (Ollama)

### Quick Start with Ollama

For cost-effective and private local AI model usage:

```bash
# 1. Install Ollama
curl -fsSL https://ollama.ai/install.sh | sh

# 2. Start Ollama service
ollama serve

# 3. Download recommended coding models
ollama pull codellama:7b-instruct        # Meta's CodeLlama - great for code generation
ollama pull deepseek-coder:6.7b-instruct # Excellent performance for coding tasks
ollama pull qwen2.5-coder:7b-instruct    # Latest high-performance coding model

# 4. Check available models
./glens models list

# 5. Test with local model
./glens analyze --ai-models ollama https://petstore3.swagger.io/api/v3/openapi.json
```

### Model Management

```bash
# List all available models (local + cloud)
./glens models list

# Check provider status
./glens models status

# Ollama-specific commands
./glens models ollama list
./glens models ollama status

# Download additional models
ollama pull llama3.1:8b-instruct         # General purpose with good coding ability
ollama pull starcoder2:7b                # Specialized code completion model
```

### Recommended Ollama Models for Code Generation

| Model                          | Size   | Best For                     | Performance |
| ------------------------------ | ------ | ---------------------------- | ----------- |
| `codellama:7b-instruct`        | ~3.8GB | General code generation      | ⭐⭐⭐⭐    |
| `deepseek-coder:6.7b-instruct` | ~3.7GB | High-quality code generation | ⭐⭐⭐⭐⭐  |
| `qwen2.5-coder:7b-instruct`    | ~4.1GB | Latest, best performance     | ⭐⭐⭐⭐⭐  |
| `starcoder2:7b`                | ~4.0GB | Code completion              | ⭐⭐⭐      |

### Benefits of Local LLMs

✅ **Cost-effective** - No API costs
✅ **Privacy** - Code stays on your machine
✅ **Offline** - Works without internet
✅ **Fast** - No network latency
✅ **Customizable** - Fine-tune models for your needs

## 🔧 Development

### Prerequisites

- [Micromamba](https://mamba.readthedocs.io/en/latest/installation.html) (recommended) or Go 1.25+
- Git
- Access to AI model APIs (optional with Ollama)

### Quick Setup with Micromamba

```bash
# Clone the repository
git clone <repository>
cd glens

# Run the automated setup script
./setup.sh

# Or manually with make
make setup

# Activate the environment
micromamba activate glens

# Run with example
make run
```

### Manual Development Setup

```bash
# Clone and setup
git clone <repository>
cd glens

# Create isolated environment
make env

# Install dependencies and build
make all

# Run tests
make test

# Run with local config
make run-local
```

### Available Make Commands

```bash
make help        # Show all available commands
make env         # Create/update micromamba environment
make shell       # Enter the environment shell
make build       # Build the binary
make test        # Run tests with coverage
make lint        # Run code linters
make fmt         # Format Go code
make clean       # Clean build artifacts
make setup       # Full development environment setup
```

### Project Structure

```text
glens/
├── cmd/                    # CLI commands
│   ├── root.go            # Root command and configuration
│   └── analyze.go         # Main analyze command
├── pkg/
│   ├── parser/            # OpenAPI specification parsing
│   ├── github/            # GitHub API integration
│   ├── ai/                # AI model clients
│   ├── generator/         # Test generation and execution
│   └── reporter/          # Report generation
├── configs/               # Configuration templates
├── templates/             # Issue and test templates
└── .github/workflows/     # GitHub Actions workflows
```

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Run the test suite
6. Submit a pull request

### Code Style

- Follow Go formatting conventions (`gofmt`)
- Add comprehensive tests
- Include documentation for new features
- Use conventional commit messages

## 📖 Examples

### Example 1: Basic Analysis

```bash
./glens analyze https://petstore3.swagger.io/api/v3/openapi.json
```

### Example 2: Full Integration with GitHub

```bash
./glens analyze \
  --ai-models gpt4,sonnet4,flash-pro \
  --github-repo myorg/myapi \
  --test-framework testify \
  --create-issues \
  --run-tests \
  --output report.html \
  https://api.myservice.com/openapi.json
```

### Example 3: Configuration File

```bash
./glens analyze --config production.yaml ./specs/api-v2.yaml
```

## 🔒 Security Considerations

- **API Keys**: Store all API keys as GitHub secrets, never in code
- **Permissions**: Use minimal required permissions for GitHub tokens
- **Rate Limiting**: Built-in rate limiting for AI model APIs
- **Validation**: Input validation for OpenAPI specifications
- **Secrets**: Automatic secret detection and prevention in generated code

## 📈 Performance

- **Parallel Processing**: AI models run in parallel for faster results
- **Caching**: Intelligent caching of API responses
- **Optimization**: Configurable timeouts and retry logic
- **Resource Management**: Efficient memory and CPU usage

## 🐛 Troubleshooting

### Common Issues

#### "API key missing" error

- Ensure all required API keys are set as GitHub secrets
- Check secret names match exactly: `OPENAI_API_KEY`, `ANTHROPIC_API_KEY`, etc.

#### "OpenAPI spec not accessible" error

- Verify the OpenAPI URL is publicly accessible
- Check if the specification is valid JSON/YAML
- Ensure proper authentication if required

#### "Test execution failed" error

- Check Go version compatibility (requires Go 1.21+)
- Verify test framework dependencies
- Review generated test code for syntax errors

#### GitHub integration issues

- Ensure `GITHUB_TOKEN` has required permissions
- Check repository name format: `owner/repo`
- Verify repository access permissions

### Debug Mode

Enable debug logging for troubleshooting:

```bash
./glens analyze --debug --log-format json https://api.example.com/openapi.json
```

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Cobra](https://github.com/spf13/cobra) for CLI framework
- [Viper](https://github.com/spf13/viper) for configuration management
- [Zerolog](https://github.com/rs/zerolog) for structured logging
- [Testify](https://github.com/stretchr/testify) for testing framework
- OpenAI, Anthropic, and Google for AI model APIs

## 📞 Support

- 📚 [Documentation](https://github.com/your-org/glens/wiki)
- 🐛 [Issues](https://github.com/your-org/glens/issues)
- 💬 [Discussions](https://github.com/your-org/glens/discussions)
- 📧 [Email Support](mailto:support@yourorg.com)
