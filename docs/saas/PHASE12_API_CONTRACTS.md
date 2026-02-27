# Phase 12 — API Contracts & Swagger Publishing

> OpenAPI spec as single source of truth, published via API Gateway.

## Requirements Covered

AC-01 (contract-first), AC-02 (Swagger UI on Cloud Function),
AC-03 (auto-generated SDKs), AC-04 (contract tests in CI).

## Contract-First Development

Every endpoint defined in `openapi.yaml` **before** code. Code gen
ensures handlers match the spec; compiler fails on drift.

```text
cmd/api/openapi.yaml      # source of truth
├── paths/                 # endpoint definitions
├── components/schemas/    # shared models + event payloads
└── components/securitySchemes/
```

## Code Generation

```text
openapi.yaml
  ├──► oapi-codegen (Go)     → server interfaces + types
  ├──► openapi-typescript     → frontend API client types
  └──► Swagger UI             → hosted docs (Cloud Function)
```

### Go Server

```bash
oapi-codegen -package api -generate types,server \
  cmd/api/openapi.yaml > cmd/api/internal/api/openapi_gen.go
```

Handlers implement generated interface. Compiler fails if endpoint
exists in spec but has no handler implementation.

### Frontend Client

```bash
npx openapi-typescript cmd/api/openapi.yaml \
  -o frontend/src/lib/api-types.ts
```

Type-safe fetch; TS compiler catches spec/client drift.

## Swagger UI (AC-02)

Cloud Function serves Swagger UI pointing at the live spec.
Live at `https://api.glens.dev/docs` — always matches deployed API.

## API Gateway (GCP)

```hcl
resource "google_api_gateway_api" "glens" { api_id = "glens-api" }
resource "google_api_gateway_api_config" "v1" {
  api = google_api_gateway_api.glens.api_id
  openapi_documents {
    document { path = "openapi.yaml"
               contents = filebase64("cmd/api/openapi.yaml") }
  }
}
```

Gateway validates requests against spec — free validation layer.

## Contract Tests in CI (AC-04)

```yaml
contract-check:
  steps:
    - run: oapi-codegen ... > /tmp/gen.go && diff /tmp/gen.go openapi_gen.go
    - run: npx openapi-typescript ... -o /tmp/t.ts && diff /tmp/t.ts api-types.ts
```

PR blocked if spec and code diverge.

## Steps

1. Write `cmd/api/openapi.yaml` with all endpoints + schemas
2. Set up `oapi-codegen` + `openapi-typescript` generation
3. Deploy Swagger UI Cloud Function at `/docs`
4. Add API Gateway Terraform module
5. Add contract-diff CI check

## Success Criteria

- [ ] `openapi.yaml` is single source of truth
- [ ] Go server interfaces auto-generated; compiler catches drift
- [ ] Frontend types auto-generated; TS catches drift
- [ ] Swagger UI live at `/docs`; CI blocks spec/code divergence
