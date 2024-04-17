package git_tools

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func FindGitRepoRootUsingGit(dir string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Dir = dir
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// FindGitRepoRootByTraversal finds the Git repository root for the given directory by manually checking for a .git directory.
func FindGitRepoRootByTraversal(dir string) (string, error) {
	currentPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path for %s: %v", dir, err)
	}

	// Check up to the filesystem root
	for currentPath != filepath.Dir(currentPath) { // Continue until the root directory is reached
		if _, err := os.Stat(filepath.Join(currentPath, ".git")); err == nil {
			return currentPath, nil
		}

		// Move up one directory level
		currentPath = filepath.Dir(currentPath)
	}

	return "", fmt.Errorf("no Git repository found starting from directory %s", dir)
}
