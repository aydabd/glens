# Makefile for module glens/tools/accuracy
# Works standalone (can be moved to its own repo) or inside the monorepo.
# Targets: fmt, fmt-check, vet, tidy, lint, test, build, all
# Micromamba is used when available (local dev); plain go is used as fallback (CI).

MODULE      := glens/tools/accuracy
ENV_NAME    := glens-dev
LINT_VER    := v2.4.0

MAMBA := $(shell command -v micromamba 2>/dev/null)
ifdef MAMBA
  GO  := micromamba run -n $(ENV_NAME) go
  RUN := micromamba run -n $(ENV_NAME) bash -c
else
  GO  := go
  RUN := bash -c
endif

.DEFAULT_GOAL := help

.PHONY: all fmt fmt-check vet tidy lint test build clean help

all: fmt-check vet lint test ## Run all checks (CI equivalent)

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format Go source
	$(GO) fmt ./...

fmt-check: ## Check formatting (fails if unformatted; same check as CI)
	@if [ -n "$$(find . -name '*.go' | xargs gofmt -l)" ]; then \
		echo "Unformatted files (run: make fmt):"; \
		find . -name '*.go' | xargs gofmt -l; \
		exit 1; \
	fi

vet: ## Run go vet
	$(GO) vet ./...

tidy: ## Run go mod tidy
	$(GO) mod tidy

lint: ## Run golangci-lint (auto-installs if missing)
	@$(RUN) 'if ! command -v golangci-lint >/dev/null 2>&1; then \
		go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(LINT_VER); \
	fi && export PATH="$$(go env GOPATH)/bin:$$PATH" && golangci-lint run --timeout=3m'

test: ## Run tests with race detector
	$(GO) test -short -v -race ./...

build: ## Build the accuracy binary
	$(GO) build -o accuracy .

clean: ## Remove build artifacts
	$(GO) clean ./...
	rm -f accuracy
