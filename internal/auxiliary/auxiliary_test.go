package auxiliary

import (
	"fmt"
	"os"
	"testing"
)

// TODO Make tests for auxiliary funcs
const SetValue = "set_val"

func giveMeAFunc(isSet bool, k string) func(string) error {
	return func(k string) error {
		if isSet {
			fmt.Printf("TESTDEBUG: Will be used Setenv(%s)", k)
			return os.Setenv(k, SetValue)
		} else {
			fmt.Printf("TESTDEBUG: Will be used Unsetenv(%s)", k)
			return os.Unsetenv(k)
		}
	}
}

func TestGetEnvVariable(t *testing.T) {
	type inputs struct {
		envVarName         string
		envVarDefaultValue string
	}

	tests := []struct {
		name            string
		input           inputs
		isSetEnvVarFunc bool
		want            string
	}{ //Test table
		{
			name: "Positive test",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value1",
			},
			isSetEnvVarFunc: true,
			//f:    giveMeAFunc(true, "TEST_ENV_VAR"),
			want: SetValue,
		},
		{
			name: "Default value test when var value is empty",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value2",
			},
			isSetEnvVarFunc: true,
			want:            "test_default_value2",
		},
		{
			name: "Default value test when not set var",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value2",
			},
			isSetEnvVarFunc: false,
			want:            "test_default_value2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if GetEnvVariable(tt.input.envVarName, tt.input.envVarDefaultValue) != tt.want {
				giveMeAFunc(tt.isSetEnvVarFunc, tt.input.envVarName)
				t.Errorf("TEST_ERROR: input value is %v,  want is %s ", tt.input, tt.want)
			}
		})
	}

	t.Setenv("TEST_ENV_VAR", "test_val")
}

func TestTrimQuotes(t *testing.T) {
	tests := []struct {
		name     string
		inputVal string
		wantVal  string
	}{ //Test table
		{
			name:     "Positive test",
			inputVal: "\"test_string1\"",
			wantVal:  "test_string1",
		},
		{
			name:     "First quote mark is missed",
			inputVal: "test_string2\"",
			wantVal:  "test_string2",
		},
		{
			name:     "Last quote mark is missed",
			inputVal: "\"test_string3",
			wantVal:  "test_string3",
		},
		{
			name:     "No quote marks",
			inputVal: "test_string4",
			wantVal:  "test_string4",
		},
	}

	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			if TrimQuotes(tt.inputVal) != tt.wantVal {
				t.Errorf("TEST_ERROR: input value is %s,  want is %s ", tt.inputVal, tt.wantVal)
			}
		})
	}
}
