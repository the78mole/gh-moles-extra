package cmd

import (
	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Manage GitHub Actions workflow runs",
	Long:  `Commands for managing GitHub Actions workflow runs in a repository.`,
}

func init() {
	runCmd.AddCommand(cleanupCmd)
}
