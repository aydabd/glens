resource "google_storage_bucket" "reports" {
  name     = "${var.project}-glens-reports"
  location = var.region
  project  = var.project

  uniform_bucket_level_access = true

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      age = 90
    }
    action {
      type = "Delete"
    }
  }
}
