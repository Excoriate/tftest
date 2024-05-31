package git_tools

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Excoriate/tftest/pkg/utils"
)

// IsAGitRepository checks if the given directory or any of its parent directories up to `levels` is a Git repository.
// It returns the git root directory, the subdirectory passed relative to the git root, and any error encountered.
//
// Parameters:
//   - repoRoot: The directory to start the search from.
//   - levels: The number of parent directories to check upwards.
//
// Returns:
//   - gitRoot: The absolute path to the root of the Git repository.
//   - subDir: The relative path from the Git root to the original directory.
//   - err: An error if the Git repository root could not be found or if any other error occurred.
//
// Example:
//
//	gitRoot, subDir, err := IsAGitRepository("/path/to/start", 5)
//	if err != nil {
//	    log.Fatalf("Error finding Git repository: %v", err)
//	}
//	fmt.Printf("Git repository root: %s, Subdirectory: %s\n", gitRoot, subDir)
func IsAGitRepository(repoRoot string, levels int) (gitRoot, subDir string, err error) {
	if repoRoot == "" {
		return "", "", fmt.Errorf("directory path cannot be empty")
	}
	if levels < 0 {
		return "", "", fmt.Errorf("levels must be non-negative")
	}

	originalPath, err := filepath.Abs(repoRoot)
	if err != nil {
		return "", "", fmt.Errorf("failed to resolve absolute path for %s: %v", repoRoot, err)
	}

	if err := utils.DirExistAndHasContent(originalPath); err != nil {
		return "", "", err
	}

	currentPath := originalPath
	for i := 0; i <= levels; i++ {
		gitDir := filepath.Join(currentPath, ".git")
		if info, err := os.Stat(gitDir); err == nil && info.IsDir() {
			relPath, err := filepath.Rel(currentPath, originalPath)
			if err != nil {
				return "", "", fmt.Errorf("failed to calculate relative path from %s to %s: %v", currentPath, originalPath, err)
			}
			return currentPath, relPath, nil
		} else if !os.IsNotExist(err) {
			// If the error is not because the .git directory doesn't exist, return it.
			return "", "", fmt.Errorf("unexpected error when checking the directory %s: %v", gitDir, err)
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
