// Package git provides git operations for the cherry-pick utility.
package git

import (
	"fmt"
	"strings"
)

// GenerateCherryPickBranchName creates a descriptive branch name for cherry-pick operations.
func GenerateCherryPickBranchName(originalBranchName, targetBranch string, prNumber int) string {
	if originalBranchName != "" {
		cleanOriginal := strings.ReplaceAll(originalBranchName, "/", "-")
		cleanTarget := strings.ReplaceAll(targetBranch, "/", "-")
		return fmt.Sprintf("cherry-pick-to/%s/from/%s", cleanTarget, cleanOriginal)
	}

	return fmt.Sprintf("cherry-pick-pr-%d", prNumber)
}
