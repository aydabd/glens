output "topic_ids" {
  description = "Map of Pub/Sub topic names to IDs"
  value       = { for k, v in google_pubsub_topic.topics : k => v.id }
}
