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