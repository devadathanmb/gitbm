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

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show",
	Short: "Display the current bookmark group",
	Long: `
Show the currently active bookmark group in the repository.

This command fetches and displays the name of the active bookmark group, allowing
you to confirm which group you're working with. If no bookmark group is currently
set, it suggests using 'gitbm destroy' to reset and start over.

Usage:
  gitbm show

Example:
  gitbm show

Note: This command must be run within a Git repository initialized with gitbm.`,
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

		// Get current bookmark group id
		bookmarkGroupRepo := models.NewBookmarkGroupRepository(db)
		bookmarkGrp, err := bookmarkGroupRepo.GetCurrent()

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
		}

		if bookmarkGrp == nil {
			logger.PrintInfo("No bookmark group set. Better `gitbm destory` and start over.")
			return
		}

		logger.PrintSuccess("Current bookmark group: %s*", bookmarkGrp.Name)

	},
}

func init() {
	rootCmd.AddCommand(showCmd)

}
