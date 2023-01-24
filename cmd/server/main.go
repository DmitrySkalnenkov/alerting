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
	value, ok := m.gauges[metricName]
	if ok == true {
		value = m.gauges[metricName]
		return value
	} else {
		fmt.Println("Metric with name  %s is not found", metricName)
		return 0
	}
}

func (m MemStorage) PushCounter(metricName string, value int64) {
	_, ok := m.counters[metricName]
	if ok {
		m.counters[metricName] = m.counters[metricName] + value
	} else {
		m.counters[metricName] = value
	}
}

func (m MemStorage) PopCounter(metricName string) int64 {
	var value int64
	value = m.counters[metricName]
	return value
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

	hg := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc gauge\n")
		//fmt.Printf("Req: %s", r.URL.Path)
		urlPath := r.URL.Path
		matched, err := regexp.MatchString(`\/update\/gauge\/[A-Za-z]+\/[0-9]+$`, urlPath)
		if (matched == true) && (err == nil) {
			fmt.Println("Match")
			pathSlice := strings.Split(urlPath, "/")
			mName := string(pathSlice[3])
			mValue, err := strconv.ParseFloat(pathSlice[4], 64)
			if contains(MetricNameArray, mName) && (err == nil) {
				mstorage.PushGauge(mName, mValue)
				fmt.Printf("Mstorage gauges is: %s \n", mstorage.gauges)
			}
		} else {
			fmt.Printf("URL is : %s\n", r.URL.Path)
		}

	}

	hc := func(w http.ResponseWriter, _ *http.Request) {
		io.WriteString(w, "Hello from a HandleFunc counter\n")
	}

	http.HandleFunc("/update/gauge/", hg)
	http.HandleFunc("/update/counter/", hc)
	server.ListenAndServe()

}
