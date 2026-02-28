variable "region" {
  description = "GCP region for resource deployment"
  type        = string
  default     = "us-central1"
}

variable "alert_email" {
  description = "Email address for monitoring alert notifications (prod only)"
  type        = string
  default     = ""
}

variable "api_image_tag" {
  description = "Container image tag for the API service"
  type        = string
  default     = "latest"
}
