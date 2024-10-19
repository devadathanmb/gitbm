/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// gitbm add <branch-alias>
		err := utils.ValidateBasic()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Get branch name
		var branchName string
		branchName, err = utils.GetCurrentGitBranch()
		if err != nil {
			fmt.Println("Error getting current branch name:", err)
			return
		}

		var branchAlias string
		// Branch alias
		if len(args) == 0 {
			fmt.Println("No branch alias specified. Using the branch name as the alias.")
			branchAlias = branchName
		} else {
			branchAlias = args[0]
		}

		// Get db connection
		currentDir, _ := os.Getwd()
		dbFilePath := utils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			fmt.Println("Error getting db connection:", err)
			return
		}

		// Get current bookmark group id
		currentBookmarkGroupId, err := models.GetCurrentBookmarkGroupId(db)

		if err != nil {
			fmt.Println("Error getting current bookmark group id:", err)
			return
		}

		// Add the branch to the db
		branch := &models.Branch{
			BookmarkGroupID: currentBookmarkGroupId,
			Name:            branchName,
			Alias:           branchAlias,
		}
		err = branch.Create(db)
		if err != nil {
			// handle error
		}

		if err != nil {
			fmt.Println("Error creating branch:", err)
			return
		}

		fmt.Println("Branch added successfully")

	},
}

func init() {
	rootCmd.AddCommand(addCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// addCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
