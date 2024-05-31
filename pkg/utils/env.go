package utils

import (
	"os"
	"strings"
)

// cleanValue cleans and unescapes the provided string value.
//
// Parameters:
//   - value: The string value to clean and unescape.
//
// Returns:
//   - string: The cleaned and unescaped string.
//
// Example:
//
//	cleanedValue := cleanValue("\"example\\nvalue\"")
//	fmt.Printf("Cleaned value: %s\n", cleanedValue)
func cleanValue(value string) string {
	// Remove leading and trailing double quotes
	cleanedValue := strings.Trim(value, "\"")
	// Unescape characters. This example handles common escaped characters.
	// Extend this based on specific needs.
	cleanedValue = strings.ReplaceAll(cleanedValue, "\\n", "\n")
	cleanedValue = strings.ReplaceAll(cleanedValue, "\\t", "\t")
	return cleanedValue
}

// GetAllEnvVarsFromHost retrieves all environment variables from the host and returns them as a map.
//
// Returns:
//   - map[string]string: A map containing all environment variables from the host.
//
// Example:
//
//	envVars := GetAllEnvVarsFromHost()
//	fmt.Printf("Host environment variables: %v\n", envVars)
func GetAllEnvVarsFromHost() map[string]string {
	envVars := make(map[string]string)

	for _, envVar := range os.Environ() {
		pair := strings.SplitN(envVar, "=", 2) // Ensure only the first '=' is used for splitting.
		if len(pair) == 2 {
			envVars[pair[0]] = cleanValue(pair[1])
		}
	}

	return envVars
}

// EnvVar represents an environment variable with a name and value.
type EnvVar struct {
	Name  string
	Value string
}

// GetAllEnvVarsFromHostAsStruct retrieves the specified environment variables from the host and returns them as a slice of EnvVar structs.
//
// Parameters:
//   - envVarNames: A list of environment variable names to retrieve.
//
// Returns:
//   - []EnvVar: A slice of EnvVar structs representing the specified environment variables.
//
// Example:
//
//	envVarNames := []string{"PATH", "HOME"}
//	envVars := GetAllEnvVarsFromHostAsStruct(envVarNames)
//	fmt.Printf("Host environment variables as structs: %v\n", envVars)
func GetAllEnvVarsFromHostAsStruct(envVarNames []string) []EnvVar {
	var envVars []EnvVar
	envs := make(map[string]bool)

	for _, name := range envVarNames {
		envs[name] = true
	}

	for _, envVar := range os.Environ() {
		pair := strings.SplitN(envVar, "=", 2) // Ensure only the first '=' is used for splitting.
		if len(pair) == 2 && envs[pair[0]] {
			envVars = append(envVars, EnvVar{
				Name:  pair[0],
				Value: cleanValue(pair[1]),
			})
		}
	}

	return envVars
}
