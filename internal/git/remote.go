package git

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
)

// getRemotes returns a list of all configured git remotes.
func getRemotes() ([]string, error) {
	cmd := exec.Command("git", "remote")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("get remotes: %w", err)
	}

	remotesStr := strings.TrimSpace(string(output))
	if remotesStr == "" {
		return []string{}, nil
	}

	return strings.Split(remotesStr, "\n"), nil
}

// getDefaultRemote returns the first available remote.
func getDefaultRemote() (string, error) {
	remotes, err := getRemotes()
	if err != nil {
		return "", fmt.Errorf("get remotes: %w", err)
	}

	if len(remotes) == 0 {
		return "", fmt.Errorf("no git remotes configured")
	}

	return remotes[0], nil
}

// FetchRemoteBranch executes git fetch for the specified remote and branch.
func FetchRemoteBranch(remote, branch string) error {
	args := []string{"git", "fetch", remote, branch}
	cmd := exec.Command(args[0], args[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("fetch remote %s branch %s with command '%s': %w\nOutput: %s", remote, branch, strings.Join(args, " "), err, string(output))
	}
	return nil
}

// ParseRemoteAndBranch parses a target branch string and returns the remote and branch names.
// For "remote/branch" format, it validates the remote exists and returns (remote, branch).
// For plain "branch" format, it returns ("origin", branch) as default.
// Returns error if the remote doesn't exist or if the format is ambiguous.
func ParseRemoteAndBranch(targetBranch string) (remote, branch string, err error) {
	parts := strings.SplitN(targetBranch, "/", 2)

	// If no slash, use default remote
	if len(parts) == 1 {
		defaultRemote, err := getDefaultRemote()
		if err != nil {
			return "", "", fmt.Errorf("get default remote: %w", err)
		}
		return defaultRemote, targetBranch, nil
	}

	// If slash exists, check if the first part is a remote
	candidateRemote := parts[0]
	candidateBranch := parts[1]

	remotes, err := getRemotes()
	if err != nil {
		return "", "", fmt.Errorf("get remotes: %w", err)
	}

	if slices.Contains(remotes, candidateRemote) {
		// It's a remote/branch format
		return candidateRemote, candidateBranch, nil
	}

	// It's a branch name with slash, use default remote
	defaultRemote, err := getDefaultRemote()
	if err != nil {
		return "", "", fmt.Errorf("get default remote: %w", err)
	}

	return defaultRemote, targetBranch, nil
}
