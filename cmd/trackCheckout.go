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

// Non-user facing command to track the checkouts of a branch
// This command should track two things:
// 1. Increment the number of times the branch was checked out
// 2. Update the last_checked_out_at timestamp
var trackCheckoutCmd = &cobra.Command{
	Use:    "track-checkout",
	Short:  "Internal command to track the checkouts of a branch",
	Long:   `You should not be using this!`,
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate basic
		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		// Validate incoming args
		if len(args) != 2 {
			logger.PrintError("Invalid number of arguments")
			os.Exit(1)
		}

		// Get db
		currentDir, _ := os.Getwd()
		dbPath := dbutils.GetDBPath(currentDir)
		db, err := db.GetDB(dbPath)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		defer db.Close()

		// Get the branch name
		branchCheckoutRepo := models.NewBranchCheckoutRepository(db)
		branchCheckout := &models.BranchCheckout{
			Name:            args[0],
			LatestCommitMsg: args[1],
		}
		err = branchCheckoutRepo.Upsert(branchCheckout)

		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(trackCheckoutCmd)
}
