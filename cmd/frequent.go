package cmd

import (
	"fmt"
	"os"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	fzfutils "github.com/devadathanmb/gitbm/internal/utils/fzfUtils"
	gitutils "github.com/devadathanmb/gitbm/internal/utils/gitUtils"
	"github.com/spf13/cobra"
)

var frequentCmd = &cobra.Command{
	Use:   "frequent",
	Short: "List and checkout frequently used branches",
	Long: `The 'frequent' command lists the most frequently checked out branches in your repository.

It provides an interactive fuzzy-finder interface to select and checkout a branch from the list.

By default, it shows the top 10 most frequently used branches. You can modify this behavior
using the available flags.

Usage:
  gitbm frequent [flags]

Examples:
  # List and select from the 10 most frequently used branches
  gitbm frequent

  # List and select from the 5 most frequently used branches
  gitbm frequent --limit 5

  # List and select from the 10 least frequently used branches
  gitbm frequent --reverse

  # List and select from the 5 least frequently used branches
  gitbm frequent --limit 5 --reverse`,
	Run: func(cmd *cobra.Command, args []string) {
		// Validate basic
		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
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

		var limit int
		var isReverse bool

		if cmd.Flags().Changed("limit") {
			limit, err = cmd.Flags().GetInt("limit")
			if err != nil {
				logger.PrintError(fmt.Sprint(err))
				os.Exit(1)
			}
		} else {
			limit = 10
		}

		if cmd.Flags().Changed("reverse") {
			isReverse, _ = cmd.Flags().GetBool("reverse")
		} else {
			isReverse = false
		}

		// Get the branch names and fzf them

		branchCheckoutRepo := models.NewBranchCheckoutRepository(db)
		branches, err := branchCheckoutRepo.GetFrequent(limit, isReverse)
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}
		selectedBranch, err := fzfutils.FuzzyFind(
			branches,
			func(b models.BranchCheckout) string {
				if b.LatestCommitMsg != "" {
					return fmt.Sprintf("%s -- %d -- %s", b.Name, b.CheckoutCount, b.LatestCommitMsg)
				}
				return b.Name
			},
			"Select a branch to checkout to",
		)
		if err != nil {
			if err == fzfutils.ErrSelectionCancelled {
				logger.PrintInfo("Branch selection cancelled")
				os.Exit(0)
			}
			logger.PrintError("Error selecting branch: %v", err)
			os.Exit(1)
		}

		err = gitutils.GitCheckout(selectedBranch.Name)
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(frequentCmd)
	frequentCmd.Flags().IntP("limit", "l", 10, "Limit the number of recent branches to show")
	frequentCmd.Flags().BoolP("reverse", "r", false, "Show the least recent branches")
	// gitbm frequent - should fzf with 10 most recent branches
	// gitbm frequent --rever - should fzf with 10 least recent branches
	// gitbm frequent --limit 5 - should fzf with 5 most recent branches
	// gitbm frequent --limit 5 --reverse - should fzf with 5 least recent branches
}
