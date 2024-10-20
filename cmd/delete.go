package cmd

import (
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	fzfutils "github.com/devadathanmb/gitbm/internal/utils/fzfUtils"
	"github.com/spf13/cobra"
)

var bookmarkGroupNameFlag string

var deleteCmd = &cobra.Command{
	Use:   "delete [group-name]",
	Short: "Delete a bookmark group",
	Long: `
Delete an existing bookmark group from the current Git project.
This command removes a specified bookmark group and all its associated branch bookmarks. 
If no group name is provided, an interactive selection using fzf will be presented.
Warning: This action is irreversible. All bookmarks in the deleted group will be lost.
Examples:
  gitbm delete old-feature
  gitbm delete "Completed Tasks"
  gitbm delete -g current
  gitbm delete (for interactive selection)
If the deleted group was the active group, no group will be active after deletion.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := utils.ValidateBasic(); err != nil {
			logger.PrintError("%v", err)
			os.Exit(1)
		}

		currentDir, _ := os.Getwd()
		dbPath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbPath)
		if err != nil {
			logger.PrintError("Error getting db connection: %v", err)
			os.Exit(1)
		}
		defer db.Close()

		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)
		var bookmarkGroupName string

		if cmd.Flags().Changed("group") {
			if bookmarkGroupNameFlag == "current" {
				logger.PrintInfo("Using current bookmark group.")
				bookmarkGroup, err := bookmarkGroupRepo.GetCurrent()
				if err != nil {
					logger.PrintError("Error getting current bookmark group: %v", err)
					os.Exit(1)
				}
				if bookmarkGroup == nil {
					logger.PrintInfo("No bookmark group set. Use `gitbm switch` to switch bookmark group.")
					os.Exit(0)
				}
				bookmarkGroupName = bookmarkGroup.Name
			} else {
				bookmarkGroupName = bookmarkGroupNameFlag
			}
		} else if len(args) > 0 {
			bookmarkGroupName = args[0]
		} else {
			// No group specified, use fzf to select
			bookmarkGroupsList, err := bookmarkGroupRepo.List()
			if err != nil {
				logger.PrintError("Error getting bookmark groups: %v", err)
				os.Exit(1)
			}
			if len(bookmarkGroupsList) == 0 {
				logger.PrintInfo("No bookmark groups found. Use `gitbm add` to add a bookmark group.")
				os.Exit(0)
			}

			selected, err := fzfutils.FuzzyFind(
				bookmarkGroupsList,
				func(bg models.BookmarkGroup) string { return bg.Name },
				"Select a bookmark group to delete",
			)
			if err != nil {
				if err == fzfutils.ErrSelectionCancelled {
					logger.PrintInfo("Selection cancelled")
					os.Exit(0)
				}
				logger.PrintError("Error in fuzzy selection: %v", err)
				os.Exit(1)
			}
			bookmarkGroupName = selected.Name
		}

		err = bookmarkGroupRepo.Delete(bookmarkGroupName)
		if err != nil {
			logger.PrintError("Error deleting bookmark group: %v", err)
			os.Exit(1)
		}

		logger.PrintSuccess("Bookmark group '%s' deleted successfully", bookmarkGroupName)
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&bookmarkGroupNameFlag, "group", "g", "", "Remove the specified bookmark group")
	deleteCmd.Flag("group").NoOptDefVal = "current"
}
