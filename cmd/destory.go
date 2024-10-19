package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/spf13/cobra"
)

var forceDestroy bool

var destroyCmd = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy gitbm data",
	Long:  `WARNING: There is no going back if you run this. It destroys your gitbm data.`,
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

		// Get the current directory
		currentDir, err := os.Getwd()
		if err != nil {
			fmt.Println("Error getting current directory:", err)
			return
		}

		// Check if it's a git repository
		isGitDir, err := utils.IsGitDir(currentDir)
		if err != nil {
			fmt.Println("Error checking if directory is a git repository:", err)
			return
		}
		if !isGitDir {
			fmt.Println("The current directory is not a git repository")
			return
		}

		// Get the database file path
		dbPath := utils.GetDBPath(currentDir)

		// Check if db file already exists
		doesDBExist, err := utils.DoesDBExist(dbPath)
		if err != nil {
			fmt.Println("Error checking if database file exists:", err)
			return
		}

		// Check if the database file exists
		if !doesDBExist {
			fmt.Println("No gitbm database found. Why destroy nothing?")
			return
		}

		// Remove the database file
		err = os.Remove(dbPath)
		if err != nil {
			fmt.Println("Error removing gitbm database:", err)
			return
		}

		fmt.Println("gitbm data has been successfully destroyed.")
	},
}

func init() {
	rootCmd.AddCommand(destroyCmd)
	destroyCmd.Flags().BoolVarP(&forceDestroy, "force", "f", false, "Force destroy without confirmation")
}
