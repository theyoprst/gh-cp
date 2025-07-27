// Package github provides GitHub API integration via the gh CLI.
package github

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func FetchPRData(prNumber int) (*PRData, error) {
	cmd := exec.Command("gh", "pr", "view", strconv.Itoa(prNumber), "--json", "number,title,body,state,baseRefName,mergeCommit,commits,labels")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("fetch PR data: %w", err)
	}

	var prData PRData
	if err := json.Unmarshal(output, &prData); err != nil {
		return nil, fmt.Errorf("parse PR data: %w", err)
	}

	return &prData, nil
}

func ValidatePRMerged(prData *PRData) error {
	if strings.ToLower(prData.State) != "merged" {
		return fmt.Errorf("PR #%d is not merged (state: %s)", prData.Number, prData.State)
	}

	if prData.MergeCommit == nil {
		return fmt.Errorf("PR #%d has no merge commit", prData.Number)
	}

	return nil
}

func GetCommitSHAs(prData *PRData) []string {
	shas := make([]string, len(prData.Commits))
	for i, commit := range prData.Commits {
		shas[i] = commit.SHA
	}
	return shas
}

func FormatLabels(labels []Label) []string {
	labelNames := make([]string, len(labels))
	for i, label := range labels {
		labelNames[i] = label.Name
	}
	return labelNames
}
