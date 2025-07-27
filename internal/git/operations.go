package git

import (
	"fmt"
	"os/exec"
	"strings"
)

func CheckoutTargetBranch(targetBranch string) error {
	return CheckoutTargetBranchInDir(targetBranch, "")
}

func CheckoutTargetBranchInDir(targetBranch, workingDir string) error {
	cmd := exec.Command("git", "checkout", targetBranch)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("checkout target branch %s: %w", targetBranch, err)
	}
	return nil
}

func CreateAndCheckoutBranch(branchName string) error {
	return CreateAndCheckoutBranchInDir(branchName, "")
}

func CreateAndCheckoutBranchInDir(branchName, workingDir string) error {
	cmd := exec.Command("git", "checkout", "-B", branchName)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("create and checkout branch %s: %w", branchName, err)
	}
	return nil
}

func CherryPickCommits(commitSHAs []string) error {
	return CherryPickCommitsInDir(commitSHAs, "")
}

func CherryPickCommitsInDir(commitSHAs []string, workingDir string) error {
	for _, sha := range commitSHAs {
		cmd := exec.Command("git", "cherry-pick", "-x", sha)
		if workingDir != "" {
			cmd.Dir = workingDir
		}
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("cherry-pick commit %s: %w", sha, err)
		}
	}
	return nil
}

func PushBranch(branchName string, dryRun bool) error {
	return PushBranchFromDir(branchName, dryRun, "")
}

func PushBranchFromDir(branchName string, dryRun bool, workingDir string) error {
	pushCmd := fmt.Sprintf("git push --force -u origin %s", branchName)

	if dryRun {
		if workingDir != "" {
			fmt.Printf("[DRY RUN] Would execute in %s: %s\n", workingDir, pushCmd)
		} else {
			fmt.Printf("[DRY RUN] Would execute: %s\n", pushCmd)
		}
		return nil
	}

	cmd := exec.Command("git", "push", "--force", "-u", "origin", branchName)
	if workingDir != "" {
		cmd.Dir = workingDir
	}
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("push branch %s with command '%s': %w\nOutput: %s", branchName, pushCmd, err, string(output))
	}
	return nil
}

func DeleteBranch(branchName string) error {
	cmd := exec.Command("git", "branch", "-D", branchName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("delete branch %s: %w\nOutput: %s", branchName, err, string(output))
	}
	return nil
}

func GetCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}
	return strings.TrimSpace(string(output)), nil
}

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}
