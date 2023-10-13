package auxiliary

import (
	"flag"
	"fmt"
	"os"
)

// GetParameterValue gets parameter value from flag or environment variable (environment more prior)
func GetParamValue(envVarName string, flagName, defaultValue string, usageStr string) string {
	var paramValue string
	_, isEnvVarExists := os.LookupEnv(envVarName)
	if isEnvVarExists {
		paramValue = GetEnvVariable(envVarName, defaultValue)
		return paramValue
	} else {
		flag.StringVar(&paramValue, flagName, defaultValue, usageStr)
		flag.Parse()
		if isFlagPassed(flagName) {
			fmt.Printf("DEBUG: Flag with name '%s' has value '%s'.\n", flagName, paramValue)
			return paramValue
		} else {
			fmt.Printf("DEBUG: There is no value set  for flag '%s', "+
				"so will be used default value '%s'.\n", flagName, defaultValue)
			return paramValue
		}

	}
}

func GetEnvVariable(envVarName string, defaultValue string) string {
	envVarValue := defaultValue
	if os.Getenv(envVarName) != "" {
		envVarValue = os.Getenv(envVarName)
	}
	fmt.Printf("DEBUG: Environment variable '%s' has value '%s'.\n", envVarName, envVarValue)
	return envVarValue
}

// isFlagPassed checks if flag with name sets in command line
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}

// TrimQuotes deletes first and last quotes marks from string if exists
func TrimQuotes(s string) string {
	if len(s) >= 2 {
		if s[0] == '"' && s[len(s)-1] == '"' {
			return s[1 : len(s)-1]
		} else if s[0] == '"' {
			return s[1:len(s)]
		} else if s[len(s)-1] == '"' {
			return s[0 : len(s)-1]
		}
	}
	return s
}
