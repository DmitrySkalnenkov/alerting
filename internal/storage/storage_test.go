package storage

import (
	"fmt"
	"testing"
)

func TestSetMetric(t *testing.T) {
	var mStorage MetricsStorage

	tests := []struct {
		name  string
		input Metrics
		want  Metrics
	}{
		{
			name: "Set gauge",
			input: Metrics{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
			},
			want: Metrics{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
			},
		},
		{
			name: "Set counter",
			input: Metrics{
				ID:    "TestMetric2",
				MType: "counter",
				Delta: PointOf(int64(123)),
			},
			want: Metrics{
				ID:    "TestMetric2",
				MType: "counter",
				Delta: PointOf(int64(123)),
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mStorage.SetMetric(tt.input)
			if !isMetricsEqual(mStorage[i], tt.want) {
				t.Errorf("TEST_ERROR: Current metric is %v", mStorage[i])
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	var ms MetricsStorage = MetricsStorage{
		Metrics{ID: "TestMetric1", MType: "gauge", Value: PointOf(123.321)},
		Metrics{ID: "TestMetric2", MType: "counter", Delta: PointOf(int64(123))},
	}
	fmt.Printf("DEBUG: mStorage is %v. \n", ms)

	type inputs struct {
		MetricName string
		MetricType string
	}

	tests := []struct {
		name  string
		input inputs
		want  Metrics
	}{
		{
			name: "Get gauge",
			input: inputs{
				MetricName: "TestMetric1",
				MetricType: "gauge",
			},
			want: Metrics{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curMetric := ms.GetMetric(tt.input.MetricName, tt.input.MetricType)
			if !isMetricsEqual(curMetric, tt.want) {
				t.Errorf("TEST_ERROR: Current metric is %v, want is %v ", curMetric, tt.want)
			}
		})
	}
}

func TestPushGauge(t *testing.T) {
	Mstorage = NewMemStorage()

	type inputs struct {
		MetricName  string
		MetricValue float64
	}

	tests := []struct {
		name  string
		input inputs
		want  float64
	}{
		{
			name: "Positive test",
			input: inputs{
				MetricName:  "TestMetric",
				MetricValue: 123,
			},
			want: 123,
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
			Mstorage.PushGauge(tt.input.MetricName, tt.input.MetricValue)
			if Mstorage.Gauges[tt.input.MetricName] != tt.want {
				t.Errorf("TEST_ERROR: MetricName is %s , want is %f", tt.input.MetricName, tt.want)
			}
		})
	}
}

func TestPushCounter(t *testing.T) {
	Mstorage.Gauges = make(map[string]float64)
	Mstorage.Counters = make(map[string]int64)

	type inputs struct {
		MetricName  string
		MetricValue int64
	}

	tests := []struct {
		name  string
		input inputs
		want  int64
	}{ //Test table
		{
			name: "Positve test",
			input: inputs{
				MetricName:  "TestCounterMetric1",
				MetricValue: 12,
			},
			want: 12,
		},
		{
			name: "Zero value",
			input: inputs{
				MetricName:  "TestCounterMetric2",
				MetricValue: 0,
			},
			want: 0,
		},
		{
			name: "Negative value",
			input: inputs{
				MetricName:  "TestCounterMetric3",
				MetricValue: -3,
			},
			want: -3,
		},
		{
			name: "Big value",
			input: inputs{
				MetricName:  "TestCounterMetric4",
				MetricValue: 31231231221,
			},
			want: 31231231221,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Mstorage.PushCounter(tt.input.MetricName, tt.input.MetricValue)
			if Mstorage.Counters[tt.input.MetricName] != tt.want {
				t.Errorf("TEST_ERROR: MetricName is %s , actual is %d, want is %d", tt.input.MetricName, Mstorage.Counters[tt.input.MetricName], tt.want)
			}
		})
	}
}
