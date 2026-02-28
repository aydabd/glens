# cmd/glens — Architecture

> Detailed diagrams for the main Glens CLI module.
> Master diagrams: [docs/diagrams/architecture.md](../../../docs/diagrams/architecture.md)

## Analyze Command — Pipeline Flow

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#1e40af', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
flowchart TD
    START["glens analyze &lt;spec&gt;"] --> INIT["Load config\nSetup logging"]
    INIT --> PARSE["parser.ParseOpenAPISpec\nfile or URL"]
    PARSE --> ENDPOINTS["Extract endpoints\nmethod · path · params · schemas"]
    ENDPOINTS --> LOOP["For each endpoint"]

    LOOP --> AIGEN["AI: generate test code"]
    AIGEN --> EXEC["generator: compile + run test"]
    EXEC --> CHECK{Passed?}

    CHECK -->|Yes| NEXT["Next endpoint"]
    CHECK -->|No| REAL{Real failure?}
    REAL -->|No| LOG["Log warning\nskip issue"]
    REAL -->|Yes| GHISSUE["github: create issue"]
    LOG --> NEXT
    GHISSUE --> NEXT
    NEXT --> LOOP

    LOOP -->|All done| REPORT["reporter: generate\nMarkdown + HTML"]
    REPORT --> DONE["✅ Output files"]

    style START fill:#1e40af,stroke:#1e3a8a,color:#fff
    style INIT fill:#475569,stroke:#334155,color:#fff
    style PARSE fill:#0369a1,stroke:#075985,color:#fff
    style ENDPOINTS fill:#0369a1,stroke:#075985,color:#fff
    style LOOP fill:#6d28d9,stroke:#5b21b6,color:#fff
    style AIGEN fill:#b45309,stroke:#92400e,color:#fff
    style EXEC fill:#6d28d9,stroke:#5b21b6,color:#fff
    style CHECK fill:#475569,stroke:#334155,color:#fff
    style REAL fill:#475569,stroke:#334155,color:#fff
    style LOG fill:#d97706,stroke:#b45309,color:#fff
    style GHISSUE fill:#b91c1c,stroke:#991b1b,color:#fff
    style NEXT fill:#64748b,stroke:#475569,color:#fff
    style REPORT fill:#047857,stroke:#065f46,color:#fff
    style DONE fill:#15803d,stroke:#166534,color:#fff
```

## Internal Package Layout

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'primaryColor': '#1e40af', 'primaryTextColor': '#fff', 'lineColor': '#475569', 'fontSize': '14px'}}}%%
graph TB
    subgraph CMD ["cmd/ — Cobra commands"]
        ROOT["root.go\nconfig · logging"]
        ANALYZE["analyze.go\nmain pipeline"]
        CLEANUP["cleanup.go\nissue cleanup"]
        MODELS["models.go\nAI model mgmt"]
    end

    subgraph INT ["internal/ — private packages"]
        PARSER["parser\nOpenAPI extraction"]
        AI["ai\nClient interface\nOpenAI · Claude · Gemini · Ollama"]
        GEN["generator\ntest exec in temp dir"]
        GH["github\nissue creation"]
        REP["reporter\nMarkdown · HTML · JSON"]
    end

    ANALYZE --> PARSER
    ANALYZE --> AI
    ANALYZE --> GEN
    ANALYZE --> GH
    ANALYZE --> REP
    CLEANUP --> GH
    MODELS --> AI

    style CMD fill:#eff6ff,stroke:#93c5fd,color:#0f172a
    style INT fill:#f1f5f9,stroke:#94a3b8,color:#0f172a
    style ROOT fill:#1e40af,stroke:#1e3a8a,color:#fff
    style ANALYZE fill:#1e40af,stroke:#1e3a8a,color:#fff
    style CLEANUP fill:#1e40af,stroke:#1e3a8a,color:#fff
    style MODELS fill:#1e40af,stroke:#1e3a8a,color:#fff
    style PARSER fill:#0369a1,stroke:#075985,color:#fff
    style AI fill:#b45309,stroke:#92400e,color:#fff
    style GEN fill:#6d28d9,stroke:#5b21b6,color:#fff
    style GH fill:#b91c1c,stroke:#991b1b,color:#fff
    style REP fill:#047857,stroke:#065f46,color:#fff
```

## AI Provider — Sequence Diagram

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'actorTextColor': '#0f172a', 'actorBkg': '#e2e8f0', 'actorBorder': '#475569', 'signalColor': '#1e40af', 'signalTextColor': '#0f172a', 'noteBkgColor': '#fef3c7', 'noteTextColor': '#0f172a', 'noteBorderColor': '#d97706', 'fontSize': '14px'}}}%%
sequenceDiagram
    participant Analyze as analyze.go
    participant Mgr as AI Manager
    participant Client as AI Client
    participant API as Provider API

    Analyze ->> Mgr: GenerateTest(model, endpoint)
    Mgr ->> Client: select client for model
    Client ->> API: HTTP POST (prompt + schema)
    API -->> Client: generated Go test code
    Client -->> Mgr: TestGenerationResult
    Mgr -->> Analyze: test code + metadata
```

## Test Execution — Sequence Diagram

```mermaid
%%{init: {'theme': 'base', 'themeVariables': {'actorTextColor': '#0f172a', 'actorBkg': '#e2e8f0', 'actorBorder': '#475569', 'signalColor': '#1e40af', 'signalTextColor': '#0f172a', 'noteBkgColor': '#fef3c7', 'noteTextColor': '#0f172a', 'noteBorderColor': '#d97706', 'fontSize': '14px'}}}%%
sequenceDiagram
    participant Analyze as analyze.go
    participant Gen as generator
    participant FS as temp directory
    participant Go as go test

    Analyze ->> Gen: ExecuteTest(code, endpoint)
    Gen ->> FS: create temp dir
    Gen ->> FS: write test file + go.mod
    Gen ->> Go: run go test ./...
    Go -->> Gen: stdout + stderr + exit code
    Gen ->> Gen: parse results
    Gen -->> Analyze: ExecutionResult
    Note over Gen: compiled? passed? failures?
```
