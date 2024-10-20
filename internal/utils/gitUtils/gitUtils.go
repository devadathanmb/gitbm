package gitutils

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

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

func GitCheckout(branchName string) error {
	cmd := exec.Command("git", "checkout", branchName)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("error checking out branch: %v", err)
	}
	return nil
}
