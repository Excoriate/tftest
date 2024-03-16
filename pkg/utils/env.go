package utils

import (
	"os"
	"strings"
)

func cleanValue(value string) string {
	// Remove leading and trailing double quotes
	cleanedValue := strings.Trim(value, "\"")
	// Unescape characters. This example handles common escaped characters.
	// Extend this based on specific needs.
	cleanedValue = strings.ReplaceAll(cleanedValue, "\\n", "\n")
	cleanedValue = strings.ReplaceAll(cleanedValue, "\\t", "\t")
	return cleanedValue
}

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

type EnvVar struct {
	Name  string
	Value string
}

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
