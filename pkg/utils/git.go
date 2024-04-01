package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsAGitRepository checks if the given directory or any of its parent directories up to `levels` is a git repository.
// It returns the git root directory, the subdirectory passed relative to the git root, and any error encountered.
func IsAGitRepository(repoRoot string, levels int) (gitRoot string, subDir string, err error) {
	if repoRoot == "" {
		return "", "", fmt.Errorf("directory path cannot be empty")
	}

	originalPath, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve absolute path for %s: %v", repoRoot, err)
	}

	if err := DirExistAndHasContent(originalPath); err != nil {
		return "", "", err
	}

	currentPath := originalPath
	for i := 0; i <= levels; i++ {
		_, err := os.Stat(filepath.Join(currentPath, ".git"))
		if err == nil {
			relPath, _ := filepath.Rel(currentPath, originalPath)
			return currentPath, relPath, nil
		}

		if !os.IsNotExist(err) {
			// If the error is not because the .git directory doesn't exist, return it.
			return "", "", fmt.Errorf("unexpected error when checking the directory %s: %v", currentPath, err)
		}

		// Move up one directory level
		parentDir := filepath.Dir(currentPath)
		if parentDir == currentPath {
			// If the parent directory is the same as the current one, we've reached the filesystem root
			break
		}
		currentPath = parentDir
	}

	return "", "", fmt.Errorf("no git repository found within %d levels of directory %s", levels, repoRoot)
}
