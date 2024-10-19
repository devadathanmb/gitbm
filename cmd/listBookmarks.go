package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/spf13/cobra"
)

// listBookmarksCmd represents the listBookmarks command
var listBookmarksCmd = &cobra.Command{
	Use:   "bookmarks",
	Short: "A brief description of your command",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}

		isGitDir, err := utils.IsGitDir(currentDir)
		if err != nil {
			fmt.Println("Error checking if directory is a git repository:", err)
			return
		}
		if !isGitDir {
			fmt.Println("The specified directory is not a git repository")
			return
		}

		// DB file path
		dbFilePath := utils.GetDBPath(currentDir)

		// Validate if db file already exists
		doesDBExist, err := utils.DoesDBExist(dbFilePath)

		if err != nil {
			fmt.Println("Error checking if database file exists:", err)
			return
		}

		if !doesDBExist {
			fmt.Println("Gitbm does not seem to be initialized for this directory. Run 'gitbm init' to initialize the database")
			return
		}

		// Get db connection
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			fmt.Println("Error getting db connection:", err)
			return
		}

		// Now we can list the bookmarks
		bg := &models.BookmarkGroup{}
		bookmarksList, err := bg.List(db)

		if err != nil {
			fmt.Println("Error listing bookmark groups:", err)
			return
		}
		if len(bookmarksList) == 0 {
			fmt.Println("No bookmark groups found")
			fmt.Println("Create a new bookmark group using 'gitbm create <group-name>'")
			return
		}

		fmt.Println("Found the following bookmark groups:")
		for _, bookmark := range bookmarksList {
			fmt.Println("- ", bookmark)
		}
	},
}

func init() {
	listCmd.AddCommand(listBookmarksCmd)
}
