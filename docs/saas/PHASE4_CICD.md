# Phase 4 — CI/CD Pipeline (GitHub Actions → GCP)

> Automated build, test, and deploy to GCP via GitHub Actions.

## Requirements Covered

IN-06 (deploy pipeline), IN-07 (preview envs), IN-08 (secrets).

## GCP Auth

Use **Workload Identity Federation** (no service account keys):
WIF Pool → GitHub OIDC Provider → Service Account with deploy perms.

## GitHub Secrets

| Secret | Description |
|--------|-------------|
| `GCP_PROJECT` | GCP project ID |
| `GCP_REGION` | e.g. `us-central1` |
| `WIF_PROVIDER` | Workload Identity Federation provider |
| `SA_EMAIL` | Service account for deployments |

## Workflows

### Backend (`api.yml`) — triggers on `cmd/api/**`

```yaml
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.23' }
      - run: cd cmd/api && make all
  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    steps:
      - uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}
      - run: gcloud builds submit cmd/api/ --tag ...
      - uses: google-github-actions/deploy-cloudrun@v2
        with: { service: glens-api, region: ${{ secrets.GCP_REGION }} }
```

### Frontend (`frontend.yml`) — triggers on `frontend/**`

```yaml
jobs:
  test:
    steps:
      - run: cd frontend && npm ci && npm run lint && npm test
      - run: cd frontend && npx lighthouse-ci
  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    steps:
      - run: cd frontend && npm ci && npm run build
      - uses: google-github-actions/deploy-cloud-functions@v3
        with: { name: glens-frontend, runtime: nodejs20 }
```

### Infrastructure (`infra.yml`) — triggers on `infra/**`

```yaml
jobs:
  plan:
    steps:
      - uses: hashicorp/setup-terraform@v3
      - run: cd infra && terraform init && terraform plan
  apply:
    needs: plan
    if: github.ref == 'refs/heads/main'
    environment: production
    steps:
      - run: cd infra && terraform apply -auto-approve
```

## Preview Environments (IN-07)

On each PR: deploy Cloud Run revision with PR tag, Cloud Functions
with PR suffix. Auto-cleanup on PR close. Terraform plan posted as
PR comment.

## Steps

1. Create `.github/workflows/{api,frontend,infra}.yml`
2. Set up Workload Identity Federation in GCP
3. Add secrets to GitHub repo settings
4. Test full pipeline: push → build → deploy
5. Add PR preview environment logic

## Success Criteria

- [ ] Push to `main` auto-deploys backend + frontend
- [ ] PRs get Terraform plan comments; no SA keys stored
- [ ] Deploy < 5 min; rollback = `git revert` + push
