resource "google_project_service" "secretmanager" {
  project = var.project
  service = "secretmanager.googleapis.com"

  disable_on_destroy = false
}
