# GH-CP Developer Guide

## Project Overview

This is a CLI utility written in Go that cherry-picks entire GitHub pull requests to destination branches, preserving original PR metadata. The tool integrates with GitHub CLI (`gh`) and git commands to automate the cherry-pick workflow.

## Essential Commands

### Building and Testing
```bash
# Build the binary
make build

# Run linters (critical - always run before commits)
make lint

# Auto-fix linting issues
make lint-fix

# Run tests
make test

# Clean build artifacts
make clean

# Install to GOPATH/bin
make install
```

### Development Workflow
```bash
# Test in dry-run mode (always test this first)
./gh-cp <pr-number> <target-branch> --dry-run

# Real execution after dry-run verification
./gh-cp <pr-number> <target-branch>
```

## Testing

### Test Data
For development and testing, use these reference values:
- **Test PR**: #10 "Add CLAUDE.md developer guide"
- **Target branch**: `before-claude-md`

### Dry-run Testing
```bash
# Basic dry-run test
./gh-cp 10 before-claude-md --dry-run

# Expected output:
# ✓ Generated unique branch name: cherry-pick-to/before-claude-md/from/main/0
# ✓ Created worktree at: /tmp/gh-cp-worktree-<random>
# ✓ Cherry-picked 1 commits successfully
# [DRY RUN] Would execute: git push --force -u origin ...
# [DRY RUN] Would execute: gh pr create ...
# ✓ Cleaning up worktree...
# ✓ Cleaning up branch: cherry-pick-to/before-claude-md/from/main/0
```

### Branch Naming Verification
Multiple runs should create incremental suffixes:
```bash
# First run (if no conflicts)
./gh-cp 10 before-claude-md --dry-run
# → cherry-pick-to/before-claude-md/from/main/0

# Second run (increments suffix)
./gh-cp 10 before-claude-md --dry-run
# → cherry-pick-to/before-claude-md/from/main/1
```

### Cleanup Verification
After each run, verify complete cleanup:
```bash
# Check no cherry-pick branches remain
git branch | grep cherry-pick
# Should return nothing

# Check no temporary worktrees remain
ls /tmp | grep gh-cp-worktree
# Should return nothing
```

## Architecture

### Package Structure
- **`cmd/gh-cp/main.go`**: CLI entry point with argument parsing
- **`internal/cherry/`**: Core orchestration logic
  - `picker.go`: Main workflow coordinator (`CherryPickPR` function)
  - `pr_creator.go`: GitHub PR creation with metadata preservation
- **`internal/github/`**: GitHub API integration via `gh` CLI
  - `client.go`: PR data fetching and validation
  - `types.go`: GitHub API response structures
- **`internal/git/`**: Git operations wrapper
  - `operations.go`: Git commands (checkout, cherry-pick, push)
  - `branch.go`: Branch naming logic (`cherry-pick-to/target/from/source`)

### Data Flow
1. **Main** parses CLI args and calls `cherry.CherryPickPR()`
2. **Cherry picker** orchestrates the entire workflow:
   - Fetches PR data via GitHub client
   - Validates PR is merged
   - Generates branch name using original branch from `baseRefName`
   - Executes git operations (checkout target, create branch, cherry-pick)
   - Pushes branch and creates new PR with preserved metadata
3. **Dry-run mode**: Executes all local operations but shows what remote commands would run instead of executing them

### Key Design Patterns
- **External tool integration**: Uses `exec.Command` to call `gh` and `git` CLIs rather than direct API/library calls
- **Error handling with cleanup**: Failed operations attempt to restore original git state
- **Metadata preservation**: New PRs include `[cherry-pick]` prefix, original PR reference, and copy labels
- **Branch naming**: Uses format `cherry-pick-to/{target}/from/{original}` with slash-to-hyphen replacement

### Critical Implementation Details
- Cherry-pick uses `-x` flag for commit attribution
- Dry-run mode creates local branches but cleans them up after showing commands

### Dependencies
- **External**: GitHub CLI (`gh`) must be installed and authenticated
- **Runtime**: Must be run from within a git repository
- **Go modules**: Pure Go with no external dependencies in go.mod

## Go Style Guide
- When adding context to returned errors, keep the context succinct by avoiding phrases like "failed to", which state the obvious and pile up as the error percolates up through the stack
- Do not add obvious comments in the code
