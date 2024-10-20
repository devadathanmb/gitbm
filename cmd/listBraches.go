package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	"github.com/spf13/cobra"
)

var listBranchesCmd = &cobra.Command{
	Use:   "branches",
	Short: "List all branches under the current bookmark group",
	Long: `
List all branches associated with the current bookmark group in the repository.

This command fetches and displays the branches tied to the currently active bookmark group
in the repository's database. If no branches are found, it will suggest adding them
via 'gitbm branch add' command.

Usage:
  gitbm list branches

Example:
  gitbm list branches

Note: This command must be run within a Git repository initialized with gitbm.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.ValidateBasic()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Get db connection
		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			fmt.Println("Error getting db connection:", err)
			return
		}

		defer db.Close()

		// Get current bookmark group id
		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		currentBookmarkGroupId, err := currentBookmarkGrpRepo.GetCurrentBookmarkGroupId()

		if err != nil {
			fmt.Println("Error getting current bookmark group id:", err)
			return
		}

		// Get the branches
		branchRepo := models.NewBranchRepository(db)
		branch := &models.Branch{BookmarkGroupID: currentBookmarkGroupId}
		branches, err := branchRepo.ListByBookmarkGroupId(branch.BookmarkGroupID)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(branches) == 0 {
			logger.PrintError("No branches found. Use `gitbm add` to add a branch.")
			os.Exit(1)
		}

		for _, b := range branches {
			fmt.Printf("- Branch: %s, Alias: %s\n", b.Name, b.Alias)
		}
	},
}

func init() {
	listCmd.AddCommand(listBranchesCmd)
}
