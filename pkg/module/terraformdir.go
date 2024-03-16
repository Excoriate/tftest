package module

import (
	"fmt"
	"github.com/excoriate/tftest/pkg/validation"
)

func GetTerraformDir(path string) (string, error) {
	if err := validation.IsValidTFDir(path); err != nil {
		return "", fmt.Errorf("invalid terraform directory: %v", err)
	}

	if err := validation.HasTerraformFiles(path, []string{".tf"}); err != nil {
		return "", fmt.Errorf("no terraform files found in directory: %v", err)
	}

	return path, nil
}
