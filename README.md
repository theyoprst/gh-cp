# GitHub Pull Request Cherry-Pick CLI

A CLI utility written in Go that allows you to cherry-pick entire GitHub pull requests to a destination branch. GitHub doesn't have this functionality ready-to-use, so this project solves that problem by automating the cherry-pick process while preserving the original PR's metadata.

## Prerequisites

- Go (for building the utility)
- GitHub CLI (`gh`) - must be pre-configured and authenticated
- Git repository with proper access permissions

## Usage

```bash
gh-cp <pull-request-number> <target-branch>
```

**Parameters:**
- `pull-request-number` - The number of the merged pull request to cherry-pick
- `target-branch` - The destination branch to cherry-pick the changes to

**Example:**
```bash
gh-cp 1319 release/v2.1
```

## How It Works

1. **Fetches PR Information**: Uses `gh pr view 1319 --json "mergeCommit,commits,state,title,labels,body"` to retrieve all information about the pull request

2. **Validates PR State**: Ensures the PR is merged (cherry-picking unmerged changes is not supported)

3. **Creates New Branch**: Checks out to the target branch and creates a new branch with the naming pattern `cherry-pick/original-pr-branch-name`

4. **Cherry-picks Changes**: Applies all commits from the original PR to the new branch

5. **Creates New PR**: Pushes changes and creates a new pull request with:
   - Title prefixed with `[cherry-pick]`
   - Body prefixed with a reference message linking to the original PR
   - All original labels copied over
   - Same description as the original PR

## Current Limitations

- **Conflict Resolution**: The tool does not automatically resolve merge conflicts. If conflicts occur during cherry-picking, the process will stop and require manual intervention
- **Merged PRs Only**: Only works with merged pull requests

## Example Workflow

```bash
# Cherry-pick PR #1319 to the release/v2.1 branch
gh-cp 1319 release/v2.1
```

This will:
1. Fetch information about PR #1319
2. Create a new branch `cherry-pick/feature-branch-name` based on `release/v2.1`
3. Cherry-pick all commits from PR #1319
4. Push the new branch
5. Create a new PR titled `[cherry-pick] Original PR Title` with a link back to PR #1319