// Package main implements the gh-cp CLI tool for cherry-picking GitHub pull requests.
package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/theyoprst/gh-cp/internal/cherry"
	"github.com/theyoprst/gh-cp/internal/github"
	"github.com/theyoprst/gh-cp/internal/version"
)

var (
	dryRun          bool
	skipMergedCheck bool
)

var rootCmd = &cobra.Command{
	Use:   "gh-cp <pull-request-number> <target-branch>",
	Short: "Cherry-pick GitHub pull requests to destination branches",
	Long: `gh-cp is a CLI utility that cherry-picks entire GitHub pull requests
to destination branches, preserving original PR metadata.

The tool integrates with GitHub CLI (gh) and git commands to automate
the cherry-pick workflow including creating a new branch, cherry-picking
commits, and creating a new pull request.

Target branch can be specified as 'branch' or 'remote/branch' format.`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		prNumber, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid pull request number: %s", args[0])
		}

		if prNumber <= 0 {
			return fmt.Errorf("pull request number must be positive, got: %d", prNumber)
		}

		targetBranch := args[1]
		if targetBranch == "" {
			return fmt.Errorf("target branch cannot be empty")
		}

		config := &github.Config{
			DryRun:          dryRun,
			SkipMergedCheck: skipMergedCheck,
		}

		return cherry.CherryPickPR(prNumber, targetBranch, config)
	},
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		version.PrintVersion()
	},
}

func main() {
	rootCmd.Flags().BoolVar(&dryRun, "dry-run", false, "show what would be done without executing remote operations")
	rootCmd.Flags().BoolVar(&skipMergedCheck, "skip-merged-check", false, "skip check that PR is merged (use with caution)")
	rootCmd.AddCommand(versionCmd)
	rootCmd.InitDefaultCompletionCmd()

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
