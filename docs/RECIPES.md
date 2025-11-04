# gh-talk Recipes

Common workflows and patterns for managing GitHub PR conversations.

## Quick Start

### First Time Setup

```bash
# Install gh-talk
gh extension install hamishmorgan/gh-talk

# Verify installation
gh talk --help
```

### Your First Conversation

```bash
# List all threads in a PR
gh talk list threads --pr 123

# Reply to a thread
gh talk reply PRRT_kwDOQN97u85gQeTN "Fixed in latest commit"

# Add emoji reaction
gh talk react PRRC_kwDOQN97u86UHqK7 👍
```

## ID Format Guide

### Understanding Node IDs

gh-talk uses GitHub's global node IDs, not numeric database IDs.

**Thread IDs** (Pull Request Review Threads):

- Format: `PRRT_kwDOQN97u85gQeTN`
- Prefix: `PRRT_`
- Use for: `gh talk reply`, `gh talk resolve`

**Comment IDs** (Review Thread Comments or Issue Comments):

- Format: `PRRC_kwDOQN97u86UHqK7` or `IC_kwDOQN97u87PVA8l`
- Prefix: `PRRC_` (PR review comment) or `IC_` (issue comment)
- Use for: `gh talk react`, `gh talk hide`

**How to find IDs:**

```bash
# List all thread IDs
gh talk list threads --pr 123

# Get thread details with comment IDs
gh talk show PRRT_kwDOQN97u85gQeTN
```

**Common mistake:**

```bash
# ❌ Wrong: Using numeric database ID
gh talk react 2486231843 👍
# Error: invalid comment ID: 2486231843
# You provided a numeric database ID...

# ✅ Right: Using node ID
gh talk react PRRC_kwDOQN97u86UHqK7 👍
```

## Common Workflows

### Recipe 1: Address All Review Feedback

**Scenario**: PR has review feedback that needs addressing.

**Steps**:

```bash
# Step 1: See what needs attention
gh talk list threads --pr <PR>

# Step 2: For each unresolved thread, reply and resolve
gh talk reply PRRT_xxx "Fixed in commit abc123" --resolve

# Step 3: Add reaction to show you've addressed it
gh talk reply PRRT_xxx "Fixed" --resolve --react 👍

# Step 4: Hide the original comment to declutter
gh talk hide PRRC_original_comment --reason resolved
```

**Complete example:**

```bash
# View unresolved feedback
gh talk list threads --pr 137

# Address each point
gh talk reply PRRT_kwDOQN97u85gQeTN "Renamed variable as suggested" --resolve --react ❤️
gh talk reply PRRT_kwDOQN97u85gQfTo "Added error handling" --resolve --react 👍
gh talk reply PRRT_kwDOQN97u85gRzXw "Tests updated" --resolve --react 🎉
```

### Recipe 2: Clean Up Resolved Conversations

**Scenario**: PR has many resolved threads cluttering the view.

**Steps**:

```bash
# Step 1: List resolved threads
gh talk list threads --pr <PR> | grep "✓ Resolved"

# Step 2: Hide comments in resolved threads
# Find comment IDs with: gh talk show PRRT_xxx
gh talk hide PRRC_comment1 PRRC_comment2 PRRC_comment3 --reason resolved

# Bulk hide multiple comments at once
gh talk hide PRRC_aaa PRRC_bbb PRRC_ccc PRRC_ddd --reason resolved
```

**Using jq for automation:**

```bash
# Extract resolved thread IDs (requires custom script)
gh talk list threads --pr 137 --json | \
  jq -r '.[] | select(.isResolved == true) | .id'
```

### Recipe 3: Acknowledge Feedback Without Resolving

**Scenario**: Reviewer made a good point, but not fixing it now.

**Steps**:

```bash
# Acknowledge the feedback
gh talk reply PRRT_xxx "Great point! Will address in a follow-up PR"

# Add reaction to show you read it
gh talk react PRRC_xxx 👀

# Do NOT use --resolve flag
# Thread stays open as a reminder
```

**Example:**

```bash
gh talk reply PRRT_kwDOQN97u85gQeTN "Good catch! This is out of scope for this PR, but I created issue #45 to track it"
gh talk react PRRC_kwDOQN97u86UHqK7 👍
```

### Recipe 4: Batch Operations

**Scenario**: Need to perform the same operation on multiple comments.

**React to multiple comments:**

```bash
# Add thumbs up to several comments
gh talk react PRRC_aaa PRRC_bbb PRRC_ccc 👍

# Add heart to all reviewer comments
gh talk react PRRC_comment1 PRRC_comment2 👀
```

**Hide multiple comments:**

```bash
# Hide all spam comments at once
gh talk hide IC_spam1 IC_spam2 IC_spam3 --reason spam

# Hide resolved thread comments
gh talk hide PRRC_resolved1 PRRC_resolved2 --reason resolved
```

### Recipe 5: Review and React to Feedback

**Scenario**: Quickly review and acknowledge all new feedback.

**Steps**:

```bash
# Step 1: See all threads
gh talk list threads --pr <PR>

# Step 2: Read each thread
gh talk show PRRT_xxx

# Step 3: React to each comment
gh talk react PRRC_comment1 👀  # Acknowledged
gh talk react PRRC_comment2 ❤️  # Agreed
gh talk react PRRC_comment3 🎉  # Excellent suggestion
```

**Reaction guide:**

- 👀 - "I've seen this"
- 👍 - "Will do" / "Done"
- ❤️ - "Thanks!" / "Great idea"
- 🎉 - "Excellent feedback!"
- 😕 - "Not sure about this"

### Recipe 6: Resolve Thread and Hide Original

**Scenario**: Clean resolution - thread resolved and decluttered.

**Steps**:

```bash
# Step 1: Reply and resolve
gh talk reply PRRT_xxx "Fixed" --resolve

# Step 2: Get the original comment ID
gh talk show PRRT_xxx  # Look for PRRC_ IDs

# Step 3: Hide the original comment
gh talk hide PRRC_original --reason resolved
```

**One-liner workflow:**

```bash
# Reply, resolve, and hide in sequence
gh talk reply PRRT_xxx "Fixed" --resolve && \
gh talk hide PRRC_original --reason resolved
```

## Advanced Patterns

### Automation with Scripts

**Find all unresolved threads:**

```bash
#!/bin/bash
# Save as: gh-talk-unresolved.sh

PR=$1
gh talk list threads --pr $PR | grep "○ Unresolved" | awk '{print $1}'
```

**Reply to all threads in a file:**

```bash
#!/bin/bash
# Save as: gh-talk-bulk-reply.sh

MESSAGE=$1
shift
for THREAD in "$@"; do
  gh talk reply $THREAD "$MESSAGE"
done
```

**Usage:**

```bash
./gh-talk-bulk-reply.sh "Fixed" PRRT_xxx PRRT_yyy PRRT_zzz
```

### Integration with PR Workflow

**After pushing a fix:**

```bash
# 1. Push your changes
git push

# 2. Reply to relevant threads
gh talk reply PRRT_addressed_in_commit "Fixed in latest commit" --resolve

# 3. Add status comment on PR
gh pr comment <PR> --body "Addressed all review feedback"
```

**Before merging:**

```bash
# 1. Check all threads are resolved
gh talk list threads --pr <PR>

# 2. Ensure no unresolved threads
gh talk list threads --pr <PR> | grep -c "○ Unresolved"
# Should output: 0

# 3. Merge when clean
gh pr merge <PR> --squash
```

## Troubleshooting

### Problem: Invalid ID format error

**Error message:**

```
Error: invalid comment ID: 2486231843
You provided a numeric database ID, but gh-talk requires node IDs
```

**Solution:**

Don't use numeric IDs from API responses. Use `gh talk list` to get node IDs:

```bash
# Find correct IDs
gh talk list threads --pr <PR>

# Get thread details with comment IDs
gh talk show PRRT_xxx
```

### Problem: Can't find thread ID

**Solution:**

```bash
# List all threads with full details
gh talk list threads --pr <PR>

# Search GitHub PR page
# Click "View conversation" on any review comment
# The URL will contain thread identifiers
```

### Problem: Thread won't resolve

**Check**:

1. Are you using the correct thread ID (starts with `PRRT_`)?
2. Do you have permission to resolve threads?
3. Is the thread already resolved?

```bash
# Check thread status
gh talk show PRRT_xxx
# Look for: "Status: Resolved" or "Status: Unresolved"
```

### Problem: Command not found

**Solution:**

```bash
# Reinstall extension
gh extension remove hamishmorgan/gh-talk
gh extension install hamishmorgan/gh-talk

# Verify
gh talk --version
```

## Tips and Tricks

### Use Aliases

Add to your shell config:

```bash
# ~/.bashrc or ~/.zshrc
alias ghl='gh talk list threads --pr'
alias ghr='gh talk reply'
alias ghreact='gh talk react'
alias ghh='gh talk hide'
```

**Usage:**

```bash
ghl 137                    # List threads in PR 137
ghr PRRT_xxx "Fixed"       # Quick reply
ghreact PRRC_xxx 👍        # Quick reaction
```

### Copy Thread IDs

**macOS:**

```bash
gh talk list threads --pr 137 | grep "Unresolved" | pbcopy
```

**Linux:**

```bash
gh talk list threads --pr 137 | grep "Unresolved" | xclip -selection clipboard
```

### Check Before Bulk Operations

```bash
# Dry run by echoing first
for ID in PRRC_aaa PRRC_bbb PRRC_ccc; do
  echo "Would hide: $ID"
done

# If looks good, actually hide
for ID in PRRC_aaa PRRC_bbb PRRC_ccc; do
  gh talk hide $ID --reason resolved
done
```

### Quick Status Check

```bash
# See thread counts
gh talk list threads --pr <PR> | grep -c "Resolved"
gh talk list threads --pr <PR> | grep -c "Unresolved"
```

## Next Steps

- Read [SPEC.md](SPEC.md) for complete feature list
- See [API.md](API.md) for GraphQL details
- Check [WORKFLOWS.md](WORKFLOWS.md) for integration patterns
- Review [USER-FEEDBACK.md](USER-FEEDBACK.md) for enhancement ideas

## Getting Help

**Issues or questions:**

```bash
# Check existing issues
gh issue list --repo hamishmorgan/gh-talk

# Create new issue
gh issue create --repo hamishmorgan/gh-talk
```

**Debug information:**

```bash
# Show version
gh talk --version

# Show help
gh talk --help
gh talk <command> --help
```

---

**Contribute**: Found a useful recipe? [Submit a PR](https://github.com/hamishmorgan/gh-talk/pulls) to add it here!
