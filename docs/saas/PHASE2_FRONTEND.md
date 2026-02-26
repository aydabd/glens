# Phase 2 — Frontend Tech Stack

> Clean, performant, WCAG-accessible frontend on Cloud Functions.

## Requirements Covered

FE-01 – FE-12, DB-02 (historical charts).

## Framework Evaluation (Feb 2026)

| Framework | Bundle | SSR/SSG | Cloud Fn fit | A11y | Perf |
|-----------|--------|---------|-------------|------|------|
| **SvelteKit 2** | ~15 KB | ✅ native | ✅ adapter-node | ✅ compiler | ⭐⭐⭐ |
| Flutter Web | ~2 MB | ❌ canvas | ❌ no SSR | ⚠️ manual | ⭐ |
| Next.js 15 | ~90 KB | ✅ | ✅ adapter | ⚠️ plugin | ⭐⭐ |
| Astro 5 | ~0 KB | ✅ islands | ✅ node | ⚠️ plugin | ⭐⭐⭐ |

### Why Not Flutter Web

- **Bundle size** — ~2 MB Dart-to-JS baseline vs ~15 KB SvelteKit.
- **No SSR** — Flutter renders to `<canvas>`, invisible to crawlers.
- **Cloud Functions** — no adapter; cold starts are heavy (~3 s).
- **A11y** — canvas breaks screen readers; WCAG 2.2 AA impossible.
- **SEO** — no semantic HTML; crawlers see an empty `<body>`.

Flutter excels for native mobile (iOS/Android) but is a poor fit for
a serverless web SaaS requiring accessibility and fast cold starts.

### Recommendation: SvelteKit 2

SvelteKit wins: smallest JS, compiler a11y warnings, progressive
enhancement, `adapter-node` for Cloud Functions, SSE support.

## WCAG 2.2 AA Compliance

- **Perceivable** — semantic HTML, compiler a11y lints
- **Operable** — keyboard nav, focus management, skip-links
- **Understandable** — clear errors, form validation
- **Robust** — works without JS; axe-core + Lighthouse in CI

## Project Structure

```
frontend/
├── src/routes/
│   ├── +page.svelte              # Spec upload (FE-01)
│   ├── +layout.svelte            # Shell: nav, skip-link
│   ├── analyze/+page.svelte      # Progress + results (FE-03/04)
│   ├── analyze/approve/+page.svelte # Destructive-test dialog (FE-11)
│   ├── settings/auth/+page.svelte  # Auth config (FE-12)
│   └── reports/[id]/+page.svelte # Reports + charts (DB-02)
├── src/lib/
│   ├── api.ts        # Typed client (auto-gen from OpenAPI)
│   ├── sse.ts        # EventSource helper
│   └── charts.ts     # Chart.js wrapper for result history
├── svelte.config.js  # adapter-node for Cloud Functions
└── package.json
```

## Result Visualisation (DB-02)

Use **Chart.js 4** (tree-shakeable, ~10 KB for bar/line charts):

- Pass/fail trend over time (line chart)
- Endpoint coverage by category (stacked bar)
- AI model accuracy comparison (radar chart)

Data fetched from `GET /api/v1/results?workspace=X&range=30d`.

## Auth Config UI (FE-12, SE-03, SE-04)

Users configure target-API credentials via a form that stores only
**Secret Manager references**, never raw secrets:

1. User enters credential in frontend
2. Frontend sends to `POST /api/v1/secrets` (backend-only)
3. Backend stores in Secret Manager, returns a `ref` path
4. Frontend stores only the `ref`; raw value is never persisted

## Destructive-Test Approval (FE-11)

Before running, the analyze-preview response lists endpoint risks.
The UI shows a modal grouping endpoints by risk level (🟢🟡🔴).
User can batch-approve or deselect individual endpoints.

## Steps

1. `npm create svelte@latest frontend` (skeleton, TypeScript)
2. Install shadcn-svelte 1.x + Chart.js 4
3. Generate API types from backend OpenAPI spec
4. Build pages: upload, analyze, approve, results, auth config
5. Add axe-core tests + Lighthouse CI config

## Success Criteria

- [ ] Lighthouse a11y ≥ 95; all pages keyboard-navigable
- [ ] SSE progress renders in real-time
- [ ] Historical charts display test trends (DB-02)
- [ ] Auth config stores refs only — no secret leaks (SE-04)
- [ ] Bundle < 30 KB gzipped (excl. Chart.js lazy chunk)
