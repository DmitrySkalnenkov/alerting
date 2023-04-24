package handlers

import (
	//"alerting/internal"
	"alerting/internal/storage"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

//var MS = storage.NewMetricStorage()

func GaugesHandlerAPI2(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	fmt.Printf("DEBUG: Gauge handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/gauge/[A-Za-z0-9]+/[0-9.-]+$`, urlPath)
	if matched && (err == nil) {
		var curMetric storage.Metrics
		pathSlice := strings.Split(urlPath, "/")
		curMetric.ID = string(pathSlice[3])
		curMetric.MType = "gauge"
		var v float64
		v, err = strconv.ParseFloat(pathSlice[4], 64)
		curMetric.Value = &v
		fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", curMetric.ID, *curMetric.Value)
		if err == nil {
			storage.MetStorage.SetMetric(curMetric)
			//storage.Mstorage.PushGauge(mName, mValue)
			fmt.Printf("DEBUG: Mstorage value for gauge metric %v is %v.\n", curMetric.ID, *storage.MetStorage.GetMetric(curMetric.ID, "gauge").Value)
			//io.WriteString(w, "DEBUG: Hello from gauge handler (Status OK).\n")
		} else {
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

func CounterHandlerAPI2(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	fmt.Printf("DEBUG: Counter handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/counter/[A-Za-z0-9]+/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		var curMetric storage.Metrics
		pathSlice := strings.Split(urlPath, "/")
		curMetric.ID = string(pathSlice[3])
		curMetric.MType = "counter"
		var v int64
		v, err = strconv.ParseInt(pathSlice[4], 10, 64)
		curMetric.Delta = &v
		fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", curMetric.ID, *curMetric.Delta)
		if err == nil {
			storage.MetStorage.SetMetric(curMetric)
			fmt.Printf("DEBUG: Mstorage value for counter metric %v is %v.\n", curMetric.ID, *storage.MetStorage.GetMetric(curMetric.ID, "counter").Delta)
		} else {
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
		//fmt.Printf("INFO: URL is : %s.\n", r.URL.Path)
		//io.WriteString(w, "DEBUG: Hello from counter handler (Bad Request). \n")
	}
}

func GetCounterHandlerAPI2(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	matched, err := regexp.MatchString(`/value/counter/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		curMetric := storage.MetStorage.GetMetric(curMetricName, "counter")
		if curMetric != storage.NilMetric {
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			w.Header().Set("Content-Type", "application/json")
			txJson, err := json.Marshal(curMetric)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			_, err = io.WriteString(w, fmt.Sprintf("%v", string(txJson)))
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("DEBUG: Value of JSON response is %v:", string(txJson))
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

// Handler for getting all current metric values API2
func GetAllMetricsHandlerAPI2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	txJsonMetricList, err := json.Marshal(*storage.MetStorage)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
	}
	_, err = io.WriteString(w, fmt.Sprintf("%v", string(txJsonMetricList)))
	if err != nil {
		log.Fatal(err)
	}
}

//Handler for updating gauge value
// /update/gauges/<MetricName>/<MetricValue> then status -- OK (200) and save MetricValue into map with key MetricName
// /update/gauges/ then status -- NotFound (404)
// /update/gauges then status -- BadRequest (400)

func GaugesHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: Gauge handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/gauge/[A-Za-z0-9]+/[0-9.-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseFloat(pathSlice[4], 64)
		//fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", mName, mValue)
		if err == nil {
			storage.Mstorage.PushGauge(mName, mValue)
			//fmt.Printf("DEBUG: Mstorage gauges is %v.\n", storage.Mstorage.Gauges)
			//io.WriteString(w, "DEBUG: Hello from gauge handler (Status OK).\n")
		} else {
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

// Handler for updating counter value
func CountersHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: Counter handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`/update/counter/[A-Za-z0-9]+/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseInt(pathSlice[4], 10, 64)
		//fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", mName, mValue)
		if err == nil {
			storage.Mstorage.PushCounter(mName, mValue)
			//fmt.Printf("DEBUG: Mstorage counter is %v.\n", storage.Mstorage.Counters)
			_, err = io.WriteString(w, "Hello from counter handler (Status OK).\n")
			if err != nil {
				log.Fatal(err)
			}
		} else {
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
func GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	matched, err := regexp.MatchString(`/value/gauge/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Gauges)
		curMetric := storage.MetStorage.GetMetric(curMetricName, "gauge")
		if curMetric != storage.NilMetric {
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			w.Header().Set("Content-Type", "application/json")
			txJson, err := json.Marshal(curMetric)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			}
			w.WriteHeader(http.StatusOK)
			_, err = io.WriteString(w, fmt.Sprintf("%v", string(txJson)))
			if err != nil {
				log.Fatal()
			}
			//fmt.Println("DEBUG: Value of JSON response is %v:", string(txJson))
		} else {
			w.Header().Set("Content-Type", "application/json")
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

//// Handler for getting counter value
//func GetCounterHandler(w http.ResponseWriter, r *http.Request) {
//	urlPath := r.URL.Path
//	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
//	matched, err := regexp.MatchString(`/value/counter/[A-Za-z0-9]+`, urlPath)
//	if matched && (err == nil) {
//		curMetricName := chi.URLParam(r, "MetricName")
//		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Gauges)
//		curMetric := storage.MetStorage.GetMetric(curMetricName, "counter")
//		if curMetric != storage.NilMetric {
//			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
//			w.Header().Set("Content-Type", "application/json")
//			txJson, err := json.Marshal(curMetric)
//			if err != nil {
//				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
//			}
//			w.WriteHeader(http.StatusOK)
//			io.WriteString(w, fmt.Sprintf("%v", string(txJson)))
//			//fmt.Println("DEBUG: Value of JSON response is %v:", string(txJson))
//		} else {
//			w.Header().Set("Content-Type", "application/json")
//			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
//		}
//	}
//}

//// Handler for getting all current metric values
//func GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
//	//	urlPath := r.URL.Path
//	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
//	//fmt.Printf("DEBUG: MemStorage gauges map %v.\n", storage.Mstorage.Gauges)
//	//fmt.Printf("DEBUG: MemStorage counters map is %v.\n", storage.Mstorage.Counters)
//	for mName, mValue := range storage.Mstorage.Gauges {
//		fmt.Printf("%v - %v\n", mName, mValue)
//		io.WriteString(w, fmt.Sprintf("%v - %v\n", mName, mValue))
//	}
//	for mName, mValue := range storage.Mstorage.Counters {
//		fmt.Printf("%v - %v\n", mName, mValue)
//		io.WriteString(w, fmt.Sprintf("%v - %v\n", mName, mValue))
//	}
//}
