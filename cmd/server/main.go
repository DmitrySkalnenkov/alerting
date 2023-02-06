package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var MetricNameArray = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
	"PollCount",
}

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

var mstorage = new(MemStorage)

func (m MemStorage) PushGauge(metricName string, value float64) {
	m.gauges[metricName] = value
}

/*func (m MemStorage) PopGauge(metricName string) float64 {
	_, ok := m.gauges[metricName]
	if ok {
		return m.gauges[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
}*/

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.counters[metricName]
	if ok {
		mstorage.counters[metricName] = mstorage.counters[metricName] + value
	} else {
		mstorage.counters[metricName] = value
	}
}

/*func (m MemStorage) PopCounter(metricName string) int64 {
	_, ok := m.counters[metricName]
	if ok {
		return m.counters[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
}*/

//Handler for gauges
// /update/gauges/<MetricName>/<MetricValue> then status -- OK (200) and save MetricValue into map with key MetricName
// /update/gauges/ then status -- NotFound (404)
// /update/gauges then status -- BadRequest (400)
func GaugesHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	matched, err := regexp.MatchString(`\/update\/gauge\/[A-Za-z]+\/[0-9.-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseFloat(pathSlice[4], 64)
		if contains(MetricNameArray, mName) && (err == nil) {
			mstorage.PushGauge(mName, mValue)
			fmt.Printf("Mstorage gauges is: %v.\n", mstorage.gauges)
		}
	} else if urlPath == "/update/gauge/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		fmt.Printf("URL is : %s\n", r.URL.Path)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Printf("URL is : %s\n", r.URL.Path)
	}
	io.WriteString(w, "Hello from gauge handler. \n")

}

//Handler for counters
func CountersHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	matched, err := regexp.MatchString(`\/update\/counter\/[A-Za-z]+\/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseInt(pathSlice[4], 10, 64)
		if contains(MetricNameArray, mName) && (err == nil) {
			mstorage.PushCounter(mName, mValue)
			fmt.Printf("Mstorage counter is: %v.\n", mstorage.counters)
		}
	} else if urlPath == "/update/counter/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		fmt.Printf("URL is : %s\n", r.URL.Path)
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Printf("URL is : %s\n", r.URL.Path)
	}
	io.WriteString(w, "Hello from counter handler.\n")
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}

func main() {

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}
	mstorage.gauges = make(map[string]float64)
	mstorage.counters = make(map[string]int64)

	hg := GaugesHandler
	hc := CountersHandler

	hni := func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
		io.WriteString(w, "Hello from not implemented handler.\n")

	}

	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/", hni)
	http.HandleFunc("/update/gauge/", hg)
	http.HandleFunc("/update/counter/", hc)
	server.ListenAndServe()

}
