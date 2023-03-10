package handlers

import (
	//"alerting/internal"
	"alerting/internal/storage"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

//Handler for updating gauge value
// /update/gauges/<MetricName>/<MetricValue> then status -- OK (200) and save MetricValue into map with key MetricName
// /update/gauges/ then status -- NotFound (404)
// /update/gauges then status -- BadRequest (400)

func GaugesHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: Gauge handler. URL is %s.\n", string(urlPath))
	matched, err := regexp.MatchString(`\/update\/gauge\/[A-Za-z0-9]+\/[0-9.-]+$`, urlPath)
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
			io.WriteString(w, fmt.Sprintf("Value parsing error. %s.\n", err))
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
	matched, err := regexp.MatchString(`\/update\/counter\/[A-Za-z0-9]+\/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseInt(pathSlice[4], 10, 64)
		//fmt.Printf("DEBUG: Metric name matched. MetricName is %s, MetricValue is %v.\n", mName, mValue)
		if err == nil {
			storage.Mstorage.PushCounter(mName, mValue)
			//fmt.Printf("DEBUG: Mstorage counter is %v.\n", storage.Mstorage.Counters)
			io.WriteString(w, "Hello from counter handler (Status OK).\n")
		} else {
			io.WriteString(w, fmt.Sprintf("Value parsing error. %s.\n", err))
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
	matched, err := regexp.MatchString(`\/value\/gauge\/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Gauges)
		_, ok := storage.Mstorage.Gauges[curMetricName]
		if ok {
			curMetricValue := storage.Mstorage.PopGauge(curMetricName)
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			io.WriteString(w, fmt.Sprintf("%v", curMetricValue))
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

// Handler for getting counter value
func GetCounterHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	//fmt.Printf("DEBUG: URL is : %s.\n", urlPath)
	matched, err := regexp.MatchString(`\/value\/counter\/[A-Za-z0-9]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		//fmt.Printf("DEBUG: MemStorage map is %v.\n", storage.Mstorage.Counters)
		_, ok := storage.Mstorage.Counters[curMetricName]
		if ok {
			curMetricValue := storage.Mstorage.PopCounter(curMetricName)
			//fmt.Printf("DEBUG: Value for %s is %v.\n", curMetricName, curMetricValue)
			io.WriteString(w, fmt.Sprintf("%v", curMetricValue))
		} else {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		}
	}
}

//Handler for getting all current metric values
func GetAllMetricsHandler(w http.ResponseWriter, r *http.Request) {
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
