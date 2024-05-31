package validation

import (
	"fmt"
	"path/filepath"

	"github.com/Excoriate/tftest/pkg/utils"
)

// IsATerragruntModule checks if the given path is a valid Terragrunt module.
// A valid Terragrunt module is a directory that contains a terragrunt.hcl file.
//
// Parameters:
//   - path: The path to the directory to check.
//
// Returns:
//   - error: An error if the path is not a valid Terragrunt module.
//
// Example:
//
//	err := IsATerragruntModule("/path/to/module")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid Terragrunt module.")
func IsATerragruntModule(path string) error {
	if err := utils.IsValidDirE(path); err != nil {
		return fmt.Errorf("the terragrunt module does not exist: %s", path)
	}

	if err := IsAHCLFile(filepath.Join(path, "terragrunt.hcl")); err != nil {
		return fmt.Errorf("the terragrunt module is not a valid terragrunt module: %s", path)
	}

	return nil
}
