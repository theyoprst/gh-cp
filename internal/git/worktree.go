package git

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func CreateWorktree(branchName, targetBranch string) (string, error) {
	worktreePath, err := GetUniqueWorktreePath()
	if err != nil {
		return "", err
	}
	
	cmd := exec.Command("git", "worktree", "add", worktreePath, targetBranch)
	if err := cmd.Run(); err != nil {
		// Clean up the created directory if worktree creation fails
		if removeErr := os.RemoveAll(worktreePath); removeErr != nil {
			return "", fmt.Errorf("create worktree: %w (cleanup failed: %w)", err, removeErr)
		}
		return "", fmt.Errorf("create worktree: %w", err)
	}
	
	return worktreePath, nil
}

func RemoveWorktree(worktreePath string) error {
	cmd := exec.Command("git", "worktree", "remove", worktreePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("remove worktree: %w", err)
	}
	return nil
}

func GetUniqueWorktreePath() (string, error) {
	worktreePath, err := os.MkdirTemp("", "gh-cp-worktree-")
	if err != nil {
		return "", fmt.Errorf("create temp directory: %w", err)
	}
	return worktreePath, nil
}


func CheckBranchExists(branchName string) (bool, error) {
	cmd := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branchName)
	err := cmd.Run()
	if err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) && exitError.ExitCode() == 1 {
			return false, nil
		}
		return false, fmt.Errorf("check branch exists: %w", err)
	}
	return true, nil
}

func IsWorktreeClean(worktreePath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = worktreePath
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("check worktree status: %w", err)
	}
	return len(strings.TrimSpace(string(output))) == 0, nil
}