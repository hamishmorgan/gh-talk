# gh-talk Documentation

Welcome to the gh-talk documentation. This directory contains all technical and development documentation for the project.

## Documentation Overview

### For Users

- **[Main README](../README.md)** - Quick start guide and feature overview

### For Developers

- **[SPEC.md](SPEC.md)** - Complete specification and technical design
  - Use cases and target users
  - Command structure and examples
  - Technical architecture
  - Implementation phases
  - Performance goals

- **[API.md](API.md)** - GitHub API capabilities and reference
  - GraphQL schema details
  - Available queries and mutations
  - Data types and structures
  - Rate limits and best practices
  - Code examples and workflows

- **[REAL-DATA.md](REAL-DATA.md)** - Real API response data and analysis
  - Actual ID formats from production
  - Complete data structures from live testing
  - Thread, comment, and reaction examples
  - Error responses and edge cases
  - Design implications from real data

- **[GH-CLI.md](GH-CLI.md)** - GitHub CLI (`gh`) analysis
  - Core `gh` capabilities and commands
  - PR and issue comment features
  - Strengths and limitations
  - Extension integration patterns
  - What gh-talk adds to the ecosystem

- **[WORKFLOWS.md](WORKFLOWS.md)** - Real-world usage patterns and workflows
  - How developers interact with PR conversations
  - Common code review patterns
  - Bot and automation workflows
  - AI agent integration scenarios
  - Best practices and team collaboration

- **[STRUCTURE.md](STRUCTURE.md)** - Project structure and design decisions
  - Directory layout rationale
  - Package responsibilities
  - Design decision explanations
  - Evolution path

## Quick Links

### Getting Started
1. Read the [main README](../README.md) for installation
2. Review [SPEC.md](SPEC.md) for the full feature set
3. Check [STRUCTURE.md](STRUCTURE.md) to understand the codebase layout

### Contributing
1. Understand the [project structure](STRUCTURE.md)
2. Follow the implementation plan in [SPEC.md](SPEC.md)
3. Review the [API reference](API.md) for GitHub API details

## External Resources

- [GitHub CLI Extension Documentation](https://docs.github.com/en/github-cli/github-cli/creating-github-cli-extensions)
- [Effective Go](https://go.dev/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)

---

**Last Updated**: 2025-11-02

