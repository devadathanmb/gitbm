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

var listBookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "List all bookmark groups",
	Long: `
List all bookmark groups in the current Git repository.

This command displays all the bookmark groups that have been created
in the current repository. If no bookmark groups exist, it will suggest
using the 'gitbm add' command to create one.

Usage:
  gitbm list bookmarks

Example:
  gitbm list bookmarks

Note: This command must be run from within a Git repository initialized with gitbm.`,
	Run: func(cmd *cobra.Command, args []string) {

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

		// Now we can list the bookmarks
		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)
		// bg := &models.BookmarkGroup{}
		bookmarksList, err := bookmarkGroupRepo.List()

		if err != nil {
			logger.PrintError("Error getting bookmark groups: %v", err)
			os.Exit(1)
		}
		if len(bookmarksList) == 0 {
			logger.PrintError("No bookmark groups found. Use `gitbm add` to add a bookmark group.")
			os.Exit(1)
		}

		logger.PrintSuccess("Found bookmarks:")
		for _, bookmark := range bookmarksList {
			logger.Print(bookmark.Name)
		}
	},
}

func init() {
	listCmd.AddCommand(listBookmarksCmd)
}
