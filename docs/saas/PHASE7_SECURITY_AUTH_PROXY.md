# Phase 7 — Target-API Authentication Proxy

> Securely inject credentials into test requests without leaking.

## Requirements Covered

SE-01 (auth types), SE-02 (custom headers), SE-03 (server-side
secrets), SE-04 (frontend never sees raw creds), FE-12 (auth UI).

## Supported Auth Types (SE-01)

| Type | Config | Injected as |
|------|--------|-------------|
| Bearer token | `credential_ref` | `Authorization: Bearer <val>` |
| API key (header) | `credential_ref` + `header_name` | `X-API-Key: <val>` |
| API key (query) | `credential_ref` + `param_name` | `?api_key=<val>` |
| OAuth2 client creds | `client_id_ref` + `secret_ref` + `token_url` | Auto-token |
| mTLS | `cert_ref` + `key_ref` | TLS client cert |
| Custom headers | map of `header → credential_ref` | Arbitrary |

## Security Flow

```
Frontend                        Backend (Cloud Run)
  │ POST /api/v1/secrets          │
  │ {"name":"tok","value":"sk…"} ►│► Secret Manager (store)
  │ ◄ {"ref":"projects/…/v/1"}    │
  │                               │
  │ POST /api/v1/analyze          │
  │ {"target_auth":{"type":       │
  │   "bearer","credential_ref":  │
  │   "projects/…/v/1"}}        ──►│ Resolve ref → inject header
  │                               │ → call target API
```

**Key rule**: frontend stores only `ref` paths. Raw values never
appear in Firestore, logs, spans, or API responses.

## Secret Resolution (Go)

```go
func resolveSecret(ctx context.Context, ref string) (string, error) {
    client, err := secretmanager.NewClient(ctx)
    if err != nil { return "", fmt.Errorf("resolve: new client: %w", err) }
    defer client.Close()
    result, err := client.AccessSecretVersion(ctx,
        &smpb.AccessSecretVersionRequest{Name: ref})
    if err != nil { return "", fmt.Errorf("resolve: %w", err) }
    return string(result.Payload.Data), nil
}
```

## Custom Gateway Headers (SE-02)

For Kong, Apigee, or custom gateways:

```json
{ "extra_headers": {
    "X-Kong-Key": "ref:projects/p/secrets/kong/versions/1",
    "X-Tenant": "ref:projects/p/secrets/tenant/versions/1" } }
```

Values prefixed `ref:` are resolved server-side. Plain-string values
are passed through unchanged and MUST NOT contain secrets; use `ref:`
for any sensitive header values.

## Frontend Auth Config (FE-12)

1. Auth-type dropdown (Bearer, API key, OAuth2, mTLS, custom)
2. User enters credential → `POST /api/v1/secrets` (once)
3. Backend returns `ref` → stored in workspace config
4. Raw value never in browser storage or cookies

## OAuth2 Client Credentials

Backend handles full flow: resolve `client_id_ref` + `secret_ref`
→ POST to `token_url` → cache token → inject Bearer header.

## Steps

1. Add Secret Manager Go SDK to `cmd/api`
2. Implement `POST /api/v1/secrets` (store + return ref)
3. Implement auth-proxy middleware (resolve → inject)
4. Add OAuth2 client-credentials handler
5. Build frontend auth config page; add mTLS support

## Success Criteria

- [ ] Bearer, API-key, OAuth2, mTLS auth types work end-to-end
- [ ] Raw secrets never in logs, spans, or responses
- [ ] Frontend stores only ref paths, never raw values
- [ ] Custom gateway headers injected correctly
