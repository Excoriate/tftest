package tmplutils

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

// ProcessTemplFile processes a template file with the provided data and writes the output to the destination file.
// It returns an error if any of the operations fail.
// The function map is used to define custom functions that can be called from the template.
func ProcessTemplFile(templatePath, destPath string, funcMap template.FuncMap, data interface{}) error {
	content, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("failed to read template file: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
		return fmt.Errorf("failed to create parent directory for %s: %w", destPath, err)
	}

	tmpl, tmplErr := template.New(filepath.Base(templatePath)).Funcs(funcMap).Parse(string(content))
	if tmplErr != nil {
		return fmt.Errorf("failed to parse template: %w", tmplErr)
	}

	// Write to destination with processed content
	file, fileErr := os.Create(destPath)
	if fileErr != nil {
		return fmt.Errorf("failed to create destination file: %w", fileErr)
	}

	defer file.Close()

	// Execute the template with provided data
	if err := tmpl.Execute(file, data); err != nil {
		return fmt.Errorf("failed to execute template with provided data: %w", err)
	}

	return nil
}
