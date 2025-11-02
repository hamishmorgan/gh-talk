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
- [STRUCTURE.md](docs/STRUCTURE.md) - Project structure and design decisions
- [AGENTS.md](AGENTS.md) - AI agent development guidelines
- [GitHub CLI Extension Docs](https://docs.github.com/en/github-cli/github-cli/using-github-cli-extensions)

## Development

```bash
# Clone the repository
git clone https://github.com/hamishmorgan/gh-talk.git
cd gh-talk

# Build
go build

# Install locally
gh extension install .

# Test
./gh-talk
```

## License

MIT

## Author

Hamish Morgan

