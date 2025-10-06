# Glens Architecture

## Problem It Solves

```mermaid
graph TB
    subgraph "Without Glens"
        P1[OpenAPI Spec Created] --> P2[Manual Test Writing]
        P2 --> P3[Time Consuming]
        P2 --> P4[Inconsistent Coverage]
        P2 --> P5[Human Error]
        P3 --> P6[Delayed QA]
        P4 --> P6
        P5 --> P6
    end

    subgraph "With Glens"
        S1[OpenAPI Spec] --> S2[AI Test Generation]
        S2 --> S3[Automated Execution]
        S3 --> S4{Tests Pass?}
        S4 -->|Yes| S5[Report Only]
        S4 -->|No| S6[GitHub Issue Created]
        S6 --> S7[Detailed Failure Info]
        S5 --> S8[Fast Feedback]
        S7 --> S8
    end

    style P6 fill:#ff6b6b
    style S8 fill:#51cf66
```

## High-Level Flow

```mermaid
flowchart LR
    A[OpenAPI Spec] --> B[Parse Spec]
    B --> C[Extract Endpoints]
    C --> D[AI Test Generation]
    D --> E[Execute Tests]
    E --> F{Results}
    F -->|Pass| G[Report]
    F -->|Fail| H[Create GitHub Issue]
    H --> G
    G --> I[Done]

    style A fill:#4dabf7
    style D fill:#fab005
    style H fill:#ff6b6b
    style I fill:#51cf66
```

## Detailed Component Architecture

```mermaid
graph TB
    subgraph "Input Layer"
        URL[OpenAPI Spec URL/File]
        CFG[Config File]
        ENV[Environment Variables]
    end

    subgraph "Parser Layer"
        PARSE[OpenAPI Parser]
        VAL[Spec Validator]
    end

    subgraph "AI Layer"
        MGR[AI Manager]
        GPT[OpenAI GPT-4]
        CLAUDE[Anthropic Claude]
        GEMINI[Google Gemini]
        OLLAMA[Ollama Local LLM]
    end

    subgraph "Generator Layer"
        GEN[Test Generator]
        EXEC[Test Executor]
        ANAL[Result Analyzer]
    end

    subgraph "Integration Layer"
        GH[GitHub Client]
        ISSUE[Issue Creator]
    end

    subgraph "Output Layer"
        MD[Markdown Report]
        HTML[HTML Report]
        JSON[JSON Report]
    end

    URL --> PARSE
    CFG --> MGR
    ENV --> MGR
    PARSE --> VAL
    VAL --> MGR

    MGR --> GPT
    MGR --> CLAUDE
    MGR --> GEMINI
    MGR --> OLLAMA

    GPT --> GEN
    CLAUDE --> GEN
    GEMINI --> GEN
    OLLAMA --> GEN

    GEN --> EXEC
    EXEC --> ANAL

    ANAL --> |Failures| ISSUE
    ISSUE --> GH

    ANAL --> MD
    ANAL --> HTML
    ANAL --> JSON

    style PARSE fill:#4dabf7
    style MGR fill:#fab005
    style GEN fill:#fab005
    style ANAL fill:#ff6b6b
    style ISSUE fill:#ff6b6b
    style MD fill:#51cf66
```

## Test Generation Flow

```mermaid
sequenceDiagram
    participant User
    participant CLI
    participant Parser
    participant AI
    participant Generator
    participant Executor
    participant GitHub
    participant Reporter

    User->>CLI: analyze command
    CLI->>Parser: Load OpenAPI spec
    Parser->>Parser: Validate spec
    Parser->>CLI: Return endpoints

    loop For each endpoint
        CLI->>AI: Request test generation
        AI->>AI: Analyze endpoint schema
        AI-->>CLI: Return test code

        CLI->>Generator: Save test file
        Generator->>Executor: Run go test
        Executor-->>Generator: Test results

        alt Test Failure
            Generator->>GitHub: Create issue
            GitHub-->>Generator: Issue created
        else Test Pass
            Generator->>Generator: Log success
        end

        Generator->>Reporter: Add results
    end

    Reporter->>User: Generate report
```

## Issue Creation Decision Logic

```mermaid
flowchart TD
    A[Test Executed] --> B{Compilation OK?}
    B -->|No| C[Infrastructure Issue]
    B -->|Yes| D{Tests Run?}
    D -->|No| E[Setup Issue]
    D -->|Yes| F{Any Failures?}
    F -->|No| G[All Pass - No Issue]
    F -->|Yes| H{Real Test Failure?}
    H -->|No| I[Connection/Setup Error]
    H -->|Yes| J[Create GitHub Issue]

    C --> K[Log Error - No Issue]
    E --> K
    I --> K
    G --> L[Report Success]
    J --> M[Report with Issue Link]

    style C fill:#ffd43b
    style E fill:#ffd43b
    style I fill:#ffd43b
    style K fill:#868e96
    style G fill:#51cf66
    style J fill:#ff6b6b
    style L fill:#51cf66
    style M fill:#ff8787
```

## AI Model Comparison Flow

```mermaid
graph LR
    subgraph "Input"
        EP[Endpoint Schema]
    end

    subgraph "AI Models"
        GPT[GPT-4<br/>Cloud]
        CLAUDE[Claude<br/>Cloud]
        OLLAMA[Ollama<br/>Local]
    end

    subgraph "Generated Tests"
        T1[Test 1]
        T2[Test 2]
        T3[Test 3]
    end

    subgraph "Execution"
        E1[Run Test 1]
        E2[Run Test 2]
        E3[Run Test 3]
    end

    subgraph "Analysis"
        CMP[Compare Results]
        BEST[Identify Best Model]
    end

    EP --> GPT
    EP --> CLAUDE
    EP --> OLLAMA

    GPT --> T1
    CLAUDE --> T2
    OLLAMA --> T3

    T1 --> E1
    T2 --> E2
    T3 --> E3

    E1 --> CMP
    E2 --> CMP
    E3 --> CMP

    CMP --> BEST

    style GPT fill:#fab005
    style CLAUDE fill:#fab005
    style OLLAMA fill:#fab005
    style BEST fill:#51cf66
```

## File Structure

```mermaid
graph TB
    subgraph "Project Root"
        ROOT[glens/]
    end

    subgraph "Commands"
        CMD[cmd/]
        ROOT_GO[root.go - Config & CLI]
        ANALYZE[analyze.go - Main Logic]
        MODELS[models.go - AI Management]
    end

    subgraph "Packages"
        PKG[pkg/]
        AI[ai/ - AI Clients]
        GEN[generator/ - Test Gen]
        GH[github/ - GitHub API]
        PARSE[parser/ - OpenAPI]
        REP[reporter/ - Reports]
    end

    subgraph "Config"
        CONF[configs/]
        EXAMPLE[config.example.yaml]
        USER[config.yaml]
    end

    subgraph "Templates"
        TMPL[templates/]
        ISSUE_T[issue-templates/]
        TEST_T[test-templates/]
    end

    ROOT --> CMD
    ROOT --> PKG
    ROOT --> CONF
    ROOT --> TMPL

    CMD --> ROOT_GO
    CMD --> ANALYZE
    CMD --> MODELS

    PKG --> AI
    PKG --> GEN
    PKG --> GH
    PKG --> PARSE
    PKG --> REP

    CONF --> EXAMPLE
    CONF --> USER

    TMPL --> ISSUE_T
    TMPL --> TEST_T

    style ROOT fill:#4dabf7
    style CMD fill:#fab005
    style PKG fill:#51cf66
```
