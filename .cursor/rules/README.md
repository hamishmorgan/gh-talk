# Cursor Rules Directory

This directory contains auto-applying rules for AI agents working on gh-talk.

## Structure

**`.cursor/commands/`** - Slash-invokable workflows
- See `.cursor/README.md` for full command list
- Files: pr-create.md, pr-iterate.md, issue-implement.md, etc.
- Invoke with: `/pr-create`, `/issue-implement`, etc.

**`.cursor/rules/`** - Auto-applying rules (this directory)
- Rules that load automatically based on frontmatter
- Uses YAML frontmatter (alwaysApply, globs)
- Files: creating-rules.mdc

## Files in This Directory

**`create-rule.mdc`** - Guide for creating cursor workflow rules
- **Load strategy**: Context-based (`globs: [".cursor/**/*.mdc", "AGENTS.md"]`)
- **Purpose**: Creating full workflow rules in .cursor/rules/
- **When it loads**: When editing rule files or AGENTS.md

**`create-command.mdc`** - Guide for creating command wrappers
- **Load strategy**: Context-based (`globs: [".cursor/commands/**/*.md"]`)
- **Purpose**: Creating thin command wrappers in .cursor/commands/
- **When it loads**: When editing command files

## Frontmatter Format

Rules in this directory use YAML frontmatter:

```yaml
---
description: "Brief purpose"
alwaysApply: true    # or false
globs: ["pattern"]   # optional
---
```

## Relationship to Commands

**Commands** (`.cursor/commands/*.md`):
- No frontmatter
- Invoked via slash commands
- User triggers explicitly

**Rules** (`.cursor/rules/*.mdc`):
- Has frontmatter
- Auto-loaded by Cursor
- Based on alwaysApply or globs

## Adding New Workflows

**For slash commands**: Create in `.cursor/commands/`
- File: `command-name.md`
- Format: `# command-name` as first line
- No frontmatter needed

**For auto-applying rules**: Create in `.cursor/rules/`
- File: `rule-name.mdc`
- Must have frontmatter
- Set alwaysApply or globs

See `create-rule.mdc` for complete guide.

---

**Last Updated**: 2025-11-03
