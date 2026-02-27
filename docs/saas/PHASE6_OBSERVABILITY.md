# Phase 6 — Observability (OpenTelemetry)

> End-to-end traces, metrics, and structured logs via OpenTelemetry.

## Requirements Covered

OB-01 (OTel SDK), OB-02 (Cloud Trace + Monitoring), OB-03 (logging).

## Why OpenTelemetry

- **Vendor-neutral** — single SDK exports to any backend.
- **GCP-native** — Google maintains the Cloud Trace exporter.
- **Go SDK 1.34** — stable traces + metrics API (Feb 2026).
- **Auto-instrumentation** — `net/http`, gRPC, SQL out of the box.

## Architecture

```
Backend (Go) → OTel SDK → TracerProvider + MeterProvider
  ├── Cloud Trace   (free 2.5M spans/mo)
  ├── Cloud Monitoring (free custom metrics)
  └── Cloud Logging  (free 50 GB/mo)
```

Frontend sends no telemetry; backend correlates via `traceparent`.

## Go Integration

```go
import (
  "fmt"

  gcptrace "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/trace"
  gcpmetric "github.com/GoogleCloudPlatform/opentelemetry-operations-go/exporter/metric"
)

func initOTel(ctx context.Context, project string) (func(), error) {
  te, err := gcptrace.New(gcptrace.WithProjectID(project))
  if err != nil {
    return nil, fmt.Errorf("create trace exporter: %w", err)
  }
  tp := trace.NewTracerProvider(trace.WithBatcher(te))
  otel.SetTracerProvider(tp)

  me, err := gcpmetric.New(gcpmetric.WithProjectID(project))
  if err != nil {
    return nil, fmt.Errorf("create metric exporter: %w", err)
  }
  mp := metric.NewMeterProvider(metric.WithReader(
    metric.NewPeriodicReader(me)))
  otel.SetMeterProvider(mp)
  return func() { tp.Shutdown(ctx); mp.Shutdown(ctx) }, nil
}
```

## Key Spans

| Span | Attributes |
|------|------------|
| `analyze.parse_spec` | `spec.url`, `endpoint_count` |
| `analyze.generate_test` | `ai.model`, `endpoint.path` |
| `analyze.execute_test` | `test.passed`, `duration_ms` |
| `auth_proxy.resolve` | `secret.ref` (no raw value!) |

## Custom Metrics

| Metric | Type | Labels |
|--------|------|--------|
| `glens.analyze.duration` | Histogram | `model`, `status` |
| `glens.tests.total` | Counter | `result` (pass/fail) |
| `glens.api.requests` | Counter | `method`, `path`, `code` |

## Structured Logging (OB-03)

`zerolog` (already in `pkg/logging`) with GCP-compatible JSON.
Cloud Logging auto-correlates logs with traces via `trace_id`:

```go
log.Info().Str("trace_id", span.SpanContext().TraceID().String()).
  Str("endpoint", path).Msg("test execution complete")
```

## Dependencies (latest stable, Feb 2026)

| Package | Version |
|---------|---------|
| `go.opentelemetry.io/otel` | 1.34.0 |
| `opentelemetry-operations-go/exporter/trace` | 1.25.0 |
| `opentelemetry-operations-go/exporter/metric` | 0.49.0 |

## Steps

1. Add OTel deps to `cmd/api/go.mod`
2. Implement `initOTel()` in `internal/telemetry/`
3. Instrument handlers with spans + custom metrics
4. Inject `trace_id` into zerolog; verify in Cloud Trace

## Success Criteria

- [ ] Every analyze call produces a trace in Cloud Trace
- [ ] Custom metrics visible in Cloud Monitoring
- [ ] Logs correlated with traces via `trace_id`
- [ ] No secret values in any span attributes
