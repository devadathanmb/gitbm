/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/devadathanmb/gitbm/internal/db"
	"github.com/devadathanmb/gitbm/internal/db/models"
	"github.com/devadathanmb/gitbm/internal/utils"
	"github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a bookmark group",
	Long:  `Creates a new branch bookmark group for the git project`,
	Run: func(cmd *cobra.Command, args []string) {
		// What should this command do?
		// 1. Check if the current directory is a git directory
		// 2. Check if the database file exists
		// 3. Create a new bookmark group

		// Check if the current directory is a git directory
		var bookmarkGroupName string
		var gitDir string
		var err error

		if len(args) == 0 {
			fmt.Println("No bookmark group name specified. Using a random name.")

			bookmarkGroupName = utils.GetRandomName()
		} else {
			bookmarkGroupName = args[0]
		}

		isGitDir, err := utils.IsGitDir(gitDir)
		if err != nil {
			fmt.Println("Error checking if directory is a git repository:", err)
			return
		}
		if !isGitDir {
			fmt.Println("The specified directory is not a git repository")
			return
		}

		// DB file path
		dbFilePath := utils.GetDBPath(gitDir)

		// Validate if db file already exists
		doesDBExist, err := utils.DoesDBExist(dbFilePath)

		if err != nil {
			fmt.Println("Error checking if database file exists:", err)
			return
		}

		if !doesDBExist {
			fmt.Println("Database does not exist. Run 'gitbm init' to initialize the database")
			return
		}

		// Get DB connection
		db, err := db.GetDB(dbFilePath)

		if err != nil {
			fmt.Println("Error getting database connection:", err)
			return
		}

		// Create a new bookmark group
		bg := models.BookmarkGroup{
			Name: bookmarkGroupName,
		}

		err = bg.Create(db)
		if err != nil {
			if sqliteErr, ok := err.(sqlite3.Error); ok {
				if sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
					fmt.Printf("Error: Bookmark group '%s' already exists.\n", bookmarkGroupName)
				} else {
					fmt.Printf("SQLite error occurred: %v\n", sqliteErr)
				}
			} else {
				fmt.Printf("Error creating bookmark group: %v\n", err)
			}
			return
		}

		if err != nil {
			fmt.Println("Error creating stuff", err)
			return
		}

		fmt.Println("Bookmark group created successfully. Name:", bookmarkGroupName)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
