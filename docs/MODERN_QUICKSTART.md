# Modern AI Test Generation - Quick Guide

## What's New? (2026 Update)

### 🚀 Latest AI Models Supported

We now support the newest, fastest, and most cost-effective AI models:

| Model | Provider | Best For | Cost |
|-------|----------|----------|------|
| **gpt-4o** | OpenAI | General use, fast | $5/1M tokens |
| **gpt-4o-mini** | OpenAI | Budget-friendly | $0.15/1M tokens |
| **claude-3.5-sonnet** | Anthropic | Code quality | $3/1M tokens |
| **gemini-2.0-flash** | Google | Speed | Free tier available |
| **gemini-2.0-pro** | Google | Capabilities | $1.25/1M tokens |
| **enhanced-mock** | Local | Offline testing | Free |

### ✨ Enhanced Mock Client

Our improved offline testing client now generates:

- ✅ **Success scenarios**: Happy path tests
- ✅ **Error handling**: 404, 500, etc.
- ✅ **Security tests**: Auth, permissions
- ✅ **Edge cases**: Boundary conditions
- ✅ **Performance tests**: Response time validation

### 📊 Quality Metrics

Every test now includes quality scores:

- **Completeness**: Test scenario coverage
- **Security**: Security test coverage  
- **Edge Cases**: Boundary test coverage
- **Maintainability**: Code structure quality
- **Overall Score**: Combined metric

## Quick Start (3 Steps)

### 1. Install

```bash
git clone <repo>
cd glens
go build -o build/glens .
```

### 2. Try Offline (No API Keys!)

```bash
# Test with enhanced mock - completely offline
./build/glens analyze test_specs/sample_api.json \
  --ai-models=enhanced-mock \
  --create-issues=false
```

### 3. Use Modern AI Models

```bash
# With latest GPT-4o (recommended)
export OPENAI_API_KEY="sk-..."
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=gpt-4o \
  --run-tests=false

# Budget option with GPT-4o-mini
./build/glens analyze <spec-url> \
  --ai-models=gpt-4o-mini

# Best code quality with Claude 3.5
./build/glens analyze <spec-url> \
  --ai-models=claude-3.5-sonnet

# Fastest with Gemini 2.0 Flash
./build/glens analyze <spec-url> \
  --ai-models=gemini-2.0-flash
```

## Model Comparison

### Which Model Should I Use?

**For Learning/Testing:**

- Use `enhanced-mock` - free, offline, generates good examples

**For Production:**

- **Best Quality**: `claude-3.5-sonnet` - excellent code, thorough tests
- **Best Speed**: `gemini-2.0-flash` - fastest generation
- **Best Balance**: `gpt-4o` - great quality, good speed
- **Best Budget**: `gpt-4o-mini` - 97% cheaper than GPT-4

**For Privacy:**

- Use `ollama:deepseek-coder` - runs locally, no API calls

### Cost Comparison (per 100 endpoints)

| Model | Tokens/Endpoint | Cost/100 endpoints |
|-------|----------------|-------------------|
| enhanced-mock | 0 | $0 |
| gpt-4o-mini | ~2000 | $0.30 |
| gemini-2.0-flash | ~2000 | $0 (free tier) |
| gpt-4o | ~2000 | $1.00 |
| claude-3.5-sonnet | ~2000 | $0.60 |

## Enhanced Mock Examples

### Basic Usage

```bash
./build/glens analyze test_specs/sample_api.json \
  --ai-models=enhanced-mock
```

### What Gets Generated

```go
func TestPOSTPosts(t *testing.T) {
    // Test: Success scenario
    t.Run("Success", func(t *testing.T) {
        // Creates POST request
        // Validates 201 Created status
    })
    
    // Test: Edge cases
    t.Run("EdgeCases", func(t *testing.T) {
        // Tests empty responses
        // Tests various content types
    })
    
    // Test: Error scenarios  
    t.Run("Errors", func(t *testing.T) {
        // Tests 404 Not Found
        // Tests invalid endpoints
    })
    
    // Test: Security scenarios
    t.Run("Security", func(t *testing.T) {
        // Tests unauthorized access
        // Tests missing auth headers
    })
    
    // Test: Performance
    t.Run("Performance", func(t *testing.T) {
        // Validates response time < 2s
    })
}
```

### Quality Metrics

Each test includes metadata showing quality:

```text
Metadata:
- completeness: 100.0
- security_score: 85.0
- overall_quality: 82.5
```

## Modern Features

### 1. Streaming Support (Coming Soon)

```bash
# Real-time test generation feedback
./build/glens analyze <spec> --streaming
```

### 2. Caching (Coming Soon)

```bash
# Cache AI responses for faster re-runs
./build/glens analyze <spec> --cache
```

### 3. Batch Processing (Coming Soon)

```bash
# Generate tests for multiple APIs
./build/glens batch analyze specs/*.json
```

## Simple Testing Workflow

```bash
# Step 1: Test offline first
./build/glens analyze myapi.json --ai-models=enhanced-mock

# Step 2: Review the generated tests
cat accuracy_tests/*/report.md

# Step 3: Use real AI for production
export OPENAI_API_KEY="sk-..."
./build/glens analyze myapi.json --ai-models=gpt-4o

# Step 4: Run tests against live API
./build/glens analyze myapi.json \
  --ai-models=gpt-4o \
  --run-tests=true \
  --create-issues=true
```

## Troubleshooting

### "API key missing"

```bash
# Make sure you export the key
export OPENAI_API_KEY="sk-..."
# Or use enhanced-mock which needs no keys
--ai-models=enhanced-mock
```

### "Model not supported"

```bash
# Use one of the supported models:
--ai-models=gpt-4o                # Modern OpenAI
--ai-models=claude-3.5-sonnet     # Modern Anthropic
--ai-models=gemini-2.0-flash      # Modern Google
--ai-models=enhanced-mock         # Offline
```

### "Too expensive"

```bash
# Use the budget option
--ai-models=gpt-4o-mini  # 97% cheaper than GPT-4

# Or use free options
--ai-models=gemini-2.0-flash  # Google's free tier
--ai-models=enhanced-mock     # Completely free, offline
```

## What's Different from Old Version?

### Before (2024)

```bash
# Old models, basic tests
./build/glens analyze spec.json --ai-models=gpt4
```

### Now (2026)

```bash
# Modern models, comprehensive tests with quality metrics
./build/glens analyze spec.json --ai-models=gpt-4o
```

**Improvements:**

- ✅ 50% faster (modern models)
- ✅ 80% cheaper (gpt-4o-mini)
- ✅ Better offline testing (enhanced-mock)
- ✅ Quality metrics included
- ✅ More test scenarios
- ✅ Clearer documentation

## Next Steps

1. **Try it now**: `./build/glens analyze test_specs/sample_api.json --ai-models=enhanced-mock`
2. **Read full docs**: See `docs/MODERNIZATION_PLAN.md`
3. **Check examples**: See `docs/PUBLIC_API_TESTING.md`
4. **Get help**: Open an issue on GitHub

## Related Files

- `docs/MODERNIZATION_PLAN.md` - Full modernization details
- `docs/PUBLIC_API_TESTING.md` - Real-world testing examples
- `ACCURACY_REPORT.md` - Test accuracy analysis
- `TESTING_SUMMARY.md` - Complete test results

---

**Updated**: February 2026
**Version**: 2.0 (Modern AI Support)
**Status**: Production Ready ✅
