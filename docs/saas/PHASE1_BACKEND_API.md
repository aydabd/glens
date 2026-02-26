# Phase 1 — Backend REST API + MCP

> Stateless HTTP API on Cloud Run with MCP tooling support.

## Requirements Covered

BE-01 – BE-08, SE-01 – SE-02 (auth proxy), TS-01 (categorisation).

## Approach

New module `cmd/api` reuses `internal/` packages. Stdlib `net/http`,
SSE for streaming, JSON-RPC for MCP tool calls.

```
cmd/api/
├── main.go
├── internal/handler/       # analyze, models, health, mcp
├── internal/middleware/     # CORS, logging, auth, otel, authproxy
├── internal/safety/         # endpoint categoriser (read vs write)
├── go.mod / Makefile / Dockerfile / openapi.yaml
```

## API Endpoints

| Method | Path | Maps to |
|--------|------|---------|
| GET | `/healthz` | Liveness probe |
| POST | `/api/v1/analyze` | BE-01 – BE-03 (SSE stream) |
| GET | `/api/v1/models` | BE-06 |
| POST | `/api/v1/issues` | BE-04 |
| GET | `/api/v1/reports/:id` | BE-05, BE-07 |
| POST | `/api/v1/mcp` | MCP JSON-RPC (tool/resource calls) |
| POST | `/api/v1/analyze/preview` | TS-01 dry-run categorisation |

## MCP Integration

Expose Glens as an MCP tool server so AI agents can call it:

```json
{ "jsonrpc": "2.0", "method": "tools/call",
  "params": { "name": "analyze_spec",
    "arguments": { "spec_url": "...", "models": ["gpt4"] } } }
```

## Auth-Proxy Headers (SE-01, SE-02)

The frontend sends a `credential_ref` (Secret Manager path); the
backend resolves it server-side and injects headers into test requests:

```json
{ "spec_url": "...", "target_auth": {
    "type": "bearer", "credential_ref": "projects/x/secrets/tok/v/1"
  }, "extra_headers": { "X-Kong-Key": "ref:projects/x/secrets/k/v/1" }
}
```

Raw secrets **never** leave the backend.

## Endpoint Categorisation (TS-01)

Before test execution, the analyzer classifies endpoints:

| HTTP method | Category | Risk |
|-------------|----------|------|
| GET, HEAD, OPTIONS | `read` | 🟢 safe |
| POST (create) | `write` | 🟡 warn |
| PUT, PATCH | `mutate` | 🟡 warn |
| DELETE | `destroy` | 🔴 danger |

Response includes categories so the frontend can prompt the user.

## Dockerfile (Go 1.24)

```dockerfile
FROM golang:1.24-alpine AS build
WORKDIR /src
COPY . .
RUN cd cmd/api && go build -o /app .

FROM gcr.io/distroless/static
COPY --from=build /app /app
EXPOSE 8080
ENTRYPOINT ["/app"]
```

## Steps

1. `go mod init glens/tools/api` + add to `go.work`
2. Implement health, analyze, models, mcp handlers
3. Add auth-proxy middleware (resolves Secret Manager refs)
4. Add endpoint categoriser in `internal/safety/`
5. Write `openapi.yaml`, Dockerfile, Makefile, CI workflow

## Success Criteria

- [ ] `POST /api/v1/analyze` streams SSE progress + result
- [ ] `POST /api/v1/mcp` handles JSON-RPC tool calls
- [ ] Auth-proxy resolves credential refs without leaking secrets
- [ ] `/api/v1/analyze/preview` returns endpoint risk categories
- [ ] Docker image builds; existing `cmd/glens` tests pass
