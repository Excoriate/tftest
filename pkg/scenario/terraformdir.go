package scenario

import (
	"fmt"
	"testing"

	"github.com/Excoriate/tftest/pkg/validation"
)

func GetTerraformDir(t *testing.T, path string, isParallel bool) (string, error) {
	if err := validation.IsValidTFDir(path); err != nil {
		return "", fmt.Errorf("invalid terraform directory: %v", err)
	}

	if err := validation.HasTerraformFiles(path, []string{".tf"}); err != nil {
		return "", fmt.Errorf("no terraform files found in directory: %v", err)
	}

	if isParallel {
		return SetupTerraformDirForParallelism(t, path)
	}

	return path, nil
}
