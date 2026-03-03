# Step 2 — Connect GitHub to GCP via OIDC (No Keys Required)

GitHub Actions authenticates to GCP using **Workload Identity Federation**.
No service-account JSON keys are stored anywhere — GitHub presents a short-lived OIDC
token that GCP validates automatically.

## 2.1 Create the Workload Identity Pool

Run once **per GCP project** (dev and prod):

```bash
PROJECT=glens-dev   # repeat with glens-prod

gcloud iam workload-identity-pools create github-pool \
  --location=global \
  --project="$PROJECT" \
  --display-name="GitHub Actions Pool"
```

## 2.2 Add the GitHub OIDC provider to the pool

```bash
REPO="aydabd/glens"   # change to your GitHub org/repo

gcloud iam workload-identity-pools providers create-oidc github-provider \
  --location=global \
  --workload-identity-pool=github-pool \
  --project="$PROJECT" \
  --display-name="GitHub OIDC" \
  --issuer-uri="https://token.actions.githubusercontent.com" \
  --attribute-mapping="google.subject=assertion.sub,attribute.repository=assertion.repository" \
  --attribute-condition="assertion.repository == '${REPO}'"
```

## 2.3 Create a deployment service account

```bash
gcloud iam service-accounts create glens-deploy \
  --project="$PROJECT" \
  --display-name="Glens GitHub Deploy SA"
```

## 2.4 Grant required roles to the service account

```bash
SA="glens-deploy@${PROJECT}.iam.gserviceaccount.com"

for ROLE in \
  roles/run.admin \
  roles/artifactregistry.writer \
  roles/storage.admin \
  roles/secretmanager.admin \
  roles/datastore.owner \
  roles/pubsub.admin \
  roles/bigquery.admin \
  roles/iam.serviceAccountTokenCreator; do
  gcloud projects add-iam-policy-binding "$PROJECT" \
    --member="serviceAccount:$SA" --role="$ROLE"
done
```

## 2.5 Allow GitHub Actions to impersonate the service account

```bash
POOL_ID=$(gcloud iam workload-identity-pools describe github-pool \
  --location=global --project="$PROJECT" --format="value(name)")

gcloud iam service-accounts add-iam-policy-binding "$SA" \
  --project="$PROJECT" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${POOL_ID}/attribute.repository/${REPO}"
```

## 2.6 Note the provider resource name

```bash
gcloud iam workload-identity-pools providers describe github-provider \
  --location=global \
  --workload-identity-pool=github-pool \
  --project="$PROJECT" \
  --format="value(name)"
# Example output:
# projects/123456789/locations/global/workloadIdentityPools/github-pool/providers/github-provider
```

Save this string — you will use it as the `WIF_PROVIDER` secret in Step 3.

---

**Next:** [Step 3 — GitHub Environments, Secrets & Variables](STEP3_GITHUB_SECRETS.md)
