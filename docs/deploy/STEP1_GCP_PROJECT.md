# Step 1 — Create GCP Projects and Enable APIs

## 1.1 Create two GCP projects

Run these commands in **Cloud Shell** or your local terminal (replace billing account ID):

```bash
# Find your billing account ID
gcloud billing accounts list

BILLING="YOUR_BILLING_ACCOUNT_ID"   # e.g. 012345-ABCDEF-012345

# Dev project
gcloud projects create glens-dev --name="Glens Dev"
gcloud billing projects link glens-dev --billing-account="$BILLING"

# Prod project
gcloud projects create glens-prod --name="Glens Prod"
gcloud billing projects link glens-prod --billing-account="$BILLING"
```

> **Tip:** Use the same region (e.g. `us-central1`) for all services to avoid egress costs.

## 1.2 Enable required APIs

Run for **each** project (`glens-dev` and `glens-prod`):

```bash
for PROJECT in glens-dev glens-prod; do
  gcloud services enable \
    run.googleapis.com \
    artifactregistry.googleapis.com \
    secretmanager.googleapis.com \
    firestore.googleapis.com \
    pubsub.googleapis.com \
    bigquery.googleapis.com \
    cloudtrace.googleapis.com \
    monitoring.googleapis.com \
    iamcredentials.googleapis.com \
    sts.googleapis.com \
    cloudbuild.googleapis.com \
    --project="$PROJECT"
done
```

> `iamcredentials.googleapis.com` and `sts.googleapis.com` are required for
> Workload Identity Federation (OIDC) in Step 2.

## 1.3 Create Artifact Registry repositories

```bash
for PROJECT in glens-dev glens-prod; do
  gcloud artifacts repositories create glens \
    --repository-format=docker \
    --location=us-central1 \
    --project="$PROJECT" \
    --description="Glens container images"
done
```

Images will be stored at:

```text
us-central1-docker.pkg.dev/<project>/glens/api:<tag>
```

## 1.4 Create a Terraform state bucket

Terraform stores remote state in a Cloud Storage bucket.
Create one bucket **per project** (the bucket name must be globally unique):

```bash
for PROJECT in glens-dev glens-prod; do
  gcloud storage buckets create gs://glens-terraform-state-${PROJECT} \
    --location=us-central1 \
    --project="$PROJECT" \
    --uniform-bucket-level-access
done
```

> **Note:** The bucket name `glens-terraform-state` is configured in
> `infra/main.tf`. Update that value to match your actual bucket name
> before running Terraform in Step 4.

## 1.5 Verify

```bash
gcloud projects list --filter="projectId:glens-*"
# Should show glens-dev and glens-prod
```

---

**Next:** [Step 2 — Connect GitHub to GCP via OIDC](STEP2_GITHUB_OIDC.md)
