package scenario

import (
	"fmt"
	"testing"

	"github.com/Excoriate/tftest/pkg/tfdir"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
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

	testDirRelPathFromGitRepo, gitRepoRoot, err := tfdir.GetRelativePathFromGitRepo(tfDir, t)
	if err != nil {
		return "", fmt.Errorf("failed to resolve test directory: %v", err)
	}

	return test_structure.CopyTerraformFolderToTemp(t, gitRepoRoot, testDirRelPathFromGitRepo), nil
}
