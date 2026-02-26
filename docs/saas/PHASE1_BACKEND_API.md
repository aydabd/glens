# Phase 1 — Backend REST API

> Wrap the existing Glens CLI into a stateless HTTP API on Cloud Run.

## Requirements Covered

BE-01 (parse spec), BE-02 (AI generate), BE-03 (execute tests), BE-04
(GitHub issues), BE-05 (reports), BE-06 (list models), BE-07 (filter by
op-id), BE-08 (REST API).

## Approach

New module `cmd/api` reuses all `internal/` packages from `cmd/glens`.
Stdlib `net/http`, stateless, SSE for progress streaming.

```
cmd/api/
├── main.go                 # HTTP server
├── internal/handler/       # analyze, models, health
├── internal/middleware/     # CORS, logging, auth stub
├── go.mod / Makefile / Dockerfile / openapi.yaml
```

## API Endpoints

| Method | Path | Maps to |
|--------|------|---------|
| GET | `/healthz` | Liveness probe |
| POST | `/api/v1/analyze` | BE-01, BE-02, BE-03 (SSE stream) |
| GET | `/api/v1/models` | BE-06 |
| POST | `/api/v1/issues` | BE-04 |
| GET | `/api/v1/reports/:id` | BE-05, BE-07 |

## Example: POST `/api/v1/analyze`

```json
{ "spec_url": "https://petstore.swagger.io/v2/swagger.json",
  "ai_models": ["gpt4"], "op_id": "getPetById",
  "run_tests": true, "create_issues": false }
```

Response is SSE: `event: progress` → `event: result`.

## Dockerfile

```dockerfile
FROM golang:1.23-alpine AS build
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
2. Implement health, analyze, models handlers
3. Add SSE streaming for progress events
4. Write `openapi.yaml` for the API itself
5. Add Dockerfile, Makefile, CI workflow

## Success Criteria

- [ ] `GET /healthz` → `200 {"status":"ok"}`
- [ ] `POST /api/v1/analyze` streams SSE progress + result
- [ ] Docker image builds and runs locally
- [ ] All existing `cmd/glens` tests still pass
