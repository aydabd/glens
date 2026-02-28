locals {
  topics = [
    "glens-analyze",
    "glens-test-results",
    "glens-reports",
    "glens-secrets",
    "glens-export",
  ]
}

resource "google_pubsub_topic" "topics" {
  for_each = toset(local.topics)
  name     = each.value
  project  = var.project
}
