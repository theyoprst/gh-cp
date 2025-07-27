// Package git provides git operations for the cherry-pick utility.
package git

import (
	"fmt"
	"strings"
)

// GenerateUniqueBranchName creates a unique branch name by checking for conflicts and adding suffixes.
func GenerateUniqueBranchName(originalBranchName, targetBranch string, prNumber int) (string, error) {
	baseName := generateCherryPickBranchBaseName(originalBranchName, targetBranch, prNumber)

	suffix := 0
	for {
		candidateName := fmt.Sprintf("%s/%d", baseName, suffix)
		exists, err := CheckBranchExists(candidateName)
		if err != nil {
			return "", fmt.Errorf("check branch exists: %w", err)
		}
		if !exists {
			return candidateName, nil
		}
		suffix++
		if suffix > 100 {
			return "", fmt.Errorf("too many branch name conflicts")
		}
	}
}

func generateCherryPickBranchBaseName(originalBranchName, targetBranch string, prNumber int) string {
	if originalBranchName != "" {
		cleanOriginal := strings.ReplaceAll(originalBranchName, "/", "-")
		cleanTarget := strings.ReplaceAll(targetBranch, "/", "-")
		return fmt.Sprintf("cherry-pick-to/%s/from/%s", cleanTarget, cleanOriginal)
	}

	return fmt.Sprintf("cherry-pick-pr-%d", prNumber)
}
