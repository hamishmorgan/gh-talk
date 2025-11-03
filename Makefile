.PHONY: help
help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
	go build -o gh-talk

.PHONY: install
install: build ## Install as gh extension (removes old installation first)
	@gh extension remove talk 2>/dev/null || true
	gh extension install .

.PHONY: test
test: ## Run tests
	go test -v -race ./...

.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report: coverage.html"

.PHONY: test-e2e
test-e2e: ## Run E2E tests
	E2E_TEST=1 go test -v -tags=e2e ./e2e/...

.PHONY: lint
lint: ## Run all linters
	gofmt -w .
	goimports -w .
	go vet ./...
	golangci-lint run

.PHONY: lint-fix
lint-fix: ## Fix linting issues
	golangci-lint run --fix

.PHONY: fmt
fmt: ## Format code
	gofmt -w .
	goimports -w .

.PHONY: lint-md
lint-md: ## Lint markdown files
	markdownlint '**/*.md' '**/*.mdc' --ignore node_modules

.PHONY: lint-md-fix
lint-md-fix: ## Fix markdown issues
	markdownlint '**/*.md' '**/*.mdc' --ignore node_modules --fix

.PHONY: clean
clean: ## Clean build artifacts
	rm -f gh-talk
	rm -f coverage.out coverage.html

.PHONY: deps
deps: ## Download dependencies
	go mod download
	go mod tidy

.PHONY: update-deps
update-deps: ## Update dependencies
	go get -u ./...
	go mod tidy

.PHONY: ci
ci: lint test ## Run CI checks locally
	@echo "✓ All CI checks passed"

.PHONY: reinstall
reinstall: install ## Alias for install (always removes and reinstalls)

.PHONY: uninstall
uninstall: ## Remove the gh extension
	gh extension remove talk

.DEFAULT_GOAL := help


