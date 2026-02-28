resource "google_firestore_database" "main" {
  project     = var.project
  name        = "(default)"
  location_id = var.region
  type        = "FIRESTORE_NATIVE"
}

resource "google_firestore_index" "workspace_runs" {
  project    = var.project
  database   = google_firestore_database.main.name
  collection = "runs"

  fields {
    field_path = "workspace_id"
    order      = "ASCENDING"
  }

  fields {
    field_path = "created_at"
    order      = "DESCENDING"
  }
}
