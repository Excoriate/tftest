package common

import (
	"fmt"
	"github.com/gruntwork-io/terratest/modules/git"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
	"os"
	"path/filepath"
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

	currentDir, _ := os.Getwd()
	gitRepoPath := git.GetRepoRoot(t)
	testDirRelPathFromGitRepo, err := filepath.Rel(gitRepoPath, filepath.Join(currentDir, tfDir))

	if err != nil {
		return "", fmt.Errorf("failed to get relative path to git repo: %v", err)
	}

	testFolder := test_structure.CopyTerraformFolderToTemp(t, gitRepoPath, testDirRelPathFromGitRepo)

	return testFolder, nil
}
