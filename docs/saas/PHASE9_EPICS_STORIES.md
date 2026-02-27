# Phase 9 — Epics, Stories & Tasks (GitHub Issues)

> Product-owner breakdown of all work into parallel, non-conflicting units.

## Requirements Covered

EP-01 (epics → stories → tasks), EP-02 (acceptance criteria per task),
EP-03 (parallel without conflicts), EP-04 (definition of done).

## Epic Structure

Each epic maps 1:1 to a phase. Stories are vertical slices within a
domain. Tasks are atomic units assignable to one developer/agent.

## Epic 1 — Backend API (Phase 1)

| Story | Tasks | Acceptance |
|-------|-------|------------|
| S1.1 Health endpoint | T: handler, test, OpenAPI entry | `GET /healthz` → 200 |
| S1.2 Analyze endpoint | T: handler, SSE stream, parser integration | SSE events received |
| S1.3 Models endpoint | T: handler, test | JSON model list returned |
| S1.4 MCP JSON-RPC | T: handler, tool registration | `tools/call` executes |
| S1.5 Analyze preview | T: categoriser, handler, test | Risk categories returned |

## Epic 2 — Frontend (Phase 2)

| Story | Tasks | Acceptance |
|-------|-------|------------|
| S2.1 Spec upload page | T: route, form, API call | File/URL accepted |
| S2.2 Live progress | T: SSE client, progress UI | Real-time updates |
| S2.3 Results dashboard | T: table, charts, filters | Data renders correctly |
| S2.4 Approval modal | T: modal, risk grouping, batch buttons | Approve/reject works |
| S2.5 Auth config page | T: form, secret store, ref display | Ref stored, no leak |

## Epic 3 — GCP Infra (Phase 3)

| Story | Tasks | Acceptance |
|-------|-------|------------|
| S3.1 Cloud Run module | T: TF module, variables, outputs | `terraform plan` ok |
| S3.2 Firestore module | T: TF module, indexes, rules | Data persists |
| S3.3 BigQuery export | T: scheduled fn, TF module | Daily export runs |

## Epic 4 — CI/CD (Phase 4)

| Story | Tasks | Acceptance |
|-------|-------|------------|
| S4.1 Backend workflow | T: yaml, WIF auth, deploy job | Push → deploy |
| S4.2 Frontend workflow | T: yaml, build, deploy job | Push → deploy |
| S4.3 Infra workflow | T: yaml, plan/apply jobs | PR gets plan comment |
| S4.4 Preview envs | T: PR deploy, cleanup on close | Preview URL in PR |

## Epics 5-8 (Phases 5-8)

| Epic | Key stories | Domain |
|------|-------------|--------|
| 5 Auth | Login, API keys, limits | `internal/auth/` |
| 6 Observability | OTel init, spans, metrics | `internal/telemetry/` |
| 7 Security | Secret store, auth-proxy | `internal/authproxy/` |
| 8 Safety | Categoriser, approval, cleanup | `internal/safety/` |

## Conflict-Free Parallel Development

Each domain owns its own `internal/` package. No cross-domain imports.
Agents work on separate packages; integration via interfaces:

```text
cmd/api/internal/
├── handler/    # Epic 1    ├── auth/       # Epic 5
├── telemetry/  # Epic 6    ├── authproxy/  # Epic 7
├── safety/     # Epic 8    └── events/     # Epic 11
```

## Task Template (GitHub Issue)

```markdown
**Story**: S1.2 Analyze endpoint
**Task**: Implement SSE stream handler
**Domain**: cmd/api/internal/handler/
**Acceptance**: Unit test passes; SSE events match OpenAPI schema
**Definition of Done**:
- [ ] Code in `internal/handler/analyze.go`
- [ ] Unit test in `internal/handler/analyze_test.go`
- [ ] `make all` passes in `cmd/api`
- [ ] OpenAPI spec updated if endpoints change
```

## Steps

1. Create GitHub labels: `epic:1-api`, `epic:2-frontend`, etc.
2. Create epic issues linking to phase docs
3. Break each epic into story issues with acceptance criteria
4. Break stories into task issues with domain boundaries
5. Assign tasks to agents — no two tasks share a file

## Success Criteria

- [ ] Every phase has an epic with ≥ 3 stories
- [ ] Every story has ≥ 2 tasks with acceptance criteria
- [ ] No two tasks modify the same Go package or route file
- [ ] Each task has a clear definition of done
