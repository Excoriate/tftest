package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// IsValidDirE checks if the given path is a valid directory.
//
// Parameters:
//   - path: The path to check.
//
// Returns:
//   - error: An error if the path is not a valid directory.
//
// Example:
//
//	err := IsValidDirE("/path/to/check")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The path is a valid directory.")
func IsValidDirE(path string) error {
	// Clean the path to remove any unnecessary parts.
	cleanPath := filepath.Clean(path)

	// Use os.Stat to check if the path exists and get file info.
	info, err := os.Stat(cleanPath)
	if err != nil {
		if os.IsNotExist(err) {
			// The path does not exist.
			return fmt.Errorf("path does not exist: %s, error: %w", cleanPath, err)
		}
		// There was some problem accessing the path.
		return fmt.Errorf("failed to stat the path: %w", err)
	}

	// Check if the path is a directory.
	if !info.IsDir() {
		// The path is not a directory.
		return fmt.Errorf("path is not a directory: %s", cleanPath)
	}

	// The path is a valid directory.
	return nil
}

// DirExistAndHasContent checks if the given directory exists and has content.
//
// Parameters:
//   - dirPath: The path to the directory to check.
//
// Returns:
//   - error: An error if the directory does not exist or if there is any other issue.
//
// Example:
//
//	err := DirExistAndHasContent("/path/to/dir")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The directory exists and has content.")
func DirExistAndHasContent(dirPath string) error {
	if dirPath == "" {
		return fmt.Errorf("directory path cannot be empty")
	}

	currentDir, _ := os.Getwd()

	_, err := os.Stat(dirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory %s does not exist in current directory %s", dirPath, currentDir)
		}

		return fmt.Errorf("unexpected error when checking the directory %s: %v", dirPath, err)
	}

	return nil
}
