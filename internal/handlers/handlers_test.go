package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
			name:    "null metric name",
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
