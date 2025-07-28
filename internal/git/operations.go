package git

import (
	"fmt"
	"os/exec"
)

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

func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	return cmd.Run() == nil
}
