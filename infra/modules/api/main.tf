resource "google_cloud_run_v2_service" "api" {
  name     = "glens-api"
  location = var.region
  project  = var.project

  template {
    scaling {
      min_instance_count = var.min_instances
      max_instance_count = var.max_instances
    }

    containers {
      image = "${var.region}-docker.pkg.dev/${var.project}/glens/api:${var.api_image_tag}"

      env {
        name  = "LOG_LEVEL"
        value = var.log_level
      }

      ports {
        container_port = 8080
      }

      resources {
        limits = {
          cpu    = "1"
          memory = "512Mi"
        }
      }
    }
  }
}

resource "google_cloud_run_v2_service_iam_member" "public" {
  name     = google_cloud_run_v2_service.api.name
  location = google_cloud_run_v2_service.api.location
  project  = var.project
  role     = "roles/run.invoker"
  member   = "allUsers"
}
