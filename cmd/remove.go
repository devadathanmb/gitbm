package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	fzfutils "github.com/devadathanmb/gitbm/internal/utils/fzfUtils"
	gitutils "github.com/devadathanmb/gitbm/internal/utils/gitUtils"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:   "remove [branch_name]",
	Short: "Remove a branch from the current bookmark group",
	Long: `
Remove the specified branch from the current bookmark group.
This command removes a branch from the currently active bookmark group in the repository's database.
If no branch name is provided, an interactive selection using fzf will be presented.
Usage:
  gitbm remove [branch_name]
Examples:
  # Remove the specified branch from the current bookmark group
  gitbm remove feature-branch
  # Interactively select a branch to remove
  gitbm remove
  # Remove the current branch
  gitbm remove -b current
Note: This command must be run from within a Git repository initialized with gitbm.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateBasic(); err != nil {
			logger.PrintError("%v", err)
			os.Exit(1)
		}

		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)
		if err != nil {
			logger.PrintError("Error getting db connection: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		currentBookmarkGroupId, err := currentBookmarkGrpRepo.GetCurrentBookmarkGroupId()
		if err != nil {
			logger.PrintError("Error getting current bookmark group id: %v", err)
			os.Exit(1)
		}

		branchRepo := models.NewBranchRepository(db)
		var branchName string

		if cmd.Flags().Changed("branch") {
			if branchNameFlag == "current" {
				branchName, err = gitutils.GetCurrentGitBranch()
				if err != nil {
					logger.PrintError("Error getting current branch name: %v", err)
					os.Exit(1)
				}
				logger.PrintInfo("Using current branch: %s", branchName)
			} else {
				branchName = branchNameFlag
			}
		} else if len(args) > 0 {
			branchName = args[0]
		} else {
			// No branch specified, use fzf to select
			branches, err := branchRepo.ListByBookmarkGroupId(currentBookmarkGroupId)
			if err != nil {
				logger.PrintError("Error getting branches: %v", err)
				os.Exit(1)
			}
			if len(branches) == 0 {
				logger.PrintError("No branches found. Use `gitbm add` to add a branch.")
				os.Exit(1)
			}

			selectedBranch, err := fzfutils.FuzzyFind(
				branches,
				func(b models.Branch) string {
					if b.Alias != "" {
						return fmt.Sprintf("%s -- %s", b.Name, b.Alias)
					}
					return b.Name
				},
				"Select a branch to remove",
			)
			if err != nil {
				if err == fzfutils.ErrSelectionCancelled {
					logger.PrintInfo("Branch selection cancelled")
					os.Exit(0)
				}
				logger.PrintError("Error selecting branch: %v", err)
				os.Exit(1)
			}
			branchName = selectedBranch.Name
		}

		branch, err := branchRepo.GetByName(currentBookmarkGroupId, branchName)
		if err != nil {
			logger.PrintError("Error getting branch: %v", err)
			os.Exit(1)
		}

		if err := branchRepo.Remove(branch.BookmarkGroupID, branch.Name); err != nil {
			logger.PrintError("Error removing branch: %v", err)
			os.Exit(1)
		}

		logger.PrintSuccess("Branch '%s' removed successfully from the current bookmark group", branchName)
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().StringVarP(&branchNameFlag, "branch", "b", "", "Remove the specified branch from the current bookmark group")
	removeCmd.Flag("branch").NoOptDefVal = "current"
}
