/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [bookmarks|branches]",
	Short: "List bookmark groups or branches",
	Long: `
List either all bookmark groups or branches in the current bookmark group.

Usage:
  gitbm list bookmarks  - List all bookmark groups
  gitbm list branches   - List branches in the current bookmark group

If no argument is provided, this command will prompt you to specify 
whether to list bookmarks or branches.

Note: 
- 'bookmarks' refers to bookmark groups, not individual bookmarked branches.
- 'branches' shows only the branches in the currently active bookmark group.

Examples:
  gitbm list bookmarks
  gitbm list branches`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please specify what to list: 'bookmarks' or 'branches'")
			return
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
