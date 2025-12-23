# Test OpenAPI Specifications

This directory contains OpenAPI specification files used for testing Glens' accuracy in generating integration tests.

## Available Specifications

### `sample_api.json`

A simple OpenAPI 3.0.3 specification demonstrating basic CRUD operations.

**Endpoints:**
- `GET /users` - List all users
- `GET /users/{id}` - Get a specific user by ID
- `POST /posts` - Create a new post

**Features:**
- Path parameters
- Query parameters
- Multiple response codes
- RESTful design patterns

## Using These Specifications

### Test with Glens

```bash
# Analyze the sample API
./build/glens analyze test_specs/sample_api.json \
  --ai-models=mock \
  --create-issues=false \
  --run-tests=false \
  --output=reports/sample_api_report.md
```

### Run Accuracy Tests

```bash
# Run the full accuracy test suite
./scripts/test_accuracy.sh
```

## Adding New Test Specifications

To add a new OpenAPI spec for testing:

1. Create a JSON or YAML file in this directory
2. Ensure it's a valid OpenAPI 2.0, 3.0, or 3.1 specification
3. Add it to the test suite in `scripts/test_accuracy.sh`

### Example Template

```json
{
  "openapi": "3.0.3",
  "info": {
    "title": "Your API Name",
    "version": "1.0.0",
    "description": "API description"
  },
  "servers": [
    {
      "url": "https://api.example.com/v1"
    }
  ],
  "paths": {
    "/your-endpoint": {
      "get": {
        "summary": "Endpoint description",
        "operationId": "operationName",
        "responses": {
          "200": {
            "description": "Success response"
          }
        }
      }
    }
  }
}
```

## Validation

Validate your OpenAPI specs using:

```bash
# Using swagger-cli (if installed)
swagger-cli validate test_specs/your_spec.json

# Using openapi-generator (if installed)
openapi-generator validate -i test_specs/your_spec.json
```

## Test Complexity Levels

### Level 1: Basic (Current)
- Simple GET/POST endpoints
- Basic parameters
- Standard response codes

### Level 2: Intermediate (Future)
- Authentication (OAuth, API keys)
- Complex request bodies
- Multiple content types
- Pagination

### Level 3: Advanced (Future)
- Webhooks
- Callbacks
- Complex schemas with $ref
- Multiple servers
- Security schemes

## Related Files

- `accuracy_tests/` - Test results directory
- `scripts/test_accuracy.sh` - Testing script
- `ACCURACY_REPORT.md` - Comprehensive accuracy report

## Contributing

When adding test specs:
1. Use realistic API designs
2. Include diverse endpoint patterns
3. Add complex scenarios gradually
4. Document the purpose and features
5. Ensure specifications are valid

---

*Test specifications are used for validating Glens' parsing and test generation capabilities*
