# Phase 10 — Testing Strategy

> Unit, integration, e2e, and architecture acceptance tests.

## Requirements Covered

QA-01 (unit), QA-02 (integration), QA-03 (e2e scenarios),
QA-04 (architecture acceptance in CI), QA-05 (definition-of-done).

## Test Pyramid

```
       ┌─────────┐  E2E: 5-10 scenario tests (full flow)
      ┌┴─────────┴┐ Integration: per-boundary (API→DB, API→SecretMgr)
     ┌┴────────────┴┐ Unit: per-package (handler, safety, authproxy)
```

## Unit Tests (per Go package)

| Package | What | Example |
|---------|------|---------|
| `handler/` | HTTP status, JSON, SSE | `TestAnalyze_Returns200` |
| `safety/` | Categorisation logic | `TestCategorise_DELETE` |
| `authproxy/` | Ref resolution, headers | `TestResolveSecret` |
| `auth/` | JWT verify, middleware | `TestAuth_InvalidToken` |
| `events/` | Publish, payload shape | `TestPublish_Event` |

Table-driven with `testify`. Mock external deps via interfaces.

## Integration Tests (per boundary)

| Boundary | Fixture |
|----------|---------|
| API → Firestore | Emulator |
| API → Secret Mgr | Emulator/mock |
| API → AI provider | Mock server |
| Frontend → API | Test server |

`_integration_test.go` suffix. Run: `-tags=integration`. CI uses
Docker emulators.

## E2E Scenario Tests (QA-03)

| Scenario | Given | When/Then |
|----------|-------|-----------|
| Read-only analysis | User + workspace | Upload spec → 🟢 safe, SSE ok, report stored |
| Destructive approval | Spec with DELETE | Preview → 🔴 warn → approve → only approved run |
| Auth-proxy flow | Bearer stored | Analyze with ref → token resolved, no leak |
| Multi-tenant | Two workspaces | User A run → user B blocked by rules |

## Architecture Acceptance Tests (QA-04)

| Check | Fails if |
|-------|----------|
| No cross-domain imports | `handler/` imports `safety/` internals |
| OpenAPI ↔ handlers match | Endpoint added without spec update |
| No secrets in logs | Raw secret value found in output |
| Dependency graph valid | Circular or forbidden import |

## CI Acceptance Gate

```yaml
acceptance:
  steps:
    - run: cd cmd/api && make test              # unit
    - run: cd cmd/api && make test-integration  # boundary
    - run: cd cmd/api && make test-arch         # architecture
    - run: cd e2e && make test-e2e              # scenarios
```

Every PR must pass all four gates before merge.

## File Conventions

Unit: `analyze_test.go` next to `analyze.go`. Integration:
`analyze_integration_test.go`. E2E: `e2e/scenarios/*_test.go`.

## Steps

1. Define test interfaces for all external dependencies
2. Write unit tests per `internal/` package
3. Set up emulators in CI; write integration tests
4. Write e2e scenarios; add architecture checks to CI

## Success Criteria

- [ ] ≥ 80% unit test coverage per Go package
- [ ] Integration tests pass with emulators in CI
- [ ] 4+ e2e scenarios pass against deployed stack
- [ ] Architecture checks block PRs with violations
