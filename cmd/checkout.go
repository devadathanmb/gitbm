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

var checkoutCmd = &cobra.Command{
	Use:   "checkout [branch-name]",
	Short: "Checkout a branch from the current bookmark group",
	Long: `
Checkout a Git branch from the current bookmark group.

If a branch name is provided, it will checkout that specific branch. 
If no branch name is provided, it will open an interactive fuzzy-finder 
to select a branch from the current bookmark group.

Usage:
  gitbm checkout [branch-name]

Examples:
  gitbm checkout feature-branch
  gitbm checkout  # Opens fuzzy finder

The fuzzy finder allows searching by branch name or alias.

Note: An active bookmark group is required.`,
	Run: func(cmd *cobra.Command, args []string) {

		// Validate basic
		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		// Get db connection
		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		defer db.Close()

		// Get current bookmark group id
		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		currentBookmarkGroupId, err := currentBookmarkGrpRepo.GetCurrentBookmarkGroupId()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		if currentBookmarkGroupId == 0 {
			logger.PrintInfo("No bookmark group set.")
			return
		}

		var branchName string

		// Get the branches in the current bookmark group
		branchRepo := models.NewBranchRepository(db)
		if len(args) > 0 {
			branchName = args[0]

			// Validate if the branch exists in the current bookmark group
			branch := &models.Branch{Name: branchName, BookmarkGroupID: currentBookmarkGroupId}
			_, err := branchRepo.GetByName(currentBookmarkGroupId, branch.Name)
			if err != nil {
				logger.PrintError(fmt.Sprint(err))
				os.Exit(1)
			}
		} else {
			branch := &models.Branch{BookmarkGroupID: currentBookmarkGroupId}
			branches, err := branchRepo.ListByBookmarkGroupId(branch.BookmarkGroupID)
			if err != nil {
				logger.PrintError(fmt.Sprint(err))
				os.Exit(1)
			}

			if len(branches) == 0 {
				logger.PrintInfo("No branches in the current bookmark group.")
				return
			}

			// fzf the branches
			selectedBranch, err := fzfutils.FuzzyFind(
				branches,
				func(b models.Branch) string {
					if b.Alias != "" {
						return fmt.Sprintf("%s -- %s", b.Name, b.Alias)
					}
					return b.Name
				},
				"Select a branch:",
			)
			if err != nil {
				if err == fzfutils.ErrSelectionCancelled {
					logger.PrintInfo("Branch selection cancelled")
					os.Exit(1)
				}
				logger.PrintError("Error selecting branch: %v", err)
				os.Exit(1)
			}

			branchName = selectedBranch.Name
		}

		// Now git checkout to the branch
		err = gitutils.GitCheckout(branchName)
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		logger.PrintInfo("Checked out to branch: %s", branchName)

	},
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}
