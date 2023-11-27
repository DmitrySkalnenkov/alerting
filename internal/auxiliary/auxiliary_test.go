package auxiliary

import (
	"fmt"
	"os"
	"testing"
)

func TestGetParamValue(t *testing.T) {
	type inputs struct {
		envVarName         string
		envVarDefaultValue string
		paramFlagName      string
	}

	tests := []struct {
		name            string
		input           inputs
		isSetEnvVarFunc bool
		setEnvValue     string
		setParamValue   string
		want            string
	}{ //Test table
		{
			name: "Positive test. The parameter value is , the environment variable isn't set",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value1",
				paramFlagName:      "t1",
			},
			isSetEnvVarFunc: false,
			setEnvValue:     "set_env_val",
			setParamValue:   "set_param_val",
			want:            "set_param_val",
		},
		{
			name: "Positive test. The parameter value is set, the environment variable is set too",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value2",
				paramFlagName:      "t1",
			},
			isSetEnvVarFunc: true,
			setEnvValue:     "set_env_val",
			setParamValue:   "set_param_val",
			want:            "set_env_val",
		},
		{
			name: "Positive test. The parameter value is set, the environment variable is empty",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value3",
				paramFlagName:      "t1",
			},
			isSetEnvVarFunc: true,
			setEnvValue:     "",
			setParamValue:   "set_param_val",
			want:            "test_default_value3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isSetEnvVarFunc {
				fmt.Printf("TESTDEBUG: Will be used Setenv(%s)\n", tt.input.envVarName)
				os.Setenv(tt.input.envVarName, tt.setEnvValue)
			} else {
				fmt.Printf("TESTDEBUG: Will be used Unsetenv(%s)\n", tt.input.envVarName)
				os.Unsetenv(tt.input.envVarName)
			}
			argStr := "-" + tt.input.paramFlagName + "=" + tt.setParamValue
			os.Args = append(os.Args, argStr)
			if GetParamValue(tt.input.envVarName, tt.input.paramFlagName, tt.input.envVarDefaultValue, "") != tt.want {
				t.Errorf("TEST_ERROR: input value is %v,  want is %s ", tt.input, tt.want)
			}
		})
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
		setValue        string
		want            string
	}{ //Test table
		{
			name: "Positive test",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value1",
			},
			isSetEnvVarFunc: true,
			setValue:        "set_val",
			//f:    giveMeAFunc(true, "TEST_ENV_VAR"),
			want: "set_val",
		},
		{
			name: "Default value test when var value is empty",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value2",
			},
			isSetEnvVarFunc: true,
			setValue:        "",
			want:            "test_default_value2",
		},
		{
			name: "Default value test when not set var",
			input: inputs{
				envVarName:         "TEST_ENV_VAR",
				envVarDefaultValue: "test_default_value2",
			},
			isSetEnvVarFunc: false,
			setValue:        "set_val",
			want:            "test_default_value2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.isSetEnvVarFunc {
				fmt.Printf("TESTDEBUG: Will be used Setenv(%s)\n", tt.input.envVarName)
				os.Setenv(tt.input.envVarName, tt.setValue)
			} else {
				fmt.Printf("TESTDEBUG: Will be used Unsetenv(%s)\n", tt.input.envVarName)
				os.Unsetenv(tt.input.envVarName)
			}
			if GetEnvVariable(tt.input.envVarName, tt.input.envVarDefaultValue) != tt.want {
				t.Errorf("TEST_ERROR: input value is %v,  want is %s ", tt.input, tt.want)
			}
		})
	}

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
