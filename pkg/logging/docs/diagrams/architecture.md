# pkg/logging — Architecture

> Detailed diagrams for the shared zerolog wrapper.
> Master diagrams: [docs/diagrams/architecture.md](../../../docs/diagrams/architecture.md)

## Logging Usage Across Modules

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#15803d', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
graph TB
    subgraph PKG ["pkg/logging"]
        SETUP["Setup(Config)\ninitialise zerolog"]
        LEVEL["Level\ninfo · debug · warn · error"]
        FORMAT["Format\nconsole · json"]
    end

    GLENS["cmd/glens\nmain CLI"]
    API["cmd/api\nREST API server"]

    GLENS -->|imports| SETUP
    API -->|imports| SETUP

    DEMO["cmd/tools/demo\nzero deps — no import"]
    ACC["cmd/tools/accuracy\nzero deps — no import"]

    style PKG fill:#f0fdf4,stroke:#86efac,color:#0f172a
    style SETUP fill:#15803d,stroke:#166534,color:#fff
    style LEVEL fill:#047857,stroke:#065f46,color:#fff
    style FORMAT fill:#047857,stroke:#065f46,color:#fff
    style GLENS fill:#1e40af,stroke:#1e3a8a,color:#fff
    style API fill:#7c3aed,stroke:#6d28d9,color:#fff
    style DEMO fill:#64748b,stroke:#475569,color:#fff
    style ACC fill:#64748b,stroke:#475569,color:#fff
```

## Setup Sequence

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'actorTextColor': '#0f172a', 'actorBkg': '#e2e8f0', 'actorBorder': '#475569', 'signalColor': '#1e40af', 'signalTextColor': '#0f172a', 'noteBkgColor': '#fef3c7', 'noteTextColor': '#0f172a', 'noteBorderColor': '#d97706', 'fontSize': '14px'}}}%%
sequenceDiagram
    participant App as Application
    participant Log as pkg/logging
    participant ZL as zerolog

    App ->> Log: Setup(Config{Level, Format})
    Log ->> ZL: set global log level
    Log ->> ZL: set output writer (console or JSON)
    Log -->> App: logging ready

    App ->> ZL: log.Info().Msg("started")
    Note over ZL: output to stdout
```
