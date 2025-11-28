package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh-moles",
	Short: "GitHub CLI extension with tools for GitHub repositories",
	Long: `gh-moles-extra is a GitHub CLI extension that provides various tools
for managing GitHub repositories. Install it with:

  gh extension install the78mole/gh-moles-extra

Then use it as:

  gh moles <command>`,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(runCmd)
}
