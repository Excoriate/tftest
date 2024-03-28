package validation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Excoriate/tftest/pkg/utils"
)

// HasTerraformFiles checks if the given directory has Terraform files with the given extensions.
// If the directory does not have any Terraform files with the given extensions, it returns an error.
func HasTerraformFiles(path string, extensions []string) error {
	files, err := filepath.Glob(filepath.Join(path, path, "*"))
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

	return nil
}

// IsValidTFDir checks if the given path is a valid Terraform directory.
// A valid Terraform directory is a directory that exists and is not empty.
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
func IsValidTFModuleDir(path string) error {
	if err := IsValidTFDir(path); err != nil {
		return err
	}

	if err := HasTerraformFiles(path, []string{".tf"}); err != nil {
		return err
	}

	return nil
}

// HasTFVarFiles checks if the given directory has Terraform variable files.
// If the directory does not have any Terraform variable files, it returns an error.
func HasTFVarFiles(path string) (bool, error) {
	if path == "" {
		return false, nil
	}

	files, err := filepath.Glob(filepath.Join(path, "*.tfvars"))
	if err != nil {
		return false, fmt.Errorf("failed to list files in directory: %v", err)
	}

	return len(files) > 0, nil
}
