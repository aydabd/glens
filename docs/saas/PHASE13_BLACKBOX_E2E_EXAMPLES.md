# Phase 13 — 20 Blackbox & E2E Test Examples

> TDD-first test catalogue: 20 scenarios driving spec, design, and code.

## Requirements Covered

BB-01 – BB-20. Each test drives doc, design, implementation, and QA.

## How to Use

1. Pick a test → write the spec fixture → write the test → implement
2. Each test is a GitHub issue (story) assignable to one agent
3. Tests use mock servers — no real AI/GitHub API needed
4. All tests run in CI via `make test-e2e` in `cmd/glens/e2e/`

## Test Catalogue

### Core Analysis Pipeline (BB-01 – BB-05)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-01 | Analyze read-only spec | `petstore_readonly.json` (GET only) | Exit 0, report contains all endpoints, no issues |
| BB-02 | Analyze CRUD spec | `petstore_crud.json` (GET/POST/PUT/DELETE) | Exit 0, report covers all 4 HTTP methods |
| BB-03 | Analyze with `--op-id` filter | `petstore_crud.json --op-id getPetById` | Only 1 endpoint in report |
| BB-04 | Analyze invalid spec URL | `http://invalid.test/nope.json` | Exit non-zero, error message, no report |
| BB-05 | Analyze empty spec (no paths) | `empty_spec.json` (valid OpenAPI, 0 paths) | Exit 0, report says "0 endpoints" |

### AI Model Selection (BB-06 – BB-08)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-06 | Single model (ollama) | `--ai-models ollama` | Mock hit once per endpoint |
| BB-07 | Multiple models | `--ai-models ollama,gpt4` | Mock hit twice per endpoint |
| BB-08 | Invalid model name | `--ai-models nonexistent` | Exit non-zero, helpful error |

### Report Generation (BB-09 – BB-11)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-09 | Markdown output (default) | `--output report.md` | Valid markdown, contains spec title |
| BB-10 | Custom output path | `--output /tmp/custom/report.md` | File created at specified path |
| BB-11 | Report with multi-model | Two models on CRUD spec | Report has model comparison section |

### GitHub Issue Management (BB-12 – BB-14)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-12 | No issues on all-pass | `--create-issues=true` (tests pass) | Zero POST calls to mock GitHub |
| BB-13 | Issue on test failure | `--create-issues=true` (test fails) | Exactly 1 issue created via mock |
| BB-14 | `--create-issues=false` | Explicit opt-out | Zero GitHub API calls regardless |

### Error Resilience (BB-15 – BB-17)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-15 | AI server unreachable | Mock on closed port | Exit 0, no panic, error logged |
| BB-16 | AI returns garbage | Mock returns `{invalid}` | Exit 0, skip endpoint, continue |
| BB-17 | Spec with auth schemes | `petstore_auth.json` (OAuth2 + API key) | Endpoints parsed with security info |

### CLI UX & Config (BB-18 – BB-20)

| ID | Test | Input | Expected |
|----|------|-------|----------|
| BB-18 | `models ollama list` | Mock Ollama server | Lists model names from mock |
| BB-19 | `cleanup --dry-run` | Mock GitHub with 3 open issues | Shows 3 issues, closes 0 |
| BB-20 | Config file override | Config YAML with custom URL | Mock at custom URL is called |

## Required Test Fixtures

| Fixture | Endpoints | Purpose |
|---------|-----------|---------|
| `petstore_readonly.json` | 2× GET | BB-01: safe read-only |
| `petstore_crud.json` | GET, POST, PUT, DELETE | BB-02/03/11: full CRUD |
| `empty_spec.json` | 0 paths | BB-05: edge case |
| `petstore_auth.json` | 2× GET with security | BB-17: auth schemes |

Existing `sample_api.json` used for BB-06–10, BB-12–16, BB-18–20.

## TDD Workflow Per Test

1. Create fixture → 2. Write failing test → 3. Implement code →
4. Test passes → 5. Update docs → 6. CI gate green

## Definition of Done (per BB-XX)

- [ ] Fixture spec in `test_specs/`; E2E test passes in CI
- [ ] Feature code implemented with unit tests
- [ ] Docs updated if user-facing behaviour changed

Each BB-XX → GitHub issue (`epic:13-e2e` label), independently
runnable via `go test -run BB_XX -v` in `cmd/glens/e2e/`.

## Success Criteria

- [ ] All 20 E2E tests pass in CI
- [ ] 4 new fixture specs in `test_specs/`
- [ ] Zero flaky tests (mock servers, no network)
