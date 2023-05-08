package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// func CounterHandlerAPI2(w http.ResponseWriter, r *http.Request) {
func TestCounterHandlerAPI2(t *testing.T) {
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
			name:    "null metric name",
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
			request: "http://127.0.0.1:8080/counter/gauge/Metric/rqwer",
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
			h := http.HandlerFunc(CounterHandlerAPI2)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("TEST_ERROR: %s:", err)
			}
			defer res.Body.Close()
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}

// func GaugesHandlerAPI2(w http.ResponseWriter, r *http.Request) {
func TestGaugeHandlerAPI2(t *testing.T) {
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
			h := http.HandlerFunc(GaugeHandlerAPI2)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("TEST_ERROR: %s:", err)
			}
			defer res.Body.Close()
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}

func TestGetCounterHandlerAPI2(t *testing.T) {
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
			testMetric: storage.Metrics{ID: "TestMetric1", MType: "counter", Delta: storage.PointOf(int64(123))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric1",
			want: want{
				contentType: "application/json",
				code:        200,
				response:    `{"id":"TestMetric1","type":"counter","delta":123}`,
			},
		},
		{
			name:       "positive test #2",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "counter", Delta: storage.PointOf(int64(-321))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric2",
			want: want{
				contentType: "application/json",
				code:        200,
				response:    `{"id":"TestMetric2","type":"counter","delta":-321}`,
			},
		},
		{
			name:       "metric not found test",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "counter", Delta: storage.PointOf(int64(1111))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric3",
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
			router.HandleFunc("/value/counter/{MetricName}", GetCounterHandlerAPI2)
			router.ServeHTTP(w, getReq)
			resp := w.Result()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			//fmt.Println(resp.StatusCode)
			//fmt.Println(resp.Header.Get("Content-Type"))
			//fmt.Println(string(bodyBytes))
			//log.Fatal(resp.Body.Close())
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

// func GaugesHandler(w http.ResponseWriter, r *http.Request) {
func TestGaugeHandlerAPI1(t *testing.T) {
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

			h := http.HandlerFunc(GaugeHandlerAPI1)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Errorf("TEST_ERROR: %s:", err)
			}
			defer res.Body.Close()
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}

// func CountersHandler(w http.ResponseWriter, r *http.Request) {
func TestCounterHandlerAPI1(t *testing.T) {
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
			h := http.HandlerFunc(CounterHandlerAPI1)
			h.ServeHTTP(w, req)
			res := w.Result()
			_, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			defer res.Body.Close()
			//fmt.Printf("TEST_DEBUG: Status is %s, status code is %d, body is %s. \n", res.Status, res.StatusCode, string(resBody))
			if res.StatusCode != tt.want.code {
				t.Errorf("TEST_ERROR: Expected status code %d, got %d", tt.want.code, res.StatusCode)
			}
		})
	}
}

// func GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
func TestGetGaugeHandlerAPI1(t *testing.T) {
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
				contentType: "text/plain",
				code:        200,
				response:    `123`,
			},
		},
		{
			name:       "positive test #2",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "gauge", Value: storage.PointOf(-321.0)},
			getRequest: "http://127.0.0.1:8080/value/gauge/TestMetric2",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `-321`,
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
			router.HandleFunc("/value/gauge/{MetricName}", GetGaugeHandlerAPI1)
			router.ServeHTTP(w, getReq)
			//GetGaugeHandler(w, getReq)
			resp := w.Result()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header.Get("Content-Type"))
			fmt.Println(string(bodyBytes))
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

// func GetCounterHandler(w http.ResponseWriter, r *http.Request) {
func TestGetCounterHandlerAPI1(t *testing.T) {
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
			testMetric: storage.Metrics{ID: "TestMetric1", MType: "counter", Delta: storage.PointOf(int64(123))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric1",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `123`,
			},
		},
		{
			name:       "positive test #2",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "counter", Delta: storage.PointOf(int64(-321))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric2",
			want: want{
				contentType: "text/plain",
				code:        200,
				response:    `-321`,
			},
		},
		{
			name:       "metric not found test",
			testMetric: storage.Metrics{ID: "TestMetric2", MType: "counter", Delta: storage.PointOf(int64(1111))},
			getRequest: "http://127.0.0.1:8080/value/counter/TestMetric3",
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
			router.HandleFunc("/value/counter/{MetricName}", GetCounterHandlerAPI1)
			router.ServeHTTP(w, getReq)
			//GetGaugeHandler(w, getReq)
			resp := w.Result()
			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()
			fmt.Println(resp.StatusCode)
			fmt.Println(resp.Header.Get("Content-Type"))
			fmt.Println(string(bodyBytes))
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
