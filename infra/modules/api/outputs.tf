output "url" {
  description = "Cloud Run API service URL"
  value       = google_cloud_run_v2_service.api.uri
}
