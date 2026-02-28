output "notification_channel_id" {
  description = "Monitoring notification channel ID"
  value       = var.alert_email != "" ? google_monitoring_notification_channel.alert[0].id : ""
}
