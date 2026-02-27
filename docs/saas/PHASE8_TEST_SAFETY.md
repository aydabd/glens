# Phase 8 — Destructive Test Safety

> Categorise, warn, and get user approval before running write tests.

## Requirements Covered

TS-01 (categorise), TS-02 (warn), TS-03 (batch approve), TS-04 (cleanup).

## Problem

Glens runs tests against real APIs. Write/delete endpoints can
corrupt production data. Users must be warned and explicitly approve
destructive tests before execution.

## Endpoint Categorisation (TS-01)

| Category | Methods | Risk |
|----------|---------|------|
| `read` | GET, HEAD, OPTIONS | 🟢 safe |
| `write` | POST (create) | 🟡 medium |
| `mutate` | PUT, PATCH | 🟡 medium |
| `destroy` | DELETE | 🔴 high |

POST with `/search`, `/query`, or `/list` in path → `read`.
OpenAPI `x-safe: true` extension overrides default.

## Preview API (before execution)

`POST /api/v1/analyze/preview` returns categories without running:

```json
{ "endpoints": [
    { "path": "/pets", "method": "GET", "category": "read" },
    { "path": "/pets", "method": "POST", "category": "write" },
    { "path": "/pets/{id}", "method": "DELETE", "category": "destroy" }
  ],
  "warnings": ["1 endpoint will DELETE data — irreversible"] }
```

## User Approval Flow (TS-02, TS-03)

```
"Analyze" clicked → POST /api/v1/analyze/preview
  → Frontend: approval modal (FE-11)
    🟢 Read-only — auto-approved
    🟡 Write/Mutate — needs approval
    🔴 Destroy — explicit approval required
  → User: approve/deselect per endpoint
    "Approve All" for dev environments
  → POST /api/v1/analyze
    { "approved_endpoints": [...], "skipped_endpoints": [...] }
```

## Cleanup Hooks (TS-04)

- **Reverse ops** — POST creates → generate DELETE for cleanup
- **Transaction IDs** — tag created resources for batch cleanup
- **Cleanup endpoint** — `POST /api/v1/runs/{id}/cleanup`
- **Dry-run mode** — generate tests but skip target API calls

Cleanup is best-effort. UI shows post-run cleanup summary.

## Implementation (Go)

```go
func categorise(ep *parser.Endpoint) EndpointCategory {
    switch strings.ToUpper(ep.Method) {
    case "GET", "HEAD", "OPTIONS":
        return EndpointCategory{Category: "read", Risk: "safe"}
    case "POST":
        if isSafePost(ep) { return EndpointCategory{Category: "read", Risk: "safe"} }
        return EndpointCategory{Category: "write", Risk: "medium"}
    case "PUT", "PATCH":
        return EndpointCategory{Category: "mutate", Risk: "medium"}
    case "DELETE":
        return EndpointCategory{Category: "destroy", Risk: "high"}
    default:
        return EndpointCategory{Category: "write", Risk: "medium"}
    }
}
```

## Frontend Approval Modal (FE-11)

- Grouped by risk: 🟢 → 🟡 → 🔴; checkboxes per endpoint
- "Approve All" / "Safe Only" quick-select buttons
- Warning: "⚠ X endpoints will modify the target database"
- Environment indicator: dev / staging / production

## Steps

1. Add `internal/safety/categoriser.go` to `cmd/api`
2. Implement `POST /api/v1/analyze/preview`; update analyze handler
3. Add cleanup hook generator; build approval modal (FE-11)

## Success Criteria

- [ ] All endpoints categorised by method + heuristics
- [ ] User sees warnings before destructive tests run
- [ ] Batch approve/deselect; cleanup removes test-created resources
