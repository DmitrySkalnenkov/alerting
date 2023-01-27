package main

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauges   map[string]float64
	counters map[string]int64
}

func (m MemStorage) PushGauge(metricName string, value float64) {
	m.gauges[metricName] = value
}

func (m MemStorage) PopGauge(metricName string) float64 {
	_, ok := m.gauges[metricName]
	if ok {
		return m.gauges[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
}

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.counters[metricName]
	if ok {
		m.counters[metricName] = m.counters[metricName] + value
	} else {
		//fmt.Printf("Counter metric with name  %s is not found.\n", metricName)
		m.counters[metricName] = value
	}
}

func (m MemStorage) PopCounter(metricName string) int64 {
	_, ok := m.counters[metricName]
	if ok {
		return m.counters[metricName]
	} else {
		fmt.Printf("Gauge metric with name  %s is not found.\n", metricName)
		return 0
	}
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
	MetricNameArray := []string{
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

	server := &http.Server{
		Addr: "127.0.0.1:8080",
	}

	mstorage := new(MemStorage)
	mstorage.gauges = make(map[string]float64)
	mstorage.counters = make(map[string]int64)

	//Handler for gauges
	hg := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from gauge handler. \n")
		w.WriteHeader(http.StatusOK)
		//fmt.Printf("Req: %s", r.URL.Path)
		urlPath := r.URL.Path
		matched, err := regexp.MatchString(`\/update\/gauge\/[A-Za-z]+\/[0-9.]+$`, urlPath)
		if matched && (err == nil) {
			//fmt.Println("Match")
			pathSlice := strings.Split(urlPath, "/")
			mName := string(pathSlice[3])
			mValue, err := strconv.ParseFloat(pathSlice[4], 64)
			if contains(MetricNameArray, mName) && (err == nil) {
				mstorage.PushGauge(mName, mValue)
				fmt.Printf("Mstorage gauges is: %v.\n", mstorage.gauges)
			}
		} else {
			fmt.Printf("URL is : %s\n", r.URL.Path)
			//http.Error(w, err.Error(), 404)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}

	}
	//Handler for counters
	hc := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from counter handler.\n")
		w.WriteHeader(http.StatusOK)
		urlPath := r.URL.Path
		matched, err := regexp.MatchString(`\/update\/counter\/[A-Za-z]+\/[0-9]+$`, urlPath)
		if matched && (err == nil) {
			//fmt.Println("Match")
			pathSlice := strings.Split(urlPath, "/")
			mName := string(pathSlice[3])
			mValue, err := strconv.ParseInt(pathSlice[4], 10, 64)
			if contains(MetricNameArray, mName) && (err == nil) {
				mstorage.PushCounter(mName, mValue)
				fmt.Printf("Mstorage counter is: %v.\n", mstorage.counters)
			}
		} else {
			fmt.Printf("URL is : %s\n", r.URL.Path)
			//http.Error(w, err.Error(), 404)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		}
	}

	/* //Root handler
	hr := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from root handler.\n")
		urlPath := r.URL.Path
		fmt.Printf("DEBUG: URL is : %s\n", urlPath)
		//http.Error(w, "Unkonwn URL", 404)
		http.Error(w, "Error. 404", 404)
	}*/

	http.Handle("/", http.NotFoundHandler())
	http.HandleFunc("/update/gauge/", hg)
	http.HandleFunc("/update/counter/", hc)
	server.ListenAndServe()

}
