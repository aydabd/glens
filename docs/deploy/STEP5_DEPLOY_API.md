# Step 5 — Build and Deploy the API Container

## 5.1 Build the container image locally

From the **repository root**:

```bash
PROJECT=glens-dev          # use glens-prod for prod
REGION=us-central1
TAG=v1.0.0                 # use a semver tag or 'latest'
IMAGE="${REGION}-docker.pkg.dev/${PROJECT}/glens/api:${TAG}"

docker build \
  -f cmd/api/Dockerfile \
  -t "$IMAGE" \
  .
```

The multi-stage `Dockerfile` compiles the Go binary and produces a minimal
distroless image (~10 MB).

## 5.2 Push the image to Artifact Registry

```bash
# Authenticate Docker with GCP
gcloud auth configure-docker "${REGION}-docker.pkg.dev" --project="$PROJECT"

docker push "$IMAGE"
```

## 5.3 Deploy to Cloud Run

```bash
gcloud run deploy glens-api \
  --image="$IMAGE" \
  --region="$REGION" \
  --project="$PROJECT" \
  --platform=managed \
  --allow-unauthenticated \
  --set-env-vars="LOG_LEVEL=info"
```

> Terraform already created the Cloud Run service — this command updates the
> running image to your new tag without re-running Terraform.

## 5.4 Verify the deployment

```bash
gcloud run services describe glens-api \
  --region="$REGION" --project="$PROJECT" \
  --format="value(status.url)"
# e.g. https://glens-api-<hash>-uc.a.run.app
```

Open that URL in a browser — you should see the Glens API health response.

## 5.5 Automated CI/CD (GitHub Actions)

Pushing to `main` with changes under `cmd/api/**` triggers `.github/workflows/api.yml` which:

1. Runs `make all` (fmt, vet, lint, test)
2. Authenticates to GCP using the OIDC token + `WIF_PROVIDER` secret
3. Builds and pushes the image to Artifact Registry
4. Runs `gcloud run deploy` with the new image tag

### Add the deploy step to `api.yml`

The workflow at `.github/workflows/api.yml` currently only runs tests.
Add the following `deploy` job to enable automatic deployment
(use `environment: prod` for a separate prod deploy job):

```yaml
  deploy:
    needs: test
    if: github.ref == 'refs/heads/main'
    runs-on: ubuntu-latest
    environment: dev          # change to 'prod' for a prod deploy job
    permissions:
      id-token: write
      contents: read
    steps:
      - uses: actions/checkout@v4
      - uses: google-github-actions/auth@v2
        with:
          workload_identity_provider: ${{ secrets.WIF_PROVIDER }}
          service_account: ${{ secrets.SA_EMAIL }}
      - uses: google-github-actions/setup-gcloud@v2
      - run: gcloud auth configure-docker ${{ vars.GCP_REGION }}-docker.pkg.dev --quiet
      - run: |
          IMAGE="${{ vars.GCP_REGION }}-docker.pkg.dev/${{ vars.GCP_PROJECT }}/glens/api:${{ github.sha }}"
          docker build -f cmd/api/Dockerfile -t "$IMAGE" .
          docker push "$IMAGE"
          gcloud run deploy glens-api --image="$IMAGE" \
            --region="${{ vars.GCP_REGION }}" --project="${{ vars.GCP_PROJECT }}" \
            --platform=managed --allow-unauthenticated --quiet
```

---

**Next:** [Step 6 — Map a Friendly Domain Name](STEP6_CUSTOM_DOMAIN.md)
