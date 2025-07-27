// Package main implements the gh-cp CLI tool for cherry-picking GitHub pull requests.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/theyoprst/gh-cp/internal/cherry"
	"github.com/theyoprst/gh-cp/internal/github"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: %s <pull-request-number> <target-branch> [--dry-run]\n", os.Args[0])
		os.Exit(1)
	}

	prNumber, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Printf("Error: invalid pull request number: %s\n", os.Args[1])
		os.Exit(1)
	}

	targetBranch := os.Args[2]

	config := &github.Config{
		DryRun: len(os.Args) > 3 && os.Args[3] == "--dry-run",
	}
	if err := cherry.CherryPickPR(prNumber, targetBranch, config); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
