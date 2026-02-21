# glens/pkg/logging

Generic zerolog setup wrapper. Configures global log level, format (JSON or console), and timestamps.

Module: `glens/pkg/logging`

This library has **no imports from any `internal/` package** and no dependencies on glens internals — it can be used in any Go project or moved to a separate repository at any time.

## Install

Inside the monorepo workspace, `go.work` resolves this automatically via a `replace` directive.

To use it in an external project:

```bash
go get glens/pkg/logging@vX.Y.Z
```

## Usage

```go
import "glens/pkg/logging"

func main() {
    logging.Setup(logging.Config{
        Level:  logging.LevelInfo,
        Format: logging.FormatJSON,
    })

    log.Info().Str("key", "value").Msg("started")
}
```

### Config fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `Level` | `Level` | `LevelInfo` | Log level (`debug`, `info`, `warn`, `error`) |
| `Format` | `Format` | `FormatJSON` | Output format (`json` or `console`) |

## Makefile targets

Run from this directory (`pkg/logging/`):

| Target | Description |
|--------|-------------|
| `make all` | fmt-check + vet + lint + test (same as CI) |
| `make fmt` | Format source |
| `make fmt-check` | Fail if source is unformatted |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint |
| `make test` | Run tests with race detector |
| `make clean` | Remove build artifacts |

## Versioning

Tag releases with the `pkg/logging/` prefix:

```bash
git tag pkg/logging/v0.1.0
git push origin pkg/logging/v0.1.0
```

## Module structure

```
pkg/logging/
├── logging.go       # Setup(), Level, Format types
├── logging_test.go
├── go.mod           # Module: glens/pkg/logging
├── Makefile
└── README.md
```
