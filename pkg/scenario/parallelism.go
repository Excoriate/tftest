package scenario

import (
	"fmt"
	"testing"

	"github.com/Excoriate/tftest/pkg/tfdir"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// SetupTerraformDirForParallelism sets up the Terraform directory for parallelism by copying
// the specified Terraform directory to a temporary directory. This ensures that parallel tests
// do not interfere with each other.
//
// Parameters:
//   - t: The testing instance. This parameter is required.
//   - tfDir: The path to the Terraform directory. This parameter is required.
//
// Returns:
//   - string: The path to the temporary directory containing the copied Terraform directory.
//   - error: An error if the setup fails.
//
// Example:
//
//	tempDir, err := SetupTerraformDirForParallelism(t, "/path/to/terraform/dir")
//	if err != nil {
//	    t.Fatalf("Error setting up Terraform directory for parallelism: %v", err)
//	}
//	fmt.Printf("Temporary Terraform directory: %s\n", tempDir)
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
