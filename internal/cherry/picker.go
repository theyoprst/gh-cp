// Package cherry provides the main cherry-pick orchestration logic.
package cherry

import (
	"fmt"

	"github.com/theyoprst/gh-cp/internal/git"
	"github.com/theyoprst/gh-cp/internal/github"
)

// CherryPickPR orchestrates the complete cherry-pick workflow for a GitHub PR.
func CherryPickPR(prNumber int, targetBranch string, config *github.Config) error {
	if !git.IsGitRepo() {
		return fmt.Errorf("not in a git repository")
	}

	fmt.Printf("✓ Fetching PR #%d...\n", prNumber)
	prData, err := github.FetchPRData(prNumber)
	if err != nil {
		return fmt.Errorf("fetch PR data: %w", err)
	}

	fmt.Printf("✓ Fetched PR #%d: \"%s\"\n", prData.Number, prData.Title)

	if err := github.ValidatePRMerged(prData); err != nil {
		return err
	}
	fmt.Printf("✓ Validated PR is merged\n")

	commitSHAs := github.GetCommitSHAs(prData)
	if len(commitSHAs) == 0 {
		return fmt.Errorf("PR #%d has no commits", prNumber)
	}

	branchName, err := git.GenerateUniqueBranchName(prData.BaseRefName, targetBranch, prData.Number)
	if err != nil {
		return fmt.Errorf("generate unique branch name: %w", err)
	}
	fmt.Printf("✓ Generated unique branch name: %s\n", branchName)

	fmt.Printf("✓ Creating worktree for isolated operations...\n")
	worktreePath, err := git.CreateWorktree(branchName, targetBranch)
	if err != nil {
		return fmt.Errorf("create worktree: %w", err)
	}
	fmt.Printf("✓ Created worktree at: %s\n", worktreePath)

	defer func() {
		fmt.Printf("✓ Cleaning up worktree...\n")
		if err := git.RemoveWorktree(worktreePath); err != nil {
			fmt.Printf("⚠️ Remove worktree %s: %v\n", worktreePath, err)
		}

		fmt.Printf("✓ Cleaning up branch: %s\n", branchName)
		if err := git.DeleteBranch(branchName); err != nil {
			fmt.Printf("⚠️ Delete branch %s: %v\n", branchName, err)
		}
	}()

	if err := git.CreateAndCheckoutBranchInDir(branchName, worktreePath); err != nil {
		return fmt.Errorf("create branch in worktree: %w", err)
	}
	fmt.Printf("✓ Created and checked out branch: %s\n", branchName)

	if err := git.CherryPickCommitsInDir(commitSHAs, worktreePath); err != nil {
		if isClean, checkErr := git.IsWorktreeClean(worktreePath); checkErr == nil && !isClean {
			fmt.Printf("⚠️ Cherry-pick failed with conflicts. Manual resolution required in: %s\n", worktreePath)
			fmt.Printf("After resolving conflicts:\n")
			fmt.Printf("  1. cd %s\n", worktreePath)
			fmt.Printf("  2. git add . && git cherry-pick --continue (repeat for each commit)\n")
			fmt.Printf("  3. git push --force -u origin %s\n", branchName)
			fmt.Printf("  4. Create PR manually after completing cherry-pick\n")
			fmt.Printf("  5. Remove temporary directory: rm -rf %s\n", worktreePath)
			return fmt.Errorf("cherry-pick conflicts require manual resolution")
		}
		return fmt.Errorf("cherry-pick: %w", err)
	}
	fmt.Printf("✓ Cherry-picked %d commits successfully\n", len(commitSHAs))

	if err := git.PushBranchFromDir(branchName, config.DryRun, worktreePath); err != nil {
		return err
	}

	prURL, err := CreatePR(prData, targetBranch, branchName, config.DryRun)
	if err != nil {
		return err
	}
	if prURL != "" {
		fmt.Printf("✓ Created PR: %s\n", prURL)
	}

	if config.DryRun {
		fmt.Println("✓ Dry run completed successfully")
	} else {
		fmt.Println("✓ Cherry-pick completed successfully")
	}

	return nil
}
