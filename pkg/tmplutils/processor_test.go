package tmplutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

func TestProcessTemplFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Define test cases
	testCases := []struct {
		name         string
		templatePath string
		destPath     string
		funcMap      template.FuncMap
		data         interface{}
		expectError  bool
	}{
		{
			name:         "Template file does not exist",
			templatePath: "nonexistent",
			destPath:     filepath.Join(tempDir, "dest"),
			expectError:  true,
		},
		{
			name:         "Cannot create destination directory",
			templatePath: filepath.Join(tempDir, "template"),
			destPath:     "/invalid/path/dest",
			expectError:  true,
		},
		{
			name:         "Cannot parse template",
			templatePath: filepath.Join(tempDir, "invalid_template"),
			destPath:     filepath.Join(tempDir, "dest"),
			expectError:  true,
		},
		{
			name:         "Cannot create destination file",
			templatePath: filepath.Join(tempDir, "template"),
			destPath:     "/invalid/path/dest",
			expectError:  true,
		},
		{
			name:         "Template execution fails",
			templatePath: filepath.Join(tempDir, "template"),
			destPath:     filepath.Join(tempDir, "dest"),
			funcMap:      template.FuncMap{"fail": func() (string, error) { return "", fmt.Errorf("fail") }},
			data:         map[string]interface{}{"Value": "fail"},
			expectError:  true,
		},
		{
			name:         "Everything is correct",
			templatePath: filepath.Join(tempDir, "template"),
			destPath:     filepath.Join(tempDir, "dest"),
			expectError:  false,
		},
	}

	// Run test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create template file if it should exist
			if tc.templatePath != "nonexistent" {
				var content []byte
				if tc.templatePath == filepath.Join(tempDir, "invalid_template") {
					content = []byte("{{")
				} else if tc.name == "Template execution fails" {
					content = []byte("{{fail}}") // Call the "fail" function in the template
				} else {
					content = []byte("{{.Value}}")
				}
				err := os.WriteFile(tc.templatePath, content, 0644)
				if err != nil {
					t.Fatalf("Failed to create template file: %v", err)
				}
			}

			// Call ProcessTemplFile
			err := ProcessTemplFile(tc.templatePath, tc.destPath, tc.funcMap, tc.data)

			// Check if error is expected
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
