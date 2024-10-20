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

var deleteCmd = &cobra.Command{
	Use:   "delete [group-name]",
	Short: "Delete a bookmark group",
	Long: `
Delete an existing bookmark group from the current Git project.

This command removes a specified bookmark group and all its associated branch bookmarks. 
If no group name is provided, it will delete the current active bookmark group.

Warning: This action is irreversible. All bookmarks in the deleted group will be lost.

Examples:
  gitbm delete old-feature
  gitbm delete "Completed Tasks"

If the deleted group was the active group, no group will be active after deletion.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate basic
		err := utils.ValidateBasic()

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
		}

		// Get db connection
		currentDir, _ := os.Getwd()
		dbPath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbPath)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
		}

		defer db.Close()

		var bookmarkGroupName string
		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)

		if len(args) == 0 {
			logger.PrintWarning("No bookmark group name specified. Using current bookmark group.")
			bookmarkGroup, err := bookmarkGroupRepo.GetCurrent()
			if bookmarkGroup == nil {
				logger.PrintInfo("No bookmark group set. Use `gitbm switch` to switch bookmark group.")
				os.Exit(0)
			}
			bookmarkGroupName = bookmarkGroup.Name
			if err != nil {
				logger.PrintError("Error getting current bookmark group name: %v", err)
				os.Exit(1)
			}
		} else {
			bookmarkGroupName = args[0]
		}

		err = bookmarkGroupRepo.Delete(bookmarkGroupName)

		if err != nil {
			logger.PrintError("Error deleting bookmark group: %v", err)
			os.Exit(1)
		}

		logger.PrintSuccess("Bookmark group deleted: %s", bookmarkGroupName)

	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
