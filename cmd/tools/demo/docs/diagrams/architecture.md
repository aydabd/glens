# cmd/tools/demo — Architecture

> Detailed diagrams for the OpenAPI spec visualiser tool.
> Master diagrams: [docs/diagrams/architecture.md](../../../../docs/diagrams/architecture.md)

## Demo Tool — Flow

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#b45309', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
flowchart TD
    START["glens-demo &lt;spec&gt;"] --> FLAGS["Parse flags\n--spec · --version"]
    FLAGS --> LOAD["loader.Load\nfile or HTTP URL"]
    LOAD --> PARSE["loader.Parse\nJSON → OpenAPI struct"]
    PARSE --> RENDER["render output"]

    RENDER --> R1["Banner\ntool name + version"]
    RENDER --> R2["SpecInfo\ntitle · version · servers"]
    RENDER --> R3["Endpoints\nmethod · path · operationId"]
    RENDER --> R4["ModelComparison\nschema field table"]
    RENDER --> R5["SampleTest\nexample test snippet"]

    R1 --> OUT["stdout"]
    R2 --> OUT
    R3 --> OUT
    R4 --> OUT
    R5 --> OUT

    style START fill:#b45309,stroke:#92400e,color:#fff
    style FLAGS fill:#475569,stroke:#334155,color:#fff
    style LOAD fill:#0369a1,stroke:#075985,color:#fff
    style PARSE fill:#0369a1,stroke:#075985,color:#fff
    style RENDER fill:#6d28d9,stroke:#5b21b6,color:#fff
    style R1 fill:#15803d,stroke:#166534,color:#fff
    style R2 fill:#15803d,stroke:#166534,color:#fff
    style R3 fill:#15803d,stroke:#166534,color:#fff
    style R4 fill:#15803d,stroke:#166534,color:#fff
    style R5 fill:#15803d,stroke:#166534,color:#fff
    style OUT fill:#0f172a,stroke:#1e293b,color:#fff
```

## Internal Package Layout

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#b45309', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
graph TB
    MAIN["main.go\nflag parsing + orchestration"]

    subgraph INT ["internal/"]
        LOADER["loader/\nfetch + parse OpenAPI spec"]
        RENDER["render/\nbanner · endpoints · models"]
    end

    MAIN --> LOADER
    MAIN --> RENDER

    style INT fill:#f1f5f9,stroke:#94a3b8,color:#0f172a
    style MAIN fill:#b45309,stroke:#92400e,color:#fff
    style LOADER fill:#0369a1,stroke:#075985,color:#fff
    style RENDER fill:#15803d,stroke:#166534,color:#fff
```
