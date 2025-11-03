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

- **[GO-GH.md](GO-GH.md)** - go-gh library guide and patterns
  - GraphQL and REST client usage
  - Repository context detection
  - Terminal and table formatting
  - Error handling patterns
  - Best practices and common patterns

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

- **[DESIGN.md](DESIGN.md)** - Key design decisions and rationale
  - Thread ID system (Full IDs + Interactive + URLs)
  - Command syntax and patterns
  - Flag conventions
  - Error message patterns
  - Output formats and filtering

- **[CLI-FRAMEWORK.md](CLI-FRAMEWORK.md)** - CLI framework analysis and decision
  - What gh uses (custom, not Cobra)
  - Framework alternatives (Cobra, urfave/cli, stdlib, kong)
  - Comparison and trade-offs
  - Recommendation: Cobra (despite gh not using it)
  - Implementation strategy

- **[COBRA.md](COBRA.md)** - Cobra framework guide for gh-talk
  - Core concepts (commands, flags, subcommands)
  - gh-talk command structure with Cobra
  - Common patterns for our use case
  - Integration with go-gh library
  - Testing strategies
  - Complete implementation examples

- **[EXTENSION-PATTERNS.md](EXTENSION-PATTERNS.md)** - Analysis of successful gh extensions
  - Study of 5 popular extensions (gh-dash, gh-s, gh-poi, gh-copilot, gh-branch)
  - Project structures and patterns
  - Framework choices (gh-copilot uses Cobra!)
  - Common testing approaches
  - Best practices and anti-patterns
  - Validation of our design choices

- **[ENGINEERING.md](ENGINEERING.md)** - Engineering practices and CI/CD
  - Testing strategy (unit, integration, contract, E2E)
  - CI/CD workflows (test, lint, build, security)
  - Linting configuration (Go, Markdown, YAML)
  - Code coverage and quality gates
  - Makefile for common tasks
  - Pre-release checklist

- **[ENVIRONMENT.md](ENVIRONMENT.md)** - Environment variables reference
  - GitHub CLI variables (GH_TOKEN, GH_HOST, GH_REPO, etc.)
  - Terminal variables (NO_COLOR, CLICOLOR, etc.)
  - gh-talk specific variables
  - Priority and precedence rules
  - Usage examples and testing

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

