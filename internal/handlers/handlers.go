package handlers

import (
	"alerting/internal"
	"alerting/internal/storage"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

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
		if internal.Contains(internal.MetricNameArray, mName) && (err == nil) {
			storage.Mstorage.PushGauge(mName, mValue)
			fmt.Printf("DEBUG: Mstorage gauges is: %v.\n", storage.Mstorage.Gauges)
			io.WriteString(w, "DEBUG: Hello from gauge handler (Status OK). \n")
		}
	} else if urlPath == "/update/gauge/" {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		fmt.Printf("URL is : %s\n", r.URL.Path)
		io.WriteString(w, "DEBUG: Hello from gauge handler (Status Not Found). \n")
	} else {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		fmt.Printf("URL is : %s\n", r.URL.Path)
		io.WriteString(w, "DEBUG: Hello from gauge handler (Bad Request). \n")
	}
	io.WriteString(w, "DEBUG: Hello from gauge handler (Default). \n")

}

//Handler for counters
func CountersHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	matched, err := regexp.MatchString(`\/update\/counter\/[A-Za-z]+\/[0-9-]+$`, urlPath)
	if matched && (err == nil) {
		pathSlice := strings.Split(urlPath, "/")
		mName := string(pathSlice[3])
		mValue, err := strconv.ParseInt(pathSlice[4], 10, 64)
		if internal.Contains(internal.MetricNameArray, mName) && (err == nil) {
			storage.Mstorage.PushCounter(mName, mValue)
			fmt.Printf("Mstorage counter is: %v.\n", storage.Mstorage.Counters)
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

//Handler for getting gauge values

func GetGaugeHandler(w http.ResponseWriter, r *http.Request) {
	urlPath := r.URL.Path
	matched, err := regexp.MatchString(`\/value\/gauge\/[A-Za-z]+`, urlPath)
	if matched && (err == nil) {
		curMetricName := chi.URLParam(r, "MetricName")
		if internal.Contains(internal.MetricNameArray, curMetricName) && (err == nil) {
			curMetricValue := storage.Mstorage.PopGauge(curMetricName)
			fmt.Printf("DEBUG: Value for %s is %s", curMetricName, curMetricValue)
		}
	}
}
