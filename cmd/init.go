package cmd

import (
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	gitutils "github.com/devadathanmb/gitbm/internal/utils/gitUtils"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize gitbm for the current Git repository",
	Long: `
Initialize gitbm for the Git repository in the current directory.

This command sets up the necessary database and configurations for gitbm to manage
branch bookmarks in the current Git repository. It creates a .gitbm.db file in the
.git directory of the repository.

If gitbm is already initialized for the repository, this command will display an error.
To reinitialize, use 'gitbm destroy' first, then run 'gitbm init' again.

Note: 
- This command must be run from within a Git repository.
- It only affects the repository in the current working directory.

Example:
  cd /path/to/your/repo
  gitbm init`,

	Run: func(cmd *cobra.Command, args []string) {
		initDir, err := os.Getwd()
		if err != nil {
			logger.PrintError("Error getting current directory:", err)
			os.Exit(1)
		}

		isGitDir, err := utils.IsGitDir(initDir)
		if err != nil {
			logger.PrintError("Error checking if directory is a git repository:", err)
			os.Exit(1)
		}
		if !isGitDir {
			logger.PrintError("The specified directory is not a git repository")
			os.Exit(1)
		}

		// DB file path
		dbFilePath := dbutils.GetDBPath(initDir)

		// Validate if db file already exists
		doesDBExist, err := utils.DoesDBExist(dbFilePath)

		if err != nil {
			logger.PrintError("Error checking if database file exists:", err)
			os.Exit(1)
		}

		if doesDBExist {
			logger.PrintError("Gitbm is already initialized for this directory")
			os.Exit(1)
		}

		// Create the database file
		err = dbutils.CreateDB(dbFilePath)

		if err != nil {
			logger.PrintError("Error creating database file:", err)
			os.Exit(1)
		}

		// Initialize the database and run migrations
		err = db.InitDB(dbFilePath)
		if err != nil {
			logger.PrintError("Error initializing database:", err)
			os.Exit(1)
		}

		logger.PrintInfo("Initialized gitbm database")

		err = gitutils.InstallGitHook(initDir)

		if err != nil {
			logger.PrintError("Error installing gitbm hook:", err)
			os.Exit(1)
		}

		logger.PrintInfo("Installed gitbm hook")

		logger.PrintSuccess("Gitbm initialized successfully. Ready to use! 🚀")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
