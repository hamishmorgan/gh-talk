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
- `--format <type>` - Output format: table (default), json, tsv
- `--json <fields>` - JSON output with specific fields (like gh CLI)
- `--pr <number>` - PR number (or infer from current branch)
- `--issue <number>` - Issue number
- `-R, --repo OWNER/REPO` - Repository (or infer from git)

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
ID                         File:Line          Status        Reactions  Preview
-------------------------  -----------------  ------------  ---------  ---------------------
PRRT_kwDOQN97u85gQeTN     src/api.go:42      ○ OPEN        👍 2       Consider using...
PRRT_kwDOQN97u85gQecu     src/db.go:89       ✓ RESOLVED    🎉 1       Great refactor!
```

**Note:** Real thread IDs are 25-30 characters (base64-encoded Node IDs).
See [REAL-DATA.md](REAL-DATA.md) for actual ID formats from GitHub API.

### 2. Reply Commands

#### `gh talk reply [<thread-id>] [<message>]`

Reply to a review thread.

**Arguments (both optional for interactive mode):**

- `<thread-id>` - Thread ID (PRRT_...), URL, or omit for interactive selection
- `<message>` - Reply message (use quotes for multi-word), or use --editor

**Flags:**

- `--resolve` - Resolve thread after replying
- `-e, --editor` - Open editor for message composition
- `-m, --message <text>` - Message text (alternative to positional argument)

**Examples:**

```bash
# Interactive mode (prompts for thread and message)
gh talk reply

# With full Node ID (real format from GitHub)
gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in commit abc123"

# Reply and resolve
gh talk reply PRRT_kwDOQN97u85gQeTN "Addressed by refactoring" --resolve

# Compose in editor
gh talk reply PRRT_kwDOQN97u85gQeTN --editor

# URL support (Phase 2)
gh talk reply https://github.com/owner/repo/pull/123#discussion_r456 "Fixed"
```

**Thread ID Formats Supported:**

1. **Full Node ID** - `PRRT_kwDOQN97u85gQeTN` (for scripting)
2. **Interactive** - Empty argument prompts selection
3. **URL** - `https://github.com/.../pull/123#discussion_r456` (Phase 2)

**Note:** Short numeric IDs (1, 2, 3) are NOT supported due to caching complexity.
See [DESIGN.md](DESIGN.md) for thread ID system rationale.

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
# React to multiple comments (space-separated)
gh talk react PRRC_kwDOQN97u86UHqK7 PRRC_kwDOQN97u86UHqOo 👍
```

**Emoji Formats Supported:**
- Unicode emoji: `👍`, `🎉`, `❤️`
- GraphQL names: `THUMBS_UP`, `HOORAY`, `HEART`
- Lowercase: `thumbs_up`, `hooray`, `heart`
- Slack-style: `:thumbs_up:`, `:tada:`
- Shorthand: `+1`, `-1`

See [DESIGN.md](DESIGN.md#emoji-handling) for complete mapping.

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
# Resolve multiple threads (space-separated)
gh talk resolve PRRT_kwDO...123 PRRT_kwDO...456 PRRT_kwDO...789

# With confirmation (default for multiple)
? Resolve 3 threads? (y/N)

# Skip confirmation
gh talk resolve PRRT_kwDO...123 PRRT_kwDO...456 --yes
```

**Interactive Mode:**
```bash
# No arguments prompts for selection
gh talk resolve
? Select threads to resolve:
  [x] test_file.go:7  - Consider using a constant...
  [ ] test_file.go:14 - This loop could be optimized...
  [x] test_file.go:18 - Consider extracting this...
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

**GitHub CLI (automatically used via go-gh):**
- `GH_TOKEN` - GitHub authentication token
- `GH_HOST` - GitHub host (default: github.com)
- `GH_REPO` - Repository context (OWNER/REPO)
- `GH_FORCE_TTY` - Force terminal mode
- `GH_DEBUG` - Enable debug logging

**gh-talk Specific:**
- `GH_TALK_CONFIG` - Config file location (default: ~/.config/gh-talk/config.yml)
- `GH_TALK_CACHE_DIR` - Cache directory (default: ~/.cache/gh-talk)
- `GH_TALK_CACHE_TTL` - Cache TTL in minutes (default: 5)
- `GH_TALK_FORMAT` - Default output format (table, json, tsv)
- `GH_TALK_EDITOR` - Editor for message composition

**Terminal:**
- `NO_COLOR` - Disable colors
- `CLICOLOR` - Color support (0 or 1)
- `EDITOR` - Default text editor

See [ENVIRONMENT.md](ENVIRONMENT.md) for complete reference.

## Data Sources

### GraphQL API Queries

All data fetched via GitHub GraphQL API v4.

**Core queries:**

1. **Review threads** - `pullRequest.reviewThreads` (PR only)
2. **Comments** - Thread comments with reactions
3. **Issue comments** - `issue.comments` (separate from review threads)
4. **Reactions** - `reactionGroups` (always includes all 8 types)

**Real Data Structures:**
- Thread IDs: `PRRT_kwDOQN97u85gQeTN` (25-30 chars)
- Comment IDs: `PRRC_kwDOQN97u86UHqK7` or `IC_kwDOQN97u87PVA8l`
- Issue IDs: `I_kwDOQN97u87VYpUq`
- Review IDs: `PRR_kwDOQN97u87LMeCy`

See [API.md](API.md) for GraphQL schema and [REAL-DATA.md](REAL-DATA.md) for actual response structures from live testing.

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

### Test Types

**Unit Tests:**
- GraphQL query construction
- Response parsing and error handling
- Filter logic
- Emoji handling and mapping
- ID validation and parsing
- **Target:** 90%+ coverage for API package

**Integration Tests:**
- Full command execution with mocked API
- Flag parsing and validation
- Output formatting (table, JSON, TSV)
- Error message generation
- **Target:** All commands tested

**Contract Tests:**
- GraphQL query structs match GitHub schema
- Response parsing with real fixtures
- Using testdata/ real API responses
- **Target:** All queries/mutations validated

**E2E Tests (Optional):**
- Real API calls (expensive, rate-limited)
- Manual testing on test PR #1 and Issue #2
- Cross-platform compatibility (Linux, macOS, Windows)
- **Target:** Smoke tests for main workflows

**Test Fixtures:**
- `testdata/pr_full_response.json` - Complete PR with threads
- `testdata/issue_full_response.json` - Complete issue with comments
- `testdata/pr_with_resolved_threads.json` - Mixed resolution states

**Overall Target:** 80%+ code coverage

See [ENGINEERING.md](ENGINEERING.md) for complete testing strategy and CI/CD setup.

## Future Enhancements

### Phase 2 Features

- **Suggested replies** - AI-generated response templates
- **Batch editing** - Bulk reply to multiple threads
- **Templates** - Saved reply templates
- **Notifications** - Desktop notifications for new comments
- **Diff preview** - Show code changes in context

### Phase 3 Features

- **Full TUI mode** - Bubble Tea interactive interface (like gh-dash)
- **Discussion support** - GitHub Discussions integration
- **Search** - Full-text search across conversations
- **Analytics** - Conversation metrics and insights
- **URL → Node ID conversion** - Direct URL support (requires API lookup)

**Note:** Issue support is included in Phase 1 (MVP). Issues and PRs share the same comment/reaction model, with issues being simpler (no review threads).

See [WORKFLOWS.md](WORKFLOWS.md) for detailed usage patterns and [REAL-DATA.md](REAL-DATA.md) for Issue vs PR differences.

## Dependencies

### Required

- `gh` CLI v2.0+
- Go 1.21+ (for building from source)
- GitHub GraphQL API access

### Go Dependencies

- `github.com/cli/go-gh/v2` v2.12.2+ - GitHub CLI library
- `github.com/spf13/cobra` v1.8+ - CLI framework
- `github.com/charmbracelet/bubbletea` v0.25+ - TUI framework (Phase 3)
- `github.com/charmbracelet/lipgloss` v0.9+ - Terminal styling (Phase 3)

**Note:** JSON processing uses Go stdlib `encoding/json` (no external jq dependency)

### Optional

- `fzf` - Enhanced fuzzy finding (if available, used for better interactive selection)
- Terminal with Unicode support for emoji rendering
- Terminal with true color support for best visual experience

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

**Actual Structure (refined from research):**

```
gh-talk/
├── main.go                    # Entry point (minimal)
├── internal/                  # All implementation (private)
│   ├── api/                  # GitHub GraphQL API client wrapper
│   │   ├── client.go         # Client wrapper around go-gh
│   │   ├── threads.go        # Thread operations
│   │   ├── comments.go       # Comment operations
│   │   ├── reactions.go      # Reaction operations
│   │   ├── types.go          # GraphQL type definitions
│   │   └── mock.go           # Test mocks
│   ├── commands/             # Cobra command implementations
│   │   ├── root.go           # Root command setup
│   │   ├── list.go           # List commands
│   │   ├── reply.go          # Reply command
│   │   ├── resolve.go        # Resolve/unresolve commands
│   │   ├── react.go          # React command
│   │   ├── hide.go           # Hide/unhide commands
│   │   └── show.go           # Show command
│   ├── format/               # Output formatting
│   │   ├── table.go          # Table output (terminal)
│   │   ├── json.go           # JSON output
│   │   └── tsv.go            # TSV output (non-TTY)
│   ├── filter/               # Filtering logic
│   │   └── threads.go        # Client-side filtering
│   ├── config/               # Configuration management
│   │   ├── config.go         # Config file handling
│   │   └── env.go            # Environment variables
│   ├── cache/                # Caching layer
│   │   └── cache.go          # Thread data caching (5min TTL)
│   └── tui/                  # Terminal UI (Phase 3)
│       └── tui.go            # Bubble Tea implementation
├── testdata/                  # Test fixtures (real API responses)
│   ├── README.md
│   ├── pr_full_response.json
│   ├── issue_full_response.json
│   └── pr_with_resolved_threads.json
├── docs/                      # Documentation
│   ├── API.md                # GitHub API reference
│   ├── REAL-DATA.md          # Real API responses
│   ├── GO-GH.md              # go-gh library guide
│   ├── COBRA.md              # Cobra implementation guide
│   ├── DESIGN.md             # Design decisions
│   ├── ENGINEERING.md        # Testing & CI/CD
│   └── ...                   # Additional docs
├── .github/
│   └── workflows/
│       ├── test.yml          # Test, lint, coverage
│       ├── build.yml         # Multi-platform builds
│       └── release.yml       # Automated releases
├── .golangci.yml             # Linter configuration
├── Makefile                   # Development tasks
├── go.mod
├── go.sum
└── README.md                  # User documentation
```

**Key Architectural Decisions:**
- **No cmd/** for binary location (main.go in root per gh extension convention)
- **No pkg/** (code is not meant to be imported externally)
- **internal/** for all implementation (enforces encapsulation)
- **testdata/** for real API response fixtures

See [STRUCTURE.md](STRUCTURE.md) for architecture rationale and [EXTENSION-PATTERNS.md](EXTENSION-PATTERNS.md) for validation from successful extensions.

### Key Libraries

**Core:**
- `github.com/cli/go-gh/v2` v2.12.2 - GitHub CLI library (GraphQL/REST clients, auth, terminal)
- `github.com/spf13/cobra` v1.8 - CLI framework (validated by gh-copilot usage)

**Future (Phase 3):**
- `github.com/charmbracelet/bubbletea` v0.25 - TUI framework (interactive mode)
- `github.com/charmbracelet/lipgloss` v0.9 - Terminal styling

**Testing:**
- Standard library `testing` package
- Real API fixtures in `testdata/`
- Mock patterns for API client

**Note:** Cobra choice validated - GitHub's own gh-copilot extension uses Cobra.

## Documentation

### User Documentation

- `README.md` - Installation, quick start, feature overview
- `docs/ENVIRONMENT.md` - Environment variables reference
- `docs/WORKFLOWS.md` - Real-world usage patterns and examples

### Developer Documentation

**Architecture & Design:**
- `docs/SPEC.md` - This specification (complete feature set)
- `docs/DESIGN.md` - Key design decisions and rationale
- `docs/STRUCTURE.md` - Project structure and organization

**APIs & Implementation:**
- `docs/API.md` - GitHub API capabilities and reference
- `docs/REAL-DATA.md` - Real API responses from live testing (1,885 lines!)
- `docs/GO-GH.md` - go-gh library guide and patterns
- `docs/COBRA.md` - Cobra implementation guide

**Analysis & Validation:**
- `docs/GH-CLI.md` - GitHub CLI analysis (what gh does/doesn't do)
- `docs/CLI-FRAMEWORK.md` - Framework choice analysis
- `docs/EXTENSION-PATTERNS.md` - Successful extension patterns
- `docs/WORKFLOWS.md` - User workflows and personas

**Quality & Testing:**
- `docs/ENGINEERING.md` - Testing strategy, CI/CD, quality practices
- `AGENTS.md` - AI agent development guidelines
- `testdata/README.md` - Test fixture documentation

**Total:** 11,882 lines of comprehensive documentation

## License

MIT License (following gh CLI extension conventions)

## Implementation Status

### Research Phase: ✅ Complete

**Completed:**
- ✅ Comprehensive API research and live testing
- ✅ Real data structure documentation (PR #1, Issue #2)
- ✅ Design decisions finalized
- ✅ Framework choice validated (Cobra used by gh-copilot!)
- ✅ Extension pattern analysis (5 successful extensions)
- ✅ Complete engineering infrastructure (CI/CD, testing, linting)
- ✅ 11,882 lines of documentation

**Ready for Implementation:**
- ✅ All critical decisions made
- ✅ Thread ID system designed (Full IDs + Interactive)
- ✅ Command syntax finalized
- ✅ Testing strategy defined
- ✅ CI/CD pipelines created
- ✅ Quality gates established

### Development Phase: Next

**Phase 1: MVP (Weeks 1-3)**
- Add Cobra dependency
- Implement core commands (list, reply, resolve, react)
- Basic filtering and formatting
- Unit and integration tests
- 80%+ code coverage

**Phase 2: Enhancement (Weeks 4-6)**
- Issue support (comments, reactions)
- Hide/unhide commands
- Advanced filtering
- Bulk operations
- Shell completion

**Phase 3: Polish (Weeks 7+)**
- Interactive TUI mode (Bubble Tea)
- Configuration file support
- URL → Node ID conversion
- Advanced features

See [ENGINEERING.md](ENGINEERING.md) for detailed implementation roadmap.

---

**Version**: 0.1.0 (Specification)  
**Last Updated**: 2025-11-02  
**Status**: Specification Complete - Ready for Implementation  
**Research:** Complete with live API testing  
**Documentation:** 11,882 lines across 15 files

