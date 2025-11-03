# GitHub Conversation Workflows and Usage Patterns

**Understanding how people and agents interact with PR and Issue conversations**

## Overview

This document describes the real-world workflows, conversation patterns, and usage scenarios for GitHub PR and Issue conversations. It synthesizes best practices, common patterns, and automation needs to inform how `gh-talk` will be used by both humans and AI agents.

## Target Users

### 1. Developers (Primary Users)

**Profile:**

- Write code and submit PRs daily
- Respond to review feedback
- Review teammates' code
- Manage 5-20 active PRs at once
- Context-switch frequently between tasks

**Pain Points:**

- Browser context-switching breaks flow
- Finding unresolved threads is tedious
- Quick acknowledgments require full page loads
- Cannot script review responses
- Mobile reviews are difficult

**Workflow Needs:**

- Fast thread navigation
- Quick reactions and replies
- Bulk thread resolution
- Scriptable operations
- Terminal-first workflow

### 2. Code Reviewers

**Profile:**

- Review 10-50 PRs per week
- Provide detailed feedback
- Follow up on addressed comments
- Ensure quality standards
- Track review progress

**Pain Points:**

- Tracking which comments were addressed
- Knowing when threads are truly resolved
- Following up on ignored feedback
- Managing review backlog
- Filtering noise from signal

**Workflow Needs:**

- See only unresolved threads
- Filter by file or topic
- Quick approval workflow
- Thread status tracking
- Bulk resolution after verification

### 3. AI Agents & Bots

**Profile:**

- Automated code review (linting, security, coverage)
- Dependency updates (Dependabot, Renovate)
- CI/CD feedback
- Auto-merge automation
- Code generation assistance

**Pain Points:**

- Web UI requires browser automation
- Rate limits on API calls
- Complex GraphQL queries
- No CLI tools for thread management
- Difficult to script conversations

**Workflow Needs:**

- Scriptable thread creation
- Programmatic resolution
- Bulk operations
- JSON output for parsing
- Integration with CI/CD

### 4. Maintainers & Project Leads

**Profile:**

- Oversee multiple repositories
- Ensure reviews are completed
- Enforce code quality
- Manage contributor relationships
- Track project health

**Pain Points:**

- Visibility into review status
- Identifying blocked PRs
- Managing contributor feedback
- Ensuring timely responses
- Maintaining code standards

**Workflow Needs:**

- Overview of pending reviews
- Unresolved thread reporting
- Bulk thread management
- Team workflow automation
- Quality metrics

## Common Conversation Workflows

### Workflow 1: PR Author Addresses Review Feedback

**Typical Flow:**

```
1. PR submitted → Review requested
2. Reviewer leaves 15 comments across 5 files
3. Author receives notification
4. Author addresses each comment:
   - Fix code
   - Reply to thread
   - Mark thread as resolved
5. Reviewer verifies fixes
6. Reviewer approves PR
7. PR merged
```

**Current Experience (Web UI):**

```
✗ Load PR in browser
✗ Click "Files changed" tab
✗ Scroll to find comments
✗ Read comment
✗ Switch to editor
✗ Fix code
✗ Commit changes
✗ Switch back to browser
✗ Find comment again (lost scroll position)
✗ Click "Reply"
✗ Type message
✗ Click "Resolve conversation"
✗ Repeat 14 more times
```

**With gh-talk:**

```bash
# List unresolved threads
gh talk list threads --unresolved

# Output:
# Thread ID         File              Line  Author      Preview
# PRRT_abc123      src/api.go         42   reviewer1   Consider using...
# PRRT_def456      src/db.go          89   reviewer2   This could be...

# Address in editor, commit, then:
gh talk reply PRRT_abc123 "Fixed in commit abc123" --resolve

# Or bulk resolve after addressing all
gh talk list threads --unresolved --json id | \
  xargs -I {} gh talk resolve {}
```

**Time Saved:** ~60% (from 30 min to 12 min for 15 comments)

### Workflow 2: Quick Acknowledgments

**Scenario:** Reviewer leaves helpful suggestions, author wants to acknowledge without lengthy replies.

**Current Experience (Web UI):**

```
✗ Load PR in browser
✗ Find comment
✗ Click 👍 reaction
✗ Wait for page to update
✗ Repeat for each comment
```

**With gh-talk:**

```bash
# React to multiple comments quickly
gh talk list threads --author reviewer1 --json comments | \
  jq -r '.[].comments[0].id' | \
  xargs -I {} gh talk react {} 👍

# Or individually
gh talk react PRRC_xyz789 👍
gh talk react PRRC_abc456 🎉
```

**Time Saved:** ~80% (from 5 min to 1 min for 10 reactions)

### Workflow 3: Reviewer Verifies Fixes

**Typical Flow:**

```
1. Author marks threads as resolved
2. Reviewer needs to verify fixes
3. Check each resolved thread
4. Verify code changes
5. Approve if all fixed
```

**Current Experience (Web UI):**

```
✗ Load PR
✗ Click "Conversation" tab
✗ Manually find resolved threads
✗ Click "Show resolved"
✗ Check each one individually
✗ No way to filter "recently resolved"
```

**With gh-talk:**

```bash
# See recently resolved threads
gh talk list threads --resolved --since yesterday

# Check specific thread
gh talk show PRRT_abc123 --with-diff

# If satisfied, approve via gh
gh pr review --approve -b "All feedback addressed!"
```

**Time Saved:** ~50% (from 10 min to 5 min for verification)

### Workflow 4: Managing Review Backlog

**Scenario:** Reviewer has 20 PRs to review, needs to prioritize.

**Current Experience (Web UI):**

```
✗ Open each PR in browser
✗ Manually check for unresolved threads
✗ Try to remember which need attention
✗ No way to batch check status
```

**With gh-talk:**

```bash
# Check all your PRs for unresolved threads
gh pr list --assignee @me --json number | \
  jq -r '.[].number' | \
  while read pr; do
    echo "PR #$pr:"
    gh talk list threads --pr $pr --unresolved | head -3
  done

# Or create a script to find PRs needing attention
gh talk status --assignee @me
# Shows:
# PR #123: 5 unresolved threads (awaiting author)
# PR #456: 12 unresolved threads (awaiting your review)
# PR #789: 0 unresolved threads (ready to merge)
```

**Time Saved:** ~70% (from 20 min to 6 min to triage 20 PRs)

### Workflow 5: Hiding Noise

**Scenario:** PR has bot comments, outdated threads, and spam that clutters the view.

**Current Experience (Web UI):**

```
✗ Manually minimize each comment
✗ No bulk operations
✗ Hidden comments still count in totals
✗ No filtering for hidden comments
```

**With gh-talk:**

```bash
# Hide all bot comments from a specific bot
gh talk list comments --author dependabot --json id | \
  jq -r '.[].id' | \
  xargs -I {} gh talk hide {} --reason outdated

# Hide spam
gh talk hide PRRC_spam123 --reason spam

# Hide resolved discussions
gh talk list threads --resolved --json comments | \
  jq -r '.[].comments[].id' | \
  xargs -I {} gh talk hide {} --reason resolved
```

**Time Saved:** ~90% (from 10 min to 1 min for bulk hiding)

## Bot & Automation Workflows

### Bot Workflow 1: Automated Code Review

**Use Case:** CI bot runs linters and posts review comments.

**Implementation:**

```bash
#!/bin/bash
# CI script for automated review

# Run linter and get issues
lint_results=$(npm run lint --json)

# Parse results and create review threads
echo "$lint_results" | jq -r '.[] | @json' | while read issue; do
  file=$(echo "$issue" | jq -r '.file')
  line=$(echo "$issue" | jq -r '.line')
  message=$(echo "$issue" | jq -r '.message')
  
  # Post review comment
  gh api graphql -f query="
    mutation {
      addPullRequestReviewThread(input: {
        pullRequestId: \"$PR_ID\"
        path: \"$file\"
        line: $line
        body: \"⚠️ Linter: $message\"
      }) {
        thread { id }
      }
    }
  "
done

# If all checks pass, approve
if [ "$lint_errors" = "0" ]; then
  gh pr review --approve -b "✅ All automated checks passed"
fi
```

**With gh-talk Enhancements:**

```bash
# Check if automated comments were addressed
gh talk list threads --author ci-bot --unresolved --json count

# Auto-resolve threads where code was fixed
gh talk list threads --author ci-bot --json id,path,line | \
  while read thread; do
    # Check if issue still exists
    if ! lint_check "$thread"; then
      gh talk resolve "$thread" --message "✅ Issue fixed"
    fi
  done
```

### Bot Workflow 2: Dependency Update Management

**Use Case:** Dependabot creates PRs, team reviews and merges.

**Current Challenge:**

- Dependabot PRs pile up
- Need to verify compatibility
- Want to batch merge safe updates
- Manual review is time-consuming

**With gh-talk:**

```bash
#!/bin/bash
# Auto-merge safe dependency updates

# Get all Dependabot PRs
gh pr list --author app/dependabot --json number,title | \
  jq -r '.[] | select(.title | contains("patch")) | .number' | \
  while read pr; do
    # Check if CI passed
    if gh pr checks $pr | grep -q "All checks have passed"; then
      # Check for unresolved threads
      unresolved=$(gh talk list threads --pr $pr --unresolved --json count)
      
      if [ "$unresolved" = "0" ]; then
        # Add approval reaction
        gh talk react $(gh pr view $pr --json comments | \
          jq -r '.comments[0].id') 🚀
        
        # Auto-merge
        gh pr merge $pr --auto --squash
      fi
    fi
  done
```

### Bot Workflow 3: AI Code Assistant Response

**Use Case:** AI assistant (like GitHub Copilot or Claude) responds to code review questions.

**Scenario:**

```
Reviewer: "Can you explain why you used this algorithm here?"
```

**AI Agent Workflow:**

```bash
#!/bin/bash
# AI agent responds to questions

# Monitor for questions directed at bot
gh talk list threads --unresolved --json id,comments | \
  jq -r '.[] | select(.comments[-1].body | contains("@ai-assistant")) | .id' | \
  while read thread_id; do
    # Get thread context
    thread_data=$(gh talk show $thread_id --json body,path,line,diffHunk)
    
    # Generate AI response
    response=$(ai_generate_response "$thread_data")
    
    # Post response
    gh talk reply $thread_id "$response"
    
    # Add reaction to original question
    gh talk react $(echo "$thread_data" | jq -r '.comments[-1].id') 👀
  done
```

### Bot Workflow 4: Auto-Resolution After CI

**Use Case:** CI fixes issues automatically, bot resolves threads.

**Implementation:**

```bash
#!/bin/bash
# Auto-resolve threads after automated fixes

# Get threads about formatting
gh talk list threads --unresolved --json id,body | \
  jq -r '.[] | select(.body | contains("formatting")) | .id' | \
  while read thread; do
    # Run auto-formatter
    npm run format
    
    # Commit changes
    git add -A
    git commit -m "style: auto-format code"
    git push
    
    # Resolve thread
    gh talk resolve $thread --message "🤖 Auto-formatted by CI"
  done
```

## Team Collaboration Patterns

### Pattern 1: Pair Programming Follow-up

**Scenario:** Two developers pair on a feature, submit PR, address reviews together.

**Workflow:**

```bash
# Developer A lists unresolved threads
gh talk list threads --unresolved --format table

# Developer B shares screen, they discuss each thread
# Developer A addresses comments while B watches

# Quick acknowledgments as they go
gh talk react PRRC_123 👍  # "Good point"
gh talk react PRRC_456 🎉  # "Great suggestion!"

# Reply with explanations
gh talk reply PRRT_789 "We considered that but chose X because..."

# Resolve as addressed
gh talk resolve PRRT_789
```

### Pattern 2: Async Code Review

**Scenario:** Team distributed across timezones, reviews happen async.

**Workflow:**

```bash
# Morning routine (reviewer in EU):
# Check PRs assigned for review
gh pr list --review-requested @me --json number | \
  jq -r '.[].number' | \
  while read pr; do
    echo "=== PR #$pr ==="
    gh talk list threads --pr $pr --unresolved
  done

# Evening routine (author in US):
# Check feedback received during the day
gh talk list threads --pr $(gh pr view --json number | jq .number) \
  --since 8hours --format table

# Address each comment
gh talk reply PRRT_123 "Fixed in latest commit"
gh talk resolve PRRT_123

# Leave note for reviewer
gh pr comment "All feedback addressed, PTAL @reviewer"
```

### Pattern 3: Mentorship

**Scenario:** Senior dev reviews junior dev's PR, provides learning opportunities.

**Workflow:**

```bash
# Senior dev adds educational comments
gh api graphql -f query='...'  # Creates review thread

# Junior dev acknowledges and asks follow-up questions
gh talk reply PRRT_edu123 "Thanks! Can you explain more about why async is better here?"

# Senior dev provides detailed explanation
gh talk reply PRRT_edu123 "Sure! The key difference is..."

# Junior dev marks as understood
gh talk react PRRC_last456 ❤️  # "Learned something!"
gh talk resolve PRRT_edu123 --message "Got it, thanks for explaining!"
```

### Pattern 4: Community Open Source

**Scenario:** Maintainer reviews external contributor PR.

**Workflow:**

```bash
# Maintainer receives PR from external contributor
# Reviews and leaves feedback
gh pr review 789 --comment

# Contributor addresses some but not all feedback
# Maintainer checks progress
gh talk list threads --pr 789 --unresolved --author @me

# Maintainer provides encouraging reactions
gh talk react PRRC_fixed1 🎉  # Contributor fixed something!
gh talk react PRRC_fixed2 👍  # Another fix!

# Remaining issues
gh talk list threads --pr 789 --unresolved | \
  gh talk reply PRRT_remaining "No worries if you're stuck, I can help with this one"

# Eventually all resolved
gh pr review 789 --approve
gh pr comment 789 "Thanks for your contribution! 🎉"
```

## AI Agent Specific Workflows

### AI Workflow 1: Automated PR Creation with Context

**Scenario:** AI agent (like Cursor, GitHub Copilot Workspace) creates PR with full context.

**Implementation:**

```bash
#!/bin/bash
# AI agent creates PR with relevant context

# Create PR
pr_num=$(gh pr create --title "$title" --body "$body" --json number | jq -r '.number')

# Add context comments for reviewers
gh pr comment $pr_num "🤖 This PR was generated by AI to address issue #123"

# Add specific explanations for complex changes
for file in $changed_files; do
  explanation=$(ai_explain_changes "$file")
  gh api graphql -f query="
    mutation {
      addPullRequestReviewComment(input: {
        pullRequestId: \"$pr_id\"
        path: \"$file\"
        position: $position
        body: \"📝 **AI Explanation**: $explanation\"
      }) {
        comment { id }
      }
    }
  "
done

# Monitor for questions
while true; do
  # Check for mentions
  mentions=$(gh talk list threads --pr $pr_num --json comments | \
    jq -r '.[] | select(.comments[-1].body | contains("@ai-agent"))')
  
  if [ -n "$mentions" ]; then
    # Respond to questions
    # ... (see Bot Workflow 3)
  fi
  
  sleep 300  # Check every 5 minutes
done
```

### AI Workflow 2: Continuous Code Improvement

**Scenario:** AI agent monitors PRs and suggests improvements.

**Implementation:**

```bash
#!/bin/bash
# AI agent suggests improvements on new PRs

# Monitor new PRs
gh pr list --state open --json number,createdAt | \
  jq -r '.[] | select(.createdAt > (now - 3600)) | .number' | \
  while read pr; do
    # Get PR diff
    diff=$(gh pr diff $pr)
    
    # AI analyzes code
    suggestions=$(ai_analyze_code "$diff")
    
    # Post suggestions as review comments
    echo "$suggestions" | jq -r '.[] | @json' | while read suggestion; do
      file=$(echo "$suggestion" | jq -r '.file')
      line=$(echo "$suggestion" | jq -r '.line')
      body=$(echo "$suggestion" | jq -r '.body')
      
      # Create thread with suggestion
      gh api graphql -f query="
        mutation {
          addPullRequestReviewThread(input: {
            pullRequestId: \"$pr_id\"
            path: \"$file\"
            line: $line
            body: \"💡 **AI Suggestion**: $body\"
          }) {
            thread { id }
          }
        }
      "
    done
  done
```

### AI Workflow 3: Learning from Review Patterns

**Scenario:** AI agent learns from team's review patterns to provide better suggestions.

**Implementation:**

```bash
#!/bin/bash
# AI agent analyzes review history

# Get all resolved threads from last month
threads=$(gh talk list threads --resolved --since 30days --json body,path,line,comments)

# Extract patterns
patterns=$(echo "$threads" | ai_extract_patterns)

# Use patterns to pre-review new PRs
gh pr list --state open --json number | \
  jq -r '.[].number' | \
  while read pr; do
    # Get PR diff
    diff=$(gh pr diff $pr)
    
    # Check against learned patterns
    issues=$(ai_check_patterns "$diff" "$patterns")
    
    if [ -n "$issues" ]; then
      # Post proactive review
      gh pr review $pr --comment -b "
🤖 **Automated Pre-Review**

Based on historical review patterns, I noticed:
$issues

These are just suggestions - feel free to ignore if not applicable!
      "
    fi
  done
```

### AI Workflow 4: Context-Aware Assistance

**Scenario:** AI agent provides context-aware help during code review.

**Implementation:**

```bash
#!/bin/bash
# AI provides context when developer is stuck

# Monitor for specific keywords in threads
gh talk list threads --unresolved --json id,comments | \
  jq -r '.[] | select(.comments[-1].body | contains("not sure") or contains("confused")) | .id' | \
  while read thread; do
    # Get full context
    context=$(gh talk show $thread --with-diff --json)
    
    # Generate helpful explanation
    help=$(ai_generate_help "$context")
    
    # Post as thread reply
    gh talk reply $thread "
🤖 **AI Assistant**

I noticed you might be unsure about this. Here's some context:

$help

Let me know if you'd like more details!
    "
  done
```

## Best Practices Enabled by gh-talk

### 1. Quick Feedback Loop

**Practice:** Respond to all comments within 24 hours.

**With gh-talk:**

```bash
# Daily routine - check pending feedback
alias check-reviews='gh talk list threads --unresolved --format table'

# Quick acknowledgment
alias ack='gh talk react $1 👍'

# Set up notification
gh talk list threads --unresolved --since 24h && \
  notify-send "You have $(gh talk list threads --unresolved | wc -l) unresolved threads"
```

### 2. Meaningful Reactions

**Practice:** Use reactions to convey meaning, not just "seen".

**Conventions:**

```bash
👍  # "I agree" / "Will do"
❤️  # "Thanks for the help" / "Great suggestion"
🎉  # "Excellent point!" / "This is much better"
👀  # "I'm looking into this"
🚀  # "Ready to merge" / "Let's ship it"
😕  # "I'm not sure about this"
```

**Quick Commands:**

```bash
alias agree='gh talk react $1 👍'
alias thanks='gh talk react $1 ❤️'
alias looking='gh talk react $1 👀'
```

### 3. Clean Resolution

**Practice:** Only resolve threads when truly addressed, with explanation.

**With gh-talk:**

```bash
# Bad: Bulk resolve without checking
# gh talk list threads --json id | xargs -I {} gh talk resolve {}

# Good: Resolve with context
gh talk resolve PRRT_123 --message "Fixed by using the builder pattern as suggested"

# Good: Verify before resolving
gh talk show PRRT_123 --with-diff && \
  gh talk resolve PRRT_123 --message "Confirmed fix works in tests"
```

### 4. Thread Hygiene

**Practice:** Hide outdated/resolved threads to keep conversation focused.

**With gh-talk:**

```bash
# After resolving, hide to reduce clutter
gh talk list threads --resolved --json comments | \
  jq -r '.[].comments[] | .id' | \
  xargs -I {} gh talk hide {} --reason resolved

# Hide bot spam
gh talk list comments --author bot-account --json id | \
  jq -r '.[].id' | \
  xargs -I {} gh talk hide {} --reason outdated
```

### 5. Bulk Operations for Efficiency

**Practice:** Process similar feedback in batches.

**With gh-talk:**

```bash
# Address all formatting comments at once
npm run format
git commit -m "style: address formatting feedback"
gh talk list threads --unresolved --json id,body | \
  jq -r '.[] | select(.body | contains("format")) | .id' | \
  xargs -I {} gh talk resolve {} --message "Fixed by running formatter"

# React to all helpful suggestions
gh talk list threads --author reviewer1 --json comments | \
  jq -r '.[].comments[0].id' | \
  xargs -I {} gh talk react {} 👍
```

## Common Conversation Patterns

### Pattern: Acknowledge → Address → Resolve

```
Reviewer: "This function is too complex, consider breaking it down"
Author: 👍 (acknowledgment via reaction)
Author: *refactors code*
Author: "Broke it into 3 smaller functions, see latest commit"
Author: *resolves thread*
```

**With gh-talk:**

```bash
gh talk react PRRC_complex123 👍
# ... make changes ...
gh talk reply PRRT_complex123 "Broke it into 3 smaller functions" --resolve
```

### Pattern: Question → Answer → Thanks

```
Author: "Why do we need to handle this edge case?"
Reviewer: "Good question! Here's why..." 
Author: ❤️ (thanks reaction)
Author: *resolves thread*
```

**With gh-talk:**

```bash
gh talk reply PRRT_question456 "Good question! Here's why..."
# Author responds
gh talk react PRRC_answer789 ❤️
```

### Pattern: Suggestion → Discussion → Consensus

```
Reviewer: "Consider using pattern X instead"
Author: "I tried that but ran into issue Y"
Reviewer: "Ah, makes sense. Current approach is fine"
Reviewer: *resolves thread*
```

**With gh-talk:**

```bash
gh talk reply PRRT_suggestion "I tried that but ran into issue Y"
# Reviewer reads
gh talk reply PRRT_suggestion "Ah, makes sense. Current approach is fine"
gh talk resolve PRRT_suggestion
```

### Pattern: Multiple Reviewers → Consensus

```
Reviewer1: "Should use approach A"
Reviewer2: "I disagree, approach B is better"
Author: "What if we use approach C which combines both?"
Reviewer1: 👍
Reviewer2: 🎉
Author: *resolves thread*
```

**With gh-talk:**

```bash
gh talk reply PRRT_debate "What if we use approach C which combines both?"
# Reviewers react
# Author sees consensus
gh talk show PRRT_debate
# Sees both reviewers reacted positively
gh talk resolve PRRT_debate
```

## Metrics & Insights

### Developer Productivity Metrics

**Trackable with gh-talk:**

```bash
# Time to first response
gh talk list threads --author @me --json createdAt,comments | \
  jq '.[] | {
    created: .createdAt,
    first_response: .comments[1].createdAt,
    delta: (.comments[1].createdAt - .createdAt)
  }'

# Resolution rate
resolved=$(gh talk list threads --resolved --since 7days | wc -l)
total=$(gh talk list threads --since 7days | wc -l)
echo "Resolution rate: $(($resolved * 100 / $total))%"

# Active threads per PR
gh pr list --json number | \
  jq -r '.[].number' | \
  while read pr; do
    count=$(gh talk list threads --pr $pr --unresolved | wc -l)
    echo "PR #$pr: $count unresolved threads"
  done
```

### Team Health Metrics

**Insights from gh-talk data:**

```bash
# Review response time by team member
gh talk list threads --json author,comments | \
  jq -r 'group_by(.author.login) | .[] | {
    reviewer: .[0].author.login,
    avg_response_time: (map(.comments[1].createdAt - .comments[0].createdAt) | add / length)
  }'

# Most active reviewers
gh talk list threads --json comments | \
  jq -r '[.[].comments[].author.login] | group_by(.) | 
    map({reviewer: .[0], count: length}) | 
    sort_by(.count) | reverse'

# Threads requiring most discussion
gh talk list threads --json id,comments | \
  jq -r 'map({id: .id, comment_count: (.comments | length)}) | 
    sort_by(.comment_count) | reverse | .[0:10]'
```

## Integration with Development Workflows

### Git Hooks Integration

**Pre-push hook:**

```bash
#!/bin/bash
# .git/hooks/pre-push

# Check for unresolved threads before pushing
pr_num=$(gh pr view --json number | jq -r '.number')

if [ -n "$pr_num" ]; then
  unresolved=$(gh talk list threads --pr $pr_num --unresolved --json count)
  
  if [ "$unresolved" -gt 0 ]; then
    echo "⚠️  Warning: $unresolved unresolved review threads"
    echo "Consider addressing feedback before pushing more changes"
    read -p "Continue anyway? (y/n) " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
      exit 1
    fi
  fi
fi
```

### CI/CD Integration

**GitHub Actions workflow:**

```yaml
name: Auto-resolve CI threads
on:
  pull_request:
    types: [synchronize]

jobs:
  auto-resolve:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      
      - name: Install gh-talk
        run: gh extension install hamishmorgan/gh-talk
      
      - name: Check CI threads
        run: |
          # Get threads from CI bot
          threads=$(gh talk list threads --author github-actions --unresolved --json id,body)
          
          # Check if issues are fixed
          echo "$threads" | jq -r '.[].id' | while read thread; do
            # Run check
            if check_is_fixed "$thread"; then
              gh talk resolve "$thread" --message "✅ Fixed in latest commit"
            fi
          done
```

### Editor Integration

**VS Code task:**

```json
{
  "version": "2.0.0",
  "tasks": [
    {
      "label": "Show unresolved threads",
      "type": "shell",
      "command": "gh talk list threads --unresolved",
      "problemMatcher": [],
      "presentation": {
        "reveal": "always",
        "panel": "new"
      }
    },
    {
      "label": "Reply to thread",
      "type": "shell",
      "command": "gh talk reply ${input:threadId} '${input:message}'",
      "problemMatcher": []
    }
  ],
  "inputs": [
    {
      "id": "threadId",
      "type": "promptString",
      "description": "Thread ID"
    },
    {
      "id": "message",
      "type": "promptString",
      "description": "Reply message"
    }
  ]
}
```

## Conclusion

### Summary of Usage Patterns

**Humans Use gh-talk For:**

1. ✅ Quick thread navigation and filtering
2. ✅ Rapid acknowledgments via reactions
3. ✅ Bulk operations on similar feedback
4. ✅ Scripting review workflows
5. ✅ Terminal-native PR management
6. ✅ Reducing context switches

**Bots Use gh-talk For:**

1. ✅ Automated review posting
2. ✅ Thread resolution based on fixes
3. ✅ CI/CD integration
4. ✅ Dependency management automation
5. ✅ Code quality enforcement
6. ✅ Metrics and reporting

**AI Agents Use gh-talk For:**

1. ✅ Context-aware assistance
2. ✅ Automated PR creation with explanations
3. ✅ Learning from review patterns
4. ✅ Proactive code suggestions
5. ✅ Responding to developer questions
6. ✅ Pattern recognition and improvement

### Key Differentiators

**What Makes gh-talk Valuable:**

- **Speed**: 60-90% time savings on common tasks
- **Automation**: Scriptable for bots and workflows
- **Bulk Operations**: Process many threads at once
- **Terminal Native**: Never leave the command line
- **AI-Friendly**: Designed for agent integration
- **Team Workflows**: Supports collaboration patterns

### Success Criteria

**gh-talk succeeds when:**

- Developers spend 80% less time in browser for reviews
- Review response times decrease by 50%
- Thread resolution rates increase
- Bots can fully automate review workflows
- AI agents can participate naturally in code review
- Teams adopt it as primary review tool

---

**Last Updated**: 2025-11-02  
**Context**: Real-world usage patterns and workflows for gh-talk  
**Sources**: Developer workflows, bot patterns, AI agent needs, best practices
