package auxiliary

import (
	"fmt"
	"os"
)

func GetEnvVariable(envVarName string, defaultValue string) string {
	envVarValue := defaultValue
	if os.Getenv(envVarName) != "" {
		envVarValue = os.Getenv(envVarName)
	}
	fmt.Printf("DEBUG: Variable '%s' has value '%s'.\n", envVarName, envVarValue)
	return envVarValue
}

// TrimQuotes deletes first and last quotes marks from string if exists
func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		}
	}
	return s
}
