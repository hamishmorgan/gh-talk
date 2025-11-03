# go-gh Library Guide

**Comprehensive guide to using `go-gh` for GitHub CLI extensions**

## Overview

`go-gh` is the official Go library for authoring GitHub CLI extensions. It provides authenticated API clients, terminal utilities, and helpers that automatically integrate with `gh` CLI conventions.

**Repository:** https://github.com/cli/go-gh  
**Documentation:** https://pkg.go.dev/github.com/cli/go-gh/v2  
**Examples:** https://github.com/cli/go-gh/blob/trunk/example_gh_test.go

## Why go-gh?

### Automatic Integration with `gh`

**go-gh modules obey GitHub CLI conventions by default:**

1. **Repository Context** - `repository.Current()` respects:
   - `GH_REPO` environment variable
   - Git remote configuration (fallback)

2. **Authentication** - API clients automatically use:
   - `GH_TOKEN` environment variable
   - `GH_HOST` environment variable
   - User's stored OAuth token (from `gh auth`)

3. **Terminal Capabilities** - Determined from environment:
   - `GH_FORCE_TTY` - Force terminal mode
   - `NO_COLOR` - Disable colors
   - `CLICOLOR`, `CLICOLOR_FORCE` - Color control
   - `TERM`, `COLORTERM` - Terminal type

4. **Output Formatting** - Same engines as `gh`:
   - Table printer with auto-truncation
   - Go template support
   - JSON formatting

5. **Browser Integration** - Activates user's preferred browser

### Benefits Over Manual Implementation

✅ **No Authentication Code** - Handles tokens automatically  
✅ **No Context Detection** - Finds current repo from git  
✅ **No Terminal Handling** - Detects capabilities automatically  
✅ **No HTTP Boilerplate** - Clean API client interface  
✅ **Consistent with `gh`** - Users get expected behavior  

## Core Modules

### 1. API Clients (`pkg/api`)

**Two Client Types:**
- `GraphQLClient` - For GraphQL API (recommended for gh-talk)
- `RESTClient` - For REST API (simpler endpoints)

#### GraphQL Client

**Basic Usage:**
```go
import "github.com/cli/go-gh/v2/pkg/api"

// Create client with defaults (uses gh auth)
client, err := api.DefaultGraphQLClient()
if err != nil {
    return err
}

// Define response structure
var query struct {
    Repository struct {
        PullRequest struct {
            ReviewThreads struct {
                Nodes []struct {
                    ID         string
                    IsResolved bool
                }
            } `graphql:"reviewThreads(first: $first)"`
        } `graphql:"pullRequest(number: $number)"`
    } `graphql:"repository(owner: $owner, name: $name)"`
}

// Set variables
variables := map[string]interface{}{
    "owner":  graphql.String("hamishmorgan"),
    "name":   graphql.String("gh-talk"),
    "number": graphql.Int(1),
    "first":  graphql.Int(20),
}

// Execute query
err = client.Query("ReviewThreads", &query, variables)
if err != nil {
    return err
}

// Access results
threads := query.Repository.PullRequest.ReviewThreads.Nodes
```

**With Custom Options:**
```go
opts := api.ClientOptions{
    EnableCache:    true,              // Enable response caching
    CacheTTL:       5 * time.Minute,   // Cache for 5 minutes
    Timeout:        30 * time.Second,  // Request timeout
    Log:            os.Stdout,         // Log requests
    LogColorize:    true,              // Colorize logs
    LogVerboseHTTP: true,              // Log headers/bodies
}

client, err := api.NewGraphQLClient(opts)
```

#### GraphQL Mutations

**Example from go-gh:**
```go
client, err := api.DefaultGraphQLClient()

// Define mutation structure
var mutation struct {
    ResolveReviewThread struct {
        Thread struct {
            ID         string
            IsResolved bool
        }
    } `graphql:"resolveReviewThread(input: $input)"`
}

// Define input type
type ResolveInput struct {
    ThreadID string `json:"threadId"`
}

variables := map[string]interface{}{
    "input": ResolveInput{
        ThreadID: "PRRT_kwDOQN97u85gQeTN",
    },
}

// Execute mutation
err = client.Mutate("ResolveReviewThread", &mutation, variables)
if err != nil {
    return err
}

// Access result
thread := mutation.ResolveReviewThread.Thread
fmt.Printf("Thread %s resolved: %v\n", thread.ID, thread.IsResolved)
```

#### GraphQL Pagination

**Pattern from go-gh examples:**
```go
client, err := api.DefaultGraphQLClient()

var query struct {
    Repository struct {
        PullRequest struct {
            ReviewThreads struct {
                Nodes []struct {
                    ID string
                }
                PageInfo struct {
                    HasNextPage bool
                    EndCursor   string
                }
            } `graphql:"reviewThreads(first: $first, after: $endCursor)"`
        } `graphql:"pullRequest(number: $number)"`
    } `graphql:"repository(owner: $owner, name: $name)"`
}

variables := map[string]interface{}{
    "owner":     graphql.String("owner"),
    "name":      graphql.String("repo"),
    "number":    graphql.Int(1),
    "first":     graphql.Int(30),
    "endCursor": (*graphql.String)(nil), // nil for first page
}

allThreads := []struct{ ID string }{}

for {
    err := client.Query("ReviewThreads", &query, variables)
    if err != nil {
        return err
    }
    
    // Collect results
    allThreads = append(allThreads, query.Repository.PullRequest.ReviewThreads.Nodes...)
    
    // Check for more pages
    if !query.Repository.PullRequest.ReviewThreads.PageInfo.HasNextPage {
        break
    }
    
    // Update cursor for next page
    variables["endCursor"] = graphql.String(
        query.Repository.PullRequest.ReviewThreads.PageInfo.EndCursor,
    )
}

fmt.Printf("Fetched %d threads across multiple pages\n", len(allThreads))
```

#### REST Client

**Basic Usage:**
```go
client, err := api.DefaultRESTClient()
if err != nil {
    return err
}

// GET request
response := []struct{ Name string }{}
err = client.Get("repos/cli/cli/tags", &response)
if err != nil {
    return err
}
```

**When to Use REST vs GraphQL:**

| Use REST For | Use GraphQL For |
|-------------|-----------------|
| Simple GET requests | Complex nested data |
| Standard endpoints | Custom queries |
| File downloads | Filtering and pagination |
| Legacy API support | Efficient data fetching |

**For gh-talk:** Use GraphQL (better for threads/comments/reactions)

### 2. Repository Context (`pkg/repository`)

**Get Current Repository:**
```go
import "github.com/cli/go-gh/v2/pkg/repository"

// Automatically determines from:
// 1. GH_REPO environment variable
// 2. Git remotes in current directory
repo, err := repository.Current()
if err != nil {
    return err
}

fmt.Printf("Host: %s\n", repo.Host)     // github.com
fmt.Printf("Owner: %s\n", repo.Owner)   // hamishmorgan
fmt.Printf("Name: %s\n", repo.Name)     // gh-talk
```

**Parse Repository String:**
```go
// Flexible formats:
// - "OWNER/REPO"
// - "HOST/OWNER/REPO"
// - "https://github.com/OWNER/REPO"

repo, err := repository.Parse("cli/cli")
repo, err := repository.Parse("github.com/cli/cli")
repo, err := repository.Parse("https://github.com/cli/cli")

// All produce same result:
// repo.Host  = "github.com"
// repo.Owner = "cli"
// repo.Name  = "cli"
```

**Parse with Default Host:**
```go
// If string doesn't include host, use provided default
repo, err := repository.ParseWithHost("owner/repo", "github.example.com")
// repo.Host = "github.example.com"
```

**Use in gh-talk:**
```go
// Auto-detect current PR's repository
repo, err := repository.Current()
if err != nil {
    // Fall back to explicit --repo flag
}

// Use in GraphQL query
variables := map[string]interface{}{
    "owner": graphql.String(repo.Owner),
    "name":  graphql.String(repo.Name),
}
```

### 3. Terminal (`pkg/term`)

**Initialize from Environment:**
```go
import "github.com/cli/go-gh/v2/pkg/term"

terminal := term.FromEnv()

// Check if connected to terminal
if terminal.IsTerminalOutput() {
    // Show fancy tables with colors
} else {
    // Output machine-readable format (TSV, JSON)
}

// Get terminal dimensions
width, height, err := terminal.Size()

// Check color support
if terminal.IsColorEnabled() {
    // Use ANSI color codes
}

// Access standard streams
stdin := terminal.In()      // io.Reader
stdout := terminal.Out()    // io.Writer
stderr := terminal.ErrOut() // io.Writer
```

**Use in gh-talk:**
```go
terminal := term.FromEnv()

// Adapt output format
if terminal.IsTerminalOutput() {
    // Show pretty tables
    renderTable(threads, terminal.Out())
} else {
    // Output JSON for scripts
    json.NewEncoder(terminal.Out()).Encode(threads)
}
```

### 4. Table Printer (`pkg/tableprinter`)

**Basic Table:**
```go
import "github.com/cli/go-gh/v2/pkg/tableprinter"

terminal := term.FromEnv()
width, _, _ := terminal.Size()

// Create table printer
t := tableprinter.New(
    terminal.Out(),
    terminal.IsTerminalOutput(),
    width,
)

// Add header
t.AddField("ID")
t.AddField("File")
t.AddField("Line")
t.AddField("Status")
t.EndRow()

// Add data rows
for _, thread := range threads {
    t.AddField(thread.ShortID)
    t.AddField(thread.Path)
    t.AddField(fmt.Sprintf("%d", thread.Line))
    t.AddField(thread.Status)
    t.EndRow()
}

// Render (auto-formats for terminal vs non-terminal)
if err := t.Render(); err != nil {
    return err
}
```

**With Formatting:**
```go
// Color function
green := func(s string) string {
    return "\x1b[32m" + s + "\x1b[m"
}

// Add field with color (auto-disabled in non-TTY)
t.AddField("RESOLVED", tableprinter.WithColor(green))

// Prevent truncation for important fields
t.AddField(thread.ID, tableprinter.WithTruncate(nil))

// Custom truncation
truncate := func(maxWidth int, s string) string {
    if len(s) > maxWidth {
        return s[:maxWidth-3] + "..."
    }
    return s
}
t.AddField(thread.Body, tableprinter.WithTruncate(truncate))
```

**Output Modes:**
- **Terminal**: Pretty-printed columns, colors, auto-truncated to fit width
- **Non-Terminal**: TSV format, no truncation, no colors (perfect for piping)

### 5. Prompter (`pkg/prompter`)

**Interactive Selection:**
```go
import "github.com/cli/go-gh/v2/pkg/prompter"

p := prompter.New(os.Stdin, os.Stdout, os.Stderr)

// Select single option
options := []string{"Thread 1", "Thread 2", "Thread 3"}
selectedIndex, err := p.Select("Which thread?", "", options)
if err != nil {
    return err
}
fmt.Printf("Selected: %s\n", options[selectedIndex])

// Multi-select
indices, err := p.MultiSelect("Select threads to resolve", nil, options)

// Text input
message, err := p.Input("Reply message:", "")

// Confirm action
confirmed, err := p.Confirm("Resolve this thread?", true)
```

**Use in gh-talk:**
```go
// Interactive thread selection when ID not provided
if threadID == "" {
    threads := listThreads()
    options := make([]string, len(threads))
    for i, t := range threads {
        options[i] = fmt.Sprintf("%s:%d - %s", t.Path, t.Line, t.Preview)
    }
    
    idx, err := p.Select("Select thread:", "", options)
    if err != nil {
        return err
    }
    threadID = threads[idx].ID
}

// Get reply message interactively
if message == "" {
    message, err = p.Input("Reply message:", "")
    if err != nil {
        return err
    }
}
```

### 6. Exec (`gh.Exec`)

**Shell Out to `gh` Commands:**
```go
import gh "github.com/cli/go-gh/v2"

// Execute gh command and capture output
stdout, stderr, err := gh.Exec("issue", "list", "-R", "cli/cli", "--limit", "5")
if err != nil {
    return err
}

fmt.Println(stdout.String())
```

**Use Cases:**
- Leverage existing `gh` commands
- Don't reimplement functionality
- Combine with custom logic

**Example for gh-talk:**
```go
// Use gh to get current PR number
stdout, _, err := gh.Exec("pr", "view", "--json", "number")
if err != nil {
    return fmt.Errorf("no current PR found")
}

var result struct{ Number int }
json.Unmarshal(stdout.Bytes(), &result)
prNumber := result.Number
```

## Error Handling

### GraphQL Errors

**Error Type:**
```go
type GraphQLError struct {
    Errors []GraphQLErrorItem
}

type GraphQLErrorItem struct {
    Message    string
    Locations  []struct { Line, Column int }
    Path       []interface{}
    Extensions map[string]interface{}
    Type       string  // "NOT_FOUND", "FORBIDDEN", etc.
}
```

**Handling:**
```go
var gqlErr *api.GraphQLError
if errors.As(err, &gqlErr) {
    for _, e := range gqlErr.Errors {
        switch e.Type {
        case "NOT_FOUND":
            return fmt.Errorf("thread not found: %s", e.Message)
        case "FORBIDDEN":
            return fmt.Errorf("permission denied: %s", e.Message)
        default:
            return fmt.Errorf("GraphQL error: %s", e.Message)
        }
    }
}
```

**Match Specific Errors:**
```go
// Check for specific error type on specific path
if gqlErr.Match("NOT_FOUND", "resolveReviewThread") {
    return fmt.Errorf("thread not found or already deleted")
}

// Match subpaths with trailing "."
if gqlErr.Match("FORBIDDEN", "repository.pullRequest.") {
    return fmt.Errorf("no access to this PR")
}
```

### HTTP Errors

**Error Type:**
```go
type HTTPError struct {
    Errors     []HTTPErrorItem
    Headers    http.Header
    Message    string
    RequestURL *url.URL
    StatusCode int
}

type HTTPErrorItem struct {
    Code     string
    Field    string
    Message  string
    Resource string
}
```

**Common Status Codes:**
```go
var httpErr *api.HTTPError
if errors.As(err, &httpErr) {
    switch httpErr.StatusCode {
    case 401:
        return fmt.Errorf("authentication failed")
    case 403:
        return fmt.Errorf("forbidden: check permissions")
    case 404:
        return fmt.Errorf("resource not found")
    case 422:
        return fmt.Errorf("validation failed: %s", httpErr.Message)
    case 502, 503, 504:
        return fmt.Errorf("GitHub API unavailable, try again")
    default:
        return fmt.Errorf("HTTP %d: %s", httpErr.StatusCode, httpErr.Message)
    }
}
```

## Best Practices for gh-talk

### 1. Client Initialization

**Recommended Pattern:**
```go
package api

import (
    "github.com/cli/go-gh/v2/pkg/api"
)

// Client wraps the GraphQL client with our methods
type Client struct {
    graphql *api.GraphQLClient
}

// NewClient creates a new API client using gh defaults
func NewClient() (*Client, error) {
    gql, err := api.DefaultGraphQLClient()
    if err != nil {
        return nil, fmt.Errorf("create GraphQL client: %w", err)
    }
    
    return &Client{graphql: gql}, nil
}

// NewClientWithOptions creates client with custom options (for testing)
func NewClientWithOptions(opts api.ClientOptions) (*Client, error) {
    gql, err := api.NewGraphQLClient(opts)
    if err != nil {
        return nil, fmt.Errorf("create GraphQL client: %w", err)
    }
    
    return &Client{graphql: gql}, nil
}
```

### 2. Query Organization

**Recommended Structure:**
```go
package api

// Types match GraphQL schema
type ReviewThread struct {
    ID         string
    IsResolved bool
    Path       string
    Line       int
    Comments   struct {
        Nodes []ReviewComment
    }
}

type ReviewComment struct {
    ID     string
    Body   string
    Author struct {
        Login string
    }
}

// ListThreads fetches review threads for a PR
func (c *Client) ListThreads(owner, name string, number int) ([]ReviewThread, error) {
    var query struct {
        Repository struct {
            PullRequest struct {
                ReviewThreads struct {
                    Nodes []ReviewThread
                } `graphql:"reviewThreads(first: 100)"`
            } `graphql:"pullRequest(number: $number)"`
        } `graphql:"repository(owner: $owner, name: $name)"`
    }
    
    variables := map[string]interface{}{
        "owner":  graphql.String(owner),
        "name":   graphql.String(name),
        "number": graphql.Int(number),
    }
    
    err := c.graphql.Query("ListThreads", &query, variables)
    if err != nil {
        return nil, fmt.Errorf("query threads: %w", err)
    }
    
    return query.Repository.PullRequest.ReviewThreads.Nodes, nil
}
```

### 3. Repository Context Detection

**Recommended Pattern:**
```go
package commands

import (
    "github.com/cli/go-gh/v2/pkg/repository"
)

// GetRepository gets repo from flag or auto-detects
func GetRepository(repoFlag string) (owner, name string, err error) {
    var repo repository.Repository
    
    if repoFlag != "" {
        // Use explicit --repo flag
        repo, err = repository.Parse(repoFlag)
        if err != nil {
            return "", "", fmt.Errorf("invalid repository: %w", err)
        }
    } else {
        // Auto-detect from git context
        repo, err = repository.Current()
        if err != nil {
            return "", "", fmt.Errorf("not in a git repository, use --repo flag")
        }
    }
    
    return repo.Owner, repo.Name, nil
}
```

### 4. Terminal-Aware Output

**Recommended Pattern:**
```go
package format

import (
    "encoding/json"
    "github.com/cli/go-gh/v2/pkg/tableprinter"
    "github.com/cli/go-gh/v2/pkg/term"
}

// OutputThreads renders threads appropriately for context
func OutputThreads(threads []Thread, format string) error {
    terminal := term.FromEnv()
    
    // Explicit format
    if format == "json" {
        return json.NewEncoder(terminal.Out()).Encode(threads)
    }
    
    // Auto-detect: TTY = table, non-TTY = TSV
    if terminal.IsTerminalOutput() {
        return renderTable(threads, terminal)
    }
    return renderTSV(threads, terminal.Out())
}

func renderTable(threads []Thread, terminal term.Term) error {
    width, _, _ := terminal.Size()
    t := tableprinter.New(terminal.Out(), true, width)
    
    // Headers
    t.AddField("ID")
    t.AddField("File:Line")
    t.AddField("Status")
    t.AddField("Preview")
    t.EndRow()
    
    // Rows
    for _, thread := range threads {
        t.AddField(thread.ShortID)
        t.AddField(fmt.Sprintf("%s:%d", thread.Path, thread.Line))
        t.AddField(formatStatus(thread))
        t.AddField(thread.Preview, tableprinter.WithTruncate(truncate))
        t.EndRow()
    }
    
    return t.Render()
}
```

### 5. Error Messages

**User-Friendly Errors:**
```go
func handleError(err error) error {
    var gqlErr *api.GraphQLError
    if errors.As(err, &gqlErr) {
        for _, e := range gqlErr.Errors {
            switch e.Type {
            case "NOT_FOUND":
                return fmt.Errorf(
                    "thread not found\n\n" +
                    "The thread may have been deleted or the ID is invalid.\n" +
                    "Run 'gh talk list threads' to see available threads.",
                )
            case "FORBIDDEN":
                return fmt.Errorf(
                    "permission denied\n\n" +
                    "You don't have permission to resolve this thread.\n" +
                    "Only the PR author, comment author, or repo admins can resolve threads.",
                )
            }
        }
    }
    
    return fmt.Errorf("operation failed: %w", err)
}
```

## Advanced Features

### Caching

**Enable Response Caching:**
```go
opts := api.ClientOptions{
    EnableCache: true,
    CacheTTL:    5 * time.Minute,  // Cache for 5 minutes
    CacheDir:    "/custom/cache",  // Optional, defaults to gh cache dir
}

client, err := api.NewGraphQLClient(opts)

// Subsequent identical queries served from cache
// Reduces API calls and improves performance
```

**Use Cases:**
- List commands (threads don't change often)
- User data (names, avatars)
- Repository metadata

**Don't Cache:**
- Real-time data (isResolved status)
- After mutations (data changed)

### Context Support

**Using Context for Cancellation:**
```go
import "context"

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

// Use context-aware methods
err := client.QueryWithContext(ctx, "QueryName", &query, variables)
err := client.MutateWithContext(ctx, "MutationName", &mutation, variables)

// Supports:
// - Timeouts
// - Cancellation
// - Deadline propagation
```

### Logging

**Debug Logging:**
```go
opts := api.ClientOptions{
    Log:            os.Stderr,      // Where to log
    LogColorize:    true,           // Color output
    LogVerboseHTTP: true,           // Include headers/bodies
    LogIgnoreEnv:   false,          // Respect GH_DEBUG env var
}

client, err := api.NewGraphQLClient(opts)

// Logs all requests when GH_DEBUG=1 or when enabled
// Example log output:
// → POST https://api.github.com/graphql
// → Headers: ...
// → Body: {"query": "...", "variables": {...}}
// ← 200 OK (245ms)
// ← Body: {"data": {...}}
```

**Use for Development:**
```bash
# Enable debug logging
export GH_DEBUG=1
gh talk list threads

# Or in code for specific client
opts := api.ClientOptions{Log: os.Stderr, LogVerboseHTTP: true}
```

### Custom HTTP Transport

**For Testing:**
```go
// Mock transport for tests
mockTransport := &mockRoundTripper{
    responses: map[string]*http.Response{
        "graphql": {
            StatusCode: 200,
            Body:       io.NopCloser(strings.NewReader(`{"data": {...}}`)),
        },
    },
}

opts := api.ClientOptions{
    Transport: mockTransport,
    AuthToken: "test-token",
    Host:      "github.com",
}

client, err := api.NewGraphQLClient(opts)
// Client now uses mock transport
```

## Package Overview

### Available Packages

| Package | Purpose | Use in gh-talk |
|---------|---------|----------------|
| `api` | GraphQL/REST clients | ✅ Core - all API calls |
| `repository` | Repo context parsing | ✅ High - auto-detect repo |
| `term` | Terminal capabilities | ✅ High - output formatting |
| `tableprinter` | Table formatting | ✅ High - list commands |
| `prompter` | Interactive prompts | ✅ Medium - optional UX |
| `auth` | Authentication | ⚠️ Low - handled by api |
| `config` | gh config access | ⚠️ Low - mostly internal |
| `browser` | Open URLs | ⚠️ Low - maybe for --web |
| `jq` | JSON processing | ❌ Not needed (use encoding/json) |
| `template` | Go templates | ❌ Not needed (use text/template) |
| `markdown` | Markdown rendering | ❌ Not needed |
| `jsonpretty` | JSON formatting | ❌ Not needed |

### Less Common Packages

**browser:** Open URLs in user's browser
```go
import "github.com/cli/go-gh/v2/pkg/browser"

// Open thread in browser
url := fmt.Sprintf("https://github.com/%s/%s/pull/%d#discussion_r%s", 
    owner, repo, pr, discussionID)
browser.Open(url)
```

**config:** Access gh configuration
```go
import "github.com/cli/go-gh/v2/pkg/config"

cfg, err := config.Read(nil)
editor, _ := cfg.Get([]string{"editor"})
```

## GraphQL Type Annotations

### Struct Tags for GraphQL

**Pattern:**
```go
type Query struct {
    Repository struct {
        PullRequest struct {
            ReviewThreads struct {
                Nodes []Thread
                PageInfo struct {
                    HasNextPage bool
                    EndCursor   string
                }
            } `graphql:"reviewThreads(first: $first, after: $after)"`
        } `graphql:"pullRequest(number: $number)"`
    } `graphql:"repository(owner: $owner, name: $name)"`
}
```

**Rules:**
- Use `graphql` struct tags for field mapping
- Parameters in parentheses with `$` prefix
- Nested structs for nested queries
- Array types for lists

**Variable Types:**
```go
variables := map[string]interface{}{
    "owner":  graphql.String("value"),      // String
    "number": graphql.Int(123),             // Int
    "first":  graphql.Int(50),              // Int
    "after":  (*graphql.String)(nil),       // Nullable String
    "input":  InputStruct{...},             // Input object
}
```

**Inline Fragments:**
```go
type UnionType struct {
    TypeName string `graphql:"__typename"`
    OnIssue  struct {
        Title string
    } `graphql:"... on Issue"`
    OnPullRequest struct {
        Number int
    } `graphql:"... on PullRequest"`
}
```

## Common Patterns

### Pattern 1: Query with Pagination

```go
func (c *Client) ListAllThreads(owner, name string, pr int) ([]Thread, error) {
    var allThreads []Thread
    var endCursor *graphql.String
    
    for {
        var query struct {
            Repository struct {
                PullRequest struct {
                    ReviewThreads struct {
                        Nodes    []Thread
                        PageInfo struct {
                            HasNextPage bool
                            EndCursor   string
                        }
                    } `graphql:"reviewThreads(first: 100, after: $after)"`
                } `graphql:"pullRequest(number: $number)"`
            } `graphql:"repository(owner: $owner, name: $name)"`
        }
        
        variables := map[string]interface{}{
            "owner":  graphql.String(owner),
            "name":   graphql.String(name),
            "number": graphql.Int(pr),
            "after":  endCursor,
        }
        
        err := c.graphql.Query("ListThreads", &query, variables)
        if err != nil {
            return nil, err
        }
        
        threads := query.Repository.PullRequest.ReviewThreads
        allThreads = append(allThreads, threads.Nodes...)
        
        if !threads.PageInfo.HasNextPage {
            break
        }
        
        cursor := graphql.String(threads.PageInfo.EndCursor)
        endCursor = &cursor
    }
    
    return allThreads, nil
}
```

### Pattern 2: Mutation with Input

```go
func (c *Client) ResolveThread(threadID string) error {
    var mutation struct {
        ResolveReviewThread struct {
            Thread struct {
                ID         string
                IsResolved bool
            }
        } `graphql:"resolveReviewThread(input: $input)"`
    }
    
    type ResolveInput struct {
        ThreadID string `json:"threadId"`
    }
    
    variables := map[string]interface{}{
        "input": ResolveInput{
            ThreadID: threadID,
        },
    }
    
    err := c.graphql.Mutate("ResolveThread", &mutation, variables)
    if err != nil {
        return fmt.Errorf("resolve thread: %w", err)
    }
    
    return nil
}
```

### Pattern 3: Context Detection

```go
func GetPRNumber(prArg string) (int, error) {
    // Explicit PR number
    if prArg != "" {
        return strconv.Atoi(prArg)
    }
    
    // Infer from current branch
    stdout, _, err := gh.Exec("pr", "view", "--json", "number")
    if err != nil {
        return 0, fmt.Errorf("no PR found for current branch")
    }
    
    var result struct{ Number int }
    if err := json.Unmarshal(stdout.Bytes(), &result); err != nil {
        return 0, err
    }
    
    return result.Number, nil
}
```

### Pattern 4: Terminal-Adaptive Output

```go
func OutputResults(data interface{}, jsonFields []string) error {
    terminal := term.FromEnv()
    
    // JSON requested explicitly
    if len(jsonFields) > 0 {
        filtered := filterFields(data, jsonFields)
        return json.NewEncoder(terminal.Out()).Encode(filtered)
    }
    
    // Auto-detect output mode
    if terminal.IsTerminalOutput() {
        // Human-readable table
        return renderTable(data, terminal)
    } else {
        // Machine-readable TSV
        return renderTSV(data, terminal.Out())
    }
}
```

### Pattern 5: Error Handling Chain

```go
func (c *Client) ReplyToThread(threadID, message string) error {
    err := c.addReply(threadID, message)
    if err != nil {
        return c.handleError(err, threadID)
    }
    return nil
}

func (c *Client) handleError(err error, threadID string) error {
    var gqlErr *api.GraphQLError
    if errors.As(err, &gqlErr) {
        if gqlErr.Match("NOT_FOUND", "addPullRequestReviewThreadReply") {
            return &ThreadNotFoundError{ID: threadID}
        }
        if gqlErr.Match("FORBIDDEN", "") {
            return &PermissionError{Operation: "reply", ID: threadID}
        }
    }
    
    return fmt.Errorf("API error: %w", err)
}
```

## Testing with go-gh

### Mocking Clients

**For Unit Tests:**
```go
package api_test

import (
    "net/http"
    "testing"
    
    "github.com/cli/go-gh/v2/pkg/api"
)

func TestListThreads(t *testing.T) {
    // Create mock transport
    mockTransport := &mockRoundTripper{
        response: &http.Response{
            StatusCode: 200,
            Body:       loadFixture(t, "testdata/pr_full_response.json"),
        },
    }
    
    // Create client with mock
    opts := api.ClientOptions{
        Transport: mockTransport,
        AuthToken: "test-token",
        Host:      "github.com",
    }
    
    client, err := NewClientWithOptions(opts)
    if err != nil {
        t.Fatal(err)
    }
    
    // Test using fixture data
    threads, err := client.ListThreads("owner", "repo", 1)
    if err != nil {
        t.Fatal(err)
    }
    
    if len(threads) != 2 {
        t.Errorf("expected 2 threads, got %d", len(threads))
    }
}
```

### Using Test Fixtures

**Load Real Responses:**
```go
func loadFixture(t *testing.T, path string) io.ReadCloser {
    data, err := os.ReadFile(path)
    if err != nil {
        t.Fatal(err)
    }
    return io.NopCloser(bytes.NewReader(data))
}

// Use in tests
func TestParseThreads(t *testing.T) {
    // testdata/pr_full_response.json from our testing
    fixture := loadFixture(t, "testdata/pr_full_response.json")
    
    var response struct {
        Data struct {
            Repository struct {
                PullRequest struct {
                    ReviewThreads struct {
                        Nodes []Thread
                    }
                }
            }
        }
    }
    
    err := json.NewDecoder(fixture).Decode(&response)
    if err != nil {
        t.Fatal(err)
    }
    
    threads := response.Data.Repository.PullRequest.ReviewThreads.Nodes
    
    // Assertions on real data
    if len(threads) != 2 {
        t.Errorf("expected 2 threads, got %d", len(threads))
    }
}
```

## Environment Variables

### Supported by go-gh

| Variable | Purpose | Default |
|----------|---------|---------|
| `GH_TOKEN` | Auth token | From `gh auth` |
| `GH_HOST` | GitHub host | `github.com` |
| `GH_REPO` | Current repository | From git remotes |
| `GH_FORCE_TTY` | Force terminal mode | Auto-detect |
| `GH_DEBUG` | Enable debug logging | Off |
| `NO_COLOR` | Disable colors | Auto-detect |
| `CLICOLOR` | Color support | Auto-detect |
| `CLICOLOR_FORCE` | Force colors | Off |

### Use in gh-talk

**Respect All Conventions:**
```go
// Auto-uses GH_TOKEN
client, _ := api.DefaultGraphQLClient()

// Auto-detects GH_REPO
repo, _ := repository.Current()

// Auto-adapts to terminal
terminal := term.FromEnv()

// Auto-respects GH_DEBUG
// (no code needed, built into api.ClientOptions)
```

## ClientOptions Reference

### All Available Options

```go
type ClientOptions struct {
    // Authentication
    AuthToken        string              // Override auth token
    
    // Host
    Host             string              // GitHub host (default: github.com)
    UnixDomainSocket string              // Unix socket for requests
    
    // Caching
    EnableCache      bool                // Enable response caching
    CacheDir         string              // Cache directory
    CacheTTL         time.Duration       // Cache time-to-live
    
    // Headers
    Headers          map[string]string   // Custom headers
    SkipDefaultHeaders bool              // Don't set default headers
    
    // Logging
    Log              io.Writer           // Log destination
    LogIgnoreEnv     bool                // Ignore GH_DEBUG env var
    LogColorize      bool                // Colorize log output
    LogVerboseHTTP   bool                // Log headers and bodies
    
    // Performance
    Timeout          time.Duration       // Request timeout
    Transport        http.RoundTripper   // Custom transport (testing)
}
```

### Typical Configurations

**Development (with logging):**
```go
opts := api.ClientOptions{
    Log:            os.Stderr,
    LogColorize:    true,
    LogVerboseHTTP: true,
    Timeout:        30 * time.Second,
}
```

**Production (with caching):**
```go
opts := api.ClientOptions{
    EnableCache: true,
    CacheTTL:    5 * time.Minute,
    Timeout:     30 * time.Second,
}
```

**Testing (with mocks):**
```go
opts := api.ClientOptions{
    Transport: mockTransport,
    AuthToken: "test-token",
    Host:      "github.com",
}
```

## gh-talk Implementation Strategy

### Recommended Architecture

```go
// internal/api/client.go
type Client struct {
    graphql *api.GraphQLClient
}

func NewClient() (*Client, error) {
    gql, err := api.DefaultGraphQLClient()
    if err != nil {
        return nil, err
    }
    return &Client{graphql: gql}, nil
}

// internal/api/threads.go
func (c *Client) ListThreads(ctx context.Context, owner, name string, pr int) ([]Thread, error) {
    // Implementation
}

// internal/api/mutations.go
func (c *Client) ResolveThread(ctx context.Context, threadID string) error {
    // Implementation
}

// internal/format/table.go
func RenderThreads(threads []Thread, format string) error {
    terminal := term.FromEnv()
    // Implementation
}

// cmd/list.go (future)
func runList(cmd *cobra.Command, args []string) error {
    client, _ := api.NewClient()
    repo, _ := repository.Current()
    threads, _ := client.ListThreads(ctx, repo.Owner, repo.Name, prNumber)
    return format.RenderThreads(threads, formatFlag)
}
```

### Error Handling Strategy

```go
// Define custom error types
type ThreadNotFoundError struct {
    ID string
}

func (e *ThreadNotFoundError) Error() string {
    return fmt.Sprintf("thread %s not found", e.ID)
}

// Wrap go-gh errors into custom types
func handleAPIError(err error) error {
    var gqlErr *api.GraphQLError
    if errors.As(err, &gqlErr) {
        for _, e := range gqlErr.Errors {
            switch e.Type {
            case "NOT_FOUND":
                // Extract ID from error context if possible
                return &ThreadNotFoundError{ID: "unknown"}
            case "FORBIDDEN":
                return &PermissionError{}
            }
        }
    }
    return err
}
```

## Best Practices Summary

### DO Use

✅ **`api.DefaultGraphQLClient()`** - Automatic authentication  
✅ **`repository.Current()`** - Auto-detect repository  
✅ **`term.FromEnv()`** - Terminal-aware output  
✅ **`tableprinter.New()`** - Professional tables  
✅ **`gh.Exec()`** - Leverage existing gh commands  
✅ **Context methods** - Support timeouts and cancellation  
✅ **Error type checking** - Use `errors.As()` for go-gh errors  

### DON'T

❌ **Manual HTTP clients** - Use go-gh clients instead  
❌ **Hardcode github.com** - Respect `GH_HOST`  
❌ **Ignore `GH_REPO`** - Use `repository.Current()`  
❌ **Print to stdout directly** - Use `term.Out()`  
❌ **Reimplement auth** - Trust go-gh's auth  
❌ **Skip error types** - Check for GraphQLError, HTTPError  

### Testing

✅ **Use testdata/ fixtures** - Real API responses  
✅ **Mock transport** - For client tests  
✅ **Test with real IDs** - From testdata/ files  
✅ **Test error paths** - Mock error responses  

## References

- **go-gh Repository:** https://github.com/cli/go-gh
- **API Documentation:** https://pkg.go.dev/github.com/cli/go-gh/v2
- **Examples:** https://github.com/cli/go-gh/blob/trunk/example_gh_test.go
- **shurcooL/graphql:** https://github.com/shurcooL/graphql (underlying GraphQL library)

### Related Documentation

- [API.md](API.md) - GitHub API capabilities
- [REAL-DATA.md](REAL-DATA.md) - Real response structures
- [GH-CLI.md](GH-CLI.md) - How gh CLI works
- [STRUCTURE.md](STRUCTURE.md) - Project structure

---

**Last Updated**: 2025-11-02  
**go-gh Version**: v2.12.2  
**Context**: Implementation guide for gh-talk


