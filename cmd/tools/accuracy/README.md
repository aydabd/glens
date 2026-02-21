# glens-accuracy

Loads one or more OpenAPI specs, counts endpoints per path/method, and emits a markdown accuracy report.

Replaces `scripts/test_accuracy.sh`. Module: `glens/tools/accuracy`

## Install

### Download binary

Download `glens-accuracy` for your platform from the [releases page](https://github.com/aydabd/glens/releases).

### Build from source

```bash
cd cmd/tools/accuracy
make build        # builds to ../../build/glens-accuracy
```

## Usage

```bash
# Single spec — prints report to stdout
./build/glens-accuracy path/to/openapi.json

# Multiple specs
./build/glens-accuracy spec1.json spec2.json https://example.com/spec.json

# Write to file
./build/glens-accuracy --output report.md spec.json
```

## Makefile targets

Run from this directory (`cmd/tools/accuracy/`):

| Target | Description |
|--------|-------------|
| `make all` | fmt-check + vet + lint + test (same as CI) |
| `make fmt` | Format source |
| `make fmt-check` | Fail if source is unformatted |
| `make vet` | Run `go vet` |
| `make lint` | Run golangci-lint |
| `make test` | Run tests with race detector |
| `make build` | Build binary |
| `make clean` | Remove build artifacts |

## Module structure

```
cmd/tools/accuracy/
├── main.go                       # Entry point (flag parsing, exit codes)
├── internal/
│   ├── analyze/
│   │   └── analyze.go            # Load specs, count endpoints, cross-platform paths
│   └── report/
│       └── report.go             # Build markdown accuracy report
├── go.mod                        # Module: glens/tools/accuracy (zero external deps)
├── Makefile
└── README.md
```

Zero external dependencies — pure Go stdlib.
