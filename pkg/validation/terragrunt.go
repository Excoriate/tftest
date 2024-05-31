package validation

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/tftest/pkg/utils"
)

// IsATerragruntModule checks if the given path is a valid terragrunt module
//
// The path is expected to be a directory
// The directory must contain a terragrunt.hcl file
//
// Returns an error if the path is not a valid terragrunt module
func IsATerragruntModule(path string) error {
	if err := utils.IsValidDirE(path); err != nil {
		return fmt.Errorf("the terragrunt module does not exist: %s", path)
	}

	if err := IsAHCLFile(filepath.Join(path, "terragrunt.hcl")); err != nil {
		return fmt.Errorf("the terragrunt module is not a valid terragrunt module: %s", path)
	}

	return nil
}
