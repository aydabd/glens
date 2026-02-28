variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "region" {
  description = "GCP region for resource deployment"
  type        = string
  default     = "us-central1"
}

variable "log_level" {
  description = "Application log level (debug, info, warn, error)"
  type        = string
  default     = "info"
}

variable "min_instances" {
  description = "Minimum number of Cloud Run instances"
  type        = number
  default     = 0
}

variable "max_instances" {
  description = "Maximum number of Cloud Run instances"
  type        = number
  default     = 3
}

variable "alert_channel" {
  description = "Monitoring notification channel name"
  type        = string
  default     = "dev-slack"
}

variable "alert_email" {
  description = "Email address for monitoring alert notifications"
  type        = string
  default     = ""
}

variable "api_image_tag" {
  description = "Container image tag for the API service"
  type        = string
  default     = "latest"
}
