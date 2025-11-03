# pr-merge

@.cursor/rules/pr-merge.mdc

Final verification and merge workflow for PRs.

**Usage:**
- `/pr-merge <PR>` - Verify CI passes, then merge (squash and delete branch)

**Behavior:**
- Verifies all CI checks pass (exit code 0)
- Posts status update
- Merges PR with `--squash --delete-branch`
- Updates local main branch

**When invoked:** Your use of `/pr-merge` IS your instruction to merge

