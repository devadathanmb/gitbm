/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List bookmarks or branches",
	Long: `List bookmarks or branches. For example:
	gitbm list bookmarks - List all the bookmarks
	gitbm list branches - List all the branches
`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("List what? bookmarks or branches?")
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
