package storage

import (
	"fmt"
	"testing"
)

func TestHmacSha256(t *testing.T) {
	tests := []struct {
		name    string
		dataStr string
		keyStr  string
		want    string
	}{
		{
			name:    "Positive test. Test case1 (RFC 4231)",
			dataStr: "4869205468657265",
			keyStr:  "0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b0b",
			want:    "b0344c61d8db38535ca8afceaf0bf12b881dc200c9833da726e9376c2e32cff7",
		},
		{
			name:    "Positive test. Test case2 (RFC 4231)",
			dataStr: "7768617420646f2079612077616e7420666f72206e6f7468696e673f",
			keyStr:  "4a656665",
			want:    "5bdcc146bf60754e6a042426089575c75a003f089d2739839dec58b964ec3843",
		},
		{
			name:    "Positive test. Test case3 (RFC 4231)",
			dataStr: "dddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddddd",
			keyStr:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:    "773ea91e36800e46854db8ebd09181a72959098b3ef8c122d9635514ced565fe",
		},
		{
			name:    "Positive test. Test case7 (RFC 4231)",
			dataStr: "5468697320697320612074657374207573696e672061206c6172676572207468616e20626c6f636b2d73697a65206b657920616e642061206c6172676572207468616e20626c6f636b2d73697a6520646174612e20546865206b6579206e6565647320746f20626520686173686564206265666f7265206265696e6720757365642062792074686520484d414320616c676f726974686d2e",
			keyStr:  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			want:    "9b09ffa71b942fcb27635fbcd5b0e944bfdc63644f0713938a7f51535c3a35e2",
		},
		{
			name:    "Negative test. Wrong key",
			dataStr: "12345678901234567890",
			keyStr:  "abcdefgh",
			want:    "",
		},
		{
			name:    "Negative test. Wrong data",
			dataStr: "abcdefgh",
			keyStr:  "12345678901234567890",
			want:    "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := HmacSha256(tt.dataStr, tt.keyStr)
			if output != tt.want {
				t.Errorf("TEST_ERROR: Hash function output is %s, but expected result is %s.\n", output, tt.want)
			}
		})
	}
}

func TestSetMetric(t *testing.T) {
	var ms = MetricsStorage{
		Metric{ID: "TestMetric1", MType: "gauge", Value: PointOf(123.321)},
		Metric{ID: "TestMetric2", MType: "counter", Delta: PointOf(int64(123))},
		Metric{ID: "TestMetric3", MType: "gauge", Value: PointOf(234.567)},
	}

	tests := []struct {
		name  string
		input Metric
		want  Metric
	}{
		{
			name: "Set gauge",
			input: Metric{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
			},
			want: Metric{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
				Hash:  "a47a6ad075b438dece782009a19c45b18050dcd9e9d85aafc02ee2cfb58b9e0d",
			},
		},
		{
			name: "Set counter",
			input: Metric{
				ID:    "TestMetric2",
				MType: "counter",
				Delta: PointOf(int64(123)),
			},
			want: Metric{
				ID:    "TestMetric2",
				MType: "counter",
				Delta: PointOf(int64(246)),
				Hash:  "1ecc232bd6cb1b65061ae636abe23e936d2e984b963603afb204cced3611f81a",
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms.SetMetric(tt.input)
			if !IsMetricsEqual(ms[i], tt.want) {
				t.Errorf("TEST_ERROR: Current metric is %v", ms[i])
			}
		})
	}
}

func TestGetMetric(t *testing.T) {
	var ms = MetricsStorage{
		Metric{ID: "TestMetric1", MType: "gauge", Value: PointOf(123.321)},
		Metric{ID: "TestMetric2", MType: "counter", Delta: PointOf(int64(123))},
	}
	fmt.Printf("DEBUG: mStorage is %v. \n", ms)

	type inputs struct {
		MetricName string
		MetricType string
	}

	tests := []struct {
		name  string
		input inputs
		want  Metric
	}{
		{
			name: "Get gauge",
			input: inputs{
				MetricName: "TestMetric1",
				MetricType: "gauge",
			},
			want: Metric{
				ID:    "TestMetric1",
				MType: "gauge",
				Value: PointOf(123.321),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			curMetric := ms.GetMetric(tt.input.MetricName, tt.input.MetricType)
			if !IsMetricsEqual(curMetric, tt.want) {
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
