package git

import (
	"fmt"
	"strings"
)

func GenerateCherryPickBranchName(originalBranchName, targetBranch string, prNumber int) string {
	if originalBranchName != "" {
		cleanOriginal := strings.ReplaceAll(originalBranchName, "/", "-")
		cleanTarget := strings.ReplaceAll(targetBranch, "/", "-")
		return fmt.Sprintf("cherry-pick-to/%s/from/%s", cleanTarget, cleanOriginal)
	}
	
	return fmt.Sprintf("cherry-pick-pr-%d", prNumber)
}