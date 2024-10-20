package cmd

import (
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	gitutils "github.com/devadathanmb/gitbm/internal/utils/gitUtils"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [branch_name]",
	Short: "Remove a branch from the current bookmark group",
	Long: `
Remove the specified branch from the current bookmark group.

This command removes a branch from the currently active bookmark group in the repository's database.
If no branch name is provided, the current branch will be used.

Usage:
  gitbm remove [branch_name]

Examples:
  # Remove the specified branch from the current bookmark group
  gitbm remove feature-branch

  # Remove the current branch from the current bookmark group
  gitbm remove

Note: This command must be run from within a Git repository initialized with gitbm.`,

	Run: func(cmd *cobra.Command, args []string) {
		err := utils.ValidateBasic()

		if err != nil {
			logger.PrintError("%v", err)
			os.Exit(1)
		}

		var branchName string

		if len(args) == 0 {
			logger.PrintWarning("No branch name specified. Using current branch.")
			branchName, err = gitutils.GetCurrentGitBranch()
			if err != nil {
				logger.PrintError("Error getting current branch name: %v", err)
				os.Exit(1)
			}
		} else {
			branchName = args[0]
		}

		// Get the db connection
		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			logger.PrintError("Error getting db connection: %v", err)
			os.Exit(1)
		}

		defer db.Close()

		// Get current bookmark group id
		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		currentBookmarkGroupId, err := currentBookmarkGrpRepo.GetCurrentBookmarkGroupId()

		if err != nil {
			logger.PrintError("Error getting current bookmark group id: %v", err)
			os.Exit(1)
		}

		branchRepo := models.NewBranchRepository(db)

		branch := &models.Branch{
			BookmarkGroupID: currentBookmarkGroupId,
			Name:            branchName,
		}

		branch, err = branchRepo.GetByName(currentBookmarkGroupId, branchName)

		if err != nil {
			logger.PrintError("Error getting branch: %v", err)
			os.Exit(1)
		}

		err = branchRepo.Remove(branch.BookmarkGroupID, branch.Name)

		if err != nil {
			logger.PrintError("Error removing branch: %v", err)
			os.Exit(1)
		}

		logger.PrintSuccess("Branch %s removed successfully", branchName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().BoolVarP(&forceDestroy, "force", "f", false, "Force destroy without confirmation")
}

// gitbm remove --branch -b [branch_name]
