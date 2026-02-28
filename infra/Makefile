# Makefile for Terraform infrastructure
# Targets: fmt, fmt-check, validate, lint, all

.DEFAULT_GOAL := help

.PHONY: all fmt fmt-check validate lint init clean help

all: fmt-check validate ## Run all checks (CI equivalent)

help: ## Show available targets
	@awk 'BEGIN {FS = ":.*##"} /^[a-zA-Z_-]+:.*##/ {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

fmt: ## Format Terraform files
	terraform fmt -recursive .

fmt-check: ## Check formatting (fails if unformatted; same check as CI)
	terraform fmt -check -recursive .

init: ## Initialize Terraform (no backend)
	terraform init -backend=false

validate: init ## Validate Terraform configuration
	terraform validate

lint: fmt-check validate ## Run all linting checks

clean: ## Remove Terraform cache
	rm -rf .terraform .terraform.lock.hcl
