# Glens SaaS Transformation Plan

> Master plan for converting the Glens CLI into a cloud-hosted SaaS on GCP.

## Goal

Expose Glens as a multi-tenant web service with accessible frontend,
observability, secure API-auth proxying, and destructive-test safety.

## Guiding Principles

1. **Reuse** — existing Go backend wrapped, not rewritten.
2. **Serverless-first** — Cloud Run + Cloud Functions, GCP free tier.
3. **Accessibility** — WCAG 2.2 AA from day one.
4. **Security-first** — secrets never reach the frontend.
5. **Observability** — OpenTelemetry end-to-end.
6. **Safety** — destructive tests categorised and user-approved.
7. **Independent phases** — all 8 phases built in parallel.

## Requirement Map

### Backend (BE)

| ID | Requirement | Source |
|----|-------------|--------|
| BE-01 – BE-07 | Parse spec, AI generate, execute, issues, reports, models, filter | existing |
| BE-08 | REST API / MCP exposing BE-01 – BE-07 | Phase 1 |
| BE-09 | Authentication & API key management | Phase 5 |
| BE-10 | Multi-tenant workspace isolation | Phase 5 |

### Frontend (FE)

| ID | Requirement | Depends on |
|----|-------------|------------|
| FE-01 – FE-03 | Spec upload, model select, live progress | BE-01/06/02, BE-08 |
| FE-04 | Test results dashboard + charts | BE-03, DB-01 |
| FE-05 – FE-07 | Issue links, filter, download reports | BE-04/07/05, BE-08 |
| FE-08 – FE-09 | WCAG 2.2 AA, responsive | — |
| FE-10 | Login / workspace dashboard | BE-09 |
| FE-11 | Destructive-test approval dialog | TS-02 |
| FE-12 | Target-API auth config (no secret leak) | SE-03 |

### Security (SE), Observability (OB), Database (DB), Test Safety (TS)

| ID | Requirement | Phase |
|----|-------------|-------|
| SE-01 | Target-API auth (Bearer, API-key, OAuth2, mTLS) | 7 |
| SE-02 | Custom header injection (Kong, gateway) | 7 |
| SE-03 | Secrets server-side only (Secret Manager) | 7 |
| SE-04 | Frontend never sees raw credentials | 7 |
| OB-01 | OpenTelemetry SDK (traces + metrics) | 6 |
| OB-02 | Cloud Trace + Cloud Monitoring export | 6 |
| OB-03 | Structured JSON logging | 6 |
| DB-01 | Persist test results (Firestore) | 3 |
| DB-02 | Historical result charts | 2 |
| DB-03 | BigQuery analytics export | 3 |
| TS-01 | Categorise endpoints: read vs write/delete | 8 |
| TS-02 | Warn before destructive tests | 8 |
| TS-03 | Batch-approve / reject | 8 |
| TS-04 | Post-test cleanup hooks | 8 |

## Phases (all independent, parallel)

| # | Document | Focus |
|---|----------|-------|
| 1 | [PHASE1_BACKEND_API.md](PHASE1_BACKEND_API.md) | REST / MCP API |
| 2 | [PHASE2_FRONTEND.md](PHASE2_FRONTEND.md) | Frontend (Flutter eval) |
| 3 | [PHASE3_GCP_INFRA.md](PHASE3_GCP_INFRA.md) | GCP IaC + DB |
| 4 | [PHASE4_CICD.md](PHASE4_CICD.md) | CI/CD deploy |
| 5 | [PHASE5_AUTH_MULTITENANCY.md](PHASE5_AUTH_MULTITENANCY.md) | User auth |
| 6 | [PHASE6_OBSERVABILITY.md](PHASE6_OBSERVABILITY.md) | OpenTelemetry |
| 7 | [PHASE7_SECURITY_AUTH_PROXY.md](PHASE7_SECURITY_AUTH_PROXY.md) | Target-API auth |
| 8 | [PHASE8_TEST_SAFETY.md](PHASE8_TEST_SAFETY.md) | Test safety |

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
| Secrets | GCP Secret Manager | — |
