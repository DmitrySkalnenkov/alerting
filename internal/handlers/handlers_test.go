package handlers

import (
	"alerting/internal/storage"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

// func GaugesHandlerAPI2(w http.ResponseWriter, r *http.Request) {
func TestGaugesHandlerAPI2(t *testing.T) {
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
			h := http.HandlerFunc(GaugesHandlerAPI2)
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

// func GaugesHandler(w http.ResponseWriter, r *http.Request) {
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

// func CountersHandler(w http.ResponseWriter, r *http.Request) {
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

// func GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
func TestGetGaugeHandler(t *testing.T) {
	type want struct {
		contentType string
		code        int
		response    string
	}
	tests := []struct {
		name       string
		testMetric storage.Metrics
		getRequest string
		want       want
	}{ //Test table
		{
			name:       "positive test #1",
			testMetric: storage.Metrics{ID: "TestMetric1", MType: "gauge", Value: storage.PointOf(123.0)},
			getRequest: "http://127.0.0.1:8080/value/gauge/TestMetric1",
			want: want{
				contentType: "application/json",
				code:        200,
				response:    `{"id":"TestMetric1","type":"gauge","value":123}`,
			},
		},
		{
			name:       "positive test #2",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "gauge", Value: storage.PointOf(-321.0)},
			getRequest: "http://127.0.0.1:8080/value/gauge/TestMetric2",
			want: want{
				contentType: "application/json",
				code:        200,
				response:    `{"id":"TestMetric2","type":"gauge","value":-321}`,
			},
		},
		{
			name:       "metric not found test",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "gauge", Value: storage.PointOf(1111.0)},
			getRequest: "http://127.0.0.1:8080/value/gauge/TestMetric3",
			want: want{
				contentType: "text/plain",
				code:        404,
				response:    http.StatusText(http.StatusNotFound) + "\n",
			},
		},
	}
	for _, tt := range tests {
		// запускаем каждый тест
		t.Run(tt.name, func(t *testing.T) {
			storage.MetStorage.SetMetric(tt.testMetric)
			getReq := httptest.NewRequest(http.MethodPost, tt.getRequest, nil)
			w := httptest.NewRecorder()
			router := chi.NewRouter()
			router.Use(middleware.Logger)
			router.HandleFunc("/value/gauge/{MetricName}", GetGaugeHandler)
			router.ServeHTTP(w, getReq)
			//GetGaugeHandler(w, getReq)
			resp := w.Result()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header.Get("Content-Type"))
			fmt.Println(string(bodyBytes))
			resp.Body.Close()
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if resp.StatusCode == tt.want.code {
				bodyString := string(bodyBytes)
				if bodyString != tt.want.response {
					t.Errorf("TEST_ERROR: Expected response %v, got %v", tt.want.response, bodyString)
				}
			} else {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, resp.StatusCode)
			}
		})
	}
}
