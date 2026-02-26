# Phase 2 — Frontend Tech Stack

> Modern, accessible frontend deployed as a Cloud Function (SSR).

## Requirements Covered

FE-01 (spec upload), FE-02 (model select), FE-03 (live progress),
FE-04 (results dashboard), FE-05 (GitHub issues), FE-06 (filter),
FE-07 (download reports), FE-08 (WCAG 2.2 AA), FE-09 (responsive).

## Why SvelteKit

| Framework | Bundle | SSR | A11y built-in | Complexity |
|-----------|--------|-----|---------------|------------|
| **SvelteKit** | ~15 KB | ✅ | ✅ compiler warns | Low ✅ |
| Next.js | ~90 KB | ✅ | ⚠️ plugin | Medium |
| Astro | ~0 KB | ✅ | ⚠️ plugin | Low |

SvelteKit wins: smallest JS, compiler a11y warnings, progressive
enhancement (works without JS), `adapter-node` for Cloud Functions.

## WCAG 2.2 AA Compliance

- **Perceivable** — semantic HTML, compiler a11y warnings
- **Operable** — keyboard nav, focus management, skip-links
- **Understandable** — clear errors, form validation
- **Robust** — works without JS; axe-core + Lighthouse CI in pipeline

## Project Structure

```
frontend/
├── src/routes/
│   ├── +page.svelte              # Home — spec upload (FE-01)
│   ├── +layout.svelte            # Shell: nav, skip-link, footer
│   ├── analyze/+page.svelte      # Progress + results (FE-03, FE-04)
│   ├── models/+page.svelte       # Model selector (FE-02)
│   └── reports/[id]/+page.svelte # Report view + download (FE-07)
├── src/lib/
│   ├── api.ts          # Typed client (auto-gen from BE OpenAPI spec)
│   └── sse.ts          # EventSource helper (FE-03)
├── svelte.config.js    # adapter-node for Cloud Functions
├── package.json
└── Dockerfile
```

## UI Components

Use **shadcn-svelte** — accessible, unstyled primitives with ARIA
support, keyboard nav, and CSS-variable theming (light/dark).

## API Client

Auto-generated from backend OpenAPI spec for type safety:

```bash
npx openapi-typescript cmd/api/openapi.yaml \
  -o frontend/src/lib/api-types.ts
```

## Deployment

SvelteKit + `adapter-node` → Node.js server on Cloud Functions:

```bash
npm run build
gcloud functions deploy glens-frontend \
  --gen2 --runtime=nodejs20 --source=build/ \
  --entry-point=handler --region=us-central1
```

## Steps

1. `npm create svelte@latest frontend` (skeleton, TypeScript)
2. Install shadcn-svelte + configure adapter-node
3. Generate API types from backend OpenAPI spec
4. Build pages: home (upload), analyze (SSE), reports (download)
5. Add axe-core tests + Lighthouse CI config

## Success Criteria

- [ ] Lighthouse a11y score ≥ 95
- [ ] All pages keyboard-navigable
- [ ] SSE progress renders in real-time
- [ ] Works without JavaScript (progressive enhancement)
- [ ] Bundle < 30 KB gzipped
