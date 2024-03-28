package scenario

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/git"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// SetupTerraformDirForParallelism sets up the Terraform directory for parallelism.
// It copies the Terraform directory to a temporary directory and returns the path to the temporary directory.
func SetupTerraformDirForParallelism(t *testing.T, tfDir string) (string, error) {
	if t == nil {
		return "", fmt.Errorf("t is required")
	}

	if tfDir == "" {
		return "", fmt.Errorf("tfDir is required")
	}

	// if the path is absolute we can just use it, otherwise we calculate the path relative to the git repo root
	getRepoRoot := git.GetRepoRoot(t)
	testDirRelPathFromGitRepo, err := resolveTfDir(tfDir, getRepoRoot)
	if err != nil {
		return "", fmt.Errorf("failed to resolve test directory: %v", err)
	}

	return test_structure.CopyTerraformFolderToTemp(t, getRepoRoot, testDirRelPathFromGitRepo), nil
}

func resolveTfDir(tfDir string, repoRoot string) (string, error) {
	if !strings.HasPrefix(tfDir, "/") {
		currentDir, _ := os.Getwd()
		resolvedPath, err := filepath.Rel(repoRoot, filepath.Join(currentDir, tfDir))
		if err != nil {
			return "", fmt.Errorf("failed to get relative path to git repo: %v", err)
		}
		return resolvedPath, nil
	}
	return strings.TrimPrefix(tfDir, repoRoot), nil
}
