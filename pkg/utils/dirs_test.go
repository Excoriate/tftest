package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidDirE(t *testing.T) {
	// Create a temporary directory to simulate a valid directory.
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tempDir) // Clean up after the test.

	// Create a temporary file to simulate an invalid directory (not a directory).
	tempFile, err := ioutil.TempFile("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	tempFilePath := tempFile.Name()
	defer func(name string) {
		_ = os.Remove(name)
	}(tempFilePath) // Clean up after the test.

	tests := []struct {
		name    string
		path    string
		wantErr bool
	}{
		{"Valid Directory", tempDir, false},
		{"Invalid Path (File)", tempFilePath, true},
		{"Nonexistent Path", filepath.Join(tempDir, "nonexistent"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := IsValidDirE(tt.path)
			if (err != nil) != tt.wantErr {
				t.Errorf("IsValidDirE() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDirExistAndHasContent(t *testing.T) {
	// Reuse tempDir from the previous test setup.
	tempDir, err := ioutil.TempDir("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}
	defer func(path string) {
		_ = os.RemoveAll(path)
	}(tempDir) // Clean up after the test.

	tests := []struct {
		name    string
		dirPath string
		wantErr bool
	}{
		{"Empty Path", "", true},
		{"Existing Directory", tempDir, false},
		{"Nonexistent Directory", filepath.Join(tempDir, "nonexistent"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DirExistAndHasContent(tt.dirPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("DirExistAndHasContent() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
