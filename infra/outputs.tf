output "api_url" {
  description = "Cloud Run API service URL"
  value       = module.api.url
}

output "frontend_url" {
  description = "Cloud Functions frontend URL"
  value       = module.storage.frontend_url
}
