package tfvars

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetTFVarsFromWorkdir scans the provided workdir directory for all .tfvars files
// and returns their filenames. If workdir is empty, it returns an error.
//
// Parameters:
//   - workdir: The directory to scan for .tfvars files. This parameter is required.
//
// Returns:
//   - []string: A slice of filenames of all .tfvars files found in the workdir.
//   - error: An error if the workdir is empty or if there is an error during the directory traversal.
//
// Example:
//
//	tfvarFiles, err := GetTFVarsFromWorkdir("/path/to/terraform/dir")
//	if err != nil {
//	    log.Fatalf("Error getting .tfvars files: %v", err)
//	}
//	fmt.Printf(".tfvars files: %v\n", tfvarFiles)
func GetTFVarsFromWorkdir(workdir string) ([]string, error) {
	if workdir == "" {
		return nil, fmt.Errorf("workdir cannot be empty")
	}

	var tfvarFiles []string

	// Use filepath.Walk to traverse the directory tree rooted at workdir
	err := filepath.Walk(workdir, func(path string, _ os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".tfvars" {
			// If it does, add its filename to the slice
			tfvarFiles = append(tfvarFiles, filepath.Base(path))
		}

		// Return nil to continue the walk
		return nil
	})

	return tfvarFiles, err
}
