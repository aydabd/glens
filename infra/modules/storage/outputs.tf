output "bucket_name" {
  description = "Cloud Storage bucket name"
  value       = google_storage_bucket.reports.name
}

# TODO: Replace with actual Cloud Functions or Firebase Hosting URL once frontend is deployed.
output "frontend_url" {
  description = "Cloud Functions frontend URL (placeholder — not yet deployed)"
  value       = "https://${var.project}.web.app"
}
