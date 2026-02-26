# Glens SaaS Transformation Plan

> Master plan for converting the Glens CLI into a cloud-hosted SaaS on GCP.

## Goal

Expose Glens capabilities (OpenAPI analysis, AI test generation, GitHub issue
creation) as a multi-tenant web service with a clean, accessible frontend.

## Guiding Principles

1. **Reuse** — the existing Go backend is wrapped, not rewritten.
2. **Serverless-first** — Cloud Run + Cloud Functions on GCP free tier.
3. **Accessibility** — WCAG 2.2 AA compliance from day one.
4. **Independent phases** — every phase can be built in parallel.

## Requirement Map

Each requirement has a unique ID. Frontend requirements (`FE-*`) are linked to
the backend requirements (`BE-*`) they depend on.

### Backend Requirements (BE)

| ID | Requirement | Current CLI feature |
|----|-------------|---------------------|
| BE-01 | Parse OpenAPI spec from URL or upload | `parser.ParseOpenAPISpec` |
| BE-02 | Generate tests via AI models | `ai.NewManager` / `GenerateTest` |
| BE-03 | Execute generated tests | `generator.ExecuteTest` |
| BE-04 | Create GitHub issues on failure | `github.CreateEndpointIssue` |
| BE-05 | Generate markdown/HTML reports | `reporter.GenerateReport` |
| BE-06 | List/select AI models | `models` command |
| BE-07 | Filter endpoints by operation ID | `--op-id` flag |
| BE-08 | REST API exposing BE-01 – BE-07 | **New — Phase 1** |
| BE-09 | Authentication & API key management | **New — Phase 5** |
| BE-10 | Multi-tenant workspace isolation | **New — Phase 5** |

### Frontend Requirements (FE)

| ID | Requirement | Depends on |
|----|-------------|------------|
| FE-01 | Upload / paste OpenAPI spec URL | BE-01, BE-08 |
| FE-02 | Select AI models for generation | BE-06, BE-08 |
| FE-03 | View real-time test generation progress | BE-02, BE-08 |
| FE-04 | Display test results dashboard | BE-03, BE-05, BE-08 |
| FE-05 | Link to / create GitHub issues | BE-04, BE-08 |
| FE-06 | Filter endpoints by operation ID | BE-07, BE-08 |
| FE-07 | Download reports (Markdown / HTML) | BE-05, BE-08 |
| FE-08 | WCAG 2.2 AA accessible UI | — |
| FE-09 | Responsive design (mobile-first) | — |
| FE-10 | Login / workspace dashboard | BE-09, BE-10 |

### Infrastructure Requirements (IN)

| ID | Requirement | Phase |
|----|-------------|-------|
| IN-01 | GCP project + free-tier resource plan | Phase 3 |
| IN-02 | Terraform IaC for all GCP resources | Phase 3 |
| IN-03 | Cloud Run service for backend API | Phase 3 |
| IN-04 | Cloud Functions for frontend SSR/BFF | Phase 3 |
| IN-05 | Cloud Storage for static assets | Phase 3 |
| IN-06 | GitHub Actions deploy pipeline | Phase 4 |
| IN-07 | Preview environments per PR | Phase 4 |
| IN-08 | Secret management (Secret Manager) | Phase 4 |

## Phase Overview

| Phase | Document | Can start | Depends on |
|-------|----------|-----------|------------|
| 1 | [PHASE1_BACKEND_API.md](PHASE1_BACKEND_API.md) | Immediately | — |
| 2 | [PHASE2_FRONTEND.md](PHASE2_FRONTEND.md) | Immediately | — |
| 3 | [PHASE3_GCP_INFRA.md](PHASE3_GCP_INFRA.md) | Immediately | — |
| 4 | [PHASE4_CICD.md](PHASE4_CICD.md) | Immediately | — |
| 5 | [PHASE5_AUTH_MULTITENANCY.md](PHASE5_AUTH_MULTITENANCY.md) | Immediately | — |

All phases are designed to be **completely independent**. Integration happens
at well-defined API contracts (OpenAPI spec for the REST API, Terraform
outputs for infrastructure).

## Tech Stack Summary

| Layer | Technology | Why |
|-------|-----------|-----|
| Backend API | Go + net/http (stdlib) | Reuses existing code |
| Frontend | SvelteKit (SSR on Cloud Functions) | Minimal JS, built-in a11y, fast |
| Hosting | GCP Cloud Run + Cloud Functions | Free tier generous |
| IaC | Terraform | Industry standard, GCP provider |
| CI/CD | GitHub Actions | Already used in repo |
| Auth | Firebase Authentication | Free tier, GCP-native |
| Database | Firestore | Serverless, free tier |
| Storage | Cloud Storage | Static assets, reports |

## How to Read This Plan

1. Start with this file for the full picture.
2. Open the phase you want to implement.
3. Each phase lists its own deliverables, requirements covered, and steps.
4. Requirement IDs link back to this master table.
