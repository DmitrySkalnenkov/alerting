package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	//"strconv"
	"time"

	"github.com/DmitrySkalnenkov/alerting/internal/auxiliary"
	"github.com/DmitrySkalnenkov/alerting/internal/handlers"
	"github.com/DmitrySkalnenkov/alerting/internal/storage"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	//time.Sleep(500 * time.Millisecond)
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	//hni := func(w http.ResponseWriter, r *http.Request) {
	//	http.Error(w, http.StatusText(http.StatusNotImplemented), http.StatusNotImplemented)
	//	_, err := io.WriteString(w, "Hello from not implemented handler.\n")
	//	if err != nil {
	//		log.Fatal(err)
	//	}
	//}
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

	hostportStr := auxiliary.GetEnvVariable("ADDRESS", "localhost:8080")
	hostportStr = auxiliary.TrimQuotes(hostportStr)
	storeIntervalStr := auxiliary.GetEnvVariable("STORE_INTERVAL", "300s")
	//storeIntervalStr += "s"
	storeFilePath := ""

	valueStoreFilePath, isStoreFilePath := os.LookupEnv("STORE_FILE")
	if !(isStoreFilePath && valueStoreFilePath == "") {
		storeFilePath = auxiliary.GetEnvVariable("STORE_FILE", "/tmp/devops-metrics-db.json")
	}
	isRestoreStr := auxiliary.GetEnvVariable("RESTORE", "true")

	storeIntervalTime := 0 * time.Second
	if storeIntervalStr != "0" {
		var err error
		storeIntervalTime, err = time.ParseDuration(storeIntervalStr)
		if err != nil {
			fmt.Printf("ERROR: Cannot conver STORE_INTERVAL value (%s) to second. Will be used 0 value. \n",
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
