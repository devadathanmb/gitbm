package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	"github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [name]",
	Short: "Create a new bookmark group",
	Long: `
Create a new bookmark group for the current Git project.

This command allows you to create a named collection of branch bookmarks. 
If no name is provided, it will use a random name.

Examples:
  gitbm create
  gitbm create migrate-to-kafka
  gitbm create "Feature: Add user profile"

The newly created bookmark group becomes the active group.`,
	Run: func(cmd *cobra.Command, args []string) {
		var bookmarkGroupName string

		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		if len(args) == 0 {
			logger.PrintWarning("No bookmark group name specified. Using a random name.")

			bookmarkGroupName = utils.GetRandomName()
		} else {
			bookmarkGroupName = args[0]
		}

		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)

		db, err := db.GetDB(dbFilePath)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		defer db.Close()

		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)
		// Create a new bookmark group
		bg := models.BookmarkGroup{
			Name: bookmarkGroupName,
		}

		err = bookmarkGroupRepo.Create(&bg)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					logger.PrintError("Bookmark group with the same name already exists")
					os.Exit(1)
				} else {
					logger.PrintError("Error creating bookmark group: %v", err)
					os.Exit(1)
				}
			} else {
				logger.PrintError("Error creating bookmark group: %v", err)
				os.Exit(1)
			}
		}

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		logger.PrintSuccess("Bookmark group created: %s", bookmarkGroupName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
