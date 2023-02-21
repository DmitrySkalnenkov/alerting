package internal

import (
	"testing"
)

func TestContains(t *testing.T) {
	//func contains(s []string, str string) bool {}
	type inputs struct {
		stringArray []string
		stringValue string
	}

	tests := []struct {
		name  string
		input inputs
		want  bool
	}{ //Test table
		{
			name: "Positve test",
			input: inputs{
				stringArray: []string{
					"Alloc",
					"BuckHashSys",
					"Frees"},
				stringValue: "BuckHashSys",
			},
			want: true,
		},
		{
			name: "Negative test",
			input: inputs{
				stringArray: []string{
					"Alloc",
					"BuckHashSys",
					"Frees"},
				stringValue: "UnknownMetric",
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Contains(tt.input.stringArray, tt.input.stringValue) != tt.want {
				t.Errorf("TEST_ERROR: StringArray is %s , StringValue is %s, want is %t", tt.input.stringArray, tt.input.stringValue, tt.want)
			}
		})
	}
}
