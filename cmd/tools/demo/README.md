# glens-demo

Renders a human-readable summary of an OpenAPI spec: endpoint list, model comparison table, and a sample test snippet.

Replaces `scripts/demo_modern.sh`. Module: `glens/tools/demo`

## Install

### Download binary

Download `glens-demo` for your platform from the [releases page](https://github.com/aydabd/glens/releases).

### Build from source

```bash
cd cmd/tools/demo
make build        # builds to ../../build/glens-demo
```

## Usage

```bash
# Load from a local file
./build/glens-demo path/to/openapi.json

# Load from a URL
./build/glens-demo https://petstore3.swagger.io/api/v3/openapi.json
```

## Makefile targets

Run from this directory (`cmd/tools/demo/`):

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
cmd/tools/demo/
├── main.go                       # Entry point (flag parsing, exit codes)
├── internal/
│   ├── loader/
│   │   └── loader.go             # Fetch + parse OpenAPI spec (file or HTTP URL)
│   └── render/
│       └── render.go             # Banner, endpoint list, model table, snippet
├── go.mod                        # Module: glens/tools/demo (zero external deps)
├── Makefile
└── README.md
```

Zero external dependencies — pure Go stdlib.
