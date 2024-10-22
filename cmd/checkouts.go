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

// checkoutsCmd represents the checkouts command
var checkoutsCmd = &cobra.Command{
	Use:   "checkouts",
	Short: "Reset the branch checkouts data",
	Long: `The 'reset checkouts' command resets the branch checkouts data in the repository.

It clears all the branch checkouts data stored in the database. This action is irreversible.

By default, it shows the top 10 most frequently used branches. You can modify this behavior
using the available flags.

Usage:
  gitbm reset checkouts [flags]

Examples:
  # Reset all branch checkouts data
  gitbm reset checkouts
`,

	Run: func(cmd *cobra.Command, args []string) {
		// Validate basic
		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		// Get the db connection
		cwd, _ := os.Getwd()
		dbPath := dbutils.GetDBPath(cwd)
		db, err := db.GetDB(dbPath)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		defer db.Close()

		// Remove all checkouts
		checkoutsRepo := models.NewBranchCheckoutRepository(db)
		err = checkoutsRepo.DeleteAll()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		logger.PrintSuccess("Checkouts data has been reset")
	},
}

func init() {
	resetCmd.AddCommand(checkoutsCmd)
}
