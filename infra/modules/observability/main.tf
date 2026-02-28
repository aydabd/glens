resource "google_project_service" "cloudtrace" {
  project = var.project
  service = "cloudtrace.googleapis.com"

  disable_on_destroy = false
}

resource "google_monitoring_notification_channel" "alert" {
  count        = var.alert_email != "" ? 1 : 0
  project      = var.project
  display_name = var.alert_channel
  type         = "email"

  labels = {
    email_address = var.alert_email
  }
}
