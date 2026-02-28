variable "project" {
  description = "GCP project ID"
  type        = string
}

variable "alert_channel" {
  description = "Monitoring notification channel name"
  type        = string
}

variable "alert_email" {
  description = "Email address for monitoring alert notifications"
  type        = string
  default     = ""
}
