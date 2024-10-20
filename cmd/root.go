package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gitbm",
	Short: "A CLI tool for bookmarking Git branches",
	Long: `
gitbm is a command-line tool designed to help you bookmark, manage, and switch between Git branches easily.

For more detailed documentation on each command, use 'gitbm <command> --help'.`,

	// Uncomment this if you have any persistent flags or configuration
	// PersistentPreRun: func(cmd *cobra.Command, args []string) {
	//   // You can add common setup or validation code here
	// },
}

// Execute starts the root command
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
