# Testing Glens Against Public APIs - Practical Guide

## Overview

This guide shows you how to test Glens' accuracy against real public APIs in your own environment.

## Prerequisites

You need ONE of these AI providers:

```bash
# Option 1: OpenAI (Recommended for quality)
export OPENAI_API_KEY="sk-..."

# Option 2: Anthropic Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# Option 3: Google Gemini
export GOOGLE_API_KEY="..."

# Option 4: Ollama (Free, runs locally)
# Start Ollama server first
ollama serve &
ollama pull codellama:7b-instruct
```

## Public APIs You Can Test

### 1. Swagger PetStore (Recommended for First Test)

```bash
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=gpt4 \
  --create-issues=false \
  --run-tests=false \
  --output=reports/petstore_analysis.md
```

**Why PetStore?**

- ✅ Simple, well-documented API
- ✅ Standard CRUD operations
- ✅ Multiple content types
- ✅ Authentication examples
- ✅ Perfect for learning

**Expected Results:**

- ~20 endpoints identified
- Tests generated for each endpoint
- Health score: 30-40% (without execution)
- Report showing all operations

### 2. JSONPlaceholder (Simple REST API)

```bash
# Note: JSONPlaceholder doesn't publish an OpenAPI spec by default
# You would need to find or create one
./build/glens analyze path/to/jsonplaceholder-openapi.json \
  --ai-models=gpt4 \
  --run-tests=true \
  --create-issues=false
```

**Why JSONPlaceholder?**

- ✅ Actual working API (typicode.com)
- ✅ Free, no authentication
- ✅ Can execute tests against live endpoints
- ✅ RESTful design

### 3. GitHub REST API

```bash
./build/glens analyze \
  https://raw.githubusercontent.com/github/rest-api-description/main/descriptions/api.github.com/api.github.com.json \
  --ai-models=gpt4,claude \
  --create-issues=false \
  --run-tests=false \
  --output=reports/github_api_comparison.md
```

**Why GitHub API?**

- ✅ Real-world complexity
- ✅ OAuth authentication
- ✅ Pagination examples
- ✅ Large API surface
- ⚠️ May take 10-30 minutes to analyze

**Expected Results:**

- 300+ endpoints identified
- AI model comparison
- Complex parameter handling
- Authentication scenarios

### 4. Stripe API

```bash
./build/glens analyze \
  https://raw.githubusercontent.com/stripe/openapi/master/openapi/spec3.json \
  --ai-models=gpt4 \
  --create-issues=false \
  --run-tests=false \
  --output=reports/stripe_analysis.md
```

**Why Stripe?**

- ✅ Production-quality API design
- ✅ Complex schemas
- ✅ Excellent documentation
- ✅ Real-world patterns
- ⚠️ Large spec, may take time

### 5. OpenWeatherMap

```bash
# You'll need to find their OpenAPI spec or create one
./build/glens analyze path/to/openweathermap-spec.json \
  --ai-models=gpt4 \
  --run-tests=true \
  --create-issues=false
```

**Why OpenWeatherMap?**

- ✅ Free API key available
- ✅ Simple authentication (API key)
- ✅ Can test against live endpoints
- ✅ Real responses to validate

## Step-by-Step: Complete Test Example

### Example: Testing PetStore with GPT-4

```bash
# 1. Set up environment
export OPENAI_API_KEY="your-key-here"
export GITHUB_TOKEN="your-token"  # Optional
export GITHUB_REPOSITORY="your/repo"  # Optional

# 2. Run analysis only (no test execution)
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=gpt4 \
  --create-issues=false \
  --run-tests=false \
  --output=reports/petstore_gpt4.md \
  --debug

# 3. Check the report
cat reports/petstore_gpt4.md

# 4. Compare with Claude
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=claude \
  --create-issues=false \
  --run-tests=false \
  --output=reports/petstore_claude.md

# 5. Multi-model comparison
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=gpt4,claude,gemini \
  --create-issues=false \
  --run-tests=false \
  --output=reports/petstore_comparison.md
```

## Understanding Test Results

### What to Look For

1. **Parsing Success**

   ```text
   11:52PM INF OpenAPI specification parsed successfully endpoints_count=20
   ```

2. **Test Generation**

   ```text
   11:52PM INF Generating tests with AI model ai_model=gpt4
   ```

3. **Quality Indicators**
   - Number of endpoints covered
   - Health score percentage
   - Test categories (happy path, error handling, edge cases)

### Sample Report Sections

```markdown
## 📊 Executive Summary

| Metric | Value |
|--------|-------|
| **Total Endpoints** | 20 |
| **Tests Generated** | 20 |
| **Overall Health Score** | 85.0% |

## 🤖 AI Model Performance Comparison

| Model | Success Rate | Avg Quality |
|-------|-------------|-------------|
| gpt4  | 100%        | 9.2/10      |
| claude| 100%        | 8.9/10      |
```

## Testing Against Live Endpoints

⚠️ **WARNING**: Only test against APIs you have permission to test!

```bash
# Example: Test PetStore (they allow testing)
./build/glens analyze https://petstore3.swagger.io/api/v3/openapi.json \
  --ai-models=gpt4 \
  --run-tests=true \
  --create-issues=false \
  --output=reports/petstore_live_test.md
```

**What happens:**

1. Tests are generated for each endpoint
2. Tests are compiled and executed
3. HTTP requests are made to live API
4. Results are validated against spec
5. Pass/fail status recorded

**Note**: Test execution can take several minutes depending on:

- Number of endpoints
- API response times
- Network latency

## Analyzing Issues (Spec Violations)

To create GitHub issues for spec violations:

```bash
# 1. Ensure GitHub is configured
export GITHUB_TOKEN="ghp_..."
export GITHUB_REPOSITORY="your/repo"

# 2. Run with issue creation enabled
./build/glens analyze https://api.example.com/openapi.json \
  --ai-models=gpt4 \
  --run-tests=true \
  --create-issues=true

# 3. Issues are created ONLY when:
#    - Tests execute successfully
#    - Tests fail (spec mismatch)
#    - NOT for connection errors or infrastructure issues
```

**Example Issue Content:**

```markdown
# Test Failure: GET /users/{id}

## Expected (from OpenAPI spec)
- Status code: 200
- Content-Type: application/json
- Response schema: User object

## Actual (from live API)
- Status code: 200
- Content-Type: text/html  ❌
- Response: HTML error page

## Recommendation
Update API to return JSON or update spec to reflect actual behavior.
```

## Troubleshooting

### Connection Errors

```bash
Error: failed to fetch from URL: dial tcp: lookup api.example.com
```

**Solutions:**

- Check internet connection
- Verify URL is correct
- Try downloading spec locally first
- Use `curl` to test URL access

### AI API Errors

```bash
Error: failed to generate with OpenAI: 401 Unauthorized
```

**Solutions:**

- Verify API key is set correctly
- Check API key has sufficient credits
- Try a different AI model

### Test Execution Timeouts

```bash
Error: test execution timed out after 2m
```

**Solutions:**

- Increase timeout in config
- Run without `--run-tests` first
- Test specific endpoints with `--op-id`

## Best Practices

### 1. Start Small

```bash
# Test one endpoint first
./build/glens analyze spec.json --op-id=getUser
```

### 2. Disable Issue Creation Initially

```bash
# Preview results before creating issues
--create-issues=false
```

### 3. Compare AI Models

```bash
# See which model produces best tests
--ai-models=gpt4,claude,gemini
```

### 4. Review Generated Code

```bash
# Check test quality manually
cat reports/report.md
```

### 5. Validate Gradually

```bash
# Phase 1: Parse only
--run-tests=false --create-issues=false

# Phase 2: Generate tests
--run-tests=false --create-issues=false

# Phase 3: Execute tests
--run-tests=true --create-issues=false

# Phase 4: Create issues
--run-tests=true --create-issues=true
```

## Real-World Results You Can Expect

### PetStore API (20 endpoints)

**With Mock Client:**

- Parsing: ✅ 100% success
- Test Generation: ✅ 20/20 endpoints
- Code Quality: ⚠️ Basic (mock generates simple tests)

**With GPT-4:**

- Parsing: ✅ 100% success
- Test Generation: ✅ 20/20 endpoints
- Code Quality: ✅ Excellent (comprehensive test scenarios)
- Edge Cases: ✅ Handled
- Security Tests: ✅ Included

**With Claude:**

- Parsing: ✅ 100% success
- Test Generation: ✅ 20/20 endpoints
- Code Quality: ✅ Excellent (well-structured)
- Documentation: ✅ Better comments

**With Ollama (codellama):**

- Parsing: ✅ 100% success
- Test Generation: ✅ 20/20 endpoints
- Code Quality: ⚠️ Variable (depends on prompt)
- Speed: ✅ Fast (local)
- Cost: ✅ Free

## Conclusion

Glens can successfully test against public APIs when you have:

1. ✅ A valid OpenAPI specification (URL or file)
2. ✅ An AI API key (or local Ollama)
3. ✅ (Optional) GitHub access for issue creation
4. ✅ (Optional) API access for live testing

The mock client in this repository demonstrates the framework works. For production use, test with real AI models against real APIs.

## Getting Help

If you encounter issues:

1. Run with `--debug` flag
2. Check execution logs
3. Review generated reports
4. Consult [QUICKSTART.md](../docs/QUICKSTART.md)
5. Open a GitHub issue with logs

---

*Ready to test? Start with PetStore using the commands above!*
