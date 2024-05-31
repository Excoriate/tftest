package scenario

import (
	"testing"

	"github.com/Excoriate/tftest/pkg/validation"
)

// GetTerraformDir returns the Terraform directory path for the scenario.
// If the scenario is running in parallel, it sets up the Terraform directory for parallelism.
//
// Parameters:
//   - t: The testing instance.
//   - path: The path to the Terraform directory.
//   - isParallel: A boolean flag indicating whether the scenario is running in parallel.
//
// Returns:
//   - string: The path to the Terraform directory (or a temporary directory if running in parallel).
//   - error: An error if the Terraform directory is not valid or if parallel setup fails.
//
// Example:
//
//	terraformDir, err := GetTerraformDir(t, "/path/to/terraform/dir", true)
//	if err != nil {
//	    t.Fatalf("Error getting Terraform directory: %v", err)
//	}
//	fmt.Printf("Terraform directory: %s\n", terraformDir)
func GetTerraformDir(t *testing.T, path string, isParallel bool) (string, error) {
	if err := validation.IsValidTFModuleDir(path); err != nil {
		return "", err
	}

	if isParallel {
		return SetupTerraformDirForParallelism(t, path)
	}

	return path, nil
}
