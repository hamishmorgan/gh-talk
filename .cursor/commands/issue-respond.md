# issue-respond

# Responding to GitHub Issues

**Purpose**: Handle issue discussions, questions, and comments.

**When to Use**: When responding to or managing issue conversations (invoke with `/issue-respond`)

**Prerequisites**:

- Issue number
- `gh` CLI and `gh-talk` installed

## Response Workflow

### 1. Read Issue and Comments

```bash
# View issue
gh issue view <NUMBER>

# See all comments
gh issue view <NUMBER> --comments

# Or list for multiple issues
gh issue list --assignee @me
```

### 2. Understand Context

**Check:**

- Original issue description
- All comments and discussion
- Related PRs or issues
- Labels and milestones

### 3. Respond Appropriately

**Types of responses:**

**Answer a question:**

```bash
gh issue comment <NUMBER> --body "The reason this happens is...

[Explanation]

Does this help?

🤖 Cursor"
```

**Provide update:**

```bash
gh issue comment <NUMBER> --body "## Progress Update

✓ Completed X
⏳ Working on Y
📋 Next: Z

ETA: [timeframe]

🤖 Cursor"
```

**Request more information:**

```bash
gh issue comment <NUMBER> --body "To help debug this, could you provide:

- Version of gh-talk
- Output of: gh talk list threads --pr <PR>
- Full error message

Thanks!

🤖 Cursor"
```

**Close as completed:**

```bash
gh issue close <NUMBER> --reason completed --comment "Fixed in PR #<PR>

This is now available in v0.2.0.

🤖 Cursor"
```

**Close as won't fix:**

```bash
gh issue close <NUMBER> --reason "not planned" --comment "After discussion, we've decided not to implement this because [reason].

Alternative: [suggestion if applicable]

Thanks for the suggestion!

🤖 Cursor"
```

## Using gh-talk for Issues

**Add reactions:**

```bash
# React to issue body
gh talk react I_xxx 👍

# React to helpful comments
gh talk react IC_xxx ❤️
```

**Note**: Issues don't have review threads, only comments

## Best Practices

### DO ✅

**Be helpful:**

- Provide clear explanations
- Link to relevant docs
- Offer alternatives
- Thank for contributions

**Be timely:**

- Respond within 24-48 hours
- Set expectations for resolution time
- Update if delays occur

**Be specific:**

- Quote relevant parts when responding
- Reference line numbers or files
- Provide examples

### DON'T ❌

**Be dismissive:**
❌ "This won't work"
✅ "This approach has challenges: [explain]"

**Leave hanging:**
❌ No response for weeks
✅ "Looking into this, will update by [date]"

**Be vague:**
❌ "Maybe later"
✅ "Not in current roadmap, but open to PR if you want to implement"

## Common Scenarios

### Bug Report Response

```bash
gh issue comment <NUM> --body "Thanks for the report!

I can reproduce this. The issue is [explanation].

I'll fix this in the next release.

Tracking: #<PR-when-created>

🤖 Cursor"
```

### Feature Request Response

```bash
gh issue comment <NUM> --body "Interesting idea!

This aligns with [project goal]. I'll add to roadmap.

Priority: Medium
Timeline: Likely v0.3.0

🤖 Cursor"
```

### Cannot Reproduce

```bash
gh issue comment <NUM> --body "I'm unable to reproduce this.

Could you provide:
- Exact command run
- Full output/error
- gh-talk version

This will help me investigate.

🤖 Cursor"
```

### Duplicate Issue

```bash
gh issue close <NUM> --reason duplicate --comment "Duplicate of #<OTHER>

Discussion continuing there.

🤖 Cursor"
```

## Agent Signature

**ALL issue comments MUST end with:**

```
🤖 Cursor
```

## Quick Reference

```bash
# View and respond
gh issue view <NUM>
gh issue comment <NUM> --body "...

🤖 Cursor"

# Close with reason
gh issue close <NUM> --reason completed
gh issue close <NUM> --reason "not planned"

# React with gh-talk
gh talk react I_xxx 👍
gh talk react IC_xxx ❤️
```

---

**Related**:

- `issue-create.mdc` - Creating issues
- `pr-create.mdc` - Creating PRs to address issues

---
Note: This content duplicated in .cursor/rules/issue-respond.mdc for auto-loading
