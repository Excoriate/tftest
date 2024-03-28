package scenario

import (
	"testing"

	"github.com/Excoriate/tftest/pkg/validation"
)

// GetTerraformDir returns the Terraform directory path.
// It's used to get the Terraform directory path for the scenario.
// If the scenario is running in parallel, it sets up the Terraform directory for parallelism.
func GetTerraformDir(t *testing.T, path string, isParallel bool) (string, error) {
	if err := validation.IsValidTFModuleDir(path); err != nil {
		return "", err
	}

	if isParallel {
		return SetupTerraformDirForParallelism(t, path)
	}

	return path, nil
}
