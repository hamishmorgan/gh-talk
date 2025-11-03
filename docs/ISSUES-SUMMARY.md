# GitHub Issues Summary

**Created**: 2025-11-03  
**Source**: Real-world user feedback from production usage  
**Total Issues Created**: 12

## Issues Created from USER-FEEDBACK.md

### High Priority (Immediate Value)

#### Issue #4: JSON Output Format ✅ IMPLEMENTED

**Status**: Closed - Implemented in commit afa5b0e  
**Implementation**: Added proper JSON output with all fields, comment IDs included  
**Test**: `gh talk list threads --format json` works perfectly

#### Issue #5: Bulk Operations for Hide and React

**Status**: Open  
**Priority**: High  
**Description**: Accept multiple IDs or filter-based bulk operations  
**Example**: `gh talk hide --threads resolved --reason resolved`  
**Impact**: Would eliminate 80% of repetitive commands

#### Issue #7: Status/Summary Command

**Status**: Open  
**Priority**: High  
**Description**: Overview command showing PR review progress  
**Example**: `gh talk status --pr 137` shows threads resolved, comments hidden, etc.  
**Impact**: Verification without leaving terminal

### Medium Priority (Quality of Life)

#### Issue #6: Improve Show Command ✅ IMPLEMENTED  

**Status**: Closed - Implemented in commit afa5b0e  
**Implementation**: Comment IDs now prominently displayed with numbers  
**Test**: `gh talk show PRRT_xxx` shows `[1] PRRC_xxx` format

#### Issue #8: Add --react Flag to Reply Command

**Status**: Open  
**Priority**: Medium  
**Description**: Combine reply + react + resolve in one command  
**Example**: `gh talk reply PRRT_xxx "Done" --resolve --react 👍`  
**Impact**: Reduces command count for common workflow

#### Issue #9: Support Numeric Database IDs

**Status**: Open  
**Priority**: Medium  
**Description**: Auto-convert numeric IDs to node IDs  
**Challenge**: Requires API lookup to convert  
**Alternative**: Better error message (Issue #14)

#### Issue #10: Cleanup Workflow Command

**Status**: Open  
**Priority**: Medium  
**Description**: Single command to hide all resolved thread comments  
**Example**: `gh talk cleanup --pr 137`  
**Impact**: Common end-of-review workflow

#### Issue #11: Documentation - Common Workflows and Recipes

**Status**: Open  
**Priority**: Medium  
**Type**: Documentation  
**Description**: Add docs/RECIPES.md with common patterns  
**Content**: Address all feedback, clean up conversations, troubleshooting

#### Issue #14: Better Error Messages for Numeric IDs

**Status**: Open  
**Priority**: Medium  
**Description**: Improve error message with conversion hint  
**Impact**: Helps users self-correct without being stuck

### Low Priority (Nice to Have)

#### Issue #12: Interactive TUI Mode

**Status**: Open  
**Priority**: Low (Phase 3)  
**Description**: Full TUI with Bubble Tea for reviewing threads  
**Scope**: Major feature for Phase 3  
**Related**: Already planned in SPEC.md

#### Issue #13: Template Support for Replies

**Status**: Open  
**Priority**: Low  
**Description**: Template system for common reply messages  
**Example**: `gh talk reply --template fixed --commit abc123`  
**Use Case**: Teams with standard responses

#### Issue #15: Dry-Run Mode

**Status**: Open  
**Priority**: Low  
**Description**: Preview bulk operations before executing  
**Example**: `gh talk hide --threads resolved --dry-run`  
**Impact**: Safety for bulk operations

## Implementation Priority

### Already Implemented ✅ (Closed)

1. ✅ Issue #4: JSON output
2. ✅ Issue #6: Show command with comment IDs

### Next to Implement (High Priority)

3. ⏭️ Issue #5: Bulk operations (hide, react with multiple IDs)
4. ⏭️ Issue #7: Status command
5. ⏭️ Issue #8: --react flag on reply

### Short-term (Medium Priority)

6. ⏭️ Issue #14: Better error messages
7. ⏭️ Issue #11: Documentation/recipes
8. ⏭️ Issue #10: Cleanup command

### Long-term (Low Priority / Phase 3)

9. ⏭️ Issue #9: Numeric ID conversion (complex)
10. ⏭️ Issue #12: Interactive TUI mode (Phase 3)
11. ⏭️ Issue #13: Template support
12. ⏭️ Issue #15: Dry-run mode

## Feature Gaps Identified

### Already Implemented But Not Discovered

- `--unresolved` flag (exists in code)
- `--author` filter (exists in code)
- `--file` filter (exists in code)

**Action**: Ensure these are well-documented and visible in help text.

### Truly Missing (From Feedback)

- Bulk operations with filters
- JSON output (NOW IMPLEMENTED)
- Status/summary command
- --react on reply command
- Comment ID display (NOW IMPLEMENTED)
- Workflow commands (cleanup, acknowledge)

## User Satisfaction

**Overall Rating**: Excellent 🌟  
**Success Rate**: 100% (all commands worked)  
**Time Saved**: 10-15 minutes vs web UI  
**Would Use Again**: Absolutely yes  
**Would Recommend**: Yes, enthusiastically

### Most Valuable Features (Per User)

1. Reply + Resolve in one command
2. Thread listing with clear overview
3. Emoji reactions
4. Comment hiding
5. Clear error messages

### Biggest Opportunities (Per User)

1. Bulk operations (now Issue #5)
2. Integrated verification (now Issue #7)
3. Workflow shortcuts (now Issue #10)
4. JSON output (NOW IMPLEMENTED ✅)

## Implementation Plan

### Immediate (This Session)

- ✅ JSON output - DONE
- ✅ Show command improvements - DONE
- ⏭️ Bulk operations - NEXT

### This Week

- Issue #5: Bulk operations
- Issue #7: Status command  
- Issue #8: --react flag

### This Month

- Issue #10: Cleanup command
- Issue #11: Documentation
- Issue #14: Error messages

### Phase 3

- Issue #12: TUI mode
- Issue #13: Templates
- Issue #15: Dry-run

---

**Last Updated**: 2025-11-03  
**Issues Created**: 12  
**Issues Closed**: 2  
**Issues Open**: 10  
**User Feedback**: Highly positive with clear improvement path
