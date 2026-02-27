# Phase 11 — Event-Driven Architecture

> Domain events via Pub/Sub triggering Cloud Functions for async work.

## Requirements Covered

EV-01 (domain events), EV-02 (Pub/Sub topics), EV-03 (event-triggered
Cloud Functions), EV-04 (event schema contracts).

## Why Events

Long-running or decoupled workflows: report generation, BigQuery
export, GitHub issue creation, notifications. API responds immediately;
Cloud Functions handle side effects asynchronously.

## Event Catalog

| Event | Topic | Subscriber |
|-------|-------|------------|
| `analyze.completed` | `glens-analyze` | Report generator |
| `test.failed` | `glens-test-results` | Issue creator |
| `report.generated` | `glens-reports` | Notification fn |
| `secret.stored` | `glens-secrets` | Audit logger |
| `export.scheduled` | `glens-export` | BigQuery exporter |

## Event Schema (JSON)

```json
{ "event_type": "analyze.completed",
  "event_id": "uuid-v4",
  "timestamp": "2026-02-26T12:00:00Z",
  "workspace_id": "ws_abc",
  "payload": { "run_id": "run_123", "passed": 8, "failed": 2 } }
```

All events share `event_type`, `event_id`, `timestamp`, `workspace_id`.
Payload varies by type. Schemas published in `openapi.yaml`.

## Architecture

```text
API (Cloud Run) ──publish──► Pub/Sub Topic
                                │
  ┌─────────────────────────────┼──────────────────┐
  ▼                             ▼                  ▼
fn-report-generator    fn-issue-creator    fn-notification
```

Each function is independently deployable and testable.

## Publisher (Go)

```go
func publishEvent(ctx context.Context, topic string, evt Event) error {
    data, _ := json.Marshal(evt)
    result := pubsubClient.Topic(topic).Publish(ctx,
        &pubsub.Message{Data: data,
            Attributes: map[string]string{
                "event_type": evt.Type, "workspace_id": evt.WorkspaceID}})
    _, err := result.Get(ctx)
    return err
}
```

## Terraform

```hcl
resource "google_pubsub_topic" "analyze" { name = "glens-analyze" }
resource "google_cloudfunctions2_function" "report_gen" {
  name = "fn-report-generator"
  event_trigger {
    event_type   = "google.cloud.pubsub.topic.v1.messagePublished"
    pubsub_topic = google_pubsub_topic.analyze.id
  }
}
```

## Testing

- **Unit**: mock Pub/Sub client; verify message shape
- **Integration**: Pub/Sub emulator; publish → function triggers
- **E2E**: analyze spec → verify report + issue created async

## Steps

1. Define event schema in `cmd/api/internal/events/schema.go`
2. Add Pub/Sub publisher to API handlers
3. Create Cloud Functions for each subscriber
4. Add topics + functions to Terraform
5. Write unit + integration tests with emulator

## Success Criteria

- [ ] `analyze.completed` triggers report generation
- [ ] `test.failed` creates GitHub issue asynchronously
- [ ] Events match published schema contract
- [ ] Each function independently testable with emulator
