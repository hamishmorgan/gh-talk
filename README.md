# gh-talk

**GitHub CLI Extension for Comprehensive PR & Issue Conversation Management**

A terminal-native tool for managing all aspects of GitHub PR and Issue conversations: replies, reactions, thread resolution, and filtering.

## Status

🚧 **Early Development** - See [SPEC.md](docs/SPEC.md) for full specification

## Quick Start

```bash
# Install the extension
gh extension install hamishmorgan/gh-talk

# List unresolved review threads
gh talk list threads --unresolved

# Reply to a thread
gh talk reply PRRT_abc123 "Fixed in commit abc123"

# Add emoji reaction
gh talk react PRRC_xyz789 👍

# Resolve a thread
gh talk resolve PRRT_abc123
```

## Features

- 💬 **Reply to review threads** - Never leave the terminal
- 😄 **Emoji reactions** - Quick acknowledgments and responses
- ✅ **Resolve/unresolve threads** - Manage conversation state
- 🙈 **Hide comments** - Minimize noise
- 📋 **Filter conversations** - Find what needs attention
- 🚀 **Bulk operations** - Handle multiple threads efficiently
- 🎨 **Interactive mode** - TUI for visual workflow (planned)

## Documentation

- [SPEC.md](docs/SPEC.md) - Complete specification and technical design
- [API.md](docs/API.md) - GitHub API capabilities and reference
- [REAL-DATA.md](docs/REAL-DATA.md) - Real API responses and data structures
- [GO-GH.md](docs/GO-GH.md) - go-gh library guide and patterns
- [DESIGN.md](docs/DESIGN.md) - Key design decisions and rationale
- [CLI-FRAMEWORK.md](docs/CLI-FRAMEWORK.md) - CLI framework choice (Cobra)
- [COBRA.md](docs/COBRA.md) - Cobra implementation guide
- [EXTENSION-PATTERNS.md](docs/EXTENSION-PATTERNS.md) - Successful extension analysis
- [ENGINEERING.md](docs/ENGINEERING.md) - Testing, CI/CD, and quality practices
- [ENVIRONMENT.md](docs/ENVIRONMENT.md) - Environment variables reference
- [GH-CLI.md](docs/GH-CLI.md) - GitHub CLI analysis and integration
- [WORKFLOWS.md](docs/WORKFLOWS.md) - Real-world usage patterns and workflows
- [STRUCTURE.md](docs/STRUCTURE.md) - Project structure and design decisions
- [GitHub CLI Extension Docs](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions)

## For AI Agents

- [AGENTS.md](AGENTS.md) - Development guidelines and standards for AI agents working on this project

## Project Status

- [CRITICAL-REVIEW.md](docs/CRITICAL-REVIEW.md) - Honest assessment of problems and risks
- [ACTION-PLAN.md](docs/ACTION-PLAN.md) - Immediate actions to start implementation

**Current Status:** ✅ **MVP Complete and Enhanced!**

- 9 working commands
- Real-world tested and iterated
- 5 user feedback issues implemented
- All pushed to GitHub

## Usage

### List Review Threads

```bash
# List unresolved threads (default)
gh talk list threads

# List all threads
gh talk list threads --all

# List resolved threads
gh talk list threads --resolved

# Filter by file
gh talk list threads --file src/api.go

# Specific PR
gh talk list threads --pr 123
```

### Reply to Threads

```bash
# Interactive mode (prompts for selection)
gh talk reply

# With thread ID
gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in latest commit"

# Reply and resolve
gh talk reply PRRT_kwDOQN97u85gQeTN "Done!" --resolve
```

### Resolve Threads

```bash
# Interactive selection
gh talk resolve

# Single thread
gh talk resolve PRRT_kwDOQN97u85gQeTN

# Multiple threads
gh talk resolve PRRT_abc PRRT_def PRRT_ghi

# With message
gh talk resolve PRRT_abc --message "All fixed"
```

### Add Reactions

```bash
# Add thumbs up
gh talk react PRRC_kwDOQN97u86UHqK7 👍

# Add to multiple comments (bulk operation)
gh talk react PRRC_aaa PRRC_bbb PRRC_ccc 👍

# Add rocket (by name)
gh talk react PRRC_kwDOQN97u86UHqK7 ROCKET

# Remove reaction
gh talk react PRRC_kwDOQN97u86UHqK7 👍 --remove
```

### Check PR Status

```bash
# Show review progress
gh talk status --pr 137

# Quick one-line summary
gh talk status --compact
```

### Combined Workflow

```bash
# Reply, react, and resolve in one command
gh talk reply PRRT_xxx "Fixed!" --react 👍 --resolve

# Get JSON output for scripting
gh talk list threads --format json | jq '.[] | select(.isResolved == false)'
```

### View Thread Details

```bash
# Show full conversation
gh talk show PRRT_kwDOQN97u85gQeTN
```

### Hide Comments

```bash
# Hide as spam
gh talk hide IC_kwDOQN97u87PVA8l --reason spam

# Unhide
gh talk unhide IC_kwDOQN97u87PVA8l
```

## Development

```bash
# Clone the repository
git clone https://github.com/hamishmorgan/gh-talk.git
cd gh-talk

# Download dependencies
go mod download

# Build
go build

# Run tests
go test ./...

# Install locally
gh extension install .

# Try it out
gh talk --help
gh talk list threads
```

## Development Commands

```bash
make help           # Show available commands
make build          # Build binary
make test           # Run tests
make lint           # Run linters
make ci             # Run all CI checks
```

## License

MIT

## Author

Hamish Morgan
