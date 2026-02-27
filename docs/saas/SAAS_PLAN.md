# Glens SaaS Transformation Plan

> Master plan: Glens CLI → multi-tenant SaaS on GCP (17 phases).

## Guiding Principles

1. **Reuse** — wrap existing Go backend, don't rewrite.
2. **Serverless-first** — Cloud Run + Cloud Functions, GCP free tier.
3. **Contract-first** — OpenAPI spec is the single source of truth.
4. **Test-driven** — acceptance tests define done; CI enforces them.
5. **Event-driven** — Pub/Sub + Cloud Functions for async side effects.
6. **Domain isolation** — own package per domain; no conflicts.
7. **Security-first** — secrets never reach the frontend.
8. **Env parity** — dev/prod differ only in log level + scaling.

## Requirement Map

| ID | Requirement | Phase |
|----|-------------|-------|
| BE-01 – BE-08 | Parse, generate, execute, issues, reports, API/MCP | 1 |
| BE-09 – BE-10 | Auth, multi-tenancy | 5 |
| FE-01 – FE-12 | Upload, dashboard, a11y, login, approval, auth UI | 2, 5 |
| SE-01 – SE-04 | Auth types, headers, server-side secrets | 7 |
| OB-01 – OB-03 | OTel SDK, Cloud Trace, structured logs | 6 |
| DB-01 – DB-03 | Firestore results, charts, BigQuery | 3 |
| TS-01 – TS-04 | Categorise, warn, approve, cleanup | 8 |
| EP-01 – EP-04 | Epics → stories → tasks, acceptance criteria | 9 |
| QA-01 – QA-05 | Unit, integration, e2e, arch checks, DoD gates | 10 |
| EV-01 – EV-04 | Domain events, Pub/Sub, Cloud Function triggers | 11 |
| AC-01 – AC-04 | OpenAPI contract-first, Swagger UI, SDK gen, drift | 12 |
| BB-01 – BB-20 | 20 blackbox/E2E test examples for TDD | 13 |
| EM-01 – EM-04 | Local emulators, CI emulators, Terraform validate | 14 |
| EP-01 – EP-05 | Dev/prod parity, no env-specific code paths | 15 |
| RV-01 – RV-06 | Semver, pre-release, promote, regression gates | 16 |
| IP-01 – IP-05 | Issue provider interface, GitHub/GitLab/Jira | 17 |

## Phases (all independent, parallel)

| # | Document | Focus |
|---|----------|-------|
| 1 | [PHASE1_BACKEND_API.md](PHASE1_BACKEND_API.md) | REST / MCP API |
| 2 | [PHASE2_FRONTEND.md](PHASE2_FRONTEND.md) | Frontend |
| 3 | [PHASE3_GCP_INFRA.md](PHASE3_GCP_INFRA.md) | GCP IaC + DB |
| 4 | [PHASE4_CICD.md](PHASE4_CICD.md) | CI/CD deploy |
| 5 | [PHASE5_AUTH_MULTITENANCY.md](PHASE5_AUTH_MULTITENANCY.md) | User auth |
| 6 | [PHASE6_OBSERVABILITY.md](PHASE6_OBSERVABILITY.md) | OpenTelemetry |
| 7 | [PHASE7_SECURITY_AUTH_PROXY.md](PHASE7_SECURITY_AUTH_PROXY.md) | API auth proxy |
| 8 | [PHASE8_TEST_SAFETY.md](PHASE8_TEST_SAFETY.md) | Test safety |
| 9 | [PHASE9_EPICS_STORIES.md](PHASE9_EPICS_STORIES.md) | Epics → tasks |
| 10 | [PHASE10_TESTING_STRATEGY.md](PHASE10_TESTING_STRATEGY.md) | Test strategy |
| 11 | [PHASE11_EVENT_ARCHITECTURE.md](PHASE11_EVENT_ARCHITECTURE.md) | Events / Pub/Sub |
| 12 | [PHASE12_API_CONTRACTS.md](PHASE12_API_CONTRACTS.md) | API contracts |
| 13 | [PHASE13_BLACKBOX_E2E_EXAMPLES.md](PHASE13_BLACKBOX_E2E_EXAMPLES.md) | 20 E2E tests |
| 14 | [PHASE14_LOCAL_EMULATORS.md](PHASE14_LOCAL_EMULATORS.md) | GCP emulators |
| 15 | [PHASE15_ENV_PARITY.md](PHASE15_ENV_PARITY.md) | Dev/prod parity |
| 16 | [PHASE16_RELEASES_VERSIONING.md](PHASE16_RELEASES_VERSIONING.md) | Semver releases |
| 17 | [PHASE17_ISSUE_PROVIDER.md](PHASE17_ISSUE_PROVIDER.md) | Issue providers |

## Tech Stack (latest stable, Feb 2026)

| Layer | Technology | Version |
|-------|-----------|---------|
| Backend | Go + `net/http` | 1.24 |
| Frontend | SvelteKit + shadcn-svelte | 2.x |
| Hosting | Cloud Run / Cloud Functions Gen2 | — |
| IaC | Terraform + google provider | 1.10 / 6.x |
| CI/CD | GitHub Actions + release-please | v4 |
| Auth | Firebase Auth | 5.x |
| DB | Firestore + BigQuery (export) | — |
| Observability | OpenTelemetry → Cloud Trace | 1.34 |
| Events | Cloud Pub/Sub | — |
| Contracts | oapi-codegen + openapi-typescript | — |
| Secrets | GCP Secret Manager | — |
| Emulators | gcloud beta emulators + Docker Compose | — |
| Issues | GitHub API (GitLab, Jira planned) | — |
