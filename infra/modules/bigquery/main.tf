resource "google_bigquery_dataset" "analytics" {
  dataset_id = "glens_analytics"
  project    = var.project
  location   = var.region

  labels = {
    managed_by = "terraform"
  }
}

resource "google_bigquery_table" "run_results" {
  dataset_id = google_bigquery_dataset.analytics.dataset_id
  table_id   = "run_results"
  project    = var.project

  schema = jsonencode([
    { name = "run_id", type = "STRING", mode = "REQUIRED" },
    { name = "workspace_id", type = "STRING", mode = "REQUIRED" },
    { name = "status", type = "STRING", mode = "REQUIRED" },
    { name = "duration_ms", type = "INTEGER", mode = "NULLABLE" },
    { name = "created_at", type = "TIMESTAMP", mode = "REQUIRED" },
  ])
}
