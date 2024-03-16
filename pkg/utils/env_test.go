package utils

import (
	"os"
	"reflect"
	"testing"
)

func TestCleanValue(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"No quotes or escape chars", "value", "value"},
		{"Quotes around", "\"value\"", "value"},
		{"Escaped newline", "\\n", "\n"},
		{"Escaped tab", "\\t", "\t"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cleanValue(tt.input); got != tt.expected {
				t.Errorf("cleanValue() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func setEnvVars(vars map[string]string) func() {
	originalVars := make(map[string]string)
	for k, v := range vars {
		originalVars[k] = os.Getenv(k)
		_ = os.Setenv(k, v)
	}

	return func() {
		for k, v := range originalVars {
			_ = os.Setenv(k, v)
		}
	}
}

func TestGetAllEnvVarsFromHost(t *testing.T) {
	// Setting up environment variables for testing
	cleanup := setEnvVars(map[string]string{"TEST_VAR": "\"value\"", "TEST_NEWLINE": "\\n"})
	defer cleanup()

	got := GetAllEnvVarsFromHost()
	if got["TEST_VAR"] != "value" || got["TEST_NEWLINE"] != "\n" {
		t.Errorf("GetAllEnvVarsFromHost() did not clean values correctly, got: %v", got)
	}
}

func TestGetAllEnvVarsFromHostAsStruct(t *testing.T) {
	// Setting up environment variables for testing
	cleanup := setEnvVars(map[string]string{"TEST_VAR": "value", "TEST_NEWLINE": "\n"})
	defer cleanup()

	envVarNames := []string{"TEST_VAR", "TEST_NEWLINE"}
	got := GetAllEnvVarsFromHostAsStruct(envVarNames)
	want := []EnvVar{
		{Name: "TEST_VAR", Value: "value"},
		{Name: "TEST_NEWLINE", Value: "\n"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("GetAllEnvVarsFromHostAsStruct() = %v, want %v", got, want)
	}
}
