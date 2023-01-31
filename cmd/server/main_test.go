package main

import (
	"testing"
)

func TestPushGauge(t *testing.T) {
	mstorage := new(MemStorage)
	mstorage.gauges = make(map[string]float64)
	mstorage.counters = make(map[string]int64)

	type inputs struct {
		MetricName  string
		MetricValue float64
	}

	tests := []struct {
		name  string
		input inputs
		want  float64
	}{ //Test table
		{
			name: "Positve test",
			input: inputs{
				MetricName:  "TestMetric",
				MetricValue: 123.0,
			},
			want: 123.0,
		},
		{
			name: "Zero value",
			input: inputs{
				MetricName:  "TestMetric",
				MetricValue: 0,
			},
			want: 0,
		},
		{
			name: "Negative value",
			input: inputs{
				MetricName:  "TestMetric",
				MetricValue: -321.0,
			},
			want: -321.0,
		},
		{
			name: "Big value",
			input: inputs{
				MetricName:  "TestMetric",
				MetricValue: 31231231221.0,
			},
			want: 31231231221.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mstorage.PushGauge(tt.input.MetricName, tt.input.MetricValue)
			if mstorage.gauges[tt.input.MetricName] != tt.want {
				t.Errorf("MetricName is %s , want is %f", tt.input.MetricName, tt.want)
			}
		})
	}
}
