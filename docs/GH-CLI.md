# GitHub CLI (`gh`) Capabilities and Analysis

**Comprehensive guide to understanding the `gh` CLI and its relevance to `gh-talk`**

## Overview

The GitHub CLI (`gh`) is the official command-line interface for GitHub, providing seamless integration with GitHub from the terminal. Understanding its capabilities, strengths, and limitations is crucial for building `gh-talk` as a complementary extension.

**Official Documentation:** [https://cli.github.com/manual](https://cli.github.com/manual)

## Core Architecture

### Command Structure

```bash
gh <command> <subcommand> [flags]
```

**Command Categories:**

- **Core Commands**: auth, browse, codespace, gist, issue, org, pr, project, release, repo
- **GitHub Actions**: cache, run, workflow
- **Extension Commands**: User-installed extensions (like gh-talk)
- **Additional Commands**: alias, api, attestation, completion, config, extension, gpg-key, label, preview, ruleset, search, secret, ssh-key, status, variable

### Extension Model

**Key Characteristics:**

```bash
# Extensions are repositories named gh-<extname>
# They contain an executable: gh-<extname>
# All arguments passed through: gh talk <args> → gh-talk <args>
```

**Extension Rules:**

- Cannot override core `gh` commands
- Must start with `gh-` prefix
- Executed directly as binaries
- Auto-checked for updates every 24 hours
- Can use `gh extension exec <extname>` to force execution

**Discovery:** <https://github.com/topics/gh-extension>

## PR Commands (`gh pr`)

### General Commands

#### `gh pr create`

**Purpose:** Create a pull request

**Capabilities:**

- Interactive mode (prompts for details)
- Can auto-fill from commits
- Supports draft PRs
- Can add reviewers, assignees, labels, projects

**Limitations:**

- Doesn't support creating review threads
- No line-specific commenting during creation

#### `gh pr list`

**Purpose:** List pull requests in a repository

**JSON Fields Available:**

```
additions, assignees, author, autoMergeRequest, baseRefName, baseRefOid,
body, changedFiles, closed, closedAt, closingIssuesReferences, comments,
commits, createdAt, deletions, files, fullDatabaseId, headRefName,
headRefOid, headRepository, headRepositoryOwner, id, isCrossRepository,
isDraft, labels, latestReviews, maintainerCanModify, mergeCommit,
mergeStateStatus, mergeable, mergedAt, mergedBy, milestone, number,
potentialMergeCommit, projectCards, projectItems, reactionGroups,
reviewDecision, reviewRequests, reviews, state, statusCheckRollup,
title, updatedAt, url
```

**Notable:**

- `reactionGroups` - Shows emoji reactions summary
- `reviews` - Shows review summaries
- `comments` - Shows comment count
- No direct access to review threads or individual thread comments

#### `gh pr view`

**Purpose:** Display PR details

**Flags:**

- `--comments` - View PR comments (top-level only)
- `--json` - Output as JSON
- `--web` - Open in browser

**Limitations:**

- `--comments` shows only top-level PR comments
- Does NOT show review thread comments
- No way to view individual review threads
- Cannot filter by resolved/unresolved status

### Targeted Commands

#### `gh pr comment`

**Purpose:** Add a comment to a pull request

**Capabilities:**

```bash
gh pr comment <pr> --body "Comment text"
gh pr comment <pr> --editor              # Open editor
gh pr comment <pr> --web                 # Open browser
gh pr comment <pr> --edit-last           # Edit your last comment
gh pr comment <pr> --delete-last         # Delete your last comment
```

**Type of Comment:** Top-level PR comment only

**Limitations:**

- ❌ Cannot reply to review threads
- ❌ Cannot add line-specific comments
- ❌ Cannot create new review threads
- ❌ Only adds general PR comments
- ❌ No emoji reaction support

#### `gh pr review`

**Purpose:** Add a review to a pull request

**Capabilities:**

```bash
gh pr review <pr> --approve              # Approve
gh pr review <pr> --comment              # General comment
gh pr review <pr> --request-changes      # Request changes
gh pr review <pr> --body "Review text"   # Add review body
```

**Review Types:**

- Approve (`--approve`)
- Comment (`--comment`)
- Request Changes (`--request-changes`)

**Limitations:**

- ❌ Cannot add line-specific comments
- ❌ Cannot create review threads
- ❌ Only creates review-level comments
- ❌ No way to interact with existing review threads
- ❌ Cannot resolve/unresolve threads
- ❌ No emoji reactions

**Workflow Gap:**
The review command submits a review but cannot add the line-by-line comments that make up review threads. Those must be added via the web UI or API.

#### `gh pr checks`

**Purpose:** Show CI status

**Note:** Useful for context but not conversation-related.

#### Other PR Commands

- `gh pr checkout` - Check out PR locally
- `gh pr close/reopen` - Change PR state
- `gh pr diff` - View changes
- `gh pr edit` - Edit PR metadata
- `gh pr merge` - Merge PR
- `gh pr ready` - Mark as ready for review
- `gh pr lock/unlock` - Lock conversation
- `gh pr update-branch` - Update PR branch

## Issue Commands (`gh issue`)

### General Commands

#### `gh issue create`

**Purpose:** Create a new issue

**Capabilities:**

- Interactive prompts
- Template support
- Add labels, assignees, projects

#### `gh issue list`

**Purpose:** List issues

**JSON Fields Available:**

```
assignees, author, body, closed, closedAt, comments, createdAt, id,
labels, milestone, number, projectCards, projectItems, reactionGroups,
state, title, updatedAt, url
```

**Notable:**

- `reactionGroups` - Emoji reaction summary
- `comments` - Comment count
- No access to individual comment details

#### `gh issue view`

**Purpose:** View issue details

**JSON Fields Available:**

```
assignees, author, body, closed, closedAt, closedByPullRequestsReferences,
comments, createdAt, id, isPinned, labels, milestone, number, projectCards,
projectItems, reactionGroups, state, stateReason, title, updatedAt, url
```

**Limitations:**

- Cannot view individual comments with JSON output
- `reactionGroups` shows summary, not per-comment reactions
- No comment filtering or sorting

### Targeted Commands

#### `gh issue comment`

**Purpose:** Add a comment to an issue

**Capabilities:**

```bash
gh issue comment <issue> --body "Comment"
gh issue comment <issue> --editor
gh issue comment <issue> --web
gh issue comment <issue> --edit-last
gh issue comment <issue> --delete-last
```

**Similarities to `gh pr comment`:**

- Same interface
- Same editing capabilities
- Same limitations

**Limitations:**

- ❌ No emoji reactions
- ❌ Cannot minimize/hide comments
- ❌ Cannot view comment metadata (author, timestamp)
- ❌ No threading support

#### Other Issue Commands

- `gh issue close/reopen` - Change state
- `gh issue edit` - Edit metadata
- `gh issue lock/unlock` - Lock conversation
- `gh issue pin/unpin` - Pin to repository
- `gh issue transfer` - Move to another repo

## API Command (`gh api`)

### Purpose

Direct access to GitHub's REST and GraphQL APIs with authentication handled automatically.

### Capabilities

**REST API:**

```bash
# GET request
gh api repos/{owner}/{repo}/pulls

# POST request
gh api repos/{owner}/{repo}/issues/123/comments -f body='Comment'

# With parameters
gh api -X GET search/issues -f q='repo:cli/cli is:open'
```

**GraphQL API:**

```bash
gh api graphql -f query='
  query {
    repository(owner: "owner", name: "repo") {
      pullRequest(number: 123) {
        title
      }
    }
  }
'
```

**Advanced Features:**

- `--paginate` - Auto-fetch all pages
- `--slurp` - Combine pages into array
- `--cache` - Cache responses
- `--jq` - Filter with jq
- `--template` - Format with Go templates
- Placeholder replacement: `{owner}`, `{repo}`, `{branch}`

### JSON Processing

**Built-in jq Support:**

```bash
gh api repos/{owner}/{repo}/pulls --jq '.[].title'
```

**Go Template Support:**

```bash
gh api repos/{owner}/{repo}/pulls --template '{{range .}}{{.title}}{{"\n"}}{{end}}'
```

**Template Functions:**

- `autocolor` - Colorize for terminals
- `color <style> <input>` - Apply colors
- `join <sep> <list>` - Join arrays
- `pluck <field> <list>` - Extract field from objects
- `tablerow` / `tablerender` - Create tables
- `timeago` - Relative timestamps
- `timefmt` - Format timestamps
- `truncate` - Limit length
- `hyperlink` - Terminal hyperlinks

### Pagination

**GraphQL Pagination:**

```bash
gh api graphql --paginate -f query='
  query($endCursor: String) {
    repository(owner: "owner", name: "repo") {
      pullRequests(first: 100, after: $endCursor) {
        nodes { number title }
        pageInfo { hasNextPage endCursor }
      }
    }
  }
'
```

**Requirements:**

- Query must accept `$endCursor: String`
- Must fetch `pageInfo{ hasNextPage, endCursor }`

**REST Pagination:**

- Automatically follows `Link` headers
- Use `--paginate` flag

### Authentication

**Automatic:**

- Uses `GH_TOKEN` or `GITHUB_TOKEN` environment variable
- Uses `gh auth` credentials
- Supports GitHub Enterprise with `GH_ENTERPRISE_TOKEN`

**Custom Host:**

```bash
gh api --hostname github.example.com /repos/owner/repo
```

## Formatting and Output

### JSON Output

**Available on Many Commands:**

```bash
gh pr list --json number,title,author
gh issue view <num> --json body,comments,reactionGroups
```

**Requirements:**

- Must specify fields: `--json field1,field2`
- Fields are command-specific
- Run with `--json` alone to see available fields

### JQ Filtering

**Inline Processing:**

```bash
gh pr list --json author --jq '.[].author.login'
```

**Complex Queries:**

```bash
gh issue list --json number,title,labels --jq \
  'map(select((.labels | length) > 0))
  | map(.labels = (.labels | map(.name)))
  | .[:3]'
```

### Go Templates

**Basic:**

```bash
gh pr list --json number,title --template '{{range .}}#{{.number}}: {{.title}}{{"\n"}}{{end}}'
```

**With Helpers:**

```bash
gh pr list --json number,title,updatedAt --template \
  '{{range .}}{{tablerow (printf "#%v" .number | autocolor "green") .title (timeago .updatedAt)}}{{end}}'
```

## Strengths of `gh` CLI

### 1. **Excellent Core Functionality**

✅ **Repository Management**

- Create, clone, fork, archive repos
- Manage settings, deploy keys, autolinks
- View, sync, rename operations

✅ **Issue Management**

- Full CRUD operations
- Pin, lock, transfer
- Label management
- Project integration

✅ **PR Basics**

- Create, list, view PRs
- Merge, close, reopen
- Update branches
- Check CI status

✅ **Authentication**

- Seamless login/logout
- Multi-account support (`gh auth switch`)
- Token management
- Git integration (`gh auth setup-git`)

### 2. **Powerful API Access**

✅ **`gh api` Command**

- Direct REST/GraphQL access
- Automatic authentication
- Built-in pagination
- JSON processing (jq, templates)
- Response caching

✅ **JSON Output**

- Many commands support `--json`
- Structured data extraction
- Scriptable workflows

✅ **Formatting Tools**

- jq integration (no installation needed)
- Go templates with helpers
- Color support
- Table rendering

### 3. **Extension Ecosystem**

✅ **Easy to Extend**

- Simple extension model
- No complex API
- Arguments passed through
- Discovery via GitHub topics

✅ **Integration**

- Extensions feel native
- Automatic update checks
- Built-in browsing (`gh extension browse`)

### 4. **Developer Experience**

✅ **Context Awareness**

- Auto-detects current repo
- Infers PR from branch
- Placeholder replacement (`{owner}`, `{repo}`)

✅ **Interactive Mode**

- Prompts for missing info
- Editor integration
- Web fallback (`--web`)

✅ **Helpful Defaults**

- Smart argument parsing
- Multiple ID formats (number, URL, branch)
- Consistent interface across commands

## Limitations of `gh` CLI

### 1. **PR Review Thread Management** ❌

**No Support For:**

- ❌ Listing review threads
- ❌ Viewing thread comments
- ❌ Replying to threads
- ❌ Resolving/unresolving threads
- ❌ Creating review threads
- ❌ Filtering threads (resolved, unresolved, by file)

**Why This Matters:**
Review threads are the core of code review conversations. The lack of thread management forces users to:

- Switch to web UI for detailed reviews
- Cannot handle review feedback from terminal
- Cannot see which comments are resolved
- Cannot reply to specific line comments

**Workaround:**
Must use `gh api` with GraphQL queries to access review threads.

### 2. **Comment Limitations** ❌

**What `gh pr comment` Does:**

- ✅ Adds top-level PR comments
- ✅ Edits last comment
- ✅ Deletes last comment

**What It Doesn't Do:**

- ❌ Add emoji reactions
- ❌ Reply to review threads
- ❌ Create line-specific comments
- ❌ View comment reactions
- ❌ Minimize/hide comments
- ❌ View comment metadata

**Issue Comments Similar:**

- Same capabilities
- Same limitations
- No reaction support
- No threading

### 3. **Emoji Reactions** ❌

**Completely Missing:**

- ❌ No command to add reactions
- ❌ No command to remove reactions
- ❌ No command to list reactions
- ❌ `reactionGroups` in JSON output shows summary only

**Why This Matters:**
Emoji reactions are a quick, low-friction way to:

- Acknowledge comments
- Show agreement/disagreement
- Indicate you've seen something
- Communicate without lengthy responses

**Workaround:**
Must use `gh api` with mutations:

```bash
gh api graphql -f query='mutation {
  addReaction(input: {subjectId: "...", content: THUMBS_UP}) {
    reaction { id }
  }
}'
```

### 4. **Review Management** ❌

**`gh pr review` Limitations:**

- ❌ Cannot add line-specific comments
- ❌ Cannot create review threads
- ❌ Cannot resolve threads
- ❌ No way to dismiss reviews
- ❌ Cannot view review threads

**What It Can Do:**

- ✅ Approve PR (review-level)
- ✅ Request changes (review-level)
- ✅ Add general comment (review-level)

**The Gap:**
A "review" in GitHub includes:

1. Review-level comment (✅ supported)
2. Line-specific comments creating threads (❌ not supported)
3. Thread resolution (❌ not supported)

### 5. **Filtering and Search** ⚠️

**Limited Filtering:**

- `gh pr list` - Basic state, label, assignee filters
- `gh issue list` - Similar basic filters
- No filter for comment activity
- No filter for review status
- No filter for reaction type

**No Thread-Level Operations:**

- Cannot list unresolved threads
- Cannot filter by file path
- Cannot search thread content
- Cannot bulk operations on threads

### 6. **Comment Metadata** ❌

**Cannot Access:**

- Individual comment timestamps
- Comment authors (except in JSON for top-level)
- Comment edit history
- Comment reactions per-comment
- Hidden/minimized status

**JSON Output Limitations:**

- `comments` field is often just a count
- No detailed comment list in most commands
- Must make separate API calls

### 7. **Bulk Operations** ❌

**No Batch Support:**

- Cannot resolve multiple threads at once
- Cannot add reactions to multiple comments
- Cannot hide multiple comments
- Each operation requires separate command

### 8. **Conversation Context** ❌

**Missing Context:**

- Cannot see thread context (surrounding code)
- Cannot see diff when viewing comments
- No inline preview of file content
- Must open PR to see code context

## Comparison: What gh-talk Adds

### Core Value Proposition

**What `gh` Does Well:**

- PR creation and basic management
- Issue creation and basic management
- General commenting
- Repository operations
- Powerful API access

**What `gh-talk` Adds:**

1. **Review Thread Management** ⭐
   - List threads (all, resolved, unresolved)
   - View thread details
   - Reply to threads
   - Resolve/unresolve threads
   - Filter by file, status, author

2. **Emoji Reactions** ⭐
   - Add/remove reactions easily
   - View reaction counts
   - Quick acknowledgments
   - All 8 GitHub reaction types

3. **Comment Hiding** ⭐
   - Minimize off-topic comments
   - Hide spam
   - Mark as outdated
   - Reduce noise

4. **Advanced Filtering**
   - Unresolved threads only
   - By file path
   - By author
   - By date range
   - By reaction type

5. **Bulk Operations**
   - Resolve multiple threads
   - Add reactions to multiple comments
   - Batch comment hiding

6. **Better UX for Reviews**
   - See all threads at once
   - Context-aware display
   - Markdown preview
   - Quick actions

## How gh-talk Leverages gh

### Built on `gh` Foundation

**Uses `gh` Infrastructure:**

1. **Authentication** - `gh auth` credentials
2. **API Access** - `gh api` for GraphQL
3. **Extension Model** - Installed as `gh talk`
4. **JSON Processing** - Can use `gh`'s jq/template features
5. **Context Detection** - Leverage `{owner}/{repo}` placeholders

### Complementary, Not Competing

**gh-talk Should:**

- ✅ Use `gh api` for GraphQL queries
- ✅ Respect `GH_TOKEN` environment variable
- ✅ Follow `gh` command conventions
- ✅ Integrate with `gh`'s ecosystem
- ✅ Feel like a native `gh` command

**gh-talk Should NOT:**

- ❌ Reimplement `gh pr create`
- ❌ Replace `gh issue create`
- ❌ Duplicate repository operations
- ❌ Try to handle authentication directly

### Integration Points

**Use `gh` For:**

```bash
# Get current repo/PR context
gh repo view --json nameWithOwner
gh pr view --json number

# Execute GraphQL queries
gh api graphql -f query='...'

# Check authentication
gh auth status
```

**gh-talk Adds:**

```bash
# List review threads
gh talk list threads --unresolved

# Reply to thread
gh talk reply PRRT_abc123 "Fixed!"

# Add reaction
gh talk react PRRC_xyz789 👍

# Resolve thread
gh talk resolve PRRT_abc123
```

## Technical Insights

### Extension Execution Model

**How Extensions Run:**

1. User runs `gh talk <command>`
2. `gh` looks for executable `gh-talk`
3. `gh` passes all arguments to `gh-talk`
4. `gh-talk` runs independently
5. Output goes directly to user

**Implications:**

- Extensions are standalone binaries
- Full control over UX
- Can use any language (we use Go)
- Access to all `gh` features via shell out
- No IPC overhead

### go-gh Library

**What It Provides:**

```go
import "github.com/cli/go-gh/v2/pkg/api"

// Get authenticated GraphQL client
client, _ := api.DefaultGraphQLClient()

// Execute query
var response struct { ... }
client.Query("repository", query, &response, vars)
```

**Benefits:**

- Handles authentication automatically
- Uses same config as `gh` CLI
- GraphQL and REST clients
- JSON response parsing
- Respects `GH_TOKEN`, `GH_HOST`, etc.

### Command Conventions

**Following `gh` Patterns:**

```bash
# Verb-noun structure
gh pr create    # not: gh create pr
gh talk reply   # not: gh reply talk

# Consistent flags
--help          # Always available
--repo OWNER/REPO  # Override current repo
--json fields   # JSON output
--jq query      # Filter output
--web           # Open in browser

# Argument formats
gh pr view 123              # Number
gh pr view URL              # Full URL
gh pr view branch-name      # Branch name
```

**gh-talk Should:**

- Use similar patterns
- Support thread ID formats
- Provide `--json` output
- Support `--help` thoroughly

## Recommendations for gh-talk

### 1. Lean on `gh api`

**Don't Reimplement:**

```go
// Bad: Manual HTTP client
client := &http.Client{...}
req, _ := http.NewRequest("POST", "https://api.github.com/graphql", ...)

// Good: Use go-gh
client, _ := api.DefaultGraphQLClient()
client.Query("repo", query, &response, vars)
```

### 2. Follow `gh` Conventions

**Command Structure:**

```bash
# Good: Matches gh patterns
gh talk list threads --unresolved
gh talk reply <id> <message>
gh talk resolve <id>

# Bad: Different patterns
gh talk threads:list --unresolved
gh talk <id> reply <message>
gh talk <id> --resolve
```

### 3. Provide JSON Output

**Where It Makes Sense:**

```bash
gh talk list threads --json id,isResolved,path
gh talk show <id> --json comments,reactions

# Can be piped to jq
gh talk list threads --json id,isResolved | jq '.[] | select(.isResolved == false)'
```

### 4. Context Awareness

**Infer from Environment:**

```bash
# When in repo directory
gh talk list threads  # Infers owner/repo from git

# Override when needed
gh talk list threads --repo owner/repo --pr 123
```

### 5. Respect `gh` Authentication

**Use Existing Credentials:**

```go
// Automatic with go-gh
client, err := api.DefaultGraphQLClient()

// Respects:
// - gh auth login credentials
// - GH_TOKEN environment variable
// - GH_ENTERPRISE_TOKEN for GHES
```

### 6. Helpful Error Messages

**Like `gh`:**

```
✗ Error: Thread not found
│ 
│ Could not find thread with ID: PRRT_invalid123
│ 
│ To see available threads, run:
│   gh talk list threads
```

### 7. Interactive Fallbacks

**When Missing Arguments:**

```go
// If user doesn't provide message
if message == "" {
    // Prompt interactively
    message = promptForMessage()
}

// Or suggest --editor
fmt.Println("No message provided. Use --editor to compose in your editor.")
```

## References

### Official Documentation

- **GitHub CLI Manual**: [https://cli.github.com/manual](https://cli.github.com/manual)
- **go-gh Library**: [https://github.com/cli/go-gh](https://github.com/cli/go-gh)
- **Extension Guide**: [https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions)

### Command References

- `gh --help` - Main help
- `gh pr --help` - PR commands
- `gh issue --help` - Issue commands
- `gh api --help` - API access
- `gh help formatting` - JSON/template formatting
- `gh help reference` - Complete command reference

### Extension Resources

- **Extension Topic**: [https://github.com/topics/gh-extension](https://github.com/topics/gh-extension)
- **go-gh Examples**: [https://github.com/cli/go-gh/blob/trunk/example_gh_test.go](https://github.com/cli/go-gh/blob/trunk/example_gh_test.go)

---

**Last Updated**: 2025-11-02  
**gh Version**: Based on gh CLI v2.x  
**Context**: Analysis for gh-talk extension development
