package validation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Excoriate/tftest/pkg/utils"
)

// HasTerraformFiles checks if the given directory has Terraform files with the given extensions.
// If the directory does not have any Terraform files with the given extensions, it returns an error.
//
// Parameters:
//   - path: The path to the directory to check.
//   - extensions: A list of file extensions to look for.
//
// Returns:
//   - error: An error if the directory does not have any Terraform files with the given extensions.
//
// Example:
//
//	err := HasTerraformFiles("/path/to/dir", []string{".tf"})
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The directory has Terraform files.")
func HasTerraformFiles(path string, extensions []string) error {
	files, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return fmt.Errorf("failed to list files in directory: %v", err)
	}

	for _, file := range files {
		for _, ext := range extensions {
			if filepath.Ext(file) == ext {
				return nil
			}
		}
	}

	return fmt.Errorf("no Terraform files with extensions %v found in directory: %s", extensions, path)
}

// IsValidTFDir checks if the given path is a valid Terraform directory.
// A valid Terraform directory is a directory that exists and is not empty.
//
// Parameters:
//   - path: The path to the directory to check.
//
// Returns:
//   - error: An error if the path is not a valid Terraform directory.
//
// Example:
//
//	err := IsValidTFDir("/path/to/dir")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid Terraform directory.")
func IsValidTFDir(path string) error {
	tfDirInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the terraform directory does not exist: %s", path)
	}

	if !tfDirInfo.IsDir() {
		return fmt.Errorf("the terraform directory is not a directory: %s", path)
	}

	return nil
}

// IsValidTFVarFile checks if the given path is a valid Terraform variable file.
// A valid Terraform variable file is a file with the .tfvars extension that is not empty.
//
// Parameters:
//   - path: The path to the file to check.
//
// Returns:
//   - error: An error if the path is not a valid Terraform variable file.
//
// Example:
//
//	err := IsValidTFVarFile("/path/to/variables.tfvars")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid Terraform variable file.")
func IsValidTFVarFile(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the terraform variable file does not exist: %s", path)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("the terraform variable file is a directory: %s", path)
	}

	if filepath.Ext(path) != ".tfvars" {
		return fmt.Errorf("the terraform variable file does not have a .tfvars extension: %s", path)
	}

	if err := utils.FileHasContent(path); err != nil {
		return fmt.Errorf("the terraform variable file is empty: %s", path)
	}

	return nil
}

// IsValidTFModuleDir checks if the given path is a valid Terraform module directory.
// A valid Terraform module directory is a directory that contains at least one Terraform file with the .tf extension.
// The path must also be a valid directory.
//
// Parameters:
//   - path: The path to the directory to check.
//
// Returns:
//   - error: An error if the path is not a valid Terraform module directory.
//
// Example:
//
//	err := IsValidTFModuleDir("/path/to/module")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid Terraform module directory.")
func IsValidTFModuleDir(path string) error {
	if err := IsValidTFDir(path); err != nil {
		return err
	}

	if err := HasTerraformFiles(path, []string{".tf"}); err != nil {
		return err
	}

	return nil
}

// HasTFVarFiles checks if the given directory has Terraform variable files with the .tfvars extension.
// If the directory does not have any Terraform variable files, it returns an error.
//
// Parameters:
//   - path: The path to the directory to check.
//
// Returns:
//   - bool: True if the directory has Terraform variable files, false otherwise.
//   - error: An error if the directory does not have any Terraform variable files or if there is any other issue.
//
// Example:
//
//	hasTFVars, err := HasTFVarFiles("/path/to/dir")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	}
//	if hasTFVars {
//	    fmt.Println("The directory has Terraform variable files.")
//	} else {
//	    fmt.Println("The directory does not have Terraform variable files.")
func HasTFVarFiles(path string) (bool, error) {
	if path == "" {
		return false, fmt.Errorf("path cannot be empty")
	}

	files, err := filepath.Glob(filepath.Join(path, "*.tfvars"))
	if err != nil {
		return false, fmt.Errorf("failed to list files in directory: %v", err)
	}

	return len(files) > 0, nil
}

// IsAHCLFile checks if the given path is a valid .hcl file.
// A valid .hcl file is a file with the .hcl extension that is not empty.
//
// Parameters:
//   - path: The path to the file to check.
//
// Returns:
//   - error: An error if the path is not a valid .hcl file.
//
// Example:
//
//	err := IsAHCLFile("/path/to/config.hcl")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid .hcl file.")
func IsAHCLFile(path string) error {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("the .hcl file does not exist: %s", path)
	}

	if fileInfo.IsDir() {
		return fmt.Errorf("the .hcl file is a directory: %s", path)
	}

	if filepath.Ext(path) != ".hcl" {
		return fmt.Errorf("the .hcl file does not have a .hcl extension: %s", path)
	}

	if err := utils.FileHasContent(path); err != nil {
		return fmt.Errorf("the .hcl file is empty: %s", path)
	}

	return nil
}
