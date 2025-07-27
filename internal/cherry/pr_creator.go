package cherry

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/theyoprst/gh-cp/internal/github"
)

// CreatePR creates a new GitHub pull request with the cherry-picked changes.
func CreatePR(prData *github.PRData, targetBranch string, dryRun bool) (prURL string, err error) {
	title := fmt.Sprintf("[cherry-pick] %s", prData.Title)

	body := fmt.Sprintf(`Cherry-picked from #%d

%s`, prData.Number, prData.Body)

	labels := github.FormatLabels(prData.Labels)

	args := []string{"pr", "create", "--title", title, "--body", body, "--base", targetBranch}

	if len(labels) > 0 {
		args = append(args, "--label", strings.Join(labels, ","))
	}

	if dryRun {
		cmdStr := fmt.Sprintf("gh %s", strings.Join(args, " "))
		fmt.Printf("[DRY RUN] Would execute: %s\n", cmdStr)
		return "", nil
	}

	cmd := exec.Command("gh", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		cmdStr := fmt.Sprintf("gh %s", strings.Join(args, " "))
		return "", fmt.Errorf("create PR with command '%s': %w\nOutput: %s", cmdStr, err, string(output))
	}

	return strings.TrimSpace(string(output)), nil
}