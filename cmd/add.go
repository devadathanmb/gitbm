package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	gitutils "github.com/devadathanmb/gitbm/internal/utils/gitUtils"
	"github.com/spf13/cobra"
)

var branchNameFlag string // Variable to hold the value of --branch flag

var addCmd = &cobra.Command{
	Use:   "add [branch-alias]",
	Short: "Add current branch (or specified branch) to the active bookmark group",
	Long: `
Add the current Git branch or a specified branch to the active bookmark group.

This command bookmarks the current Git branch in the active bookmark group. 
You can optionally specify an alias for the branch, and you can also specify a branch name using the --branch (-b) flag.

If no alias is provided, the branch name will be used as the alias.

Usage:
  gitbm add [branch-alias] [--branch <branch-name>]

Examples:
  gitbm add                        # Adds the current branch
  gitbm add "Feature X"             # Adds the current branch with an alias
  gitbm add --branch feature/1234   # Adds a specific branch
  gitbm add --branch feature/1234 "Alias for Feature X"

Note:
- This command must be run from within a Git repository.
- A bookmark group must be active (use 'gitbm switch' if none is active).
- If the branch is already bookmarked in the current group, this command will fail.

The command will:
1. Get the current Git branch name or use the provided branch name.
2. Use the provided alias or the branch name if no alias is given.
3. Add the branch to the active bookmark group in the database.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate basic requirements
		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		// Get branch name
		var branchName string
		if branchNameFlag != "" {
			// Use the branch specified by --branch flag
			branchName = branchNameFlag
		} else {
			// Get the current branch if no --branch flag is provided
			branchName, err = gitutils.GetCurrentGitBranch()
			if err != nil {
				logger.PrintError("Error getting current branch name: %v", err)
				os.Exit(1)
			}
		}

		// Handle branch alias
		var branchAlias string
		if len(args) == 0 {
			logger.PrintWarning("No branch alias specified. Using branch name as alias.")
			branchAlias = branchName
		} else {
			branchAlias = args[0]
		}

		// Get db connection
		currentDir, _ := os.Getwd()
		dbFilePath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			logger.PrintError("Error getting db connection: %v", err)
			os.Exit(1)
		}

		defer db.Close()

		// Get current bookmark group id
		currentBookmarkGrpRepo := models.NewCurrentBookmarkGroupRepository(db)
		currentBookmarkGroupId, err := currentBookmarkGrpRepo.GetCurrentBookmarkGroupId()

		if err != nil {
			logger.PrintError("Error getting current bookmark group id: %v", err)
			os.Exit(1)
		}

		branchRepo := models.NewBranchRepository(db)

		// Add the branch to the db
		branch := &models.Branch{
			BookmarkGroupID: currentBookmarkGroupId,
			Name:            branchName,
			Alias:           branchAlias,
		}

		err = branchRepo.Create(branch)
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		logger.PrintSuccess("Branch %s added successfully", branchName)
	},
}

func init() {
	// Register the --branch (-b) flag
	addCmd.Flags().StringVarP(&branchNameFlag, "branch", "b", "", "Specify a branch name to bookmark (default is current branch)")
	rootCmd.AddCommand(addCmd)
}
