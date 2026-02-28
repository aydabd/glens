# Phase 19 — Real-Time Analysis & Live Monitoring

> Stream analysis progress and test results to users in real time via
> Server-Sent Events (SSE). Critical for user experience — users see
> each test pass/fail as it happens instead of waiting for the full run.

## Why

The current `/api/v1/analyze` endpoint returns `202 Accepted` with a run ID,
and clients must poll for results. This creates a poor UX for long-running
analyses (dozens of endpoints × multiple AI models). Real-time streaming:

- Shows each test result as it completes (seconds, not minutes of waiting).
- Enables live dashboards for team monitoring.
- Reduces server load (no polling loops).
- Allows the client to cancel a run mid-flight.
- Provides progress indicators (3/42 tests completed).

## Requirements

| ID | Requirement | Notes |
|----|-------------|-------|
| RT-01 | SSE endpoint streams analysis events | `GET /api/v1/runs/{id}/events` |
| RT-02 | Events include progress, result, error, done | Typed event stream |
| RT-03 | `responseWriter` forwards `http.Flusher` | Required for SSE |
| RT-04 | Client can cancel a run via API | `POST /api/v1/runs/{id}/cancel` |
| RT-05 | Run status persisted in Firestore | Survives reconnect |
| RT-06 | SSE reconnect via `Last-Event-ID` | RFC 9110 / EventSource spec |
| RT-07 | Live dashboard component (frontend) | SvelteKit EventSource |
| RT-08 | Test result events emitted to Pub/Sub | Downstream consumers |
| RT-09 | Heartbeat keepalive every 15 s | Prevent proxy/LB timeout |
| RT-10 | Rate limit concurrent runs per org | Prevent resource abuse |

## SSE Event Types

```text
event: progress
data: {"run_id":"abc","completed":3,"total":42,"percent":7}

event: test.pass
data: {"run_id":"abc","endpoint":"/pets","method":"GET","model":"gpt-4o","duration_ms":1240}

event: test.fail
data: {"run_id":"abc","endpoint":"/pets","method":"POST","model":"gpt-4o","error":"status 500","duration_ms":890}

event: error
data: {"run_id":"abc","message":"AI provider rate limit exceeded","retryable":true}

event: done
data: {"run_id":"abc","passed":38,"failed":4,"duration_ms":47200,"report_url":"/api/v1/reports/abc"}

: keepalive
```

## Architecture

```text
Client (browser / CLI)
  │
  ├── POST /api/v1/analyze          → 202 {run_id}
  │
  └── GET /api/v1/runs/{id}/events  → SSE stream
        │
        │  ┌──────────────────────────────────────┐
        │  │ SSE Handler                           │
        │  │  1. Validate run ownership            │
        │  │  2. Set Content-Type: text/event-stream│
        │  │  3. Flush headers                     │
        │  │  4. Subscribe to run channel          │
        │  │  5. Stream events until done/cancel   │
        │  └──────────────────────────────────────┘
        │
  Analysis Worker (goroutine or Cloud Function)
        │
        ├── For each endpoint × model:
        │     1. Generate test via AI
        │     2. Execute test
        │     3. Publish result event → Pub/Sub
        │     4. Push to run channel (in-process or Redis)
        │
        └── Emit done event → Pub/Sub + run channel
```

## `responseWriter` Flusher Support

The current `responseWriter` wrapper in `middleware.go` drops `http.Flusher`.
SSE requires flushing after each event. Implementation:

```go
// Flush forwards to the underlying writer if it supports http.Flusher.
func (rw *responseWriter) Flush() {
    if f, ok := rw.ResponseWriter.(http.Flusher); ok {
        f.Flush()
    }
}
```

This also enables future WebSocket upgrade support via `http.Hijacker`.

## API Endpoints

### `GET /api/v1/runs/{id}/events` — SSE stream

```text
Accept: text/event-stream
Last-Event-ID: 12          (optional, for reconnect)

200 OK
Content-Type: text/event-stream
Cache-Control: no-cache
Connection: keep-alive

id: 1
event: progress
data: {"run_id":"abc","completed":0,"total":42,"percent":0}

id: 2
event: test.pass
data: {"run_id":"abc","endpoint":"/pets","method":"GET","model":"gpt-4o","duration_ms":1240}
```

### `POST /api/v1/runs/{id}/cancel` — Cancel a run

```json
// Response
{
  "run_id": "abc",
  "status": "cancelled",
  "completed": 12,
  "total": 42
}
```

### `GET /api/v1/runs/{id}` — Run status (polling fallback)

```json
{
  "run_id": "abc",
  "status": "running",
  "completed": 12,
  "total": 42,
  "started_at": "2026-02-28T10:00:00Z",
  "results": []
}
```

## Data Model (Firestore)

```text
workspaces/{wsId}/runs/{runId}
  ├── status: "accepted" | "running" | "completed" | "failed" | "cancelled"
  ├── spec_url: string
  ├── models: string[]
  ├── total_tests: int
  ├── completed_tests: int
  ├── passed: int
  ├── failed: int
  ├── started_at: timestamp
  ├── completed_at: timestamp | null
  └── results/{idx}
       ├── endpoint_path: string
       ├── endpoint_method: string
       ├── model: string
       ├── status: "pass" | "fail" | "error"
       ├── error_message: string | null
       ├── duration_ms: int
       └── created_at: timestamp
```

## Event Flow (Pub/Sub integration)

```text
Analysis Worker
  │
  ├── test.started → Pub/Sub topic: glens-test-events
  ├── test.passed  → Pub/Sub topic: glens-test-events
  ├── test.failed  → Pub/Sub topic: glens-test-events
  └── run.completed → Pub/Sub topic: glens-run-events
        │
        ├── Cloud Function: generate report
        ├── Cloud Function: create GitHub issues (if enabled)
        └── Cloud Function: update dashboard metrics
```

## Frontend Integration

```svelte
<!-- RunMonitor.svelte -->
<script>
  const source = new EventSource(`/api/v1/runs/${runId}/events`);

  source.addEventListener('progress', (e) => {
    progress = JSON.parse(e.data);
  });

  source.addEventListener('test.pass', (e) => {
    results = [...results, { ...JSON.parse(e.data), status: 'pass' }];
  });

  source.addEventListener('test.fail', (e) => {
    results = [...results, { ...JSON.parse(e.data), status: 'fail' }];
  });

  source.addEventListener('done', (e) => {
    source.close();
    summary = JSON.parse(e.data);
  });
</script>

<ProgressBar value={progress.completed} max={progress.total} />
{#each results as r}
  <TestResult {r} />
{/each}
```

## Implementation Phases (sub-phases)

| Step | Description | Depends on |
|------|-------------|------------|
| 19a | `responseWriter` forwards `http.Flusher` | Phase 1 |
| 19b | Run model in Firestore + CRUD | Phase 3 |
| 19c | SSE handler + in-process event channel | 19a |
| 19d | Analysis worker emits events per test | 19b, Phase 11 |
| 19e | `Last-Event-ID` reconnect support | 19c |
| 19f | Run cancel endpoint | 19b |
| 19g | Pub/Sub integration for test events | 19d, Phase 11 |
| 19h | Frontend live dashboard component | 19c, Phase 2 |
| 19i | Concurrent run rate limiting | 19b, Phase 18 |

## Implementation Steps

1. **19a** — Add `Flush()` to `responseWriter` in `middleware.go`
2. **19b** — Add `Run` model + status tracking in Firestore
3. **19c** — Implement SSE handler at `GET /api/v1/runs/{id}/events`
4. **19d** — Wire analysis worker to emit per-test events
5. **19e** — Support `Last-Event-ID` header for reconnect
6. **19f** — Add `POST /api/v1/runs/{id}/cancel` endpoint
7. **19g** — Publish test events to Pub/Sub for downstream consumers
8. **19h** — Build SvelteKit `RunMonitor` component with `EventSource`
9. **19i** — Add per-org concurrent run limits (ties into Phase 18 rate limiting)

## Success Criteria

- [ ] SSE endpoint streams test results in real time (19c)
- [ ] `responseWriter` supports `http.Flusher` (19a)
- [ ] Run status persisted; survives reconnect (19b, 19e)
- [ ] Analysis worker emits events per test (19d)
- [ ] Client can cancel a running analysis (19f)
- [ ] Frontend shows live progress + results (19h)
- [ ] Test events published to Pub/Sub (19g)
- [ ] Heartbeat prevents proxy timeouts (19c)
- [ ] Concurrent run limits enforced (19i)
