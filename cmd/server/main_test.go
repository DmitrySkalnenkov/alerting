package main

import (
	//	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPushGauge(t *testing.T) {
	//mstorage := new(MemStorage)
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
			mstorage.PushGauge(tt.input.MetricName, tt.input.MetricValue)
			if mstorage.gauges[tt.input.MetricName] != tt.want {
				t.Errorf("TEST_ERROR: MetricName is %s , want is %f", tt.input.MetricName, tt.want)
			}
		})
	}
}

func TestPushCounter(t *testing.T) {
	//mstorage := new(MemStorage)
	mstorage.gauges = make(map[string]float64)
	mstorage.counters = make(map[string]int64)

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
			mstorage.PushCounter(tt.input.MetricName, tt.input.MetricValue)
			if mstorage.counters[tt.input.MetricName] != tt.want {
				t.Errorf("TEST_ERROR: MetricName is %s , actual is %d, want is %d", tt.input.MetricName, mstorage.counters[tt.input.MetricName], tt.want)
			}
		})
	}
}

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
			if contains(tt.input.stringArray, tt.input.stringValue) != tt.want {
				t.Errorf("TEST_ERROR: StringArray is %s , StringValue is %s, want is %t", tt.input.stringArray, tt.input.stringValue, tt.want)
			}
		})
	}
}

//func GaugesHandler(w http.ResponseWriter, r *http.Request) {
func TestGaugesHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
		response    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{ //Test table
		{
			name:    "positive test #1",
			request: "http://127.0.0.1:8080/update/gauge/TestMetric/12421234123.0",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "positive test #2",
			request: "http://127.0.0.1:8080/update/gauge/TestMetric/-232131123.0",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "null metric name ",
			request: "http://127.0.0.1:8080/update/gauge/",
			want: want{
				contentType: "text/plain",
				code:        404,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "bad request",
			request: "http://127.0.0.1:8080/update/gauge",
			want: want{
				contentType: "text/plain",
				code:        400,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "symbolic value",
			request: "http://127.0.0.1:8080/update/gauge/Metric/rqwer",
			want: want{
				contentType: "text/plain",
				code:        400,
				response:    `{"status":"ok"}`,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(GaugesHandler)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Errorf("TEST_ERROR: %s:", err)
			}
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}

//func CountersHandler(w http.ResponseWriter, r *http.Request) {
func TestCountersHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
		response    string
	}
	tests := []struct {
		name    string
		request string
		want    want
	}{ //Test table
		{
			name:    "positive test #1",
			request: "http://127.0.0.1:8080/update/counter/TestMetric/12421234123",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "positive test #2",
			request: "http://127.0.0.1:8080/update/counter/TestMetric/-232131123",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "null metric name ",
			request: "http://127.0.0.1:8080/update/counter/",
			want: want{
				contentType: "text/plain",
				code:        404,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "bad request",
			request: "http://127.0.0.1:8080/update/counter",
			want: want{
				contentType: "text/plain",
				code:        400,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "symbolic value",
			request: "http://127.0.0.1:8080/update/counter/Metric/rqwer",
			want: want{
				contentType: "text/plain",
				code:        400,
				response:    `{"status":"ok"}`,
			},
		},
		{
			name:    "float value",
			request: "http://127.0.0.1:8080/update/counter/Metric/666.3",
			want: want{
				contentType: "text/plain",
				code:        400,
				response:    `{"status":"ok"}`,
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()
			h := http.HandlerFunc(CountersHandler)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				t.Fatal(err)
			}
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}
