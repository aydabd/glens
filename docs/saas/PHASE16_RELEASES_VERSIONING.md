# Phase 16 — Semantic Versioning & Release Pipeline

> Pre-release to dev, promote to prod with regression gates.

## Requirements Covered

RV-01 (semver), RV-02 (pre-release to dev), RV-03 (promote to prod),
RV-04 (regression gate), RV-05 (smoke tests), RV-06 (deployment e2e).

## Version Format

```
v1.2.3          — production release
v1.2.3-rc.1     — release candidate (deployed to dev)
v1.2.3-alpha.1  — early preview (CI only, no deploy)
```

Tags trigger builds. No manual version bumps — use `release-please`.

## Release Flow

```
Feature PR → main (CI: unit + integration)
    │
    ├── release-please creates Release PR (bumps version)
    │
    └── Merge Release PR → tag v1.2.3-rc.1
          │
          ├── Deploy to dev (auto)
          ├── Run regression tests against dev
          ├── Run smoke tests against dev
          ├── Run e2e business flow tests
          │
          └── All pass → promote: tag v1.2.3
                │
                ├── Deploy to prod (auto)
                ├── Run smoke tests against prod
                └── Monitor error rate for 15 min
```

## GitHub Actions Workflow

```yaml
on: { push: { tags: ['v*'] } }
jobs:
  release:
    steps:
      - run: |
          [[ "${{ github.ref_name }}" == *-rc* ]] \
            && echo "ENV=dev" >> $GITHUB_ENV \
            || echo "ENV=prod" >> $GITHUB_ENV
      - uses: google-github-actions/auth@v2
      - run: |
          cd infra && terraform apply \
            -var-file=environments/$ENV.tfvars \
            -var="image_tag=${{ github.ref_name }}" -auto-approve
```

## Regression Testing Gate (RV-04)

After dev deploy, before prod promotion:

| Suite | Target | Pass criteria |
|-------|--------|---------------|
| Unit | CI | 100% pass |
| Integration | Emulators | 100% pass |
| Smoke | Dev environment | Health + core flow |
| E2E business | Dev environment | All scenarios pass |
| E2E deployment | Dev environment | Deploy + rollback ok |
| E2E events | Dev environment | Pub/Sub round-trip |

## Smoke Tests (RV-05)

```bash
BASE_URL="${1:-https://dev.glens.app}"
curl -sf "$BASE_URL/healthz" || exit 1
curl -sf "$BASE_URL/api/v1/models" | jq '.models | length > 0' || exit 1
```

## Deployment E2E (RV-06)

1. Deploy rc → verify health → full analysis flow → Firestore write
2. Verify Pub/Sub event → trigger rollback → re-deploy → clean state

Rollback = route traffic to previous Cloud Run revision (instant).

## Steps

1. Configure `release-please` + tag-triggered deploy workflow
2. Write smoke script + regression/deployment e2e suites
3. Add post-deploy monitoring (error rate < 0.1%)

## Success Criteria

- [ ] `release-please` creates versioned release PRs
- [ ] RC tags auto-deploy to dev; stable tags to prod
- [ ] Regression suite gates prod promotion
- [ ] Rollback restores previous version in < 2 min
