# Phase 3 — GCP Infrastructure (IaC)

> All GCP resources defined in Terraform. Targets GCP free tier.

## Requirements Covered

IN-01 (GCP project), IN-02 (Terraform IaC), IN-03 (Cloud Run),
IN-04 (Cloud Functions), IN-05 (Cloud Storage).

## GCP Free Tier Fit

| Service | Free allowance | Usage |
|---------|---------------|-------|
| Cloud Run | 2M req/mo, 360K vCPU-sec | Backend API |
| Cloud Functions | 2M invocations/mo | Frontend SSR |
| Cloud Storage | 5 GB | Static assets |
| Firestore | 1 GB, 50K reads/day | User data |
| Secret Manager | 6 active versions | API keys |
| Firebase Auth | 50K MAU | Auth |

All within free tier for < 1K users.

## Terraform Layout

```
infra/
├── main.tf / variables.tf / outputs.tf
├── modules/
│   ├── api/        # Cloud Run backend
│   ├── frontend/   # Cloud Functions frontend
│   ├── storage/    # Cloud Storage buckets
│   ├── firestore/  # Database
│   └── secrets/    # Secret Manager
└── environments/
    ├── dev.tfvars
    └── prod.tfvars
```

## Core Resources (abbreviated HCL)

```hcl
# Cloud Run — backend API
resource "google_cloud_run_v2_service" "api" {
  name     = "glens-api"
  location = var.region
  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/glens/api:latest"
      ports { container_port = 8080 }
      resources { limits = { cpu = "1", memory = "512Mi" } }
    }
    scaling { min_instance_count = 0; max_instance_count = 3 }
  }
}

# Cloud Functions — frontend SSR
resource "google_cloudfunctions2_function" "frontend" {
  name     = "glens-frontend"
  location = var.region
  build_config { runtime = "nodejs20"; entry_point = "handler" }
  service_config {
    min_instance_count = 0; max_instance_count = 3
    environment_variables = { API_URL = google_cloud_run_v2_service.api.uri }
  }
}

# Cloud Storage — static assets
resource "google_storage_bucket" "assets" {
  name     = "${var.project}-static"
  location = var.region
  website  { main_page_suffix = "index.html" }
}
```

## Network Topology

```
User ──► Cloud CDN/LB ──┬── Cloud Functions (Frontend SSR)
                         └── Cloud Storage (Static)
                              │
                         Cloud Run (Backend API)
                              │
                         Firestore (Data)
```

## Steps

1. Create `infra/` with provider + remote state config
2. Write Cloud Run, Cloud Functions, Storage, Firestore modules
3. Write Secret Manager module
4. Create `dev.tfvars` / `prod.tfvars`
5. Validate with `terraform plan`

## Success Criteria

- [ ] `terraform plan` succeeds with no errors
- [ ] `terraform apply` creates all resources
- [ ] Services accessible at generated URLs
- [ ] Secrets in Secret Manager, not env vars
- [ ] Total cost within GCP free tier
