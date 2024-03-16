package validation

import (
	"fmt"
	"os"
	"path/filepath"
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
