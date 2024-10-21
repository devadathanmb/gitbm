package gitutils

import (
	"bytes"
	_ "embed"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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

//go:embed git-hooks/post-checkout
var PostCheckoutHook string

func InstallGitHook(gitDir string) error {
	hooksDir, err := GetGitHooksDir()
	if err != nil {
		return fmt.Errorf("failed to get hooks directory: %w", err)
	}
	if err := os.MkdirAll(hooksDir, 0755); err != nil {
		return fmt.Errorf("failed to create hooks directory: %w", err)
	}

	hookPath := filepath.Join(hooksDir, "post-checkout")
	if err := os.WriteFile(hookPath, []byte(PostCheckoutHook), 0755); err != nil {
		return fmt.Errorf("failed to write post-checkout hook: %w", err)
	}

	return nil
}

func GetGitDir(dir string) string {
	return filepath.Join(dir, ".git")
}

func GetGitHooksDir() (string, error) {
	// Run the git command to get the core.hooksPath value
	cmd := exec.Command("git", "config", "--get", "core.hooksPath")
	output, err := cmd.Output() // Capture output here
	if err != nil {
		// If git config --get fails, return the default hooks directory
		return filepath.Join(".git", "hooks"), nil
	}

	// Trim any trailing spaces/newlines from the output
	dir := strings.TrimSpace(string(output))

	// If no custom hooks path is set, fallback to the default .git/hooks directory
	if dir == "" {
		dir = filepath.Join(".git", "hooks")
	}

	return dir, nil
}
