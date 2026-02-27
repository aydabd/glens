# Phase 5 — Authentication & Multi-Tenancy

> Secure, multi-tenant SaaS using Firebase Auth and Firestore.

## Requirements Covered

BE-09 (auth & API keys), BE-10 (multi-tenant isolation), FE-10 (login).

## Auth Flow

```
User ── Login ──► Firebase Auth ──► JWT (ID token)
  │                                       │
  └──► Frontend (httpOnly cookie) ──► Backend (verify via Admin SDK)
                                           │
                                      Firestore (scoped by user_id)
```

## Data Model (Firestore)

```
users/{userId}         — email, plan, settings
workspaces/{wsId}      — owner_id, members[]
  └── runs/{runId}     — spec_url, models, status, results
  └── api_keys/{keyId} — name, key_hash, last_used
```

## Security Rules (excerpt)

```javascript
match /databases/{database}/documents {
  match /workspaces/{wsId} {
    allow read: if request.auth.uid == resource.data.owner_id
                || request.auth.uid in resource.data.members;
    allow write: if request.auth.uid == resource.data.owner_id;
    match /runs/{runId} {
      allow read, write: if request.auth.uid ==
        get(/databases/$(database)/documents/workspaces/$(wsId)).data.owner_id;
    }
  }
}
```

## Backend Auth Middleware (Go)

```go
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := extractToken(r) // reads httpOnly cookie or Authorization header
        if token == "" {
            http.Error(w, "unauthorized", http.StatusUnauthorized)
            return
        }
        decoded, err := firebaseAuth.VerifyIDToken(r.Context(), token)
        if err != nil {
            http.Error(w, "invalid token", http.StatusUnauthorized)
            return
        }
        ctx := context.WithValue(r.Context(), userIDKey, decoded.UID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

## API Key Endpoints

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/v1/keys` | Create (returns raw key once) |
| GET | `/api/v1/keys` | List keys |
| DELETE | `/api/v1/keys/:id` | Revoke key |

Keys are bcrypt-hashed before storage.

## Free vs Pro Plan

| Feature | Free | Pro |
|---------|------|-----|
| Analyses/month | 10 | Unlimited |
| AI models | GPT-4o-mini | All |
| Team members | 1 | 10 |
| Report retention | 7 days | 90 days |
| API keys | 1 | 10 |

Plan enforcement is in backend middleware, not frontend.

## Steps

1. Enable Firebase Auth + configure GitHub/Google OAuth providers
2. Add Firebase Admin SDK; implement auth middleware
3. Design Firestore model + security rules
4. Build login/signup pages in SvelteKit
5. Implement API key management; add usage tracking and plan limits

## Success Criteria

- [ ] Sign up / log in via GitHub or Google
- [ ] Isolated workspace data; security rules block cross-tenant access
- [ ] API keys work; free tier limits enforced; auth < 50 ms
