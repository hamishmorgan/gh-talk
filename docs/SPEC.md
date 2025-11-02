# gh-talk Specification

**GitHub CLI Extension for Comprehensive PR & Issue Conversation Management**

## Overview

`gh-talk` is a GitHub CLI extension that provides a complete workflow for managing conversations on Pull Requests and Issues. It handles the full lifecycle of GitHub discussions: viewing, responding (text + emoji), triaging (hide/resolve), and filtering.

## Core Philosophy

- **Short commands** - Minimize typing for common operations
- **Comprehensive scope** - Handle all conversation interactions in one tool
- **Terminal-native** - Never leave the terminal for conversation management
- ***AI Native** - Usable as a skill by AI agents
- **Workflow-optimized** - Support the natural flow of code review
- **Feels like gh-cli*** - Interface is as similar to `gh` usage as possible.

## Use Cases

### Primary Use Cases

1. **Reply to review comments** - Respond to code review feedback without opening browser
2. **Emoji reactions** - Quick acknowledgments and emotional responses
3. **Thread management** - Hide noise, resolve addressed comments
4. **Filter conversations** - Find what needs attention
5. **Bulk operations** - Resolve multiple threads, dismiss reviews

### Target Users

- AI Agents so they can interact fully in the issue/PR review process
- Developers who live in the terminal
- Teams with high PR volume
- Anyone tired of browser context-switching for simple reactions

## Command Structure

```bash
gh talk <command> [subcommand] [flags] [args]
```

## Commands

### 1. List / View Commands

#### `gh talk list`

List threads and comments with filtering.

**Subcommands:**

```bash
gh talk list threads              # List all review threads
gh talk list threads --unresolved # Only unresolved threads
gh talk list threads --resolved   # Only resolved threads
gh talk list comments             # List all comments
gh talk list reviews              # List all reviews
```

**Flags:**

- `--unresolved` - Only show unresolved threads
- `--resolved` - Only show resolved threads  
- `--unhidden` - Exclude hidden/minimized comments
- `--reactions <emoji>` - Filter by reactions (e.g., `--reactions 👍,🚀`)
- `--author <username>` - Filter by comment author
- `--since <date>` - Only show comments since date
- `--file <path>` - Only show comments on specific file
- `--format <type>` - Output format: table (default), json, markdown

**Examples:**

```bash
# Show unresolved threads needing attention
gh talk list threads --unresolved

# Show comments with thumbs up reactions
gh talk list comments --reactions 👍

# Show all comments from specific user
gh talk list comments --author reviewer-name

# Show threads on specific file
gh talk list threads --file src/main.go --unresolved
```

**Output format:**

```
Thread ID    File                  Line  Author     Status      Reactions  Preview
-----------  -------------------  ----  ---------  ----------  ---------  --------
PRRT_abc123  src/api.go            42   reviewer1  unresolved  👍 x2      Consider using...
PRRT_def456  src/db.go             89   reviewer2  unresolved  🎉        Great refactor!
```

### 2. Reply Commands

#### `gh talk reply <thread-id> <message>`

Reply to a review thread.

**Arguments:**

- `<thread-id>` - Thread ID (from list command)
- `<message>` - Reply message (use quotes for multi-word)

**Flags:**

- `--resolve` - Resolve thread after replying
- `--react <emoji>` - Add reaction after replying
- `--editor` - Open editor for message composition

**Examples:**

```bash
# Simple reply
gh talk reply PRRT_abc123 "Fixed in commit abc123"

# Reply and resolve
gh talk reply PRRT_abc123 "Addressed by refactoring the function" --resolve

# Reply, react, and resolve
gh talk reply PRRT_abc123 "Good catch, fixed!" --react 👍 --resolve

# Compose in editor
gh talk reply PRRT_abc123 --editor
```

**Special handling:**

- Escape special characters automatically
- Support multi-line messages via editor
- Validate thread exists before posting

### 3. Reaction Commands

#### `gh talk react <comment-id> <emoji>`

Add emoji reaction to a comment.

**Arguments:**

- `<comment-id>` - Comment ID (from list command)
- `<emoji>` - Emoji to react with

**Supported reactions:**

- `👍` `THUMBS_UP` - Agreement, approval
- `👎` `THUMBS_DOWN` - Disagreement  
- `😄` `LAUGH` - Humor, appreciation
- `🎉` `HOORAY` - Celebration
- `😕` `CONFUSED` - Need clarification
- `❤️` `HEART` - Appreciation, love
- `🚀` `ROCKET` - Excitement, ready to ship
- `👀` `EYES` - Acknowledged, watching

**Flags:**

- `--remove` - Remove reaction instead of adding

**Examples:**

```bash
# Add thumbs up
gh talk react PRRC_xyz789 👍

# Add rocket emoji (text shortcode)
gh talk react PRRC_xyz789 🚀

# Remove reaction
gh talk react PRRC_xyz789 👍 --remove
```

**Bulk reactions:**

```bash
# React to multiple comments
gh talk react PRRC_abc123,PRRC_def456 👍
```

### 4. Resolution Commands

#### `gh talk resolve <thread-id>`

Mark a review thread as resolved.

**Arguments:**

- `<thread-id>` - Thread ID to resolve

**Flags:**

- `--message <msg>` - Add reply before resolving
- `--unresolve` - Unresolve instead of resolve

**Examples:**

```bash
# Resolve thread
gh talk resolve PRRT_abc123

# Resolve with explanation
gh talk resolve PRRT_abc123 --message "Fixed in commit abc123"

# Unresolve thread (reopen discussion)
gh talk resolve PRRT_abc123 --unresolve
```

**Bulk resolution:**

```bash
# Resolve multiple threads
gh talk resolve PRRT_abc123,PRRT_def456,PRRT_ghi789

# Resolve all resolved threads matching filter
gh talk list threads --file src/api.go | gh talk resolve --batch
```

### 5. Hide Commands

#### `gh talk hide <comment-id>`

Minimize/hide a comment.

**Arguments:**

- `<comment-id>` - Comment ID to hide

**Flags:**

- `--reason <type>` - Hide reason: off-topic, spam, outdated, resolved
- `--unhide` - Unhide instead of hide

**Examples:**

```bash
# Hide comment as off-topic
gh talk hide PRRC_xyz789 --reason off-topic

# Hide as outdated
gh talk hide PRRC_xyz789 --reason outdated

# Unhide comment
gh talk hide PRRC_xyz789 --unhide
```

### 6. Review Management Commands

#### `gh talk dismiss <review-id>`

Dismiss a review (when all threads resolved).

**Arguments:**

- `<review-id>` - Review ID to dismiss

**Flags:**

- `--message <msg>` - Dismissal message (required)
- `--auto` - Auto-dismiss reviews with all threads resolved

**Examples:**

```bash
# Dismiss specific review
gh talk dismiss PRRE_abc123 --message "All comments addressed"

# Auto-dismiss all reviews with resolved threads
gh talk dismiss --auto --message "All feedback incorporated"
```

**Safety checks:**

- Verify all threads are resolved before dismissing
- Prevent accidental dismissal of active reviews
- Require explicit message explaining dismissal

### 7. Context Commands

#### `gh talk show <id>`

Show detailed information about a thread, comment, or review.

**Arguments:**

- `<id>` - Thread, comment, or review ID

**Flags:**

- `--with-diff` - Include code diff context
- `--format <type>` - Output format: text (default), json, markdown

**Examples:**

```bash
# Show thread details
gh talk show PRRT_abc123

# Show comment with code diff
gh talk show PRRC_xyz789 --with-diff

# Export as JSON
gh talk show PRRT_abc123 --format json
```

**Output includes:**

- Thread/comment metadata (author, timestamp)
- Full conversation history
- Reactions summary
- Resolution status
- Code context (line numbers, file path)

### 8. Interactive Mode

#### `gh talk interactive`

Launch interactive TUI for conversation management.

**Features:**

- Browse threads with arrow keys
- Preview conversations
- Quick reply/react/resolve actions
- Keyboard shortcuts for all operations
- Real-time updates

**Keyboard shortcuts:**

- `r` - Reply to selected thread
- `e` - React with emoji picker
- `x` - Resolve/unresolve thread
- `h` - Hide/unhide comment
- `/` - Filter/search
- `q` - Quit

## Configuration

### Config File

Location: `~/.config/gh-talk/config.yml`

```yaml
# Default settings
defaults:
  format: table
  reactions: true
  auto_resolve: false

# Output preferences
output:
  color: true
  unicode: true
  compact: false

# Filters
filters:
  exclude_resolved: true
  exclude_dismissed: true
  exclude_bots: true

# Aliases
aliases:
  # Custom reaction shortcuts
  +1: 👍
  ship: 🚀
  thanks: ❤️
```

### Environment Variables

- `GH_TALK_FORMAT` - Default output format
- `GH_TALK_REACTIONS` - Enable/disable reaction display
- `GH_TALK_EDITOR` - Editor for message composition

## Data Sources

### GraphQL API Queries

All data fetched via GitHub GraphQL API using the existing workflow patterns from `.cursor/rules/graphql-comments-workflow.mdc`.

**Core queries:**

1. **Review threads** - `pullRequest.reviewThreads`
2. **Comments** - Thread comments with reactions
3. **Reviews** - PR reviews with state
4. **Reactions** - Emoji reactions on comments

### Caching Strategy

- Cache thread/comment data for 5 minutes
- Invalidate on write operations
- Background refresh in interactive mode

## Error Handling

### Graceful Degradation

- Continue on partial failures (e.g., one comment fails to post)
- Clear error messages with suggested fixes
- Retry logic for transient failures

### Common Errors

1. **Permission denied** - Not authorized to reply/resolve
2. **Not found** - Thread/comment deleted or invalid ID
3. **Rate limit** - GraphQL API throttling
4. **Network error** - Connection issues

**Error message format:**

```
✗ Error: Permission denied
│ You don't have permission to resolve this thread
│ 
│ Possible solutions:
│ • Request write access to the repository
│ • Ask the thread author or PR owner to resolve
│ 
│ Thread ID: PRRT_abc123
│ PR: #123 in owner/repo
```

## Performance Goals

- **List command**: < 2s for 100 threads
- **Reply command**: < 1s response time
- **Interactive mode**: < 500ms UI updates
- **Bulk operations**: Process 50 items in < 10s

## Testing Strategy

### Unit Tests

- GraphQL query construction
- Response parsing and error handling
- Filter logic
- Emoji handling and escaping

### Integration Tests

- End-to-end command execution
- API interaction (mocked)
- Error scenarios

### Manual Testing

- Real PR workflows
- Cross-platform compatibility (macOS, Linux)
- Different repository types (public, private, org)

## Future Enhancements

### Phase 2 Features

- **Suggested replies** - AI-generated response templates
- **Batch editing** - Bulk reply to multiple threads
- **Templates** - Saved reply templates
- **Notifications** - Desktop notifications for new comments
- **Diff preview** - Show code changes in context

### Phase 3 Features

- **Issue support** - Extend to issue conversations
- **Discussion support** - GitHub Discussions integration
- **Search** - Full-text search across conversations
- **Analytics** - Conversation metrics and insights

## Dependencies

### Required

- `gh` CLI v2.0+
- GitHub GraphQL API access
- `jq` for JSON processing (Go implementation: no external dependency)

### Optional

- `fzf` for fuzzy finding in interactive mode
- Terminal with Unicode support for emoji rendering

## Distribution

### Installation

```bash
# Install from GitHub
gh extension install hamishmorgan/gh-talk

# Install from source
cd ~/src/gh-talk
gh extension install .
```

### Updates

```bash
# Update to latest version
gh extension upgrade gh-talk

# Uninstall
gh extension remove gh-talk
```

### Packaging

- Single binary distribution
- Cross-compile for macOS (Intel/ARM), Linux (x64/ARM)
- Automated releases via GitHub Actions

## Success Metrics

### Adoption

- 100+ installs in first month
- 5+ daily active users
- Positive feedback from team

### Efficiency

- Reduce browser context switches by 80%
- Average 30s saved per comment interaction
- Support 10+ PRs per day workflow

## Technical Architecture

### Project Structure

```
gh-talk/
├── cmd/
│   ├── list.go       # List commands
│   ├── reply.go      # Reply commands
│   ├── react.go      # Reaction commands
│   ├── resolve.go    # Resolution commands
│   ├── hide.go       # Hide commands
│   ├── dismiss.go    # Dismiss commands
│   └── show.go       # Show commands
├── pkg/
│   ├── api/          # GraphQL API client
│   ├── filter/       # Filtering logic
│   ├── format/       # Output formatting
│   └── config/       # Configuration management
├── internal/
│   ├── tui/          # Terminal UI (interactive mode)
│   └── cache/        # Caching layer
├── main.go           # Entry point
├── go.mod
├── go.sum
├── SPEC.md           # This file
└── README.md         # User documentation
```

### Key Libraries

- `github.com/cli/go-gh` - Official GitHub CLI library
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework (interactive mode)
- `github.com/charmbracelet/lipgloss` - Terminal styling

## Documentation

### User Documentation

- README.md - Installation, quick start, common workflows
- Man page - Detailed command reference
- Examples directory - Real-world usage scenarios

### Developer Documentation

- CONTRIBUTING.md - How to contribute
- ARCHITECTURE.md - Technical design decisions
- API.md - GraphQL query reference

## License

MIT License (following gh CLI extension conventions)

---

**Version**: 0.1.0 (Specification)  
**Last Updated**: 2025-11-02  
**Status**: Draft - Pre-implementation

