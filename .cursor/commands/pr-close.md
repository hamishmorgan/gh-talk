# pr-close

# Closing Pull Requests

**Purpose**: Properly close a PR when deciding not to merge (different approach, superseded, etc.)

**When to Use**: When abandoning PR work or changing direction (invoke with `/pr-close`)

**Prerequisites**:

- Open PR to close
- Clear reason for closing

## When to Close (Not Merge)

**Close when:**

- Taking different approach (will create new PR)
- Work is superseded by another PR
- Requirements changed
- Decided not to implement
- Stale/outdated (no longer relevant)

**Don't close if:**

- Just needs more work → keep open, continue with `/pr-iterate`
- CI failing → fix it, don't abandon
- Review feedback → address it, don't close

## Closing Process

### 1. Add Closing Comment

**ALWAYS explain why:**

```bash
gh pr comment <PR> --body "## Closing This PR

**Reason**: [Why you're closing]

**What's Next**: [Alternative approach or new PR]

Example:
- Taking different approach in PR #<NEW>
- Requirements changed, no longer needed
- Superseded by PR #<OTHER>
- Decided against this implementation

Thanks for the review feedback received.

🤖 Cursor"
```

### 2. Close the PR

```bash
gh pr close <PR>
```

**Result**: PR closed but preserved for reference

### 3. Clean Up Local Branch

```bash
# Switch back to main
git checkout main
git pull origin main

# Delete feature branch
git branch -d <branch-name>  # Soft delete (safe)
# or
git branch -D <branch-name>  # Force delete (if unmerged)

# Delete remote branch
git push origin --delete <branch-name>
```

## Common Scenarios

### Scenario 1: Taking Different Approach

```bash
gh pr comment <PR> --body "## Closing This PR

**Reason**: After review, taking a different approach

The feedback showed this design has issues. Creating
new PR with better approach.

**New PR**: Will create #<NEW> with revised design

Thanks @reviewer for the feedback!

🤖 Cursor"

gh pr close <PR>
```

**Then**: Use `/pr-create` for new approach

### Scenario 2: Superseded by Another PR

```bash
gh pr comment <PR> --body "## Closing This PR

**Reason**: Superseded by PR #<OTHER>

That PR implements the same feature with better approach.

Closing this as duplicate work.

🤖 Cursor"

gh pr close <PR>
```

### Scenario 3: Requirements Changed

```bash
gh pr comment <PR> --body "## Closing This PR

**Reason**: Requirements changed

After discussion, we're not implementing this feature.

See: #<ISSUE> for context

Thanks for the work that went into this.

🤖 Cursor"

gh pr close <PR>
```

### Scenario 4: Stale PR

```bash
gh pr comment <PR> --body "## Closing This PR

**Reason**: Stale - no activity for [timeframe]

This PR hasn't been updated in [time]. Closing to keep
PR list clean.

Can reopen if still needed.

🤖 Cursor"

gh pr close <PR>
```

## Etiquette

**DO ✅:**

- Always explain why closing
- Thank reviewers for feedback
- Indicate what's next (if anything)
- Be respectful of work done
- Close cleanly (comment first, then close)

**DON'T ❌:**

- Close without comment (rude and confusing)
- Close just because it's hard (persist or ask for help)
- Leave orphan branches (clean up)

## Clean Up Checklist

After closing PR:

- [ ] Closing comment added with clear reason
- [ ] PR closed
- [ ] Local branch deleted
- [ ] Remote branch deleted
- [ ] Related issue updated (if any)

## Quick Reference

```bash
# Close with explanation
gh pr comment <PR> --body "Reason for closing...

🤖 Cursor"
gh pr close <PR>

# Clean up
git checkout main
git branch -D <branch>
git push origin --delete <branch>
```

## Agent Signature

**ALL closing comments MUST include:**

```
🤖 Cursor
```

---

**Related**:

- `pr-create.mdc` - Creating new PR with different approach
- `pr-merge.mdc` - Normal completion path (merge, not close)

---
Note: This content duplicated in .cursor/rules/pr-close.mdc for auto-loading
