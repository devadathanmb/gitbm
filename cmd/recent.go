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

var recentCmd = &cobra.Command{
	Use:   "recent [frequent]",
	Short: "List and checkout recently used branches",
	Long: `The 'recent' command lists the most recently checked out branches in your repository.
It provides an interactive fuzzy-finder interface to select and checkout a branch from the list.
By default, it shows the 10 most recently used branches. You can modify this behavior
using the available flags.

When used with the 'frequent' argument, it shows branches ordered by frequency of checkouts
rather than recency, helping you access your most commonly used branches.

This command is useful for quickly switching between branches you've been working on lately
or frequently use.`,
	Example: `  # List and select from the 10 most recently used branches
  gitbm recent
  # List and select from the 10 most frequently used branches
  gitbm recent frequent
  # List and select from the 5 most recently used branches
  gitbm recent --limit 5
  # List and select from the 5 most frequently used branches
  gitbm recent frequent --limit 5
  # List and select from the 10 least recently used branches
  gitbm recent --reverse
  # List and select from the 10 least frequently used branches
  gitbm recent frequent --reverse`,
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
		var branches []models.BranchCheckout
		if len(args) > 0 && args[0] == "frequent" {
			branches, err = branchCheckoutRepo.GetRecentFrequent(limit, isReverse)
		} else {
			branches, err = branchCheckoutRepo.GetRecent(limit, isReverse)
		}
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}
		selectedBranch, err := fzfutils.FuzzyFind(
			branches,
			func(b models.BranchCheckout) string {
				if b.LatestCommitMsg != "" {
					return fmt.Sprintf("%s -- %s", b.Name, b.LatestCommitMsg)
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
	rootCmd.AddCommand(recentCmd)
	recentCmd.Flags().IntP("limit", "l", 10, "Limit the number of branches to show (default 10)")
	recentCmd.Flags().BoolP("reverse", "r", false, "Show the least recently used branches instead of the most recent")
}
