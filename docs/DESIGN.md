# Design Decisions

**Key design choices and their rationale**

## Thread ID System

### Decision: Hybrid Multi-Format Support

**Supports:**
- ✅ Full Node IDs (for scripting/automation)
- ✅ Interactive selection (for exploration)
- ✅ URLs (for convenience)

**Does NOT Support:**
- ❌ Short numeric IDs (1, 2, 3) - Requires persistent caching, complex invalidation

### Rationale

**Why Full IDs:**
- Accurate and unambiguous
- No caching required
- Works in any context
- Perfect for scripts and automation
- Copy-paste from API responses

**Why Interactive:**
- Best UX for exploratory use
- No need to copy/paste IDs
- See context while selecting
- Natural for terminal workflows

**Why URLs:**
- Copy-paste from browser
- Shareable references
- Human-readable context
- Natural for collaboration

**Why NOT Short IDs:**
- Requires persistent cache (complexity)
- Cache invalidation is hard (new threads, deleted threads)
- Context-dependent (which PR?)
- Breaks in new sessions
- False convenience (more complexity than value)

### Implementation Strategy

#### 1. Argument Parsing

```go
// GetThreadID resolves thread ID from various formats
func GetThreadID(arg string) (string, error) {
    // Case 1: Empty argument → Interactive selection
    if arg == "" {
        return selectThreadInteractively()
    }
    
    // Case 2: Full Node ID → Use directly
    if strings.HasPrefix(arg, "PRRT_") {
        return arg, nil
    }
    
    // Case 3: URL → Extract ID
    if strings.Contains(arg, "github.com") || strings.Contains(arg, "http") {
        return extractThreadIDFromURL(arg)
    }
    
    // Case 4: Invalid format
    return "", fmt.Errorf(
        "invalid thread reference: %s\n\n" +
        "Supported formats:\n" +
        "  - Full ID: PRRT_kwDOQN97u85gQeTN\n" +
        "  - URL: https://github.com/owner/repo/pull/123#discussion_r456\n" +
        "  - Empty (interactive selection)\n",
        arg,
    )
}
```

#### 2. Interactive Selection

```go
func selectThreadInteractively() (string, error) {
    // Get current PR context
    repo, _ := repository.Current()
    prNum, _ := getCurrentPRNumber()
    
    // Fetch threads
    client, _ := api.NewClient()
    threads, err := client.ListThreads(repo.Owner, repo.Name, prNum)
    if err != nil {
        return "", err
    }
    
    if len(threads) == 0 {
        return "", fmt.Errorf("no review threads found")
    }
    
    // Build options
    options := make([]string, len(threads))
    for i, t := range threads {
        status := "○"
        if t.IsResolved {
            status = "✓"
        }
        options[i] = fmt.Sprintf(
            "%s %s:%d - %s (%d comments)",
            status,
            t.Path,
            t.Line,
            truncate(t.Comments[0].Body, 50),
            len(t.Comments),
        )
    }
    
    // Prompt user
    p := prompter.New(os.Stdin, os.Stdout, os.Stderr)
    idx, err := p.Select("Select thread:", "", options)
    if err != nil {
        return "", err
    }
    
    return threads[idx].ID, nil
}
```

#### 3. URL Parsing

**GitHub URL Formats:**
```
Review Thread:
https://github.com/owner/repo/pull/123#discussion_r1234567890

Review Comment:
https://github.com/owner/repo/pull/123#discussion_r1234567890

Issue Comment:
https://github.com/owner/repo/issues/456#issuecomment-7890123456
```

**Implementation:**
```go
func extractThreadIDFromURL(urlStr string) (string, error) {
    u, err := url.Parse(urlStr)
    if err != nil {
        return "", fmt.Errorf("invalid URL: %w", err)
    }
    
    // Extract discussion ID from fragment
    // Format: #discussion_r1234567890
    fragment := u.Fragment
    
    if strings.HasPrefix(fragment, "discussion_r") {
        // Review thread/comment
        discussionID := strings.TrimPrefix(fragment, "discussion_r")
        return convertDiscussionIDToNodeID(discussionID)
    }
    
    if strings.HasPrefix(fragment, "issuecomment-") {
        // Issue comment
        commentID := strings.TrimPrefix(fragment, "issuecomment-")
        return convertCommentIDToNodeID(commentID)
    }
    
    return "", fmt.Errorf(
        "could not extract thread ID from URL\n" +
        "URL fragment: %s\n" +
        "Expected format: #discussion_r... or #issuecomment-...",
        fragment,
    )
}
```

**Challenge: Discussion ID ≠ Node ID**

The URL contains a **discussion ID** (integer), but we need the **Node ID** (base64 string).

**Solution: Query API to Convert:**
```go
func convertDiscussionIDToNodeID(discussionID string) (string, error) {
    // Must query GitHub to convert discussion ID to Node ID
    // This is a limitation - URLs don't contain the full Node ID
    
    // Option 1: Query by database ID (if API supports it)
    // Option 2: Fetch all threads, match by discussion ID
    // Option 3: Use a different URL format
    
    // For MVP: Return error, ask user to use full ID
    return "", fmt.Errorf(
        "URL-based thread reference not yet supported\n" +
        "Please use the full thread ID instead.\n\n" +
        "Run 'gh talk list threads' to see thread IDs.",
    )
}
```

**Phase 2 Enhancement:** Implement URL → Node ID conversion

### Usage Examples

```bash
# Interactive (no argument)
gh talk reply
? Select thread:
  ○ test_file.go:7  - Consider using a constant... (2 comments)
  ✓ test_file.go:10 - Good naming for variables... (2 comments)
> ○ test_file.go:14 - This loop could be optimized... (1 comment)
  ✓ test_file.go:18 - Consider extracting this... (1 comment)

Message: Fixed in latest commit

# Full ID (scripting)
gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in latest commit"

# URL (future - Phase 2)
gh talk reply https://github.com/hamishmorgan/gh-talk/pull/1#discussion_r1610088127 "Fixed"
# Note: Requires API lookup to convert discussion ID → Node ID
```

### Implementation Phases

**Phase 1 (MVP):**
- ✅ Full Node IDs
- ✅ Interactive selection
- ❌ URLs (deferred - requires conversion logic)

**Phase 2:**
- ✅ URL support with API conversion
- ✅ Short ID caching (optional enhancement)

## Command Syntax

### Decision: Consistent Verb-Noun Pattern

**Pattern:** `gh talk <verb> [<object>] [<args>] [flags]`

**Following `gh` conventions:**
```bash
gh pr view <number>       # gh pattern
gh talk show <id>         # gh-talk equivalent

gh pr list                # gh pattern
gh talk list threads      # gh-talk equivalent

gh pr comment <pr> <msg>  # gh pattern  
gh talk reply <id> <msg>  # gh-talk equivalent
```

### Command Structure

#### List Commands

```bash
gh talk list threads [flags]
gh talk list comments [flags]
gh talk list reviews [flags]  # Future

Flags:
  --pr <number>          PR number (or infer from git)
  --issue <number>       Issue number
  --repo OWNER/REPO      Repository (or infer from git)
  --unresolved           Only unresolved threads (default for threads)
  --resolved             Only resolved threads
  --all                  All threads (resolved + unresolved)
  --author <username>    Filter by author
  --file <path>          Filter by file path
  --since <date>         Since date
  --format <type>        Output format: table, json, tsv
  --json <fields>        JSON output with specific fields
```

#### Action Commands

```bash
gh talk reply [<thread-id>] [<message>] [flags]
gh talk resolve [<thread-id>] [flags]
gh talk unresolve [<thread-id>] [flags]
gh talk react <comment-id> <emoji> [flags]
gh talk hide <comment-id> [flags]

Arguments:
  <thread-id>    Thread ID, URL, or empty for interactive
  <message>      Message text (or use --editor)
  <comment-id>   Comment ID or URL
  <emoji>        Emoji or name (👍 or THUMBS_UP)

Flags:
  --resolve              Resolve after replying
  --message, -m <text>   Message text
  --editor, -e           Open editor for message
  --reason <type>        Hide reason (spam, off-topic, etc.)
  --remove               Remove reaction instead of add
```

#### Show Command

```bash
gh talk show [<id>] [flags]

Arguments:
  <id>    Thread ID, Comment ID, Issue, or PR (or empty for current)

Flags:
  --type <type>     Force type: thread, comment, issue, pr
  --with-diff       Include code diff context
  --format <type>   Output format: text, json
  --json <fields>   JSON output with fields
```

### Context Inference

**Automatic Detection:**
```bash
# In repo with current branch having PR
gh talk list threads
# → Automatically uses current PR

# Explicit override
gh talk list threads --pr 123
gh talk list threads --repo owner/repo --pr 123
```

**Priority:**
1. Explicit flags (--pr, --issue, --repo)
2. Environment variables (GH_REPO)
3. Git context (current repository + branch)
4. Error if none available

### Argument vs Flag Philosophy

**Arguments (Positional):**
- Primary object: thread ID, comment ID
- Primary action: message text
- Required for command to work

**Flags (Named):**
- Modifiers: --resolve, --editor
- Filters: --unresolved, --author
- Context: --pr, --repo
- Output: --format, --json

**Examples:**
```bash
# ID and message as arguments (most common)
gh talk reply PRRT_abc123 "Fixed!"

# Modifiers as flags
gh talk reply PRRT_abc123 "Fixed!" --resolve

# Interactive (no arguments)
gh talk reply
# Prompts for: thread selection, message

# Message via editor (no message argument)
gh talk reply PRRT_abc123 --editor
```

## Error Messages

### Decision: Helpful, Actionable, Contextual

**Pattern:**
```
<emoji> <short description>

<detailed explanation>

<suggested actions>

<contextual information>
```

**Examples:**

**Thread Not Found:**
```
✗ Thread not found

Could not find thread with ID: PRRT_kwDOQN97u85gQeTN

The thread may have been:
  • Deleted
  • From a different PR
  • Invalid ID format

To see available threads, run:
  gh talk list threads
```

**Permission Denied:**
```
✗ Permission denied

You don't have permission to resolve this thread.

Only the following can resolve threads:
  • PR author
  • Comment author  
  • Repository administrators

Thread: PRRT_kwDOQN97u85gQeTN
PR: hamishmorgan/gh-talk#1
```

**No PR Context:**
```
✗ Not in a PR context

Cannot determine which PR to use.

You can either:
  • Run this command from a branch with an open PR
  • Specify a PR explicitly: --pr <number>
  • Set GH_REPO: export GH_REPO=owner/repo

Current directory: /Users/hamish/src/gh-talk
Git branch: main
```

### Error Message Guidelines

**DO:**
- ✅ Use emoji for quick visual recognition (✓ ✗ ⚠️ 💡)
- ✅ Provide context (what failed, why)
- ✅ Suggest actionable next steps
- ✅ Include relevant IDs/values
- ✅ Format for readability (bullets, sections)

**DON'T:**
- ❌ Just print raw API errors
- ❌ Use technical jargon unnecessarily
- ❌ Leave user stuck (always suggest action)
- ❌ Be verbose (keep it concise)

## Output Formats

### Decision: Terminal-Adaptive with Explicit Override

**Auto-Detection:**
```go
terminal := term.FromEnv()

if terminal.IsTerminalOutput() {
    // Human-readable table
    renderTable(threads)
} else {
    // Machine-readable TSV
    renderTSV(threads)
}
```

**Explicit Formats:**
```bash
--format table    # Force table (even in pipe)
--format json     # JSON output
--format tsv      # Tab-separated values
--json <fields>   # JSON with specific fields (like gh)
```

### Table Format

**For List Threads:**
```
ID                       File:Line       Status      Comments  Preview
-----------------------  --------------  ----------  --------  ----------------------------
PRRT_kwDOQN97u85gQeTN   test_file.go:7  ○ OPEN           2  Consider using a constant...
PRRT_kwDOQN97u85gQecu   test_file.go:14 ○ OPEN           1  This loop could be optimized...
PRRT_kwDOQN97u85gQfgh   test_file.go:10 ✓ RESOLVED       2  Good naming for variables...
```

**With Colors:**
- Green ✓ for resolved
- Default ○ for unresolved
- Gray for file paths
- Bold for preview

**Auto-Truncation:**
- Preview column truncates to fit terminal
- ID column never truncates
- File:Line column never truncates

### JSON Format

**Like `gh` CLI:**
```bash
# List specific fields
gh talk list threads --json id,path,line,isResolved

# Output:
[
  {
    "id": "PRRT_kwDOQN97u85gQeTN",
    "path": "test_file.go",
    "line": 7,
    "isResolved": false
  },
  {
    "id": "PRRT_kwDOQN97u85gQecu",
    "path": "test_file.go",
    "line": 14,
    "isResolved": false
  }
]
```

**Can pipe to jq:**
```bash
gh talk list threads --json id,isResolved | jq '.[] | select(.isResolved == false) | .id'
```

### TSV Format

**Non-TTY Default:**
```bash
gh talk list threads | cat
# Outputs TSV automatically

ID	Path	Line	IsResolved	CommentCount	Preview
PRRT_kwDOQN97u85gQeTN	test_file.go	7	false	2	Consider using a constant...
PRRT_kwDOQN97u85gQecu	test_file.go	14	false	1	This loop could be optimized...
```

**Perfect for:**
- Piping to other commands
- Processing in scripts
- No truncation
- Parseable

## Flag Conventions

### Decision: Follow `gh` Patterns

**Standard Flags (like gh):**
```bash
-R, --repo OWNER/REPO      Repository
-q, --jq <expression>      Filter JSON with jq
-t, --template <tmpl>      Format with Go template
    --json <fields>        JSON output
-w, --web                  Open in browser
-h, --help                 Show help
```

**gh-talk Specific:**
```bash
    --pr <number>          PR number
    --issue <number>       Issue number
    --unresolved           Unresolved only (default for threads)
    --resolved             Resolved only
    --all                  All (resolved + unresolved)
    --author <user>        Filter by author
    --file <path>          Filter by file
    --since <date>         Since date/time
    --format <type>        Output format (table, json, tsv)
-m, --message <text>       Message text
-e, --editor               Open editor
    --resolve              Resolve after action
    --reason <type>        Reason (for hiding)
    --remove               Remove instead of add
```

### Flag Groups

**Filter Flags (for list commands):**
- `--unresolved`, `--resolved`, `--all`
- `--author <user>`
- `--file <path>`
- `--since <date>`

**Context Flags (global):**
- `--repo OWNER/REPO`
- `--pr <number>`
- `--issue <number>`

**Output Flags (global):**
- `--format <type>`
- `--json <fields>`
- `--jq <expression>`

**Action Modifiers:**
- `--resolve` (for reply)
- `--reason <type>` (for hide)
- `--remove` (for react)
- `--editor` (for reply)

## Emoji Handling

### Decision: Accept Both Emoji and Names

**Supported Formats:**
```bash
gh talk react <id> 👍              # Unicode emoji
gh talk react <id> THUMBS_UP       # GraphQL enum name
gh talk react <id> thumbs_up       # Lowercase variant
gh talk react <id> :thumbs_up:     # Slack-style
gh talk react <id> "+1"            # Shorthand
```

**Mapping:**
```go
var emojiMap = map[string]string{
    "👍":         "THUMBS_UP",
    "THUMBS_UP":  "THUMBS_UP",
    "thumbs_up":  "THUMBS_UP",
    ":thumbs_up:": "THUMBS_UP",
    "+1":         "THUMBS_UP",
    "👎":         "THUMBS_DOWN",
    "-1":         "THUMBS_DOWN",
    "😄":         "LAUGH",
    ":laugh:":    "LAUGH",
    "🎉":         "HOORAY",
    ":hooray:":   "HOORAY",
    ":tada:":     "HOORAY",
    "😕":         "CONFUSED",
    ":confused:": "CONFUSED",
    "❤️":         "HEART",
    ":heart:":    "HEART",
    "🚀":         "ROCKET",
    ":rocket:":   "ROCKET",
    "👀":         "EYES",
    ":eyes:":     "EYES",
}

func parseEmoji(input string) (string, error) {
    // Normalize
    input = strings.TrimSpace(strings.ToUpper(input))
    
    // Check map (case-insensitive)
    if graphqlEnum, ok := emojiMap[strings.ToLower(input)]; ok {
        return graphqlEnum, nil
    }
    
    // Try direct uppercase (THUMBS_UP → THUMBS_UP)
    if isValidReactionContent(input) {
        return input, nil
    }
    
    return "", fmt.Errorf(
        "invalid emoji: %s\n\n" +
        "Supported reactions:\n" +
        "  👍 THUMBS_UP     😄 LAUGH      ❤️ HEART\n" +
        "  👎 THUMBS_DOWN   🎉 HOORAY     🚀 ROCKET\n" +
        "  😕 CONFUSED      👀 EYES\n",
        input,
    )
}
```

## Repository Context Detection

### Decision: Follow `gh` Precedence

**Priority Order:**
1. Explicit `--repo` flag
2. `GH_REPO` environment variable  
3. Git remotes in current directory
4. Error if none available

**Implementation:**
```go
func getRepository(repoFlag string) (repository.Repository, error) {
    // 1. Explicit flag
    if repoFlag != "" {
        repo, err := repository.Parse(repoFlag)
        if err != nil {
            return repository.Repository{}, fmt.Errorf(
                "invalid repository format: %s\n" +
                "Expected: OWNER/REPO or HOST/OWNER/REPO",
                repoFlag,
            )
        }
        return repo, nil
    }
    
    // 2 & 3. Environment variable or git remotes
    // (repository.Current() handles both)
    repo, err := repository.Current()
    if err != nil {
        return repository.Repository{}, fmt.Errorf(
            "could not determine repository\n\n" +
            "You can either:\n" +
            "  • Run this command from a git repository\n" +
            "  • Use --repo OWNER/REPO flag\n" +
            "  • Set GH_REPO environment variable\n",
        )
    }
    
    return repo, nil
}
```

## PR/Issue Detection

### Decision: Explicit Context Required for Some Commands

**Commands Needing PR Context:**
- `list threads` - PRs only
- `resolve` - PR threads only
- `reply` - Needs to know if PR thread or issue comment

**Two Approaches:**

**Approach 1: Infer from Current Branch**
```go
func getCurrentPR() (int, error) {
    // Use gh to check current branch
    stdout, _, err := gh.Exec("pr", "view", "--json", "number")
    if err != nil {
        return 0, fmt.Errorf("no PR for current branch")
    }
    
    var result struct{ Number int }
    json.Unmarshal(stdout.Bytes(), &result)
    return result.Number, nil
}
```

**Approach 2: Require Explicit Flag**
```bash
gh talk list threads --pr 123
```

**Decision: Hybrid**
- Try to infer from current branch
- Fall back to requiring flag
- Clear error message if neither works

**Example:**
```bash
# In branch with PR → Works
git checkout my-feature-branch
gh talk list threads

# In main branch → Error with suggestion
git checkout main
gh talk list threads
# Error: No PR for current branch
#   Use --pr <number> to specify a PR

# Explicit works anywhere
gh talk list threads --pr 123
```

## Comment Type Detection

### Decision: Auto-Detect from ID Prefix

**ID Prefix → Type Mapping:**
```go
func detectType(id string) string {
    switch {
    case strings.HasPrefix(id, "PRRT_"):
        return "review_thread"
    case strings.HasPrefix(id, "PRRC_"):
        return "review_comment"
    case strings.HasPrefix(id, "IC_"):
        return "comment"        // Issue or PR top-level
    case strings.HasPrefix(id, "I_"):
        return "issue"
    case strings.HasPrefix(id, "PR_"):
        return "pull_request"
    default:
        return "unknown"
    }
}
```

**Ambiguity: `IC_` Comments**

`IC_` can be:
- Issue comment
- PR top-level comment

**Resolution:**
- Query API to determine parent
- Or accept both work the same (reactions, hiding)

**For commands that care:**
```go
func getCommentContext(commentID string) (Context, error) {
    // Query to determine parent
    var query struct {
        Node struct {
            TypeName string `graphql:"__typename"`
            OnIssueComment struct {
                Issue struct {
                    Number int
                }
            } `graphql:"... on IssueComment"`
        } `graphql:"node(id: $id)"`
    }
    
    variables := map[string]interface{}{
        "id": graphql.ID(commentID),
    }
    
    // Determine if issue or PR from parent
    // ...
}
```

## Reaction Display

### Decision: Show Non-Zero Only (with --all Flag)

**Default Behavior:**
```
Comments:
  PRRC_kwDO...  "Consider using a constant..."  by hamishmorgan
  Reactions: 🚀 1

  PRRC_kwDO...  "Good point! I will refactor..."  by hamishmorgan
  (no reactions)
```

**With --all Flag:**
```bash
gh talk show <id> --reactions all

Comments:
  PRRC_kwDO...  "Consider using a constant..."
  Reactions: 👍 0  👎 0  😄 0  🎉 0  😕 0  ❤️ 0  🚀 1  👀 0
```

**With --verbose:**
```bash
gh talk show <id> --verbose

Comments:
  PRRC_kwDO...  "Consider using a constant..."
  Reactions:
    🚀 hamishmorgan (2025-11-02 21:43:33Z)
```

**In Lists:**
```
ID                       File:Line       Status    Reactions  Preview
-----------------------  --------------  --------  ---------  --------------------
PRRT_kwDOQN97u85gQeTN   test_file.go:7  ○ OPEN    🚀 1       Consider using...
```

## Filtering Strategy

### Decision: Client-Side Filtering

**Why:**
- GitHub API doesn't support server-side thread filtering
- Must fetch all, filter locally
- Cache to avoid repeated fetches

**Implementation:**
```go
func (c *Client) ListThreads(owner, name string, pr int, filter Filter) ([]Thread, error) {
    // 1. Fetch all threads (with caching)
    allThreads, err := c.fetchAllThreads(owner, name, pr)
    if err != nil {
        return nil, err
    }
    
    // 2. Apply filters client-side
    filtered := []Thread{}
    for _, thread := range allThreads {
        if filter.Match(thread) {
            filtered = append(filtered, thread)
        }
    }
    
    return filtered, nil
}

type Filter struct {
    Resolved   *bool    // nil = all, true = resolved, false = unresolved
    Author     string   // Filter by comment author
    File       string   // Filter by file path
    Since      time.Time
}

func (f Filter) Match(t Thread) bool {
    if f.Resolved != nil && t.IsResolved != *f.Resolved {
        return false
    }
    if f.Author != "" && !threadHasAuthor(t, f.Author) {
        return false
    }
    if f.File != "" && t.Path != f.File {
        return false
    }
    if !f.Since.IsZero() && t.CreatedAt.Before(f.Since) {
        return false
    }
    return true
}
```

**Caching:**
```go
// Cache all threads for 5 minutes
// Invalidate on mutations (resolve, reply, etc.)
type ThreadCache struct {
    Key       string    // "owner/repo/pr/123"
    Threads   []Thread
    FetchedAt time.Time
    TTL       time.Duration
}

func (c *ThreadCache) IsValid() bool {
    return time.Since(c.FetchedAt) < c.TTL
}
```

## Bulk Operations

### Decision: Multiple Arguments, Not Stdin

**Pattern:**
```bash
# Multiple IDs as arguments
gh talk resolve PRRT_abc123 PRRT_def456 PRRT_ghi789

# Not: stdin piping (too magical)
# gh talk list threads --json id | gh talk resolve
```

**Why:**
- More explicit
- Clear what's happening
- Safer (can review before executing)
- Follows shell conventions

**With Confirmation:**
```bash
gh talk resolve PRRT_abc123 PRRT_def456 PRRT_ghi789

⚠️  About to resolve 3 threads:
  • test_file.go:7  - Consider using a constant...
  • test_file.go:14 - This loop could be optimized...
  • test_file.go:18 - Consider extracting this...

? Proceed? (y/N)
```

**Skip Confirmation:**
```bash
gh talk resolve PRRT_abc123 PRRT_def456 --yes
```

## Configuration

### Decision: Minimal Config for MVP

**Phase 1: No Config File**
- Use environment variables
- Use flags
- Keep it simple

**Phase 2: Optional Config**
```yaml
# ~/.config/gh-talk/config.yml

defaults:
  format: table
  unresolved_only: true

output:
  color: true
  compact: false

aliases:
  "+1": "THUMBS_UP"
  "ship": "ROCKET"
  "thanks": "HEART"
```

**Not Needed for MVP**

## Interactive Mode

### Decision: Defer to Phase 3

**MVP (Phase 1-2):**
- Focus on command-line arguments
- Interactive selection for missing IDs
- No full TUI

**Future (Phase 3):**
- Full TUI with Bubble Tea
- Keyboard navigation
- Real-time updates
- Visual thread browser

**Rationale:**
- TUI is complex (2-3 weeks of work)
- Command-line covers 80% of use cases
- Can add TUI later without breaking CLI
- Focus on core functionality first

## Summary of Key Decisions

| Decision | Choice | Phase |
|----------|--------|-------|
| Thread ID Input | Full ID, Interactive, URL | 1 (URL in 2) |
| Short IDs | Not supported | Deferred |
| Command Pattern | Verb-noun (like gh) | 1 |
| Context Detection | Auto-infer with flag override | 1 |
| Error Messages | Helpful with suggestions | 1 |
| Output Format | Terminal-adaptive + explicit | 1 |
| Filtering | Client-side | 1 |
| Caching | 5 min TTL, invalidate on write | 1 |
| Bulk Operations | Multiple args + confirmation | 1 |
| Config File | Not in MVP | 3 |
| Full TUI | Deferred | 3 |

## Implementation Implications

### What This Enables

**Clear Implementation Path:**
1. ✅ Thread ID parsing logic defined
2. ✅ Command structure finalized
3. ✅ Flag conventions established
4. ✅ Error handling patterns set
5. ✅ Output formats specified

**What to Build First:**
1. ID parsing and validation
2. Repository context detection
3. PR/Issue number inference
4. GraphQL client wrapper
5. Basic list command
6. Interactive selection
7. Reply command
8. Resolve command

**What Can Wait:**
- URL → Node ID conversion (Phase 2)
- Config file (Phase 3)
- Full TUI (Phase 3)
- Advanced filtering (Phase 2)

## Validation

### These Decisions Support:

✅ **Scripting:** Full IDs work in automated scripts  
✅ **Exploration:** Interactive selection for ad-hoc use  
✅ **Convenience:** URL support (Phase 2) for browser workflow  
✅ **Consistency:** Follows `gh` patterns  
✅ **Simplicity:** No complex caching for MVP  
✅ **Clarity:** Explicit is better than implicit  
✅ **Scalability:** Can add features without breaking changes  

---

**Last Updated**: 2025-11-02  
**Status**: Finalized for MVP implementation  
**Review**: Ready to proceed with implementation

