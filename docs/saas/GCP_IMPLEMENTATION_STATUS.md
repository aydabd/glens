# GCP SaaS — Implementation Status

> Tracks what is built, what is partially done, and what remains
> across all 19 SaaS transformation phases.

Last updated: 2026-03-03

## Status Legend

| Symbol | Meaning |
|--------|---------|
| ✅ | Implemented and working |
| 🔶 | Partially implemented |
| ⬜ | Not started |

---

## Phase 1 — Backend REST API + MCP

| Item | Status | Notes |
|------|--------|-------|
| `cmd/api` module + `go.work` entry | ✅ | Module exists, registered in workspace |
| `GET /healthz` | ✅ | Health handler with tests |
| `POST /api/v1/analyze` | ✅ | Analyze handler with tests |
| `POST /api/v1/analyze/preview` | ✅ | Preview handler with tests |
| `GET /api/v1/models` | ✅ | Models handler with tests |
| `POST /api/v1/mcp` | ✅ | MCP JSON-RPC handler with tests |
| CORS + logging + recovery middleware | ✅ | Middleware package with tests |
| `openapi.yaml` | ✅ | Full spec v1.0.0 |
| Dockerfile (distroless) | ✅ | Multi-stage build |
| Makefile + CI workflow | ✅ | `api.yml` workflow |
| Auth-proxy middleware (Secret Manager) | ⬜ | Planned in Phase 7 |
| SSE streaming from analyze | ⬜ | Planned in Phase 19 |

**Summary**: Core API is **complete**. Auth-proxy and SSE streaming are deferred
to their respective phases.

---

## Phase 2 — Frontend (SvelteKit)

| Item | Status | Notes |
|------|--------|-------|
| SvelteKit project scaffold | ⬜ | No `frontend/` directory exists |
| shadcn-svelte UI components | ⬜ | — |
| Spec upload page | ⬜ | — |
| Live progress (SSE) | ⬜ | Depends on Phase 19 |
| Results dashboard + charts | ⬜ | — |
| Auth config page | ⬜ | Depends on Phase 5 + 7 |
| Approval modal | ⬜ | Depends on Phase 8 |
| WCAG 2.2 AA compliance | ⬜ | — |

**Summary**: Frontend is **not started**. The entire SvelteKit application
remains to be built.

---

## Phase 3 — GCP Infrastructure (Terraform)

| Item | Status | Notes |
|------|--------|-------|
| `infra/main.tf` + provider config | ✅ | GCS backend, workspace-driven |
| Cloud Run module (`modules/api/`) | ✅ | Container hosting + LB |
| Firestore module (`modules/firestore/`) | ✅ | Database + indexes |
| BigQuery module (`modules/bigquery/`) | ✅ | Analytics export |
| Cloud Storage module (`modules/storage/`) | ✅ | Static assets + reports |
| Secret Manager module (`modules/secrets/`) | ✅ | Credential storage |
| Observability module (`modules/observability/`) | ✅ | Trace + monitoring |
| Pub/Sub module (`modules/pubsub/`) | ✅ | Event streaming |
| `dev.tfvars` / `prod.tfvars` | ✅ | Workspace-based env config |
| Frontend Cloud Functions module | ⬜ | No frontend module yet |

**Summary**: GCP infrastructure is **nearly complete**. All seven Terraform
modules are built. Only the frontend hosting module is missing (blocked on
Phase 2).

---

## Phase 4 — CI/CD Pipeline

| Item | Status | Notes |
|------|--------|-------|
| Backend workflow (`api.yml`) | ✅ | Test + deploy on push |
| CLI workflow (`glens.yml`) | ✅ | Test + build |
| Infrastructure workflow (`infra.yml`) | ✅ | Plan + apply |
| Workload Identity Federation (OIDC) | ✅ | No SA keys |
| Release-please config | ✅ | Multi-module semver |
| Build + sign + release workflow | ✅ | Multi-platform binaries |
| Docs lint workflow | ✅ | Markdown validation |
| OpenAPI integration test workflow | ✅ | Contract validation |
| Tool-specific workflows (demo, accuracy) | ✅ | Per-module CI |
| Frontend workflow | ⬜ | Blocked on Phase 2 |
| PR preview environments | ⬜ | Not yet implemented |

**Summary**: CI/CD is **largely complete** with 16 workflows. Frontend
pipeline and PR preview environments remain.

---

## Phase 5 — Authentication & Multi-Tenancy

| Item | Status | Notes |
|------|--------|-------|
| Firebase Auth setup | ⬜ | Not configured |
| Auth middleware (JWT verify) | ⬜ | Spec notes "security disabled" |
| Firestore data model (users, workspaces) | ⬜ | Schema designed, not applied |
| Login/signup pages | ⬜ | Blocked on Phase 2 |
| API key management | ⬜ | — |
| Plan limits (free/pro) | ⬜ | — |

**Summary**: Authentication is **not started**. The OpenAPI spec explicitly
notes that global security is disabled until auth middleware is implemented.

---

## Phase 6 — Observability (OpenTelemetry)

| Item | Status | Notes |
|------|--------|-------|
| Observability Terraform module | ✅ | Cloud Trace + Monitoring |
| Structured logging (zerolog) | ✅ | `pkg/logging` module |
| OTel SDK integration in API | ⬜ | Not wired in `cmd/api` |
| Custom spans + metrics | ⬜ | — |
| Trace-log correlation (`trace_id`) | ⬜ | — |

**Summary**: Infrastructure for observability is **ready** (Terraform module +
logging library). OTel SDK integration in the API code is not yet done.

---

## Phase 7 — Target-API Auth Proxy

| Item | Status | Notes |
|------|--------|-------|
| Secret Manager Go SDK | ⬜ | Module exists in Terraform |
| `POST /api/v1/secrets` endpoint | ⬜ | — |
| Auth-proxy middleware | ⬜ | — |
| OAuth2 client credentials handler | ⬜ | — |
| Frontend auth config page | ⬜ | Blocked on Phase 2 |

**Summary**: Auth proxy is **not started**. The Secret Manager Terraform
module is provisioned but the Go integration is not built.

---

## Phase 8 — Destructive Test Safety

| Item | Status | Notes |
|------|--------|-------|
| Endpoint categoriser (`safety/`) | ✅ | `internal/safety/` with tests |
| `POST /api/v1/analyze/preview` | ✅ | Returns risk categories |
| Approval flow (frontend modal) | ⬜ | Blocked on Phase 2 |
| Cleanup hooks | ⬜ | — |

**Summary**: Backend categorisation is **implemented**. The frontend approval
UX and cleanup hooks are not built.

---

## Phase 9 — Epics, Stories & Tasks

| Item | Status | Notes |
|------|--------|-------|
| Epic/story structure documented | ✅ | Phase doc defines breakdown |
| GitHub labels created | ⬜ | — |
| Issues created per task | ⬜ | — |

**Summary**: The breakdown is **documented** but GitHub issues are not yet
created as formal tracking items.

---

## Phase 10 — Testing Strategy

| Item | Status | Notes |
|------|--------|-------|
| Unit tests per package | ✅ | 12+ test files in `cmd/api` |
| Integration tests | 🔶 | Accuracy + demo tools have them |
| E2E scenario tests | 🔶 | Local LLM e2e test exists |
| Architecture acceptance tests | ⬜ | — |
| CI acceptance gate (4-tier) | ⬜ | — |

**Summary**: Unit testing is **solid**. Integration and E2E testing exist
but the full 4-tier CI gate is not configured.

---

## Phase 11 — Event-Driven Architecture

| Item | Status | Notes |
|------|--------|-------|
| Event schema (`events/schema.go`) | ✅ | Domain events defined |
| Event publisher | ✅ | `events/publisher.go` with tests |
| Pub/Sub Terraform module | ✅ | Topics provisioned |
| Cloud Function subscribers | ⬜ | Not built |

**Summary**: Event publishing is **implemented** on the backend. Cloud
Function consumers (report gen, issue creator) are not built.

---

## Phase 12 — API Contracts & Swagger

| Item | Status | Notes |
|------|--------|-------|
| `openapi.yaml` (source of truth) | ✅ | Full spec v1.0.0 |
| OpenAPI integration test workflow | ✅ | CI validates spec |
| `oapi-codegen` server generation | ⬜ | Handlers are hand-written |
| `openapi-typescript` client gen | ⬜ | Blocked on Phase 2 |
| Swagger UI Cloud Function | ⬜ | — |
| API Gateway Terraform | ⬜ | — |

**Summary**: The OpenAPI spec exists and is CI-validated. Code generation
and Swagger UI hosting are not set up.

---

## Phase 13 — Blackbox & E2E Test Examples

| Item | Status | Notes |
|------|--------|-------|
| Test fixture specs | 🔶 | `sample_api.json` exists |
| 20 E2E test scenarios | ⬜ | Only 1 e2e test exists |
| `make test-e2e` target | ⬜ | — |

**Summary**: The test catalogue is **documented** but the 20 scenarios are
not implemented. One local LLM e2e test exists.

---

## Phase 14 — Local & CI Emulators

| Item | Status | Notes |
|------|--------|-------|
| `docker-compose.emulators.yml` | ✅ | Firestore + PubSub + mock secrets |
| Mock secret server (`test/mock-secrets/`) | ✅ | Docker-based mock |
| CI emulator services | ⬜ | Not in CI workflows yet |
| Terraform validation in CI | ⬜ | — |

**Summary**: Local emulator setup is **ready**. CI integration of emulators
is not yet wired.

---

## Phase 15 — Environment Parity

| Item | Status | Notes |
|------|--------|-------|
| Workspace-driven Terraform | ✅ | `dev` / `prod` workspaces |
| Identical modules, variable-only diffs | ✅ | Scaling + log level differ |
| Parity-check CI job | ⬜ | — |
| No env-specific code paths | ✅ | Config-driven runtime |

**Summary**: Infrastructure parity is **achieved**. The automated CI parity
check is not implemented.

---

## Phase 16 — Semantic Versioning & Releases

| Item | Status | Notes |
|------|--------|-------|
| `release-please` configured | ✅ | Multi-module manifest |
| Release workflows per module | ✅ | 6 release workflows |
| RC → dev auto-deploy | ⬜ | — |
| Regression gate before prod promote | ⬜ | — |
| Smoke test scripts | ⬜ | — |
| Deployment e2e tests | ⬜ | — |

**Summary**: Release tooling is **configured**. The promotion pipeline
(RC → dev → regression → prod) is not wired.

---

## Phase 17 — Issue Tracker Provider Abstraction

| Item | Status | Notes |
|------|--------|-------|
| `IssueProvider` interface | ✅ | `internal/issues/provider.go` |
| GitHub provider | ✅ | `github.go` with tests |
| GitLab provider (stub) | ✅ | `gitlab.go` stub |
| Jira provider (stub) | ✅ | `jira.go` stub |
| Provider registry | ✅ | Factory pattern |
| Per-workspace config | ⬜ | Needs auth + Firestore |

**Summary**: The provider abstraction is **fully implemented** with GitHub
as the working provider and GitLab/Jira stubs ready for future work.

---

## Phase 18 — RFC 9457 & Enterprise Auth

| Item | Status | Notes |
|------|--------|-------|
| RFC 9457 `ProblemDetail` struct | ✅ | `internal/handler/problem.go` |
| `application/problem+json` responses | ✅ | All error handlers |
| Error type tests | ✅ | Comprehensive test coverage |
| RBAC types + middleware | ⬜ | — |
| Organisation entity | ⬜ | — |
| OIDC / SAML SSO | ⬜ | — |
| API key scoping per org | ⬜ | — |
| Audit logging | ⬜ | — |
| Rate limiting per org | ⬜ | — |

**Summary**: RFC 9457 error responses are **complete**. Enterprise auth
(RBAC, SSO, orgs) is not started.

---

## Phase 19 — Real-Time Analysis & Live Monitoring

| Item | Status | Notes |
|------|--------|-------|
| `responseWriter` Flusher support | ⬜ | — |
| SSE endpoint (`/runs/{id}/events`) | ⬜ | — |
| Run model in Firestore | ⬜ | — |
| Analysis worker event emission | ⬜ | — |
| Run cancel endpoint | ⬜ | — |
| Frontend live dashboard | ⬜ | Blocked on Phase 2 |
| Pub/Sub test event integration | ⬜ | — |

**Summary**: Real-time streaming is **not started**.

---

## Overall Progress

| Phase | Focus | Status |
|-------|-------|--------|
| 1 | Backend API | ✅ Core complete |
| 2 | Frontend | ⬜ Not started |
| 3 | GCP Infrastructure | ✅ 7/7 modules |
| 4 | CI/CD | ✅ 16 workflows |
| 5 | Auth & Multi-tenancy | ⬜ Not started |
| 6 | Observability | 🔶 Infra ready, code not wired |
| 7 | Auth Proxy | ⬜ Not started |
| 8 | Test Safety | 🔶 Backend done, frontend not |
| 9 | Epics & Stories | 🔶 Documented, not tracked |
| 10 | Testing Strategy | 🔶 Unit tests good, gates missing |
| 11 | Event Architecture | 🔶 Publisher done, consumers not |
| 12 | API Contracts | 🔶 Spec exists, codegen not set up |
| 13 | E2E Tests | ⬜ 1/20 tests exist |
| 14 | Emulators | 🔶 Local ready, CI not wired |
| 15 | Env Parity | 🔶 Infra done, CI check missing |
| 16 | Releases | 🔶 Tooling done, pipeline not |
| 17 | Issue Providers | ✅ Interface + GitHub provider |
| 18 | RFC 9457 + Enterprise | 🔶 Errors done, enterprise not |
| 19 | Real-time Streaming | ⬜ Not started |

### What is fully implemented (GCP)

- REST API with 5 endpoints, MCP support, and OpenAPI spec
- Full Terraform IaC (7 modules) for Cloud Run, Firestore, BigQuery,
  Pub/Sub, Secret Manager, Storage, and Observability
- 16 CI/CD workflows with release-please
- Issue provider abstraction (GitHub, GitLab stub, Jira stub)
- RFC 9457 error responses
- Event publishing + schema
- Endpoint safety categoriser
- Docker + local emulator setup

### What remains for GCP SaaS

- **Frontend** — entire SvelteKit application (Phase 2)
- **Authentication** — Firebase Auth, API keys, plan limits (Phase 5)
- **Auth Proxy** — Secret Manager integration in Go (Phase 7)
- **Observability** — OTel SDK wiring in API code (Phase 6)
- **Event Consumers** — Cloud Functions for reports, issues, notifications (Phase 11)
- **Code Generation** — oapi-codegen + openapi-typescript (Phase 12)
- **E2E Tests** — 19 of 20 blackbox scenarios (Phase 13)
- **CI Emulators** — Wire emulators into CI workflows (Phase 14)
- **Release Pipeline** — RC → dev → regression → prod promotion (Phase 16)
- **Enterprise Auth** — RBAC, SSO, orgs, audit (Phase 18)
- **Real-time Streaming** — SSE, run model, cancel (Phase 19)
