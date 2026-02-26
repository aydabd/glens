# Glens SaaS Transformation Plan

> Master plan for converting the Glens CLI into a cloud-hosted SaaS on GCP.

## Goal

Expose Glens as a multi-tenant web service with accessible frontend,
observability, secure API-auth proxying, destructive-test safety,
contract-first APIs, event-driven architecture, and comprehensive
automated testing at every level (unit, integration, e2e).

## Guiding Principles

1. **Reuse** — existing Go backend wrapped, not rewritten.
2. **Serverless-first** — Cloud Run + Cloud Functions, GCP free tier.
3. **Contract-first** — OpenAPI spec is the single source of truth.
4. **Test-driven** — acceptance tests define done; CI enforces them.
5. **Event-driven** — async side effects via Pub/Sub + Cloud Functions.
6. **Domain isolation** — each domain in its own package; no conflicts.
7. **Security-first** — secrets never reach the frontend.
8. **Independent phases** — all 13 phases built in parallel.

## Requirement Map

### Backend (BE), Frontend (FE)

| ID | Requirement | Phase |
|----|-------------|-------|
| BE-01 – BE-07 | Parse, generate, execute, issues, reports, models, filter | existing |
| BE-08 | REST API / MCP | 1 |
| BE-09 – BE-10 | Auth, multi-tenancy | 5 |
| FE-01 – FE-09 | Upload, progress, dashboard, a11y, responsive | 2 |
| FE-10 – FE-12 | Login, approval dialog, auth config | 2, 5 |

### Security (SE), Observability (OB), Database (DB), Test Safety (TS)

| ID | Requirement | Phase |
|----|-------------|-------|
| SE-01 – SE-04 | Auth types, headers, server-side secrets | 7 |
| OB-01 – OB-03 | OTel SDK, Cloud Trace, structured logs | 6 |
| DB-01 – DB-03 | Firestore results, charts, BigQuery | 3 |
| TS-01 – TS-04 | Categorise, warn, approve, cleanup | 8 |

### Epics & Stories (EP), Quality (QA), Events (EV), Contracts (AC)

| ID | Requirement | Phase |
|----|-------------|-------|
| EP-01 | Epics → stories → tasks (GitHub Issues) | 9 |
| EP-02 | Acceptance criteria per task | 9 |
| EP-03 | Parallel work without file conflicts | 9 |
| EP-04 | Definition of done = tests pass | 9 |
| QA-01 | Unit tests per domain package | 10 |
| QA-02 | Integration tests per boundary | 10 |
| QA-03 | E2E scenario-based tests | 10 |
| QA-04 | Architecture acceptance checks in CI | 10 |
| QA-05 | Definition-of-done gates in CI | 10 |
| EV-01 | Domain events (analyze, test, report) | 11 |
| EV-02 | Pub/Sub topics per event type | 11 |
| EV-03 | Event-triggered Cloud Functions | 11 |
| EV-04 | Event schema contracts | 11 |
| AC-01 | OpenAPI contract-first development | 12 |
| AC-02 | Swagger UI published via Cloud Function | 12 |
| AC-03 | Auto-generated client SDKs | 12 |
| AC-04 | Contract-drift tests in CI | 12 |
| BB-01 – BB-20 | 20 blackbox/E2E test examples for TDD | 13 |

## Phases (all independent, parallel)

| # | Document | Focus |
|---|----------|-------|
| 1 | [PHASE1_BACKEND_API.md](PHASE1_BACKEND_API.md) | REST / MCP API |
| 2 | [PHASE2_FRONTEND.md](PHASE2_FRONTEND.md) | Frontend |
| 3 | [PHASE3_GCP_INFRA.md](PHASE3_GCP_INFRA.md) | GCP IaC + DB |
| 4 | [PHASE4_CICD.md](PHASE4_CICD.md) | CI/CD deploy |
| 5 | [PHASE5_AUTH_MULTITENANCY.md](PHASE5_AUTH_MULTITENANCY.md) | User auth |
| 6 | [PHASE6_OBSERVABILITY.md](PHASE6_OBSERVABILITY.md) | OpenTelemetry |
| 7 | [PHASE7_SECURITY_AUTH_PROXY.md](PHASE7_SECURITY_AUTH_PROXY.md) | Target-API auth |
| 8 | [PHASE8_TEST_SAFETY.md](PHASE8_TEST_SAFETY.md) | Test safety |
| 9 | [PHASE9_EPICS_STORIES.md](PHASE9_EPICS_STORIES.md) | Epics → tasks |
| 10 | [PHASE10_TESTING_STRATEGY.md](PHASE10_TESTING_STRATEGY.md) | Testing strategy |
| 11 | [PHASE11_EVENT_ARCHITECTURE.md](PHASE11_EVENT_ARCHITECTURE.md) | Events / Pub/Sub |
| 12 | [PHASE12_API_CONTRACTS.md](PHASE12_API_CONTRACTS.md) | API contracts |
| 13 | [PHASE13_BLACKBOX_E2E_EXAMPLES.md](PHASE13_BLACKBOX_E2E_EXAMPLES.md) | 20 E2E tests |

## Tech Stack (latest stable, Feb 2026)

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend | Go + `net/http` | 1.24 |
| Frontend | SvelteKit + shadcn-svelte | 2.x |
| Hosting | Cloud Run / Cloud Functions Gen2 | — |
| IaC | Terraform + google provider | 1.10 / 6.x |
| CI/CD | GitHub Actions | v4 |
| Auth | Firebase Auth | 5.x |
| DB | Firestore + BigQuery (export) | — |
| Observability | OpenTelemetry → Cloud Trace | 1.34 |
| Events | Cloud Pub/Sub | — |
| API Gateway | GCP API Gateway | — |
| Contracts | oapi-codegen + openapi-typescript | — |
| Secrets | GCP Secret Manager | — |
