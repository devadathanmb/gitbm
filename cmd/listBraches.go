package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/spf13/cobra"
)

// listBrachesCmd represents the listBraches command
var listBrachesCmd = &cobra.Command{
	Use:   "branches",
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
			fmt.Println("gitbm does not seem to be initialized for this directory. Run 'gitbm init' to initialize the database")
			return
		}

		// Get db connection
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			fmt.Println("Error getting db connection:", err)
			return
		}

		// Get current bookmark group id
		currentBookmarkGroupId, err := models.GetCurrentBookmarkGroupId(db)

		if err != nil {
			fmt.Println("Error getting current bookmark group id:", err)
			return
		}
		fmt.Println("Current bookmark group id:", currentBookmarkGroupId)

		// Get the branches
		branch := &models.Branch{BookmarkGroupID: currentBookmarkGroupId}
		branches, err := branch.ListByBookmarkGroupId(db)
		if err != nil {
			fmt.Println("Error getting branches:", err)
			return
		}

		for _, b := range branches {
			fmt.Printf("- Branch: %s, Alias: %s\n", b.Name, b.Alias)
		}
	},
}

func init() {
	listCmd.AddCommand(listBrachesCmd)
}
