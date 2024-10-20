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
	"github.com/spf13/cobra"
)

// switchCmd represents the switch command
var switchCmd = &cobra.Command{
	Use:   "switch [bookmark_group_name]",
	Short: "Switch to a different bookmark group",
	Long: `
Switch to a specified bookmark group, or select one interactively.

This command allows you to change the active bookmark group in your Git repository.
You can provide the bookmark group name as an argument, or if no name is given, it will 
present a list of bookmark groups for you to choose from using an interactive menu (fzf).

Usage:
  gitbm switch [bookmark_group_name]

Examples:
  # Switch to a specific bookmark group by name
  gitbm switch feature-group

  # Select a bookmark group interactively using fzf
  gitbm switch

Note: This command must be run from within a Git repository initialized with gitbm.`,
	Run: func(cmd *cobra.Command, args []string) {
		err := utils.ValidateBasic()

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
		}

		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		defer db.Close()

		var bookmarkGroupName string
		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)
		if len(args) > 0 {
			bookmarkGroupName = args[0]
		} else {
			bookmarkGroupsList, err := bookmarkGroupRepo.List()

			if err != nil {
				logger.PrintError("Error getting bookmark groups: %v", err)
				os.Exit(1)
			}
			if len(bookmarkGroupsList) == 0 {
				logger.PrintInfo("No bookmark groups found. Use `gitbm add` to add a bookmark group.")
				os.Exit(1)
			}

			// Show fzf menu
			bookmarkGroupNames := []string{}
			for _, bookmark := range bookmarkGroupsList {
				bookmarkGroupNames = append(bookmarkGroupNames, bookmark.Name)
			}

			selected, err := fzfutils.FuzzyFind(bookmarkGroupNames, func(bg string) string { return bg }, "Select bookmark group")

			if err != nil {
				if err == fzfutils.ErrSelectionCancelled {
					logger.PrintInfo("Selection cancelled")
					os.Exit(1)
				}
				logger.PrintError("Error in fuzzy selection: %v", err)
				os.Exit(1)
			}

			bookmarkGroupName = selected
			logger.PrintInfo("Selected bookmark group: %s", bookmarkGroupName)

		}
		// Get bookmark group id
		bookmarkGroup, err := bookmarkGroupRepo.GetByName(bookmarkGroupName)
		if err != nil {
			logger.PrintError("Error getting bookmark group: %v", err)
			os.Exit(1)
		}

		// set current bookmark group
		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		err = currentBookmarkGrpRepo.SetCurrentBookmarkGroupId(bookmarkGroup.ID)

		if err != nil {
			logger.PrintError("Error setting current bookmark group: %v", err)
			os.Exit(1)
		}

		logger.PrintSuccess("Bookmark group switched to: %s*", bookmarkGroupName)
	},
}

func init() {
	rootCmd.AddCommand(switchCmd)
}
