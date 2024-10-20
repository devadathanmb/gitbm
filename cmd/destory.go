package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/devadathanmb/gitbm/internal/utils"
	dbutils "github.com/devadathanmb/gitbm/internal/utils/dbUtils"
	"github.com/spf13/cobra"
)

var forceDestroy bool

var destroyCmd = &cobra.Command{
	Use:   "destroy [-f]",
	Short: "Remove all gitbm data",
	Long: `
Completely remove all gitbm data from the current Git repository.

WARNING: This action is irreversible. It will delete all bookmark groups,
branch bookmarks, and any other data created by gitbm in this repository.

This command:
- Deletes the gitbm database file
- Removes all stored bookmark groups and their associated branches
- Resets any gitbm-related configurations

By default, this command will prompt for confirmation before proceeding.
Use the -f or --force flag to bypass the confirmation prompt.

Use this command with extreme caution, typically when you want to start 
fresh or remove gitbm entirely from your project.

Examples:
  gitbm destroy
  gitbm destroy -f`,
	Run: func(cmd *cobra.Command, args []string) {
		if !forceDestroy {
			fmt.Print("Are you sure you want to destroy all gitbm data? This action cannot be undone. (Y/N): ")
			reader := bufio.NewReader(os.Stdin)
			response, _ := reader.ReadString('\n')
			response = strings.TrimSpace(strings.ToLower(response))

			if response != "y" && response != "yes" {
				fmt.Println("Operation cancelled.")
				return
			}
		}

		err := utils.ValidateBasic()
		if err != nil {
			logger.PrintError(fmt.Sprint(err))
			os.Exit(1)
		}

		// Remove the database file
		currentDir, _ := os.Getwd()
		dbPath := dbutils.GetDBPath(currentDir)
		err = os.Remove(dbPath)
		if err != nil {
			fmt.Println("Error removing gitbm database:", err)
			return
		}

		logger.PrintSuccess("gitbm data has been successfully destroyed.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().BoolVarP(&forceDestroy, "force", "f", false, "Force destroy without confirmation")
}
