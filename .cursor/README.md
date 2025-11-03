# Cursor Configuration

This directory contains configuration for Cursor IDE AI agents.

## Structure

```
.cursor/
├── commands/         # Slash-invokable workflows (/pr-create, /issue-respond, etc.)
│   ├── pr-create.md
│   ├── pr-iterate.md
│   └── ...
└── rules/           # Auto-applying rules with frontmatter (alwaysApply, globs)
    ├── creating-rules.mdc
    └── README.md
```

## Commands vs Rules

**`.cursor/commands/`** - Slash command workflows
- Invoked explicitly: `/pr-create`, `/issue-implement`
- No frontmatter needed
- Plain markdown files (.md extension)
- User triggers when needed

**`.cursor/rules/`** - Auto-applying rules
- Loaded automatically based on frontmatter (alwaysApply, globs)
- Has YAML frontmatter
- .mdc extension
- Cursor loads based on context

## Available Commands

Type `/` in Cursor chat to see all available commands.

**PR Workflows:**
- `/pr-create` - Create and submit PR
- `/pr-iterate` - Handle feedback and fixes
- `/pr-merge` - Final review and merge
- `/pr-close` - Abandon PR with explanation
- `/pr-review` - Review others' PRs

**Issue Workflows:**
- `/issue-create` - Create well-structured issues
- `/issue-respond` - Handle issue discussions
- `/issue-implement` - Issue → Implementation → PR
- `/issue-close` - Close issues properly

**Reference:**
- `/emoji-semantics` - Emoji reaction meanings

**Meta:**
- `/creating-rules` - Guide for creating new workflows

## Available Rules

Rules in `.cursor/rules/` are loaded automatically based on frontmatter configuration.

See `.cursor/rules/README.md` for details.

---

**Official Documentation**: [Cursor Custom Slash Commands](https://cursor.com/changelog/1-6)

**Last Updated**: 2025-11-03

