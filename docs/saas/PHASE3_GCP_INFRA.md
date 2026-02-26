# Phase 3 — GCP Infrastructure (IaC)

> All GCP resources via Terraform. Free tier + database + analytics.

## Requirements Covered

IN-01 – IN-05, DB-01 (persist results), DB-03 (BigQuery export).

## GCP Free Tier + Services

| Service | Free allowance | Usage |
|---------|---------------|-------|
| Cloud Run | 2M req/mo, 360K vCPU-sec | Backend API |
| Cloud Functions | 2M invocations/mo | Frontend SSR |
| Cloud Storage | 5 GB | Static assets, reports |
| Firestore | 1 GB, 50K reads/day | User data + results |
| Secret Manager | 6 active versions | API keys + creds |
| Firebase Auth | 50K MAU | Authentication |
| Cloud Trace | 2.5M spans/mo | OTel traces |
| BigQuery | 1 TB query/mo, 10 GB | Analytics export |

Firestore over AlloyDB/Cloud SQL: serverless, scales to zero, no
idle cost, fits our document model. BigQuery for batch analytics.

## Terraform Layout

```
infra/
├── main.tf / variables.tf / outputs.tf
├── modules/
│   ├── api/          # Cloud Run
│   ├── frontend/     # Cloud Functions
│   ├── storage/      # Cloud Storage
│   ├── firestore/    # DB + indexes
│   ├── bigquery/     # Analytics
│   ├── secrets/      # Secret Manager
│   └── observability/ # Trace + Monitoring
└── environments/ { dev.tfvars, prod.tfvars }
```

## Test Results Schema (DB-01, Firestore)

```
workspaces/{wsId}/runs/{runId}
  ├── spec_url, models[], status, created_at
  ├── endpoints[]
  │     ├── path, method, category (read/write/destroy)
  │     └── results{ model → { passed, duration, output } }
  └── summary { total, passed, failed, duration }
```

Index: `(workspace_id, created_at DESC)` for dashboard queries.

## BigQuery Export (DB-03)

Daily scheduled Cloud Function exports Firestore → BigQuery dataset.

## Core Terraform (abbreviated)

```hcl
resource "google_cloud_run_v2_service" "api" {
  name = "glens-api"; location = var.region
  template {
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/glens/api:latest"
      ports { container_port = 8080 }
    }
    scaling { min_instance_count = 0; max_instance_count = 3 }
  }
}
```

## Network Topology

```
User → CDN/LB ─┬─ Cloud Functions (Frontend)
                └─ Cloud Storage (Static)
                     │
                Cloud Run (API) → Firestore + Secret Mgr
                                       │
                                  BigQuery (daily)
```

## Steps

1. Create `infra/` with provider + remote state (GCS bucket)
2. Write modules: api, frontend, storage, firestore, bigquery
3. Write secrets + observability modules
4. Create `dev.tfvars` / `prod.tfvars`; validate `terraform plan`

## Success Criteria

- [ ] `terraform plan` succeeds; `apply` creates all resources
- [ ] Firestore stores results with proper indexes (DB-01)
- [ ] BigQuery receives daily export (DB-03)
- [ ] Total cost within GCP free tier for < 1K users
