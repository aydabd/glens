# cmd/tools/accuracy — Architecture

> Detailed diagrams for the endpoint accuracy reporter tool.
> Master diagrams: [docs/diagrams/architecture.md](../../../../docs/diagrams/architecture.md)

## Accuracy Tool — Flow

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#b45309', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
flowchart TD
    START["glens-accuracy &lt;specs&gt;"] --> FLAGS["Parse flags\n--output · --version"]
    FLAGS --> LOAD["analyze.Load\nfile or HTTP URL"]
    LOAD --> COUNT["analyze.Count\nendpoints per path/method"]
    COUNT --> BUILD["report.Build\nmarkdown table"]
    BUILD --> WRITE{--output?}
    WRITE -->|Yes| FILE["Write to file"]
    WRITE -->|No| STDOUT["Print to stdout"]

    style START fill:#b45309,stroke:#92400e,color:#fff
    style FLAGS fill:#475569,stroke:#334155,color:#fff
    style LOAD fill:#0369a1,stroke:#075985,color:#fff
    style COUNT fill:#6d28d9,stroke:#5b21b6,color:#fff
    style BUILD fill:#047857,stroke:#065f46,color:#fff
    style WRITE fill:#475569,stroke:#334155,color:#fff
    style FILE fill:#15803d,stroke:#166534,color:#fff
    style STDOUT fill:#15803d,stroke:#166534,color:#fff
```

## Internal Package Layout

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#b45309', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
graph TB
    MAIN["main.go\nflag parsing + orchestration"]

    subgraph INT ["internal/"]
        ANALYZE["analyze/\nload specs · count endpoints"]
        REPORT["report/\nbuild markdown report"]
    end

    MAIN --> ANALYZE
    MAIN --> REPORT

    style INT fill:#f1f5f9,stroke:#94a3b8,color:#0f172a
    style MAIN fill:#b45309,stroke:#92400e,color:#fff
    style ANALYZE fill:#0369a1,stroke:#075985,color:#fff
    style REPORT fill:#047857,stroke:#065f46,color:#fff
```
