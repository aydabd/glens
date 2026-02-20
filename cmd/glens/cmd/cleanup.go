package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"glens/tools/glens/internal/github"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up test issues from GitHub repository",
	Long: `Closes all test-related issues in the specified GitHub repository.

This is useful for cleaning up issues created during integration testing.
By default, it closes all issues with the "ai-generated" label.

Example:
  glens cleanup --github-repo aydabd/test-agent-ideas
  glens cleanup --github-repo aydabd/test-agent-ideas --labels test-failure,integration-test
  glens cleanup --github-repo aydabd/test-agent-ideas --dry-run`,
	RunE: runCleanup,
}

func init() {
	rootCmd.AddCommand(cleanupCmd)

	cleanupCmd.Flags().String("github-repo", "", "GitHub repository for cleanup (owner/repo)")
	cleanupCmd.Flags().StringSlice("labels", []string{"ai-generated"}, "Labels to filter issues for cleanup")
	cleanupCmd.Flags().Bool("dry-run", false, "List issues that would be closed without actually closing them")

	_ = viper.BindPFlag("github.repository", cleanupCmd.Flags().Lookup("github-repo"))
	_ = viper.BindPFlag("cleanup.labels", cleanupCmd.Flags().Lookup("labels"))
	_ = viper.BindPFlag("cleanup.dry_run", cleanupCmd.Flags().Lookup("dry-run"))
}

func runCleanup(_ *cobra.Command, _ []string) error {
	ctx := context.Background()

	// Get GitHub token
	githubToken := os.Getenv("GITHUB_TOKEN")
	if githubToken == "" {
		return fmt.Errorf("GITHUB_TOKEN environment variable is required")
	}

	// Get repository
	githubRepo := viper.GetString("github.repository")
	if githubRepo == "" {
		githubRepo = os.Getenv("GITHUB_REPOSITORY")
	}
	if githubRepo == "" {
		return fmt.Errorf("github repository is required (use --github-repo flag or GITHUB_REPOSITORY env var)")
	}

	// Get labels
	labels := viper.GetStringSlice("cleanup.labels")
	if len(labels) == 0 {
		labels = []string{"ai-generated"}
	}

	// Get dry-run flag
	dryRun := viper.GetBool("cleanup.dry_run")

	log.Info().
		Str("repository", githubRepo).
		Strs("labels", labels).
		Bool("dry_run", dryRun).
		Msg("Starting cleanup operation")

	// Create GitHub client
	githubClient, err := github.NewClient(githubToken)
	if err != nil {
		return fmt.Errorf("failed to create GitHub client: %w", err)
	}

	if err := githubClient.SetRepository(githubRepo); err != nil {
		return fmt.Errorf("failed to set repository: %w", err)
	}

	// List issues
	issues, err := githubClient.ListIssuesByLabel(ctx, labels)
	if err != nil {
		return fmt.Errorf("failed to list issues: %w", err)
	}

	if len(issues) == 0 {
		log.Info().Msg("No issues found matching the specified labels")
		return nil
	}

	// Count open issues
	openCount := 0
	for _, issue := range issues {
		if issue.GetState() == "open" {
			openCount++
		}
	}

	log.Info().
		Int("total_issues", len(issues)).
		Int("open_issues", openCount).
		Int("closed_issues", len(issues)-openCount).
		Msg("Found issues")

	if dryRun {
		fmt.Println("\n🔍 Dry-run mode: The following issues would be closed:")
		fmt.Println()
		for _, issue := range issues {
			if issue.GetState() == "open" {
				fmt.Printf("  #%-4d [%s] %s\n",
					issue.GetNumber(),
					issue.GetState(),
					issue.GetTitle())
			}
		}
		fmt.Printf("\nTotal: %d open issue(s) would be closed\n", openCount)
		return nil
	}

	// Close issues
	if openCount > 0 {
		fmt.Printf("\n🧹 Closing %d open issue(s)...", openCount)
		fmt.Println()
		fmt.Println()
		closedCount, err := githubClient.CloseTestIssues(ctx, labels)
		if err != nil {
			return fmt.Errorf("failed to close issues: %w", err)
		}

		fmt.Printf("✅ Successfully closed %d issue(s)\n", closedCount)
	} else {
		fmt.Println("\n✨ All matching issues are already closed!")
	}

	return nil
}
