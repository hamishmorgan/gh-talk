# Agent Instructions

Instructions for AI agents working with the gh-talk GitHub CLI extension.

## Project Context

`gh-talk` is a GitHub CLI extension for managing PR and Issue conversations from the terminal.

**Core Features:**

- Reply to review threads
- Add emoji reactions
- Resolve/unresolve threads
- Hide comments
- Dismiss reviews
- Filter and list conversations

**Repository Structure:**

```
gh-talk/
├── cmd/           # Command implementations (list, reply, react, resolve)
├── pkg/           # Public library code (api, filter, format)
├── internal/      # Private application code (config, cache)
├── main.go        # Entry point
└── SPEC.md        # Complete specification
```

## Dependencies

**Required:**

- Go 1.21+
- `gh` CLI v2.0+
- Git 2.0+

**Key Libraries:**

- `github.com/cli/go-gh` - GitHub CLI library
- `github.com/spf13/cobra` - CLI framework

## Documentation Standards

- Use formal, minimal tone
- Use technically precise language
- Eliminate unnecessary words
- Be direct and concise

## Go Code Standards

### Style Guide

Follow [Effective Go](https://go.dev/doc/effective_go) and [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md).

**Formatting:**

```bash
gofmt -w .
goimports -w .
```

**Naming:**

- **Packages**: lowercase, single word (`api`, `filter`)
- **Files**: lowercase with underscores (`thread_list.go`)
- **Functions**: `MixedCaps` (exported), `mixedCaps` (unexported)
- **Variables**: `camelCase`
- **Constants**: `MixedCaps`
- **Interfaces**: `-er` suffix for single-method

**Package Organization:**

```go
package cmd

import (
    "context"
    "fmt"
    
    "github.com/spf13/cobra"
    
    "github.com/hamishmorgan/gh-talk/pkg/api"
)

const defaultTimeout = 30 * time.Second

var rootCmd = &cobra.Command{
    Use:   "gh-talk",
    Short: "Manage GitHub PR conversations",
}

type Config struct {
    Timeout time.Duration
}

func Execute() error {
    return rootCmd.Execute()
}
```

### Comments

Only add comments when they provide non-obvious information.

**Good comments explain:**

- **Why** something is done
- **Context** not clear from code
- **Workarounds** for bugs
- **Public API** (godoc)

```go
// Bad: Redundant
// Set timeout to 30 seconds
timeout := 30 * time.Second

// Good: Provides context
// GitHub GraphQL API has 30s default timeout; extend for large queries
timeout := 60 * time.Second

// Good: Package documentation
// Package api provides GitHub GraphQL API client for PR conversation management.
package api
```

### Clean Code Principles

**1. Meaningful Names**

```go
// Bad
func GetT(id string) (*T, error)

// Good
func GetThread(id string) (*Thread, error)
```

**2. Single Responsibility**

```go
// Bad: Too many responsibilities
type Manager struct {
    // Handles API, formatting, caching, logging
}

// Good: Separate concerns
type APIClient struct { }
type Formatter struct { }
type Cache struct { }
```

**3. Keep It Simple**

```go
// Bad: Overly complex
result := map[bool]string{true: "resolved", false: "unresolved"}[thread.IsResolved]

// Good: Clear
var result string
if thread.IsResolved {
    result = "resolved"
} else {
    result = "unresolved"
}
```

**4. Fail Fast**

```go
func ProcessThread(ctx context.Context, id string) error {
    if id == "" {
        return fmt.Errorf("thread ID required")
    }
    if ctx == nil {
        return fmt.Errorf("context required")
    }
    
    // Main logic here
    return nil
}
```

### Error Handling

**Always wrap errors with context:**

```go
func FetchThread(id string) (*Thread, error) {
    resp, err := api.Get(id)
    if err != nil {
        return nil, fmt.Errorf("fetch thread %s: %w", id, err)
    }
    
    thread, err := parseResponse(resp)
    if err != nil {
        return nil, fmt.Errorf("parse response: %w", err)
    }
    
    return thread, nil
}
```

**Return errors, don't panic:**

```go
// Bad: Library code shouldn't panic
func MustGetThread(id string) *Thread {
    thread, err := GetThread(id)
    if err != nil {
        panic(err)
    }
    return thread
}

// Good: Return errors
func GetThread(id string) (*Thread, error) {
    // implementation
}
```

**Custom error types when needed:**

```go
type NotFoundError struct {
    ID string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("thread %s not found", e.ID)
}

// Check error types
if errors.As(err, &NotFoundError{}) {
    // Handle not found
}
```

## Git Commit Attribution

**All AI agent commits must use `--author` flag:**

```bash
# Cursor AI
git commit --author="Cursor <cursor@noreply.local>" -m "Implement GraphQL client"

# GitHub Copilot  
git commit --author="Copilot <copilot@noreply.local>" -m "Add error handling"

# Claude
git commit --author="Claude <claude@noreply.local>" -m "Refactor API types"
```

**Purpose:**

- Maintain transparency about AI contributions
- Identify which AI tool made changes
- Allow filtering commits by author
- Distinguish automated from manual work

## Testing

### Test Structure

```go
func TestClient_ResolveThread(t *testing.T) {
    tests := []struct {
        name    string
        id      string
        wantErr bool
    }{
        {"valid thread ID", "PRRT_abc123", false},
        {"empty thread ID", "", true},
        {"invalid prefix", "INVALID_123", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            c := NewClient("token")
            err := c.ResolveThread(context.Background(), tt.id)
            if (err != nil) != tt.wantErr {
                t.Errorf("ResolveThread() error = %v, wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Test Categories

1. **Unit Tests** - Test individual functions
2. **Integration Tests** - Test complete workflows with mocked API
3. **Contract Tests** - Validate GraphQL query structure

### Best Practices

- Use table-driven tests
- Test error paths thoroughly
- Mock external dependencies
- Keep tests fast (< 1s per test)
- Use `testdata/` for fixtures

## GitHub CLI Extension Guidelines

### Requirements

- Executable named `gh-talk` (invoked as `gh talk`)
- Handle `--help` flag
- Exit codes: 0 for success, non-zero for errors
- Respect `GH_TOKEN`, `GH_HOST` environment variables

### Using go-gh Library

```go
import "github.com/cli/go-gh/v2/pkg/api"

// Get authenticated GraphQL client
client, err := api.DefaultGraphQLClient()
if err != nil {
    return err
}

// Execute GraphQL query
var response struct {
    Repository struct {
        PullRequest struct {
            ReviewThreads struct {
                Nodes []Thread
            }
        }
    }
}

err = client.Query("repo", query, &response, variables)
```

### Command Structure

```go
import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
    Use:   "gh-talk",
    Short: "Manage GitHub PR conversations",
}

var listCmd = &cobra.Command{
    Use:   "list",
    Short: "List review threads",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Implementation
        return nil
    },
}

func init() {
    rootCmd.AddCommand(listCmd)
}
```

### Error Messages

```go
// Print errors to stderr
if err != nil {
    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
    os.Exit(1)
}

// Provide helpful context
if errors.Is(err, ErrNotFound) {
    fmt.Fprintf(os.Stderr, "Thread not found. Run 'gh talk list' to see available threads.\n")
    os.Exit(1)
}
```

## Code Quality

**Pre-commit checks:**

```bash
gofmt -w .           # Format code
goimports -w .       # Organize imports
golangci-lint run    # Lint
go test ./...        # Test
go vet ./...         # Vet
```

---

**Version**: 0.1.0  
**Last Updated**: 2025-11-02
