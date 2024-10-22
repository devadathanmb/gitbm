package cmd

import (
	"github.com/devadathanmb/gitbm/internal/logger"
	"github.com/spf13/cobra"
)

var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset the branch checkouts data",
	Long:  `The 'reset' command resets the branch checkouts data in the repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			logger.Print("Reset what?")
			return
		}

	},
}

func init() {
	rootCmd.AddCommand(resetCmd)
}
