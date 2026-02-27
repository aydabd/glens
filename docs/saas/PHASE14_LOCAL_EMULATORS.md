# Phase 14 — Local & CI Testing with GCP Emulators

> Test Cloud Functions, Firestore, and Pub/Sub locally and in CI.

## Requirements Covered

EM-01 (local emulators), EM-02 (CI emulators), EM-03 (Terraform
validation), EM-04 (deployment smoke tests).

## Problem

Cloud Functions + Firestore + Pub/Sub only testable in real GCP
without emulators. We need deterministic local + CI testing without
cloud accounts.

## GCP Emulators

| Emulator | Image / Tool | Port |
|----------|-------------|------|
| Firestore | `google/cloud-sdk` (gcloud beta) | 8085 |
| Pub/Sub | `google/cloud-sdk` (gcloud beta) | 8086 |
| Cloud Functions | `functions-framework-go` | 8087 |
| Secret Manager | Mock server (Go) | 8088 |

## Docker Compose (`docker-compose.emulators.yml`)

```yaml
services:
  firestore:
    image: google/cloud-sdk:latest
    command: gcloud beta emulators firestore start --host-port=0.0.0.0:8085
    ports: ["8085:8085"]  # 0.0.0.0 required inside container
  pubsub:
    image: google/cloud-sdk:latest
    command: gcloud beta emulators pubsub start --host-port=0.0.0.0:8086
    ports: ["8086:8086"]
  secret-mock:
    build: ./test/mock-secrets
    ports: ["8088:8088"]
```

## Integration Test Setup (Go)

```go
func TestMain(m *testing.M) {
    os.Setenv("FIRESTORE_EMULATOR_HOST", "localhost:8085")
    os.Setenv("PUBSUB_EMULATOR_HOST", "localhost:8086")
    os.Exit(m.Run())
}
```

GCP client libraries auto-detect `*_EMULATOR_HOST` env vars.

## CI Workflow (`ci-integration.yml`)

```yaml
jobs:
  integration:
    services:
      firestore:
        image: google/cloud-sdk:latest
        ports: ["8085:8085"]
      pubsub:
        image: google/cloud-sdk:latest
        ports: ["8086:8086"]
    env:
      FIRESTORE_EMULATOR_HOST: localhost:8085
      PUBSUB_EMULATOR_HOST: localhost:8086
    steps:
      - run: cd cmd/api && make test-integration
```

## Terraform Validation (EM-03)

- `terraform validate` — syntax check (no GCP needed)
- `terraform plan` with mock provider — structure check
- Terratest in CI — create + verify + destroy in test project

## Deployment Smoke Tests (EM-04)

After deploy to dev: health check, Firestore r/w, Pub/Sub round-trip,
Secret Manager resolution, Cloud Function trigger — all automated.

## Steps

1. Create `docker-compose.emulators.yml` at repo root
2. Add mock secret server in `test/mock-secrets/`
3. Update integration tests to use emulator env vars
4. Add CI workflow with emulator services + Terratest

## Success Criteria

- [ ] `docker compose up` starts all emulators locally
- [ ] Integration tests pass against emulators in CI
- [ ] `terraform validate` passes without GCP credentials
- [ ] Smoke tests verify deployed dev environment
