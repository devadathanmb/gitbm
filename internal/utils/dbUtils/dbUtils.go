package dbutils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Create a gitbm.db sqlite database file in the given directory
func CreateDB(dir string) error {

	// Create the database file
	file, err := os.Create(dir)

	if err != nil {
		return fmt.Errorf("error creating database file: %v", err)
	}

	defer file.Close()
	return nil
}

// Helper function to get the path of the sqlite database file
func GetDBPath(dir string) string {
	dbFilePath := filepath.Join(dir, ".git", "gitbm.db")
	return dbFilePath
}
