package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/goombaio/namegenerator"
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

// Gets random name as the title says
func GetRandomName() (name string) {
	seed := time.Now().UTC().UnixNano()
	nameGenerator := namegenerator.NewNameGenerator(seed)

	name = nameGenerator.Generate()

	return name
}

// Get current git branch name
func GetCurrentGitBranch() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
