# Phase 18 — RFC 9457 Problem Details & Enterprise Authentication

> Standardise API errors per RFC 9457 and provide SSO / RBAC / multi-tenant
> auth for enterprise SaaS customers.

## Part A — RFC 9457 Problem Details for HTTP APIs

### Why

[RFC 9457](https://www.rfc-editor.org/rfc/rfc9457) (formerly RFC 7807)
defines a machine-readable format for HTTP error responses. Adopting it:

- Gives clients a single, predictable error schema.
- Allows typed error categories via `type` URIs.
- Carries structured context (which field failed, what limit was hit).
- Is an IETF standard; most API gateways and client libraries understand it.

### Requirements

| ID | Requirement | Notes |
|----|-------------|-------|
| PD-01 | All non-2xx responses use `application/problem+json` | Content-Type header |
| PD-02 | Response body has `type`, `title`, `status`, `detail` | Required fields |
| PD-03 | `instance` field carries the request path | Debugging aid |
| PD-04 | Validation errors include `errors[]` extension member | Per-field details |
| PD-05 | `type` URIs are stable, documented, and resolve to docs | e.g. `/errors/validation` |
| PD-06 | Existing tests updated to assert new error shape | No regressions |

### Error Types

| Type URI suffix | Title | Typical status |
|----------------|-------|----------------|
| `/errors/validation` | Validation Error | 400 |
| `/errors/unauthorized` | Unauthorized | 401 |
| `/errors/forbidden` | Forbidden | 403 |
| `/errors/not-found` | Not Found | 404 |
| `/errors/rate-limit` | Rate Limit Exceeded | 429 |
| `/errors/internal` | Internal Server Error | 500 |

### Response Example

```json
{
  "type": "https://glens.dev/errors/validation",
  "title": "Validation Error",
  "status": 400,
  "detail": "spec_url is required",
  "instance": "/api/v1/analyze"
}
```

### Validation Error with Extension

```json
{
  "type": "https://glens.dev/errors/validation",
  "title": "Validation Error",
  "status": 400,
  "detail": "One or more fields failed validation",
  "instance": "/api/v1/analyze",
  "errors": [
    {"field": "spec_url", "reason": "required"},
    {"field": "models[0]", "reason": "unknown model"}
  ]
}
```

### Implementation Steps

1. Add `ProblemDetail` struct to `cmd/api/internal/handler/`
2. Replace `writeError()` with `writeProblem()`
3. Update all handlers to use `writeProblem()`
4. Update `ErrorResponse` → `ProblemDetail` in `openapi.yaml`
5. Update all existing tests to assert the new error shape
6. Add dedicated tests for each error type

### Success Criteria

- [x] All error responses use `application/problem+json` content type
- [x] All error responses include `type`, `title`, `status`, `detail`, `instance`
- [x] Existing tests pass with updated assertions
- [x] New tests cover each error type

---

## Part B — Enterprise Authentication & Authorisation

### Context

Phase 5 covers Firebase Auth for individual users.
Phase 7 covers target-API credential proxying.

This phase extends both to support **enterprise / B2B customers** who need:

- SSO integration (their own IdP)
- Role-based access control (RBAC) within workspaces
- Organisation-level billing and user management

### Requirements

| ID | Requirement | Notes |
|----|-------------|-------|
| EA-01 | OIDC / SAML SSO for enterprise tenants | Customer brings own IdP |
| EA-02 | RBAC roles: `owner`, `admin`, `member`, `viewer` | Per-workspace |
| EA-03 | Organisation entity above workspaces | Billing + user directory |
| EA-04 | API key scoping per organisation | Keys belong to org, not user |
| EA-05 | Audit log for auth events | Login, role change, key use |
| EA-06 | Session management | Revoke, expiry, refresh |
| EA-07 | MFA support | TOTP or WebAuthn via IdP |
| EA-08 | Rate limiting per org / plan tier | Prevent abuse |

### SSO Flow (OIDC)

```text
Enterprise User ──► Glens Login Page ──► "Sign in with SSO"
                                            │
                                    Redirect to customer IdP
                                    (Okta / Azure AD / Auth0)
                                            │
                                    IdP authenticates user
                                            │
                                    Callback → Glens Backend
                                            │
                              Verify ID token → create/link account
                                            │
                              Issue session (httpOnly cookie + CSRF token)
```

### RBAC Model

```text
Organisation
  ├── owner (1+)          — billing, delete org, manage admins
  ├── admin (0+)          — manage members, API keys, workspaces
  └── Workspace
       ├── admin          — configure workspace, approve endpoints
       ├── member         — run analyses, view reports
       └── viewer         — read-only access to reports
```

### Data Model Extension (Firestore)

```text
organisations/{orgId}
  ├── name, plan, sso_config, billing_email
  ├── members/{userId}    — role, joined_at, invited_by
  └── workspaces/{wsId}   — (existing model from Phase 5)
       └── runs/{runId}

sso_configs/{orgId}
  ├── provider: "oidc" | "saml"
  ├── issuer_url
  ├── client_id_ref       — Secret Manager ref
  ├── client_secret_ref   — Secret Manager ref
  └── allowed_domains[]   — email domain allow-list
```

### Auth Middleware Stack

```text
Request
  │
  ├── 1. Rate Limiter       — per IP / per org / per plan
  ├── 2. Auth (JWT verify)  — Firebase or OIDC token
  ├── 3. Org Resolver       — map user → org
  ├── 4. RBAC Check         — role ≥ required for route
  └── 5. Handler
```

### API Key Auth (machine-to-machine)

```text
POST /api/v1/analyze
Authorization: Bearer glens_sk_live_abc123...

Backend:
  1. Hash key → lookup in Firestore api_keys collection
  2. Resolve org_id from key
  3. Check key scope includes "analyze:write"
  4. Enforce org rate limit
  5. Proceed to handler
```

### RBAC Permission Matrix

| Action | Owner | Admin | Member | Viewer |
|--------|-------|-------|--------|--------|
| Manage org settings | ✅ | ❌ | ❌ | ❌ |
| Manage billing | ✅ | ❌ | ❌ | ❌ |
| Add/remove admins | ✅ | ❌ | ❌ | ❌ |
| Add/remove members | ✅ | ✅ | ❌ | ❌ |
| Create API keys | ✅ | ✅ | ❌ | ❌ |
| Create workspace | ✅ | ✅ | ❌ | ❌ |
| Run analysis | ✅ | ✅ | ✅ | ❌ |
| View reports | ✅ | ✅ | ✅ | ✅ |
| Approve endpoints | ✅ | ✅ | ✅ | ❌ |
| Configure SSO | ✅ | ❌ | ❌ | ❌ |

### Plan Tiers (extended from Phase 5)

| Feature | Free | Pro | Enterprise |
|---------|------|-----|------------|
| Users | 1 | 10 | Unlimited |
| Workspaces | 1 | 5 | Unlimited |
| Analyses/month | 10 | 500 | Custom |
| AI models | GPT-4o-mini | All | All + custom |
| SSO | ❌ | ❌ | ✅ |
| RBAC | ❌ | Basic | Full |
| Audit log | ❌ | 7 days | 90 days |
| API keys | 1 | 10 | Unlimited |
| SLA | — | — | 99.9% |
| Support | Community | Email | Dedicated |

### Implementation Phases (sub-phases)

| Step | Description | Depends on |
|------|-------------|------------|
| 18a | RFC 9457 error responses | Phase 1 (done) |
| 18b | RBAC types + middleware | Phase 5 |
| 18c | Organisation entity + data model | Phase 5, Phase 3 |
| 18d | OIDC SSO integration | 18c |
| 18e | SAML SSO integration | 18d |
| 18f | API key scoping per org | 18c |
| 18g | Audit logging | 18c, Phase 11 |
| 18h | Rate limiting per org/plan | 18c |

### Implementation Steps

1. **18a** — Implement RFC 9457 in `cmd/api/internal/handler/` (this PR)
2. **18b** — Add RBAC types (`Role`, `Permission`) + `AuthzMiddleware`
3. **18c** — Add `Organisation` Firestore model + CRUD endpoints
4. **18d** — Integrate OIDC provider (generic — works with Okta, Azure AD, Auth0)
5. **18e** — Add SAML 2.0 SP support (optional, enterprise-only)
6. **18f** — Scope API keys to organisations; enforce in auth middleware
7. **18g** — Emit auth events to Pub/Sub; store audit log in Firestore
8. **18h** — Token-bucket rate limiter keyed by org ID and plan tier

### Success Criteria

- [ ] RFC 9457 error responses on all endpoints (18a — this PR)
- [ ] RBAC roles enforced; viewers cannot mutate (18b)
- [ ] Organisation CRUD; users belong to orgs (18c)
- [ ] OIDC SSO login works with at least one IdP (18d)
- [ ] API keys scoped per org; revocation immediate (18f)
- [ ] Auth events in audit log; queryable by org (18g)
- [ ] Rate limits applied per org and plan (18h)
