# Step 4 — Provision GCP Infrastructure with Terraform

Terraform creates every GCP resource: Cloud Run service, Firestore database,
Pub/Sub topics, Secret Manager secrets, BigQuery dataset, and storage buckets.

## 4.1 Update the Terraform state bucket name

Open `infra/main.tf` and set the bucket name to match what you created in Step 1.4:

```hcl
backend "gcs" {
  bucket = "glens-terraform-state-glens-dev"  # use glens-terraform-state-glens-prod for prod
  prefix = "terraform/state"
}
```

> **Tip:** You can also pass `--backend-config` at init time to avoid editing the file.

## 4.2 Authenticate locally

```bash
gcloud auth application-default login
```

## 4.3 Initialise Terraform

```bash
cd infra
terraform init
```

## 4.4 Deploy dev environment

```bash
terraform workspace new dev 2>/dev/null || terraform workspace select dev
terraform plan -out=dev.tfplan
terraform apply dev.tfplan
```

Terraform prints the API URL when done:

```text
Outputs:
api_url = "https://glens-api-<hash>-uc.a.run.app"
```

Save this URL — you will map a friendly domain to it in Step 6.

## 4.5 Deploy prod environment

```bash
# Switch the state bucket to the prod bucket first (edit infra/main.tf or use -backend-config)
terraform init -reconfigure \
  -backend-config="bucket=glens-terraform-state-glens-prod"

terraform workspace new prod 2>/dev/null || terraform workspace select prod
terraform plan -out=prod.tfplan
terraform apply prod.tfplan
```

## 4.6 What Terraform creates

| Resource | Purpose |
|----------|---------|
| `google_cloud_run_v2_service.api` | Runs the Glens API container |
| `google_firestore_database.main` | Stores workspace runs and results |
| `google_storage_bucket.reports` | Stores generated reports |
| `google_project_service.*` | Enables required GCP APIs |
| `google_monitoring_notification_channel` | Sends alerts (prod only) |

## 4.7 CI auto-apply (optional)

Pushing changes to `infra/**` on `main` triggers `.github/workflows/infra.yml`
which runs `terraform apply` automatically for both workspaces.
The `prod` apply is gated by the GitHub environment approval configured in Step 3.

---

**Next:** [Step 5 — Build and Deploy the API Container](STEP5_DEPLOY_API.md)
