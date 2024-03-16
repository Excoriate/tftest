package utils

import (
	"fmt"
	"os"
)

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
