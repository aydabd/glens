# Glens Makefile

# Core Configuration
BINARY_NAME := glens
GO_VERSION := 1.25
MAIN_PACKAGE := .
BUILD_DIR := build
DIST_DIR := $(BUILD_DIR)/dist
REPORTS_DIR := $(BUILD_DIR)/reports
PROFILES_DIR := $(BUILD_DIR)/profiles
DOCS_DIR := $(BUILD_DIR)/docs
ENV_NAME := glens-dev
OPENAPI_URL ?= https://petstore3.swagger.io/api/v3/openapi.json
GITHUB_REPO_ISSUE_CREATION_TEST ?= aydabd/test-agent-ideas
OP_ID ?= getPetById

# Build Configuration
LDFLAGS := -ldflags="-s -w"
BUILD_FLAGS := -v $(LDFLAGS)

# Platform Targets
PLATFORMS := linux/amd64 windows/amd64 darwin/amd64 darwin/arm64

# File Patterns
COVERAGE_OUT := $(REPORTS_DIR)/coverage.out
COVERAGE_HTML := $(REPORTS_DIR)/coverage.html
CPU_PROF := $(PROFILES_DIR)/cpu.prof
MEM_PROF := $(PROFILES_DIR)/mem.prof
BUILD_SUBDIRS := $(DIST_DIR) $(REPORTS_DIR) $(PROFILES_DIR) $(DOCS_DIR)

# Micromamba Environment
MAMBA_CMD := micromamba
MAMBA_RUN := $(MAMBA_CMD) run -n $(ENV_NAME)
GOCMD := $(MAMBA_RUN) go

# Tool Commands
GOBUILD := $(GOCMD) build
GOCLEAN := $(GOCMD) clean
GOTEST := $(GOCMD) test
GOGET := $(GOCMD) get
GOMOD := $(GOCMD) mod
GOFMT := $(GOCMD) fmt

# Development Tools
TOOLS := golangci-lint gosec govulncheck air
TOOL_GOLANGCI := github.com/golangci/golangci-lint/cmd/golangci-lint@latest
TOOL_GOSEC := github.com/securecodewarrior/gosec/v2/cmd/gosec@latest
TOOL_GOVULN := golang.org/x/vuln/cmd/govulncheck@latest
TOOL_AIR := github.com/cosmtrek/air@latest

# Messages
MSG_PREFIX := 🚀
SUCCESS := ✅
ERROR := ❌
WARNING := ⚠️
INFO := ℹ️

# Helper Functions
define check_tool
	@command -v $(1) >/dev/null 2>&1 || { \
		echo "$(ERROR) $(1) not found. $(2)"; \
		exit 1; \
	}
endef

define install_tool
	@if ! $(MAMBA_RUN) bash -c 'export PATH="$$(go env GOPATH)/bin:$$PATH" && which $(1)' >/dev/null 2>&1; then \
		echo "$(WARNING) Installing $(1)..."; \
		$(MAMBA_RUN) go install $(2); \
		echo "$(SUCCESS) $(1) installed"; \
	fi
endef

define run_tool
	@$(MAMBA_RUN) bash -c 'export PATH="$$(go env GOPATH)/bin:$$PATH" && $(1)'
endef

define build_platform
	@echo "$(MSG_PREFIX) Building for $(1)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(word 1,$(subst /, ,$(1))) GOARCH=$(word 2,$(subst /, ,$(1))) \
		$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(subst /,-,$(1))$(if $(findstring windows,$(1)),.exe) $(MAIN_PACKAGE)
endef

define ensure_dirs
	@mkdir -p $(BUILD_SUBDIRS)
endef

define env_exists
	$(MAMBA_CMD) info -e --json | grep -q "/$(ENV_NAME)\""
endef

define env_action
	if $(call env_exists); then \
		echo "$(MSG_PREFIX) $(1)..."; \
		$(2); \
	else \
		echo "$(MSG_PREFIX) $(3)..."; \
		$(4); \
	fi
endef

# export github token with gh cli - returns the token value
define get_github_token
$(shell gh auth token 2>/dev/null)
endef

# Default target
.PHONY: all
all: env clean deps test build

.DEFAULT_GOAL := help

# Help target
.PHONY: help
help: ## Show this help message
	@echo "$(BINARY_NAME) - Glens"
	@echo ""
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} \
		/^[a-zA-Z_-]+:.*##/ { printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
	@echo ""
	@echo "Variables:"
	@echo "  OPENAPI_URL - OpenAPI specification URL (default: $(OPENAPI_URL))"
	@echo "  OP_ID       - Operation ID for targeted testing"
	@echo ""
	@echo "Build Structure:"
	@echo "  $(BUILD_DIR)/           - Main build directory"
	@echo "  $(DIST_DIR)/       - Release distributions"
	@echo "  $(REPORTS_DIR)/    - Test reports, coverage, lint results"
	@echo "  $(PROFILES_DIR)/   - Performance profiles"
	@echo "  $(DOCS_DIR)/       - Generated documentation"

# ==============================================================================
# ENVIRONMENT MANAGEMENT
# ==============================================================================

.PHONY: check-mamba
check-mamba: ## Check if micromamba is installed
	$(call check_tool,micromamba,Please install micromamba: https://mamba.readthedocs.io/en/latest/installation.html)

.PHONY: check-env
check-env: check-mamba ## Check if micromamba environment exists
	@$(call env_exists) || { \
		echo "$(ERROR) Environment '$(ENV_NAME)' not found"; \
		exit 1; \
	}

.PHONY: env
env: check-mamba ## Create/update micromamba environment
	@$(call env_action,Updating environment,$(MAMBA_CMD) env update -n $(ENV_NAME) -f environment.yml,Creating environment,$(MAMBA_CMD) env create -f environment.yml)
	@echo "$(SUCCESS) Environment '$(ENV_NAME)' ready"

.PHONY: clean-env
clean-env: check-mamba ## Remove micromamba environment
	@$(call env_action,Removing environment,$(MAMBA_CMD) env remove -n $(ENV_NAME) -y && echo "$(SUCCESS) Environment removed",Environment does not exist,echo "$(INFO) Environment does not exist")

.PHONY: shell
shell: check-env ## Enter micromamba environment shell
	@$(MAMBA_CMD) run -n $(ENV_NAME) bash

# ==============================================================================
# BUILD & DEVELOPMENT
# ==============================================================================

.PHONY: deps
deps: check-env ## Download dependencies
	@$(GOMOD) download && $(GOMOD) tidy

.PHONY: build
build: check-env deps ## Build the binary
	@echo "$(MSG_PREFIX) Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)
	@echo "$(SUCCESS) Build completed: $(BUILD_DIR)/$(BINARY_NAME)"

.PHONY: build-all
build-all: check-env deps ## Build for all platforms
	@echo "$(MSG_PREFIX) Starting multi-platform build..."
	@$(foreach platform,$(PLATFORMS),$(MAKE) build-platform PLATFORM=$(platform);)
	@echo "$(SUCCESS) All platform builds completed"

.PHONY: build-platform
build-platform: ## Build for specific platform (internal target)
	@echo "$(MSG_PREFIX) Building for $(PLATFORM)..."
	@mkdir -p $(BUILD_DIR)
	@GOOS=$(word 1,$(subst /, ,$(PLATFORM))) GOARCH=$(word 2,$(subst /, ,$(PLATFORM))) \
		$(GOBUILD) $(BUILD_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-$(subst /,-,$(PLATFORM))$(if $(findstring windows,$(PLATFORM)),.exe) $(MAIN_PACKAGE)

.PHONY: clean
clean: ## Clean build artifacts
	@$(GOCLEAN)
	@rm -rf $(BUILD_DIR)

.PHONY: fmt
fmt: check-env ## Format Go code
	@$(GOFMT) ./...

.PHONY: test
test: check-env ## Run tests with coverage
	$(call ensure_dirs)
	@$(GOTEST) -v -race -coverprofile=$(COVERAGE_OUT) ./...
	@$(GOCMD) tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "$(SUCCESS) Coverage report: $(COVERAGE_HTML)"

.PHONY: test-integration
test-integration: check-env ## Run integration tests (requires GITHUB_TOKEN)
	@echo "$(MSG_PREFIX) Running integration tests..."
	@GITHUB_TOKEN="$(call get_github_token)" \
		GITHUB_TEST_REPO="$(GITHUB_REPO_ISSUE_CREATION_TEST)" \
		$(GOTEST) -v -tags=integration ./pkg/github/... || true
	@echo "$(SUCCESS) Integration tests completed"

.PHONY: test-short
test-short: check-env ## Run only unit tests (skip integration)
	$(call ensure_dirs)
	@$(GOTEST) -v -short -race -coverprofile=$(COVERAGE_OUT) ./...
	@echo "$(SUCCESS) Unit tests completed"

.PHONY: bench
bench: check-env ## Run benchmarks
	$(call ensure_dirs)
	@$(GOTEST) -bench=. -benchmem ./... | tee $(REPORTS_DIR)/bench.out

.PHONY: lint
lint: check-env ## Run linters
	@$(MAMBA_RUN) pre-commit run --all-files
	@echo "$(SUCCESS) Linting completed"

# ==============================================================================
# EXECUTION & ANALYSIS
# ==============================================================================

.PHONY: run
run: build ## Run with default OpenAPI spec
	@$(BUILD_DIR)/$(BINARY_NAME) analyze $(OPENAPI_URL)

.PHONY: run-local
run-local: build ## Run with local config
	@$(BUILD_DIR)/$(BINARY_NAME) analyze --config configs/config.yaml $(OPENAPI_URL)

.PHONY: run-ollama
run-ollama: build ollama-serve ## Run with Ollama model (no GitHub issues)
	@echo "$(MSG_PREFIX) Running with Ollama (local test, no GitHub issues)..."
	@$(BUILD_DIR)/$(BINARY_NAME) analyze --ai-models ollama --create-issues=false $(OPENAPI_URL)

.PHONY: run-ollama-issues
run-ollama-issues: build ollama-serve ## Run with Ollama and create GitHub issues on test failures
	@echo "$(MSG_PREFIX) Running with Ollama and creating GitHub issues to $(GITHUB_REPO_ISSUE_CREATION_TEST) [op-id: $(OP_ID)]..."
	@GITHUB_TOKEN="$(call get_github_token)" \
		$(BUILD_DIR)/$(BINARY_NAME) analyze \
		--ai-models ollama \
		--github-repo $(GITHUB_REPO_ISSUE_CREATION_TEST) \
		--op-id $(OP_ID) \
		$(OPENAPI_URL)

.PHONY: cleanup-test-issues
cleanup-test-issues: build ## Clean up test issues from GitHub repository (dry-run by default)
	@echo "$(MSG_PREFIX) Cleaning up test issues from $(GITHUB_REPO_ISSUE_CREATION_TEST)..."
	@GITHUB_TOKEN="$(call get_github_token)" $(BUILD_DIR)/$(BINARY_NAME) cleanup \
		--github-repo $(GITHUB_REPO_ISSUE_CREATION_TEST) \
		--dry-run

.PHONY: cleanup-test-issues-confirm
cleanup-test-issues-confirm: build ## Clean up test issues from GitHub repository (actually closes them)
	@echo "$(MSG_PREFIX) Closing test issues in $(GITHUB_REPO_ISSUE_CREATION_TEST)..."
	@GITHUB_TOKEN="$(call get_github_token)" $(BUILD_DIR)/$(BINARY_NAME) cleanup \
		--github-repo $(GITHUB_REPO_ISSUE_CREATION_TEST)

.PHONY: test-endpoint
test-endpoint: build ## Test specific endpoint (requires OP_ID=<operationId>)
	@test -n "$(OP_ID)" || { echo "$(ERROR) Please specify: make test-endpoint OP_ID=getPetById"; exit 1; }
	$(call ensure_dirs)
	@GITHUB_TOKEN="$(call get_github_token)" $(BUILD_DIR)/$(BINARY_NAME) analyze --ai-models ollama --op-id $(OP_ID) --create-issues=false --run-tests=false --output $(REPORTS_DIR)/$(OP_ID)-test.md $(OPENAPI_URL)

.PHONY: test-api
test-api: build ## Test entire API with different OpenAPI spec
	@GITHUB_TOKEN="$(call get_github_token)" $(BUILD_DIR)/$(BINARY_NAME) analyze --ai-models ollama --create-issues=false --run-tests=false $(OPENAPI_URL)

# ==============================================================================
# OLLAMA MANAGEMENT
# ==============================================================================

.PHONY: ollama-serve
ollama-serve: check-env ## Start Ollama server
	@$(MAMBA_RUN) ollama serve &>/dev/null & echo "$(SUCCESS) Ollama server started"

.PHONY: ollama-status
ollama-status: check-env ## Check Ollama status
	@$(MAMBA_RUN) ollama --version 2>/dev/null || echo "$(ERROR) Ollama not available"

.PHONY: ollama-pull-codellama
ollama-pull-codellama: check-env ## Download CodeLlama model
	@$(MAMBA_RUN) ollama pull codellama:7b-instruct

.PHONY: ollama-list
ollama-list: check-env ## List Ollama models
	@$(MAMBA_RUN) ollama list

.PHONY: test-ollama
test-ollama: build ## Test Ollama integration
	@$(BUILD_DIR)/$(BINARY_NAME) models ollama status

# ==============================================================================
# SECURITY & QUALITY
# ==============================================================================

.PHONY: security
security: check-env ## Run security scan
	$(call install_tool,gosec,$(TOOL_GOSEC))
	$(call ensure_dirs)
	$(call run_tool,gosec -fmt json -out $(REPORTS_DIR)/security.json ./... 2>/dev/null || gosec ./...)

.PHONY: vuln-check
vuln-check: check-env ## Check for vulnerabilities
	$(call install_tool,govulncheck,$(TOOL_GOVULN))
	$(call ensure_dirs)
	$(call run_tool,govulncheck ./... | tee $(REPORTS_DIR)/vulnerabilities.out)

# ==============================================================================
# PACKAGING & RELEASE
# ==============================================================================

.PHONY: docker
docker: ## Build Docker image
	@docker build -t glens:latest .

.PHONY: release
release: clean deps test build-all ## Create release builds
	$(call ensure_dirs)
	@cd $(BUILD_DIR) && for file in $(BINARY_NAME)-*; do \
		if [[ $$file == *".exe" ]]; then \
			zip dist/$${file%.exe}.zip $$file; \
		else \
			tar -czf dist/$$file.tar.gz $$file; \
		fi; \
	done
	@cd $(DIST_DIR) && sha256sum * > checksums.txt
	@echo "$(SUCCESS) Release artifacts created in $(DIST_DIR)/"

# ==============================================================================
# DEVELOPMENT UTILITIES
# ==============================================================================

.PHONY: dev
dev: ## Start development mode with hot reload
	$(call install_tool,air,$(TOOL_AIR))
	@$(MAMBA_RUN) air

.PHONY: setup
setup: env ## Setup development environment
	@$(GOMOD) download && $(GOMOD) tidy
	$(call install_tool,golangci-lint,$(TOOL_GOLANGCI))
	$(call install_tool,air,$(TOOL_AIR))
	@echo "$(SUCCESS) Development environment setup complete"

.PHONY: update-deps
update-deps: check-env ## Update dependencies
	@$(GOGET) -u ./... && $(GOMOD) tidy

.PHONY: profile
profile: build ## Generate performance profiles
	$(call ensure_dirs)
	@$(BUILD_DIR)/$(BINARY_NAME) analyze --cpuprofile=$(CPU_PROF) --memprofile=$(MEM_PROF) $(OPENAPI_URL)
	@echo "$(SUCCESS) Profiles generated: $(CPU_PROF), $(MEM_PROF)"

.PHONY: docs
docs: build ## Generate documentation
	$(call ensure_dirs)
	@$(BUILD_DIR)/$(BINARY_NAME) --help > $(DOCS_DIR)/cli-help.txt
	@echo "$(SUCCESS) Documentation generated in $(DOCS_DIR)/"

.PHONY: reports
reports: test lint security vuln-check ## Generate all reports
	@echo "$(SUCCESS) All reports generated in $(REPORTS_DIR)/"

.PHONY: ci
ci: clean fmt lint test build ## Simulate CI pipeline
	@echo "$(SUCCESS) CI pipeline simulation completed"
