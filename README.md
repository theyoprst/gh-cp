# GitHub Pull Request Cherry-Pick CLI

A CLI utility written in Go that allows you to cherry-pick entire GitHub pull requests to a destination branch. GitHub doesn't have this functionality ready-to-use, so this project solves that problem by automating the cherry-pick process while preserving the original PR's metadata.

## Prerequisites

- Go (for building the utility)
- GitHub CLI (`gh`) - must be pre-configured and authenticated
- Git repository with proper access permissions

## Installation

### Quick Install (Recommended)

```bash
go install github.com/theyoprst/gh-cp/cmd/gh-cp@latest
```

This will download, build, and install the `gh-cp` binary to your `$(go env GOPATH)/bin` directory.

### Alternative: Build from Source

```bash
# Clone the repository
git clone https://github.com/theyoprst/gh-cp.git
cd gh-cp

# Build and install
make install
```

### Verify Installation

```bash
# Check if gh-cp is available
gh-cp --help
```

**Note**: Make sure your `$(go env GOPATH)/bin` is in your system's `$PATH` environment variable.

## Upgrading

The standard `go install github.com/theyoprst/gh-cp/cmd/gh-cp@latest` command won't upgrade if you already have the same version installed. Use one of these methods to force an upgrade:

### Recommended: Force Rebuild
```bash
go install -a github.com/theyoprst/gh-cp/cmd/gh-cp@latest
```

### Alternative: Clear Build Cache
```bash
go clean -cache && go install github.com/theyoprst/gh-cp/cmd/gh-cp@latest
```

### Nuclear Option: Clear Module Cache
```bash
go clean -modcache && go install github.com/theyoprst/gh-cp/cmd/gh-cp@latest
```

**Note**: The `-modcache` option removes all cached Go modules and will cause slower builds for all Go projects until modules are re-downloaded.

### Third-Party Tools
For users managing many Go binaries, consider tools like [gup](https://github.com/nao1215/gup) which can update all Go binaries automatically:
```bash
# Install gup
go install github.com/nao1215/gup@latest

# Update all Go binaries
gup update
```

## Usage

```bash
gh-cp <pull-request-number> <target-branch> [--dry-run]
```

**Parameters:**
- `pull-request-number` - The number of the merged pull request to cherry-pick
- `target-branch` - The destination branch to cherry-pick the changes to (can be `branch` or `remote/branch` format)
- `--dry-run` - (Optional) Preview mode: shows what would be done without making any remote changes
- `--skip-merged-check` - (Optional) Skip check that PR is merged (use with caution when cherry-picking unmerged PRs)

**Examples:**
```bash
# Cherry-pick PR #1319 to release/v2.1 branch
gh-cp 1319 release/v2.1

# Cherry-pick PR #1319 to origin/release/v2.1 branch
gh-cp 1319 origin/release/v2.1

# Preview the cherry-pick operation without making changes
gh-cp 1319 release/v2.1 --dry-run

# Cherry-pick an unmerged PR (use with caution)
gh-cp 1319 release/v2.1 --skip-merged-check
```

## How It Works

1. **Fetches PR Information**: Uses `gh pr view 1319 --json "number,title,body,state,baseRefName,mergeCommit,commits,labels"` to retrieve all information about the pull request

2. **Validates PR State**: Ensures the PR is merged (cherry-picking unmerged changes is not supported)

3. **Fetches Remote State**: Automatically fetches the latest state of the target branch from the remote

4. **Creates Isolated Worktree**: Creates a temporary git worktree to perform operations without affecting your current working directory, allowing cherry-picking even with uncommitted changes

5. **Creates New Branch**: Creates a new branch with the naming pattern `cherry-pick-to/target-branch/from/original-branch`. If the branch already exists, adds incremental suffixes (`/0`, `/1`, etc.)

6. **Cherry-picks Changes**: Applies all commits from the original PR to the new branch in the isolated worktree

7. **Creates New PR**: Pushes changes and creates a new pull request targeting the specified branch with:
   - Title prefixed with `[cherry-pick]`
   - Body prefixed with a reference message linking to the original PR
   - All original labels copied over
   - Same description as the original PR
   - Base branch set to the target branch you specified

## Current Limitations

- **Conflict Resolution**: The tool does not automatically resolve merge conflicts. If conflicts occur during cherry-picking, the tool will leave the worktree in place and provide instructions for manual resolution
- **Merged PRs Only**: By default, only works with merged pull requests (use `--skip-merged-check` to override)

## Conflict Resolution

When cherry-pick conflicts occur, the tool will:
1. Leave the temporary worktree in place (usually in `/tmp/gh-cp-worktree-*`)
2. Display the worktree path and provide step-by-step manual resolution instructions
3. Exit with an error message containing all necessary commands to complete the process

## Example Workflow

```bash
# Cherry-pick PR #1319 to the release/v2.1 branch
gh-cp 1319 release/v2.1
```

This will:
1. Fetch information about PR #1319
2. Create a temporary worktree to isolate operations
3. Create a unique branch name like `cherry-pick-to/release-v2.1/from/feature-branch-name/0` (incrementing the suffix if needed)
4. Cherry-pick all commits from PR #1319 in the isolated worktree
5. Push the new branch (with force to handle conflicts)
6. Create a new PR from the cherry-pick branch → `release/v2.1` titled `[cherry-pick] Original PR Title` with a link back to PR #1319
7. Clean up the temporary worktree

## Dry-Run Mode

Use `--dry-run` to preview operations without making any remote changes:

```bash
gh-cp 1319 release/v2.1 --dry-run
```

In dry-run mode, the tool will:
- ✅ Fetch PR information and validate it's merged
- ✅ Create temporary worktree and cherry-pick commits
- 🔍 Show push command that would be executed (from worktree)
- 🔍 Show PR creation command that would be executed
- ✅ Clean up temporary worktree (your working directory remains untouched)

## Development

### Building from Source

```bash
# Build the binary
make build

# Clean build artifacts
make clean
```

### Code Quality

```bash
# Run linters
make lint

# Auto-fix linting issues where possible
make lint-fix
```

### Testing

```bash
# Run tests
make test
```

### Installation

```bash
# Build and install to $(go env GOPATH)/bin
make install
```

### Project Structure

- `cmd/gh-cp/` - Main application entry point
- `internal/github/` - GitHub API integration
- `internal/git/` - Git operations
- `internal/cherry/` - Core cherry-pick logic