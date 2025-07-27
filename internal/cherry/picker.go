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

	branchName := git.GenerateCherryPickBranchName(prData.BaseRefName, targetBranch, prData.Number)
	fmt.Printf("✓ Generated branch name: %s\n", branchName)

	originalBranch, err := git.GetCurrentBranch()
	if err != nil {
		return fmt.Errorf("get current branch: %w", err)
	}

	if err := git.CheckoutTargetBranch(targetBranch); err != nil {
		return err
	}
	fmt.Printf("✓ Checked out target branch: %s\n", targetBranch)

	if err := git.CreateAndCheckoutBranch(branchName); err != nil {
		return err
	}
	fmt.Printf("✓ Created branch: %s\n", branchName)

	if err := git.CherryPickCommits(commitSHAs); err != nil {
		if err := git.CheckoutTargetBranch(originalBranch); err != nil {
			fmt.Printf("⚠️ Failed to switch back to original branch %s: %v\n", originalBranch, err)
		}
		if err := git.DeleteBranch(branchName); err != nil {
			fmt.Printf("⚠️ Failed to delete branch %s: %v\n", branchName, err)
		}
		return fmt.Errorf("cherry-pick failed: %w", err)
	}
	fmt.Printf("✓ Cherry-picked %d commits successfully\n", len(commitSHAs))

	if err := git.PushBranch(branchName, config.DryRun); err != nil {
		return err
	}

	prURL, err := CreatePR(prData, targetBranch, config.DryRun)
	if err != nil {
		return err
	}
	fmt.Printf("✓ Created PR: %s\n", prURL)

	if config.DryRun {
		fmt.Printf("✓ Cleaning up dry-run: switching back to %s and deleting %s\n", originalBranch, branchName)
		if err := git.CheckoutTargetBranch(originalBranch); err != nil {
			fmt.Printf("⚠️ Failed to switch back to original branch %s: %v\n", originalBranch, err)
		}
		if err := git.DeleteBranch(branchName); err != nil {
			fmt.Printf("⚠️ Failed to delete dry-run branch %s: %v\n", branchName, err)
		}
		fmt.Println("✓ Dry run completed successfully")
	} else {
		fmt.Println("✓ Cherry-pick completed successfully")
	}

	return nil
}