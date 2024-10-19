package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/spf13/cobra"
)

// Must handle two cases:
// 1. gitbm init - Initializes gitbm for the current dir
// 2. gitbm init <dir> - Initializes gitbm for the specified dir
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gitbm for the git project",
	Run: func(cmd *cobra.Command, args []string) {
		initDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current working directory:", err)
			return
		}

		isGitDir, err := utils.IsGitDir(initDir)
		if err != nil {
			fmt.Println("Error checking if directory is a git repository:", err)
			return
		}
		if !isGitDir {
			fmt.Println("The specified directory is not a git repository")
			return
		}

		// DB file path
		dbFilePath := utils.GetDBPath(initDir)

		// Validate if db file already exists
		doesDBExist, err := utils.DoesDBExist(dbFilePath)

		if err != nil {
			fmt.Println("Error checking if database file exists:", err)
			return
		}

		if doesDBExist {
			fmt.Println("Gitbm seems to be already initialized for this directory. If you want to reinitialize, please run gitbm destroy and then gitbm init")
			return
		}

		// Create the database file
		err = utils.CreateDB(dbFilePath)

		if err != nil {
			fmt.Println("Error creating database file:", err)
			return
		}

		// Initialize the database and run migrations
		err = db.InitDB(dbFilePath)
		if err != nil {
			fmt.Println("Error initializing database:", err)
			return
		}

		fmt.Println("Gitbm initialized successfully. Ready to use!")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
