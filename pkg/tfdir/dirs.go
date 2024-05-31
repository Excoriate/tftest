package tfdir

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/git"
)

// GetRelativePathFromGitRepo returns the relative path to the Git repository root for the specified Terraform directory.
// If the Terraform directory is an absolute path, it returns the relative path to the Git repository root.
// If the Terraform directory is a relative path, it returns the relative path to the Git repository root.
//
// Parameters:
//   - tfDir: The path to the Terraform directory. This parameter is required.
//   - t: The testing instance.
//
// Returns:
//   - relativePath: The relative path to the Git repository root from the Terraform directory.
//   - repoRoot: The root directory of the Git repository.
//   - err: An error if the relative path to the Git repository root could not be determined.
//
// Example:
//
//	relativePath, repoRoot, err := GetRelativePathFromGitRepo("/path/to/terraform/dir", t)
//	if err != nil {
//	    t.Fatalf("Error getting relative path from Git repo: %v", err)
//	}
//	fmt.Printf("Relative Path: %s, Repo Root: %s\n", relativePath, repoRoot)
func GetRelativePathFromGitRepo(tfDir string, t *testing.T) (relativePath, repoRoot string, err error) {
	if tfDir == "" {
		return "", "", fmt.Errorf("tfDir is required")
	}

	repoRoot = git.GetRepoRoot(t)

	t.Logf("The git repo root is %s", repoRoot)

	if !filepath.IsAbs(tfDir) {
		var currentDir string
		currentDir, err = os.Getwd()
		if err != nil {
			return "", "", fmt.Errorf("failed to get current directory: %v", err)
		}
		relativePath, err = filepath.Rel(repoRoot, filepath.Join(currentDir, tfDir))
		if err != nil {
			return "", "", fmt.Errorf("failed to get relative path to git repo: %v", err)
		}

		return relativePath, repoRoot, nil
	}

	if !strings.HasPrefix(tfDir, repoRoot) {
		return "", "", fmt.Errorf("the tfdir passed %s does not start with the repo root %s", tfDir, repoRoot)
	}

	relativePath = strings.TrimPrefix(tfDir, repoRoot)
	return relativePath, repoRoot, nil
}
