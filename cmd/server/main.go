package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/DmitrySkalnenkov/alerting/internal/auxiliary"
	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Compress(5))
	r.HandleFunc("/", handlers.GetAllMetricsHandler)
	r.Post("/update/", handlers.UpdateHandler)
	r.Post("/value/", handlers.ValueHandler)
	r.Get("/update/gauge/*", handlers.GaugeHandlerAPI1)
	r.Get("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/gauge/*", handlers.GaugeHandlerAPI1)
	r.Post("/update/counter/*", handlers.CounterHandlerAPI1)
	r.Post("/update/*", handlers.NotImplementedHandler)
	r.Post("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Post("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)
	r.Get("/value/gauge/{MetricName}", handlers.GetGaugeHandlerAPI1)
	r.Get("/value/counter/{MetricName}", handlers.GetCounterHandlerAPI1)

	//hostportStr := auxiliary.GetEnvVariable("ADDRESS", "localhost:8080")
	hostportStr := auxiliary.GetParamValue("ADDRESS", "a", "localhost:8080", "ADDRESS should be in 'ip:port' format")
	hostportStr = auxiliary.TrimQuotes(hostportStr)

	//storeIntervalStr := auxiliary.GetEnvVariable("STORE_INTERVAL", "300s")
	//storeIntervalStr += "s"
	storeIntervalStr := auxiliary.GetParamValue("STORE_INTERVAL", "i", "300s", "Store "+
		"interval should be 0 or <second>s format, default is '300s'.")
	storeFilePath := ""

	valueStoreFilePath, isStoreFilePath := os.LookupEnv("STORE_FILE")
	if !(isStoreFilePath && valueStoreFilePath == "") {
		//storeFilePath = auxiliary.GetEnvVariable("STORE_FILE", "/tmp/devops-metrics-db.json")
		storeFilePath = auxiliary.GetParamValue("STORE_FILE", "f", "/tmp/devops-metrics-db.json", "Store file path should be absolute path "+
			"to file. If STORE_FILE variable is empty string than storing functionality will not be used.")
	}
	//isRestoreStr := auxiliary.GetEnvVariable("RESTORE", "true")
	isRestoreStr := auxiliary.GetParamValue("RESTORE", "r", "true", "RESTORE variable or flag 'r' should be 'true' of 'false'.")
	storeIntervalTime := 0 * time.Second
	if storeIntervalStr != "0" {
		var err error
		storeIntervalTime, err = time.ParseDuration(storeIntervalStr)
		if err != nil {
			fmt.Printf("ERROR: Cannot conver STORE_INTERVAL value (%s) to time. Will be used 0 value. \n",
				storeIntervalStr)
			storeIntervalTime = 0 * time.Second
			return
		}
	}

	s := &http.Server{
		Addr:         hostportStr,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	s.Handler = r
	//s.Handler = handlers.GzipHandle(r)

	var c = make(chan storage.MetricsStorage)
	if strings.ToLower(isRestoreStr) == "true" {
		storage.RestoreMetricsFromFile(storeFilePath, storage.MetStorage)
	}
	_ = storeIntervalTime
	if storeFilePath != "" { // Disabling of storing metrics into the file
		go storage.UpdateMetricsInChannel(c)
		go storage.WriteMetricsToFile(storeFilePath, c, storeIntervalTime)
	}
	log.Fatal(s.ListenAndServe())
}
