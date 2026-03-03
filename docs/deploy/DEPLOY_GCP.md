# Deploy Glens to GCP — Overview

This guide walks you through deploying Glens as a SaaS service on Google Cloud Platform.
After completing all steps you will have:

- A live API reachable at a friendly URL (e.g. `https://api.glens.example.com`)
- Automated CI/CD that re-deploys on every push to `main`
- Separate **dev** and **prod** environments, each isolated in its own GCP project
- No long-lived service-account keys — GitHub connects to GCP via OIDC

## Steps

| Step | Document | What you do |
|------|----------|-------------|
| 1 | [STEP1_GCP_PROJECT.md](STEP1_GCP_PROJECT.md) | Create GCP projects and enable APIs |
| 2 | [STEP2_GITHUB_OIDC.md](STEP2_GITHUB_OIDC.md) | Connect GitHub → GCP without keys (OIDC) |
| 3 | [STEP3_GITHUB_SECRETS.md](STEP3_GITHUB_SECRETS.md) | Configure GitHub environments, secrets & variables |
| 4 | [STEP4_TERRAFORM.md](STEP4_TERRAFORM.md) | Provision all GCP resources with Terraform |
| 5 | [STEP5_DEPLOY_API.md](STEP5_DEPLOY_API.md) | Build and deploy the container to Cloud Run |
| 6 | [STEP6_CUSTOM_DOMAIN.md](STEP6_CUSTOM_DOMAIN.md) | Map a friendly domain name to the service |

Work through the steps in order — each step depends on the previous one.

## Prerequisites

- A Google Cloud account with billing enabled
- A domain name you control (for the friendly URL in Step 6)
- `gcloud` CLI installed locally ([install guide](https://cloud.google.com/sdk/docs/install))
- Terraform ≥ 1.10 installed locally ([install guide](https://developer.hashicorp.com/terraform/install))
- Docker installed locally ([install guide](https://docs.docker.com/get-docker/))
- Admin access to the GitHub repository

## Architecture at a Glance

```text
GitHub Actions (push to main)
  │
  ├─ OIDC token ──► Workload Identity Federation ──► GCP Service Account
  │
  ├─ docker build & push ──► Artifact Registry (us-central1-docker.pkg.dev)
  │
  └─ terraform apply / gcloud run deploy
           │
           └─► Cloud Run (glens-api) ──► Firestore · Secret Manager · Pub/Sub
                    │
           Custom Domain (api.glens.example.com)
```

## Environment Mapping

| GitHub environment | Terraform workspace | GCP project |
|-------------------|---------------------|-------------|
| `dev` | `dev` | `glens-dev` |
| `prod` | `prod` | `glens-prod` |

## Time Estimate

| Step | Approx. time |
|------|-------------|
| Steps 1–3 (one-time GCP + GitHub setup) | 30–45 min |
| Step 4 (Terraform apply) | 5–10 min |
| Step 5 (first deploy) | 5 min |
| Step 6 (custom domain) | 10–15 min + DNS propagation (up to 24 h) |
