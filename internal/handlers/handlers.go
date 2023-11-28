package handlers

import (
	"compress/gzip"
	"encoding/json"
	"fmt"

	"github.com/go-chi/chi/v5"

	//"alerting/internal"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/DmitrySkalnenkov/alerting/internal/storage"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

// GzipHandle as middleware
func GzipHandle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// проверяем, что клиент поддерживает gzip-сжатие
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			// если gzip не поддерживается, передаём управление
			// дальше без изменений
			next.ServeHTTP(w, r)
			return
		}

		// создаём gzip.Writer поверх текущего w
		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		// передаём обработчику страницы переменную типа gzipWriter для вывода данных
		next.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}

// For Not implemented handlers
func NotImplementedHandler(w http.ResponseWriter, r *http.Request) {
	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	_, err := io.WriteString(w, "Hello from not implemented handler.\n")
	if err != nil {
		log.Fatal(err)
	}
}

//TODO: Delete Get from UpdateHandler

// handler for URL /update/ (GET or POST)
func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	urlSliced := strings.Split(r.URL.Path, "/")
	if r.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(r.Body)
		var curMetric storage.Metrics
		err := decoder.Decode(&curMetric)
		if err != nil {
			log.Println(err)
			return
		}
		if (curMetric.MType == "gauge" || curMetric.MType == "counter") && curMetric.ID != "" {
			w.WriteHeader(http.StatusOK)
			fmt.Printf("DEBUG: Metric %v was stored into the storage.\n", curMetric)
			storage.MetStorage.SetMetric(curMetric)

		} else {
			NotImplementedHandler(w, r)
		}
	} else if urlSliced[2] == "gauge" {
		GaugeHandlerAPI1(w, r)
	} else if urlSliced[2] == "counter" {
		CounterHandlerAPI1(w, r)
	} else {
		NotImplementedHandler(w, r)
	}
}

// handeler for URL /value/ (GET or POST)
func ValueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") == "application/json" {
		decoder := json.NewDecoder(r.Body)
		var curMetric storage.Metrics
		err := decoder.Decode(&curMetric)
		if err != nil {
			log.Println(err)
			return
		}
		if (curMetric.MType == "gauge" || curMetric.MType == "counter") && curMetric.ID != "" {
			//fmt.Printf("DEBUG: Get metric struct from request %v.\n", curMetric)
			m := storage.MetStorage.GetMetric(curMetric.ID, curMetric.MType)
			if !storage.IsMetricsEqual(m, storage.NilMetric) {
				switch curMetric.MType {
				case "gauge":
					curMetric.Value = m.Value
				case "counter":
					curMetric.Delta = m.Delta
				}
				w.Header().Set("Content-Type", "application/json")
				txJSON, err := json.Marshal(curMetric)
				if err != nil {
					http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				}
				w.WriteHeader(http.StatusOK)
				_, err = io.WriteString(w, fmt.Sprintf("%v", string(txJSON)))
				if err != nil {
					log.Fatal(err)
				}
			} else {
				//w.Header().Set("Content-Type", "plain/json")
				w.WriteHeader(http.StatusNotFound)
				_, err = io.WriteString(w, fmt.Sprintf("Metric with ID %v and type %v is not found.", curMetric.ID, curMetric.MType))
				if err != nil {
					log.Fatal(err)
				}
			}
		} else {
			NotImplementedHandler(w, r)
		}
	} else {
		NotImplementedHandler(w, r)
	}
}

// Handler for updating gauge value. GET request
// /update/gauges/<MetricName>/<MetricValue> then status -- OK (200) and save MetricValue into map with key MetricName
// /update/gauges/ then status -- NotFound (404)
// /update/gauges then status -- BadRequest (400)
func GaugeHandlerAPI1(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: Gauge handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/gauge/[A-Za-z0-9]+/[0-9.-]+$`, urlPath)
	if matched && (err == nil) {
		curMetric := *storage.NewMetric()
		pathSlice := strings.Split(urlPath, "/")
		curMetric.ID = string(pathSlice[3])
		curMetric.MType = "gauge"
		var v float64
		v, err = strconv.ParseFloat(pathSlice[4], 64)
		*curMetric.Value = v
		fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", curMetric.ID, *curMetric.Value)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			storage.MetStorage.SetMetric(curMetric)
			//storage.Mstorage.PushGauge(mName, mValue)
			//fmt.Printf("DEBUG: Mstorage gauges is %v.\n", storage.Mstorage.Gauges)
			//io.WriteString(w, "DEBUG: Hello from gauge handler (Status OK).\n")
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err = io.WriteString(w, fmt.Sprintf("Value parsing error. %s.\n", err))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if (err == nil) && (urlPath == "/update/gauge/") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		//fmt.Printf("DEBUG: URL is %s.\n", r.URL.Path)
		//io.WriteString(w, "DEBUG: Hello from gauge handler (Status Not Found). \n")
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		//fmt.Printf("INFO: URL is : %s.\n", r.URL.Path)
		//io.WriteString(w, "DEBUG: Hello from gauge handler (Bad Request). \n")
	}
}

// Handler for updating counter value. GET request
func CounterHandlerAPI1(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: Counter handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/counter/[A-Za-z0-9]+/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		curMetric := *storage.NewMetric()
		pathSlice := strings.Split(urlPath, "/")
		curMetric.ID = string(pathSlice[3])
		curMetric.MType = "counter"
		var d int64
		d, err = strconv.ParseInt(pathSlice[4], 10, 64)
		*curMetric.Delta = d
		//fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", curMetric.ID, *curMetric.Delta)
		if err == nil {
			w.WriteHeader(http.StatusOK)
			storage.MetStorage.SetMetric(curMetric)
			//fmt.Printf("DEBUG: Mstorage counter is %v.\n", storage.Mstorage.Counters)
			_, err = io.WriteString(w, "Hello from counter handler (Status OK).\n")
			if err != nil {
				log.Fatal(err)
			}
		} else {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_, err = io.WriteString(w, fmt.Sprintf("Value parsing error. %s.\n", err))
			if err != nil {
				log.Fatal(err)
			}
		}
	} else if (err == nil) && (urlPath == "/update/counter/") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		//fmt.Printf("DEBUG: URL is %s.\n", r.URL.Path)
		//io.WriteString(w, "DEBUG: Hello from counter handler (Status Not Found). \n")
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		//fmt.Printf("DEBUG: URL is : %s.\n", r.URL.Path)
		//io.WriteString(w, "DEBUG: Hello from counter handler (Bad Request). \n")
	}
}

// Handler for getting gauge value
func GetGaugeHandlerAPI1(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	matched, err := regexp.MatchString(`/value/gauge/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Gauges)
		curMetric := storage.MetStorage.GetMetric(curMetricName, "gauge")
		if curMetric != storage.NilMetric {
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			//w.Header().Set("Content-Type", "plain/text")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf("%v", *(curMetric.Value)))
			fmt.Printf("DEBUG: Value of curMetric  is %v:\n", *(curMetric.Value))
		} else {
			//w.Header().Set("Content-Type", "plain/text")
			w.WriteHeader(http.StatusNotFound)
			fmt.Printf("DEBUG: Value of curMetric  is not found")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

// Handler for getting counter value
func GetCounterHandlerAPI1(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	matched, err := regexp.MatchString(`/value/counter/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Gauges)
		curMetric := storage.MetStorage.GetMetric(curMetricName, "counter")
		if curMetric != storage.NilMetric {
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			//w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, fmt.Sprintf("%v", *(curMetric.Delta)))
			fmt.Printf("DEBUG: Value of curMetric  is %v:\n", *(curMetric.Delta))
		} else {
			//w.Header().Set("Content-Type", "text/plain")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

func GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetAllMetricsHandlerAPI1(w, r)
	case "POST":
		GetAllMetricsHandlerAPI2(w, r)
	}
}

// Handler for getting all current metric values
func GetAllMetricsHandlerAPI1(w http.ResponseWriter, r *http.Request) {
	//	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	//fmt.Printf("DEBUG: MemStorage gauges map %v.\n", storage.Mstorage.Gauges)
	//fmt.Printf("DEBUG: MemStorage counters map is %v.\n", storage.Mstorage.Counters)
	for mName, mValue := range storage.Mstorage.Gauges {
		fmt.Printf("%v - %v\n", mName, mValue)
		io.WriteString(w, fmt.Sprintf("%v - %v\n", mName, mValue))
	}
	for mName, mValue := range storage.Mstorage.Counters {
		fmt.Printf("%v - %v\n", mName, mValue)
		io.WriteString(w, fmt.Sprintf("%v - %v\n", mName, mValue))
	}
}

// Handler for getting all current metric values API2
func GetAllMetricsHandlerAPI2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	txJSONMetricList, err := json.Marshal(*storage.MetStorage)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	_, err = io.WriteString(w, fmt.Sprintf("%v", string(txJSONMetricList)))
	if err != nil {
		log.Fatal(err)
	}
}
