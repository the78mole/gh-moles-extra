package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/cli/go-gh/v2/pkg/api"
	"github.com/cli/go-gh/v2/pkg/repository"
	"github.com/spf13/cobra"
)

var (
	autoConfirm      bool
	deleteFailedOnly bool
	keepCount        int
)

var cleanupCmd = &cobra.Command{
	Use:     "cleanup [KEEP_COUNT]",
	Aliases: []string{"clean"},
	Short:   "Delete old GitHub Actions workflow runs",
	Long: `Delete old workflow runs, keeping only the most recent ones OR delete all failed runs.

Examples:
  gh moles run cleanup           # Keep 20 most recent runs (with confirmation)
  gh moles run cleanup -y        # Keep 20 most recent runs (no confirmation)
  gh moles run cleanup 50        # Keep 50 most recent runs (with confirmation)
  gh moles run cleanup -y 50     # Keep 50 most recent runs (no confirmation)
  gh moles run cleanup --failed  # Delete all failed runs (with confirmation)
  gh moles run cleanup -y -f     # Delete all failed runs (no confirmation)`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCleanup,
}

func init() {
	cleanupCmd.Flags().BoolVarP(&autoConfirm, "yes", "y", false, "Skip confirmation prompt (auto-confirm deletion)")
	cleanupCmd.Flags().BoolVarP(&deleteFailedOnly, "failed", "f", false, "Delete all failed runs instead of keeping recent ones")
}

// workflowRun represents a GitHub Actions workflow run
type workflowRun struct {
	ID         int64  `json:"id"`
	Conclusion string `json:"conclusion"`
}

// workflowRunsResponse represents the API response for listing workflow runs
type workflowRunsResponse struct {
	TotalCount   int           `json:"total_count"`
	WorkflowRuns []workflowRun `json:"workflow_runs"`
}

func runCleanup(cmd *cobra.Command, args []string) error {
	// Default keep count
	keepCount = 20

	// Parse keep count from args if provided
	if len(args) > 0 {
		if deleteFailedOnly {
			return fmt.Errorf("KEEP_COUNT is ignored when using --failed")
		}
		var err error
		keepCount, err = strconv.Atoi(args[0])
		if err != nil || keepCount < 1 {
			return fmt.Errorf("KEEP_COUNT must be a positive integer")
		}
	}

	fmt.Println("🧹 GitHub Actions Run Cleanup")
	if deleteFailedOnly {
		fmt.Println("🚫 Deleting all failed runs...")
	} else {
		fmt.Printf("📊 Keeping the %d most recent runs...\n", keepCount)
	}
	fmt.Println()

	// Get current repository from git context
	repo, err := repository.Current()
	if err != nil {
		return fmt.Errorf("❌ Error: Could not determine current repository\nMake sure you're in a git repository: %w", err)
	}

	// Create REST API client (handles authentication automatically)
	client, err := api.DefaultRESTClient()
	if err != nil {
		return fmt.Errorf("❌ Error: Could not create GitHub API client\nPlease run: gh auth login\n%w", err)
	}

	fmt.Println("📈 Analyzing workflow runs...")

	// Get all workflow runs using the GitHub API
	runs, err := getWorkflowRuns(client, repo.Owner, repo.Name)
	if err != nil {
		return fmt.Errorf("failed to get workflow runs: %w", err)
	}

	var runsToDelete []int64

	if deleteFailedOnly {
		// Filter failed runs
		for _, run := range runs {
			if run.Conclusion == "failure" {
				runsToDelete = append(runsToDelete, run.ID)
			}
		}
		fmt.Printf("   Found %d failed runs\n", len(runsToDelete))

		if len(runsToDelete) == 0 {
			fmt.Println("✅ No failed runs found - nothing to clean up")
			return nil
		}
		fmt.Printf("🗑️  Will delete all %d failed runs\n", len(runsToDelete))
	} else {
		// Keep the newest keepCount runs
		totalRuns := len(runs)
		fmt.Printf("   Found %d total runs\n", totalRuns)

		if totalRuns <= keepCount {
			fmt.Printf("✅ No cleanup needed - only %d runs found (keeping %d)\n", totalRuns, keepCount)
			return nil
		}

		// Runs are already sorted by date (newest first) from the API
		// So we skip the first keepCount and delete the rest
		for i := keepCount; i < totalRuns; i++ {
			runsToDelete = append(runsToDelete, runs[i].ID)
		}
		fmt.Printf("🗑️  Will delete %d old runs (keeping newest %d)\n", len(runsToDelete), keepCount)
	}
	fmt.Println()

	if len(runsToDelete) == 0 {
		fmt.Println("✅ No runs to delete")
		return nil
	}

	fmt.Printf("   Found %d runs to delete\n", len(runsToDelete))
	fmt.Println()

	// Confirm deletion
	if !autoConfirm {
		if deleteFailedOnly {
			fmt.Printf("⚠️  This will permanently delete %d failed workflow runs.\n", len(runsToDelete))
		} else {
			fmt.Printf("⚠️  This will permanently delete %d workflow runs.\n", len(runsToDelete))
		}
		fmt.Print("Continue? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("❌ Cleanup cancelled")
			return nil
		}
	} else {
		fmt.Println("⚡ Auto-confirming deletion (--yes flag used)")
	}

	fmt.Println()
	if deleteFailedOnly {
		fmt.Println("🗑️  Deleting failed runs...")
	} else {
		fmt.Println("🗑️  Deleting old runs...")
	}

	// Delete runs in batches
	batchSize := 5
	deletedCount := 0
	failedCount := 0

	for i, runID := range runsToDelete {
		fmt.Printf("   Deleting run %d... ", runID)

		if err := deleteWorkflowRun(client, repo.Owner, repo.Name, runID); err != nil {
			fmt.Println("❌ (failed or already deleted)")
			failedCount++
		} else {
			fmt.Println("✅")
			deletedCount++
		}

		// Add small delay every batchSize deletions to be API-friendly
		if (i+1)%batchSize == 0 && i < len(runsToDelete)-1 {
			time.Sleep(1 * time.Second)
		}
	}

	fmt.Println()
	fmt.Println("📊 Cleanup Summary:")
	fmt.Printf("   ✅ Successfully deleted: %d runs\n", deletedCount)
	if failedCount > 0 {
		fmt.Printf("   ❌ Failed to delete: %d runs\n", failedCount)
	}
	if deleteFailedOnly {
		fmt.Println("   🚫 Deleted all failed runs")
	} else {
		fmt.Printf("   📈 Remaining runs: %d (newest)\n", keepCount)
	}
	fmt.Println()
	fmt.Println("🎉 Cleanup completed!")

	return nil
}

// getWorkflowRuns retrieves all workflow runs from the GitHub API
func getWorkflowRuns(client *api.RESTClient, owner, repo string) ([]workflowRun, error) {
	var allRuns []workflowRun
	page := 1
	perPage := 100 // Maximum allowed by API

	for {
		path := fmt.Sprintf("repos/%s/%s/actions/runs?per_page=%d&page=%d", owner, repo, perPage, page)
		var response workflowRunsResponse
		if err := client.Get(path, &response); err != nil {
			return nil, fmt.Errorf("failed to list workflow runs: %w", err)
		}

		allRuns = append(allRuns, response.WorkflowRuns...)

		// Check if we've got all runs
		if len(response.WorkflowRuns) < perPage || len(allRuns) >= response.TotalCount {
			break
		}
		page++
	}

	return allRuns, nil
}

// deleteWorkflowRun deletes a single workflow run using the GitHub API
func deleteWorkflowRun(client *api.RESTClient, owner, repo string, runID int64) error {
	path := fmt.Sprintf("repos/%s/%s/actions/runs/%d", owner, repo, runID)
	return client.Delete(path, nil)
}
