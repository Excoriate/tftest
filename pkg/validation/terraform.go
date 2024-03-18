package validation

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/excoriate/tftest/pkg/utils"
)

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
