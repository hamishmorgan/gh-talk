# Engineering Practices and CI/CD

**Comprehensive guide to testing, CI/CD, linting, and quality assurance**

## Overview

This document defines engineering practices for gh-talk to ensure high quality, reliability, and maintainability.

## Testing Strategy

### Test Types

#### 1. Unit Tests

**Scope:** Individual functions and methods

**Pattern:**

```go
// internal/api/threads_test.go
package api

import (
    "testing"
)

func TestParseThreadID(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    string
        wantErr bool
    }{
        {
            name:    "valid full ID",
            input:   "PRRT_kwDOQN97u85gQeTN",
            want:    "PRRT_kwDOQN97u85gQeTN",
            wantErr: false,
        },
        {
            name:    "invalid prefix",
            input:   "INVALID_123",
            want:    "",
            wantErr: true,
        },
        {
            name:    "empty string",
            input:   "",
            want:    "",
            wantErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := ParseThreadID(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("ParseThreadID() error = %v, wantErr %v", err, tt.wantErr)
                return
            }
            if got != tt.want {
                t.Errorf("ParseThreadID() = %v, want %v", got, tt.want)
            }
        })
    }
}
```

**Coverage Areas:**

- ID parsing and validation
- Filtering logic
- Formatting functions
- Error handling
- URL parsing
- Emoji mapping

**Target:** 80%+ code coverage

#### 2. Integration Tests

**Scope:** Full command execution with mocked API

**Pattern:**

```go
// internal/commands/list_test.go
package commands

import (
    "bytes"
    "os"
    "testing"
    
    "github.com/hamishmorgan/gh-talk/internal/api"
)

func TestListThreadsCommand(t *testing.T) {
    // Create mock client with test fixtures
    client := api.NewMockClient(api.MockOptions{
        FixtureFile: "../../testdata/pr_with_resolved_threads.json",
    })
    
    // Create command with mock client
    cmd := NewListThreadsCommand(client)
    
    // Capture output
    output := new(bytes.Buffer)
    cmd.SetOut(output)
    cmd.SetErr(output)
    
    // Set args
    cmd.SetArgs([]string{"--pr", "1", "--unresolved"})
    
    // Execute
    err := cmd.Execute()
    if err != nil {
        t.Fatalf("command failed: %v", err)
    }
    
    // Assert output
    outStr := output.String()
    if !bytes.Contains([]byte(outStr), []byte("test_file.go:7")) {
        t.Error("expected unresolved thread at line 7")
    }
    if bytes.Contains([]byte(outStr), []byte("RESOLVED")) {
        t.Error("should not show resolved threads with --unresolved")
    }
}
```

**Coverage Areas:**

- Full command execution
- Flag parsing
- Output formatting
- Error messages
- Context detection

**Target:** All commands tested

#### 3. Contract Tests

**Scope:** Validate GraphQL queries/mutations against schema

**Pattern:**

```go
// internal/api/queries_test.go
package api

import (
    "testing"
)

func TestQueryStructures(t *testing.T) {
    // These tests validate that our query structs match
    // the actual GitHub GraphQL schema
    
    tests := []struct {
        name     string
        query    interface{}
        fixture  string
    }{
        {
            name:    "list threads query",
            query:   &listThreadsQuery{},
            fixture: "testdata/pr_full_response.json",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Load fixture
            data, err := os.ReadFile(tt.fixture)
            if err != nil {
                t.Fatal(err)
            }
            
            // Unmarshal into query struct
            err = json.Unmarshal(data, tt.query)
            if err != nil {
                t.Errorf("query struct doesn't match API response: %v", err)
            }
        })
    }
}
```

**Coverage Areas:**

- Query structures match GitHub schema
- Response parsing works
- All fields accessible

**Target:** All queries/mutations tested

#### 4. End-to-End Tests

**Scope:** Real API calls (optional, can be expensive)

**Pattern:**

```go
// e2e/basic_test.go
// +build e2e

package e2e

import (
    "os"
    "testing"
)

func TestRealAPIListThreads(t *testing.T) {
    if os.Getenv("E2E_TEST") == "" {
        t.Skip("Skipping E2E test (set E2E_TEST=1 to run)")
    }
    
    // Use real API client
    client, err := api.NewClient()
    if err != nil {
        t.Fatal(err)
    }
    
    // Query test PR
    threads, err := client.ListThreads("hamishmorgan", "gh-talk", 1)
    if err != nil {
        t.Fatalf("real API call failed: %v", err)
    }
    
    // Basic assertions on real data
    if len(threads) == 0 {
        t.Error("expected at least 1 thread in test PR")
    }
}
```

**Run with:**

```bash
E2E_TEST=1 go test -tags=e2e ./e2e/...
```

**Target:** Smoke tests for main workflows

### Test Organization

```
gh-talk/
├── internal/
│   ├── api/
│   │   ├── client.go
│   │   ├── client_test.go         # Unit tests
│   │   ├── threads.go
│   │   ├── threads_test.go        # Unit tests
│   │   ├── queries_test.go        # Contract tests
│   │   └── mock.go                # Test mocks
│   ├── commands/
│   │   ├── root.go
│   │   ├── root_test.go
│   │   ├── list.go
│   │   ├── list_test.go           # Integration tests
│   │   └── testhelpers.go         # Test utilities
│   └── format/
│       ├── table.go
│       └── table_test.go          # Unit tests
├── testdata/                       # Test fixtures
│   ├── README.md
│   ├── pr_full_response.json
│   ├── issue_full_response.json
│   └── pr_with_resolved_threads.json
└── e2e/                            # E2E tests (optional)
    └── basic_test.go
```

### Test Commands

```bash
# Run all unit tests
go test ./...

# Run with coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Run specific package
go test ./internal/api/...

# Run with race detector
go test -race ./...

# Verbose output
go test -v ./...

# E2E tests (when E2E_TEST=1)
E2E_TEST=1 go test -tags=e2e ./e2e/...
```

### Test Helpers

**Mock API Client:**

```go
// internal/api/mock.go
package api

type MockClient struct {
    ListThreadsFunc  func(owner, name string, pr int) ([]Thread, error)
    ResolveThreadFunc func(threadID string) error
}

func (m *MockClient) ListThreads(owner, name string, pr int) ([]Thread, error) {
    if m.ListThreadsFunc != nil {
        return m.ListThreadsFunc(owner, name, pr)
    }
    return nil, nil
}

func NewMockClient(opts MockOptions) *MockClient {
    if opts.FixtureFile != "" {
        return newMockFromFixture(opts.FixtureFile)
    }
    return &MockClient{}
}
```

**Fixture Loading:**

```go
// internal/testutil/fixtures.go
package testutil

import (
    "encoding/json"
    "os"
    "testing"
)

func LoadFixture(t *testing.T, path string, v interface{}) {
    t.Helper()
    
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatalf("failed to load fixture %s: %v", path, err)
    }
    
    err = json.Unmarshal(data, v)
    if err != nil {
        t.Fatalf("failed to unmarshal fixture %s: %v", path, err)
    }
}
```

## Linting

### Go Linting

**Tools:**

- **golangci-lint** - Meta-linter running multiple linters
- **gofmt** - Code formatting
- **goimports** - Import organization
- **go vet** - Static analysis

**Configuration: `.golangci.yml`**

```yaml
run:
  timeout: 5m
  tests: true

linters:
  enable:
    - gofmt
    - goimports
    - govet
    - staticcheck
    - errcheck
    - gosimple
    - ineffassign
    - unused
    - typecheck
    - misspell
    - gocyclo
    - dupl
    - goconst
    - gocritic

linters-settings:
  gocyclo:
    min-complexity: 15
  dupl:
    threshold: 100
  goconst:
    min-len: 3
    min-occurrences: 3

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
```

**Pre-commit Script:**

```bash
#!/bin/bash
# .git/hooks/pre-commit

set -e

echo "Running gofmt..."
gofmt -w .

echo "Running goimports..."
goimports -w .

echo "Running go vet..."
go vet ./...

echo "Running golangci-lint..."
golangci-lint run

echo "Running tests..."
go test ./...

echo "✓ All checks passed"
```

**Make executable:**

```bash
chmod +x .git/hooks/pre-commit
```

### Markdown Linting

**Tool:** markdownlint-cli

**Configuration:** See `.markdownlint.json` in repository root

**Rules Disabled:**

- `MD013`: Line length (documentation often has long lines)
- `MD029`: Ordered list prefixes (intentional non-sequential numbering)
- `MD034`: Bare URLs (common in documentation)
- `MD036`: Emphasis as heading (intentional documentation style)
- `MD040`: Code block language (not all code blocks need language tags)
- `MD041`: First line heading (not required)

**Allowed HTML Elements:**

- Standard: `details`, `summary`, `br`, `img`
- Non-standard: `id`, `message`, `emoji` (used for angle-bracket syntax in command examples like `reply <id> <message>` or `react <id> <emoji>`)

**Commands:**

```bash
# Install
npm install -g markdownlint-cli

# Lint markdown files (via Makefile)
make lint-md

# Fix automatically
make lint-md-fix

# Manual commands
markdownlint '**/*.md' '**/*.mdc' --ignore node_modules
markdownlint '**/*.md' '**/*.mdc' --ignore node_modules --fix
```

### YAML Linting

**Tool:** yamllint

**Configuration: `.yamllint.yml`**

```yaml
extends: default

rules:
  line-length:
    max: 120
  indentation:
    spaces: 2
  comments:
    min-spaces-from-content: 1
```

**Commands:**

```bash
# Install
pip install yamllint

# Lint YAML files
yamllint .github/workflows/*.yml
yamllint .golangci.yml
```

## CI/CD Workflows

All workflows run on push to main and on pull requests. Each workflow is focused on a single concern for clarity and maintainability.

### Workflow 1: Test

**File: `.github/workflows/test.yml`**

Runs unit tests with coverage reporting.

- Runs tests with race detection
- Checks coverage threshold (minimum 5%, target 60%)
- Uploads coverage to Codecov

### Workflow 2: Lint

**File: `.github/workflows/lint.yml`**

Runs Go code linting with golangci-lint.

- Uses golangci-lint with multiple enabled linters
- 5 minute timeout
- Checks code quality, style, and common issues

### Workflow 3: Format

**File: `.github/workflows/format.yml`**

Checks code formatting with gofmt and goimports.

- Verifies gofmt compliance
- Verifies goimports compliance
- Fails if any files need formatting

### Workflow 4: Markdown Lint

**File: `.github/workflows/markdown-lint.yml`**

Lints all markdown files (.md and .mdc).

- Checks markdown formatting and style
- Uses markdownlint-cli with project configuration
- Validates documentation quality

### Workflow 5: Build Check

**File: `.github/workflows/build.yml`**

```yaml
name: Build

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    name: Build
    strategy:
      matrix:
        os: [ubuntu-latest, macos-latest, windows-latest]
        go: ['1.21', '1.22', '1.23']
    runs-on: ${{ matrix.os }}
    
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go }}
          cache: true
      
      - name: Download dependencies
        run: go mod download
      
      - name: Build
        run: go build -v
      
      - name: Verify binary
        run: ./gh-talk --version || ./gh-talk.exe --version
        shell: bash
```

### Workflow 6: Release

**File: `.github/workflows/release.yml`**

```yaml
name: release
on:
  push:
    tags:
      - "v*"
permissions:
  contents: write
  id-token: write
  attestations: write

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: cli/gh-extension-precompile@v2
        with:
          generate_attestations: true
          go_version_file: go.mod
```

**Builds For:**

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

**Triggered by:**

```bash
git tag v0.1.0
git push origin v0.1.0
```

### Workflow 7: PR Validation

**File: `.github/workflows/pr.yml`**

```yaml
name: PR Checks

on:
  pull_request:
    types: [opened, synchronize, reopened]

jobs:
  validate:
    name: Validate PR
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Full history for comparison
      
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
          cache: true
      
      - name: Check go.mod and go.sum
        run: |
          go mod tidy
          git diff --exit-code go.mod go.sum
      
      - name: Check for large files
        run: |
          if git diff --name-only origin/main | xargs ls -lh | awk '$5 ~ /M$/ {print; exit 1}'; then
            echo "Large files detected"
            exit 1
          fi
      
      - name: Run all tests
        run: go test -v -race ./...
      
      - name: Run linters
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
  
  docs-check:
    name: Documentation Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Check broken links in docs
        uses: lycheeverse/lychee-action@v2
        with:
          args: --verbose --no-progress '**/*.md'
          fail: true
      
      - name: Lint markdown
        uses: nosborn/github-action-markdown-cli@v3.3.0
        with:
          files: .
          config_file: .markdownlint.json
```

### Workflow 8: Dependency Updates

**File: `.github/workflows/dependabot-auto-merge.yml`**

```yaml
name: Dependabot Auto-merge

on:
  pull_request:
    types: [opened, synchronize]

permissions:
  contents: write
  pull-requests: write

jobs:
  auto-merge:
    name: Auto-merge Dependabot PRs
    runs-on: ubuntu-latest
    if: github.actor == 'dependabot[bot]'
    steps:
      - uses: actions/checkout@v4
      
      - name: Check if tests pass
        run: go test ./...
      
      - name: Auto-merge patch updates
        if: |
          contains(github.event.pull_request.title, 'bump') &&
          contains(github.event.pull_request.title, 'patch')
        run: gh pr merge --auto --squash "${{ github.event.pull_request.number }}"
        env:
          GH_TOKEN: ${{ secrets.GITHUB_TOKEN }}
```

### Dependabot Configuration

**File: `.github/dependabot.yml`**

```yaml
version: 2
updates:
  - package-ecosystem: gomod
    directory: /
    schedule:
      interval: weekly
    open-pull-requests-limit: 10
    groups:
      patch-updates:
        update-types:
          - patch
```

## Linting Configuration

### golangci-lint Configuration

**File: `.golangci.yml`**

```yaml
run:
  timeout: 5m
  tests: true
  skip-dirs:
    - vendor
    - testdata

output:
  formats:
    - format: colored-line-number
  print-issued-lines: true
  print-linter-name: true
  sort-results: true

linters:
  enable:
    # Enabled by default
    - errcheck
    - gosimple
    - govet
    - ineffassign
    - staticcheck
    - unused
    
    # Additional linters
    - gofmt
    - goimports
    - misspell
    - goconst
    - gocyclo
    - dupl
    - gocritic
    - revive
    - bodyclose
    - noctx
    - unparam
    - wastedassign
    - whitespace
  
  disable:
    - typecheck  # Can be slow

linters-settings:
  gocyclo:
    min-complexity: 15
  
  goconst:
    min-len: 3
    min-occurrences: 3
  
  dupl:
    threshold: 100
  
  gocritic:
    enabled-tags:
      - diagnostic
      - style
      - performance
    disabled-checks:
      - paramTypeCombine  # Can be noisy

issues:
  exclude-use-default: false
  max-issues-per-linter: 0
  max-same-issues: 0
  
  exclude-rules:
    # Exclude test files from some linters
    - path: _test\.go
      linters:
        - dupl
        - goconst
```

### Markdown Linting

**Configuration:** See `.markdownlint.json` in repository root for the complete configuration.

**Files to Lint:**

- README.md
- AGENTS.md
- docs/\*.md
- testdata/README.md
- .cursor/commands/\*.md
- .cursor/rules/\*.mdc

**Configuration highlights:**

- Disabled rules: MD013, MD029, MD034, MD036, MD040, MD041
- Allowed HTML elements include `id`, `message`, `emoji` for angle-bracket syntax in docs

### EditorConfig

**File: `.editorconfig`**

```ini
root = true

[*]
charset = utf-8
end_of_line = lf
insert_final_newline = true
trim_trailing_whitespace = true

[*.go]
indent_style = tab
indent_size = 4

[*.{yml,yaml}]
indent_style = space
indent_size = 2

[*.md]
indent_style = space
indent_size = 2
trim_trailing_whitespace = false

[Makefile]
indent_style = tab
```

## Makefile

**File: `Makefile`**

> **Note:** Makefiles require TAB characters for indentation, not spaces. The examples below show single spaces for readability in markdown, but you must use TAB characters when creating actual Makefile rules. See the actual `Makefile` in the repository for the correct format.

```makefile
.PHONY: help
help: ## Show this help
 @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## Build the binary
 go build -o gh-talk

.PHONY: install
install: ## Install as gh extension
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

.PHONY: lint-md
lint-md: ## Lint markdown files
 markdownlint '**/*.md' --ignore node_modules

.PHONY: lint-md-fix
lint-md-fix: ## Fix markdown issues
 markdownlint '**/*.md' --ignore node_modules --fix

.PHONY: fmt
fmt: ## Format code
 gofmt -w .
 goimports -w .

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

.DEFAULT_GOAL := help
```

**Usage:**

```bash
make help          # Show available commands
make build         # Build binary
make test          # Run tests
make lint          # Run linters
make ci            # Run all CI checks locally
```

## Environment Variables

### Supported Variables

**From gh CLI (Automatically Used):**

```bash
GH_TOKEN              # GitHub auth token
GH_HOST               # GitHub host (default: github.com)
GH_REPO               # Current repository (OWNER/REPO)
GH_FORCE_TTY          # Force terminal mode
GH_DEBUG              # Enable debug logging
NO_COLOR              # Disable colors
CLICOLOR              # Color support
CLICOLOR_FORCE        # Force colors
```

**gh-talk Specific:**

```bash
GH_TALK_CONFIG        # Config file location (default: ~/.config/gh-talk/config.yml)
GH_TALK_CACHE_DIR     # Cache directory (default: ~/.cache/gh-talk)
GH_TALK_CACHE_TTL     # Cache TTL in minutes (default: 5)
GH_TALK_FORMAT        # Default output format (table, json, tsv)
GH_TALK_EDITOR        # Editor for message composition (falls back to $EDITOR)
```

### Environment Variable Handling

**Implementation:**

```go
// internal/config/env.go
package config

import (
    "os"
    "strconv"
    "time"
)

type Env struct {
    // gh CLI variables (via go-gh)
    Token       string
    Host        string
    Repo        string
    Debug       bool
    ForceTTY    bool
    
    // gh-talk specific
    ConfigPath  string
    CacheDir    string
    CacheTTL    time.Duration
    Format      string
    Editor      string
}

func FromEnvironment() Env {
    env := Env{
        Token:      os.Getenv("GH_TOKEN"),
        Host:       os.Getenv("GH_HOST"),
        Repo:       os.Getenv("GH_REPO"),
        Debug:      os.Getenv("GH_DEBUG") == "1",
        ForceTTY:   os.Getenv("GH_FORCE_TTY") != "",
        ConfigPath: os.Getenv("GH_TALK_CONFIG"),
        CacheDir:   os.Getenv("GH_TALK_CACHE_DIR"),
        Format:     os.Getenv("GH_TALK_FORMAT"),
        Editor:     os.Getenv("GH_TALK_EDITOR"),
    }
    
    // Parse cache TTL
    if ttlStr := os.Getenv("GH_TALK_CACHE_TTL"); ttlStr != "" {
        if minutes, err := strconv.Atoi(ttlStr); err == nil {
            env.CacheTTL = time.Duration(minutes) * time.Minute
        }
    }
    
    // Defaults
    if env.ConfigPath == "" {
        env.ConfigPath = defaultConfigPath()
    }
    if env.CacheDir == "" {
        env.CacheDir = defaultCacheDir()
    }
    if env.CacheTTL == 0 {
        env.CacheTTL = 5 * time.Minute
    }
    if env.Format == "" {
        env.Format = "table"
    }
    if env.Editor == "" {
        env.Editor = os.Getenv("EDITOR")
    }
    if env.Editor == "" {
        env.Editor = "vim"
    }
    
    return env
}
```

### Documentation

**File: `docs/ENVIRONMENT.md`**

```markdown
# Environment Variables

## GitHub CLI Variables

These are automatically used by go-gh:

- `GH_TOKEN` - GitHub authentication token
- `GH_HOST` - GitHub host (default: github.com)
- `GH_REPO` - Repository context (OWNER/REPO)
- `GH_FORCE_TTY` - Force terminal mode (any value)
- `GH_DEBUG` - Enable debug logging (set to 1)

## Terminal Variables

- `NO_COLOR` - Disable all colors
- `CLICOLOR` - Color support (0 or 1)
- `CLICOLOR_FORCE` - Force colors (any value)
- `TERM` - Terminal type
- `COLORTERM` - True color support

## gh-talk Specific Variables

- `GH_TALK_CONFIG` - Config file location
  - Default: `~/.config/gh-talk/config.yml`
  - Example: `export GH_TALK_CONFIG=/path/to/config.yml`

- `GH_TALK_CACHE_DIR` - Cache directory
  - Default: `~/.cache/gh-talk`
  - Example: `export GH_TALK_CACHE_DIR=/tmp/gh-talk-cache`

- `GH_TALK_CACHE_TTL` - Cache TTL in minutes
  - Default: `5`
  - Example: `export GH_TALK_CACHE_TTL=10`

- `GH_TALK_FORMAT` - Default output format
  - Values: `table`, `json`, `tsv`
  - Default: `table`
  - Example: `export GH_TALK_FORMAT=json`

- `GH_TALK_EDITOR` - Editor for message composition
  - Default: Value of `$EDITOR`, falls back to `vim`
  - Example: `export GH_TALK_EDITOR=nano`

## Examples

### Disable Colors
```bash
NO_COLOR=1 gh talk list threads
```

### Force JSON Output

```bash
GH_TALK_FORMAT=json gh talk list threads
```

### Use Custom Editor

```bash
GH_TALK_EDITOR=code gh talk reply --editor
```

### Debug Mode

```bash
GH_DEBUG=1 gh talk list threads
```

### Custom Cache Location

```bash
GH_TALK_CACHE_DIR=/tmp/cache gh talk list threads
```

```

## Code Coverage

### Coverage Goals

**Targets:**
- Overall: 80%+
- API package: 90%+ (critical)
- Commands: 70%+ (harder to test)
- Format/Filter: 85%+

### Coverage Commands

```bash
# Generate coverage
go test -coverprofile=coverage.out ./...

# View in terminal
go tool cover -func=coverage.out

# View in browser
go tool cover -html=coverage.out

# Per-package coverage
go test -cover ./internal/api
go test -cover ./internal/commands
go test -cover ./internal/format
```

### Coverage in CI

```yaml
# .github/workflows/test.yml (already shown above)
- name: Run tests with coverage
  run: go test -v -race -coverprofile=coverage.out ./...

- name: Check coverage threshold
  run: |
    total=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    if (( $(echo "$total < 80" | bc -l) )); then
      echo "Coverage $total% is below 80%"
      exit 1
    fi

- name: Upload to Codecov
  uses: codecov/codecov-action@v4
  with:
    file: ./coverage.out
```

## Quality Gates

### Pre-Merge Requirements

**All PRs Must:**

1. ✅ Pass all tests
2. ✅ Pass all linters
3. ✅ Pass format check (gofmt, goimports)
4. ✅ Maintain 80%+ coverage
5. ✅ Build on all platforms
6. ✅ Pass markdown linting
7. ✅ No broken links in docs

**GitHub Branch Protection:**

```yaml
# Settings → Branches → main
Protection rules:
  - Require pull request reviews (1+)
  - Require status checks to pass:
    - test
    - lint
    - format-check
    - build (all matrix combinations)
    - docs-check
  - Require conversation resolution before merge
  - Require linear history (rebase)
```

## Documentation Standards

### Code Documentation

**Package Documentation:**

```go
// Package api provides GitHub GraphQL API client for PR conversation management.
//
// The Client type wraps go-gh's GraphQLClient and provides domain-specific
// methods for interacting with review threads, comments, and reactions.
//
// Example:
//
// client, err := api.NewClient()
// if err != nil {
//     return err
// }
//
// threads, err := client.ListThreads("owner", "repo", 123)
package api
```

**Function Documentation:**

```go
// ListThreads fetches all review threads for a pull request.
//
// It queries the GitHub GraphQL API and returns threads with their
// comments, reactions, and metadata. Results are not filtered -
// use the filter package for client-side filtering.
//
// Parameters:
//   - owner: Repository owner
//   - name: Repository name
//   - pr: Pull request number
//
// Returns:
//   - []Thread: List of review threads
//   - error: API error or nil
//
// Example:
//
// threads, err := client.ListThreads("hamishmorgan", "gh-talk", 1)
// if err != nil {
//     return err
// }
func (c *Client) ListThreads(owner, name string, pr int) ([]Thread, error) {
    // Implementation
}
```

### Markdown Documentation

**Required Files:**

- `README.md` - User-facing quick start
- `AGENTS.md` - AI agent guidelines
- `docs/*.md` - Technical documentation
- `CHANGELOG.md` - Version history (after first release)
- `CONTRIBUTING.md` - Contribution guidelines (optional)

**Standards:**

- Clear headings
- Code examples
- Links between docs
- Keep up to date

## Git Practices

### Commit Messages

**Format:**

```
<type>: <subject>

<body>

<footer>
```

**Types:**

- `feat` - New feature
- `fix` - Bug fix
- `docs` - Documentation only
- `style` - Formatting (no code change)
- `refactor` - Code restructuring
- `test` - Adding tests
- `chore` - Maintenance

**Examples:**

```
feat: add list threads command

Implement basic list threads functionality with filtering
support for resolved/unresolved status.

Closes #15

---

fix: handle empty thread list gracefully

Previously would panic on empty thread list. Now returns
empty array and helpful message.

Fixes #23

---

docs: add Cobra implementation guide

Comprehensive guide to using Cobra for gh-talk commands
with examples and patterns.
```

### AI Agent Attribution

**For AI commits (as per AGENTS.md):**

```bash
git commit --author="Cursor <cursor@noreply.local>" -m "feat: implement list threads command"
```

### Branch Strategy

**Main Branch:**

- Protected
- All changes via PR
- Linear history (rebase)

**Feature Branches:**

```bash
feat/list-command
fix/thread-resolution
docs/update-readme
refactor/api-client
```

**Release Branches:**

```bash
release/v0.1.0
release/v0.2.0
```

## Pre-release Checklist

**Before v0.1.0:**

- [ ] All MVP commands implemented
- [ ] Comprehensive tests (80%+ coverage)
- [ ] All linters passing
- [ ] Documentation complete
- [ ] README examples work
- [ ] Manual testing on real PRs
- [ ] Cross-platform builds work
- [ ] Shell completion generated
- [ ] CHANGELOG.md created
- [ ] Release notes written

## Monitoring and Metrics

### Test Metrics

**Track:**

- Coverage percentage
- Test execution time
- Flaky test rate
- Test count

**Report in CI:**

```yaml
- name: Test Metrics
  run: |
    echo "Test Count: $(go test -v ./... | grep -c 'RUN')"
    echo "Coverage: $(go tool cover -func=coverage.out | grep total | awk '{print $3}')"
```

### Build Metrics

**Track:**

- Build time
- Binary size
- Dependency count

**Report:**

```yaml
- name: Build Metrics
  run: |
    time go build
    ls -lh gh-talk
    go list -m all | wc -l
```

## Security

### Dependency Scanning

**Dependabot already configured** (in release.yml)

**Additional: govulncheck**

```yaml
# .github/workflows/security.yml
name: Security

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]
  schedule:
    - cron: '0 0 * * 0'  # Weekly

jobs:
  vuln-check:
    name: Vulnerability Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - uses: actions/setup-go@v5
        with:
          go-version-file: go.mod
      
      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest
      
      - name: Run govulncheck
        run: govulncheck ./...
```

### Secret Scanning

**GitHub Features:**

- Secret scanning (enabled by default)
- Push protection
- Dependency review

**In Code:**

```go
// Never log tokens
if debug {
    log.Printf("Query: %s", query)  // ✅ OK
    log.Printf("Token: %s", token)  // ❌ NEVER
}
```

## Summary

### Complete Engineering Setup

**Testing:**

- ✅ Unit tests with table-driven patterns
- ✅ Integration tests with fixtures
- ✅ Contract tests for GraphQL
- ✅ E2E tests (optional)
- ✅ 80%+ coverage target
- ✅ Coverage reporting

**Linting:**

- ✅ golangci-lint (15+ linters)
- ✅ gofmt, goimports
- ✅ markdownlint for docs
- ✅ yamllint for configs
- ✅ Pre-commit hooks

**CI/CD:**

- ✅ Test workflow (all tests + coverage)
- ✅ Lint workflow (Go + Markdown)
- ✅ Build workflow (multi-platform)
- ✅ PR validation workflow
- ✅ Release workflow (already exists)
- ✅ Security scanning

**Quality Gates:**

- ✅ Branch protection
- ✅ Required status checks
- ✅ Coverage threshold
- ✅ Format validation

**Tooling:**

- ✅ Makefile for common tasks
- ✅ EditorConfig for consistency
- ✅ Dependabot for updates
- ✅ Pre-commit hooks

**Documentation:**

- ✅ Code documentation (godoc)
- ✅ README (user-facing)
- ✅ Technical docs (docs/)
- ✅ Environment variables (docs/ENVIRONMENT.md)
- ✅ Changelog (after v0.1.0)

### Implementation Order

**Phase 0: Setup (Before Coding)**

1. Add .golangci.yml
2. Add .markdownlint.json
3. Add .editorconfig
4. Add Makefile
5. Create CI workflows
6. Create docs/ENVIRONMENT.md

**Phase 1: MVP (With Testing)**

1. Implement feature
2. Write unit tests (same PR)
3. Write integration test
4. Verify coverage

**Phase 2: Polish**

1. Add pre-commit hooks
2. Security scanning
3. Performance tests
4. Documentation review

---

**Last Updated**: 2025-11-02  
**Status**: Engineering practices defined and ready to implement  
**Coverage Target**: 80%+ overall, 90%+ for critical packages
