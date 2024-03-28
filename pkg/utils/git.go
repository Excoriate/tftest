package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsAGitRepository checks if the given directory is a git repository.
func IsAGitRepository(repoRoot string) error {
	if repoRoot == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	if err := DirExistAndHasContent(repoRoot); err != nil {
		return err
	}

	_, err := os.Stat(filepath.Join(repoRoot, ".git"))
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s is not a git repository", repoRoot)
		}

		return fmt.Errorf("unexpected error when checking the directory %s: %v", repoRoot, err)
	}

	return nil
}
