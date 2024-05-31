package utils

import (
	"fmt"
	"os"
)

// FileHasContent checks if the specified file exists and has content.
//
// Parameters:
//   - file: The path to the file to check.
//
// Returns:
//   - error: An error if the file does not exist, is empty, or if there is any other issue.
//
// Example:
//
//	err := FileHasContent("/path/to/file.txt")
//	if err != nil {
//	    log.Fatalf("Error: %v", err)
//	} else {
//	    fmt.Println("The file exists and has content.")
func FileHasContent(file string) error {
	fileInfo, err := os.Stat(file)
	if err != nil {
		return fmt.Errorf("failed to stat file: %v", err)
	}

	if fileInfo.Size() == 0 {
		return fmt.Errorf("file is empty: %s", file)
	}

	return nil
}
