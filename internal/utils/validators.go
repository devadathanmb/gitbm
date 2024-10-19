package utils

import (
	"os"
	"path/filepath"
)

func fileExists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// Function to validate if the current directory is a git directory
func IsGitDir(path string) (bool, error) {
	// Check if .git directory exists
	gitDir := filepath.Join(path, ".git")
	doesFileExist, err := fileExists(gitDir)

	return doesFileExist, err
}

func DoesDBExist(dbDir string) (bool, error) {
	// Check if the database file already exists
	doesDBExist, err := fileExists(dbDir)
	return doesDBExist, err
}

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

// Does basic and necessary validations for all commands
func ValidateBasic() (err error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return ValidationError{
			Message: "Error getting current working directory: " + err.Error(),
		}
	}

	isGitDir, err := IsGitDir(currentDir)
	if err != nil {
		return ValidationError{
			Message: "Error checking if directory is a git repository: " + err.Error(),
		}
	}

	if !isGitDir {
		return ValidationError{
			Message: "The specified directory is not a git repository",
		}
	}

	// DB file path
	dbFilePath := GetDBPath(currentDir)

	// Validate if db file already exists
	doesDBExist, err := DoesDBExist(dbFilePath)

	if err != nil {
		return ValidationError{
			Message: "Error checking if database file exists: " + err.Error(),
		}
	}

	if !doesDBExist {
		return ValidationError{
			Message: "Gitbm does not seem to be initialized for this directory. Run 'gitbm init' to initialize the database",
		}
	}

	return nil
}
